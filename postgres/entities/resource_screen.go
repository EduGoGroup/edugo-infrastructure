package entities

import (
	"time"

	"github.com/google/uuid"
)

// ResourceScreen representa el mapeo recurso-pantalla en ui_config.resource_screens
type ResourceScreen struct {
	ID          uuid.UUID `db:"id"`
	ResourceID  uuid.UUID `db:"resource_id"`
	ResourceKey string    `db:"resource_key"`
	ScreenKey   string    `db:"screen_key"`
	ScreenType  string    `db:"screen_type"`
	IsDefault   bool      `db:"is_default"`
	SortOrder   int       `db:"sort_order"`
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
