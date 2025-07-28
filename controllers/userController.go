package controllers

import (
	"fmt"
	"net/http"
	"social-backend/database"
	"social-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	db := database.GetDB()

	// Example query to fetch all users
	query := `SELECT email, username, created_on FROM users`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error querying users:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	fmt.Println("CONTINUE")

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Email, &user.Username, &user.Created_on)
		if err != nil {
			fmt.Println("Error scanning user row:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over user rows:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func PostUser(c *gin.Context) {
	var newUser models.User

	// Bind the JSON body to newUser
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	fmt.Printf("connection success", db)

	query := `INSERT INTO users (email, username, password, created_at) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, newUser.Email, newUser.Username, newUser.Password, time.Now())
	if err != nil {
		fmt.Println("Error while inserting", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "email": newUser.Email, "password": newUser.Password, "username": newUser.Username})
}

func GetUserInfo(c *gin.Context) {
	// _, currentUserId := middleware.GetCurrentUser(c)

	// db := database.GetDB()

	// fmt.Println(currentUserId)

	// query := `
	// 		  SELECT u.id, u.email, u.username, p.image, p.headline, p.name
	// 		  FROM users u
	// 		  INNER JOIN profile p
	// 		  ON u.id = p.user_id
	// 		  WHERE u.id = $1
	// 		`

	// rows, err := db.Query(query, currentUserId)
	// if err != nil {
	// 	fmt.Println("Error while inserting", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
	// 	return
	// }

	// var userProfile models.Profile
	// if rows.Next() {
	// 	err = rows.Scan(&userProfile.User_id, &userProfile.Email, &userProfile.Username, &userProfile.Image, &userProfile.Headline, &userProfile.Name)
	// 	if err != nil {
	// 		fmt.Println("Error while scanning", err)
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
	// 		return
	// 	}
	// } else {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"message": "Get User successfully", "row": "userProfile"})
}
