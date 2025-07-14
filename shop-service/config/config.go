package config

import (
	"os"

	"github.com/kadekchresna/ecommerce/shop-service/helper/env"
	"gorm.io/gorm"
)

const (
	STAGING    = `staging`
	PRODUCTION = `production`
)

type Config struct {
	AppName        string
	AppPort        int
	AppEnv         string
	AppStaticToken string
	AppJWTSecret   string

	DatabaseDSN string
	RedisURL    string
	KafkaURL    string

	WarehouseServiceURL string
	ShopServiceURL      string
	ProductServiceURL   string
}

type DB struct {
	MasterDB   *gorm.DB
	SlaveDB    *gorm.DB
	AnalyticDB *gorm.DB
}

func InitConfig() Config {
	return Config{
		AppName:        os.Getenv("APP_NAME"),
		AppEnv:         os.Getenv("APP_ENV"),
		AppPort:        env.GetEnvInt("APP_PORT"),
		AppJWTSecret:   os.Getenv("APP_JWT_SECRET"),
		AppStaticToken: os.Getenv("APP_STATIC_TOKEN"),

		DatabaseDSN:         os.Getenv("DB_DSN"),
		KafkaURL:            os.Getenv("KAFKA_URL"),
		RedisURL:            os.Getenv("REDIS_URL"),
		WarehouseServiceURL: os.Getenv("WAREHOUSE_SERVICE_URL"),
		ShopServiceURL:      os.Getenv("SHOP_SERVICE_URL"),
		ProductServiceURL:   os.Getenv("PRODUCT_SERVICE_URL"),
	}
}
