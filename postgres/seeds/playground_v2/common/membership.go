package common

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/catalog"
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
// Role es la key del tipo de invitación de la membership (NO uuid): "admin",
// "teacher", "student", etc. El helper la resuelve a invitation_type_id (MP-08).
type MembershipSpec struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	SchoolID       uuid.UUID
	AcademicUnitID *uuid.UUID      // nil = alcance colegio; set = alcance unidad
	Role           string          // key del tipo ("admin"|"teacher"|"student"|...)
	Metadata       json.RawMessage // default `{}` si nil
}

func buildMembership(spec MembershipSpec, invitationTypeID uuid.UUID) entities.Membership {
	metadata := spec.Metadata
	if metadata == nil {
		metadata = json.RawMessage(`{}`)
	}
	return entities.Membership{
		ID:               spec.ID,
		UserID:           spec.UserID,
		SchoolID:         spec.SchoolID,
		AcademicUnitID:   spec.AcademicUnitID,
		InvitationTypeID: invitationTypeID,
		Metadata:         metadata,
		Status:           "active",
		EnrolledAt:       time.Now().UTC(),
	}
}

// SeedMembership inserta la membership aplicando defaults. El alcance (colegio
// vs unidad) lo decide AcademicUnitID. La key de Role se resuelve a
// invitation_type_id vía catalog.ResolveInvitationTypeID (MP-08). Idempotente
// por id.
func SeedMembership(tx *gorm.DB, spec MembershipSpec) error {
	invitationTypeID, err := catalog.ResolveInvitationTypeID(tx, spec.Role)
	if err != nil {
		return fmt.Errorf("SeedMembership: %w", err)
	}
	m := buildMembership(spec, invitationTypeID)
	return onConflictIgnore(tx, &m)
}
