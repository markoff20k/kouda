package usecases

import (
	"gorm.io/gorm"

	"github.com/zsmartex/pkg/v2/repository"

	"github.com/zsmartex/kouda/internal/models"
)

type iconUsecase struct {
	usecase[models.Icon]
}

type IconUsecase interface {
	Usecase[models.Icon]

	WithTrx(trxHandle *gorm.DB) IconUsecase
}

func NewIconUsecase(db *gorm.DB) IconUsecase {
	return iconUsecase{
		usecase: usecase[models.Icon]{
			repository: repository.New(db, models.Icon{}),
		},
	}
}

func (u iconUsecase) WithTrx(trxHandle *gorm.DB) IconUsecase {
	u.repository = u.repository.WithTrx(trxHandle)

	return u
}
