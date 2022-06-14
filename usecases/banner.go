package usecases

import (
	"github.com/zsmartex/kouda/infrastucture/repository"
	"github.com/zsmartex/kouda/internal/models"
)

type bannerUsecase struct {
	reader[repository.Reader[models.Banner], models.Banner]
	writer[repository.Writer[models.Banner], models.Banner]

	repository repository.BannerRepository
}

type BannerUsecase interface {
	Reader[models.Banner]
	Writer[models.Banner]
}

func NewBannerUsecase(repo repository.BannerRepository) BannerUsecase {
	return bannerUsecase{
		reader: reader[repository.Reader[models.Banner], models.Banner]{
			repository: repo,
		},
		writer: writer[repository.Writer[models.Banner], models.Banner]{
			repository: repo,
		},
		repository: repo,
	}
}
