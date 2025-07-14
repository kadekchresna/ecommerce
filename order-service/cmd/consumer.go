package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/kadekchresna/ecommerce/order-service/config"
	"github.com/kadekchresna/ecommerce/order-service/helper/logger"
	helper_messaging "github.com/kadekchresna/ecommerce/order-service/helper/messaging"
	helper_time "github.com/kadekchresna/ecommerce/order-service/helper/time"
	helper_uuid "github.com/kadekchresna/ecommerce/order-service/helper/uuid"
	driver_db "github.com/kadekchresna/ecommerce/order-service/infrastructure/db"
	"github.com/kadekchresna/ecommerce/order-service/infrastructure/kafka"
	"github.com/kadekchresna/ecommerce/order-service/infrastructure/lock"
	handler "github.com/kadekchresna/ecommerce/order-service/internal/v1/delivery/consumer"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/repository"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/usecase"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

func init() {
	if os.Getenv("APP_ENV") != PRODUCTIOON {
		cobra.OnInitialize(initConfigConsumer)
	}

	rootCmd.AddCommand(versionCmdConsumer)
}

var versionCmdConsumer = &cobra.Command{
	Use:   "consumer",
	Short: "Running Consumer Service",
	Run: func(cmd *cobra.Command, args []string) {
		runConsumer()
	},
}

func initConfigConsumer() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("error load ENV, %s", err.Error()))
	}
}

func runConsumer() {

	config := config.InitConfig()
	consumer := kafka.NewKafkaConsumer(
		[]string{config.KafkaURL},
		helper_messaging.TOPIC_WAREHOUSE_EVENTS,
		helper_messaging.ORDER_CONSUMER_GROUP,
	)

	defer consumer.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: config.RedisURL,
		DB:   0,
	})

	redisLock := lock.NewRedisLock(rdb)

	ctx, cancel := context.WithCancel(context.Background())
	go handleShutdown(cancel)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"),
	)
	db := driver_db.InitDB(dsn)
	timer := helper_time.NewTime(nil)
	uuidHelper := helper_uuid.NewUUID(nil)

	producer := kafka.NewKafkaProducer(
		[]string{config.KafkaURL},
		helper_messaging.TOPIC_ORDER_EVENTS,
	)
	defer producer.Close()

	logger.LogWithContext(context.Background()).Info(fmt.Sprintf("%s consumer service started...", config.AppName))

	handler := handler.NewConsumerHandler(usecase.NewOrdersUsecase(
		timer,
		uuidHelper,
		redisLock,
		repository.NewOrdersRepository(db, timer, uuidHelper, producer),
		repository.NewProductRepository(config.ProductServiceURL, config.AppStaticToken),
	))

	for {
		select {
		case <-ctx.Done():
			logger.LogWithContext(context.Background()).Info(fmt.Sprintf("%s consumer service finished", config.AppName))
			return

		default:
			msg, err := consumer.ReadMessage(ctx)
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}

			handler.HandleMessage(ctx, msg)

			consumer.Commit(ctx, msg)
		}
	}

}

func handleShutdown(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	cancel()
}
