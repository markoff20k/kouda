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
	ID        int64       `gorm:"primaryKey;not null"`
	UUID      uuid.UUID   `gorm:"type:uuid;not null;uniqueIndex:index_banners_on_uuid"`
	URL       string      `gorm:"type:character varying;not null"`
	State     BannerState `gorm:"type:character varying;not null;"`
	Type      string      `gorm:"type:character varying;not null;"`
	CreatedAt time.Time   `gorm:"type:timestamp;not null"`
	UpdatedAt time.Time   `gorm:"type:timestamp;not null"`
}

func (b Banner) TableName() string {
	return "banners"
}
