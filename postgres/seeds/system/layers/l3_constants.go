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
//     en el menú vía el mapping `materials:list`.
//   - 1.4.0: corrección F2 (2026-06-08) — la poda 1.3.0 quitó de más:
//     `materials-list` tiene mapping resource_screen y la FK
//     fk_resource_screens_screen_key exige su screen_instance, así que un
//     recreate limpio fallaba en L3 (23503). Se RESTAURA la screen_instance
//     MÍNIMA `materials-list` (no se renderiza; pantalla NATIVA), mismo
//     patrón que batch-enroll/join-requests-inbox en L4. `material-form`
//     SIGUE PODADO (sin mapping → sin FK). +1 fila screen_instances.
//   - 1.4.1 (2026-06-19, bug 0069): los 3 permisos materials del contrato L3
//     se siembran con IsSystem=true (apply + accessor espejo).
const L3_SEED_VERSION = "1.4.1"

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

	// Identificadores de pantalla de material.
	//   - materials-list (ID ...0008): RESTAURADA en F2 (2026-06-08) como
	//     screen_instance MÍNIMA (no se renderiza; pantalla NATIVA Compose).
	//     Existe SOLO para satisfacer la FK fk_resource_screens_screen_key
	//     del mapping de menú `materials:list`. La poda 2026-06-07 la había
	//     eliminado por error y rompía el recreate limpio (23503).
	//     Ver l3_screens.go.
	//   - material-form (ID ...0009): SIGUE PODADA — no tiene mapping
	//     resource_screen → no hay FK que satisfacer. La constante se
	//     conserva solo como identificador de referencia (export E2E).
	L3_SCREEN_INSTANCE_MATERIALS_LIST_ID = "b3000000-0000-0000-0000-000000000008"
	L3_SCREEN_INSTANCE_MATERIAL_FORM_ID  = "b3000000-0000-0000-0000-000000000009"
	L3_SCREEN_KEY_MATERIALS_LIST         = "materials-list"
	L3_SCREEN_KEY_MATERIAL_FORM          = "material-form"
)
