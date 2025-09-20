package models

var (
	// General Errors
	ErrInternalServer = ApiError{Code: "INTERNAL_SERVER_ERROR", Message: "Something went wrong, please try again later", StatusCode: 500}
	ErrUnauthorized   = ApiError{Code: "UNAUTHORIZED", Message: "You are not authorized to perform this action", StatusCode: 401}
	ErrInvalidInput   = ApiError{Code: "INVALID_INPUT", Message: "The input provided is invalid", StatusCode: 400}
	ErrNotFound       = ApiError{Code: "NOT_FOUND", Message: "The requested resource was not found", StatusCode: 404}

	// Post Errors
	ErrPostNotFound     = ApiError{Code: "POST_NOT_FOUND", Message: "The requested post was not found", StatusCode: 404}
	ErrPostCreateFailed = ApiError{Code: "POST_CREATE_FAILED", Message: "Failed to create post", StatusCode: 500}

	// User Errors
	ErrUserNotFound      = ApiError{Code: "USER_NOT_FOUND", Message: "User does not exist", StatusCode: 404}
	ErrFollowCheckFailed = ApiError{Code: "FOLLOW_CHECK_FAILED", Message: "Error while checking follow status", StatusCode: 500}
)
