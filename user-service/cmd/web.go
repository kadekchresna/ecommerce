package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kadekchresna/ecommerce/user-service/config"
	"github.com/kadekchresna/ecommerce/user-service/helper/logger"
	driver_db "github.com/kadekchresna/ecommerce/user-service/infrastructure/db"
	handler "github.com/kadekchresna/ecommerce/user-service/internal/v1/delivery/http"
	"github.com/kadekchresna/ecommerce/user-service/internal/v1/repository"
	"github.com/kadekchresna/ecommerce/user-service/internal/v1/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

const (
	STAGING     = `stg`
	PRODUCTIOON = `prd`
)

func init() {
	if os.Getenv("APP_ENV") != PRODUCTIOON {

		// init invoke env before everything
		cobra.OnInitialize(initConfig)

	}

	// adding command invokable
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "web",
	Short: "Running Web Service",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func initConfig() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("error load ENV, %s", err.Error()))
	}
}

func run() {
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

	// V1 Endpoints
	v1 := e.Group("/api/v1")

	handler.NewHealthHandler(v1, db)

	handler.NewAuthHandler(v1, usecase.NewAuthUsecase(config, repository.NewUsersRepository(db)))
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
