package entities

import (
	"time"

	"github.com/google/uuid"
)

// Role representa un rol del sistema RBAC
type Role struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	DisplayName string    `db:"display_name"`
	Description *string   `db:"description"`
	Scope       string    `db:"scope"` // 'system', 'school', 'unit'
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
