package entities

import (
	"time"

	"github.com/google/uuid"
)

// AssessmentAssignment representa la tabla 'assessment_assignments' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migraciones: 056_assessment_assignments.sql
// Usada por: api-mobile
type AssessmentAssignment struct {
	ID             uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey"`
	AssessmentID   uuid.UUID  `db:"assessment_id" gorm:"type:uuid;index;not null"`
	StudentID      *uuid.UUID `db:"student_id" gorm:"type:uuid"`
	AcademicUnitID *uuid.UUID `db:"academic_unit_id" gorm:"type:uuid"`
	AssignedBy     uuid.UUID  `db:"assigned_by" gorm:"type:uuid;not null"`
	AssignedAt     time.Time  `db:"assigned_at" gorm:"not null;autoCreateTime"`
	DueDate        *time.Time `db:"due_date" gorm:"default:null"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AssessmentAssignment) TableName() string {
	return "assessment.assessment_assignments"
}
