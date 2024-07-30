package routes

import (
	controller "social-backend/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/users", controller.GetUsers)
	incomingRoutes.POST("/users", controller.PostUser)
	incomingRoutes.GET("/user", controller.GetUserInfo)
}
