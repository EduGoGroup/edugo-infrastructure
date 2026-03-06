package entities

import (
	"time"

	"github.com/google/uuid"
)

// SchoolConcept representa la tabla 'school_concepts' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD definido en:
// - postgres/migrations/structure/037_academic_school_concepts.sql
//
// Representa los terminos activos por institucion (copia de definitions al crear).
type SchoolConcept struct {
	ID        uuid.UUID `db:"id" gorm:"type:uuid;primaryKey"`
	SchoolID  uuid.UUID `db:"school_id" gorm:"type:uuid;not null"`
	TermKey   string    `db:"term_key" gorm:"not null"`
	TermValue string    `db:"term_value" gorm:"not null"`
	Category  string    `db:"category" gorm:"not null;default:general"`
	CreatedAt time.Time `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (SchoolConcept) TableName() string {
	return "academic.school_concepts"
}
