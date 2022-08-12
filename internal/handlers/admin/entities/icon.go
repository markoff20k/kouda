package entities

import (
	"time"

	"github.com/zsmartex/kouda/internal/models"
)

type Icon struct {
	Code      string           `json:"code,omitempty"`
	State     models.IconState `json:"state,omitempty"`
	CreatedAt time.Time        `json:"created_at,omitempty"`
	UpdatedAt time.Time        `json:"updated_at,omitempty"`
}

func IconToEntity(icon *models.Icon) *Icon {
	return &Icon{
		Code:      icon.Code,
		State:     icon.State,
		CreatedAt: icon.CreatedAt,
		UpdatedAt: icon.UpdatedAt,
	}
}
