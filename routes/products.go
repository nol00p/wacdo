package routes

import (
	"wacdo/controllers"
	"wacdo/middlewares"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.Engine) {
	// Read access: all roles (needed for order creation)
	readGroup := router.Group("/products")
	readGroup.Use(middlewares.Authentication())
	{
		readGroup.GET("/", controllers.GetProducts)
		readGroup.GET("/:id", controllers.GetProduct)
		readGroup.GET("/category/:category_id", controllers.GetProductsByCategory)
	}

	// Write access: admin only
	writeGroup := router.Group("/products")
	writeGroup.Use(middlewares.Authentication(), middlewares.Authorization("admin"))
	{
		writeGroup.POST("/", controllers.CreateProduct)
		writeGroup.PUT("/:id", controllers.UpdateProduct)
		writeGroup.DELETE("/:id", controllers.DeleteProduct)
		writeGroup.PATCH("/:id/availability", controllers.ToggleProductAvailability)
		writeGroup.PATCH("/:id/stock", controllers.UpdateProductStock)
	}
}

func CategoriesRoutes(router *gin.Engine) {
	// Read access: all roles
	readGroup := router.Group("/categories")
	readGroup.Use(middlewares.Authentication())
	{
		readGroup.GET("/", controllers.GetCategories)
		readGroup.GET("/:id", controllers.GetCategory)
	}

	// Write access: admin only
	writeGroup := router.Group("/categories")
	writeGroup.Use(middlewares.Authentication(), middlewares.Authorization("admin"))
	{
		writeGroup.POST("/", controllers.CreateCategory)
		writeGroup.PUT("/:id", controllers.UpdateCategory)
		writeGroup.DELETE("/:id", controllers.DeleteCategory)
	}
}

func OptionRoutes(router *gin.Engine) {
	// Read access: all roles
	readGroup := router.Group("/options")
	readGroup.Use(middlewares.Authentication())
	{
		readGroup.GET("/", controllers.GetOptions)
		readGroup.GET("/:id", controllers.GetOption)
		readGroup.GET("/product/:product_id", controllers.GetOptionsByProduct)
	}

	// Write access: admin only
	writeGroup := router.Group("/options")
	writeGroup.Use(middlewares.Authentication(), middlewares.Authorization("admin"))
	{
		writeGroup.POST("/", controllers.CreateOption)
		writeGroup.PUT("/:id", controllers.UpdateOption)
		writeGroup.DELETE("/:id", controllers.DeleteOption)
	}
}

func OptionValueRoutes(router *gin.Engine) {
	// Read access: all roles
	readGroup := router.Group("/options")
	readGroup.Use(middlewares.Authentication())
	{
		readGroup.GET("/:id/values/", controllers.GetValuesByOption)
		readGroup.GET("/values/:id", controllers.GetOptionValue)
	}

	// Write access: admin only
	writeGroup := router.Group("/options")
	writeGroup.Use(middlewares.Authentication(), middlewares.Authorization("admin"))
	{
		writeGroup.POST("/:id/values/", controllers.CreateOptionValue)
		writeGroup.PUT("/values/:id", controllers.UpdateOptionValue)
		writeGroup.DELETE("/values/:id", controllers.DeleteOptionValue)
	}
}
