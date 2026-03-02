package entities

import (
	"time"

	"github.com/google/uuid"
)

// RolePermission representa la asignación de un permiso a un rol
type RolePermission struct {
	ID           uuid.UUID `db:"id" gorm:"type:uuid;primaryKey"`
	RoleID       uuid.UUID `db:"role_id" gorm:"type:uuid;index;not null"`
	PermissionID uuid.UUID `db:"permission_id" gorm:"type:uuid;index;not null"`
	CreatedAt    time.Time `db:"created_at" gorm:"not null;autoCreateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (RolePermission) TableName() string {
	return "iam.role_permissions"
}
