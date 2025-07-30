package controllers

import (
	"fmt"
	"net/http"
	"social-backend/database"
	"social-backend/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	id := c.Query("user_id")

	db := database.GetDB()
	query := `SELECT p.profile_id, p.user_id, p.avatar_url , p.background_url, p.biodata, p.created_on , u.username
	FROM profile p JOIN users u 
	ON p.user_id = u.user_id
	WHERE u.user_id = $1`

	row := db.QueryRow(query, id)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "error while query profile"})
	// }

	// for rows.Next() {
	var profile_id int
	var user_id int
	var avatar_url string
	var background_url string
	var biodata string
	var created_on string
	var username string

	err := row.Scan(&profile_id, &user_id, &avatar_url, &background_url, &biodata, &created_on, &username)
	if err != nil {
		fmt.Println("*********", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating profile"})
		return
	}

	// followers
	followers_query := `SELECT u.user_id, u.username FROM followers f JOIN users u 
					ON f.following_user_id =  u.user_id
					WHERE following_user_id = $1`
	rows, err := db.Query(followers_query, id)
	if err != nil {
		fmt.Println("error while retrieving followers ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while retrieving followers"})
		return
	}

	fmt.Print("LOW_KEY", rows)

	var followers []gin.H
	for rows.Next() {
		var username string
		var user_id int

		err := rows.Scan(&user_id, &username)
		if err != nil {
			fmt.Println("____", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while scanning username and user_id of follower"})
			return
		}

		follower := gin.H{
			"username": username,
			"user_id":  user_id,
		}

		followers = append(followers, follower)
	}

	// following
	following_query := `SELECT u.user_id, u.username FROM followers f JOIN users u 
					ON f.followed_user_id = u.user_id
					WHERE followed_user_id = $1`
	rows, err = db.Query(following_query, id)
	if err != nil {
		fmt.Println("error while retrieving followers ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while retrieving following"})
		return
	}

	var followings []gin.H
	for rows.Next() {
		var username string
		var user_id int

		err := rows.Scan(&user_id, &username)
		if err != nil {
			fmt.Println("____", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while scanning username of following"})
			return
		}

		following := gin.H{
			"username": username,
			"user_id":  user_id,
		}

		followings = append(followings, following)
	}

	profile := gin.H{
		"profile_id":     profile_id,
		"user_id":        user_id,
		"avatar_url":     avatar_url,
		"background_url": background_url,
		"created_on":     created_on,
		"username":       username,
		"biodata":        biodata,
		"followers":      followers,
		"followings":     followings,
	}

	// profile = append(profile, p)
	// }/

	// if err = rows.Err(); err != nil {
	// 	fmt.Println("*********", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error after retrieving posts"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"profile": profile,
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
