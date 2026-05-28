package layers

import (
	"encoding/json"
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// announcementsListSlotData es el slot_data JSON canónico de la
// pantalla L0 "announcements-list". Coincide bit-a-bit con design §4
// de phase-2-layer-l0; cualquier cambio aquí altera el hash del seed
// (ver seeds/CLAUDE.md sobre bump de SchemaVersion).
// F3.1 (plan 004): migrada al patrón delta. Hereda los 3 default_actions
// de list-basic-v1 (create/edit/delete con $resource$ resuelto a
// "academic.announcements" desde required_permission). No declara
// actions_removed porque conserva las tres acciones. El label del
// template ("Nuevo"/"Editar"/"Eliminar") coincide con el legacy, así que
// la normalización es transparente; scope (header/row/row) y permisos
// quedan idénticos al estado pre-migración (verificado por el harness
// screen_actions_roundtrip_test.go).
const announcementsListSlotData = `{
  "title": "Anuncios",
  "page_title": "Anuncios",
  "filter_ready_label": "Fijados",
  "filter_processing_label": "No fijados",
  "columns": [
    { "key": "title", "label": "Título" },
    { "key": "created_at", "label": "Fecha" }
  ],
  "api_prefix": "platform"
}`

// listBasicV1Definition es el JSON canónico de la `definition` del
// template "list-basic-v1" v1. Excepción documentada a la regla
// general "no tocar L0..L3" prevista en phase-6-layer-l4/design.md
// (§4 tabla y motivación): el seed L0 original guardaba `{}` para
// los 3 templates base, lo que rompía el SDUI engine del KMP (exige
// el campo `zones` en `ScreenTemplate`, ver
// modules/sdui-engine/.../model/ScreenDefinition.kt:53). Estos JSON
// derivan del inventario legacy ([archivado pre-Fase-6] data.go:607)
// que era la forma canónica usada por los renderers KMP.
//
// default_actions[] es el contrato declarativo de actions por defecto
// del template. Se materializa por el composer en backend
// (api-platform/internal/core/usecase/screen_instance/compose.go); el
// frontend solo ve la lista final ya resuelta en slot_data.actions.
// El placeholder $resource$ se expande con el prefijo derivado de
// screen_instance.required_permission (p.ej. "academic.announcements").
const listBasicV1Definition = `{
  "navigation": {"topBar": {"title": "slot:page_title", "showBack": false}},
  "default_actions": [
    {"id": "create", "scope": "header", "label": "Nuevo",    "icon": "plus",   "permission": "$resource$.create", "condition": "always", "event_id": "create", "style": "icon",        "order": 10},
    {"id": "edit",   "scope": "row",    "label": "Editar",   "icon": "pencil", "permission": "$resource$.update", "condition": "always", "event_id": "edit",   "style": "icon",        "order": 10},
    {"id": "delete", "scope": "row",    "label": "Eliminar", "icon": "trash",  "permission": "$resource$.delete", "condition": "always", "event_id": "delete", "style": "destructive", "order": 20}
  ],
  "zones": [
    {"id": "search_zone", "type": "container", "slots": [
      {"id": "search_bar", "controlType": "search-bar", "bind": "slot:search_placeholder", "default": "Buscar..."}
    ]},
    {"id": "list_actions", "type": "action-group", "scope": "header", "slots": []},
    {"id": "row_actions",  "type": "action-group", "scope": "row",    "slots": []},
    {"id": "filters", "type": "container", "distribution": "flow-row", "slots": [
      {"id": "filter_all", "controlType": "chip", "bind": "slot:filter_all_label", "selected": true, "default": "Todos"},
      {"id": "filter_ready", "controlType": "chip", "bind": "slot:filter_ready_label", "default": "Activos"},
      {"id": "filter_processing", "controlType": "chip", "bind": "slot:filter_processing_label", "default": "Otros"}
    ]},
    {"id": "empty_state", "type": "container", "condition": "data.isEmpty", "slots": [
      {"id": "empty_icon", "controlType": "icon", "bind": "slot:empty_icon"},
      {"id": "empty_title", "controlType": "label", "style": "headline", "bind": "slot:empty_state_title", "default": "Sin resultados"},
      {"id": "empty_desc", "controlType": "label", "style": "body", "bind": "slot:empty_state_description", "default": "No hay datos para mostrar."},
      {"id": "empty_action", "controlType": "filled-button", "bind": "slot:empty_action_label", "event_id": "create", "default": "Crear el primero"}
    ]},
    {"id": "list_content", "type": "simple-list", "condition": "!data.isEmpty", "itemLayout": {"slots": [
      {"id": "item_icon", "controlType": "icon", "field": "file_type_icon"},
      {"id": "item_title", "controlType": "label", "style": "headline-small", "field": "title"},
      {"id": "item_subtitle", "controlType": "label", "style": "body-small", "field": "subtitle"},
      {"id": "item_status", "controlType": "chip", "field": "status"},
      {"id": "item_date", "controlType": "label", "style": "caption", "field": "created_at"}
    ]}}
  ]
}`

// detailBasicV1Definition — ver doc de listBasicV1Definition.
// Derivado de [archivado pre-Fase-6] data.go:608.
const detailBasicV1Definition = `{
  "navigation": {"topBar": {"title": "slot:page_title", "showBack": true, "actions": []}},
  "zones": [
    {"id": "hero", "type": "container", "slots": [
      {"id": "file_type_icon", "controlType": "icon", "style": "large", "field": "file_type"},
      {"id": "status_badge", "controlType": "chip", "field": "status"}
    ]},
    {"id": "header", "type": "container", "slots": [
      {"id": "title", "controlType": "label", "style": "headline-large", "field": "title"},
      {"id": "subject", "controlType": "label", "style": "body", "field": "subject"},
      {"id": "grade", "controlType": "label", "style": "body-small", "field": "grade"}
    ]},
    {"id": "details", "type": "container", "slots": [
      {"id": "file_size", "controlType": "label", "bind": "slot:file_size_label", "field": "file_size_display", "default": "Tamaño"},
      {"id": "uploaded_date", "controlType": "label", "bind": "slot:uploaded_label", "field": "created_at", "default": "Subido"},
      {"id": "status", "controlType": "label", "bind": "slot:status_label", "field": "status", "default": "Estado"}
    ]},
    {"id": "description", "type": "container", "slots": [
      {"id": "section_title", "controlType": "label", "style": "title-medium", "bind": "slot:description_title", "default": "Descripción"},
      {"id": "description_text", "controlType": "label", "style": "body", "field": "description"}
    ]},
    {"id": "summary", "type": "container", "condition": "data.summary != null", "slots": [
      {"id": "summary_title", "controlType": "label", "style": "title-medium", "bind": "slot:summary_title", "default": "Resumen"},
      {"id": "summary_content", "controlType": "label", "style": "body", "field": "summary.main_ideas"}
    ]},
    {"id": "actions", "type": "action-group", "slots": [
      {"id": "download_btn", "controlType": "filled-button", "bind": "slot:download_label", "icon": "download", "default": "Descargar"},
      {"id": "take_quiz_btn", "controlType": "outlined-button", "bind": "slot:quiz_label", "icon": "quiz", "default": "Hacer quiz"}
    ]}
  ]
}`

// formBasicV1Definition — ver doc de listBasicV1Definition.
// Derivado de [archivado pre-Fase-6] data.go:610.
//
// La zona `form_header` (label "Formulario" + texto "Completa los campos.")
// se eliminó: el TopBar ya muestra el page_title de cada form, por lo que
// el header redundante solo ocupaba espacio sin aportar contexto.
//
// La zona `form_actions` usa scope=form-submit (separación semántica del
// snapshot 002): solo botones save/save_new/delete. La zona
// `resource_toolbar` (scope=resource-toolbar) queda intencionalmente
// ausente en form-basic-v1 — los forms simples no exhiben acciones de
// recurso. master-detail-v1 sí declara ambas zonas.
const formBasicV1Definition = `{
  "navigation": {"topBar": {"title": "slot:page_title", "showBack": true}},
  "default_actions": [
    {"id": "save_new", "scope": "form-submit", "label": "Guardar",  "icon": "save",  "permission": "$resource$.create", "condition": "create-only", "event_id": "submit-form", "style": "filled",      "order": 10},
    {"id": "save",     "scope": "form-submit", "label": "Guardar",  "icon": "save",  "permission": "$resource$.update", "condition": "edit-only",   "event_id": "submit-form", "style": "filled",      "order": 10},
    {"id": "delete",   "scope": "form-submit", "label": "Eliminar", "icon": "trash", "permission": "$resource$.delete", "condition": "edit-only",   "event_id": "delete",      "style": "destructive", "order": 20}
  ],
  "zones": [
    {"id": "form_fields", "type": "form-section", "slots": []},
    {"id": "form_submit", "type": "action-group", "scope": "form-submit", "layout_strategy": "row", "slots": []}
  ]
}`

// masterDetailV1Definition es el template "master-detail-v1": form
// con CRUD estándar (save_new/save/delete con scope=form-submit) más
// una zona resource-toolbar para acciones del recurso en estado edit
// (detail/publish/archive/etc., declaradas por la instancia via
// actions_added). La instancia debe declarar también `detail_config`
// con la screen_key destino y los params para el panel de detalle
// (el frontend KMP es quien interpreta detail_config; el backend solo
// lo persiste como blob dentro de slot_data).
//
// layout_strategy="flow-row" + overflow_threshold=3 en resource-toolbar:
// hint declarativo para el renderer KMP. El backend solo lo persiste
// dentro del JSON; no muta el entity Zone de Go.
const masterDetailV1Definition = `{
  "navigation": {"topBar": {"title": "slot:page_title", "showBack": true}},
  "default_actions": [
    {"id": "save_new", "scope": "form-submit",     "label": "Guardar",  "icon": "save",  "permission": "$resource$.create", "condition": "create-only", "event_id": "submit-form", "style": "filled",      "order": 10},
    {"id": "save",     "scope": "form-submit",     "label": "Guardar",  "icon": "save",  "permission": "$resource$.update", "condition": "edit-only",   "event_id": "submit-form", "style": "filled",      "order": 10},
    {"id": "delete",   "scope": "form-submit",     "label": "Eliminar", "icon": "trash", "permission": "$resource$.delete", "condition": "edit-only",   "event_id": "delete",      "style": "destructive", "order": 20},
    {"id": "detail",   "scope": "resource-toolbar","label": "Detalle",  "icon": "list",  "permission": "$resource$.read",   "condition": "edit-only",   "event_id": "view-detail", "style": "icon",        "order": 10}
  ],
  "zones": [
    {"id": "resource_toolbar", "type": "action-group", "scope": "resource-toolbar", "layout_strategy": "flow-row", "overflow_threshold": 3, "slots": []},
    {"id": "form_fields",      "type": "form-section", "slots": []},
    {"id": "form_submit",      "type": "action-group", "scope": "form-submit",      "layout_strategy": "row",     "slots": []}
  ]
}`

// applyL0Screens siembra las 3 ScreenTemplates compartidas (list,
// detail, form) v1, la ScreenInstance "announcements-list" y el
// mapping ResourceScreen announcements↔announcements-list de L0.
// Idempotente: UPSERT DoNothing por las claves naturales del schema.
func applyL0Screens(tx *gorm.DB) error {
	if err := upsertL0ScreenTemplates(tx); err != nil {
		return err
	}
	if err := upsertL0ScreenInstances(tx); err != nil {
		return err
	}
	if err := upsertL0ResourceScreens(tx); err != nil {
		return err
	}
	return nil
}

// upsertL0ScreenTemplates inserta las 3 templates base (list, detail,
// form) en ui_config.screen_templates. Conflict target (name, version):
// las templates son compartidas entre L0..L4 y no deben re-pisarse si
// otra capa ya las insertó.
//
// Definition: cada template embebe su JSON canónico de `zones`
// (list/detail/formBasicV1Definition). Esto es excepción documentada
// a la regla "no tocar L0..L3" prevista en
// phase-6-layer-l4/design.md (§4): el seed original usaba `{}` y
// rompía el SDUI engine del KMP (que valida el campo `zones`).
func upsertL0ScreenTemplates(tx *gorm.DB) error {
	templates := []entities.ScreenTemplate{
		{
			ID:         mustParseUUID(L0_SCREEN_TPL_LIST_ID, "L0_SCREEN_TPL_LIST_ID"),
			Pattern:    "list",
			Name:       "list-basic-v1",
			Version:    1,
			Definition: json.RawMessage([]byte(listBasicV1Definition)),
			IsActive:   true,
		},
		{
			ID:         mustParseUUID(L0_SCREEN_TPL_DETAIL_ID, "L0_SCREEN_TPL_DETAIL_ID"),
			Pattern:    "detail",
			Name:       "detail-basic-v1",
			Version:    1,
			Definition: json.RawMessage([]byte(detailBasicV1Definition)),
			IsActive:   true,
		},
		{
			ID:         mustParseUUID(L0_SCREEN_TPL_FORM_ID, "L0_SCREEN_TPL_FORM_ID"),
			Pattern:    "form",
			Name:       "form-basic-v1",
			Version:    1,
			Definition: json.RawMessage([]byte(formBasicV1Definition)),
			IsActive:   true,
		},
		{
			ID:         mustParseUUID(L0_SCREEN_TPL_MASTER_DETAIL_ID, "L0_SCREEN_TPL_MASTER_DETAIL_ID"),
			Pattern:    "master-detail",
			Name:       "master-detail-v1",
			Version:    1,
			Definition: json.RawMessage([]byte(masterDetailV1Definition)),
			IsActive:   true,
		},
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}, {Name: "version"}},
		DoNothing: true,
	}).Create(&templates).Error; err != nil {
		return fmt.Errorf("applyL0Screens: upsert screen_templates: %w", err)
	}
	return nil
}

