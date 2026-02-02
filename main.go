package main

import (
	"log"
	"wacdo/config"
	"wacdo/models"
	"wacdo/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()

	// Proxie rules
	router.SetTrustedProxies(nil)

	// Security middleware
	router.Use(config.SecurityMiddleware())
	router.Use(config.CORSMiddleware())
	router.Use(config.RateLimit(100))

	// Loading .ENV vrariables
	err := godotenv.Load()
	if err != nil {
		log.Println("file not found: .ENV")
	}
	// API router definition
	routes.UsersRoutes(router)
	routes.RolesRoutes(router)

	// Connect to DB
	config.ConnectDB()

	// DB migration
	config.DB.AutoMigrate(
		&models.Users{},
		&models.Roles{},
	)

	// Start Server on port 8000
	router.Run(":8000")
}
