package controllers

import (
	"fmt"
	"net/http"
	"social-backend/database"
	"social-backend/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

type Comment struct {
	PostId  int    `json:"postId"`
	Comment string `json:"comment"`
}

func AddComment(c *gin.Context) {
	_, currentUserId := middleware.GetCurrentUser(c)

	var comment Comment
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	db := database.GetDB()
	query := `INSERT INTO comments (post_id, user_id, comment, created_at) VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(query, comment.PostId, currentUserId, comment.Comment, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment added successfully"})
}

func GetCommentsForPost(c *gin.Context) {

	postId := c.Query("id")

	fmt.Println("postId", postId)

	db := database.GetDB()
	query := `SELECT id, user_id, comment, created_at FROM comments WHERE post_id = $1`

	rows, err := db.Query(query, postId)
	if err != nil {
		fmt.Println("*********", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while retrieving posts"})
		return
	}
	defer rows.Close()

	var comments []gin.H

	for rows.Next() {
		var id int
		var user_id int
		var comment string
		var createdAt string

		err := rows.Scan(&id, &user_id, &comment, &createdAt)
		if err != nil {
			fmt.Println("*********", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating post"})
			return
		}

		post := gin.H{
			"post_id":    id,
			"comment":    comment,
			"user_id":    user_id,
			"created_at": createdAt,
		}
		comments = append(comments, post)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("*********", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error after retrieving posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}
