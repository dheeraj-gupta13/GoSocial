package models

type User struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Created_at string `json:"created_at"`
}
