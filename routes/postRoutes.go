package routes

import (
	controller "social-backend/controllers"

	"github.com/gin-gonic/gin"
)

func PostRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/post", controller.AddPost)
	incomingRoutes.GET("/post/getUserPosts", controller.GetUsersPost)
	incomingRoutes.GET("/post/getAllPosts", controller.GetAllPosts)
	incomingRoutes.DELETE("/post", controller.DeletePost)
	incomingRoutes.POST("/save", controller.SavePost)
}
