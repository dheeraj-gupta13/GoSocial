package routes

import (
	controller "social-backend/controllers"

	"github.com/gin-gonic/gin"
)

func ProfileRoutes(incomingRoutes *gin.RouterGroup) {

	incomingRoutes.GET("/profile", controller.GetProfile)
	incomingRoutes.POST("/follow", controller.Follow)
	incomingRoutes.DELETE("/unfollow", controller.UnFollow)
	// incomingRoutes.POST("/users", controller.PostUser)
	// incomingRoutes.GET("/user", controller.GetUserInfo)
}
