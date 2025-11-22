package entities

import (
	"time"

	"github.com/google/uuid"
)

// AcademicUnit representa la tabla 'academic_units' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD con soporte para jerarquía.
//
// Migración: 003_create_academic_units.up.sql
// Usada por: api-administracion (jerarquía), api-mobile (plano), worker
type AcademicUnit struct {
	ID           uuid.UUID  `db:"id"`
	ParentUnitID *uuid.UUID `db:"parent_unit_id"` // NULL = raíz, soporta jerarquía
	SchoolID     uuid.UUID  `db:"school_id"`
	Name         string     `db:"name"`
	Code         string     `db:"code"`
	Type         string     `db:"type"` // school, grade, class, section, club, department
	Description  *string    `db:"description"`
	Level        *string    `db:"level"`
	AcademicYear int        `db:"academic_year"` // 0 = sin año específico
	Metadata     []byte     `db:"metadata"`      // JSONB stored as []byte
	IsActive     bool       `db:"is_active"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AcademicUnit) TableName() string {
	return "academic_units"
}
