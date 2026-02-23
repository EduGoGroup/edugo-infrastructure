package entities

import (
	"time"

	"github.com/google/uuid"
)

// Unit representa una unidad organizacional simplificada.
// DEPRECATED: Usar AcademicUnit para nuevas funcionalidades.
// Se mantiene para compatibilidad con edugo-api-administracion.
type Unit struct {
	ID           uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey"`
	SchoolID     uuid.UUID  `db:"school_id" gorm:"type:uuid;index;not null"`
	ParentUnitID *uuid.UUID `db:"parent_unit_id" gorm:"type:uuid;index"`
	Name         string     `db:"name" gorm:"not null"`
	Description  *string    `db:"description" gorm:"default:null"`
	IsActive     bool       `db:"is_active" gorm:"not null;default:true"`
	CreatedAt    time.Time  `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Unit) TableName() string {
	return "academic.academic_units"
}
