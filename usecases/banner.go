package usecases

import (
	"gorm.io/gorm"

	"github.com/zsmartex/pkg/v2/repository"

	"github.com/zsmartex/kouda/internal/models"
)

type bannerUsecase struct {
	usecase[models.Banner]
}

type BannerUsecase interface {
	Usecase[models.Banner]

	WithTrx(trxHandle *gorm.DB) BannerUsecase
}

func NewBannerUsecase(db *gorm.DB) BannerUsecase {
	return bannerUsecase{
		usecase: usecase[models.Banner]{
			repository: repository.New(db, models.Banner{}),
		},
	}
}

func (u bannerUsecase) WithTrx(trxHandle *gorm.DB) BannerUsecase {
	u.repository = u.repository.WithTrx(trxHandle)

	return u
}
