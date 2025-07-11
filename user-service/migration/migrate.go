package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/kadekchresna/go-boilerplate/config"

	_ "github.com/lib/pq"
)

func main() {

	if os.Getenv("APP_ENV") != config.PRODUCTION {

		// init invoke env before everything
		if err := godotenv.Load(); err != nil {
			panic(fmt.Errorf("error load ENV, %s", err.Error()))
		}

	}

	cfg := config.InitConfig()

	db, err := sql.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		log.Fatalf("Could not create database driver: %v", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/migration", currentDir),
		"postgres", driver)
	if err != nil {
		log.Fatalf("Could not create migrate instance: %v", err)
	}

	switch command := os.Args[1]; command {
	case "up":
		err := m.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Could not apply up migrations: %v", err)
		}
		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			log.Fatalf("Could not get migration version: %v", err)
		}
		log.Printf("Successfully migrated up to version %d (dirty: %v)\n", version, dirty)

	case "down":
		err := m.Down()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Could not apply down migrations: %v", err)
		}
		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			log.Printf("All migrations have been rolled back")
		} else {
			log.Printf("Successfully migrated down to version %d (dirty: %v)\n", version, dirty)
		}
	default:
		log.Fatalf("Unknown command: %s", command)
	}

	log.Println("Migration completed successfully")
}
