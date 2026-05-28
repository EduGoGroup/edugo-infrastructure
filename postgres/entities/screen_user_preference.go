package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ScreenUserPreference struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	ScreenKey   string          `gorm:"not null;index:idx_screen_user,unique" validate:"required"`
	UserID      uuid.UUID       `gorm:"type:uuid;not null;index:idx_screen_user,unique;constraint:fk_screen_user_prefs_user" validate:"required,uuid"`
	Preferences json.RawMessage `gorm:"type:jsonb;not null;default:'{}'"`
	CreatedAt   time.Time       `gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt   time.Time       `gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (ScreenUserPreference) TableName() string {
	return "ui_config.screen_user_preferences"
}
