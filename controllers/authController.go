package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"social-backend/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key")

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username`
}

type Claims struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
	jwt.StandardClaims
}

func Login(c *gin.Context) {

	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	db := database.GetDB()
	var storedPassword string
	var userId int
	err := db.QueryRow("SELECT password, id FROM users WHERE username = $1", creds.Username).Scan(&storedPassword, &userId)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching user"})
		return
	}

	// Compare the stored hashed password with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	expirationTime := time.Now().Add(8 * time.Hour)

	claims := &Claims{
		Username: creds.Username,
		Id:       userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}

func Register(c *gin.Context) {

	var creds Credentials

	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while hashing password"})
		return
	}

	db := database.GetDB()

	newPassword := string(hashedPassword)
	query := `INSERT INTO users (email, username, password, created_at) VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(query, creds.Email, creds.Username, newPassword, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating user"})
		return
	}

	var userID int
	err3 := db.QueryRow("SELECT id  FROM users WHERE username = $1", creds.Username).Scan(&userID)
	fmt.Println("*****", userID)
	if err3 != nil {
		fmt.Println("Error while fetching user", err3)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetching user"})
		return
	}

	query2 := `INSERT INTO profiles (user_id, image, headline, name, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err2 := db.Exec(query2, userID, "", "", "", time.Now())
	if err2 != nil {
		fmt.Println("Error while inserting", err2)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert profile "})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
