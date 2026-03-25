// Package routes registers API endpoints and applies authentication and authorization middleware.
package routes

import (
	"wacdo/controllers"
	"wacdo/middlewares"

	"github.com/gin-gonic/gin"
)

func UsersRoutes(router *gin.Engine) {
	public := router.Group("/users")
	{
		public.POST("/login", controllers.Login)
	}

	// Password change — any authenticated user (controller enforces own-password-only for non-admins)
	authenticated := router.Group("/users")
	authenticated.Use(middlewares.Authentication())
	{
		authenticated.PATCH("/:id/password", controllers.ChangePassword)
	}

	// User management is admin-only
	protected := router.Group("/users")
	protected.Use(middlewares.Authentication(), middlewares.Authorization("admin"))
	{
		protected.POST("/", controllers.CreateUser)
		protected.DELETE("/:id", controllers.DeleteUser)
		protected.GET("/", controllers.GetUsers)
		protected.GET("/:id", controllers.GetUser)
		protected.PATCH("/:id/status", controllers.ToggleUserStatus)
		protected.PATCH("/:id/reset-password", controllers.ResetPassword)
	}
}
