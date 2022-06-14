package usecases

import (
	"database/sql"

	"github.com/zsmartex/pkg/gpa"
	"gorm.io/gorm"

	"github.com/zsmartex/kouda/infrastucture/repository"
)

type Reader[T any] interface {
	First(filters ...gpa.Filter) (*T, error)
	Last(filters ...gpa.Filter) (*T, error)
	Find(filters ...gpa.Filter) []*T
}

type Writer[T any] interface {
	DoTrx(opts ...*sql.TxOptions) *gorm.DB
	HandleTrx(tx *gorm.DB, handler func(tx *gorm.DB) error) error
	FirstOrCreate(model interface{}, filters ...gpa.Filter)
	Create(model interface{})
	Updates(model interface{}, value T, filters ...gpa.Filter)
	UpdateColumns(model interface{}, value T, filters ...gpa.Filter)
	Delete(filters ...gpa.Filter)
}

type reader[R repository.Reader[V], V any] struct {
	repository R
}

func (r reader[R, V]) First(filters ...gpa.Filter) (*V, error) {
	return r.repository.First(filters...)
}

func (r reader[R, V]) Last(filters ...gpa.Filter) (*V, error) {
	return r.repository.Last(filters...)
}

func (r reader[R, V]) Find(filters ...gpa.Filter) []*V {
	return r.repository.Find(filters...)
}

type writer[R repository.Writer[V], V any] struct {
	repository R
}

func (r writer[R, V]) DoTrx(opts ...*sql.TxOptions) *gorm.DB {
	return r.repository.DoTrx(opts...)
}

func (r writer[R, V]) WithTrx(tx *gorm.DB) {
	r.repository.WithTrx(tx)
}

func (r writer[R, V]) HandleTrx(tx *gorm.DB, handler func(tx *gorm.DB) error) error {
	return r.repository.HandleTrx(tx, handler)
}

func (w writer[R, V]) FirstOrCreate(model interface{}, filters ...gpa.Filter) {
	w.repository.FirstOrCreate(model, filters...)
}

func (w writer[R, V]) Create(model interface{}) {
	w.repository.Create(model)
}

func (w writer[R, V]) Updates(model interface{}, value V, filters ...gpa.Filter) {
	w.repository.Updates(model, value, filters...)
}

func (w writer[R, V]) UpdateColumns(model interface{}, value V, filters ...gpa.Filter) {
	w.repository.UpdateColumns(model, value, filters...)
}

func (w writer[R, V]) Delete(filters ...gpa.Filter) {
	w.repository.Delete(filters...)
}
