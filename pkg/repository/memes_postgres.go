package repository

import (
	vezdecodebackend "Gurov-R/vezdecode-backend"
	"errors"
	"fmt"
	"math/rand"

	"github.com/jmoiron/sqlx"
)

type MemesPostgres struct {
	db *sqlx.DB
}

func NewMemesPostgres(db *sqlx.DB) *MemesPostgres {
	return &MemesPostgres{db: db}
}

func (m *MemesPostgres) GetAll() ([]vezdecodebackend.Meme, error) {
	var memes []vezdecodebackend.Meme

	query := fmt.Sprintf("SELECT * FROM %s", memesTable)

	err := m.db.Select(&memes, query)

	return m.WithLocalLikesArray(memes), err
}

func (m *MemesPostgres) LoadMemes(memes []vezdecodebackend.Meme) error {
	for _, meme := range memes {
		query := fmt.Sprintf("INSERT INTO %s (vk_id, image_url, timestamp_, likes_count, promoted) values ($1, $2, $3, $4, $5) RETURNING id", memesTable)
		m.db.QueryRow(query, meme.VkId, meme.ImageUrl, meme.Timestamp, meme.LikesCount, meme.Promoted)
	}
	return nil
}

func (m *MemesPostgres) GetMeme(page int) (vezdecodebackend.Meme, error) {
	var countResult []int

	query := fmt.Sprintf("SELECT COUNT(id) FROM %s", memesTable)
	err := m.db.Select(&countResult, query)

	count := countResult[0]

	if err != nil {
		return vezdecodebackend.Meme{}, err
	}

	if page > count-1 {
		page = count - 1
	}
	var memes []vezdecodebackend.Meme

	query = fmt.Sprintf("SELECT * FROM %s ORDER BY timestamp_ DESC LIMIT 1 OFFSET %d", memesTable, page)
	err = m.db.Select(&memes, query)
	return m.WithLocalLikes(memes[0]), err
}

func (m *MemesPostgres) GetPromotedMeme(page int) (vezdecodebackend.Meme, error) {
	var countPromotedResult []int

	query := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE promoted=true", memesTable)
	err := m.db.Select(&countPromotedResult, query)

	countPromoted := countPromotedResult[0]

	if err != nil {
		return vezdecodebackend.Meme{}, err
	}
	if countPromoted > 0 {
		var memes []vezdecodebackend.Meme
		query = fmt.Sprintf("SELECT * FROM %s WHERE promoted=true", memesTable)
		err := m.db.Select(&memes, query)
		if err != nil {
			return vezdecodebackend.Meme{}, err
		}

		return m.WithLocalLikes(memes[rand.Intn(len(memes))]), nil
	} else {
		return m.GetMeme(page)
	}
}

func (m *MemesPostgres) Like(memeId int, remoteAddr string) error {
	var likesResult []int

	query := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE meme_id = %d AND user_ip = '%s'", likesTable, memeId, remoteAddr)
	err := m.db.Select(&likesResult, query)

	if err != nil {
		return err
	}

	if likesResult[0] > 0 {
		return errors.New("cannot")
	}

	query = fmt.Sprintf("INSERT INTO %s (meme_id, user_ip) values ($1, $2)", likesTable)
	m.db.QueryRow(query, memeId, remoteAddr)

	return nil
}

func (m *MemesPostgres) Promote(memeId int) error {
	query := fmt.Sprintf("UPDATE %s SET promoted=true WHERE id=$1", memesTable)
	m.db.QueryRow(query, memeId)

	return nil
}

func (m *MemesPostgres) WithLocalLikes(meme vezdecodebackend.Meme) vezdecodebackend.Meme {
	var likesResult []int

	query := fmt.Sprintf("SELECT count(id) FROM %s WHERE meme_id=%d", likesTable, meme.Id)
	err := m.db.Select(&likesResult, query)

	if err != nil {
		return meme
	}

	meme.LikesCount += likesResult[0]

	return meme
}

func (m *MemesPostgres) WithLocalLikesArray(memes []vezdecodebackend.Meme) []vezdecodebackend.Meme {
	for i := range memes {
		memes[i] = m.WithLocalLikes(memes[i])
	}

	return memes
}
