package fixtures

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
)

// PartialCrudRoleCodeConstant es la clave que la fixture publica en
// ctx.Constants para que los tests Kotlin descubran el code del rol
// overlay creado.
const PartialCrudRoleCodeConstant = "E2EFixturePartialCrudRoleCode"

// PartialCrud crea un rol overlay similar a ReadonlyRole pero asigna
// los pares (create, read) por cada resource — sin update ni delete.
//
// Reproduce el caso "puedo listar y agregar pero no editar". Útil para
// los tests Rol × Pantalla × Acción que diferencian "listado vs
// edición" como dos pruebas independientes.
type PartialCrud struct {
	// Resources enumera las keys de los recursos sobre los que se
	// concederán los permisos (create, read).
	Resources []string
}

// Manifest implementa framework.Fixture.
func (f *PartialCrud) Manifest() framework.FixtureManifest {
	return framework.FixtureManifest{
		Name:     "partial_crud",
		Provides: []string{"partial_crud_role"},
		Tables: []string{
			"iam.role_grants",
			"iam.roles",
		},
		Constants: map[string]string{
			PartialCrudRoleCodeConstant: "{{.RoleCode}}",
		},
		Description: "Crea un rol overlay con pares (create, read) por resource — útil para reproducir 'list+create sin edit ni delete'.",
	}
}

// Apply implementa framework.Fixture.
func (f *PartialCrud) Apply(tx *gorm.DB, ctx *framework.ApplyContext) error {
	if len(f.Resources) == 0 {
		return fmt.Errorf("partial_crud: Resources requerido (al menos 1 resource key)")
	}
	for i, key := range f.Resources {
		if key == "" {
			return fmt.Errorf("partial_crud: Resources[%d] vacío", i)
		}
	}
	if ctx == nil {
		return fmt.Errorf("partial_crud: nil ApplyContext")
	}
	if tx == nil {
		return fmt.Errorf("partial_crud: nil transaction")
	}

	hash := schemaHashFromCtx(ctx)
	roleName := "pcrud_" + hash
	roleID := framework.MakeUUID(ctx, "0000-0000-0000-f0c000000001")
	if err := framework.AssertNotProductionNamespace(roleID); err != nil {
		return err
	}
	parsedRoleID, err := uuid.Parse(roleID)
	if err != nil {
		return fmt.Errorf("partial_crud: role UUID inválido (%q): %w", roleID, err)
	}

	role := entities.Role{
		ID:          parsedRoleID,
		Name:        roleName,
		DisplayName: "Partial CRUD (E2E)",
		Scope:       "unit",
		IsActive:    true,
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&role).Error; err != nil {
		return fmt.Errorf("partial_crud: insert role: %w", err)
	}
	if err := framework.UpsertBool(tx, role.TableName(), "id", role.ID, "is_active", true); err != nil {
		return err
	}

	// Asignar grants (create, read) por cada resource. P4-1 (plan B):
	// se siembra iam.role_grants con patterns `<resource_path>.<action>`.
	actions := []string{"create", "read"}
	idx := 0
	for _, resourceKey := range f.Resources {
		for _, action := range actions {
			permName, err := lookupPermissionName(tx, resourceKey, action)
			if err != nil {
				return fmt.Errorf("partial_crud: lookup permission %s:%s: %w", resourceKey, action, err)
			}
			idx++
			grantID := framework.MakeUUID(ctx, fmt.Sprintf("0000-0000-0000-f0d%09d", idx))
			if err := framework.AssertNotProductionNamespace(grantID); err != nil {
				return err
			}
			parsedGrant, err := uuid.Parse(grantID)
			if err != nil {
				return fmt.Errorf("partial_crud: role_grant UUID inválido (%q): %w", grantID, err)
			}
			grant := entities.RoleGrant{
				ID:      parsedGrant,
				RoleID:  parsedRoleID,
				Pattern: permName,
				Effect:  "allow",
			}
			// OnConflict por UNIQUE (role_id, pattern, effect): evita
			// choques con fixtures siblings que ya hayan sembrado el
			// mismo pattern.
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "role_id"}, {Name: "pattern"}, {Name: "effect"}},
				DoNothing: true,
			}).Create(&grant).Error; err != nil {
				return fmt.Errorf("partial_crud: insert role_grant %s.%s: %w", resourceKey, action, err)
			}
		}
	}

	ctx.Provide("partial_crud_role", framework.ProvidedEntity{
		Kind:  "role",
		ID:    roleID,
		Extra: map[string]string{"name": roleName, "scope": "unit"},
	})
	ctx.SetConstant(PartialCrudRoleCodeConstant, roleName)
	ctx.SetConstant("E2EFixturePartialCrudRoleID", roleID)
	return nil
}

// Cleanup implementa framework.Fixture.
func (f *PartialCrud) Cleanup(tx *gorm.DB, ctx *framework.ApplyContext) error {
	if tx == nil {
		return fmt.Errorf("partial_crud cleanup: nil transaction")
	}
	if ctx == nil || ctx.SchemaPrefix == "" {
		return fmt.Errorf("partial_crud cleanup: SchemaPrefix vacío")
	}
	prefix := ctx.SchemaPrefix
	tables := []struct {
		name string
		col  string
	}{
		{"iam.role_grants", "id"},
		{"iam.roles", "id"},
	}
	for _, t := range tables {
		if _, err := framework.DeleteByPrefix(tx, t.name, t.col, prefix); err != nil {
			return fmt.Errorf("partial_crud cleanup %s: %w", t.name, err)
		}
	}
	return nil
}
