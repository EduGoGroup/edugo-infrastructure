package entities

import (
	"time"

	"github.com/google/uuid"
)

// Resource representa un recurso/modulo del sistema para RBAC y menu
type Resource struct {
	ID            uuid.UUID  `db:"id"`
	Key           string     `db:"key"`
	DisplayName   string     `db:"display_name"`
	Description   *string    `db:"description"`
	Icon          *string    `db:"icon"`
	ParentID      *uuid.UUID `db:"parent_id"`
	SortOrder     int        `db:"sort_order"`
	IsMenuVisible bool       `db:"is_menu_visible"`
	Scope         string     `db:"scope"` // 'system', 'school', 'unit'
	IsActive      bool       `db:"is_active"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
}
