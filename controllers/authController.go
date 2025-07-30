package controllers

import (
	"database/sql"
	"net/http"
	"social-backend/database"
	"social-backend/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key")

// Register the user
func Register(c *gin.Context) {

	var creds models.Credentials

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
	query := `INSERT INTO users (email, username, password, created_on) VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(query, creds.Email, creds.Username, newPassword, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating user"})
		return
	}

	var user_id int
	err3 := db.QueryRow("SELECT user_id  FROM users WHERE username = $1", creds.Username).Scan(&user_id)
	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetching user"})
		return
	}

	query2 := `INSERT INTO profile (user_id, avatar_url, background_url, biodata, created_on) VALUES ($1, $2, $3, $4, $5)`
	_, err2 := db.Exec(query2, user_id, "", "", "", time.Now())
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert profile "})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login the user
func Login(c *gin.Context) {

	var creds models.Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	db := database.GetDB()
	var storedPassword string
	var user_id int
	err := db.QueryRow("SELECT password, user_id FROM users WHERE username = $1", creds.Username).Scan(&storedPassword, &user_id)

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

	claims := &models.Claims{
		Username: creds.Username,
		Id:       user_id,
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

func ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"validated": false, "error": "No Authorization  found"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"validated": false, "error": "No  token found"})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"validated": false, "error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"validated": true, "message": "Token is valid"})
}
