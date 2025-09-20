package models

type Profile struct {
	ProfileID     int        `json:"profile_id"`
	UserID        int        `json:"user_id"`
	AvatarURL     string     `json:"avatar_url"`
	BackgroundURL string     `json:"background_url"`
	Biodata       string     `json:"biodata"`
	CreatedOn     string     `json:"created_on"`
	Username      string     `json:"username"`
	Followers     []UserMini `json:"followers,omitempty"`
	Followings    []UserMini `json:"followings,omitempty"`
	DoIFollow     int        `json:"do_I_follow"` // -1: self, 0: not following, 1: following
}

// A lightweight user representation for followers/followings
type UserMini struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}
