package entities

import (
	"time"

	"github.com/google/uuid"
)

// Unit representa la tabla 'units' en PostgreSQL
// Esta entity es el reflejo exacto del schema de BD definido en:
// - postgres/migrations/014_create_units.up.sql (movida a _deprecated/ - tabla no existe en DB actual)
//
// Representa una unidad organizacional jerárquica (departamento, grado, grupo, etc.)
// Puede tener una unidad padre, permitiendo estructuras jerárquicas.
//
// NOTA: La migración 014 está en _deprecated/ porque la tabla 'units' no existe en la BD actual.
// Esta entidad se mantiene para compatibilidad con edugo-api-administracion hasta su migración.
type Unit struct {
	ID           uuid.UUID  `db:"id"`
	SchoolID     uuid.UUID  `db:"school_id"`
	ParentUnitID *uuid.UUID `db:"parent_unit_id"` // NULL permitido para unidades raíz
	Name         string     `db:"name"`
	Description  *string    `db:"description"` // NULL permitido
	IsActive     bool       `db:"is_active"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Unit) TableName() string {
	return "units"
}
