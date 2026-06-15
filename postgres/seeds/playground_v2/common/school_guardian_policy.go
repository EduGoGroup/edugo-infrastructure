package common

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SchoolGuardianPolicySpec describe la política de representante de una escuela
// (academic.school_guardian_policy, ADR 0026 · DEC-R-D). La fila con
// AcademicUnitID nil es el DEFAULT de la escuela; una fila con unidad la
// sobre-escribe para esa unidad. Una escuela SIN fila usa los defaults del
// esquema (no restringe), así que solo se siembra cuando se quiere apartar del
// comportamiento por defecto.
//
// Defaults si el campo viene vacío: InvitationMode="manual", GatingApprover="any",
// LinkScope="school". GatesActivation se respeta tal cual (zero value = false).
type SchoolGuardianPolicySpec struct {
	ID              uuid.UUID
	SchoolID        uuid.UUID
	AcademicUnitID  *uuid.UUID // nil = default de la escuela
	InvitationMode  string     // default "manual" ('none'|'on_enrollment'|'manual')
	GatesActivation bool
	GatingApprover  string // default "any" ('any'|'primary'|'all')
	LinkScope       string // default "school" ('school'|'school_unit')
}

func buildSchoolGuardianPolicy(spec SchoolGuardianPolicySpec) entities.SchoolGuardianPolicy {
	invitationMode := spec.InvitationMode
	if invitationMode == "" {
		invitationMode = "manual"
	}
	gatingApprover := spec.GatingApprover
	if gatingApprover == "" {
		gatingApprover = "any"
	}
	linkScope := spec.LinkScope
	if linkScope == "" {
		linkScope = "school"
	}
	return entities.SchoolGuardianPolicy{
		ID:              spec.ID,
		SchoolID:        spec.SchoolID,
		AcademicUnitID:  spec.AcademicUnitID,
		InvitationMode:  invitationMode,
		GatesActivation: spec.GatesActivation,
		GatingApprover:  gatingApprover,
		LinkScope:       linkScope,
	}
}

// SeedSchoolGuardianPolicy inserta la política de representante aplicando
// defaults. Idempotente por id.
func SeedSchoolGuardianPolicy(tx *gorm.DB, spec SchoolGuardianPolicySpec) error {
	policy := buildSchoolGuardianPolicy(spec)
	return onConflictIgnore(tx, &policy)
}
