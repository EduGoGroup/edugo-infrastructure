package fixtures

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
)

// L3IsolationConstants es una fixture pasiva que valida la presencia y
// la forma de las filas L3 (sembradas por system.ApplySystem cuando L3
// está registrada en system.Layers()) y exporta sus identificadores al
// ApplyContext para que tests downstream y el JSON
// fixtures-constants.json puedan referenciarlos sin hardcodear.
//
// NO escribe filas: L3 vive en el namespace del production seed, y la
// regla del framework prohíbe que las fixtures lo modifiquen
// (ver framework.Fixture docstring, C-REQ-10.2). Por la misma razón
// Manifest no declara Provides, Requires ni Tables: la fixture no
// genera entidades ni participa en el cleanup selectivo por prefijo.
//
// Además de la presencia de filas, la fixture realiza las assertions
// SQL focales de la Fase 5 (cubren F5-REQ-1.1, F5-REQ-2.1, F5-REQ-2.2,
// F5-REQ-3.1, F5-REQ-3.2, F5-REQ-3.3 y no-regresión sobre la cadena L1
// viewer→permisos). Las assertions HTTP/UI (F5-REQ-2.3, F5-REQ-4.1,
// F5-REQ-4.2, F5-REQ-4.3, F5-REQ-6.2, F5-REQ-6.3) quedan diferidas —
// ver docstring del scenario L3Isolation.
//
// Refs: phase-5-layer-l3/{requirements,design}.md.
type L3IsolationConstants struct{}

// Manifest implementa framework.Fixture.
func (f *L3IsolationConstants) Manifest() framework.FixtureManifest {
	return framework.FixtureManifest{
		Name:        "l3_constants_export",
		Description: "Verifica filas L3 sembradas por system.ApplySystem (resource materials + 3 permisos sin :delete + 1 ResourceScreen `list` default + 1 ScreenInstance mínima materials-list FK-satisfying; material-form podado — pantallas nativas) y exporta sus identificadores al JSON.",
		Constants: map[string]string{
			"E2EFixtureL3ResourceMaterialsID":           layers.L3_RESOURCE_MATERIALS_ID,
			"E2EFixtureL3ResourceMaterialsKey":          layers.L3_RESOURCE_MATERIALS_KEY,
			"E2EFixtureL3PermMaterialsReadID":           layers.L3_PERM_MATERIALS_READ_ID,
			"E2EFixtureL3PermMaterialsCreateID":         layers.L3_PERM_MATERIALS_CREATE_ID,
			"E2EFixtureL3PermMaterialsUpdateID":         layers.L3_PERM_MATERIALS_UPDATE_ID,
			"E2EFixtureL3ScreenInstanceMaterialsListID": layers.L3_SCREEN_INSTANCE_MATERIALS_LIST_ID,
			"E2EFixtureL3ScreenInstanceMaterialFormID":  layers.L3_SCREEN_INSTANCE_MATERIAL_FORM_ID,
			"E2EFixtureL3ScreenKeyMaterialsList":        layers.L3_SCREEN_KEY_MATERIALS_LIST,
			"E2EFixtureL3ScreenKeyMaterialForm":         layers.L3_SCREEN_KEY_MATERIAL_FORM,
		},
	}
}

