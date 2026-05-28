package layers

import (
	"encoding/json"
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// materialsListSlotData es el slot_data JSON canónico de la pantalla
// L3 "materials-list". Coincide bit-a-bit con design §3 de
// phase-5-layer-l3 (F5-REQ-3.1).
//
// CRÍTICO: NO declara acción `delete`. L3 valida CRUD parcial
// (F5-REQ-2.1 / F5-REQ-3.1). El FE debe tratar `permissionFor(DELETE)`
// como `null` y no renderizar el botón.
//
// `api_prefix=academic` (no `platform`): valida que el FE usa el
// prefix correcto por pantalla, no un único prefix global (design §3).
//
// F3.1 (plan 004): migrada al patrón delta. Hereda los default_actions
// de list-basic-v1 (create/edit/delete con $resource$ resuelto a
// "content.materials") y suprime `delete` vía actions_removed para
// conservar el CRUD parcial (sin borrado) que valida L3. create/edit
// quedan con scope header/row y permisos idénticos al estado
// pre-migración (verificado por el harness).
const materialsListSlotData = `{
  "title": "Materiales",
  "actions_removed": ["delete"],
  "columns": [
    { "key": "title", "label": "Título" },
    { "key": "description", "label": "Descripción" }
  ],
  "api_prefix": "academic"
}`

// materialFormSlotData es el slot_data JSON canónico de la pantalla
// L3 "material-form" (F5-REQ-3.2). Campos: title, description, file_url.
//
// F3.1 (2ª pasada, plan 004): migrada al patrón delta. El legacy
// declaraba una sola acción `save` (permission=content.materials.update,
// event_id=submit-form). Hereda form-basic-v1 y suprime `save_new` y
// `delete` vía actions_removed para preservar EXACTAMENTE el contrato
// semántico previo: el único botón es {submit-form|content.materials.update}.
// Se mantiene así el CRUD parcial de materials (sin create-only ni delete;
// F5-REQ-3.2 prohíbe DELETE y el catálogo L3 no define
// content.materials.delete). El campo `condition` (presentación) lo
// aporta el template (`save` edit-only); el guard semántico solo protege
// {event_id, permission}, que queda idéntico (verificado por el harness).
const materialFormSlotData = `{
  "title": "Material",
  "page_title": "Material",
  "edit_title": "Editar material",
  "fields": [
    { "key": "title", "label": "Título", "type": "text", "required": true },
    { "key": "description", "label": "Descripción", "type": "textarea", "required": false },
    { "key": "file_url", "label": "URL del archivo", "type": "text", "required": false }
  ],
  "actions_removed": ["save_new", "delete"],
  "api_prefix": "academic"
}`

// applyL3Screens siembra las 2 ScreenInstances de L3
// (materials-list y material-form) en ui_config.screen_instances.
// Reutiliza los templates list-basic-v1 y form-basic-v1 sembrados
// por L0. Idempotente: UPSERT DoNothing sobre id.
//
// Scope `unit`: alineado con el scope del resource materials
// (F5-REQ-1.1). RequiredPermission solo en la list; el form sigue el
// patrón L2 (announcement-form usa announcements:read), por lo que
// material-form usa materials:read.
func applyL3Screens(tx *gorm.DB) error {
	descList := "Listado de materiales educativos"
	descForm := "Formulario de creación/edición de materiales"
	requiredPermissionList := "content.materials.read"
	requiredPermissionForm := "content.materials.read"

	instances := []entities.ScreenInstance{
		{
			ID:                 mustParseUUID(L3_SCREEN_INSTANCE_MATERIALS_LIST_ID, "L3_SCREEN_INSTANCE_MATERIALS_LIST_ID"),
			ScreenKey:          L3_SCREEN_KEY_MATERIALS_LIST,
			TemplateID:         mustParseUUID(L0_SCREEN_TPL_LIST_ID, "L0_SCREEN_TPL_LIST_ID"),
			Name:               "Listado de materiales",
			Description:        &descList,
			SlotData:           json.RawMessage([]byte(materialsListSlotData)),
			Scope:              "unit",
			RequiredPermission: &requiredPermissionList,
			HandlerKey:         nil,
			IsActive:           true,
		},
		{
			ID:                 mustParseUUID(L3_SCREEN_INSTANCE_MATERIAL_FORM_ID, "L3_SCREEN_INSTANCE_MATERIAL_FORM_ID"),
			ScreenKey:          L3_SCREEN_KEY_MATERIAL_FORM,
			TemplateID:         mustParseUUID(L0_SCREEN_TPL_FORM_ID, "L0_SCREEN_TPL_FORM_ID"),
			Name:               "Formulario de material",
			Description:        &descForm,
			SlotData:           json.RawMessage([]byte(materialFormSlotData)),
			Scope:              "unit",
			RequiredPermission: &requiredPermissionForm,
			HandlerKey:         nil,
			IsActive:           true,
		},
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&instances).Error; err != nil {
		return fmt.Errorf("applyL3Screens: upsert screen_instances: %w", err)
	}
	return nil
}
