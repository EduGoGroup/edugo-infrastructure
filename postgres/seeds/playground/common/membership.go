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

// MembershipSpec describe la membership a sembrar. Role es la key del tipo de
// invitación de la membership (NO uuid): "admin", "teacher", "student", etc.
type MembershipSpec struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	SchoolID       uuid.UUID
	AcademicUnitID *uuid.UUID      // puede ser nil
	Role           string          // key del tipo ("admin", "teacher", etc.)
	Metadata       json.RawMessage // default `{}` si nil
}

// buildMembership mapea MembershipSpec a entities.Membership aplicando
// defaults (Metadata "{}" si nil, IsActive true, EnrolledAt = now UTC).
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
		IsActive:         true,
		EnrolledAt:       time.Now().UTC(),
	}
}

// SeedMembership inserta la membership aplicando defaults. La key de Role se
// resuelve a invitation_type_id vía catalog.ResolveInvitationTypeID (MP-08).
// Idempotente por PK.
func SeedMembership(tx *gorm.DB, spec MembershipSpec) error {
	invitationTypeID, err := catalog.ResolveInvitationTypeID(tx, spec.Role)
	if err != nil {
		return fmt.Errorf("SeedMembership: %w", err)
	}
	m := buildMembership(spec, invitationTypeID)
	return OnConflictIgnore(tx, &m)
}
