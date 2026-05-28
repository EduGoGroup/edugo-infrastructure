package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Course representa la tabla 'courses' en el schema 'content' de PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Usada por: api-learning
type Course struct {
	ID          uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	UnitID      uuid.UUID      `db:"unit_id" gorm:"type:uuid;index;not null" validate:"required,uuid"`
	Name        string         `db:"name" gorm:"not null;size:255" validate:"required,min=2,max=255"`
	Description *string        `db:"description" gorm:"default:null" validate:"omitempty"`
	Status      string         `db:"status" gorm:"not null;type:varchar(50);default:'active'" validate:"required"`
	CreatedAt   time.Time      `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt   time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
	DeletedAt   gorm.DeletedAt `db:"deleted_at" gorm:"index" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Course) TableName() string {
	return "content.courses"
}
