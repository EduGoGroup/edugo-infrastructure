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
	ID            uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	ConceptTypeID uuid.UUID `db:"concept_type_id" gorm:"type:uuid;not null;index;constraint:fk_concept_definitions_type,OnDelete:CASCADE;uniqueIndex:concept_definitions_type_key_unique" validate:"required,uuid"`
	TermKey       string    `db:"term_key" gorm:"not null;size:100;uniqueIndex:concept_definitions_type_key_unique" validate:"required,max=100"`
	TermValue     string    `db:"term_value" gorm:"not null;size:200" validate:"required,max=200"`
	Category      string    `db:"category" gorm:"not null;default:general;size:50" validate:"required,max=50"`
	SortOrder     int       `db:"sort_order" gorm:"not null;default:0" validate:"required"`
	CreatedAt     time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt     time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (ConceptDefinition) TableName() string {
	return "academic.concept_definitions"
}
