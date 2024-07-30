package models

type User struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Created_at string `json:"created_at"`
}

type UserProfile struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Image    string `json:"image"`
	Headline string `json:"headline`
	Name     string `json:"name"`
}
