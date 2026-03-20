package routes

import (
	"wacdo/controllers"
	"wacdo/middlewares"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(router *gin.Engine) {
	// Read access: all roles (needed for order creation)
	readGroup := router.Group("/menus")
	readGroup.Use(middlewares.Authentication())
	{
		readGroup.GET("/", controllers.GetMenus)
		readGroup.GET("/:id", controllers.GetMenu)
		readGroup.GET("/:id/products/", controllers.GetMenuProducts)
	}

	// Write access: admin only
	writeGroup := router.Group("/menus")
	writeGroup.Use(middlewares.Authentication(), middlewares.Authorization("admin"))
	{
		writeGroup.POST("/", controllers.CreateMenu)
		writeGroup.PUT("/:id", controllers.UpdateMenu)
		writeGroup.DELETE("/:id", controllers.DeleteMenu)
		writeGroup.PATCH("/:id/availability", controllers.ToggleMenuAvailability)
		writeGroup.POST("/:id/products/", controllers.AddProductToMenu)
		writeGroup.DELETE("/products/:id", controllers.RemoveProductFromMenu)
	}
}
