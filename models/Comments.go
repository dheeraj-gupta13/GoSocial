package models

type UserComment struct {
	Post_id int    `json:"post_id"`
	Comment string `json:"comment"`
}

type Comment struct {
	CommentId int    `json:"comment_id"`
	Comment   string `json:"comment"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	AvatarURL string `json:"avatar_url"`
}
