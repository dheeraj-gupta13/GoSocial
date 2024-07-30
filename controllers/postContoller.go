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

// type PostID struct {
// 	PostID int `json:"postId"`
// }

func AddPost(c *gin.Context) {

	currentUsername, currentUserId := middleware.GetCurrentUser(c)

	var post Post
	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fmt.Println(post)

	db := database.GetDB()
	query := `INSERT INTO post (createdby, content, imageUrl, created_at) VALUES ($1, $2, $3, $4)`

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

	query := `SELECT id, content, imageUrl, created_at FROM post WHERE createdby = $1`

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

func GetAllPosts(c *gin.Context) {
	currentUsername, currentUserId := middleware.GetCurrentUser(c)

	db := database.GetDB()

	query := `
		SELECT post.id, post.content, post.imageUrl, post.created_at, p.user_id, p.image, p.name, u.username
		FROM post 
		INNER JOIN profile p ON post.createdby = p.user_id
		INNER JOIN users u ON p.user_id = u.id
	`

	rows, err := db.Query(query)
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
		var userId int
		var userProfilePic string
		var name string
		var username string

		err := rows.Scan(&id, &content, &imageUrl, &createdAt, &userId, &userProfilePic, &name, &username)
		if err != nil {
			fmt.Println("*********", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating post"})
			return
		}

		post := gin.H{
			"post_id":          id,
			"content":          content,
			"image_url":        imageUrl,
			"created_at":       createdAt,
			"user_id":          userId,
			"user_profile_pic": userProfilePic,
			"name":             name,
			"username":         username,
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

func SavePost(c *gin.Context) {

	_, currentUserId := middleware.GetCurrentUser(c)

	var postId PostID
	if err := c.BindJSON(&postId); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fmt.Println(postId)

	db := database.GetDB()

	query := `INSERT INTO Saved (userId, postId, created_at) VALUES ($1, $2, $3)`

	_, err := db.Exec(query, currentUserId, postId.PostID, time.Now())
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating like"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post saved successfully"})
}

// func GetAllSavedPost(c *gin.Context){
// 	 _, currentUserId := middleware.GetCurrentUser(c);

// 	 db := database.GetDB();

// 	 query := `SELECT `

// }

// func GetLikesInfo(c *gin.Context) {
// 	// like count and weather I have liked a post or not
// 	currentUsername, currentUserId := middleware.GetCurrentUser(c)

// 	db := database.GetDB()

// 	query := `
// 		SELECT post.id, post.content, post.imageUrl, post.created_at, p.user_id, p.image, p.name, u.username
// 		FROM post
// 		INNER JOIN profile p ON post.createdby = p.user_id
// 		INNER JOIN users u ON p.user_id = u.id
// 	`

// 	c.JSON(http.StatusOK, gin.H{
// 		"username": currentUsername,
// 		"user_id":  currentUserId,
// 	})

// }
