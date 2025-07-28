package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"social-backend/database"
	"social-backend/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

type ReactedPost struct {
	PostID        int `json:"post_id"`
	Reaction_type int `json:"reaction_type"`
}

// React a post
func AddLike(c *gin.Context) {

	_, currentUserId := middleware.GetCurrentUser(c)

	var postId ReactedPost
	if err := c.BindJSON(&postId); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fmt.Println(postId)

	db := database.GetDB()

	query1 := `DELETE FROM postReactions WHERE post_id = $1 AND user_id = $2`
	_, err1 := db.Exec(query1, postId.PostID, currentUserId)
	if err1 != nil {
		fmt.Println(err1)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while checking if the post is already liked"})
		return
	}

	query := `INSERT INTO postReactions (post_id, user_id, reaction_type, created_on) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, postId.PostID, currentUserId, postId.Reaction_type, time.Now())
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating like"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "like added successfully"})

}

// UnReact a post
func Unlike(c *gin.Context) {

	_, currentUserId := middleware.GetCurrentUser(c)

	postId := c.Query("post_id")
	if postId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "postId is required"})
		return
	}

	fmt.Println("1", postId)

	db := database.GetDB()
	query := `DELETE FROM postReactions WHERE post_id = $1 AND user_id = $2`

	_, err := db.Exec(query, postId, currentUserId)
	if err != nil {
		fmt.Println("2", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while removing like"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "like removed successfully"})

}

// get all reactions of a posts.
func GetLikes(c *gin.Context) {

	postId := c.Query("post_id")

	fmt.Println("DHEERAJ", postId)

	db := database.GetDB()

	query := `SELECT u.username
        FROM postReactions l
        JOIN users u ON l.user_id = u.user_id
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

func GetReactionId(c *gin.Context) {

	postId := c.Query("post_id")
	_, currentUserId := middleware.GetCurrentUser(c)

	db := database.GetDB()

	query := `SELECT reaction_type FROM postreactions WHERE user_id=$1  AND post_id=$2`

	var reaction_type int
	err := db.QueryRow(query, currentUserId, postId).Scan(&reaction_type)
	if err != nil {

		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"reaction_type": -1})
			return
		}

		fmt.Println("EEEEEEEEEE", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error querying likes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reaction_type": reaction_type})

}
