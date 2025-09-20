package models

type ApiError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type ApiResponse[T any] struct {
	Success bool      `json:"success"`
	Error   *ApiError `json:"error,omitempty"`
	Data    T         `json:"data"`
}

type PostsResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Posts    []Post `json:"posts"`
}

type ProfileResponse struct {
	Profile Profile `json:"profile"`
}
