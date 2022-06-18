package service

import (
	vezdecodebackend "Gurov-R/vezdecode-backend"
	"Gurov-R/vezdecode-backend/pkg/repository"
)

type Meme interface {
	GetAllMemes() ([]vezdecodebackend.Meme, error)
	LoadVezdekod() error
	LoadGroup(address string) error
	Feed(page int) (vezdecodebackend.Meme, error)
	Like(memeId int, remoteAddr string) error
	Promote(memeId int) error
}

type Service struct {
	Meme
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Meme: NewMemeService(repos.Meme),
	}
}
