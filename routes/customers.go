package routes

import (
	"wacdo/controllers"

	"github.com/gin-gonic/gin"
)

func CustomerRoutes(router *gin.Engine) {
	routesGroup := router.Group("/customers")
	{
		routesGroup.POST("/", controllers.CreateCustomer)
		routesGroup.GET("/", controllers.GetCustomers)
		routesGroup.GET("/:id", controllers.GetCustomer)
		routesGroup.PUT("/:id", controllers.UpdateCustomer)
		routesGroup.DELETE("/:id", controllers.DeleteCustomer)
	}
}
