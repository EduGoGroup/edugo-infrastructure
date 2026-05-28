package fixtures

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
)

// MenuSubtree garantiza que un sub-árbol del menú (ej. raíz
// `academic` o `admin`) esté visible para el rol overlay del scenario,
// sembrando los grants `<resource_path>.read` necesarios para que esa
// rama aparezca al filtrar el menú.
//
// La fixture descubre el subtree con un CTE recursivo contra
// iam.resources (production seed). No toca esas filas: sólo crea filas
// nuevas en iam.role_grants vinculadas al rol overlay.
//
// P4-1 (plan B): la fixture previamente sembraba iam.role_permissions
// (asignaciones 1:1 rol×permiso). Tras la eliminación de la tabla, la
// fixture siembra iam.role_grants con patterns `<resource_path>.read`
// (effect=allow). El effective permission del rol se computa con el
// matcher wildcard (edugo-shared/auth.PermissionMatches).
//
// Si el campo RoleID está vacío, MenuSubtree depende de la capacidad
// "readonly_role" provista por la fixture sibling ReadonlyRole; si
// está informado, lo usa tal cual (compatibilidad con scenarios que
// arman su propio rol overlay).
type MenuSubtree struct {
	// SubtreeRoot es el `key` del recurso raíz (ej. "academic",
	// "admin") cuyo subtree se va a iluminar.
	SubtreeRoot string

	// RoleID opcional: UUID del rol overlay al que se asignarán los
	// grants. Si vacío, se toma de ctx.Provided["readonly_role"].ID.
	RoleID string
}

// Manifest implementa framework.Fixture.
func (f *MenuSubtree) Manifest() framework.FixtureManifest {
	requires := []string{}
	if f.RoleID == "" {
		requires = append(requires, "readonly_role")
	}
	return framework.FixtureManifest{
		Name:     "menu_subtree",
		Provides: []string{"menu_subtree"},
		Requires: requires,
		Tables: []string{
			"iam.role_grants",
		},
		Constants: map[string]string{
			"E2EFixtureMenuSubtreeRoot": "{{.SubtreeRoot}}",
		},
		Description: "Siembra role_grants <resource_path>.read para iluminar una rama del menú al rol overlay del scenario.",
	}
}

// Apply implementa framework.Fixture.
func (f *MenuSubtree) Apply(tx *gorm.DB, ctx *framework.ApplyContext) error {
	if f.SubtreeRoot == "" {
		return fmt.Errorf("menu_subtree: SubtreeRoot requerido")
	}
	if ctx == nil {
		return fmt.Errorf("menu_subtree: nil ApplyContext")
	}

	// Resolver RoleID: explícito o desde ctx.Provided.
	roleIDStr := f.RoleID
	if roleIDStr == "" {
		ent, ok := ctx.Provided["readonly_role"]
		if !ok || ent.ID == "" {
			return fmt.Errorf("menu_subtree: capability %q no provista y RoleID vacío", "readonly_role")
		}
		roleIDStr = ent.ID
	}
	roleUUID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return fmt.Errorf("menu_subtree: RoleID inválido (%q): %w", roleIDStr, err)
	}
	if tx == nil {
		return fmt.Errorf("menu_subtree: nil transaction")
	}

	resourceIDs, err := lookupSubtreeResourceIDs(tx, f.SubtreeRoot)
	if err != nil {
		return fmt.Errorf("menu_subtree: lookup subtree %q: %w", f.SubtreeRoot, err)
	}
	if len(resourceIDs) == 0 {
		return fmt.Errorf("menu_subtree: subtree %q vacío (¿key inexistente en iam.resources?)", f.SubtreeRoot)
	}

	for i, resID := range resourceIDs {
		permName, err := lookupPermissionNameByResourceID(tx, resID, "read")
		if err != nil {
			return fmt.Errorf("menu_subtree: lookup permission read sobre resource_id=%s: %w", resID, err)
		}
		grantID := framework.MakeUUID(ctx, fmt.Sprintf("0000-0000-0000-f0e%09d", i+1))
		if err := framework.AssertNotProductionNamespace(grantID); err != nil {
			return err
		}
		parsedGrant, err := uuid.Parse(grantID)
		if err != nil {
			return fmt.Errorf("menu_subtree: role_grant UUID inválido (%q): %w", grantID, err)
		}
		grant := entities.RoleGrant{
			ID:      parsedGrant,
			RoleID:  roleUUID,
			Pattern: permName,
			Effect:  "allow",
		}
		// OnConflict por (role_id, pattern, effect) que es el UNIQUE
		// natural de iam.role_grants: una fixture sibling puede haber
		// sembrado el mismo pattern bajo otro id.
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "role_id"}, {Name: "pattern"}, {Name: "effect"}},
			DoNothing: true,
		}).Create(&grant).Error; err != nil {
			return fmt.Errorf("menu_subtree: insert role_grant resource_id=%s pattern=%s: %w", resID, permName, err)
		}
	}

	ctx.Provide("menu_subtree", framework.ProvidedEntity{
		Kind: "menu_subtree",
		ID:   "",
		Extra: map[string]string{
			"root":            f.SubtreeRoot,
			"resources_count": fmt.Sprintf("%d", len(resourceIDs)),
			"role_id":         roleIDStr,
		},
	})
	ctx.SetConstant("E2EFixtureMenuSubtreeRoot", f.SubtreeRoot)
	ctx.SetConstant("E2EFixtureMenuSubtreeRoleID", roleIDStr)
	return nil
}

