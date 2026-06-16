// Capa L1 — Rol read-only (primer caso de gating real)
// =====================================================
//
// L1 añade encima de L0 un único rol de CONTRATO:
// `announcement_viewer` (scope=school). Propósito: validar que el
// sistema oculta correctamente las acciones para las que el usuario no
// tiene permiso (F3-REQ-3.x). El permiso efectivo del rol
// (`academic.announcements.read`) se otorga vía iam.role_grants desde
// L4.
//
// Composición (1 fila):
//   - 1 iam.roles: announcement_viewer (b1...0001, scope=school)
//
// MP-09 F4: L1 dejó de sembrar DATO DE TENANT (escuela demo, usuario
// viewer, user_role y membership). Ese dato vivo equivalente vive en
// playground_v2/base. system/ queda como CONTRATO PURO: sólo el rol.
//
// Acumulado tras L0+L1: 18 filas (17 de L0 + 1 de L1).
//
// Refs: ADR-7 de seed-rebuild-spec/00-context/decisions.md;
//
//	phase-3-layer-l1/{design,requirements}.md.
package layers

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l4"
	"gorm.io/gorm"
)

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
//  0. academic.invitation_types (catálogo MP-08; sin deps de FK). Se adelanta
//     desde L4 porque otras capas/seeds del ecosistema resuelven
//     invitation_type_id por FK; ApplyInvitationTypes es idempotente, así que
//     reaplicarlo en L4 no duplica.
//  1. iam.roles (rol de contrato announcement_viewer, sin deps hacia L1)
func (l *l1Layer) Apply(tx *gorm.DB) error {
	if err := l4.ApplyInvitationTypes(tx); err != nil {
		return err
	}
	return applyL1Role(tx)
}
