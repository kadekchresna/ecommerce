package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/order-service/helper/jwt"
	helper_messaging "github.com/kadekchresna/ecommerce/order-service/helper/messaging"
	helper_time "github.com/kadekchresna/ecommerce/order-service/helper/time"
	helper_uuid "github.com/kadekchresna/ecommerce/order-service/helper/uuid"
	kafka_helper "github.com/kadekchresna/ecommerce/order-service/infrastructure/kafka"
	"github.com/kadekchresna/ecommerce/order-service/infrastructure/lock"
	handler "github.com/kadekchresna/ecommerce/order-service/internal/v1/delivery/http"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/dto"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/model"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/repository"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/repository/dao"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/repository/interface/mocks"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/usecase"
	"github.com/labstack/echo/v4"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	kafkaAddr string
	redisAddr string
	dbHost    string
	dbUser    string
	dbPass    string

	err error
	db  *gorm.DB

	pool      *dockertest.Pool
	redisLock *lock.RedisLock
	consumer  *kafka_helper.KafkaConsumer
	producer  *kafka_helper.KafkaProducer
)

func createKafkaTopic(brokerAddr, topic string, partitions int) error {
	conn, err := kafka.Dial("tcp", brokerAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	controllerConn, err := kafka.Dial("tcp", controller.Host+":"+strconv.Itoa(controller.Port))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: 1,
		},
	}
	return controllerConn.CreateTopics(topicConfigs...)
}

func waitForKafkaReady(brokerAddr string) error {
	const maxWait = 60 * 2 * time.Second
	const retryInterval = 2 * time.Second

	deadline := time.Now().Add(maxWait)
	for time.Now().Before(deadline) {
		conn, err := kafka.Dial("tcp", brokerAddr)
		if err == nil {

			partitions, err := conn.ReadPartitions("__consumer_offsets")
			conn.Close()
			if err == nil && len(partitions) > 0 {
				break
			}
		}
		time.Sleep(retryInterval)
	}

	for time.Now().Before(deadline) {
		conn, err := kafka.Dial("tcp", brokerAddr)
		if err == nil {

			partitions, err := conn.ReadPartitions(helper_messaging.TOPIC_ORDER_EVENTS)
			conn.Close()
			if err == nil && len(partitions) > 0 {
				return nil
			}
		}
		time.Sleep(retryInterval)
	}
	return fmt.Errorf("Kafka not fully ready")
}

