package common

import (
	"encoding/json"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OfferingSpec describe una sesión de materia (subject_offering): la materia
// dictada en una unidad y período concretos.
//
//   - SectionLabel nil → sin sección (n1_inscripcion); set → forma parte del
//     índice único natural uq_subject_offerings_natural (n4/n17).
//   - TeacherMembershipID nil → sesión sin docente (teacher_membership_id NULL);
//     set → docente asignado.
//   - Capacity: reservado para el caso universidad; nil en N1.7 (ADR 0009).
type OfferingSpec struct {
	ID                  uuid.UUID
	SchoolID            uuid.UUID
	SubjectID           uuid.UUID
	AcademicUnitID      uuid.UUID
	PeriodID            uuid.UUID
	SectionLabel        *string         // nil = sin sección
	TeacherMembershipID *uuid.UUID      // nil = sesión sin docente
	Capacity            *int            // nil = sin cupo (default N1.7)
	Metadata            json.RawMessage // default `{}` si nil
}

func buildOffering(spec OfferingSpec) entities.SubjectOffering {
	metadata := spec.Metadata
	if metadata == nil {
		metadata = json.RawMessage(`{}`)
	}
	return entities.SubjectOffering{
		ID:                  spec.ID,
		SchoolID:            spec.SchoolID,
		SubjectID:           spec.SubjectID,
		AcademicUnitID:      spec.AcademicUnitID,
		SectionLabel:        spec.SectionLabel,
		PeriodID:            spec.PeriodID,
		TeacherMembershipID: spec.TeacherMembershipID,
		Capacity:            spec.Capacity,
		IsActive:            true,
		Metadata:            metadata,
	}
}

// SeedOffering inserta la sesión de materia aplicando defaults. Idempotente por
// id.
func SeedOffering(tx *gorm.DB, spec OfferingSpec) error {
	o := buildOffering(spec)
	return onConflictIgnore(tx, &o)
}
