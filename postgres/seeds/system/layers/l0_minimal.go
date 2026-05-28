// Package layers contiene las capas L0..L4 del seed system de EduGo,
// reorganizadas tras el rebuild documentado en
// EduUI/edugo-ui-kmp/e2e-integration-plan/seed-rebuild-spec/.
//
// Capa L0 — Mínimo viable
// =======================
//
// L0 es el piso del sistema: el dataset más pequeño con el que EduGo
// arranca y responde correctamente al flujo dinámico end-to-end. Su
// propósito no es ser útil para un cliente, sino ser diagnóstico —
// si algo del código dinámico (menú filtrado, screen-config/resolve,
// derivación de permisos en ScreenContract) está roto, L0 lo detecta
// antes de cargar las ~1200 filas del sistema completo.
//
// Composición (13 filas en 8 tablas):
//   - 1 recurso: announcements (scope=school, raíz del menú)
//   - 1 rol: super_admin (scope=system)
//   - 4 permisos: academic.announcements.{read,create,update,delete}
//   - 3 screen_templates: list-basic-v1, detail-basic-v1, form-basic-v1
//   - 1 screen_instance: announcements-list
//   - 1 resource_screen: announcements ↔ announcements-list
//   - 1 user: super_admin@edugo.system (password hash bootstrap)
//   - 1 user_role: user × super_admin (scope=system)
//
// P4-1 (plan B): los 4 role_permissions del modelo legacy fueron
// eliminados. Los permisos efectivos del super_admin se obtienen vía
// iam.role_grants con el pattern wildcard `*`, sembrado por L4.
//
// Por qué "announcements" como recurso L0:
//  1. Caso de uso real, no juguete.
//  2. Historial de bug F2·H3.a (apiPrefix incorrecto entre academic:
//     y platform:) — sirve como canario natural.
//  3. Scope school, representativo del producto.
//  4. ScreenContract estándar, sin overrides exóticos.
//
// Refs: ADR-1, ADR-2, ADR-6 de seed-rebuild-spec/00-context/decisions.md;
//
//	lecciones F2·H3.a en lecciones-aprendidas.md.
package layers

import "gorm.io/gorm"

// l0Layer implementa system.Layer por duck-typing (no se importa la
// interfaz para evitar ciclo seeds/system ↔ seeds/system/layers).
type l0Layer struct{}

// NewL0 construye una instancia de la capa L0.
// Se registra en system.Layers() como única capa post-Fase-2 (ADR-6).
func NewL0() *l0Layer { return &l0Layer{} }

// Name retorna el identificador canónico de la capa, usado por la
// flag --seed-up-to-layer y por logs.
func (l *l0Layer) Name() string { return L0_LAYER_NAME }

// SeedVersion retorna la versión semántica del contenido de L0.
// Bumpear L0_SEED_VERSION en cada cambio de dato visible.
func (l *l0Layer) SeedVersion() string { return L0_SEED_VERSION }

// Apply siembra L0 en orden respetando las FK del esquema.
// Orden obligatorio:
//  1. Resources (sin dependencias)
//  2. ScreenTemplates → ScreenInstances → ResourceScreens (FK a Resources)
//  3. Permissions (FK a Resources) → Roles → RolePermissions
//  4. Users → UserRoles (FK a Users y Roles)
//
// Los applyL0_* viven en archivos separados (l0_resources.go,
// l0_roles.go, l0_screens.go, l0_users.go).
func (l *l0Layer) Apply(tx *gorm.DB) error {
	if err := applyL0Resources(tx); err != nil {
		return err
	}
	if err := applyL0Screens(tx); err != nil {
		return err
	}
	if err := applyL0Roles(tx); err != nil {
		return err
	}
	if err := applyL0Users(tx); err != nil {
		return err
	}
	return nil
}
