package entities

import (
	"time"

	"github.com/google/uuid"
)

// Role representa un rol del sistema RBAC
type Role struct {
	ID          uuid.UUID `db:"id" gorm:"type:uuid;primaryKey"`
	Name        string    `db:"name" gorm:"uniqueIndex;not null"`
	DisplayName string    `db:"display_name" gorm:"not null"`
	Description *string   `db:"description" gorm:"default:null"`
	Scope       string    `db:"scope" gorm:"not null;type:varchar(50)"`
	IsActive    bool      `db:"is_active" gorm:"not null;default:true"`
	CreatedAt   time.Time `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Role) TableName() string {
	return "iam.roles"
}
