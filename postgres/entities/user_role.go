package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserRole representa la asignación de un rol a un usuario en un contexto específico
type UserRole struct {
	ID             uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID  `db:"user_id" gorm:"type:uuid;index;not null"`
	RoleID         uuid.UUID  `db:"role_id" gorm:"type:uuid;index;not null"`
	SchoolID       *uuid.UUID `db:"school_id" gorm:"type:uuid;index"`
	AcademicUnitID *uuid.UUID `db:"academic_unit_id" gorm:"type:uuid;index"`
	IsActive       bool       `db:"is_active" gorm:"not null;default:true"`
	GrantedBy      *uuid.UUID `db:"granted_by" gorm:"type:uuid"`
	GrantedAt      time.Time  `db:"granted_at" gorm:"not null"`
	ExpiresAt      *time.Time `db:"expires_at" gorm:"default:null"`
	CreatedAt      time.Time  `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt      time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (UserRole) TableName() string {
	return "iam.user_roles"
}
