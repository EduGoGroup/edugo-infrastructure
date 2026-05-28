package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AcademicUnit representa la tabla 'academic_units' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD con soporte para jerarquía.
//
// Migración: 003_create_academic_units.up.sql
// Usada por: api-administracion (jerarquía), api-mobile (plano), worker
type AcademicUnit struct {
	ID           uuid.UUID       `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	ParentUnitID *uuid.UUID      `db:"parent_unit_id" gorm:"type:uuid;index;constraint:academic_units_parent_fkey,OnDelete:SET NULL;check:academic_units_no_self_reference,id <> parent_unit_id" validate:"omitempty,uuid"`
	SchoolID     uuid.UUID       `db:"school_id" gorm:"type:uuid;index;not null;constraint:academic_units_school_fkey,OnDelete:CASCADE;uniqueIndex:academic_units_unique_code" validate:"required,uuid"`
	Name         string          `db:"name" gorm:"not null;size:255" validate:"required,min=2,max=255"`
	Code         string          `db:"code" gorm:"not null;size:50;uniqueIndex:academic_units_unique_code" validate:"required,min=2,max=50"`
	Type         string          `db:"type" gorm:"not null;type:varchar(50);check:academic_units_type_check,type IN ('school','grade','class','section','club','department')" validate:"required,oneof=school grade class section club department"`
	Description  *string         `db:"description" gorm:"default:null" validate:"omitempty"`
	Level        *string         `db:"level" gorm:"default:null;size:50" validate:"omitempty"`
	AcademicYear int             `db:"academic_year" gorm:"not null;default:0;uniqueIndex:academic_units_unique_code" validate:"required"`
	Metadata     json.RawMessage `db:"metadata" gorm:"type:jsonb;default:'{}'"`
	IsActive     bool            `db:"is_active" gorm:"not null;default:true"`
	CreatedAt    time.Time       `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt    time.Time       `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
	DeletedAt    gorm.DeletedAt  `db:"deleted_at" gorm:"index" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AcademicUnit) TableName() string {
	return "academic.academic_units"
}
