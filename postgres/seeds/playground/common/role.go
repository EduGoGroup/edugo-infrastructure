package common

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RoleSpec describe el rol a sembrar. Si Description está vacío, el campo
// de la entidad queda en nil (puntero); si Scope está vacío, default
// "school".
type RoleSpec struct {
	ID          uuid.UUID
	Name        string
	DisplayName string
	Description string // si vacío, Description queda nil (puntero)
	Scope       string // default "school" si vacío
}

// buildRole mapea RoleSpec a entities.Role aplicando defaults.
func buildRole(spec RoleSpec) entities.Role {
	var desc *string
	if spec.Description != "" {
		d := spec.Description
		desc = &d
	}
	scope := spec.Scope
	if scope == "" {
		scope = "school"
	}
	return entities.Role{
		ID:          spec.ID,
		Name:        spec.Name,
		DisplayName: spec.DisplayName,
		Description: desc,
		Scope:       scope,
		IsActive:    true,
	}
}

// SeedRole inserta el rol aplicando defaults. Idempotente por PK.
func SeedRole(tx *gorm.DB, spec RoleSpec) error {
	role := buildRole(spec)
	return OnConflictIgnore(tx, &role)
}

// buildRoleGrant construye el RoleGrant con ID determinístico
// SHA1(roleID:pattern:allow), igual que insertGrants en focal_evaluacion_v2
// y focal_botonera.
func buildRoleGrant(roleID uuid.UUID, pattern string) entities.RoleGrant {
	const effect = "allow"
	gid := uuid.NewSHA1(uuid.NameSpaceOID, []byte(roleID.String()+":"+pattern+":"+effect))
	return entities.RoleGrant{
		ID:      gid,
		RoleID:  roleID,
		Pattern: pattern,
		Effect:  effect,
	}
}

// SeedRoleGrant inserta un allow-grant con ID determinístico y respeta el
// uniqueIndex (role_id, pattern, effect). Idempotente.
func SeedRoleGrant(tx *gorm.DB, roleID uuid.UUID, pattern string) error {
	grant := buildRoleGrant(roleID, pattern)
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "role_id"}, {Name: "pattern"}, {Name: "effect"}},
		DoNothing: true,
	}).Create(&grant).Error
}
