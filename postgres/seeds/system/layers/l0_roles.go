package layers

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// applyL0Roles siembra el rol super_admin de L0 y sus 4 permisos CRUD
// sobre el recurso announcements.
//
// P4-1 (plan B): se eliminó upsertL0RolePermissions. La asignación
// rol → permisos del modelo legacy (tabla iam.role_permissions) ya no
// existe. Los permisos efectivos del super_admin se otorgan vía
// iam.role_grants con el pattern wildcard `*`, sembrado por
// applyL4RoleGrants en seeds/system/l4/roles_permissions.go.
//
// Orden interno (FK):
//  1. iam.roles      → super_admin
//  2. iam.permissions → academic.announcements.{read,create,update,delete}
//     (FK a iam.resources(announcements), sembrado por applyL0Resources)
//
// Idempotente: usa ON CONFLICT DO NOTHING. Reaplicar no produce duplicados.
func applyL0Roles(tx *gorm.DB) error {
	if err := upsertL0Role(tx); err != nil {
		return err
	}
	if err := upsertL0Permissions(tx); err != nil {
		return err
	}
	return nil
}

// upsertL0Role inserta el rol super_admin (scope=system).
// Es el único rol de L0 y el ancla del bootstrap: el usuario L0
// (sembrado por applyL0Users) lo recibirá vía user_role.
func upsertL0Role(tx *gorm.DB) error {
	id, err := uuid.Parse(L0_ROLE_SUPER_ADMIN_ID)
	if err != nil {
		return fmt.Errorf("upsertL0Role: parse id: %w", err)
	}
	desc := "Rol con acceso total al sistema. Usuario de bootstrapping."
	role := entities.Role{
		ID:          id,
		Name:        L0_ROLE_SUPER_ADMIN_NAME,
		DisplayName: "Super Administrador",
		Description: &desc,
		Scope:       "system",
		IsActive:    true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&role).Error
}

// upsertL0Permissions inserta los 4 permisos CRUD sobre announcements.
// Todos comparten ResourceID (L0_RESOURCE_ANNOUNCEMENTS_ID) y Scope=school,
// que es el scope canónico de operación del producto.
func upsertL0Permissions(tx *gorm.DB) error {
	resourceID, err := uuid.Parse(L0_RESOURCE_ANNOUNCEMENTS_ID)
	if err != nil {
		return fmt.Errorf("upsertL0Permissions: parse resource_id: %w", err)
	}

	specs := []struct {
		idStr       string
		name        string
		displayName string
		action      string
	}{
		{L0_PERM_ANNOUNCEMENTS_READ, "academic.announcements.read", "Ver Anuncios", "read"},
		{L0_PERM_ANNOUNCEMENTS_CREATE, "academic.announcements.create", "Crear Anuncio", "create"},
		{L0_PERM_ANNOUNCEMENTS_UPDATE, "academic.announcements.update", "Editar Anuncio", "update"},
		{L0_PERM_ANNOUNCEMENTS_DELETE, "academic.announcements.delete", "Eliminar Anuncio", "delete"},
	}

	perms := make([]entities.Permission, 0, len(specs))
	for _, s := range specs {
		id, err := uuid.Parse(s.idStr)
		if err != nil {
			return fmt.Errorf("upsertL0Permissions: parse id %s: %w", s.idStr, err)
		}
		perms = append(perms, entities.Permission{
			ID:          id,
			Name:        s.name,
			DisplayName: s.displayName,
			ResourceID:  resourceID,
			Action:      s.action,
			Scope:       "school",
			IsActive:    true,
		})
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).CreateInBatches(&perms, 10).Error
}
