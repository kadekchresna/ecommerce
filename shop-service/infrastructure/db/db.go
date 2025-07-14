package driver_db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(DatabaseDSN string) *gorm.DB {
	dbConn, err := gorm.Open(postgres.Open(DatabaseDSN), &gorm.Config{
		Logger:             logger.Default.LogMode(logger.Info),
		PrepareStmt:        true,
		PrepareStmtMaxSize: 10,
	})
	if err != nil {
		panic(fmt.Errorf("error init db, %s", err.Error()))
	}

	sqlDB, err := dbConn.DB()
	if err != nil {
		panic(fmt.Errorf("error getting db from gorm, %s", err.Error()))
	}

	if err := sqlDB.Ping(); err != nil {
		panic(fmt.Errorf("error pinging db: %s", err.Error()))
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	return dbConn
}
