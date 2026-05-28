package entities

import (
	"time"

	"github.com/google/uuid"
)

// ConceptType representa la tabla 'concept_types' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD definido en:
// - postgres/migrations/structure/035_academic_concept_types.sql
//
// Representa un tipo de institucion con terminologia predefinida.
type ConceptType struct {
	ID          uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	Name        string    `db:"name" gorm:"not null;size:100" validate:"required,min=2,max=100"`
	Code        string    `db:"code" gorm:"uniqueIndex;not null;size:50" validate:"required,min=2,max=50"`
	Description *string   `db:"description" gorm:"default:null" validate:"omitempty"`
	// NOTE: partial index idx_concept_types_active (WHERE is_active = true) must be created in post_gorm.sql
	IsActive    bool      `db:"is_active" gorm:"not null;default:true"`
	CreatedAt   time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt   time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (ConceptType) TableName() string {
	return "academic.concept_types"
}
