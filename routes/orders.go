package routes

import (
	"wacdo/controllers"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(router *gin.Engine) {
	ordersGroup := router.Group("/orders")
	{
		ordersGroup.POST("/", controllers.CreateOrder)
		ordersGroup.GET("/", controllers.GetOrders)
		ordersGroup.GET("/:id", controllers.GetOrder)
		ordersGroup.PATCH("/:id/status", controllers.UpdateOrderStatus)
		ordersGroup.PATCH("/:id/cancel", controllers.CancelOrder)
	}

	// Customer orders route
	customersGroup := router.Group("/customers")
	{
		customersGroup.GET("/:id/orders", controllers.GetOrdersByCustomer)
	}
}
