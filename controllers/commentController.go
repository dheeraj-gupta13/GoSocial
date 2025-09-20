package controllers

import (
	"net/http"
	"social-backend/database"
	"social-backend/middleware"
	"social-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

/*
Add a comment on the post.
*/
func AddComment(c *gin.Context) {
	_, currentUserId := middleware.GetCurrentUser(c)

	var comment models.UserComment
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse[any]{
			Success: false,
			Error:   &models.ErrInvalidInput,
			Data:    nil,
		})
		return
	}

	query := `INSERT INTO comments (post_id, user_id, comment, created_at) VALUES ($1, $2, $3, $4)`
	db := database.GetDB()
	_, err := db.Exec(query, comment.Post_id, currentUserId, comment.Comment, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse[any]{
			Success: false,
			Error:   &models.ErrInternalServer,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse[any]{
		Success: true,
		Error:   nil,
		Data: gin.H{
			"message": "Comment added successfully",
		},
	})
}

/*
Get the comments on the post.
*/
func GetCommentsForPost(c *gin.Context) {

	postId := c.Query("id")

	query := `SELECT c.comment_id, u.username, c.comment, c.created_at, p.avatar_url 
	FROM comments c 
	JOIN users u ON c.user_id = u.user_id
	JOIN profile p ON p.user_id = u.user_id
	WHERE post_id = $1`

	db := database.GetDB()
	rows, err := db.Query(query, postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse[any]{
			Success: false,
			Error:   &models.ErrInternalServer,
			Data:    nil,
		})
		return
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var createdAt time.Time

		err := rows.Scan(&comment.CommentId, &comment.Username, &comment.Comment, &createdAt, &comment.AvatarURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse[any]{
				Success: false,
				Error:   &models.ErrInternalServer,
				Data:    nil,
			})
			return
		}

		createdFromNow := timeAgo(createdAt)
		comment.CreatedAt = createdFromNow
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse[any]{
			Success: false,
			Error:   &models.ErrInternalServer,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse[[]models.Comment]{
		Success: true,
		Error:   nil,
		Data:    comments,
	})
}