func TestMain(m *testing.M) {

	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	var dbResource, kafkaResource, redisResource *dockertest.Resource

	cleanup := func() {
		if dbResource != nil {
			_ = pool.Purge(dbResource)
		}
		if kafkaResource != nil {
			_ = pool.Purge(kafkaResource)
		}
		if redisResource != nil {
			_ = pool.Purge(redisResource)
		}
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		fmt.Println("Interrupt signal received...")
		cleanup()
		os.Exit(1)
	}()

	dbResource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_DB=omsdb",
		},
	}, func(config *docker.HostConfig) {
		// Allow container to auto-remove
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		config.PortBindings = map[docker.Port][]docker.PortBinding{
			"5432/tcp": {{HostIP: "0.0.0.0", HostPort: "5432"}},
		}
	})
	if err != nil {
		cleanup()
		log.Fatalf("Could not start resource: %s", err)
	}

	dbResource.GetPort("5432/tcp")

	kafkaResource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "bitnami/kafka",
		Tag:        "latest",
		Env: []string{
			"KAFKA_CFG_NODE_ID=1",
			"KAFKA_CFG_PROCESS_ROLES=controller,broker",
			"KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER",
			"KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093",
			"KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092",
			"KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT",
			"KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@localhost:9093",
			"ALLOW_PLAINTEXT_LISTENER=yes",
		},
		ExposedPorts: []string{"9092/tcp"},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		config.PortBindings = map[docker.Port][]docker.PortBinding{
			"9092/tcp": {{HostIP: "0.0.0.0", HostPort: "9092"}},
		}
	})
	if err != nil {
		cleanup()
		log.Fatal("Could not start Kafka container:", err)
	}

	kafkaAddr = fmt.Sprintf("localhost:%s", kafkaResource.GetPort("9092/tcp"))

	createKafkaTopic(kafkaAddr, helper_messaging.TOPIC_ORDER_EVENTS, 1)

	if err := waitForKafkaReady(kafkaAddr); err != nil {
		log.Fatalf("Kafka not ready: %v", err)
	}

	redisResource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "7-alpine",
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		config.PortBindings = map[docker.Port][]docker.PortBinding{
			"6379/tcp": {{HostIP: "0.0.0.0", HostPort: "6379"}},
		}
	})
	if err != nil {
		cleanup()
		log.Fatal("Could not start Redis container:", err)
	}

	redisAddr = fmt.Sprintf("localhost:%s", redisResource.GetPort("6379/tcp"))

	dbHost = "localhost"
	dbUser = "postgres"
	dbPass = "secret"

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?sslmode=disable",
		dbUser,
		dbPass,
		dbHost,
		"omsdb",
	)

	err = pool.Retry(func() error {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:             logger.Default.LogMode(logger.Info),
			PrepareStmt:        true,
			PrepareStmtMaxSize: 10,
		})

		if err != nil {
			return err
		}

		qs := []string{
			`CREATE TYPE public.inbox_status_type AS ENUM (
				'created',
				'in-progress',
				'failed',
				'success'
			);`,
			`CREATE TYPE public.order_status_type AS ENUM (
				'created',
				'in-payment',
				'expired',
				'completed',
				'cancelled',
				'reserving-stock',
				'failed'
			);`,
			`CREATE TYPE public.outbox_status_type AS ENUM (
				'created',
				'in-progress',
				'failed',
				'success'
			);`,
			`CREATE TABLE public.inbox (
				uuid uuid DEFAULT gen_random_uuid() NOT NULL,
				metadata jsonb DEFAULT '{}'::jsonb NOT NULL,
				status public.inbox_status_type DEFAULT 'created'::public.inbox_status_type NOT NULL,
				type character varying DEFAULT ''::character varying NOT NULL,
				created_at timestamp with time zone DEFAULT now() NOT NULL,
				updated_at timestamp with time zone DEFAULT now() NOT NULL,
				reference character varying DEFAULT ''::character varying NOT NULL,
				response jsonb DEFAULT '{}'::jsonb NOT NULL,
				retry_count integer DEFAULT 0 NOT NULL,
				action character varying DEFAULT ''::character varying NOT NULL
			);`,
			`CREATE TABLE public.orders (
				uuid uuid DEFAULT gen_random_uuid() NOT NULL,
				code character varying DEFAULT ''::character varying NOT NULL,
				user_uuid uuid NOT NULL,
				total_amount double precision DEFAULT 0.0 NOT NULL,
				expired_at timestamp with time zone NOT NULL,
				status public.order_status_type DEFAULT 'created'::public.order_status_type NOT NULL,
				created_at timestamp with time zone DEFAULT now() NOT NULL,
				updated_at timestamp with time zone DEFAULT now() NOT NULL,
				created_by uuid NOT NULL,
				updated_by uuid NOT NULL,
				metadata jsonb DEFAULT '{}'::jsonb NOT NULL
			);`,
			`CREATE TABLE public.orders_detail (
				uuid uuid DEFAULT gen_random_uuid() NOT NULL,
				product_uuid uuid NOT NULL,
				product_title character varying DEFAULT ''::character varying NOT NULL,
				product_price double precision DEFAULT 0.0 NOT NULL,
				quantity integer DEFAULT 0 NOT NULL,
				sub_total double precision DEFAULT 0.0 NOT NULL,
				order_uuid uuid NOT NULL
			);`,
			`CREATE TABLE public.outbox (
				uuid uuid DEFAULT gen_random_uuid() NOT NULL,
				metadata jsonb DEFAULT '{}'::jsonb NOT NULL,
				status public.outbox_status_type DEFAULT 'created'::public.outbox_status_type NOT NULL,
				type character varying DEFAULT ''::character varying NOT NULL,
				created_at timestamp with time zone DEFAULT now() NOT NULL,
				updated_at timestamp with time zone DEFAULT now() NOT NULL,
				reference character varying DEFAULT ''::character varying NOT NULL,
				response jsonb DEFAULT '{}'::jsonb NOT NULL,
				retry_count integer DEFAULT 0 NOT NULL,
				action character varying DEFAULT ''::character varying NOT NULL
			);`,
			`
			ALTER TABLE ONLY public.inbox
				ADD CONSTRAINT inbox_pk PRIMARY KEY (uuid);`,
			`
			ALTER TABLE ONLY public.inbox
				ADD CONSTRAINT inbox_unique UNIQUE (type, reference, action);`,
			`
			ALTER TABLE ONLY public.orders
				ADD CONSTRAINT orders_pk PRIMARY KEY (uuid);`,
			`
			ALTER TABLE ONLY public.outbox
				ADD CONSTRAINT outbox_pk PRIMARY KEY (uuid);`,
			`
			ALTER TABLE ONLY public.outbox
				ADD CONSTRAINT outbox_unique UNIQUE (type, reference, action);`,
			`CREATE INDEX inbox_status_idx ON public.inbox USING btree (status);`,
			`CREATE UNIQUE INDEX orders_code_idx ON public.orders USING btree (code);`,
			`CREATE INDEX orders_detail_order_uuid_idx ON public.orders_detail USING btree (order_uuid);`,
			`CREATE INDEX outbox_status_idx ON public.outbox USING btree (status);`,
		}

		for _, q := range qs {
			err := db.Exec(q).Error

			if err != nil {
				log.Fatalf("Failed to execute raw SQL: %v", err)
			}
		}

		return nil
	})
	if err != nil {
		cleanup()
		log.Fatalf("Could not connect to database: %s", err)
	}

	err = pool.Retry(func() error {

		rdb := redis.NewClient(&redis.Options{
			Addr: redisAddr,
			DB:   0,
		})

		redisLock = lock.NewRedisLock(rdb)

		consumer = kafka_helper.NewKafkaConsumer(
			// []string{"localhost:29092"},
			[]string{kafkaAddr},
			helper_messaging.TOPIC_ORDER_EVENTS,
			helper_messaging.ORDER_CONSUMER_GROUP,
		)
		return nil
	})
	if err != nil {
		cleanup()
		log.Fatalf("Could not connect to redis: %s", err)
	}

	code := m.Run()
	cleanup()
	os.Exit(code)

}

