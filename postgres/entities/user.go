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
	ID           uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	Email        string         `db:"email" gorm:"uniqueIndex:users_email_unique;not null;size:255" validate:"required,email"`
	PasswordHash string         `db:"password_hash" gorm:"not null;size:255" validate:"required,min=8"`
	FirstName    string         `db:"first_name" gorm:"not null;size:100" validate:"required,min=2,max=100"`
	LastName     string         `db:"last_name" gorm:"not null;size:100" validate:"required,min=2,max=100"`
	// NOTE: partial index idx_users_active (WHERE is_active = true) must be created in post_gorm.sql
	IsActive     bool           `db:"is_active" gorm:"not null;default:true"`
	TokenVersion int            `db:"token_version" gorm:"not null;default:1"`
	CreatedAt    time.Time      `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt    time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
	DeletedAt    gorm.DeletedAt `db:"deleted_at" gorm:"index" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (User) TableName() string {
	return "auth.users"
}
