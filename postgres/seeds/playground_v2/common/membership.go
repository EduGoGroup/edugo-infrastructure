package common

import (
	"encoding/json"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MembershipSpec describe la membership a sembrar. Un único helper cubre los
// dos alcances que usaban upsertMembership / upsertSchoolMembership /
// upsertUnitMembership en las copias:
//
//   - AcademicUnitID nil  → alcance COLEGIO (lo que el school_admin necesita
//     para que su JWT lleve contexto de colegio).
//   - AcademicUnitID set  → alcance UNIDAD (docentes/alumnos en una unidad).
//
// Role es el role-name de la membership (NO uuid): "admin", "teacher",
// "student", etc.
type MembershipSpec struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	SchoolID       uuid.UUID
	AcademicUnitID *uuid.UUID      // nil = alcance colegio; set = alcance unidad
	Role           string          // role-name ("admin"|"teacher"|"student"|...)
	Metadata       json.RawMessage // default `{}` si nil
}

func buildMembership(spec MembershipSpec) entities.Membership {
	metadata := spec.Metadata
	if metadata == nil {
		metadata = json.RawMessage(`{}`)
	}
	return entities.Membership{
		ID:             spec.ID,
		UserID:         spec.UserID,
		SchoolID:       spec.SchoolID,
		AcademicUnitID: spec.AcademicUnitID,
		Role:           spec.Role,
		Metadata:       metadata,
		IsActive:       true,
		EnrolledAt:     time.Now().UTC(),
	}
}

// SeedMembership inserta la membership aplicando defaults. El alcance (colegio
// vs unidad) lo decide AcademicUnitID. Idempotente por id.
func SeedMembership(tx *gorm.DB, spec MembershipSpec) error {
	m := buildMembership(spec)
	return onConflictIgnore(tx, &m)
}
