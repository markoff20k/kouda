package models

import (
	"time"

	"github.com/google/uuid"
)

type Banner struct {
	UUID        uuid.UUID `gorm:"primaryKey;type:uuid;not null;index"`
	Name        string    `gorm:"type:character varying;not null"`
	Tag         string    `gorm:"type:character varying;not null"`
	Description string    `gorm:"type:character varying;not null"`
	URL         string
	CreatedAt   time.Time `gorm:"type:timestamp;not null"`
	UpdatedAt   time.Time `gorm:"type:timestamp;not null"`
}

func (b Banner) TableName() string {
	return "banners"
}
