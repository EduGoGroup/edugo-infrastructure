package common

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GuardianRelationSpec describe el vínculo representante↔alumno a sembrar
// (academic.guardian_relations). El link es por ESCUELA: el índice único cubre
// (guardian_id, student_id, school_id), así que el mismo guardián puede colgar
// del mismo alumno en dos escuelas distintas (link_scope=school). AcademicUnitID
// solo se setea cuando la política de la escuela usa link_scope=school_unit.
//
// Defaults sensatos si el campo viene vacío: RelationshipType="parent",
// Status="active", IsActive=true. IsPrimary se respeta tal cual (zero value =
// false), porque "quién es el primario" es decisión explícita del seed.
type GuardianRelationSpec struct {
	ID               uuid.UUID
	GuardianID       uuid.UUID
	StudentID        uuid.UUID
	SchoolID         uuid.UUID
	AcademicUnitID   *uuid.UUID // nil = link a nivel escuela (link_scope=school)
	RelationshipType string     // default "parent"
	IsPrimary        bool
	IsActive         bool // ignorado para el default; ver SeedGuardianRelation
	Status           string // default "active"
}

func buildGuardianRelation(spec GuardianRelationSpec) entities.GuardianRelation {
	relationshipType := spec.RelationshipType
	if relationshipType == "" {
		relationshipType = "parent"
	}
	status := spec.Status
	if status == "" {
		status = "active"
	}
	return entities.GuardianRelation{
		ID:               spec.ID,
		GuardianID:       spec.GuardianID,
		StudentID:        spec.StudentID,
		SchoolID:         spec.SchoolID,
		AcademicUnitID:   spec.AcademicUnitID,
		RelationshipType: relationshipType,
		IsPrimary:        spec.IsPrimary,
		IsActive:         true,
		Status:           status,
	}
}

// SeedGuardianRelation inserta el vínculo representante↔alumno aplicando
// defaults. IsActive se fuerza a true (un vínculo sembrado nace activo).
// Idempotente por id.
func SeedGuardianRelation(tx *gorm.DB, spec GuardianRelationSpec) error {
	rel := buildGuardianRelation(spec)
	return onConflictIgnore(tx, &rel)
}
