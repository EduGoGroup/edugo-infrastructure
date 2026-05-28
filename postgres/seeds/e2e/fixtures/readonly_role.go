package fixtures

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
)

// ReadonlyRole crea un rol overlay con scope `unit` que agrupa
// exclusivamente grants `<resource_path>.read` para una lista de
// recursos del production seed.
//
// Diseñada para reproducir el caso "puedo entrar a la lista pero no
// edito nada". El rol vive en el namespace del scenario (UUID con
// SchemaPrefix) y nunca toca los roles ni los grants del catálogo.
//
// P4-1 (plan B): previamente la fixture insertaba iam.role_permissions
// (asignaciones 1:1). Tras la eliminación de esa tabla, siembra grants
// directos en iam.role_grants con patterns `<resource_path>.read`
// (effect=allow). El id de cada role_grant se deriva del SchemaPrefix
// para que el cleanup selectivo lo encuentre por prefijo.
type ReadonlyRole struct {
	// Resources enumera las keys de los recursos
	// (ej. ["assessments", "grades"]) sobre los que se concederá
	// el permiso `<resource>:read`. Si el slice está vacío, Apply
	// falla.
	Resources []string
}

// Manifest implementa framework.Fixture.
func (f *ReadonlyRole) Manifest() framework.FixtureManifest {
	return framework.FixtureManifest{
		Name:     "readonly_role",
		Provides: []string{"readonly_role"},
		Tables: []string{
			"iam.role_grants",
			"iam.roles",
		},
		Constants: map[string]string{
			"E2EFixtureReadonlyRoleCode": "{{.RoleCode}}",
		},
		Description: "Crea un rol overlay con scope=unit que sólo expone <resource>:read para los recursos parametrizados.",
	}
}

// Apply implementa framework.Fixture.
func (f *ReadonlyRole) Apply(tx *gorm.DB, ctx *framework.ApplyContext) error {
	if len(f.Resources) == 0 {
		return fmt.Errorf("readonly_role: Resources requerido (al menos 1 resource key)")
	}
	for i, key := range f.Resources {
		if key == "" {
			return fmt.Errorf("readonly_role: Resources[%d] vacío", i)
		}
	}
	if ctx == nil {
		return fmt.Errorf("readonly_role: nil ApplyContext")
	}
	if tx == nil {
		return fmt.Errorf("readonly_role: nil transaction")
	}

	hash := schemaHashFromCtx(ctx)
	roleName := "ro_" + hash
	roleID := framework.MakeUUID(ctx, "0000-0000-0000-f0a000000001")
	if err := framework.AssertNotProductionNamespace(roleID); err != nil {
		return err
	}
	parsedRoleID, err := uuid.Parse(roleID)
	if err != nil {
		return fmt.Errorf("readonly_role: role UUID inválido (%q): %w", roleID, err)
	}

	role := entities.Role{
		ID:          parsedRoleID,
		Name:        roleName,
		DisplayName: "ReadOnly Auditor (E2E)",
		Scope:       "unit",
		IsActive:    true,
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&role).Error; err != nil {
		return fmt.Errorf("readonly_role: insert role: %w", err)
	}
	if err := framework.UpsertBool(tx, role.TableName(), "id", role.ID, "is_active", true); err != nil {
		return err
	}

	// Asignar grant `<resource_path>.read` por cada resource pedido.
	for i, resourceKey := range f.Resources {
		permName, err := lookupPermissionName(tx, resourceKey, "read")
		if err != nil {
			return fmt.Errorf("readonly_role: lookup permission %s:read: %w", resourceKey, err)
		}
		grantID := framework.MakeUUID(ctx, fmt.Sprintf("0000-0000-0000-f0b%09d", i+1))
		if err := framework.AssertNotProductionNamespace(grantID); err != nil {
			return err
		}
		parsedGrant, err := uuid.Parse(grantID)
		if err != nil {
			return fmt.Errorf("readonly_role: role_grant UUID inválido (%q): %w", grantID, err)
		}
		grant := entities.RoleGrant{
			ID:      parsedGrant,
			RoleID:  parsedRoleID,
			Pattern: permName,
			Effect:  "allow",
		}
		// OnConflict por el UNIQUE natural (role_id, pattern, effect):
		// otra fixture del scenario (p.ej. MenuSubtree) puede haber
		// sembrado el mismo pattern bajo otro id.
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "role_id"}, {Name: "pattern"}, {Name: "effect"}},
			DoNothing: true,
		}).Create(&grant).Error; err != nil {
			return fmt.Errorf("readonly_role: insert role_grant %s.read: %w", resourceKey, err)
		}
	}

	ctx.Provide("readonly_role", framework.ProvidedEntity{
		Kind:  "role",
		ID:    roleID,
		Extra: map[string]string{"name": roleName, "scope": "unit"},
	})
	ctx.SetConstant("E2EFixtureReadonlyRoleCode", roleName)
	ctx.SetConstant("E2EFixtureReadonlyRoleID", roleID)
	return nil
}

// Cleanup implementa framework.Fixture. Borra primero role_grants y
// después el role; ambos por prefijo del scenario.
func (f *ReadonlyRole) Cleanup(tx *gorm.DB, ctx *framework.ApplyContext) error {
	if tx == nil {
		return fmt.Errorf("readonly_role cleanup: nil transaction")
	}
	if ctx == nil || ctx.SchemaPrefix == "" {
		return fmt.Errorf("readonly_role cleanup: SchemaPrefix vacío")
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
			return fmt.Errorf("readonly_role cleanup %s: %w", t.name, err)
		}
	}
	return nil
}

// lookupPermissionName resuelve el `name` (path) del permiso de un par
// (resource_key, action) consultando iam.resources e iam.permissions
// del production seed. Devuelve el name porque iam.role_grants.pattern
// almacena el path textual, no el FK al permission.
func lookupPermissionName(tx *gorm.DB, resourceKey, action string) (string, error) {
	const stmt = `
		SELECT name FROM iam.permissions
		WHERE action = ?
		  AND resource_id = (SELECT id FROM iam.resources WHERE key = ?)
	`
	var name string
	row := tx.Raw(stmt, action, resourceKey).Row()
	if err := row.Scan(&name); err != nil {
		return "", fmt.Errorf("permission %s:%s no encontrado: %w", resourceKey, action, err)
	}
	return name, nil
}

// schemaHashFromCtx extrae el hash del SchemaPrefix
// (e2e<hash>-) del scenario para componer códigos legibles
// (ej. "ro_a1b2c3d4"). Si el prefijo es el legacy (e2e00000-), devuelve
// el LegacyHash.
func schemaHashFromCtx(ctx *framework.ApplyContext) string {
	if ctx == nil {
		return framework.LegacyHash
	}
	prefix := ctx.SchemaPrefix
	// Quitar "e2e" inicial y "-" final.
	if len(prefix) > 4 && prefix[:3] == "e2e" && prefix[len(prefix)-1] == '-' {
		return prefix[3 : len(prefix)-1]
	}
	return framework.LegacyHash
}
