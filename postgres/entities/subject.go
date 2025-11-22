package entities

import (
	"time"

	"github.com/google/uuid"
)

// Subject representa la tabla 'subjects' en PostgreSQL
// Esta entity es el reflejo exacto del schema de BD definido en:
// - postgres/migrations/013_create_subjects.up.sql
//
// Representa una materia o asignatura del sistema educativo.
type Subject struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"` // NULL permitido
	Metadata    *string   `db:"metadata"`    // JSONB, NULL permitido
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Subject) TableName() string {
	return "subjects"
}
