package migrates

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gofrs/uuid"
	"github.com/volatiletech/null/v9"
	"gorm.io/gorm"

	"github.com/zsmartex/kouda/internal/models"
)

var initDatabase = gormigrate.Migration{
	ID: "20220614104412",
	Migrate: func(db *gorm.DB) error {
		type Banner struct {
			ID        int64              `gorm:"primaryKey;not null"`
			UUID      uuid.UUID          `gorm:"type:uuid;not null;uniqueIndex:index_banners_on_uuid"`
			URL       string             `gorm:"type:character varying;not null"`
			State     models.BannerState `gorm:"type:character varying;not null;"`
			Type      string             `gorm:"type:character varying;not null;"`
			CreatedAt time.Time          `gorm:"type:timestamp;not null"`
			UpdatedAt time.Time          `gorm:"type:timestamp;not null"`
		}
		type Member struct {
			ID        int64              `gorm:"primaryKey;autoIncrement"`
			UID       string             `gorm:"type:character varying(32);not null;uniqueIndex:index_members_on_uid"`
			Email     string             `gorm:"type:character varying(255);not null;uniqueIndex:index_members_on_email"`
			Username  null.String        `gorm:"type:character varying(255);uniqueIndex:index_members_on_username"`
			Level     int64              `gorm:"type:integer;not null"`
			Role      string             `gorm:"type:character varying(16);not null"`
			State     models.MemberState `gorm:"type:character varying(16);not null;default:pending"`
			CreatedAt time.Time          `gorm:"type:timestamp;not null"`
			UpdatedAt time.Time          `gorm:"type:timestamp;not null"`
		}
		return db.AutoMigrate(
			Banner{},
			Member{},
		)
	},
	Rollback: func(db *gorm.DB) error {
		return db.Migrator().DropTable(
			"banners",
			"members",
		)
	},
}
