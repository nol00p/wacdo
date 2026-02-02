package routes

import (
	"wacdo/controllers"

	"github.com/gin-gonic/gin"
)

func RolesRoutes(router *gin.Engine) {
	routesGroup := router.Group("/roles")
	{
		routesGroup.POST("/create", controllers.CreateRole)
	}
}
