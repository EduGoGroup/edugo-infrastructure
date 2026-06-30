package entities

import (
	"time"

	"github.com/google/uuid"
)

// MaterialAssignment representa la tabla 'content.material_assignment' en
// PostgreSQL. Es el PUENTE material → sesion de materia: asigna un material a
// una oferta (subject_offering) para ponerlo a disposicion de sus alumnos.
//
// Calca a assessment.assessment_assignment (el puente analogo de evaluaciones):
// el target es la OFERTA (subject_offering_id), no el alumno; los destinatarios
// se resuelven on-the-fly de academic.subject_offering_enrollments. En lugar de
// due_date, modela una ventana de disponibilidad (available_from / available_until).
//
// FKs cross-schema (subject_offering_id, assigned_by_membership_id), la FK
// same-schema material_id y el UNIQUE (material_id, subject_offering_id) se
// materializan en post_gorm.sql.
type MaterialAssignment struct {
	ID                     uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	MaterialID             uuid.UUID  `db:"material_id" gorm:"type:uuid;not null;index;constraint:material_assignment_material_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	SubjectOfferingID      uuid.UUID  `db:"subject_offering_id" gorm:"type:uuid;not null;index;constraint:material_assignment_offering_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	AssignedByMembershipID uuid.UUID  `db:"assigned_by_membership_id" gorm:"type:uuid;not null;constraint:material_assignment_assigned_by_fkey,OnDelete:RESTRICT" validate:"required,uuid"`
	AvailableFrom          *time.Time `db:"available_from" gorm:"default:null"`
	AvailableUntil         *time.Time `db:"available_until" gorm:"default:null"`
	AssignedAt             time.Time  `db:"assigned_at" gorm:"not null;default:now()"`
	CreatedAt              time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt              time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (MaterialAssignment) TableName() string {
	return "content.material_assignment"
}
