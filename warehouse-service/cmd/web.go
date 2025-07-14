package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kadekchresna/ecommerce/warehouse-service/config"

	"github.com/kadekchresna/ecommerce/warehouse-service/helper/logger"
	helper_messaging "github.com/kadekchresna/ecommerce/warehouse-service/helper/messaging"
	helper_time "github.com/kadekchresna/ecommerce/warehouse-service/helper/time"
	helper_uuid "github.com/kadekchresna/ecommerce/warehouse-service/helper/uuid"
	driver_db "github.com/kadekchresna/ecommerce/warehouse-service/infrastructure/db"
	"github.com/kadekchresna/ecommerce/warehouse-service/infrastructure/kafka"
	"github.com/kadekchresna/ecommerce/warehouse-service/infrastructure/lock"
	handler "github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/delivery/http"
	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/repository"
	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

const (
	STAGING     = `stg`
	PRODUCTIOON = `prd`
)

func init() {
	if os.Getenv("APP_ENV") != PRODUCTIOON {

		// init invoke env before everything
		cobra.OnInitialize(initConfigWeb)

	}

	// adding command invokable
	rootCmd.AddCommand(versionCmdWeb)
}

var versionCmdWeb = &cobra.Command{
	Use:   "web",
	Short: "Running Web Service",
	Run: func(cmd *cobra.Command, args []string) {
		runWeb()
	},
}

func initConfigWeb() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("error load ENV, %s", err.Error()))
	}
}

func runWeb() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PATCH, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Use(logger.RequestIDMiddleware)
	e.Use(logger.ClientIPMiddleware)

	config := config.InitConfig()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"),
	)
	db := driver_db.InitDB(dsn)

	rdb := redis.NewClient(&redis.Options{
		Addr: config.RedisURL,
		DB:   0,
	})

	redisLock := lock.NewRedisLock(rdb)

	producer := kafka.NewKafkaProducer(
		[]string{config.KafkaURL},
		helper_messaging.TOPIC_WAREHOUSE_EVENTS,
	)
	defer producer.Close()

	// V1 Endpoints
	v1 := e.Group("/api/v1")
	timer := helper_time.NewTime(nil)
	uuidHelper := helper_uuid.NewUUID(nil)

	handler.NewHealthHandler(v1, db, rdb, config.KafkaURL)
	handler.NewWarehouseHandler(v1, usecase.NewWarehouseUsecase(redisLock, repository.NewWarehouseRepository(db, &timer, &uuidHelper, producer)))
	// V1 Endpoints

	s := http.Server{
		Addr:    fmt.Sprintf(":%d", config.AppPort),
		Handler: e,
	}

	logger.LogWithContext(context.Background()).Info(fmt.Sprintf("%s service started...", config.AppName))
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

	logger.LogWithContext(context.Background()).Info(fmt.Sprintf("%s service finished", config.AppName))
}
