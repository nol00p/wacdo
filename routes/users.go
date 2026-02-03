package routes

import (
	"wacdo/controllers"

	"github.com/gin-gonic/gin"
)

func UsersRoutes(router *gin.Engine) {
	routesGroup := router.Group("/users")
	{
		routesGroup.POST("/", controllers.CreateUser)
		routesGroup.DELETE("/:id", controllers.DeleteUser)
		routesGroup.GET("/", controllers.GetUsers)
		routesGroup.GET("/:id", controllers.GetUser)
		routesGroup.POST("/login", controllers.Login)
	}
}
