package database

import (
	"github.com/zsmartex/pkg/services"
	"gorm.io/gorm"
)

func NewDatabase(host string, port int, user string, password string, dbname string) (*gorm.DB, error) {
	db, err := services.NewDatabase(host, port, user, password, dbname)
	if err != nil {
		return nil, err
	}

	return db, err
}
