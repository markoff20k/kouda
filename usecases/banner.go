package usecases

import (
	"gorm.io/gorm"

	"github.com/zsmartex/kouda/infrastucture/repository"
	"github.com/zsmartex/kouda/internal/models"
)

type bannerUsecase struct {
	usecase[models.Banner]
}

type BannerUsecase interface {
	Usecase[models.Banner]

	WithTrx(trxHandle *gorm.DB) BannerUsecase
}

func NewBannerUsecase(repo repository.Repository[models.Banner]) BannerUsecase {
	return bannerUsecase{
		usecase: usecase[models.Banner]{
			repository: repo,
		},
	}
}

func (u bannerUsecase) WithTrx(trxHandle *gorm.DB) BannerUsecase {
	u.repository = u.repository.WithTrx(trxHandle)

	return u
}
