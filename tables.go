package vezdecodebackend

type Meme struct {
	Id         int    `json:"id" db:"id"`
	VkId       int    `json:"vk_id" db:"vk_id"`
	ImageUrl   string `json:"image_url" db:"image_url"`
	Timestamp  int    `json:"timestamp" db:"timestamp_"`
	LikesCount int    `json:"likes_count" db:"likes_count"`
	Promoted   bool   `json:"promoted" db:"promoted"`
}

type Like struct {
	Id     int    `json:"id" db:"id"`
	MemeId int    `json:"meme_id" db:"meme_id"`
	UserIp string `json:"user_ip" db:"user_ip"`
}
