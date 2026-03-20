package routes

import (
	"wacdo/controllers"
	"wacdo/middlewares"

	"github.com/gin-gonic/gin"
)

func CustomerRoutes(router *gin.Engine) {
	// Customer management: admin and accueil (accueil takes orders and needs customer data)
	routesGroup := router.Group("/customers")
	routesGroup.Use(middlewares.Authentication(), middlewares.Authorization("admin", "accueil"))
	{
		routesGroup.POST("/", controllers.CreateCustomer)
		routesGroup.GET("/", controllers.GetCustomers)
		routesGroup.GET("/:id", controllers.GetCustomer)
		routesGroup.PUT("/:id", controllers.UpdateCustomer)
		routesGroup.DELETE("/:id", controllers.DeleteCustomer)
	}
}
