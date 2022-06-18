package response

type PhotosResponse struct {
	Response PhotosBody `json:"response"`
}

type PhotosBody struct {
	Count int     `json:"count"`
	Items []Photo `json:"items"`
}
type Photo struct {
	AlbumId    int    `json:"album_id"`
	Date       int    `json:"date"`
	Id         int    `json:"id"`
	OwnerId    int    `json:"owner_id"`
	CanComment int    `json:"can_comment"`
	Text       string `json:"text"`
	UserId     int    `json:"user_id"`
	HasTags    bool   `json:"has_tags"`

	Sizes    []PhotoSize `json:"sizes"`
	Likes    PhotoLikes  `json:"likes"`
	Comments PhotoCounts `json:"comments"`
	Reposts  PhotoCounts `json:"reposts"`
	Tags     PhotoCounts `json:"tags"`
}

type PhotoSize struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Url    string `json:"url"`
	Type   string `json:"type"`
}

type PhotoLikes struct {
	Count     int `json:"count"`
	UserLikes int `json:"user_likes"`
}

type PhotoCounts struct {
	Count int `json:"count"`
}

type GroupsResponse struct {
	Response []GroupBody `json:"response"`
}

type GroupBody struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
	IsClosed   int    `json:"is_closed"`
	Type       string `json:"page"`
	Photo50    string `json:"photo_50"`
	Photo100   string `json:"photo_100"`
	Photo200   string `json:"photo_200"`
}