func Test_Checkout(t *testing.T) {

	timeHelper := helper_time.NewTime(nil)
	uuidHelper := helper_uuid.NewUUID(uuid.UUID{})

	v1 := echo.New().Group("/api/v1")

	producer = kafka_helper.NewKafkaProducer(
		[]string{kafkaAddr},
		// []string{"localhost:29092"},
		helper_messaging.TOPIC_ORDER_EVENTS,
	)

	productServiceMock := mocks.NewMockIProductRepository(t)
	h := handler.NewOrdersHandler(v1, usecase.NewOrdersUsecase(
		timeHelper,
		uuidHelper,
		redisLock,
		repository.NewOrdersRepository(db, timeHelper, uuidHelper, producer),
		productServiceMock,
	))

	// START RUN CHECKOUT
	userUUID := uuid.New()
	productUUID := uuid.New()

	p := model.Products{
		UUID:  productUUID,
		Title: "title",
		Price: 10000,
	}

	productServiceMock.EXPECT().GetProduct(context.Background(), productUUID).Return(&p, nil).Once()
	body := dto.CheckoutRequest{
		Order: dto.Order{
			OrderDetails: []dto.OrderDetails{
				{
					ProductUUID: productUUID,
					Quantity:    1,
				},
			},
		}}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/checkout", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := echo.New().NewContext(req, rec)

	c.Set(jwt.USER_UUID_KEY, userUUID)

	if err := h.Checkout(c); err != nil {
		t.Errorf("Failed Checkout. %s", err.Error())
	}

	// END RUN CHECKOUT

	// START RUN OUTBOX
	req = httptest.NewRequest(http.MethodGet, "/run-outbox", nil)
	rec = httptest.NewRecorder()

	c = echo.New().NewContext(req, rec)

	if err := h.ProcessOutbox(c); err != nil {
		t.Errorf("Failed ProcessOutbox. %s", err.Error())
	}
	// END RUN OUTBOX

	// START ASSERT

	out := dao.OutboxDAO{}
	if err := db.Model(dao.OutboxDAO{}).First(&out).Error; err != nil {
		t.Errorf("Failed assert outbox. %s", err.Error())
	}

	if out.Status != model.OutboxStatusSuccess {
		t.Error("Failed to send outbox")
	}

	// END ASSERT

}
