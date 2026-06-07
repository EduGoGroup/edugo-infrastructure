package entities

import (
	"encoding/json"
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
	ID               uuid.UUID       `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	Name             string          `db:"name" gorm:"not null;size:255" validate:"required,min=2,max=255"`
	Code             string          `db:"code" gorm:"uniqueIndex:schools_code_unique;not null;size:50" validate:"required,min=2,max=50"`
	Address          *string         `db:"address" gorm:"default:null" validate:"omitempty"`
	City             *string         `db:"city" gorm:"default:null;size:100" validate:"omitempty"`
	Country          string          `db:"country" gorm:"not null;size:100;default:'Chile'" validate:"required,min=2,max=100"`
	Phone            *string         `db:"phone" gorm:"default:null;size:50" validate:"omitempty"`
	Email            *string         `db:"email" gorm:"default:null;size:255" validate:"omitempty,email"`
	ConceptTypeID    *uuid.UUID      `db:"concept_type_id" gorm:"type:uuid;constraint:fk_schools_concept_type,OnDelete:SET NULL" validate:"omitempty,uuid"`
	Metadata         json.RawMessage `db:"metadata" gorm:"type:jsonb;default:'{}'"`
	IsActive         bool            `db:"is_active" gorm:"not null;default:true"`
	SubscriptionTier string          `db:"subscription_tier" gorm:"not null;type:varchar(50);check:schools_subscription_tier_check,subscription_tier IN ('free','basic','premium','enterprise')" validate:"required,oneof=free basic premium enterprise"`
	// GradeProfile es el perfil de notas de la escuela (N4 / ADR 0020): 'basic'
	// (nota unica por sesion) o 'detailed' (componentes/grade_item). Gateado por
	// permisos en el FE; el CHECK inline vive en el tag GORM (mismo patron que
	// subscription_tier, misma tabla).
	GradeProfile string `db:"grade_profile" gorm:"not null;type:varchar(20);default:'basic';check:schools_grade_profile_check,grade_profile IN ('basic','detailed')" validate:"required,oneof=basic detailed"`
	MaxTeachers      int             `db:"max_teachers" gorm:"not null;default:0" validate:"required"`
	MaxStudents      int             `db:"max_students" gorm:"not null;default:0" validate:"required"`
	CreatedAt        time.Time       `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt        time.Time       `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
	DeletedAt        gorm.DeletedAt  `db:"deleted_at" gorm:"index" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (School) TableName() string {
	return "academic.schools"
}