// Apply verifica L3 y exporta constantes. Idempotente: sólo lee y
// llama a ctx.SetConstant.
//
// Cubre por SQL:
//   - F5-REQ-1.1: resource materials existe con scope=unit, parent_id NULL,
//     is_menu_visible=true, is_active=true.
//   - F5-REQ-2.1: 3 permisos materials:{read,create,update}; ausencia
//     explícita de materials:delete.
//   - F5-REQ-2.2: 3 role_permissions super_admin × materials; ausencia
//     explícita de super_admin × materials:delete.
//   - F5-REQ-3.1/3.2 (post-poda + F2): `material-form` screen_instance
//     sigue ELIMINADA (pantalla nativa, sin mapping). `materials-list`
//     screen_instance se conserva MÍNIMA (no renderizada) por la FK del
//     mapping de menú. No se verifica su forma de slot_data (es nativa).
//   - F5-REQ-3.3 (post-poda): 1 resource_screen (list default; el form
//     fue podado). La pantalla `materials-list` es nativa; su
//     ScreenInstance mínima existe solo para satisfacer la FK del mapping.
//   - No-regresión L1: la cadena user_roles → role_permissions →
//     permissions filtrando por viewer@edugo.demo sigue devolviendo
//     EXACTAMENTE el set {announcements:read}.
func (f *L3IsolationConstants) Apply(tx *gorm.DB, ctx *framework.ApplyContext) error {
	if err := f.verifyResourceMaterials(tx); err != nil {
		return err
	}
	if err := f.verifyMaterialsPermissions(tx); err != nil {
		return err
	}
	if err := f.verifyRolePermissions(tx); err != nil {
		return err
	}
	// Poda SDUI material (2026-06-07) + corrección F2 (2026-06-08): la
	// screen_instance `material-form` sigue ELIMINADA (sin mapping). La
	// `materials-list` se conserva MÍNIMA (no renderizada) por la FK del
	// mapping de menú. verifyResourceScreens valida 1 mapping (`list`
	// default); no se verifica la forma del slot_data (pantalla nativa).
	if err := f.verifyResourceScreens(tx); err != nil {
		return err
	}
	if err := f.verifyViewerPermissionsNoRegression(tx); err != nil {
		return err
	}

	ctx.SetConstant("E2EFixtureL3ResourceMaterialsID", layers.L3_RESOURCE_MATERIALS_ID)
	ctx.SetConstant("E2EFixtureL3ResourceMaterialsKey", layers.L3_RESOURCE_MATERIALS_KEY)
	ctx.SetConstant("E2EFixtureL3PermMaterialsReadID", layers.L3_PERM_MATERIALS_READ_ID)
	ctx.SetConstant("E2EFixtureL3PermMaterialsCreateID", layers.L3_PERM_MATERIALS_CREATE_ID)
	ctx.SetConstant("E2EFixtureL3PermMaterialsUpdateID", layers.L3_PERM_MATERIALS_UPDATE_ID)
	ctx.SetConstant("E2EFixtureL3ScreenInstanceMaterialsListID", layers.L3_SCREEN_INSTANCE_MATERIALS_LIST_ID)
	ctx.SetConstant("E2EFixtureL3ScreenInstanceMaterialFormID", layers.L3_SCREEN_INSTANCE_MATERIAL_FORM_ID)
	ctx.SetConstant("E2EFixtureL3ScreenKeyMaterialsList", layers.L3_SCREEN_KEY_MATERIALS_LIST)
	ctx.SetConstant("E2EFixtureL3ScreenKeyMaterialForm", layers.L3_SCREEN_KEY_MATERIAL_FORM)

	return nil
}

// Cleanup es no-op: las filas L3 son del system seed, no del scenario.
func (f *L3IsolationConstants) Cleanup(tx *gorm.DB, ctx *framework.ApplyContext) error {
	return nil
}

// verifyResourceMaterials cubre F5-REQ-1.1.
//
// Verifica que existe la fila en iam.resources con id=
// L3_RESOURCE_MATERIALS_ID y que sus columnas reflejan la spec
// (key=materials, display_name no vacío, scope=unit, parent_id NULL,
// is_menu_visible=true, is_active=true).
func (f *L3IsolationConstants) verifyResourceMaterials(tx *gorm.DB) error {
	type row struct {
		ID            string
		Key           string
		DisplayName   string
		Scope         string
		ParentID      *string
		IsMenuVisible bool
		IsActive      bool
	}
	const q = `
SELECT id::text          AS id,
       key               AS key,
       display_name      AS display_name,
       scope::text       AS scope,
       parent_id::text   AS parent_id,
       is_menu_visible   AS is_menu_visible,
       is_active         AS is_active
FROM iam.resources
WHERE id = ?::uuid
`
	var r row
	if err := tx.Raw(q, layers.L3_RESOURCE_MATERIALS_ID).Scan(&r).Error; err != nil {
		return fmt.Errorf("L3IsolationConstants: query iam.resources: %w", err)
	}
	if r.ID == "" {
		return fmt.Errorf(
			"L3IsolationConstants: fila L3 ausente en iam.resources id=%s — corré system.ApplySystem (con L3 registrada) antes del scenario",
			layers.L3_RESOURCE_MATERIALS_ID,
		)
	}
	if r.Key != layers.L3_RESOURCE_MATERIALS_KEY {
		return fmt.Errorf(
			"L3IsolationConstants: F5-REQ-1.1 violado — resources.key=%q, want %q",
			r.Key, layers.L3_RESOURCE_MATERIALS_KEY,
		)
	}
	if r.DisplayName == "" {
		return fmt.Errorf("L3IsolationConstants: F5-REQ-1.1 violado — resources.display_name vacío")
	}
	if r.Scope != "unit" {
		return fmt.Errorf(
			"L3IsolationConstants: F5-REQ-1.1 violado — resources.scope=%q, want %q",
			r.Scope, "unit",
		)
	}
	if r.ParentID != nil && *r.ParentID != "" {
		return fmt.Errorf(
			"L3IsolationConstants: F5-REQ-1.1 violado — resources.parent_id=%q, want NULL",
			*r.ParentID,
		)
	}
	if !r.IsMenuVisible {
		return fmt.Errorf("L3IsolationConstants: F5-REQ-1.1 violado — resources.is_menu_visible=false, want true")
	}
	if !r.IsActive {
		return fmt.Errorf("L3IsolationConstants: F5-REQ-1.1 violado — resources.is_active=false, want true")
	}
	return nil
}

