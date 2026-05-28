package common

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserSpec describe el usuario a sembrar. Password se recibe en texto
// plano; el helper aplica BcryptHash. TokenVersion no se setea: la BD
// (default 1) lo provee.
type UserSpec struct {
	ID        uuid.UUID
	Email     string
	Password  string // plaintext; el helper aplica BcryptHash internamente
	FirstName string
	LastName  string
}

// buildUser mapea UserSpec a entities.User aplicando BcryptHash al
// password.
func buildUser(spec UserSpec) entities.User {
	return entities.User{
		ID:           spec.ID,
		Email:        spec.Email,
		PasswordHash: BcryptHash(spec.Password),
		FirstName:    spec.FirstName,
		LastName:     spec.LastName,
		IsActive:     true,
	}
}

// SeedUser inserta el usuario con hash bcrypt aplicado. Idempotente por
// PK; si el ID ya existe, no actualiza (el hash no se rota).
func SeedUser(tx *gorm.DB, spec UserSpec) error {
	user := buildUser(spec)
	return OnConflictIgnore(tx, &user)
}
