package routes

import (
	controller "social-backend/controllers"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/comment", controller.AddComment)
	incomingRoutes.GET("/comment", controller.GetCommentsForPost)
}
