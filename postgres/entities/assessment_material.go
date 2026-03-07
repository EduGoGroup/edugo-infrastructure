package entities

import (
	"time"

	"github.com/google/uuid"
)

// AssessmentMaterial representa la tabla N:N entre assessments y materiales.
// Migracion: 053_assessment_materials.sql
type AssessmentMaterial struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AssessmentID uuid.UUID `gorm:"type:uuid;not null" json:"assessment_id"`
	MaterialID   uuid.UUID `gorm:"type:uuid;not null" json:"material_id"`
	SortOrder    int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AssessmentMaterial) TableName() string {
	return "assessment.assessment_materials"
}
