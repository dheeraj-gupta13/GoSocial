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
	routes.ProfileRoutes(protected)

	fmt.Printf("Listening on port %s", port)
	router.Run(":" + port)

}
