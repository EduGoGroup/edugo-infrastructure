package entities

import (
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
	ID           uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey"`
	ParentUnitID *uuid.UUID     `db:"parent_unit_id" gorm:"type:uuid;index"`
	SchoolID     uuid.UUID      `db:"school_id" gorm:"type:uuid;index;not null"`
	Name         string         `db:"name" gorm:"not null"`
	Code         string         `db:"code" gorm:"not null"`
	Type         string         `db:"type" gorm:"not null;type:varchar(50)"`
	Description  *string        `db:"description" gorm:"default:null"`
	Level        *string        `db:"level" gorm:"default:null"`
	AcademicYear int            `db:"academic_year" gorm:"not null;default:0"`
	Metadata     []byte         `db:"metadata" gorm:"type:jsonb;default:'{}'"`
	IsActive     bool           `db:"is_active" gorm:"not null;default:true"`
	CreatedAt    time.Time      `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `db:"deleted_at" gorm:"index"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AcademicUnit) TableName() string {
	return "academic.academic_units"
}
