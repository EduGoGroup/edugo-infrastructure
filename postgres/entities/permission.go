package entities

import (
	"time"

	"github.com/google/uuid"
)

// Permission representa un permiso del sistema RBAC
type Permission struct {
	ID          uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	// La profundidad del path llega a 4 segmentos (p.ej.
	// academic.join_request_approvals.school.student — SELLO × TIPO): el
	// cuantificador {0,3} alinea esta CHECK con enum.PathPermissionRegex del
	// shared (que ya permite hasta 4 segmentos). Mantener ambos en sincronía.
	Name        string    `db:"name" gorm:"uniqueIndex:permissions_name_key;not null;size:100;check:chk_permission_name_format,name ~ '^(\\*|[a-z_]+(\\.[a-z_]+){0,3}(\\.\\*)?)(:own)?$'" validate:"required,min=2,max=100"`
	DisplayName string    `db:"display_name" gorm:"not null;size:150" validate:"required,min=2,max=150"`
	Description *string   `db:"description" gorm:"default:null" validate:"omitempty"`
	ResourceID  uuid.UUID `db:"resource_id" gorm:"type:uuid;index;not null;constraint:fk_permissions_resource,OnDelete:RESTRICT;uniqueIndex:uq_permissions_resource_action" validate:"required,uuid"`
	Action      string    `db:"action" gorm:"not null;type:varchar(50);uniqueIndex:uq_permissions_resource_action" validate:"required"`
	// ENUM: created in pre_gorm.sql
	Scope     string    `db:"scope" gorm:"not null;type:iam.permission_scope;default:'school'" validate:"required"`
	IsActive  bool      `db:"is_active" gorm:"not null;default:true;index:idx_permissions_active"`
	CreatedAt time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Permission) TableName() string {
	return "iam.permissions"
}
