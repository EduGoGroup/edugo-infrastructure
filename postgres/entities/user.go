package entities

import (
	"time"

	"github.com/google/uuid"
)

// User representa la tabla 'users' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migración: 001_create_users.up.sql, 017_add_school_id_to_users.up.sql
// Usada por: api-mobile, api-administracion, worker
type User struct {
	ID           uuid.UUID  `db:"id"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	FirstName    string     `db:"first_name"`
	LastName     string     `db:"last_name"`
	Role         string     `db:"role"`      // admin, teacher, student, guardian
	SchoolID     *uuid.UUID `db:"school_id"` // Escuela principal del usuario (nullable para super_admin)
	IsActive     bool       `db:"is_active"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (User) TableName() string {
	return "users"
}

// GetSchoolIDString retorna el school_id como string o vacío si es nil
func (u *User) GetSchoolIDString() string {
	if u.SchoolID == nil {
		return ""
	}
	return u.SchoolID.String()
}
