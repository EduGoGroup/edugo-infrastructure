package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ScreenTemplate representa un template de pantalla en ui_config.screen_templates
type ScreenTemplate struct {
	ID          uuid.UUID       `db:"id" gorm:"type:uuid;primaryKey"`
	Pattern     string          `db:"pattern" gorm:"not null"`
	Name        string          `db:"name" gorm:"not null"`
	Description *string         `db:"description" gorm:"default:null"`
	Version     int             `db:"version" gorm:"not null;default:1"`
	Definition  json.RawMessage `db:"definition" gorm:"type:jsonb;default:'{}'"`
	IsActive    bool            `db:"is_active" gorm:"not null;default:true"`
	CreatedBy   *uuid.UUID      `db:"created_by" gorm:"type:uuid"`
	CreatedAt   time.Time       `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time       `db:"updated_at" gorm:"not null;autoUpdateTime"`
}
