package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole representa la asignación de un rol a un usuario en un contexto específico
type UserRole struct {
	ID             uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	UserID         uuid.UUID  `db:"user_id" gorm:"type:uuid;index;not null;constraint:fk_user_roles_user,OnDelete:CASCADE;uniqueIndex:uq_user_role_context;index:idx_user_roles_context;index:idx_user_roles_user_active" validate:"required,uuid"`
	RoleID         uuid.UUID  `db:"role_id" gorm:"type:uuid;index;not null;constraint:fk_user_roles_role,OnDelete:CASCADE;uniqueIndex:uq_user_role_context" validate:"required,uuid"`
	SchoolID       *uuid.UUID `db:"school_id" gorm:"type:uuid;index;constraint:fk_user_roles_school,OnDelete:CASCADE;uniqueIndex:uq_user_role_context;index:idx_user_roles_context" validate:"omitempty,uuid"`
	AcademicUnitID *uuid.UUID `db:"academic_unit_id" gorm:"type:uuid;index;constraint:fk_user_roles_unit,OnDelete:CASCADE;uniqueIndex:uq_user_role_context;check:chk_user_roles_unit_requires_school,academic_unit_id IS NULL OR school_id IS NOT NULL;index:idx_user_roles_context" validate:"omitempty,uuid"`
	ScopePattern   *string    `db:"scope_pattern" gorm:"type:text" validate:"omitempty"`
	IsActive       bool       `db:"is_active" gorm:"not null;default:true;index:idx_user_roles_active;index:idx_user_roles_user_active"`
	GrantedBy      *uuid.UUID `db:"granted_by" gorm:"type:uuid;constraint:fk_user_roles_granted_by,OnDelete:SET NULL" validate:"omitempty,uuid"`
	GrantedAt      time.Time  `db:"granted_at" gorm:"not null"`
	// NOTE: partial index idx_user_roles_expires (WHERE expires_at IS NOT NULL) must be created in post_gorm.sql
	// NOTE: partial index idx_user_roles_user_active_covering (WHERE is_active = true) must be created in post_gorm.sql
	ExpiresAt      *time.Time `db:"expires_at" gorm:"default:null"`
	CreatedAt      time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt      time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (UserRole) TableName() string {
	return "iam.user_roles"
}

// BeforeSave calcula scope_pattern desde school_id/academic_unit_id si
// no fue seteado explícitamente. Permite que cualquier inserción de
// UserRole (seeds, APIs) obtenga scope_pattern automáticamente.
func (ur *UserRole) BeforeSave(tx *gorm.DB) error {
	if ur.ScopePattern != nil && *ur.ScopePattern != "" {
		return nil
	}
	var pattern string
	switch {
	case ur.SchoolID == nil && ur.AcademicUnitID == nil:
		pattern = "*"
	case ur.SchoolID != nil && ur.AcademicUnitID == nil:
		pattern = "school:" + ur.SchoolID.String()
	case ur.SchoolID != nil && ur.AcademicUnitID != nil:
		pattern = "school:" + ur.SchoolID.String() + "/unit:" + ur.AcademicUnitID.String()
	default:
		return fmt.Errorf("user_role %s tiene academic_unit_id sin school_id", ur.ID)
	}
	ur.ScopePattern = &pattern
	return nil
}
