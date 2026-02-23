package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ScreenInstance representa una instancia de pantalla en ui_config.screen_instances
type ScreenInstance struct {
	ID                 uuid.UUID       `db:"id" gorm:"type:uuid;primaryKey"`
	ScreenKey          string          `db:"screen_key" gorm:"uniqueIndex;not null"`
	TemplateID         uuid.UUID       `db:"template_id" gorm:"type:uuid;index;not null"`
	Name               string          `db:"name" gorm:"not null"`
	Description        *string         `db:"description" gorm:"default:null"`
	SlotData           json.RawMessage `db:"slot_data" gorm:"type:jsonb;default:'{}'"`
	Actions            json.RawMessage `db:"actions" gorm:"type:jsonb;default:'[]'"`
	DataEndpoint       *string         `db:"data_endpoint" gorm:"default:null"`
	DataConfig         json.RawMessage `db:"data_config" gorm:"type:jsonb;default:'{}'"`
	Scope              string          `db:"scope" gorm:"not null;type:varchar(50)"`
	RequiredPermission *string         `db:"required_permission" gorm:"default:null"`
	HandlerKey         *string         `db:"handler_key" gorm:"default:null"`
	IsActive           bool            `db:"is_active" gorm:"not null;default:true"`
	CreatedBy          *uuid.UUID      `db:"created_by" gorm:"type:uuid"`
	CreatedAt          time.Time       `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt          time.Time       `db:"updated_at" gorm:"not null;autoUpdateTime"`
}
