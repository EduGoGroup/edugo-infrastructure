package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User representa la tabla 'users' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migración: 001_create_users.up.sql
// Usada por: api-mobile, api-administracion, worker
// NOTA: Los roles se manejan a través de 'roles' y 'user_roles' (RBAC).
// NOTA: La escuela del usuario se obtiene desde 'memberships.school_id', no desde este struct.
type User struct {
	ID           uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey"`
	Email        string         `db:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string         `db:"password_hash" gorm:"not null"`
	FirstName    string         `db:"first_name" gorm:"not null"`
	LastName     string         `db:"last_name" gorm:"not null"`
	IsActive     bool           `db:"is_active" gorm:"not null;default:true"`
	CreatedAt    time.Time      `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `db:"deleted_at" gorm:"index"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (User) TableName() string {
	return "auth.users"
}
