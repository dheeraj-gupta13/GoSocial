package routes

import (
	controller "social-backend/controllers"

	"github.com/gin-gonic/gin"
)

func LikeRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/like", controller.AddLike)
	incomingRoutes.DELETE("/like", controller.Unlike)
	incomingRoutes.GET("/like", controller.GetLikes)
}
