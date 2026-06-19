package entities

import (
	"time"

	"github.com/google/uuid"
)

// Role representa un rol del sistema RBAC
type Role struct {
	ID          uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	Name        string    `db:"name" gorm:"uniqueIndex:roles_name_key;not null;size:50" validate:"required,min=2,max=50"`
	DisplayName string    `db:"display_name" gorm:"not null;size:100" validate:"required,min=2,max=100"`
	Description *string   `db:"description" gorm:"default:null" validate:"omitempty"`
	// ENUM: created in pre_gorm.sql
	Scope string `db:"scope" gorm:"not null;type:iam.role_scope;default:'school'" validate:"required"`
	// ParentRoleID enlaza un rol alias con su rol canónico (herencia de
	// grants; ADR-6). NULL para roles canónicos. La herencia se resuelve
	// y aplana en el login: los grants del ancestro se funden con los
	// propios antes de emitir el JWT. FK self-referencial nullable a
	// iam.roles(id); ON DELETE SET NULL para no romper el alias si se
	// borra el canónico.
	ParentRoleID *uuid.UUID `db:"parent_role_id" gorm:"type:uuid;default:null;constraint:fk_roles_parent,OnDelete:SET NULL;index:idx_roles_parent" validate:"omitempty"`
	// LandingScreenKey es el screen_key del dashboard de inicio (landing) de
	// este rol (ADR 0024 DEC-2: string nullable, NO UUID, NO FK). NULL para
	// roles alias (caen al default de la escuela o al fallback de sistema; la
	// herencia del landing es mejora futura, no F0). Solo dato sembrado; la
	// resolución vive en el backend (F1).
	LandingScreenKey *string `db:"landing_screen_key" gorm:"type:varchar(64)" validate:"omitempty"`
	IsActive         bool    `db:"is_active" gorm:"not null;default:true;index:idx_roles_active"`
	// IsSystem marca un rol del contrato del sistema (sembrado por L0–L4).
	// Cuando es true, el rol es inmutable en runtime: los usecases de
	// mutación (delete/update) lo rechazan para proteger el contrato de
	// permisos de borrado/edición vía API (bug 0069). Los roles creados por
	// usuarios en runtime quedan en false (default).
	IsSystem  bool      `db:"is_system" gorm:"not null;default:false;index:idx_roles_system"`
	CreatedAt time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Role) TableName() string {
	return "iam.roles"
}
