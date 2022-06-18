package service

import (
	vezdecodebackend "Gurov-R/vezdecode-backend"
	"Gurov-R/vezdecode-backend/pkg/repository"
	"Gurov-R/vezdecode-backend/pkg/response"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

type MemeService struct {
	repo repository.Meme
}

func NewMemeService(repo repository.Meme) *MemeService {
	return &MemeService{
		repo: repo,
	}
}

func (m *MemeService) GetAllMemes() ([]vezdecodebackend.Meme, error) {
	return m.repo.GetAll()
}

func (m *MemeService) LoadVezdekod() error {
	return m.LoadGroupById(-197700721, "281940823", 0)
}

func (m *MemeService) LoadGroupById(groupId int, albumId string, count int) error {

	baseURL := "https://api.vk.com/method/photos.get?"
	v := url.Values{}
	v.Set("access_token", "vk1.a.O0f9zJ3VCi77juAIiF1ETtReJ0pEuwM6DE5BMprnDhiegVUSHw8PuXPWGkWxRtLDAcEpmxvDQyNGOy2a8XbS3yuzEh1uIeuMb7mAikuRVB75YNjM2QorSmbCZJ2zj7ZqEE6tTlB0FQbnC0BRWHbmLEo4pHwomZsbinKV4bOMca6y0cvNQQF1WpzIb5Acz0fb")
	v.Set("v", "5.131")
	v.Set("album_id", albumId)
	v.Set("owner_id", fmt.Sprint(groupId))
	v.Set("rev", "1")
	v.Set("extended", "1")

	if count > 0 {
		v.Set("count", fmt.Sprint(count))
	}

	url := baseURL + v.Encode()

	resp, err := http.Get(url)

	if err != nil {
		logrus.Fatalf("error in GET request (load vezdekod): %s", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	content := string(body)

	var photosResp response.PhotosResponse

	err = json.Unmarshal([]byte(content), &photosResp)

	if err != nil {
		logrus.Fatalf("error while reading json: %s", err)
	}

	return m.repo.LoadMemes(PhotosResponseToMemes(photosResp))
}

func (m *MemeService) LoadGroup(address string) error {
	baseURL := "https://api.vk.com/method/groups.getById?"
	v := url.Values{}
	v.Set("access_token", "vk1.a.O0f9zJ3VCi77juAIiF1ETtReJ0pEuwM6DE5BMprnDhiegVUSHw8PuXPWGkWxRtLDAcEpmxvDQyNGOy2a8XbS3yuzEh1uIeuMb7mAikuRVB75YNjM2QorSmbCZJ2zj7ZqEE6tTlB0FQbnC0BRWHbmLEo4pHwomZsbinKV4bOMca6y0cvNQQF1WpzIb5Acz0fb")
	v.Set("v", "5.131")

	u, err := url.Parse(address)
	if err != nil {
		logrus.Fatalf("error while parsing url: %s", err)
	}
	v.Set("group_ids", strings.ReplaceAll(u.Path, "/", ""))

	url := baseURL + v.Encode()

	resp, err := http.Get(url)

	if err != nil {
		logrus.Fatalf("error in GET request (load vezdekod): %s", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	content := string(body)

	var groupsResp response.GroupsResponse

	err = json.Unmarshal([]byte(content), &groupsResp)

	if err != nil {
		logrus.Fatalf("error while reading json: %s", err)
	}

	groupBody := groupsResp.Response[0]

	return m.LoadGroupById(-groupBody.Id, "wall", 50)
}

func PhotosResponseToMemes(resp response.PhotosResponse) []vezdecodebackend.Meme {

	var memes []vezdecodebackend.Meme

	for _, item := range resp.Response.Items {
		maxWidth := 0
		url := ""
		for _, size := range item.Sizes {
			if size.Width > maxWidth {
				url = size.Url
				maxWidth = size.Width
			}
		}
		memes = append(memes, vezdecodebackend.Meme{
			Id:         0,
			VkId:       item.Id,
			LikesCount: item.Likes.Count,
			Timestamp:  item.Date,
			ImageUrl:   url,
			Promoted:   false,
		})
	}

	return memes
}

func (m *MemeService) Feed(page int) (vezdecodebackend.Meme, error) {
	val := rand.Float64()

	if val > 0.17 {
		return m.repo.GetMeme(page)
	} else {
		return m.repo.GetPromotedMeme(page)
	}
}

func (m *MemeService) Like(memeId int, remoteAddr string) error {
	return m.repo.Like(memeId, remoteAddr)
}

func (m *MemeService) Promote(memeId int) error {
	return m.repo.Promote(memeId)
}
