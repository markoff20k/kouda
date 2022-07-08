package usecases

import (
	"errors"

	filters "github.com/zsmartex/pkg/v2/gpa"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/zsmartex/kouda/infrastucture/repository"
)

type Usecase[V schema.Tabler] interface {
	First(filters ...filters.Filter) (*V, error)
	Find(filters ...filters.Filter) []*V
	Transaction(handler func(tx *gorm.DB) error) error
	FirstOrCreate(model interface{}, filters ...filters.Filter)
	Create(model interface{})
	Updates(model interface{}, value V, filters ...filters.Filter)
	UpdateColumns(model interface{}, value V, filters ...filters.Filter)
	Delete(model interface{}, filters ...filters.Filter)
}

type usecase[V schema.Tabler] struct {
	repository repository.Repository[V]
}

func (u usecase[V]) First(filters ...filters.Filter) (model *V, err error) {
	if err := u.repository.First(&model, filters...); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	} else if err != nil {
		panic(err)
	}

	return
}

func (u usecase[V]) Find(filters ...filters.Filter) (models []*V) {
	if err := u.repository.Find(&models, filters...); err != nil {
		panic(err)
	}

	return
}

func (u usecase[V]) Transaction(handler func(tx *gorm.DB) error) error {
	return u.repository.Transaction(handler)
}

func (u usecase[V]) FirstOrCreate(model interface{}, filters ...filters.Filter) {
	if err := u.repository.FirstOrCreate(model, filters...); err != nil {
		panic(err)
	}
}

func (u usecase[V]) Create(model interface{}) {
	if err := u.repository.Create(model); err != nil {
		panic(err)
	}
}

func (u usecase[V]) Updates(model interface{}, value V, filters ...filters.Filter) {
	if err := u.repository.Updates(model, value, filters...); err != nil {
		panic(err)
	}
}

func (u usecase[V]) UpdateColumns(model interface{}, value V, filters ...filters.Filter) {
	if err := u.repository.UpdateColumns(model, value, filters...); err != nil {
		panic(err)
	}
}

func (u usecase[V]) Delete(model interface{}, filters ...filters.Filter) {
	if err := u.repository.Delete(model, filters...); err != nil {
		panic(err)
	}
}
