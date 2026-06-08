// Capa L3 — Segundo recurso (materials) con CRUD parcial
// ========================================================
//
// L3 añade encima de L0+L1+L2 el segundo recurso del sistema:
// `materials` (scope=unit), con CRUD parcial (read/create/update,
// SIN delete; F5-REQ-2.1). Es el hito conceptual del plan de rebuild:
// valida que el sistema dinámico de menú/pantalla aísla
// correctamente permisos y menús entre recursos (F5-REQ-4.x).
//
// Composición (6 filas):
//   - 1 iam.resources: materials (scope=unit, api_prefix=academic)
//   - 3 iam.permissions: materials:{read,create,update} (NO delete)
//   - 1 ui_config.screen_instances: materials-list MÍNIMA (no se renderiza;
//     la pantalla es NATIVA Compose). Existe SOLO para satisfacer la FK
//     fk_resource_screens_screen_key del mapping de menú. Ver l3_screens.go.
//   - 1 ui_config.resource_screens: list (default).
//
// Poda SDUI material (2026-06-07) + corrección F2 (2026-06-08): la poda
// quitó AMBAS ScreenInstances, pero `materials-list` tiene mapping de menú
// y la FK fk_resource_screens_screen_key exige su screen_instance → un
// recreate limpio fallaba (23503). F2 RESTAURA la screen_instance mínima
// `materials-list` (patrón batch-enroll/join-requests-inbox de L4).
// `material-form` SIGUE PODADO (sin mapping → sin FK → no se resiembra).
// Detalle en l3_screens.go.
//
// P4-1 (plan B): los 3 role_permissions super_admin × materials fueron
// eliminados; el super_admin cubre estos permisos vía pattern `*` en
// iam.role_grants desde L4.
//
// Refs: phase-5-layer-l3/{design,requirements}.md.
package layers

import "gorm.io/gorm"

// l3Layer implementa system.Layer por duck-typing (no se importa la
// interfaz para evitar ciclo seeds/system ↔ seeds/system/layers).
type l3Layer struct{}

// NewL3 construye una instancia de la capa L3.
// Se registra en system.Layers() tras NewL2() (F5-REQ-5.1).
func NewL3() *l3Layer { return &l3Layer{} }

// Name retorna el identificador canónico de la capa, usado por la
// flag --seed-up-to-layer y por logs.
func (l *l3Layer) Name() string { return L3_LAYER_NAME }

// SeedVersion retorna la versión semántica del contenido de L3.
// Bumpear L3_SEED_VERSION en cada cambio de dato visible.
func (l *l3Layer) SeedVersion() string { return L3_SEED_VERSION }

// Apply siembra L3 en orden respetando las FK del esquema.
// Orden obligatorio:
//  1. resources         (sin deps L3)
//  2. permissions       (FK a resources de L3)
//  3. screens           (screen_instance materials-list MÍNIMA; debe existir
//     ANTES del resource_screen por la FK fk_resource_screens_screen_key)
//  4. resource_screens  (FK a resources de L3 + FK screen_key →
//     screen_instances.screen_key, materializada por la fila del paso 3)
func (l *l3Layer) Apply(tx *gorm.DB) error {
	if err := applyL3Resources(tx); err != nil {
		return err
	}
	if err := applyL3Permissions(tx); err != nil {
		return err
	}
	if err := applyL3Screens(tx); err != nil {
		return err
	}
	if err := applyL3ResourceScreens(tx); err != nil {
		return err
	}
	return nil
}
