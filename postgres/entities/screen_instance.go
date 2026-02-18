package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ScreenInstance representa una instancia de pantalla en ui_config.screen_instances
type ScreenInstance struct {
	ID                 uuid.UUID       `db:"id"`
	ScreenKey          string          `db:"screen_key"`
	TemplateID         uuid.UUID       `db:"template_id"`
	Name               string          `db:"name"`
	Description        *string         `db:"description"`
	SlotData           json.RawMessage `db:"slot_data"`
	Actions            json.RawMessage `db:"actions"`
	DataEndpoint       *string         `db:"data_endpoint"`
	DataConfig         json.RawMessage `db:"data_config"`
	Scope              string          `db:"scope"`
	RequiredPermission *string         `db:"required_permission"`
	HandlerKey         *string         `db:"handler_key"`
	IsActive           bool            `db:"is_active"`
	CreatedAt          time.Time       `db:"created_at"`
	UpdatedAt          time.Time       `db:"updated_at"`
}
