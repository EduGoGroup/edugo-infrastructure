package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Subject representa la tabla 'subjects' en PostgreSQL
// Esta entity es el reflejo exacto del schema de BD definido en:
// - postgres/migrations/structure/033_academic_subjects.sql
//
// Representa una materia o asignatura del sistema educativo.
type Subject struct {
	ID uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	// SchoolID + Name forman la clave natural de la materia: una materia es
	// catalogo de ESCUELA (ADR 0016), no se repite el mismo nombre dentro de
	// la escuela. El unique compuesto uq_subjects_school_name respalda a nivel
	// BD la validacion logica ExistsByNameInSchool de la API academic.
	SchoolID       uuid.UUID      `db:"school_id" gorm:"type:uuid;index;not null;constraint:subjects_school_fkey,OnDelete:CASCADE;uniqueIndex:uq_subjects_school_name" validate:"required,uuid"`
	AcademicUnitID *uuid.UUID     `db:"academic_unit_id" gorm:"type:uuid;constraint:subjects_unit_fkey,OnDelete:SET NULL" validate:"omitempty,uuid"`
	Name           string         `db:"name" gorm:"not null;size:255;uniqueIndex:uq_subjects_school_name" validate:"required,min=2,max=255"`
	Code           *string        `db:"code" gorm:"type:varchar(50)" validate:"omitempty"`
	Description    *string        `db:"description" gorm:"default:null" validate:"omitempty"`
	IsActive       bool           `db:"is_active" gorm:"not null;default:true"`
	CreatedAt      time.Time      `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt      time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
	DeletedAt      gorm.DeletedAt `db:"deleted_at" gorm:"index" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Subject) TableName() string {
	return "academic.subjects"
}
