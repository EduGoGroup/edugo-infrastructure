package common

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EnrollmentSpec describe la inscripción de un alumno (membership) en una sesión
// de materia. SubjectID y PeriodID son copias denormalizadas e inmutables de la
// oferta que respaldan el invariante una-oferta-por-materia-por-período (bug
// 0036): deben coincidir con los de la oferta.
type EnrollmentSpec struct {
	OfferingID          uuid.UUID
	SubjectID           uuid.UUID
	PeriodID            uuid.UUID
	StudentMembershipID uuid.UUID
}

func buildEnrollment(spec EnrollmentSpec) entities.SubjectOfferingEnrollment {
	return entities.SubjectOfferingEnrollment{
		OfferingID:          spec.OfferingID,
		SubjectID:           spec.SubjectID,
		PeriodID:            spec.PeriodID,
		StudentMembershipID: spec.StudentMembershipID,
	}
}

// SeedEnrollment inscribe al alumno en la sesión de materia. La PK es compuesta
// (offering_id, student_membership_id); el OnConflict sobre ambas columnas la
// hace idempotente.
func SeedEnrollment(tx *gorm.DB, spec EnrollmentSpec) error {
	e := buildEnrollment(spec)
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "offering_id"}, {Name: "student_membership_id"}},
		DoNothing: true,
	}).Create(&e).Error
}