// verifyMaterialsPermissions cubre F5-REQ-2.1.
//
// Asserta que iam.permissions con resource_id=materials contiene
// EXACTAMENTE el set {materials:create, materials:read, materials:update}
// ordenado alfabéticamente, y que NO existe la fila materials:delete
// en toda la tabla.
func (f *L3IsolationConstants) verifyMaterialsPermissions(tx *gorm.DB) error {
	const qList = `
SELECT name
FROM iam.permissions
WHERE resource_id = ?::uuid
ORDER BY name
`
	var names []string
	if err := tx.Raw(qList, layers.L3_RESOURCE_MATERIALS_ID).Scan(&names).Error; err != nil {
		return fmt.Errorf("L3IsolationConstants: query iam.permissions: %w", err)
	}
	want := []string{"content.materials.create", "content.materials.read", "content.materials.update"}
	if len(names) != len(want) {
		return fmt.Errorf(
			"L3IsolationConstants: F5-REQ-2.1 violado — permissions por resource=materials = %v, want %v",
			names, want,
		)
	}
	for i, w := range want {
		if names[i] != w {
			return fmt.Errorf(
				"L3IsolationConstants: F5-REQ-2.1 violado — permissions por resource=materials = %v, want %v",
				names, want,
			)
		}
	}

	// Aserción explícita: materials:delete NO existe en toda la tabla.
	const qDelete = `SELECT COUNT(*) FROM iam.permissions WHERE name = ?`
	var deleteCount int64
	if err := tx.Raw(qDelete, "content.materials.delete").Scan(&deleteCount).Error; err != nil {
		return fmt.Errorf("L3IsolationConstants: query materials:delete absence: %w", err)
	}
	if deleteCount != 0 {
		return fmt.Errorf(
			"L3IsolationConstants: F5-REQ-2.1 violado — materials:delete existe en iam.permissions (count=%d), debe ser 0 — L3 valida CRUD parcial sin :delete",
			deleteCount,
		)
	}
	return nil
}

// verifyRolePermissions cubre F5-REQ-2.2.
//
// Asserta que existen los 3 role_permissions super_admin × materials
// (read, create, update) y que NO existe role_permission super_admin ×
// materials:delete.
func (f *L3IsolationConstants) verifyRolePermissions(tx *gorm.DB) error {
	const qPresent = `
SELECT COUNT(*)
FROM iam.role_permissions
WHERE role_id = ?::uuid
  AND permission_id IN (?::uuid, ?::uuid, ?::uuid)
`
	var count int64
	if err := tx.Raw(
		qPresent,
		layers.L0_ROLE_SUPER_ADMIN_ID,
		layers.L3_PERM_MATERIALS_READ_ID,
		layers.L3_PERM_MATERIALS_CREATE_ID,
		layers.L3_PERM_MATERIALS_UPDATE_ID,
	).Scan(&count).Error; err != nil {
		return fmt.Errorf("L3IsolationConstants: query role_permissions super_admin×materials: %w", err)
	}
	if count != 3 {
		return fmt.Errorf(
			"L3IsolationConstants: F5-REQ-2.2 violado — role_permissions super_admin×materials = %d, want 3",
			count,
		)
	}

	// Aserción explícita: NO existe role_permission super_admin ×
	// materials:delete (la permission misma no existe, pero validamos
	// también el join para defensa en profundidad).
	const qDelete = `
SELECT COUNT(*)
FROM iam.role_permissions rp
JOIN iam.permissions p ON p.id = rp.permission_id
WHERE rp.role_id = ?::uuid
  AND p.name = ?
`
	var deleteCount int64
	if err := tx.Raw(qDelete, layers.L0_ROLE_SUPER_ADMIN_ID, "content.materials.delete").Scan(&deleteCount).Error; err != nil {
		return fmt.Errorf("L3IsolationConstants: query role_permissions super_admin×materials:delete: %w", err)
	}
	if deleteCount != 0 {
		return fmt.Errorf(
			"L3IsolationConstants: F5-REQ-2.2 violado — role_permission super_admin×materials:delete existe (count=%d), debe ser 0",
			deleteCount,
		)
	}
	return nil
}

