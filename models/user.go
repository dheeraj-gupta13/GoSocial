package models

import "github.com/golang-jwt/jwt"

type User struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Created_on      string `json:"created_on"`
	IsEmailVerified bool   `json:"is_email_verified"`
}

type Profile struct {
	Profile_id     int    `json:"profile_id"`
	User_id        string `json:"user_id"`
	Avatar_url     string `json:"avatar_url"`
	Background_url string `json:"background_url"`
	Biodata        string `json:"biodata"`
	Created_on     string `json:"created_on"`
}

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
