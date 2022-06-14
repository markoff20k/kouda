package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zsmartex/pkg/gpa"
	"gorm.io/gorm"

	"github.com/zsmartex/kouda/log"
)

type Reader[T any] interface {
	First(filters ...gpa.Filter) (*T, error)
	Last(filters ...gpa.Filter) (model *T, err error)
	Find(filters ...gpa.Filter) []*T
}

type Writer[T any] interface {
	DoTrx(opts ...*sql.TxOptions) *gorm.DB
	WithTrx(trxHandle *gorm.DB)
	HandleTrx(tx *gorm.DB, handler func(tx *gorm.DB) error) error
	FirstOrCreate(model interface{}, filters ...gpa.Filter)
	Create(model interface{})
	Updates(model interface{}, value T, filters ...gpa.Filter)
	UpdateColumns(model interface{}, value T, filters ...gpa.Filter)
	Delete(filters ...gpa.Filter)
}

type reader[T any] struct {
	repository gpa.Repository
}

func (r reader[T]) First(filters ...gpa.Filter) (model *T, err error) {
	if err := r.repository.First(context.Background(), &model, filters...); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	} else if err != nil {
		panic(err)
	}

	return
}

func (r reader[T]) Last(filters ...gpa.Filter) (model *T, err error) {
	if err := r.repository.Last(context.Background(), &model, filters...); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	} else if err != nil {
		panic(err)
	}

	return
}

func (r reader[T]) Find(filters ...gpa.Filter) (models []*T) {
	if err := r.repository.Find(context.Background(), &models, filters...); err != nil {
		panic(err)
	}

	return
}

type writer[T any] struct {
	repository gpa.Repository
}

func (w writer[T]) DoTrx(opts ...*sql.TxOptions) *gorm.DB {
	return w.repository.Begin(opts...)
}

func (w writer[T]) HandleTrx(tx *gorm.DB, handler func(tx *gorm.DB) error) error {
	if err := handler(tx); err != nil {
		if err := tx.Rollback().Error; err != nil {
			return err
		}

		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (w writer[T]) WithTrx(trxHandle *gorm.DB) {
	if trxHandle == nil {
		log.Debug("Transaction session not found")
		return
	}

	w.repository.DB = trxHandle
}

func (w writer[T]) FirstOrCreate(model interface{}, filters ...gpa.Filter) {
	if err := w.repository.FirstOrCreate(context.Background(), model, filters...); err != nil {
		panic(err)
	}
}

func (w writer[T]) Create(model interface{}) {
	if err := w.repository.Create(context.Background(), model); err != nil {
		panic(err)
	}
}

func (w writer[T]) Updates(model interface{}, value T, filters ...gpa.Filter) {
	if err := w.repository.Updates(context.Background(), model, value, filters...); err != nil {
		panic(err)
	}
}

func (w writer[T]) UpdateColumns(model interface{}, value T, filters ...gpa.Filter) {
	if err := w.repository.UpdateColumns(context.Background(), model, value, filters...); err != nil {
		panic(err)
	}
}

func (w writer[T]) Delete(filters ...gpa.Filter) {
	// if err := w.repository.Delete(context.Background(), T{}, filters...); err != nil {
	// 	panic(err)
	// }
}
