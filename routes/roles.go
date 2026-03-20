package routes

import (
	"wacdo/controllers"
	"wacdo/middlewares"

	"github.com/gin-gonic/gin"
)

func RolesRoutes(router *gin.Engine) {
	// Role management is admin-only
	routesGroup := router.Group("/roles")
	routesGroup.Use(middlewares.Authentication(), middlewares.Authorization("admin"))
	{
		routesGroup.GET("/", controllers.GetRoles)
		routesGroup.GET("/:id", controllers.GetRole)
		routesGroup.POST("/", controllers.CreateRole)
		routesGroup.DELETE("/:id", controllers.DeleteRole)
	}
}
