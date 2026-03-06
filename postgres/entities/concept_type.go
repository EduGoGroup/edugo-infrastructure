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
	ID          uuid.UUID `db:"id" gorm:"type:uuid;primaryKey"`
	Name        string    `db:"name" gorm:"not null"`
	Code        string    `db:"code" gorm:"uniqueIndex;not null"`
	Description *string   `db:"description" gorm:"default:null"`
	IsActive    bool      `db:"is_active" gorm:"not null;default:true"`
	CreatedAt   time.Time `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (ConceptType) TableName() string {
	return "academic.concept_types"
}
