package entities

import (
	"time"

	"github.com/google/uuid"
)

// Permission representa un permiso del sistema RBAC
type Permission struct {
	ID          uuid.UUID `db:"id" gorm:"type:uuid;primaryKey"`
	Name        string    `db:"name" gorm:"uniqueIndex;not null"`
	DisplayName string    `db:"display_name" gorm:"not null"`
	Description *string   `db:"description" gorm:"default:null"`
	ResourceID  uuid.UUID `db:"resource_id" gorm:"type:uuid;index;not null"`
	ResourceKey string    `db:"resource_key" gorm:"not null"`
	Action      string    `db:"action" gorm:"not null;type:varchar(50)"`
	Scope       string    `db:"scope" gorm:"not null;type:varchar(50)"`
	IsActive    bool      `db:"is_active" gorm:"not null;default:true"`
	CreatedAt   time.Time `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Permission) TableName() string {
	return "iam.permissions"
}
