package migrates

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"

	"github.com/zsmartex/kouda/internal/models"
)

var initDatabase = gormigrate.Migration{
	ID: "20220614104412",
	Migrate: func(db *gorm.DB) error {
		type Banner struct {
			ID        int64              `gorm:"primaryKey;not null"`
			UUID      uuid.UUID          `gorm:"type:uuid;not null;index:index_banners_on_uuid"`
			Tag       string             `gorm:"type:character varying;not null"`
			URL       string             `gorm:"type:character varying;not null"`
			State     models.BannerState `gorm:"type:character varying;not null;"`
			Type      string             `gorm:"type:character varying;not null;"`
			CreatedAt time.Time          `gorm:"type:timestamp;not null"`
			UpdatedAt time.Time          `gorm:"type:timestamp;not null"`
		}
		return db.AutoMigrate(
			Banner{},
		)
	},
	Rollback: func(db *gorm.DB) error {
		return db.Migrator().DropTable(
			"banners",
		)
	},
}
