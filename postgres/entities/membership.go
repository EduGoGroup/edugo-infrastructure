package entities

import (
	"time"

	"github.com/google/uuid"
)

// Membership representa la tabla 'memberships' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migración: 004_create_memberships.up.sql
// Usada por: api-mobile, api-administracion, worker
type Membership struct {
	ID             uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID  `db:"user_id" gorm:"type:uuid;index;not null"`
	SchoolID       uuid.UUID  `db:"school_id" gorm:"type:uuid;index;not null"`
	AcademicUnitID *uuid.UUID `db:"academic_unit_id" gorm:"type:uuid;index"`
	Role           string     `db:"role" gorm:"not null;type:varchar(50)"`
	Metadata       []byte     `db:"metadata" gorm:"type:jsonb;default:'{}'"`
	IsActive       bool       `db:"is_active" gorm:"not null;default:true"`
	EnrolledAt     time.Time  `db:"enrolled_at" gorm:"not null"`
	WithdrawnAt    *time.Time `db:"withdrawn_at" gorm:"default:null"`
	CreatedAt      time.Time  `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt      time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Membership) TableName() string {
	return "academic.memberships"
}