// Cleanup implementa framework.Fixture. Sólo borra role_grants del
// scenario (las roles las gestiona la fixture que las creó).
func (f *MenuSubtree) Cleanup(tx *gorm.DB, ctx *framework.ApplyContext) error {
	if tx == nil {
		return fmt.Errorf("menu_subtree cleanup: nil transaction")
	}
	if ctx == nil || ctx.SchemaPrefix == "" {
		return fmt.Errorf("menu_subtree cleanup: SchemaPrefix vacío")
	}
	prefix := ctx.SchemaPrefix
	tables := []struct {
		name string
		col  string
	}{
		{"iam.role_grants", "id"},
	}
	for _, t := range tables {
		if _, err := framework.DeleteByPrefix(tx, t.name, t.col, prefix); err != nil {
			return fmt.Errorf("menu_subtree cleanup %s: %w", t.name, err)
		}
	}
	return nil
}

// lookupSubtreeResourceIDs ejecuta el CTE recursivo sobre
// iam.resources y devuelve todos los IDs del subtree (incluido el root).
func lookupSubtreeResourceIDs(tx *gorm.DB, rootKey string) ([]uuid.UUID, error) {
	const stmt = `
		WITH RECURSIVE subtree AS (
			SELECT id, key FROM iam.resources WHERE key = ?
			UNION
			SELECT r.id, r.key FROM iam.resources r
			INNER JOIN subtree s ON r.parent_id = s.id
		)
		SELECT id FROM subtree
	`
	rows, err := tx.Raw(stmt, rootKey).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

// lookupPermissionNameByResourceID resuelve el `name` del permiso
// (resource_id, action) sobre iam.permissions. Devuelve el name (no el
// id) porque iam.role_grants.pattern almacena el path textual, no FK.
func lookupPermissionNameByResourceID(tx *gorm.DB, resourceID uuid.UUID, action string) (string, error) {
	const stmt = `SELECT name FROM iam.permissions WHERE resource_id = ? AND action = ?`
	var name string
	row := tx.Raw(stmt, resourceID, action).Row()
	if err := row.Scan(&name); err != nil {
		return "", fmt.Errorf("permission resource_id=%s action=%s no encontrado: %w", resourceID, action, err)
	}
	return name, nil
}
