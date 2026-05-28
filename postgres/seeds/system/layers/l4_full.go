// Capa L4 — Sistema completo reorganizado por dominio
// =====================================================
//
// L4 cierra el plan de rebuild. A diferencia de L0..L3 (que llevan
// todos sus datos en seeds/system/layers/<lN>_*.go), L4 separa la
// implementación de la capa (este archivo) de los datos por dominio
// (paquete seeds/system/l4/<dominio>.go).
//
// Composición (ver phase-6-layer-l4/design.md §2):
//   - resources.go         — recursos del menú restantes
//   - roles_permissions.go — student, teacher, guardian, admin,
//     school_admin (super_admin en L0,
//     announcement_viewer en L1)
//   - screen_templates.go  — templates adicionales + refactor de
//     `definition` de los 3 templates base
//     de L0 (fix zones para SDUI engine)
//   - screen_instances.go  — ~73 instances + 14 phantom legítimas
//   - resource_screens.go  — mappings recurso↔pantalla
//   - concept_types.go     — concept_types + concept_definitions
//
// ADR-6: Layer_Legacy está desactivado de runtime desde Fase 2. El
// directorio [archivado pre-Fase-6]  permanece en disco solo como
// inventario/guía hasta el bloque C de Fase 6 (borrado físico).
//
// Refs: phase-6-layer-l4/{design,requirements,tasks}.md.
package layers

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l4"
	"gorm.io/gorm"
)

// l4Layer implementa system.Layer por duck-typing (no se importa la
// interfaz para evitar ciclo seeds/system ↔ seeds/system/layers).
type l4Layer struct{}

// NewL4 construye una instancia de la capa L4.
// Se registra en system.Layers() tras NewL3() (F6-REQ-1.2).
func NewL4() *l4Layer { return &l4Layer{} }

// Name retorna el identificador canónico de la capa, usado por la
// flag --seed-up-to-layer y por logs.
func (l *l4Layer) Name() string { return L4_LAYER_NAME }

// SeedVersion retorna la versión semántica del contenido de L4.
// Bumpear L4_SEED_VERSION en cada cambio de dato visible.
func (l *l4Layer) SeedVersion() string { return L4_SEED_VERSION }

// Apply siembra L4 delegando a los sub-paquetes por dominio.
// Orden obligatorio (respeta FK del esquema):
//  1. resources           (sin deps L4)
//  2. roles_permissions   (FK a resources)
//  3. screen_templates    (sin deps L4; refactoriza también las
//     definitions de L0)
//  4. screen_instances    (FK a screen_templates)
//  5. resource_screens    (FK a resources + screen_instances)
//  6. concept_types       (sin deps L4)
func (l *l4Layer) Apply(tx *gorm.DB) error {
	if err := l4.ApplyResources(tx); err != nil {
		return err
	}
	if err := l4.ApplyRolesPermissions(tx); err != nil {
		return err
	}
	if err := l4.ApplyScreenTemplates(tx); err != nil {
		return err
	}
	if err := l4.ApplyScreenInstances(tx); err != nil {
		return err
	}
	if err := l4.ApplyResourceScreens(tx); err != nil {
		return err
	}
	if err := l4.ApplyConceptTypes(tx); err != nil {
		return err
	}
	return nil
}
