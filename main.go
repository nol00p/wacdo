package main

import (
	"log"
	"os"
	"wacdo/config"
	"wacdo/models"
	"wacdo/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

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
	// Loading .ENV variables (must happen before middleware reads env)
	err := godotenv.Load()
	if err != nil {
		log.Println("file not found: .ENV")
	}

	router := gin.Default()

	// Proxie rules
	router.SetTrustedProxies(nil)

	// Security middleware
	router.Use(config.SecurityMiddleware())
	router.Use(config.CORSMiddleware())
	router.Use(config.RateLimit(100))
	// API router definition
	routes.UsersRoutes(router)
	routes.RolesRoutes(router)
	routes.ProductRoutes(router)
	routes.CategoriesRoutes(router)
	routes.OptionRoutes(router)
	routes.OptionValueRoutes(router)
	routes.MenuRoutes(router)
	routes.CustomerRoutes(router)
	routes.OrderRoutes(router)

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
		&models.MenuProduct{},
		&models.Customer{},
		&models.Order{},
		&models.OrderItem{},
		&models.OrderItemOption{},
	)

	// Seed default roles and admin user on first install
	seedDefaults()

	// Start Server on PORT from env (Render sets this), fallback to 8000
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router.Run(":" + port)
}

// seedDefaults creates the three default roles and an admin user on first install only.
// It checks that no roles AND no users exist — once the system has been set up, it never seeds again.
// The last-admin guard in DeleteUser and ToggleUserStatus ensures there is always at least one active admin.
func seedDefaults() {
	var roleCount int64
	config.DB.Model(&models.Roles{}).Count(&roleCount)
	var userCount int64
	config.DB.Model(&models.Users{}).Count(&userCount)

	if roleCount > 0 || userCount > 0 {
		return
	}

	log.Println("First install detected — seeding default roles and admin user...")

	roles := []models.Roles{
		{RoleName: "admin", Description: "Full access to all features"},
		{RoleName: "preparation", Description: "View orders and mark as prepared"},
		{RoleName: "accueil", Description: "Create and deliver orders"},
	}
	for i := range roles {
		config.DB.Create(&roles[i])
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Admin@1234"), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash default admin password:", err)
		return
	}

	admin := models.Users{
		Username: "admin",
		Email:    "admin@wacdo.fr",
		Password: string(hashedPassword),
		RolesID:  roles[0].ID,
		IsActive: true,
	}
	config.DB.Create(&admin)

	log.Println("Default roles created: admin, preparation, accueil")
	log.Println("Default admin user created — email: admin@wacdo.fr / password: Admin@1234")
	log.Println("Change the admin password after first login!")
}
