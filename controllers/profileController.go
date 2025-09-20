package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"social-backend/database"
	"social-backend/middleware"
	"social-backend/models"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	_, currentUserId := middleware.GetCurrentUser(c)
	input_username := c.Query("username")

	db := database.GetDB()

	query_user_id := `SELECT user_id FROM users WHERE username = $1`
	var id int
	err1 := db.QueryRow(query_user_id, input_username).Scan(&id)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting  user_id from username"})
		return
	}

	query := `SELECT p.profile_id, p.user_id, p.avatar_url , p.background_url, p.biodata, p.created_on , u.username
	FROM profile p JOIN users u 
	ON p.user_id = u.user_id
	WHERE u.user_id = $1`

	row := db.QueryRow(query, id)

	var p models.Profile // our resultant profile

	err := row.Scan(&p.ProfileID, &p.UserID, &p.AvatarURL, &p.BackgroundURL, &p.Biodata, &p.CreatedOn, &p.Username)
	if err != nil {
		fmt.Println("*********", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating profile"})
		return
	}

	// followers
	followers_query := `SELECT u.user_id, u.username FROM followers f JOIN users u 
					ON f.following_user_id =  u.user_id
					WHERE followed_user_id = $1`
	rows, err := db.Query(followers_query, id)
	if err != nil {
		fmt.Println("error while retrieving followers ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while retrieving followers"})
		return
	}

	fmt.Print("LOW_KEY", rows)

	var followers []models.UserMini
	for rows.Next() {
		var follower models.UserMini

		err := rows.Scan(&follower.UserID, &follower.Username)
		if err != nil {
			fmt.Println("____", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while scanning username and user_id of follower"})
			return
		}

		followers = append(followers, follower)
	}
	p.Followers = followers

	// following
	following_query := `SELECT u.user_id, u.username FROM followers f JOIN users u 
					ON f.followed_user_id = u.user_id
					WHERE following_user_id = $1`
	rows, err = db.Query(following_query, id)
	if err != nil {
		fmt.Println("error while retrieving followers ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while retrieving following"})
		return
	}

	var followings []models.UserMini
	for rows.Next() {
		var following models.UserMini

		err := rows.Scan(&following.UserID, &following.Username)
		if err != nil {
			fmt.Println("____", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while scanning username of following"})
			return
		}

		followings = append(followings, following)
	}
	p.Followings = followings

	var do_I_follow int
	do_I_follow_query := `SELECT id FROM followers WHERE following_user_id = $1 AND followed_user_id = $2`
	err = db.QueryRow(do_I_follow_query, currentUserId, p.UserID).Scan(&do_I_follow)

	if err == sql.ErrNoRows {
		do_I_follow = 0
	} else if err != nil {
		fmt.Println("Database error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking if I follow the user"})
		return
	} else {
		do_I_follow = 1
	}

	if currentUserId == p.UserID {
		do_I_follow = -1
	}
	p.DoIFollow = do_I_follow

	c.JSON(http.StatusOK, models.ApiResponse[models.ProfileResponse]{
		Success: true,
		Data: models.ProfileResponse{
			Profile: p,
		},
	})
}

type FollowRequest struct {
	UserID int `json:"user_id"`
}

// POST
func Follow(c *gin.Context) {

	_, currentUserId := middleware.GetCurrentUser(c)

	fmt.Println("currentUserId", currentUserId)

	var req FollowRequest
	if err := c.BindJSON(&req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found"})
		return
	}

	db := database.GetDB()

	fmt.Println("req, world ", req)

	query := `INSERT INTO Followers (following_user_id, followed_user_id, created_at) VALUES ($1, $2, $3)`
	fmt.Println("Hello, world", currentUserId, " ", req.UserID, " ", time.Now())
	_, err := db.Exec(query, currentUserId, req.UserID, time.Now())
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while following user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user followed successfully"})

}

func UnFollow(c *gin.Context) {
	_, currentUserId := middleware.GetCurrentUser(c)

	var req FollowRequest
	if err := c.BindJSON(&req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found"})
	}

	db := database.GetDB()

	query := `DELETE FROM followers WHERE following_user_id = $1 AND followed_user_id = $2`
	result, err := db.Exec(query, currentUserId, req.UserID)
	if err != nil {
		fmt.Println("Erro while unfollowing", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while unfollowing user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result, "message": "User unfollowed successfully"})
}

func UpdateProfile(c *gin.Context) {
	currentUsername, currentUserId := middleware.GetCurrentUser(c)
	db := database.GetDB()

	// Get form fields
	newUsername := c.PostForm("username")
	newBio := c.PostForm("biodata")

	fmt.Println("1", newUsername)
	fmt.Println("2", newBio)

	var currentAvatarURL, currentBackgroundURL string
	err := db.QueryRow("SELECT avatar_url, background_url FROM profile WHERE user_id = $1", currentUserId).
		Scan(&currentAvatarURL, &currentBackgroundURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch current profile"})
		return
	}

	// Initialize Cloudinary
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		fmt.Println("3", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Cloudinary"})
		return
	}

	// ---------- Handle Profile Picture ----------
	profilePicURL := currentAvatarURL
	profilePicFile, err := c.FormFile("profile_pic")
	if err == nil {
		openedFile, err := profilePicFile.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open profile picture"})
			return
		}
		defer openedFile.Close()

		uploadRes, err := cld.Upload.Upload(context.Background(), openedFile, uploader.UploadParams{
			PublicID: profilePicFile.Filename,
			Folder:   "profiles",
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Profile picture upload failed"})
			return
		}
		profilePicURL = uploadRes.SecureURL
	}

	fmt.Println("11", profilePicFile)
	fmt.Println("22", profilePicURL)

	// ---------- Handle Background Picture ----------
	backgroundPicURL := currentBackgroundURL
	backgroundFile, err := c.FormFile("background_pic")
	if err == nil {
		openedFile, err := backgroundFile.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open background picture"})
			return
		}
		defer openedFile.Close()

		uploadRes, err := cld.Upload.Upload(context.Background(), openedFile, uploader.UploadParams{
			PublicID: backgroundFile.Filename,
			Folder:   "backgrounds",
		})
		if err != nil {
			fmt.Println("4", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Background image upload failed"})
			return
		}
		backgroundPicURL = uploadRes.SecureURL
	}

	// ---------- Build Update Query ----------
	query := `UPDATE profile 
			  SET 
			      biodata = $1,
			      avatar_url = $2,
			      background_url = $3
			  WHERE user_id = $4`

	_, err = db.Exec(query,
		newBio,
		profilePicURL,
		backgroundPicURL,
		currentUserId,
	)

	if err != nil {
		fmt.Println("4", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Profile updated successfully",
		"username":       currentUsername,
		"user_id":        currentUserId,
		"new_username":   newUsername,
		"biodata":        newBio,
		"profile_pic":    profilePicURL,
		"background_pic": backgroundPicURL,
	})
}