// upsertL0ScreenInstances inserta la instancia "announcements-list"
// en ui_config.screen_instances. Conflict target screen_key (índice
// único). slot_data se inyecta como json.RawMessage para preservar
// el formato exacto del literal definido arriba.
func upsertL0ScreenInstances(tx *gorm.DB) error {
	description := "Listado de anuncios institucionales"
	requiredPermission := "academic.announcements.read"

	instance := entities.ScreenInstance{
		ID:                 mustParseUUID(L0_SCREEN_INST_ANNOUNCEMENTS_LIST_ID, "L0_SCREEN_INST_ANNOUNCEMENTS_LIST_ID"),
		ScreenKey:          L0_SCREEN_KEY_ANNOUNCEMENTS_LIST,
		TemplateID:         mustParseUUID(L0_SCREEN_TPL_LIST_ID, "L0_SCREEN_TPL_LIST_ID"),
		Name:               "Anuncios — Listado",
		Description:        &description,
		SlotData:           json.RawMessage([]byte(announcementsListSlotData)),
		Scope:              "school",
		RequiredPermission: &requiredPermission,
		HandlerKey:         nil,
		IsActive:           true,
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "screen_key"}},
		DoNothing: true,
	}).Create(&instance).Error; err != nil {
		return fmt.Errorf("applyL0Screens: upsert screen_instances: %w", err)
	}
	return nil
}

