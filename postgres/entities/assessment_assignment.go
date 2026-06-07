package entities

import (
	"time"

	"github.com/google/uuid"
)

// AssessmentAssignment representa la tabla 'assessment.assessment_assignment' en
// PostgreSQL (N4 / ADR 0019). Es el PUENTE evaluacion → sesion de materia.
//
// Cambio vs viejo: se elimina student_id (→auth.users) XOR academic_unit_id y el
// CHECK chk_assignment_target; el target es la OFERTA (subject_offering_id). La
// entrega NO crea filas por alumno: los destinatarios se resuelven de
// academic.subject_offering_enrollments (arregla A2 por construccion).
//
// FKs cross-schema (assessment_id, subject_offering_id, assigned_by_membership_id)
// y el UNIQUE (assessment_id, subject_offering_id) se materializan en post_gorm.sql.
type AssessmentAssignment struct {
	ID                     uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	AssessmentID           uuid.UUID  `db:"assessment_id" gorm:"type:uuid;not null;index;constraint:assessment_assignment_assessment_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	SubjectOfferingID      uuid.UUID  `db:"subject_offering_id" gorm:"type:uuid;not null;index;constraint:assessment_assignment_offering_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	AssignedByMembershipID uuid.UUID  `db:"assigned_by_membership_id" gorm:"type:uuid;not null;constraint:assessment_assignment_assigned_by_fkey,OnDelete:RESTRICT" validate:"required,uuid"`
	DueDate                *time.Time `db:"due_date" gorm:"default:null"`
	AssignedAt             time.Time  `db:"assigned_at" gorm:"not null;default:now()"`
	CreatedAt              time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt              time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AssessmentAssignment) TableName() string {
	return "assessment.assessment_assignment"
}
