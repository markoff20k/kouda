package repository

import (
	"context"

	"github.com/zsmartex/pkg/v2/gpa"
	"github.com/zsmartex/pkg/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Repository[T schema.Tabler] interface {
	First(model interface{}, filters ...gpa.Filter) error
	Find(models interface{}, filters ...gpa.Filter) error
	WithTrx(trxHandle *gorm.DB) Repository[T]
	Transaction(handler func(tx *gorm.DB) error) error
	FirstOrCreate(model interface{}, filters ...gpa.Filter) error
	Create(model interface{}) error
	Updates(model interface{}, value T, filters ...gpa.Filter) error
	UpdateColumns(model interface{}, value T, filters ...gpa.Filter) error
	Delete(model interface{}, filters ...gpa.Filter) error
}

type repository[T schema.Tabler] struct {
	repository gpa.Repository
}

func New[T schema.Tabler](db *gorm.DB, entity T) Repository[T] {
	return repository[T]{
		repository: gpa.New(db, entity),
	}
}

func (r repository[T]) WithTrx(trxHandle *gorm.DB) Repository[T] {
	r.repository = r.repository.WithTrx(trxHandle)

	return r
}

func (r repository[T]) First(model interface{}, filters ...gpa.Filter) (err error) {
	return r.repository.First(context.Background(), model, filters...)
}

func (r repository[T]) Find(models interface{}, filters ...gpa.Filter) error {
	return r.repository.Find(context.Background(), models, filters...)
}

func (r repository[T]) Transaction(handler func(tx *gorm.DB) error) (err error) {
	tx := r.repository.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()

			err = r.(error)
		}
	}()

	if err = handler(tx); err != nil {
		if err = tx.Rollback().Error; err != nil {
			log.Errorf("failed to rollback transaction: %v", err)
		}

		return
	}

	if err = tx.Commit().Error; err != nil {
		log.Errorf("failed to commit transaction: %v", err)
		return
	}

	return nil
}

func (r repository[T]) FirstOrCreate(model interface{}, filters ...gpa.Filter) error {
	return r.repository.FirstOrCreate(context.Background(), model, filters...)
}

func (r repository[T]) Create(model interface{}) error {
	return r.repository.Create(context.Background(), model)
}

func (r repository[T]) Updates(model interface{}, value T, filters ...gpa.Filter) error {
	return r.repository.Updates(context.Background(), model, value, filters...)
}

func (r repository[T]) UpdateColumns(model interface{}, value T, filters ...gpa.Filter) error {
	return r.repository.UpdateColumns(context.Background(), model, value, filters...)
}

func (r repository[T]) Delete(model interface{}, filters ...gpa.Filter) error {
	return r.repository.Delete(context.Background(), model, filters...)
}
