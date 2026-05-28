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
	ID             uuid.UUID       `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	UserID         uuid.UUID       `db:"user_id" gorm:"type:uuid;index;not null;constraint:memberships_user_fkey,OnDelete:CASCADE;uniqueIndex:memberships_unique_membership" validate:"required,uuid"`
	SchoolID       uuid.UUID       `db:"school_id" gorm:"type:uuid;index;not null;constraint:memberships_school_fkey,OnDelete:CASCADE;uniqueIndex:memberships_unique_membership" validate:"required,uuid"`
	AcademicUnitID *uuid.UUID      `db:"academic_unit_id" gorm:"type:uuid;index;constraint:memberships_unit_fkey,OnDelete:CASCADE;uniqueIndex:memberships_unique_membership" validate:"omitempty,uuid"`
	Role           string          `db:"role" gorm:"not null;type:varchar(50);check:memberships_role_check,role IN ('teacher','student','guardian','coordinator','admin','assistant');uniqueIndex:memberships_unique_membership" validate:"required,oneof=teacher student guardian coordinator admin assistant"`
	Metadata       json.RawMessage `db:"metadata" gorm:"type:jsonb;default:'{}'"`
	// NOTE: partial index idx_memberships_unit_role_active (WHERE is_active = true) must be created in post_gorm.sql
	IsActive    bool       `db:"is_active" gorm:"not null;default:true"`
	EnrolledAt  time.Time  `db:"enrolled_at" gorm:"not null"`
	WithdrawnAt *time.Time `db:"withdrawn_at" gorm:"default:null"`
	CreatedAt   time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt   time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Membership) TableName() string {
	return "academic.memberships"
}