// verifyResourceScreens cubre F5-REQ-3.3 (post-poda SDUI material).
//
// Verifica que ui_config.resource_screens tiene EXACTAMENTE 1 fila para
// resource_id=materials: screen_type=list, is_default=true (la pantalla
// es NATIVA, sin ScreenInstance). El mapping `form` fue podado junto con
// su ScreenInstance (poda SDUI material 2026-06-07). Aserción negativa
// explícita: NO existe fila screen_type=form.
func (f *L3IsolationConstants) verifyResourceScreens(tx *gorm.DB) error {
	const qTotal = `
SELECT COUNT(*)
FROM ui_config.resource_screens
WHERE resource_id = ?::uuid
`
	var total int64
	if err := tx.Raw(qTotal, layers.L3_RESOURCE_MATERIALS_ID).Scan(&total).Error; err != nil {
		return fmt.Errorf("L3IsolationConstants: query resource_screens count: %w", err)
	}
	if total != 1 {
		return fmt.Errorf(
			"L3IsolationConstants: F5-REQ-3.3 violado — resource_screens para resource=materials = %d, want 1 (solo list; form podado)",
			total,
		)
	}

	const qList = `
SELECT COUNT(*)
FROM ui_config.resource_screens
WHERE resource_id = ?::uuid
  AND screen_type = ?
  AND is_default = TRUE
`
	var listCount int64
	if err := tx.Raw(qList, layers.L3_RESOURCE_MATERIALS_ID, "list").Scan(&listCount).Error; err != nil {
		return fmt.Errorf("L3IsolationConstants: query resource_screens list: %w", err)
	}
	if listCount != 1 {
		return fmt.Errorf(
			"L3IsolationConstants: F5-REQ-3.3 violado — resource_screens (resource=materials, screen_type=list, is_default=true) = %d, want 1",
			listCount,
		)
	}

	const qForm = `
SELECT COUNT(*)
FROM ui_config.resource_screens
WHERE resource_id = ?::uuid
  AND screen_type = ?
`
	var formCount int64
	if err := tx.Raw(qForm, layers.L3_RESOURCE_MATERIALS_ID, "form").Scan(&formCount).Error; err != nil {
		return fmt.Errorf("L3IsolationConstants: query resource_screens form: %w", err)
	}
	if formCount != 0 {
		return fmt.Errorf(
			"L3IsolationConstants: F5-REQ-3.3 violado — resource_screens (resource=materials, screen_type=form) = %d, want 0 (mapping form podado)",
			formCount,
		)
	}
	return nil
}

// verifyViewerPermissionsNoRegression asegura que tras aplicar L3 la
// cadena L1 viewer@edugo.demo → user_role → role → role_permission →
// permission sigue devolviendo EXACTAMENTE {announcements:read}.
//
// No es un requirement explícito de Fase 5 (es no-regresión de
// F3-REQ-5.3) pero está pedido por el spec del scenario L3: una capa
// nueva no debe inflar accidentalmente los permisos del viewer ni
// concederle ningún permiso sobre materials (refuerza F5-REQ-2.3 a
// nivel SQL).
func (f *L3IsolationConstants) verifyViewerPermissionsNoRegression(tx *gorm.DB) error {
	const q = `
SELECT p.name
FROM auth.users u
JOIN iam.user_roles ur ON ur.user_id = u.id AND ur.is_active = TRUE
JOIN iam.roles r ON r.id = ur.role_id AND r.is_active = TRUE
JOIN iam.role_permissions rp ON rp.role_id = r.id
JOIN iam.permissions p ON p.id = rp.permission_id AND p.is_active = TRUE
WHERE u.email = ?
ORDER BY p.name
`
	var names []string
	if err := tx.Raw(q, layers.L1_VIEWER_EMAIL).Scan(&names).Error; err != nil {
		return fmt.Errorf("L3IsolationConstants: query viewer permissions: %w", err)
	}
	if len(names) != 1 || names[0] != "academic.announcements.read" {
		return fmt.Errorf(
			"L3IsolationConstants: no-regresión L1 violada — viewer %q tiene permisos %v, want exactamente [announcements:read] (L3 no debe filtrar permisos materials:* al viewer)",
			layers.L1_VIEWER_EMAIL, names,
		)
	}
	return nil
}
