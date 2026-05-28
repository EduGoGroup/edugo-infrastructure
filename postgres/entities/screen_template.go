package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ScreenTemplate representa un template de pantalla en ui_config.screen_templates
type ScreenTemplate struct {
	ID          uuid.UUID       `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	Pattern     string          `db:"pattern" gorm:"not null;size:50;index:idx_screen_templates_pattern" validate:"required,max=50"`
	Name        string          `db:"name" gorm:"not null;size:200;uniqueIndex:screen_templates_name_version_key" validate:"required,min=2,max=200"`
	Description *string         `db:"description" gorm:"default:null" validate:"omitempty"`
	Version     int             `db:"version" gorm:"not null;default:1;uniqueIndex:screen_templates_name_version_key" validate:"required"`
	Definition  json.RawMessage `db:"definition" gorm:"type:jsonb;default:'{}';not null"`
	IsActive    bool            `db:"is_active" gorm:"not null;default:true"`
	CreatedBy   *uuid.UUID      `db:"created_by" gorm:"type:uuid;constraint:fk_screen_templates_created_by" validate:"omitempty,uuid"`
	CreatedAt   time.Time       `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt   time.Time       `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (ScreenTemplate) TableName() string {
	return "ui_config.screen_templates"
}
