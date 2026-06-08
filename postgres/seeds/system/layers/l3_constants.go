package layers

// L3_SEED_VERSION declara la versión semántica del contenido de L3.
// Bumpear en CADA cambio de dato visible.
//
// Historial:
//   - 1.1.0: versión anterior.
//   - 1.2.0: F3 (plan 004) — materials-list/material-form migradas a
//     patrón delta (actions_removed; sin cambio semántico).
//   - 1.3.0: poda SDUI material (2026-06-07) — eliminadas las 2
//     ScreenInstances (materials-list, material-form) + slot_data y el
//     mapping resource_screen `form`. Las pantallas de material son
//     nativas; esos seeds eran código muerto. El recurso materials sigue
//     en el menú vía el mapping `materials:list` (sin ScreenInstance).
const L3_SEED_VERSION = "1.3.0"

// L3_LAYER_NAME es el nombre canónico de la capa, usado por
// --seed-up-to-layer y por logs.
const L3_LAYER_NAME = "L3"

// UUIDs hardcodeados de las entidades canónicas de L3 (prefijo b3xxx).
// Razón: tests E2E y JSON exportado a KMP los referencian por UUID.
// Política ADR-6/7: UUIDs propios, no coordinados con legacy.
const (
	L3_RESOURCE_MATERIALS_ID  = "b3000000-0000-0000-0000-000000000001"
	L3_RESOURCE_MATERIALS_KEY = "materials"

	// NOTA: NO existe materials:delete por design (F5-REQ-2.1).
	// L3 valida CRUD parcial: read/create/update solamente.
	L3_PERM_MATERIALS_READ_ID   = "b3000000-0000-0000-0000-000000000002"
	L3_PERM_MATERIALS_CREATE_ID = "b3000000-0000-0000-0000-000000000003"
	L3_PERM_MATERIALS_UPDATE_ID = "b3000000-0000-0000-0000-000000000004"

	// Identificadores de pantalla de material. Poda SDUI material
	// (2026-06-07): las ScreenInstances correspondientes (IDs ...0008 /
	// ...0009) fueron ELIMINADAS — ya no respaldan filas en
	// ui_config.screen_instances. Se conservan estas constantes solo
	// como identificadores de referencia para el export E2E
	// (e2e/fixtures/l3_constants_export.go) y para el screen_key del
	// mapping resource_screen `materials:list` (pantalla NATIVA).
	L3_SCREEN_INSTANCE_MATERIALS_LIST_ID = "b3000000-0000-0000-0000-000000000008"
	L3_SCREEN_INSTANCE_MATERIAL_FORM_ID  = "b3000000-0000-0000-0000-000000000009"
	L3_SCREEN_KEY_MATERIALS_LIST         = "materials-list"
	L3_SCREEN_KEY_MATERIAL_FORM          = "material-form"
)
