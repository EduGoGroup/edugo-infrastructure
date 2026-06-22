package entities

import (
	"time"

	"github.com/google/uuid"
)

// Resource representa un recurso/modulo del sistema para RBAC y menu
type Resource struct {
	ID            uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	Key           string     `db:"key" gorm:"uniqueIndex:resources_key_key;not null;size:50" validate:"required,max=50"`
	DisplayName   string     `db:"display_name" gorm:"not null;size:150" validate:"required,min=2,max=150"`
	Description   *string    `db:"description" gorm:"default:null" validate:"omitempty"`
	Icon          *string    `db:"icon" gorm:"default:null;size:100" validate:"omitempty"`
	ParentID      *uuid.UUID `db:"parent_id" gorm:"type:uuid;index;constraint:fk_resources_parent,OnDelete:SET NULL" validate:"omitempty,uuid"`
	SortOrder     int        `db:"sort_order" gorm:"not null;default:0;index:idx_resources_sort" validate:"required"`
	IsMenuVisible bool       `db:"is_menu_visible" gorm:"not null;default:true;index:idx_resources_menu_visible"`
	// Plan 026 (overflow de navegación): priority/pin ADITIVOS al contrato del
	// menú. Priority nullable (NULL = modo legacy → el front cae a sort_order);
	// Pin default false (recurso anclado fuera del overflow). Ambos los
	// materializa GORM AutoMigrate igual que los campos vecinos.
	Priority *int `db:"priority" gorm:"default:null;index:idx_resources_priority"`
	Pin      bool `db:"pin" gorm:"not null;default:false;index:idx_resources_pin"`
	// ENUM: created in pre_gorm.sql
	Scope     string    `db:"scope" gorm:"not null;type:iam.permission_scope;default:'school'" validate:"required"`
	IsActive  bool      `db:"is_active" gorm:"not null;default:true;index:idx_resources_active"`
	CreatedAt time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Resource) TableName() string {
	return "iam.resources"
}
