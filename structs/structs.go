package structs

type User struct {
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	DisplayName     string `json:"display_name"`
	ID              string `json:"id"`
	Login           string `json:"login"`
	OfflineImageURL string `json:"offline_image_url"`
	ProfileImageURL string `json:"profile_image_url"`
	Type            string `json:"type"`
	ViewCount       int64  `json:"view_count"`
}

type Users struct {
	Data []User
}

type Follower struct {
	FollowedAt string `json:"followed_at"`
	FromID     string `json:"from_id"`
	FromName   string `json:"from_name"`
	ToID       string `json:"to_id"`
	ToName     string `json:"to_name"`
}

type Followers struct {
	Data []Follower `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
	Total int64 `json:"total"`
}