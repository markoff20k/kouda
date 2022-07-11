package entities

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/zsmartex/kouda/internal/models"
)

type Banner struct {
	UUID      uuid.UUID          `json:"uuid,omitempty"`
	URL       string             `json:"url,omitempty"`
	ImageURL  string             `json:"image_url,omitempty"`
	State     models.BannerState `json:"state,omitempty"`
	Type      string             `json:"type,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty"`
}

func BannerToEntity(banner *models.Banner) *Banner {
	return &Banner{
		UUID:      uuid.UUID(banner.UUID),
		URL:       banner.URL,
		State:     banner.State,
		CreatedAt: banner.CreatedAt,
		UpdatedAt: banner.UpdatedAt,
	}
}
