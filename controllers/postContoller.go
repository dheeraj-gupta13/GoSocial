package controllers

import (
	"fmt"
	"net/http"
	"time"

	"social-backend/database"
	"social-backend/middleware"

	"github.com/gin-gonic/gin"
)

type Post struct {
	Content string `json:"content"`
	Image   string `json:"imageUrl"`
}

/*


"post_id":    id,
			"content":    content,
			"image_url":  imageUrl,
			"created_at": createdAt,
*/

type UserPost struct {
	PostId     string    `json:"postId"`
	Image      string    `json:"imageUrl"`
	Content    string    `json:"content"`
	Created_at string    `json:"created_at"`
	Comments   []Comment `json:"comments"`
}

func AddPost(c *gin.Context) {

	currentUsername, currentUserId := middleware.GetCurrentUser(c)

	var post Post
	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fmt.Println(post)

	db := database.GetDB()
	query := `INSERT INTO posts (createdby, content, image_url, created_at) VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(query, currentUserId, post.Content, post.Image, time.Now())
	if err != nil {
		fmt.Println("*********", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": currentUsername, "id": currentUserId, "content": post.Content})
}

func DeletePost(c *gin.Context) {

}

func GetUsersPost(c *gin.Context) {

	currentUsername, currentUserId := middleware.GetCurrentUser(c)

	db := database.GetDB()

	query := `SELECT id, content, image_url, created_at FROM posts WHERE createdby = $1`

	rows, err := db.Query(query, currentUserId)
	if err != nil {
		fmt.Println("*********", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while retrieving posts"})
		return
	}
	defer rows.Close()

	var posts []gin.H

	for rows.Next() {
		var id int
		var content string
		var imageUrl string
		var createdAt string

		err := rows.Scan(&id, &content, &imageUrl, &createdAt)
		if err != nil {
			fmt.Println("*********", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating post"})
			return
		}

		post := gin.H{
			"post_id":    id,
			"content":    content,
			"image_url":  imageUrl,
			"created_at": createdAt,
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("*********", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error after retrieving posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": currentUsername,
		"user_id":  currentUserId,
		"posts":    posts,
	})

}
