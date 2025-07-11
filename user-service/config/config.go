package config

import (
	"os"

	"github.com/kadekchresna/go-boilerplate/helper/env"
	"gorm.io/gorm"
)

const (
	STAGING    = `staging`
	PRODUCTION = `production`
)

type Config struct {
	AppName      string
	AppPort      int
	AppEnv       string
	AppJWTSecret string

	DatabaseDSN string
}

type DB struct {
	MasterDB   *gorm.DB
	SlaveDB    *gorm.DB
	AnalyticDB *gorm.DB
}

func InitConfig() Config {
	return Config{
		AppName:      os.Getenv("APP_NAME"),
		AppEnv:       os.Getenv("APP_ENV"),
		AppPort:      env.GetEnvInt("APP_PORT"),
		AppJWTSecret: os.Getenv("APP_JWT_SECRET"),

		DatabaseDSN: os.Getenv("DB_DSN"),
	}
}
