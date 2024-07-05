package routes

import (
	"fmt"
	controller "social-backend/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {

	fmt.Printf("=============>")

	incomingRoutes.POST("/login", controller.Login)
	incomingRoutes.POST("/register", controller.Register)
}
