// Capa L1 — Rol read-only (primer caso de gating real)
// =====================================================
//
// L1 añade encima de L0 un rol `announcement_viewer` (scope school)
// con un único `RolePermission` apuntando a `announcements:read`,
// más el usuario y la escuela mínima necesarios para que ese rol
// sea operativo (ADR-7). Propósito: validar que el sistema oculta
// correctamente las acciones para las que el usuario no tiene
// permiso (F3-REQ-3.x).
//
// Composición (5 filas):
//   - 1 academic.schools: Escuela Demo L1 (b1...0003)
//   - 1 iam.roles: announcement_viewer (b1...0001, scope=school)
//   - 1 auth.users: viewer@edugo.demo (b1...0002)
//   - 1 iam.user_roles: viewer × role × school (b1...0004)
//   - 1 academic.memberships: viewer × escuela L1 (b1...0006) —
//     requerido por identity API para que switch-context emita JWT
//     con permisos efectivos (CHECK constraint exige role del enum;
//     se usa "assistant" como valor mínimo).
//
// P4-1 (plan B): el role_permission viewer × announcements:read fue
// eliminado; el permiso se otorga vía iam.role_grants desde L4
// (pattern `academic.announcements.read`).
//
// Acumulado tras L0+L1: 18 filas.
//
// Refs: ADR-7 de seed-rebuild-spec/00-context/decisions.md;
//
//	phase-3-layer-l1/{design,requirements}.md.
package layers

import "gorm.io/gorm"

// l1Layer implementa system.Layer por duck-typing (no se importa la
// interfaz para evitar ciclo seeds/system ↔ seeds/system/layers).
type l1Layer struct{}

// NewL1 construye una instancia de la capa L1.
// Se registra en system.Layers() tras NewL0() (F3-REQ-4.1).
func NewL1() *l1Layer { return &l1Layer{} }

// Name retorna el identificador canónico de la capa, usado por la
// flag --seed-up-to-layer y por logs.
func (l *l1Layer) Name() string { return L1_LAYER_NAME }

// SeedVersion retorna la versión semántica del contenido de L1.
// Bumpear L1_SEED_VERSION en cada cambio de dato visible.
func (l *l1Layer) SeedVersion() string { return L1_SEED_VERSION }

// Apply siembra L1 en orden respetando las FK del esquema.
// Orden obligatorio:
//  1. academic.schools     (sin dependencias hacia L1)
//  2. iam.roles            (sin dependencias hacia L1)
//  3. auth.users           (sin dependencias hacia L1)
//  4. iam.user_roles       (FK a users de L1, roles de L1, schools de L1)
//  5. academic.memberships (FK a users de L1, schools de L1)
func (l *l1Layer) Apply(tx *gorm.DB) error {
	if err := applyL1School(tx); err != nil {
		return err
	}
	if err := applyL1Role(tx); err != nil {
		return err
	}
	if err := applyL1User(tx); err != nil {
		return err
	}
	if err := applyL1UserRole(tx); err != nil {
		return err
	}
	return applyL1Membership(tx)
}
