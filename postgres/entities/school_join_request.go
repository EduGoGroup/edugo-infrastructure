package entities

import (
	"time"

	"github.com/google/uuid"
)

// SchoolJoinRequest representa la tabla 'school_join_requests' en PostgreSQL.
//
// Modela una solicitud de ingreso a un colegio/unidad, creada al redimir un
// código de invitación (academic.school_invitations). El flujo usa un DOBLE
// GATE de aprobación: SchoolApproved* (nivel colegio) y UnitApproved* (nivel
// unidad). El estado vive en `status` ('pending','approved','rejected').
//
// Mirroreada del patrón de guardian_relation.go (CHECK inline del enum de rol
// y de status, *uuid con SET NULL para los aprobadores, autoCreate/autoUpdate).
//
// El índice UNIQUE PARCIAL (user_id, school_id, academic_unit_id)
// WHERE status='pending' NO se expresa en tag GORM (no soporta WHERE) → vive
// en sql/post_gorm.sql. El trigger set_updated_at también vive ahí.
type SchoolJoinRequest struct {
	ID               uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	UserID           uuid.UUID  `db:"user_id" gorm:"type:uuid;index;not null;constraint:school_join_requests_user_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	SchoolID         uuid.UUID  `db:"school_id" gorm:"type:uuid;not null;constraint:school_join_requests_school_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	AcademicUnitID   uuid.UUID  `db:"academic_unit_id" gorm:"type:uuid;not null;constraint:school_join_requests_unit_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	Role             string     `db:"role" gorm:"not null;type:varchar(50);check:school_join_requests_role_check,role IN ('teacher','student','guardian','coordinator','admin','assistant')" validate:"required,oneof=teacher student guardian coordinator admin assistant"`
	InvitationID     *uuid.UUID `db:"invitation_id" gorm:"type:uuid;constraint:school_join_requests_invitation_fkey,OnDelete:SET NULL" validate:"omitempty,uuid"`
	Status           string     `db:"status" gorm:"not null;type:varchar(20);default:'pending';index:idx_school_join_requests_status;check:school_join_requests_status_check,status IN ('pending','approved','rejected')" validate:"required,oneof=pending approved rejected"`
	SchoolApprovedBy *uuid.UUID `db:"school_approved_by" gorm:"type:uuid;constraint:school_join_requests_school_approver_fkey,OnDelete:SET NULL" validate:"omitempty,uuid"`
	SchoolApprovedAt *time.Time `db:"school_approved_at" gorm:"default:null" validate:"-"`
	UnitApprovedBy   *uuid.UUID `db:"unit_approved_by" gorm:"type:uuid;constraint:school_join_requests_unit_approver_fkey,OnDelete:SET NULL" validate:"omitempty,uuid"`
	UnitApprovedAt   *time.Time `db:"unit_approved_at" gorm:"default:null" validate:"-"`
	RejectedBy       *uuid.UUID `db:"rejected_by" gorm:"type:uuid;constraint:school_join_requests_rejected_by_fkey,OnDelete:SET NULL" validate:"omitempty,uuid"`
	RejectedAt       *time.Time `db:"rejected_at" gorm:"default:null" validate:"-"`
	RequestedAt      time.Time  `db:"requested_at" gorm:"not null;autoCreateTime" validate:"-"`
	CreatedAt        time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt        time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (SchoolJoinRequest) TableName() string {
	return "academic.school_join_requests"
}
