package models

// type Post struct {
// 	Post_id    int    `json:"post_id"`
// 	User_id    int    `json:"userId"`
// 	Image_url  string `json:"image_url"`
// 	Content    string `json:"content"`
// 	Location   string `json:"location"`
// 	Created_on string `json:"created_on"`
// }

type Post struct {
	Content   string `json:"content"`
	Image_url string `json:"image_url"`
}

type PostReactions struct {
	Reaction_id   string `json:"reaction_id"`
	Post_id       int    `json:"post_id"`
	User_id       int    `json:"user_id"`
	Reaction_type int    `json:"reaction_type"`
	Created_on    string `json:"created_on"`
}

type PostComment struct {
	Comment_id int    `json:"comment_id"`
	Post_id    int    `json:"post_id"`
	User_id    int    `json:"user_id"`
	Comment    string `json:"comment"`
	Created_on string `json:"created_on"`
}

type SavedPost struct {
	Saved_id   int    `json:"saved_id"`
	Post_id    int    `json:"post_id"`
	User_id    int    `json:"user_id"`
	Created_on string `json:"created_on"`
}
