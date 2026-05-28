package entities

import (
	"time"

	"github.com/google/uuid"
)

// SubjectOfferingEnrollment representa la tabla 'subject_offering_enrollments'
// en PostgreSQL.
//
// Tabla de union entre una sesion de materia (SubjectOffering) y la membresia
// del alumno inscrito. Reemplaza el sentido "alumno-cursa-materia" que antes
// vivia en academic.membership_subjects (ADR 0009 / plan 010 N1.7).
//
// PK compuesta (offering_id, student_membership_id): una inscripcion por
// (sesion, alumno). Las FKs se materializan en
// migrations/sql/post_gorm.sql (GORM no crea FKs desde el tag `constraint:`
// sin campo de relacion; mismo caso que academic.membership_subjects).
type SubjectOfferingEnrollment struct {
	OfferingID          uuid.UUID `db:"offering_id" gorm:"type:uuid;primaryKey;constraint:subject_offering_enrollments_offering_fkey,OnDelete:CASCADE;index:idx_subject_offering_enrollments_offering" validate:"required,uuid"`
	StudentMembershipID uuid.UUID `db:"student_membership_id" gorm:"type:uuid;primaryKey;constraint:subject_offering_enrollments_student_fkey,OnDelete:CASCADE;index:idx_subject_offering_enrollments_student" validate:"required,uuid"`
	EnrolledAt          time.Time `db:"enrolled_at" gorm:"not null;autoCreateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (SubjectOfferingEnrollment) TableName() string {
	return "academic.subject_offering_enrollments"
}
