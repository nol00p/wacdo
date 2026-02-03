package routes

import (
	"wacdo/controllers"

	"github.com/gin-gonic/gin"
)

func RolesRoutes(router *gin.Engine) {
	routesGroup := router.Group("/roles")
	{
		routesGroup.POST("/", controllers.CreateRole)
		routesGroup.DELETE("/:id", controllers.DeleteRole)
	}
}
