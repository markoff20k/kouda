package models

import (
	"time"
)

type IconState string

var (
	IconStateEnabled  = IconState("enabled")
	IconStateDisabled = IconState("disabled")
)

var IconStates = []IconState{
	IconStateEnabled,
	IconStateDisabled,
}

type Icon struct {
	ID        int64     `gorm:"primaryKey;not null"`
	Code      string    `gorm:"type:character varying(10);not null;uniqueIndex:index_icons_on_code"`
	State     IconState `gorm:"type:character varying(32);not null;"`
	CreatedAt time.Time `gorm:"type:timestamp;not null"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null"`
}

func (b Icon) TableName() string {
	return "icons"
}
