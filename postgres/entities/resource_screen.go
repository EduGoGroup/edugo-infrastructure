package entities

import (
	"time"

	"github.com/google/uuid"
)

// ResourceScreen representa el mapeo recurso-pantalla en ui_config.resource_screens
type ResourceScreen struct {
	ID          uuid.UUID `db:"id" gorm:"type:uuid;primaryKey"`
	ResourceID  uuid.UUID `db:"resource_id" gorm:"type:uuid;index;not null"`
	ResourceKey string    `db:"resource_key" gorm:"not null"`
	ScreenKey   string    `db:"screen_key" gorm:"not null"`
	ScreenType  string    `db:"screen_type" gorm:"not null;type:varchar(50)"`
	IsDefault   bool      `db:"is_default" gorm:"not null;default:false"`
	SortOrder   int       `db:"sort_order" gorm:"not null;default:0"`
	IsActive    bool      `db:"is_active" gorm:"not null;default:true"`
	CreatedAt   time.Time `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time `db:"updated_at" gorm:"not null;autoUpdateTime"`
}
