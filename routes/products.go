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

func CategoryRoutes(router *gin.Engine) {
	routesGroup := router.Group("/category")
	{
		routesGroup.POST("/", controllers.CreateCategory)
		routesGroup.DELETE("/:id", controllers.DeleteCategory)
		routesGroup.GET("/", controllers.GetCategories)
		routesGroup.GET("/:id", controllers.GetCategory)
		routesGroup.PUT("/:id", controllers.UpdateCategory)
	}
}
