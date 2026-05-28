package entities

import (
	"time"

	"github.com/google/uuid"
)

// Unit representa una unidad organizacional simplificada.
// DEPRECATED: Usar AcademicUnit para nuevas funcionalidades.
// Se mantiene para compatibilidad con edugo-api-administracion.
type Unit struct {
	ID           uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	SchoolID     uuid.UUID  `db:"school_id" gorm:"type:uuid;index;not null;constraint:academic_units_school_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	ParentUnitID *uuid.UUID `db:"parent_unit_id" gorm:"type:uuid;index;constraint:academic_units_parent_fkey,OnDelete:SET NULL" validate:"omitempty,uuid"`
	Name         string     `db:"name" gorm:"not null" validate:"required,min=2,max=255"`
	Description  *string    `db:"description" gorm:"default:null" validate:"omitempty"`
	IsActive     bool       `db:"is_active" gorm:"not null;default:true"`
	CreatedAt    time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt    time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Unit) TableName() string {
	return "academic.academic_units"
}
