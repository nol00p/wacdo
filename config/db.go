package config

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	user := os.Getenv("USER")
	pass := os.Getenv("DB_PASS")

    // Basic check in case env vars are missing
    if user == "" || pass == "" {
        log.Fatal("USER or DB_PASS environment variable is not set")
    }

	// Set DB connection
    dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/sharepoint?charset=utf8mb4&parseTime=True&loc=Local", user, pass)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error: Connection to database failed :", err)
	}

	DB = db
}
