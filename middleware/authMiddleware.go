package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	Id       int    `json:"id`
	jwt.StandardClaims
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("username", claims.Username)
		c.Set("id", claims.Id)
		c.Next()
	}

}

func GetCurrentUser(c *gin.Context) (string, int) {
	username, usernameExists := c.Get("username")
	id, idExists := c.Get("id")
	if !usernameExists || !idExists {
		return "", 0
	}
	return username.(string), id.(int)
}
