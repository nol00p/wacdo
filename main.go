package main

import (
	"log"
	"os"
	"wacdo/config"
	"wacdo/models"
	"wacdo/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "wacdo/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title WacDo
// @version 1.0
// @description Super Ordening System
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
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
	routes.ProductRoutes(router)
	routes.CategoriesRoutes(router)
	routes.OptionRoutes(router)
	routes.OptionValueRoutes(router)
	routes.MenuRoutes(router)

	// Swagger routes
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Connect to DB
	config.ConnectDB()

	// DB migration
	config.DB.AutoMigrate(
		&models.Users{},
		&models.Roles{},
		&models.Category{},
		&models.Products{},
		&models.ProductOptions{},
		&models.OptionValues{},
		&models.Menu{},
	)

	// Start Server on PORT from env (Render sets this), fallback to 8000
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router.Run(":" + port)
}
