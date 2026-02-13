package routes

import (
	"wacdo/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(router *gin.Engine) {
	routesGroup := router.Group("/menus")
	{
		routesGroup.POST("/", controllers.CreateMenu)
		routesGroup.GET("/", controllers.GetMenus)
		routesGroup.GET("/:id", controllers.GetMenu)
		routesGroup.PUT("/:id", controllers.UpdateMenu)
		routesGroup.DELETE("/:id", controllers.DeleteMenu)
		routesGroup.PATCH("/:id/availability", controllers.ToggleMenuAvailability)
		routesGroup.POST("/:id/products/", controllers.AddProductToMenu)
		routesGroup.GET("/:id/products/", controllers.GetMenuProducts)
		routesGroup.DELETE("/products/:id", controllers.RemoveProductFromMenu)
	}
}
