package fixtures

import (
	"encoding/json"
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
		Description: "Verifica filas L3 sembradas por system.ApplySystem (resource materials + 3 permisos sin :delete + 2 ScreenInstances + 2 ResourceScreens) y exporta sus identificadores al JSON.",
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
//   - F5-REQ-3.1: ScreenInstance materials-list con slot_data válido —
//     2 actions (create, edit) sin delete, 2 columns, api_prefix=academic.
//   - F5-REQ-3.2: ScreenInstance material-form con slot_data válido —
//     3 fields (title, description, file_url), 2 actions
//     (SAVE_NEW → :create, SAVE_EXISTING → :update) sin DELETE,
//     api_prefix=academic.
//   - F5-REQ-3.3: 2 resource_screens (list default, form no-default).
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
	if err := f.verifyMaterialsListScreen(tx); err != nil {
		return err
	}
	if err := f.verifyMaterialFormScreen(tx); err != nil {
		return err
	}
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

// verifyMaterialsListScreen cubre F5-REQ-3.1.
//
// Verifica que existe la fila en ui_config.screen_instances con
// id=L3_SCREEN_INSTANCE_MATERIALS_LIST_ID, screen_key="materials-list",
// y que su slot_data es JSON válido con:
//   - actions[] EXACTAMENTE de tamaño 2 con ids "create" y "edit"
//     (no "delete" ni permission "content.materials.delete").
//   - columns[] de tamaño 2 (title, description).
//   - api_prefix="academic".
func (f *L3IsolationConstants) verifyMaterialsListScreen(tx *gorm.DB) error {
	type row struct {
		ID        string
		ScreenKey string
		Name      string
		SlotData  []byte
	}
	const q = `
SELECT id::text   AS id,
       screen_key AS screen_key,
       name       AS name,
       slot_data  AS slot_data
FROM ui_config.screen_instances
WHERE id = ?::uuid
`
	var r row
	if err := tx.Raw(q, layers.L3_SCREEN_INSTANCE_MATERIALS_LIST_ID).Scan(&r).Error; err != nil {
		return fmt.Errorf("L3IsolationConstants: query screen_instances materials-list: %w", err)
	}
	if r.ID == "" {
		return fmt.Errorf(
			"L3IsolationConstants: ScreenInstance L3 materials-list ausente id=%s — corré system.ApplySystem (con L3 registrada) antes del scenario",
			layers.L3_SCREEN_INSTANCE_MATERIALS_LIST_ID,
		)
	}
	if r.ScreenKey != layers.L3_SCREEN_KEY_MATERIALS_LIST {
		return fmt.Errorf(
			"L3IsolationConstants: F5-REQ-3.1 violado — screen_key=%q, want %q",
			r.ScreenKey, layers.L3_SCREEN_KEY_MATERIALS_LIST,
		)
	}
	if r.Name == "" {
		return fmt.Errorf("L3IsolationConstants: F5-REQ-3.1 violado — screen_instances.name vacío")
	}

	type action struct {
		ID         string `json:"id"`
		Permission string `json:"permission"`
	}
	type column struct {
		Key string `json:"key"`
	}
	var slot struct {
		Actions   []action `json:"actions"`
		Columns   []column `json:"columns"`
		APIPrefix string   `json:"api_prefix"`
	}
	if err := json.Unmarshal(r.SlotData, &slot); err != nil {
		return fmt.Errorf("L3IsolationConstants: F5-REQ-3.1 violado — slot_data materials-list no es JSON válido: %w", err)
	}

	if got := len(slot.Actions); got != 2 {
		return fmt.Errorf("L3IsolationConstants: F5-REQ-3.1 violado — materials-list.actions[]=%d, want 2", got)
	}
	wantActions := map[string]bool{"create": false, "edit": false}
	for _, a := range slot.Actions {
		if a.ID == "delete" {
			return fmt.Errorf(
				"L3IsolationConstants: F5-REQ-3.1 violado — materials-list contiene action id=%q, prohibido por design (sin :delete)",
				a.ID,
			)
		}
		if a.Permission == "content.materials.delete" {
			return fmt.Errorf(
				"L3IsolationConstants: F5-REQ-3.1 violado — materials-list contiene action con permission=materials:delete, prohibido",
			)
		}
		if _, expected := wantActions[a.ID]; expected {
			wantActions[a.ID] = true
		}
	}
	for k, seen := range wantActions {
		if !seen {
			return fmt.Errorf("L3IsolationConstants: F5-REQ-3.1 violado — materials-list.actions[].id=%q ausente", k)
		}
	}

	if got := len(slot.Columns); got != 2 {
		return fmt.Errorf("L3IsolationConstants: F5-REQ-3.1 violado — materials-list.columns[]=%d, want 2 (title, description)", got)
	}
	wantCols := map[string]bool{"title": false, "description": false}
	for _, c := range slot.Columns {
		if _, expected := wantCols[c.Key]; expected {
			wantCols[c.Key] = true
		}
	}
	for k, seen := range wantCols {
		if !seen {
			return fmt.Errorf("L3IsolationConstants: F5-REQ-3.1 violado — materials-list.columns[].key=%q ausente", k)
		}
	}

	if slot.APIPrefix != "academic" {
		return fmt.Errorf(
			"L3IsolationConstants: F5-REQ-3.1 violado — materials-list.api_prefix=%q, want %q (design §3: prefix por pantalla)",
			slot.APIPrefix, "academic",
		)
	}
	return nil
}

// verifyMaterialFormScreen cubre F5-REQ-3.2.
//
// Verifica que existe la fila en ui_config.screen_instances con
// id=L3_SCREEN_INSTANCE_MATERIAL_FORM_ID, screen_key="material-form",
// y que su slot_data es JSON válido con:
//   - fields[] de tamaño 3 con keys title, description, file_url.
//   - actions[] EXACTAMENTE de tamaño 2: SAVE_NEW → :create,
//     SAVE_EXISTING → :update. Sin event=DELETE.
//   - api_prefix="academic".
func (f *L3IsolationConstants) verifyMaterialFormScreen(tx *gorm.DB) error {
	type row struct {
		ID        string
		ScreenKey string
		SlotData  []byte
	}
	const q = `
SELECT id::text   AS id,
       screen_key AS screen_key,
       slot_data  AS slot_data
FROM ui_config.screen_instances
WHERE id = ?::uuid
`
	var r row
	if err := tx.Raw(q, layers.L3_SCREEN_INSTANCE_MATERIAL_FORM_ID).Scan(&r).Error; err != nil {
		return fmt.Errorf("L3IsolationConstants: query screen_instances material-form: %w", err)
	}
	if r.ID == "" {
		return fmt.Errorf(
			"L3IsolationConstants: ScreenInstance L3 material-form ausente id=%s — corré system.ApplySystem (con L3 registrada) antes del scenario",
			layers.L3_SCREEN_INSTANCE_MATERIAL_FORM_ID,
		)
	}
	if r.ScreenKey != layers.L3_SCREEN_KEY_MATERIAL_FORM {
		return fmt.Errorf(
			"L3IsolationConstants: F5-REQ-3.2 violado — screen_key=%q, want %q",
			r.ScreenKey, layers.L3_SCREEN_KEY_MATERIAL_FORM,
		)
	}

	type field struct {
		Key string `json:"key"`
	}
	type action struct {
		ID         string `json:"id"`
		Event      string `json:"event"`
		Permission string `json:"permission"`
	}
	var slot struct {
		Fields    []field  `json:"fields"`
		Actions   []action `json:"actions"`
		APIPrefix string   `json:"api_prefix"`
	}
	if err := json.Unmarshal(r.SlotData, &slot); err != nil {
		return fmt.Errorf("L3IsolationConstants: F5-REQ-3.2 violado — slot_data material-form no es JSON válido: %w", err)
	}

	if got := len(slot.Fields); got != 3 {
		return fmt.Errorf("L3IsolationConstants: F5-REQ-3.2 violado — material-form.fields[]=%d, want 3 (title, description, file_url)", got)
	}
	wantFields := map[string]bool{"title": false, "description": false, "file_url": false}
	for _, fd := range slot.Fields {
		if _, expected := wantFields[fd.Key]; expected {
			wantFields[fd.Key] = true
		}
	}
	for k, seen := range wantFields {
		if !seen {
			return fmt.Errorf("L3IsolationConstants: F5-REQ-3.2 violado — material-form.fields[].key=%q ausente", k)
		}
	}

	if got := len(slot.Actions); got != 2 {
		return fmt.Errorf("L3IsolationConstants: F5-REQ-3.2 violado — material-form.actions[]=%d, want 2 (SAVE_NEW, SAVE_EXISTING)", got)
	}

	wantPerms := map[string]string{
		"SAVE_NEW":      "content.materials.create",
		"SAVE_EXISTING": "content.materials.update",
	}
	seenPerms := map[string]string{}
	for _, a := range slot.Actions {
		if a.Event == "DELETE" {
			return fmt.Errorf(
				"L3IsolationConstants: F5-REQ-3.2 violado — material-form contiene action con event=DELETE, prohibido por design",
			)
		}
		if _, expected := wantPerms[a.Event]; expected {
			seenPerms[a.Event] = a.Permission
		}
	}
	for event, want := range wantPerms {
		got, ok := seenPerms[event]
		if !ok {
			return fmt.Errorf("L3IsolationConstants: F5-REQ-3.2 violado — material-form.actions sin event=%q", event)
		}
		if got != want {
			return fmt.Errorf(
				"L3IsolationConstants: F5-REQ-3.2 violado — material-form action event=%q permission=%q, want %q",
				event, got, want,
			)
		}
	}

	if slot.APIPrefix != "academic" {
		return fmt.Errorf(
			"L3IsolationConstants: F5-REQ-3.2 violado — material-form.api_prefix=%q, want %q",
			slot.APIPrefix, "academic",
		)
	}
	return nil
}

// verifyResourceScreens cubre F5-REQ-3.3.
//
// Verifica que ui_config.resource_screens tiene EXACTAMENTE 2 filas
// para resource_id=materials: una con screen_type=list, is_default=true,
// y otra con screen_type=form, is_default=false.
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
	if total != 2 {
		return fmt.Errorf(
			"L3IsolationConstants: F5-REQ-3.3 violado — resource_screens para resource=materials = %d, want 2 (list+form)",
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
  AND is_default = FALSE
`
	var formCount int64
	if err := tx.Raw(qForm, layers.L3_RESOURCE_MATERIALS_ID, "form").Scan(&formCount).Error; err != nil {
		return fmt.Errorf("L3IsolationConstants: query resource_screens form: %w", err)
	}
	if formCount != 1 {
		return fmt.Errorf(
			"L3IsolationConstants: F5-REQ-3.3 violado — resource_screens (resource=materials, screen_type=form, is_default=false) = %d, want 1",
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
