package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Render provides a single DATABASE_URL for PostgreSQL
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		// Fallback to individual env vars for local development
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASS")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_NAME")

		if dbUser == "" || dbPass == "" {
			log.Fatal("DATABASE_URL or DB_USER/DB_PASS environment variables are not set")
		}

		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPass, dbName)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error: Connection to database failed :", err)
	}

	DB = db
}
