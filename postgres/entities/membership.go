package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Membership representa la tabla 'memberships' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migración: 004_create_memberships.up.sql
// Usada por: api-mobile, api-administracion, worker
type Membership struct {
	ID               uuid.UUID       `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	UserID           uuid.UUID       `db:"user_id" gorm:"type:uuid;index;not null;constraint:memberships_user_fkey,OnDelete:CASCADE;uniqueIndex:memberships_unique_membership" validate:"required,uuid"`
	SchoolID         uuid.UUID       `db:"school_id" gorm:"type:uuid;index;not null;constraint:memberships_school_fkey,OnDelete:CASCADE;uniqueIndex:memberships_unique_membership" validate:"required,uuid"`
	AcademicUnitID   *uuid.UUID      `db:"academic_unit_id" gorm:"type:uuid;index;constraint:memberships_unit_fkey,OnDelete:CASCADE;uniqueIndex:memberships_unique_membership" validate:"omitempty,uuid"`
	InvitationTypeID uuid.UUID       `db:"invitation_type_id" gorm:"column:invitation_type_id;type:uuid;not null;uniqueIndex:memberships_unique_membership" validate:"required,uuid"`
	Metadata         json.RawMessage `db:"metadata" gorm:"type:jsonb;default:'{}'"`
	// NOTE: la FK invitation_type_id -> academic.invitation_types(id) y el índice
	// parcial idx_memberships_unit_invitation_type_active (WHERE status = 'active')
	// se crean en post_gorm.sql (GORM no materializa la FK sin campo de relación).
	// Status es la ÚNICA fuente de verdad del estado de la membresía (plan 024
	// F4·S3·M0): 'active' (default, equivale al antiguo is_active=true), 'pending'
	// (alta sin confirmar) y 'withdrawn' (retiro; WithdrawnAt queda como timestamp
	// informativo). Reemplaza al booleano is_active, que era derivable.
	Status      string     `db:"status" gorm:"not null;type:varchar(12);default:'active';check:memberships_status_check,status IN ('pending','active','withdrawn')" validate:"required,oneof=pending active withdrawn"`
	EnrolledAt  time.Time  `db:"enrolled_at" gorm:"not null"`
	WithdrawnAt *time.Time `db:"withdrawn_at" gorm:"default:null"`
	CreatedAt   time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt   time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Membership) TableName() string {
	return "academic.memberships"
}
