package entities

import (
	"time"

	"github.com/google/uuid"
)

// Resource representa un recurso/modulo del sistema para RBAC y menu
type Resource struct {
	ID            uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey"`
	Key           string     `db:"key" gorm:"uniqueIndex;not null"`
	DisplayName   string     `db:"display_name" gorm:"not null"`
	Description   *string    `db:"description" gorm:"default:null"`
	Icon          *string    `db:"icon" gorm:"default:null"`
	ParentID      *uuid.UUID `db:"parent_id" gorm:"type:uuid;index"`
	SortOrder     int        `db:"sort_order" gorm:"not null;default:0"`
	IsMenuVisible bool       `db:"is_menu_visible" gorm:"not null;default:true"`
	Scope         string     `db:"scope" gorm:"not null;type:varchar(50)"`
	IsActive      bool       `db:"is_active" gorm:"not null;default:true"`
	CreatedAt     time.Time  `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt     time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Resource) TableName() string {
	return "iam.resources"
}
