package routes

import (
	"wacdo/controllers"

	"github.com/gin-gonic/gin"
)

func UsersRoutes(router *gin.Engine) {
	routesGroup := router.Group("/users")
	{
		routesGroup.POST("/register", controllers.Register)
		routesGroup.POST("/login", controllers.Login)
	}
}
