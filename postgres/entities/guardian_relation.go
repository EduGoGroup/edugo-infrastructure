package entities

import (
	"time"

	"github.com/google/uuid"
)

// GuardianRelation representa la tabla 'guardian_relations' en PostgreSQL
// Esta entity es el reflejo exacto del schema de BD definido en:
// - postgres/migrations/015_create_guardian_relations.up.sql
//
// Representa la relación entre un apoderado (guardian) y un estudiante.
// Define el tipo de relación familiar o legal entre ellos.
type GuardianRelation struct {
	ID               uuid.UUID `db:"id"`
	GuardianID       uuid.UUID `db:"guardian_id"`
	StudentID        uuid.UUID `db:"student_id"`
	RelationshipType string    `db:"relationship_type"` // father, mother, grandfather, etc.
	IsActive         bool      `db:"is_active"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
	CreatedBy        string    `db:"created_by"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (GuardianRelation) TableName() string {
	return "guardian_relations"
}
