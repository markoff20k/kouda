package models

import (
	"time"

	"github.com/google/uuid"
)

type BannerState string

var (
	BannerStateEnabled  = BannerState("enabled")
	BannerStateDisabled = BannerState("disabled")
)

var BannerStates = []BannerState{
	BannerStateEnabled,
	BannerStateDisabled,
}

type Banner struct {
	UUID      uuid.UUID   `gorm:"primaryKey;type:uuid;not null"`
	Tag       string      `gorm:"type:character varying;not null"`
	URL       string      `gorm:"type:character varying;not null"`
	State     BannerState `gorm:"type:character varying;not null;"`
	Type      string      `gorm:"type:character varying;not null;"`
	CreatedAt time.Time   `gorm:"type:timestamp;not null"`
	UpdatedAt time.Time   `gorm:"type:timestamp;not null"`
}

func (b Banner) TableName() string {
	return "banners"
}
