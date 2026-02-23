package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// School representa la tabla 'schools' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migración: 002_create_schools.up.sql
// Usada por: api-mobile, api-administracion, worker
type School struct {
	ID               uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey"`
	Name             string         `db:"name" gorm:"not null"`
	Code             string         `db:"code" gorm:"uniqueIndex;not null"`
	Address          *string        `db:"address" gorm:"default:null"`
	City             *string        `db:"city" gorm:"default:null"`
	Country          string         `db:"country" gorm:"not null"`
	Phone            *string        `db:"phone" gorm:"default:null"`
	Email            *string        `db:"email" gorm:"default:null"`
	Metadata         []byte         `db:"metadata" gorm:"type:jsonb;default:'{}'"`
	IsActive         bool           `db:"is_active" gorm:"not null;default:true"`
	SubscriptionTier string         `db:"subscription_tier" gorm:"not null;type:varchar(50)"`
	MaxTeachers      int            `db:"max_teachers" gorm:"not null;default:0"`
	MaxStudents      int            `db:"max_students" gorm:"not null;default:0"`
	CreatedAt        time.Time      `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt        time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `db:"deleted_at" gorm:"index"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (School) TableName() string {
	return "academic.schools"
}
