package layers

// L2_SEED_VERSION declara la versión semántica del contenido de L2.
// Bumpear en CADA cambio de dato visible.
//
// Historial:
//   - 1.0.0: versión inicial (3 fields: title, body, published_at).
//   - 1.1.0: agregado field `scope` (select school|unit, required) al
//     slot_data de announcement-form. Alinea el form con
//     CreateAnnouncementRequest.Scope (binding=required,oneof=school
//     unit). Sin él, el POST /api/v1/announcements devolvía 400.
//   - 1.2.0: el botón Guardar se desdobla en save_new (create-only,
//     permission=create) y save (edit-only, permission=update). Antes el
//     único slot pedía `update` siempre, lo que ocultaba el botón a
//     usuarios con solo `create` (caso focal-author).
//   - 1.3.0: F3 (plan 004) — announcement-form migrada a patrón delta
//     (hereda form-basic-v1; sin cambio semántico).
const L2_SEED_VERSION = "1.3.0"

// L2_LAYER_NAME es el nombre canónico de la capa, usado por
// --seed-up-to-layer y por logs.
const L2_LAYER_NAME = "L2"

// UUIDs hardcodeados de las entidades canónicas de L2 (prefijo b2xxx).
// Razón: tests E2E y JSON exportado a KMP los referencian por UUID.
const (
	L2_SCREEN_INSTANCE_ANNOUNCEMENT_FORM_ID  = "b2000000-0000-0000-0000-000000000001"
	L2_RESOURCE_SCREEN_ANNOUNCEMENTS_FORM_ID = "b2000000-0000-0000-0000-000000000002"
)

// Strings semánticos reutilizables.
const (
	L2_SCREEN_KEY_ANNOUNCEMENT_FORM = "announcement-form"
)
