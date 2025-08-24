package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"social-backend/database"
	"social-backend/middleware"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
)

// User adds a post
func AddPost(c *gin.Context) {

	currentUsername, currentUserId := middleware.GetCurrentUser(c)

	db := database.GetDB()

	// Get form values
	content := c.PostForm("content")
	location := c.PostForm("location") // optional

	// Handle file upload
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	fmt.Println("1. file", file)

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
		return
	}
	defer openedFile.Close()

	// Initialize Cloudinary
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Cloudinary"})
		return
	}

	// Upload to Cloudinary
	uploadRes, err := cld.Upload.Upload(context.Background(), openedFile, uploader.UploadParams{
		PublicID: file.Filename,
		Folder:   "posts",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Image upload failed"})
		return
	}

	fmt.Print("2. uploadRes ", uploadRes)

	// Store post in DB
	query := `INSERT INTO posts (user_id, content, image_url, created_on, location) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(query, currentUserId, content, uploadRes.SecureURL, time.Now(), location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Post created successfully",
		"username":  currentUsername,
		"user_id":   currentUserId,
		"content":   content,
		"image_url": uploadRes.SecureURL,
	})

	// var post models.Post
	// if err := c.BindJSON(&post); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	// 	return
	// }

	// db := database.GetDB()
	// query := `INSERT INTO posts (user_id, content, image_url, created_on, location) VALUES ($1, $2, $3, $4, $5)`

	// _, err := db.Exec(query, currentUserId, post.Content, post.Image_url, time.Now(), "")
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating post"})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"username": currentUsername, "user_id": currentUserId, "content": post.Content})
}

// Delete the post
func DeletePost(c *gin.Context) {

}

// Get the post of some specific user
func GetPosts(c *gin.Context) {

}

// Get the post of single user
func GetUsersPost(c *gin.Context) {

	currentUsername, currentUserId := middleware.GetCurrentUser(c)
	user_id := c.Query("user_id")

	db := database.GetDB()

	query := `
		SELECT p.post_id, p.content, p.image_url, p.created_on, p.user_id, u.username, profile.avatar_url
		FROM posts p
		INNER JOIN profile ON p.user_id = profile.user_id
		INNER JOIN users u ON p.user_id = u.user_id
		WHERE u.user_id = $1;
	`

	rows, err := db.Query(query, user_id)
	if err != nil {
		fmt.Println("error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while retrieving posts"})
		return
	}
	defer rows.Close()

	var posts []gin.H

	for rows.Next() {
		var post_id int
		var content string
		var image_url string
		var created_at time.Time
		var user_id int
		var username string
		var avatar_url string

		err := rows.Scan(&post_id, &content, &image_url, &created_at, &user_id, &username, &avatar_url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating post"})
			return
		}

		likesQuery := `SELECT u.username, p.reaction_type 
							FROM postReactions p JOIN users u  
							ON p.user_id = u.user_id
						WHERE post_id = $1`
		reactions, err := db.Query(likesQuery, post_id)
		if err != nil {
			fmt.Println(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		var reactionArr []gin.H
		for reactions.Next() {
			var who_reacted string
			var what_reacted int

			err2 := reactions.Scan(&who_reacted, &what_reacted)
			if err2 != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err2})
				return
			}

			single_reaction := gin.H{
				"who_reacted":  who_reacted,
				"what_reacted": what_reacted,
			}

			reactionArr = append(reactionArr, single_reaction)
		}

		createdFromNow := timeAgo(created_at)
		post := gin.H{
			"post_id":        post_id,
			"content":        content,
			"image_url":      image_url,
			"created_at":     createdFromNow,
			"user_id":        user_id,
			"username":       username,
			"avatar_url":     avatar_url,
			"reaction_array": reactionArr,
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error after retrieving posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": currentUsername,
		"user_id":  currentUserId,
		"posts":    posts,
	})

}

// Get all the post that are in Db, this is mainly for feed
func GetAllPosts(c *gin.Context) {
	currentUsername, currentUserId := middleware.GetCurrentUser(c)

	db := database.GetDB()

	query := `
		SELECT p.post_id, p.content, p.image_url, p.created_on, p.user_id, u.username, profile.avatar_url
		FROM posts p
		INNER JOIN profile ON p.user_id = profile.user_id
		INNER JOIN users u ON p.user_id = u.user_id
	`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while retrieving posts"})
		return
	}
	defer rows.Close()

	var posts []gin.H

	for rows.Next() {
		var post_id int
		var content string
		var image_url string
		var created_at time.Time
		var user_id int
		var username string
		var avatar_url string

		err := rows.Scan(&post_id, &content, &image_url, &created_at, &user_id, &username, &avatar_url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating post"})
			return
		}

		likesQuery := `SELECT u.username, p.reaction_type 
							FROM postReactions p JOIN users u  
							ON p.user_id = u.user_id
						WHERE post_id = $1`
		reactions, err := db.Query(likesQuery, post_id)
		if err != nil {
			fmt.Println(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		var reactionArr []gin.H
		for reactions.Next() {
			var who_reacted string
			var what_reacted int

			err2 := reactions.Scan(&who_reacted, &what_reacted)
			if err2 != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err2})
				return
			}

			single_reaction := gin.H{
				"who_reacted":  who_reacted,
				"what_reacted": what_reacted,
			}

			reactionArr = append(reactionArr, single_reaction)
		}

		var doIFollow int
		query = `SELECT id FROM followers WHERE following_user_id = $1 AND followed_user_id = $2`
		err4 := db.QueryRow(query, currentUserId, user_id).Scan(&doIFollow)
		if err4 != nil {
			if err4 == sql.ErrNoRows {
				doIFollow = 0
			} else {
				fmt.Println("Error checking follow status:", err4)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while checking follow status"})
				return
			}
		}

		if currentUserId == user_id {
			doIFollow = -1
		}

		createdFromNow := timeAgo(created_at)
		fmt.Println("createdFromNow ", createdFromNow)
		post := gin.H{
			"post_id":        post_id,
			"content":        content,
			"image_url":      image_url,
			"created_at":     createdFromNow,
			"user_id":        user_id,
			"username":       username,
			"avatar_url":     avatar_url,
			"reaction_array": reactionArr,
			"do_I_follow":    doIFollow,
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error after retrieving posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": currentUsername,
		"user_id":  currentUserId,
		"posts":    posts,
	})
}

func timeAgo(createdOn time.Time) string {
	now := time.Now()
	diff := now.Sub(createdOn)

	switch {
	case diff < time.Minute:
		return "now"
	case diff < time.Hour:
		return fmt.Sprintf("%dm", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%dh", int(diff.Hours()))
	case diff < 48*time.Hour:
		return "yesterday"
	case diff < 30*24*time.Hour:
		return fmt.Sprintf("%dd", int(diff.Hours()/24))
	case diff < 52*7*24*time.Hour: // less than a year
		return fmt.Sprintf("%dw", int(diff.Hours()/(24*7)))
	default:
		return fmt.Sprintf("%dy", int(diff.Hours()/(24*365)))
	}
}

// func SavePost(c *gin.Context) {

// 	_, currentUserId := middleware.GetCurrentUser(c)

// 	var postId PostID
// 	if err := c.BindJSON(&postId); err != nil {
// 		fmt.Println(err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	fmt.Println(postId)

// 	db := database.GetDB()

// 	query := `INSERT INTO Saved (userId, postId, created_at) VALUES ($1, $2, $3)`

// 	_, err := db.Exec(query, currentUserId, postId.PostID, time.Now())
// 	if err != nil {
// 		fmt.Println(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating like"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "post saved successfully"})
// }

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
