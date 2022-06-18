package repository

import (
	vezdecodebackend "Gurov-R/vezdecode-backend"

	"github.com/jmoiron/sqlx"
)

type Meme interface {
	GetAll() ([]vezdecodebackend.Meme, error)
	LoadMemes([]vezdecodebackend.Meme) error
	GetMeme(page int) (vezdecodebackend.Meme, error)
	GetPromotedMeme(page int) (vezdecodebackend.Meme, error)
	Like(memeId int, remoteAddr string) error
	Promote(memId int) error
}

type Repository struct {
	Meme
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Meme: NewMemesPostgres(db),
	}
}
