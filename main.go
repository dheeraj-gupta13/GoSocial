/*

Table users {
  id bigserial [pk]
  email varchar
  username varchar [not null]
  password varchar [not null]
  created_at timestampz [default: `now()`]
}

Table profile {
  id bigserial [pk]
  userId bigserial [ref:>users.id]
  image varchar [not null]
  headline varchar
  name varchar
  created_at timestampz [default: `now()`]
}

Table post {
  id bigserial [pk]
  createdby bigserial [ref:>users.id]
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
  commented_at timestampz [default: `now()`]
}

Table Followers {
  follower_id bigserial [ref :> user.id]
  followee_id bigserial [ref :> user.id]
  followed_at timestampz [default: `now()`]
}

Table Saved {
	id bigserial [pk],
	postId bigserial [ref:> post.id]
  	userId bigserial [ref:> users.id]
  	created_at timestampz [default: `now()`]
}

CREATE TABLE Saved (
    id BIGSERIAL PRIMARY KEY,
    postId BIGSERIAL REFERENCES post(id),
    userId BIGSERIAL REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT now()
);


*/

package main

import (
	"fmt"
	"log"
	"os"

	"social-backend/middleware"

	"social-backend/routes"

	"social-backend/database"

	"github.com/cloudinary/cloudinary-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var cld *cloudinary.Cloudinary

func main() {

	err1 := godotenv.Load()
	if err1 != nil {
		log.Fatal("Error loading .env file")
	}

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

	config := cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}
	router.Use(cors.New(config))
	router.Use(gin.Logger())

	cld, err = cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Fatal("Cloudinary init failed:", err)
	}

	public := router.Group("/")
	routes.AuthRoutes(public)

	protected := router.Group("/")
	protected.Use(middleware.Authentication())

	routes.UserRoutes(protected)
	routes.PostRoutes(protected)
	routes.CommentRoutes(protected)
	routes.LikeRoutes(protected)

	fmt.Printf("Listening on port %s", port)
	router.Run(":" + port)

}
