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
//
// subject_id es una COPIA DENORMALIZADA E INMUTABLE del subject_id de la oferta
// (subject_offerings.subject_id nunca cambia: UpdateSubjectOfferingInput solo
// toca docente/seccion/activo). Esta presente para garantizar a nivel BD el
// invariante "una oferta por materia por alumno" via el uniqueIndex compuesto
// uq_enrollment_student_subject (student_membership_id, subject_id): un alumno
// no puede inscribirse en dos ofertas de la MISMA materia (bug 0036). Tambien
// habilita queries directas por materia sin JOIN a subject_offerings.
type SubjectOfferingEnrollment struct {
	OfferingID          uuid.UUID `db:"offering_id" gorm:"type:uuid;primaryKey;constraint:subject_offering_enrollments_offering_fkey,OnDelete:CASCADE;index:idx_subject_offering_enrollments_offering" validate:"required,uuid"`
	StudentMembershipID uuid.UUID `db:"student_membership_id" gorm:"type:uuid;primaryKey;constraint:subject_offering_enrollments_student_fkey,OnDelete:CASCADE;index:idx_subject_offering_enrollments_student;uniqueIndex:uq_enrollment_student_subject,priority:1" validate:"required,uuid"`
	SubjectID           uuid.UUID `db:"subject_id" gorm:"type:uuid;not null;index;uniqueIndex:uq_enrollment_student_subject,priority:2" validate:"required,uuid"`
	EnrolledAt          time.Time `db:"enrolled_at" gorm:"not null;autoCreateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (SubjectOfferingEnrollment) TableName() string {
	return "academic.subject_offering_enrollments"
}
