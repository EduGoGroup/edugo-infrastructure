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
	ID             uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey"`
	SchoolID       uuid.UUID      `db:"school_id" gorm:"type:uuid;index;not null"`
	AcademicUnitID *uuid.UUID     `db:"academic_unit_id" gorm:"type:uuid"`
	Name           string         `db:"name" gorm:"not null"`
	Code           *string        `db:"code" gorm:"type:varchar(50)"`
	Description    *string        `db:"description" gorm:"default:null"`
	IsActive       bool           `db:"is_active" gorm:"not null;default:true"`
	CreatedAt      time.Time      `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt      time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `db:"deleted_at" gorm:"index"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Subject) TableName() string {
	return "academic.subjects"
}