// upsertL0ResourceScreens vincula el resource "announcements" con la
// instance "announcements-list" como pantalla default de tipo "list".
// ID derivado determinísticamente vía SHA1 sobre (resource_id,
// screen_type) para que sea estable entre re-seeds e impactando el
// hash de exports JSON. Conflict target (resource_id, screen_type).
func upsertL0ResourceScreens(tx *gorm.DB) error {
	resourceID := mustParseUUID(L0_RESOURCE_ANNOUNCEMENTS_ID, "L0_RESOURCE_ANNOUNCEMENTS_ID")
	screenType := "list"
	id := uuid.NewSHA1(uuid.NameSpaceOID, []byte(resourceID.String()+":"+screenType))

	mapping := entities.ResourceScreen{
		ID:          id,
		ResourceID:  resourceID,
		ResourceKey: L0_RESOURCE_ANNOUNCEMENTS_KEY,
		ScreenKey:   L0_SCREEN_KEY_ANNOUNCEMENTS_LIST,
		ScreenType:  screenType,
		IsDefault:   true,
		SortOrder:   0,
		IsActive:    true,
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "resource_id"}, {Name: "screen_type"}},
		DoNothing: true,
	}).Create(&mapping).Error; err != nil {
		return fmt.Errorf("applyL0Screens: upsert resource_screens: %w", err)
	}
	return nil
}

// mustParseUUID convierte un literal de UUID definido en
// l0_constants.go a uuid.UUID o devuelve panic. Los literales son
// hardcodeados y validados en compile-time-ish; un fallo aquí indica
// corrupción del archivo de constantes, no un error runtime real.
func mustParseUUID(s, name string) uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		panic(fmt.Sprintf("applyL0Screens: parse %s=%q: %v", name, s, err))
	}
	return u
}
