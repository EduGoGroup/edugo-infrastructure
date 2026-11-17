package common

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserSpec describe el usuario a sembrar. No lleva credencial: auth.users es la
// proyección local del usuario y el login lo hace identity-core (plan 053 T6.2).
type UserSpec struct {
	ID        uuid.UUID
	Email     string
	FirstName string
	LastName  string
}

func buildUser(spec UserSpec) entities.User {
	return entities.User{
		ID:        spec.ID,
		Email:     spec.Email,
		FirstName: spec.FirstName,
		LastName:  spec.LastName,
		IsActive:  true,
	}
}

// SeedUser inserta el usuario. Idempotente por id; si el id ya existe, no lo
// actualiza.
func SeedUser(tx *gorm.DB, spec UserSpec) error {
	user := buildUser(spec)
	return onConflictIgnore(tx, &user)
}
