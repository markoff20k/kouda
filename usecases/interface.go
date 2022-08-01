package usecases

import (
	"context"
	"errors"

	"github.com/zsmartex/pkg/v2/gpa"
	"github.com/zsmartex/pkg/v2/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Usecase[V schema.Tabler] interface {
	Count(context context.Context, filters ...gpa.Filter) int
	First(context context.Context, filters ...gpa.Filter) (*V, error)
	Find(context context.Context, filters ...gpa.Filter) []*V
	Transaction(handler func(tx *gorm.DB) error) error
	FirstOrCreate(context context.Context, model interface{}, filters ...gpa.Filter)
	Create(context context.Context, model interface{})
	Updates(context context.Context, model interface{}, value V, filters ...gpa.Filter)
	UpdateColumns(context context.Context, model interface{}, value V, filters ...gpa.Filter)
	Delete(context context.Context, model interface{}, filters ...gpa.Filter)
}

type usecase[V schema.Tabler] struct {
	repository repository.Repository[V]
}

func (u usecase[V]) Count(context context.Context, filters ...gpa.Filter) int {
	if count, err := u.repository.Count(context, filters...); err != nil {
		panic(err)
	} else {
		return count
	}
}

func (u usecase[V]) First(context context.Context, filters ...gpa.Filter) (model *V, err error) {
	if err := u.repository.First(context, &model, filters...); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	} else if err != nil {
		panic(err)
	}

	return
}

func (u usecase[V]) Find(context context.Context, filters ...gpa.Filter) (models []*V) {
	if err := u.repository.Find(context, &models, filters...); err != nil {
		panic(err)
	}

	return
}

func (u usecase[V]) Transaction(handler func(tx *gorm.DB) error) error {
	return u.repository.Transaction(handler)
}

func (u usecase[V]) FirstOrCreate(context context.Context, model interface{}, filters ...gpa.Filter) {
	if err := u.repository.FirstOrCreate(context, model, filters...); err != nil {
		panic(err)
	}
}

func (u usecase[V]) Create(context context.Context, model interface{}) {
	if err := u.repository.Create(context, model); err != nil {
		panic(err)
	}
}

func (u usecase[V]) Updates(context context.Context, model interface{}, value V, filters ...gpa.Filter) {
	if err := u.repository.Updates(context, model, value, filters...); err != nil {
		panic(err)
	}
}

func (u usecase[V]) UpdateColumns(context context.Context, model interface{}, value V, filters ...gpa.Filter) {
	if err := u.repository.UpdateColumns(context, model, value, filters...); err != nil {
		panic(err)
	}
}

func (u usecase[V]) Delete(context context.Context, model interface{}, filters ...gpa.Filter) {
	if err := u.repository.Delete(context, model, filters...); err != nil {
		panic(err)
	}
}
