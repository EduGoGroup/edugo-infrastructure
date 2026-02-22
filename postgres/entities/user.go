package entities

import (
	"time"

	"github.com/google/uuid"
)

// User representa la tabla 'users' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migración: 001_create_users.up.sql
// Usada por: api-mobile, api-administracion, worker
// NOTA: Los roles se manejan a través de 'roles' y 'user_roles' (RBAC).
// NOTA: La escuela del usuario se obtiene desde 'memberships.school_id', no desde este struct.
type User struct {
	ID           uuid.UUID  `db:"id"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	FirstName    string     `db:"first_name"`
	LastName     string     `db:"last_name"`
	IsActive     bool       `db:"is_active"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (User) TableName() string {
	return "users"
}
