package routes

import (
	"wacdo/controllers"
	"wacdo/middlewares"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(router *gin.Engine) {
	// View orders: all roles
	viewGroup := router.Group("/orders")
	viewGroup.Use(middlewares.Authentication())
	{
		viewGroup.GET("/", controllers.GetOrders)
		viewGroup.GET("/:id", controllers.GetOrder)
	}

	// Create and cancel orders: admin + accueil
	accueilGroup := router.Group("/orders")
	accueilGroup.Use(middlewares.Authentication(), middlewares.Authorization("admin", "accueil"))
	{
		accueilGroup.POST("/", controllers.CreateOrder)
		accueilGroup.PATCH("/:id/cancel", controllers.CancelOrder)
	}

	// Update order status (preparing → prepared → delivered): admin + preparation + accueil
	// Preparation marks as prepared, accueil marks as delivered
	statusGroup := router.Group("/orders")
	statusGroup.Use(middlewares.Authentication(), middlewares.Authorization("admin", "preparation", "accueil"))
	{
		statusGroup.PATCH("/:id/status", controllers.UpdateOrderStatus)
	}

	// Customer orders: admin + accueil
	customersGroup := router.Group("/customers")
	customersGroup.Use(middlewares.Authentication(), middlewares.Authorization("admin", "accueil"))
	{
		customersGroup.GET("/:id/orders", controllers.GetOrdersByCustomer)
	}
}
