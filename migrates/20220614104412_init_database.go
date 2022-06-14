package migrates

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

var initDatabase = gormigrate.Migration{
	ID: "20220614104412",
	Migrate: func(db *gorm.DB) error {
		type Banner struct {
			UUID        uuid.UUID `gorm:"primaryKey;type:uuid;not null;index"`
			Name        string    `gorm:"type:character varying;not null"`
			Tag         string    `gorm:"type:character varying;not null"`
			Description string    `gorm:"type:character varying;not null"`
			URL         string    `gorm:"type:character varying;not null"`
			CreatedAt   time.Time `gorm:"type:timestamp;not null"`
			UpdatedAt   time.Time `gorm:"type:timestamp;not null"`
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
