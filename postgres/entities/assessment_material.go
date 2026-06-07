package entities

import "github.com/google/uuid"

// AssessmentMaterial representa la tabla N:N 'assessment.assessment_material'
// entre evaluaciones y materiales guia (N4 / ADR 0019).
//
// PK compuesta (assessment_id, material_id): el lector deja de asumir 1:1
// (arregla A4). Las FKs (assessment_id → assessment.assessment,
// material_id → content.materials) se materializan en post_gorm.sql.
type AssessmentMaterial struct {
	AssessmentID uuid.UUID `db:"assessment_id" gorm:"type:uuid;primaryKey" json:"assessment_id" validate:"required,uuid"`
	MaterialID   uuid.UUID `db:"material_id" gorm:"type:uuid;primaryKey" json:"material_id" validate:"required,uuid"`
	SortOrder    int       `db:"sort_order" gorm:"not null;default:0" json:"sort_order"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AssessmentMaterial) TableName() string {
	return "assessment.assessment_material"
}
