package entities

import (
	"time"

	"github.com/google/uuid"
)

// AssessmentMaterial representa la tabla N:N entre assessments y materiales.
// Migracion: 053_assessment_materials.sql
type AssessmentMaterial struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id" validate:"required,uuid"`
	AssessmentID uuid.UUID `gorm:"type:uuid;not null;constraint:assessment_materials_assessment_fk,OnDelete:CASCADE;uniqueIndex:assessment_materials_unique;index:idx_assessment_materials_assessment" json:"assessment_id" validate:"required,uuid"`
	MaterialID   uuid.UUID `gorm:"type:uuid;not null;constraint:assessment_materials_material_fk,OnDelete:CASCADE;uniqueIndex:assessment_materials_unique;index:idx_assessment_materials_material" json:"material_id" validate:"required,uuid"`
	SortOrder    int       `gorm:"not null;default:0" json:"sort_order" validate:"required"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AssessmentMaterial) TableName() string {
	return "assessment.assessment_materials"
}
