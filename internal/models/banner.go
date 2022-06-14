package models

import (
	"time"

	"github.com/google/uuid"
)

type Banner struct {
	UUID      uuid.UUID `gorm:"primaryKey;type:uuid;not null;index;default:gen_random_uuid()"`
	Tag       string    `gorm:"type:character varying;not null"`
	URL       string    `gorm:"type:character varying;not null"`
	CreatedAt time.Time `gorm:"type:timestamp;not null"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null"`
}

func (b Banner) TableName() string {
	return "banners"
}
