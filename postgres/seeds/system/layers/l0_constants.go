package layers

// L0_SEED_VERSION declara la versión semántica del contenido de L0.
// Bumpear en CADA cambio de dato visible.
//
// Historial reciente:
//   - 1.2.0: announcementsListSlotData agrega filter_ready_label="Fijados"
//     y filter_processing_label="No fijados" para overridear los defaults
//     "Activos"/"Otros" del template list-basic-v1, que eran engañosos en
//     el contexto de anuncios (la entidad no tiene is_active).
//   - 1.3.0: ajustes de UI base —
//   - announcementsListSlotData incluye page_title="Anuncios" para que
//     el TopBar muestre título (antes solo tenía "title", que el
//     renderer ignora).
//   - formBasicV1Definition elimina la zona form_header (label
//     "Formulario" + texto "Completa los campos."). El TopBar ya
//     muestra el page_title; el header redundante ocupaba espacio
//     sin aportar contexto.
//   - 1.4.0: SDUI composer arquitectónico (Fase 3a) —
//   - listBasicV1Definition y formBasicV1Definition ganan
//     `default_actions[]` con placeholders `$resource$` que el
//     composer (api-platform) expande según
//     screen_instance.required_permission.
//   - formBasicV1Definition renombra scope `form` → `form-submit`
//     (separación semántica del snapshot 002) y declara
//     `layout_strategy: "row"` en la zona form_submit.
//   - Nuevo template `master-detail-v1` (pattern=master-detail) con
//     dos zonas action-group: form-submit (row) + resource-toolbar
//     (flow-row, overflow_threshold=3). Hereda los 3 defaults del
//     form más `detail` (scope=resource-toolbar, condition=edit-only).
//     Las instancias declaran `detail_config` para apuntar al panel
//     de detalle.
//   - 1.5.0: F3 (plan 004-permisologia-mvp) — patrón delta SDUI. La
//     instancia announcements-list deja de re-listar el array actions;
//     hereda default_actions del template list-basic-v1 vía el composer.
//     Sin cambio semántico (mismo set {event_id,permission} compuesto).
const L0_SEED_VERSION = "1.5.0"

// L0_LAYER_NAME es el nombre canónico de la capa, usado por
// --seed-up-to-layer y por logs.
const L0_LAYER_NAME = "L0-minimal"

// UUIDs hardcodeados de las entidades canónicas de L0.
// Razón: tests E2E y JSON exportado a KMP los referencian por UUID.
// Política ADR-6: UUIDs propios (no coordinados con legacy, que ya
// no se aplica desde Fase 2).
const (
	L0_RESOURCE_ANNOUNCEMENTS_ID = "b0000000-0000-0000-0000-000000000001"

	L0_ROLE_SUPER_ADMIN_ID = "10000000-0000-0000-0000-000000000001"

	L0_PERM_ANNOUNCEMENTS_READ   = "20000000-0000-0000-0000-000000000001"
	L0_PERM_ANNOUNCEMENTS_CREATE = "20000000-0000-0000-0000-000000000002"
	L0_PERM_ANNOUNCEMENTS_UPDATE = "20000000-0000-0000-0000-000000000003"
	L0_PERM_ANNOUNCEMENTS_DELETE = "20000000-0000-0000-0000-000000000004"

	L0_SCREEN_TPL_LIST_ID          = "30000000-0000-0000-0000-000000000001"
	L0_SCREEN_TPL_DETAIL_ID        = "30000000-0000-0000-0000-000000000002"
	L0_SCREEN_TPL_FORM_ID          = "30000000-0000-0000-0000-000000000003"
	L0_SCREEN_TPL_MASTER_DETAIL_ID = "30000000-0000-0000-0000-000000000004"

	L0_SCREEN_INST_ANNOUNCEMENTS_LIST_ID = "40000000-0000-0000-0000-000000000001"

	L0_USER_SUPER_ADMIN_ID = "50000000-0000-0000-0000-000000000001"
)

// Credenciales de bootstrapping del super_admin de L0.
//
// SEGURIDAD: Estas credenciales son SÓLO para arranque inicial del
// sistema y entornos de prueba. En cualquier deployment cloud,
// rotar antes del primer login productivo. Procedimiento de rotación
// documentado en seeds/CLAUDE.md.
const (
	L0_SUPER_ADMIN_EMAIL    = "super_admin@edugo.system"
	L0_SUPER_ADMIN_PASSWORD = "ChangeMe!2026"
)

// Strings semánticos reutilizables.
const (
	L0_RESOURCE_ANNOUNCEMENTS_KEY    = "announcements"
	L0_ROLE_SUPER_ADMIN_NAME         = "super_admin"
	L0_SCREEN_KEY_ANNOUNCEMENTS_LIST = "announcements-list"
)
