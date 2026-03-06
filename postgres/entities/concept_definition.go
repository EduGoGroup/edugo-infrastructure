package entities

import (
	"time"

	"github.com/google/uuid"
)

// ConceptDefinition representa la tabla 'concept_definitions' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD definido en:
// - postgres/migrations/structure/036_academic_concept_definitions.sql
//
// Representa los terminos predeterminados por tipo de institucion (plantilla).
type ConceptDefinition struct {
	ID            uuid.UUID `db:"id" gorm:"type:uuid;primaryKey"`
	ConceptTypeID uuid.UUID `db:"concept_type_id" gorm:"type:uuid;not null"`
	TermKey       string    `db:"term_key" gorm:"not null"`
	TermValue     string    `db:"term_value" gorm:"not null"`
	Category      string    `db:"category" gorm:"not null;default:general"`
	SortOrder     int       `db:"sort_order" gorm:"not null;default:0"`
	CreatedAt     time.Time `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt     time.Time `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (ConceptDefinition) TableName() string {
	return "academic.concept_definitions"
}
