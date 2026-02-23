package entities

import (
	"time"

	"github.com/google/uuid"
)

// Subject representa la tabla 'subjects' en PostgreSQL
// Esta entity es el reflejo exacto del schema de BD definido en:
// - postgres/migrations/013_create_subjects.up.sql
//
// Representa una materia o asignatura del sistema educativo.
type Subject struct {
	ID          uuid.UUID `db:"id" gorm:"type:uuid;primaryKey"`
	Name        string    `db:"name" gorm:"not null"`
	Description *string   `db:"description" gorm:"default:null"`
	Metadata    *string   `db:"metadata" gorm:"type:jsonb"`
	IsActive    bool      `db:"is_active" gorm:"not null;default:true"`
	CreatedAt   time.Time `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Subject) TableName() string {
	return "academic.subjects"
}
