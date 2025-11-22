package entities

import (
	"time"

	"github.com/google/uuid"
)

// Membership representa la tabla 'memberships' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migración: 004_create_memberships.up.sql
// Usada por: api-mobile, api-administracion, worker
type Membership struct {
	ID             uuid.UUID  `db:"id"`
	UserID         uuid.UUID  `db:"user_id"`
	SchoolID       uuid.UUID  `db:"school_id"`
	AcademicUnitID *uuid.UUID `db:"academic_unit_id"` // NULL permitido
	Role           string     `db:"role"`             // teacher, student, guardian, coordinator, admin, assistant
	Metadata       []byte     `db:"metadata"`         // JSONB stored as []byte
	IsActive       bool       `db:"is_active"`
	EnrolledAt     time.Time  `db:"enrolled_at"`
	WithdrawnAt    *time.Time `db:"withdrawn_at"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Membership) TableName() string {
	return "memberships"
}
