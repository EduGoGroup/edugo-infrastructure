package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ScreenTemplate representa un template de pantalla en ui_config.screen_templates
type ScreenTemplate struct {
	ID          uuid.UUID       `db:"id"`
	Pattern     string          `db:"pattern"`
	Name        string          `db:"name"`
	Description *string         `db:"description"`
	Version     int             `db:"version"`
	Definition  json.RawMessage `db:"definition"`
	IsActive    bool            `db:"is_active"`
	CreatedBy   *uuid.UUID      `db:"created_by"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
}
