package entities

import (
	"time"

	"github.com/google/uuid"
)

// Permission representa un permiso del sistema RBAC
type Permission struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	DisplayName string    `db:"display_name"`
	Description *string   `db:"description"`
	Resource    string    `db:"resource"`
	Action      string    `db:"action"`
	Scope       string    `db:"scope"` // 'system', 'school', 'unit'
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
