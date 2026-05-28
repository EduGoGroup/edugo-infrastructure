package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ScreenInstance representa una instancia de pantalla en ui_config.screen_instances
type ScreenInstance struct {
	ID                 uuid.UUID       `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	ScreenKey          string          `db:"screen_key" gorm:"uniqueIndex;not null;size:100" validate:"required,max=100"`
	TemplateID         uuid.UUID       `db:"template_id" gorm:"type:uuid;index;not null;constraint:fk_screen_instances_template" validate:"required,uuid"`
	Name               string          `db:"name" gorm:"not null;size:200" validate:"required,min=2,max=200"`
	Description        *string         `db:"description" gorm:"default:null" validate:"omitempty"`
	SlotData           json.RawMessage `db:"slot_data" gorm:"type:jsonb;default:'{}';not null"`
	Scope              string          `db:"scope" gorm:"not null;type:varchar(20);default:'school'" validate:"required"`
	RequiredPermission *string         `db:"required_permission" gorm:"default:null;size:100" validate:"omitempty"`
	HandlerKey         *string         `db:"handler_key" gorm:"default:null;size:100" validate:"omitempty"`
	IsActive           bool            `db:"is_active" gorm:"not null;default:true"`
	CreatedBy          *uuid.UUID      `db:"created_by" gorm:"type:uuid;constraint:fk_screen_instances_created_by" validate:"omitempty,uuid"`
	CreatedAt          time.Time       `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt          time.Time       `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (ScreenInstance) TableName() string {
	return "ui_config.screen_instances"
}
