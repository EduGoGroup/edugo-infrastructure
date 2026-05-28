package entities

import (
	"time"

	"github.com/google/uuid"
)

// SchoolConcept representa la tabla 'school_concepts' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD definido en:
// - postgres/migrations/structure/037_academic_school_concepts.sql
//
// Representa los términos activos configurados por institución.
type SchoolConcept struct {
	ID        uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	SchoolID  uuid.UUID `db:"school_id" gorm:"type:uuid;not null;index;constraint:fk_school_concepts_school,OnDelete:CASCADE;uniqueIndex:school_concepts_school_key_unique" validate:"required,uuid"`
	TermKey   string    `db:"term_key" gorm:"not null;size:100;uniqueIndex:school_concepts_school_key_unique" validate:"required,max=100"`
	TermValue string    `db:"term_value" gorm:"not null;size:200" validate:"required,max=200"`
	Category  string    `db:"category" gorm:"not null;default:general;size:50" validate:"required,max=50"`
	SortOrder int       `db:"sort_order" gorm:"not null;default:0" validate:"required"`
	CreatedAt time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (SchoolConcept) TableName() string {
	return "academic.school_concepts"
}
