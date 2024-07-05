package controllers

import (
	"fmt"
	"net/http"
	"social-backend/database"
	"social-backend/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

type PostID struct {
	PostID int `json:"postId"`
}

func AddLike(c *gin.Context) {

	_, currentUserId := middleware.GetCurrentUser(c)

	var postId PostID
	if err := c.BindJSON(&postId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	db := database.GetDB()
	query := `INSERT INTO likes (post_id, user_id, created_at) VALUES ($1, $2, $3)`

	_, err := db.Exec(query, postId.PostID, currentUserId, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating like"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "like added successfully"})

}

func Unlike(c *gin.Context) {

	_, currentUserId := middleware.GetCurrentUser(c)

	var postId PostID
	if err := c.BindJSON(&postId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	db := database.GetDB()
	query := `DELETE FROM likes WHERE post_id = $1 AND user_id = $2`

	_, err := db.Exec(query, postId.PostID, currentUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while removing like"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "like removed successfully"})

}

func GetLikes(c *gin.Context) {

	postId := c.Query("id")

	db := database.GetDB()

	query := ` SELECT u.username
        FROM likes l
        JOIN users u ON l.user_id = u.id
        WHERE l.post_id = $1`

	rows, err := db.Query(query, postId)
	if err != nil {
		fmt.Println("EEEEEEEEEE", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error querying likes"})
		return
	}
	defer rows.Close()

	var likes []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row"})
			return
		}
		likes = append(likes, username)
	}

	c.JSON(http.StatusOK, gin.H{"message": likes})
}
