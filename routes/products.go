package routes

import (
	"wacdo/controllers"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.Engine) {
	routesGroup := router.Group("/products")
	{
		routesGroup.POST("/", controllers.CreateProduct)
		routesGroup.DELETE("/:id", controllers.DeleteProduct)
		routesGroup.GET("/", controllers.GetProducts)
		routesGroup.GET("/:id", controllers.GetProduct)
		routesGroup.PUT("/:id", controllers.UpdateProduct)
		routesGroup.GET("/category/:category_id", controllers.GetProductsByCategory)
		routesGroup.PATCH("/:id/availability", controllers.ToggleProductAvailability)
		routesGroup.PATCH("/:id/stock", controllers.UpdateProductStock)
	}
}

func CategoriesRoutes(router *gin.Engine) {
	routesGroup := router.Group("/categories")
	{
		routesGroup.POST("/", controllers.CreateCategory)
		routesGroup.DELETE("/:id", controllers.DeleteCategory)
		routesGroup.GET("/", controllers.GetCategories)
		routesGroup.GET("/:id", controllers.GetCategory)
		routesGroup.PUT("/:id", controllers.UpdateCategory)
	}
}

func OptionRoutes(router *gin.Engine) {
	routesGroup := router.Group("/options")
	{
		routesGroup.POST("/", controllers.CreateOption)
		routesGroup.DELETE("/:id", controllers.DeleteOption)
		routesGroup.GET("/", controllers.GetOptions)
		routesGroup.GET("/:id", controllers.GetOption)
		routesGroup.PUT("/:id", controllers.UpdateOption)
		routesGroup.GET("/product/:product_id", controllers.GetOptionsByProduct)
	}
}

func OptionValueRoutes(router *gin.Engine) {
	routesGroup := router.Group("/options")
	{
		routesGroup.POST("/:id/values/", controllers.CreateOptionValue)
		routesGroup.GET("/:id/values/", controllers.GetValuesByOption)
		routesGroup.GET("/values/:id", controllers.GetOptionValue)
		routesGroup.PUT("/values/:id", controllers.UpdateOptionValue)
		routesGroup.DELETE("/values/:id", controllers.DeleteOptionValue)
	}
}
