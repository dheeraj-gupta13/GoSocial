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
	Post_id int    `json:"post_id"`
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

	_, err := db.Exec(query, comment.Post_id, currentUserId, comment.Comment, time.Now())
	if err != nil {
		fmt.Println("EEEEEEEEEE", err)
		fmt.Println("POST ID", comment.Post_id)
		fmt.Println("USER ID", comment.Comment)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment added successfully"})
}

func GetCommentsForPost(c *gin.Context) {

	postId := c.Query("id")

	fmt.Println("postId", postId)

	db := database.GetDB()
	query := `SELECT c.comment_id, u.username, c.comment, c.created_at 
	FROM comments c JOIN users u 
	ON c.user_id = u.user_id
	WHERE post_id = $1`

	rows, err := db.Query(query, postId)
	if err != nil {
		fmt.Println("*********", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while retrieving posts"})
		return
	}
	defer rows.Close()

	var comments []gin.H

	for rows.Next() {
		var comment_id int
		var username string
		var comment string
		var createdAt string

		err := rows.Scan(&comment_id, &username, &comment, &createdAt)
		if err != nil {
			fmt.Println("*********", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating post"})
			return
		}

		post := gin.H{
			"comment_id": comment_id,
			"comment":    comment,
			"username":   username,
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
