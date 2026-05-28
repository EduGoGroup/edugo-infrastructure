package entities

import (
	"time"

	"github.com/google/uuid"
)

// SchoolInvitation representa la tabla 'school_invitations' en PostgreSQL.
//
// Modela un código de invitación a un colegio/unidad con un rol predefinido.
// El código se redime para crear una solicitud de ingreso
// (academic.school_join_requests). Mirroreada del patrón de
// guardian_relation.go (CHECK inline del enum de rol, CreatedBy *uuid con
// SET NULL, autoCreate/autoUpdate).
//
// El índice/constraints que GORM no puede expresar (ninguno por ahora) y el
// trigger set_updated_at viven en sql/post_gorm.sql (sección academic).
type SchoolInvitation struct {
	ID             uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	Code           string     `db:"code" gorm:"not null;type:varchar(64);uniqueIndex:school_invitations_code_key" validate:"required,max=64"`
	SchoolID       uuid.UUID  `db:"school_id" gorm:"type:uuid;index;not null;constraint:school_invitations_school_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	AcademicUnitID uuid.UUID  `db:"academic_unit_id" gorm:"type:uuid;index;not null;constraint:school_invitations_unit_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	Role           string     `db:"role" gorm:"not null;type:varchar(50);check:school_invitations_role_check,role IN ('teacher','student','guardian','coordinator','admin','assistant')" validate:"required,oneof=teacher student guardian coordinator admin assistant"`
	Label          *string    `db:"label" gorm:"type:varchar(150);default:null" validate:"omitempty,max=150"`
	CreatedBy      *uuid.UUID `db:"created_by" gorm:"type:uuid;constraint:school_invitations_created_by_fkey,OnDelete:SET NULL" validate:"omitempty,uuid"`
	ExpiresAt      *time.Time `db:"expires_at" gorm:"default:null" validate:"-"`
	MaxUses        *int       `db:"max_uses" gorm:"default:null" validate:"omitempty"`
	UsesCount      int        `db:"uses_count" gorm:"not null;default:0" validate:"-"`
	IsActive       bool       `db:"is_active" gorm:"not null;default:true"`
	CreatedAt      time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt      time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (SchoolInvitation) TableName() string {
	return "academic.school_invitations"
}
