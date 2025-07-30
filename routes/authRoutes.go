package routes

import (
	controller "social-backend/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.POST("/login", controller.Login)
	incomingRoutes.POST("/register", controller.Register)
	incomingRoutes.GET("/validateToken", controller.ValidateToken)
}
