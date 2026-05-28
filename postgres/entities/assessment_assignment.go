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
	ID             uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	AssessmentID   uuid.UUID  `db:"assessment_id" gorm:"type:uuid;index;not null;constraint:assessment_assignments_assessment_fkey,OnDelete:CASCADE;index:idx_assignment_assessment" validate:"required,uuid"`
	StudentID      *uuid.UUID `db:"student_id" gorm:"type:uuid;constraint:assessment_assignments_student_fkey;check:chk_assignment_target,(student_id IS NOT NULL AND academic_unit_id IS NULL) OR (student_id IS NULL AND academic_unit_id IS NOT NULL)" validate:"omitempty,uuid"`
	AcademicUnitID *uuid.UUID `db:"academic_unit_id" gorm:"type:uuid;constraint:assessment_assignments_unit_fkey" validate:"omitempty,uuid"`
	AssignedBy     uuid.UUID  `db:"assigned_by" gorm:"type:uuid;not null;constraint:assessment_assignments_assigned_by_fkey" validate:"required,uuid"`
	// NOTE: partial index idx_assignment_student (WHERE student_id IS NOT NULL) must be created in post_gorm.sql
	// NOTE: partial index idx_assignment_unit (WHERE academic_unit_id IS NOT NULL) must be created in post_gorm.sql
	// NOTE: partial unique index idx_unique_student_assignment (WHERE student_id IS NOT NULL) must be created in post_gorm.sql
	// NOTE: partial unique index idx_unique_unit_assignment (WHERE academic_unit_id IS NOT NULL) must be created in post_gorm.sql
	AssignedAt     time.Time  `db:"assigned_at" gorm:"not null;autoCreateTime"`
	DueDate        *time.Time `db:"due_date" gorm:"default:null"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AssessmentAssignment) TableName() string {
	return "assessment.assessment_assignments"
}
