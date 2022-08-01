package usecases

import (
	"gorm.io/gorm"

	"github.com/zsmartex/pkg/v2/repository"

	"github.com/zsmartex/kouda/internal/models"
)

type memberUsecase struct {
	usecase[models.Member]
}

type MemberUsecase interface {
	Usecase[models.Member]

	WithTrx(trxHandle *gorm.DB) MemberUsecase
}

func NewMemberUsecase(db *gorm.DB) MemberUsecase {
	return memberUsecase{
		usecase: usecase[models.Member]{
			repository: repository.New(db, models.Member{}),
		},
	}
}

func (u memberUsecase) WithTrx(trxHandle *gorm.DB) MemberUsecase {
	u.repository = u.repository.WithTrx(trxHandle)

	return u
}
