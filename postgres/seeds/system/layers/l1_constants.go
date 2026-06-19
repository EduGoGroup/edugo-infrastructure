package layers

// L1_SEED_VERSION declara la versión semántica del contenido de L1.
// Bumpear en CADA cambio de dato visible.
//
// MP-09 F4 (1.2.0 → 1.3.0): L1 dejó de sembrar DATO DE TENANT (escuela
// demo, usuario viewer, user_role y membership). Sólo queda el rol de
// contrato announcement_viewer.
//
// 1.3.0 → 1.4.0 (2026-06-17, ADR 0024 sub-deuda "herencia del landing"):
// el rol announcement_viewer gana landing_screen_key=dashboard-schooladmin
// (antes NULL → caía a school.default "dashboard-home", el dashboard básico
// genérico; con landing propio aterriza en el dashboard de su superficie).
// 1.4.0 → 1.4.1 (2026-06-19, bug 0069): el rol de contrato announcement_viewer
// se siembra con IsSystem=true (apply + accessor espejo).
const L1_SEED_VERSION = "1.4.1"

// L1_LAYER_NAME es el nombre canónico de la capa, usado por
// --seed-up-to-layer y por logs.
const L1_LAYER_NAME = "L1-readonly"

// UUID hardcodeado del rol canónico de L1 (prefijo b1xxx).
// Razón: tests E2E y JSON exportado a KMP lo referencian por UUID.
// Política ADR-7: UUIDs propios de L1, no coordinados con legacy.
const (
	L1_ROLE_ANNOUNCEMENT_VIEWER_ID = "b1000000-0000-0000-0000-000000000001"
)

// Strings semánticos reutilizables del rol de contrato.
const (
	L1_ROLE_ANNOUNCEMENT_VIEWER_NAME = "announcement_viewer"
)
