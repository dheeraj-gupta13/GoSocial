/*

Table user {
  id bigserial [pk]
  email varchar
  username varchar [not null]
  password bigint [not null]
  created_at timestampz [default: `now()`]
}

Table profile {
  id bigserial [pk]
  userId bigserial [ref:>user.id]
  image varchar [not null]
  headline varchar
  name varchar
  created_at timestampz [default: `now()`]
}

Table post {
  id bigserial [pk]
  createdby bigserial [ref:>user.id]
  content varchar
  imageUrl string
  created_at timestampz [default: `now()`]
}

Table Likes {
  id bigserial [pk]
  postId bigserial [ref:> post.id]
  userId bigserial [ref:> user.id]

}

Table Comments {
  id bigserial [pk]
  postId bigserial [ref:> post.id]
  userId bigserial [ref:> user.id]
  comment varchar
}

Table Followers {
  follower_id bigserial [ref :> user.id]
  followee_id bigserial [ref :> user.id]
  followed_at timestampz [default: `now()`]
}


*/

package main

import (
	"fmt"
	"os"

	"social-backend/middleware"

	"social-backend/routes"

	"social-backend/database"

	"github.com/gin-gonic/gin"
)

func main() {

	// Initialize the database connection
	_, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.New()

	// config := cors.DefaultConfig()
	// config.AllowAllOrigins = true
	// router.Use(cors.New(config))

	router.Use(gin.Logger())

	routes.AuthRoutes(router)

	router.Use(middleware.Authentication())

	routes.UserRoutes(router)
	routes.PostRoutes(router)
	routes.CommentRoutes(router)
	routes.LikeRoutes(router)

	fmt.Printf("Listening on port", port)
	router.Run(":" + port)

}
