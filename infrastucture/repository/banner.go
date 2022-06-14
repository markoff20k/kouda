package repository

import (
	"github.com/zsmartex/pkg/gpa"
	"gorm.io/gorm"

	"github.com/zsmartex/kouda/internal/models"
)

type bannerRepository struct {
	reader[models.Banner]
	writer[models.Banner]

	repository gpa.Repository
}

type BannerRepository interface {
	Reader[models.Banner]
	Writer[models.Banner]
}

func NewBannerRepository(db *gorm.DB) BannerRepository {
	repo := gpa.New(db, models.Banner{})

	return bannerRepository{
		reader: reader[models.Banner]{
			repository: repo,
		},
		writer: writer[models.Banner]{
			repository: repo,
		},
		repository: repo,
	}
}
