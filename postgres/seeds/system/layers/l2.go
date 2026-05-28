// Capa L2 — Segunda pantalla por recurso (form + gating de acciones)
// ====================================================================
//
// L2 añade encima de L0+L1 una segunda pantalla del recurso
// `announcements`: el formulario `announcement-form` (template
// `form-basic-v1`) con sus 3 fields canónicos y 3 acciones cada una
// con su `permission` y `event` derivable por `ScreenContract`.
//
// Composición (2 filas):
//   - 1 ui_config.screen_instances: announcement-form
//   - 1 ui_config.resource_screens: announcements ↔ announcement-form
//     (screen_type=form, is_default=false)
//
// Acumulado tras L0+L1+L2: 24 filas.
//
// Refs: phase-4-layer-l2/{design,requirements}.md.
package layers

import "gorm.io/gorm"

// l2Layer implementa system.Layer por duck-typing (no se importa la
// interfaz para evitar ciclo seeds/system ↔ seeds/system/layers).
type l2Layer struct{}

// NewL2 construye una instancia de la capa L2.
// Se registra en system.Layers() tras NewL1().
func NewL2() *l2Layer { return &l2Layer{} }

// Name retorna el identificador canónico de la capa, usado por la
// flag --seed-up-to-layer y por logs.
func (l *l2Layer) Name() string { return L2_LAYER_NAME }

// SeedVersion retorna la versión semántica del contenido de L2.
// Bumpear L2_SEED_VERSION en cada cambio de dato visible.
func (l *l2Layer) SeedVersion() string { return L2_SEED_VERSION }

// Apply siembra L2 en orden respetando las FK del esquema.
// Orden obligatorio:
//  1. screen_instances  (FK a screen_templates de L0)
//  2. resource_screens  (FK a resources de L0 y screen_instances de L2)
func (l *l2Layer) Apply(tx *gorm.DB) error {
	if err := applyL2Screens(tx); err != nil {
		return err
	}
	if err := applyL2ResourceScreens(tx); err != nil {
		return err
	}
	return nil
}
