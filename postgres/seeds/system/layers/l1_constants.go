package layers

// L1_SEED_VERSION declara la versión semántica del contenido de L1.
// Bumpear en CADA cambio de dato visible.
const L1_SEED_VERSION = "1.2.0"

// L1_LAYER_NAME es el nombre canónico de la capa, usado por
// --seed-up-to-layer y por logs.
const L1_LAYER_NAME = "L1-readonly"

// UUIDs hardcodeados de las entidades canónicas de L1 (prefijo b1xxx).
// Razón: tests E2E y JSON exportado a KMP los referencian por UUID.
// Política ADR-7: UUIDs propios de L1, no coordinados con legacy.
const (
	L1_ROLE_ANNOUNCEMENT_VIEWER_ID = "b1000000-0000-0000-0000-000000000001"
	L1_USER_VIEWER_ID              = "b1000000-0000-0000-0000-000000000002"
	L1_SCHOOL_DEMO_ID              = "b1000000-0000-0000-0000-000000000003"
	L1_USER_ROLE_VIEWER_ID         = "b1000000-0000-0000-0000-000000000004"
	// P4-1 (plan B): L1_ROLE_PERMISSION_VIEWER_ID eliminado (b1...0005).
	// La tabla iam.role_permissions ya no existe.
	L1_MEMBERSHIP_VIEWER_ID = "b1000000-0000-0000-0000-000000000006"
)

// L1_MEMBERSHIP_ROLE es el valor que se usa en
// academic.memberships.role para el viewer. Tiene que pertenecer al
// CHECK constraint memberships_role_check
// (teacher|student|guardian|coordinator|admin|assistant). La
// semántica real del rol vive en iam.user_roles (announcement_viewer).
const L1_MEMBERSHIP_ROLE = "assistant"

// Credenciales de bootstrapping del usuario viewer de L1.
//
// SEGURIDAD: Estas credenciales son SÓLO para arranque inicial del
// sistema y entornos de prueba. En cualquier deployment cloud,
// rotar antes del primer login productivo.
const (
	L1_VIEWER_EMAIL    = "viewer@edugo.demo"
	L1_VIEWER_PASSWORD = "ChangeMe!2026"
)

// Strings semánticos reutilizables.
const (
	L1_ROLE_ANNOUNCEMENT_VIEWER_NAME = "announcement_viewer"
	L1_SCHOOL_DEMO_CODE              = "L1-DEMO"
	L1_SCHOOL_DEMO_NAME              = "Escuela Demo L1"
)
