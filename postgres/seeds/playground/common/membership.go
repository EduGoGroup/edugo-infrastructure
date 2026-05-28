package common

import (
	"encoding/json"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MembershipSpec describe la membership a sembrar. Role es el role-name
// de la membership (NO uuid): "admin", "teacher", "student", etc.
type MembershipSpec struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	SchoolID       uuid.UUID
	AcademicUnitID *uuid.UUID      // puede ser nil
	Role           string          // role-name de la membership ("admin", "teacher", etc.)
	Metadata       json.RawMessage // default `{}` si nil
}

// buildMembership mapea MembershipSpec a entities.Membership aplicando
// defaults (Metadata "{}" si nil, IsActive true, EnrolledAt = now UTC).
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

// SeedMembership inserta la membership aplicando defaults. Idempotente
// por PK.
func SeedMembership(tx *gorm.DB, spec MembershipSpec) error {
	m := buildMembership(spec)
	return OnConflictIgnore(tx, &m)
}
