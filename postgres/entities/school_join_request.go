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
// El tipo de invitación se referencia por id (invitation_type_id ->
// academic.invitation_types); el CHECK inline de rol que existía antes de MP-08
// se eliminó (la validez del tipo la garantiza la FK). El CHECK de status se
// conserva inline. *uuid con SET NULL para los aprobadores, autoCreate/autoUpdate.
//
// La FK invitation_type_id (same-schema academic) y el índice UNIQUE PARCIAL
// (user_id, school_id, academic_unit_id) WHERE status='pending' (GORM no soporta
// WHERE) viven en sql/post_gorm.sql, junto al trigger set_updated_at.
type SchoolJoinRequest struct {
	ID               uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	UserID           uuid.UUID  `db:"user_id" gorm:"type:uuid;index;not null;constraint:school_join_requests_user_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	SchoolID         uuid.UUID  `db:"school_id" gorm:"type:uuid;not null;constraint:school_join_requests_school_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	AcademicUnitID   uuid.UUID  `db:"academic_unit_id" gorm:"type:uuid;not null;constraint:school_join_requests_unit_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	InvitationTypeID uuid.UUID  `db:"invitation_type_id" gorm:"column:invitation_type_id;type:uuid;not null" validate:"required,uuid"`
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
