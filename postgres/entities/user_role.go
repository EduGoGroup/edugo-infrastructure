package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserRole representa la asignación de un rol a un usuario en un contexto específico
type UserRole struct {
	ID             uuid.UUID  `db:"id"`
	UserID         uuid.UUID  `db:"user_id"`
	RoleID         uuid.UUID  `db:"role_id"`
	SchoolID       *uuid.UUID `db:"school_id"`        // NULL = rol a nivel sistema
	AcademicUnitID *uuid.UUID `db:"academic_unit_id"` // NULL = rol a nivel escuela
	IsActive       bool       `db:"is_active"`
	GrantedBy      *uuid.UUID `db:"granted_by"`
	GrantedAt      time.Time  `db:"granted_at"`
	ExpiresAt      *time.Time `db:"expires_at"` // NULL = no expira
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}
