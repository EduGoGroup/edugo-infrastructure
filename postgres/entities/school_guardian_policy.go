package entities

import (
	"time"

	"github.com/google/uuid"
)

// SchoolGuardianPolicy = política de representante por escuela (ADR 0026 · DEC-R-D).
// La fila con academic_unit_id NULL es el DEFAULT de la escuela; una fila con unidad
// la sobre-escribe para esa unidad. Defaults = comportamiento de hoy (no restringe).
type SchoolGuardianPolicy struct {
	ID             uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SchoolID       uuid.UUID  `db:"school_id" gorm:"type:uuid;not null;index;constraint:school_guardian_policy_school_fkey,OnDelete:CASCADE"`
	AcademicUnitID *uuid.UUID `db:"academic_unit_id" gorm:"type:uuid;index;constraint:school_guardian_policy_unit_fkey,OnDelete:CASCADE"` // NULL = default de la escuela
	InvitationMode  string `db:"invitation_mode" gorm:"not null;type:varchar(20);default:'manual';check:school_guardian_policy_invmode_check,invitation_mode IN ('none','on_enrollment','manual')"`
	GatesActivation bool   `db:"gates_activation" gorm:"not null;default:false"`
	GatingApprover  string `db:"gating_approver" gorm:"not null;type:varchar(10);default:'any';check:school_guardian_policy_approver_check,gating_approver IN ('any','primary','all')"`
	LinkScope       string `db:"link_scope" gorm:"not null;type:varchar(12);default:'school';check:school_guardian_policy_scope_check,link_scope IN ('school','school_unit')"`
	CreatedAt time.Time `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

func (SchoolGuardianPolicy) TableName() string { return "academic.school_guardian_policy" }
