package routes

import (
	controller "social-backend/controllers"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.POST("/comment", controller.AddComment)
	incomingRoutes.GET("/comment", controller.GetCommentsForPost)
}
