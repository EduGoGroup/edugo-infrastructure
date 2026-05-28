// Package l4 contiene los sub-archivos por dominio de la capa L4
// del seed system (sistema completo reorganizado por dominio).
//
// L4 cierra el plan de rebuild definido en
// `EduUI/edugo-ui-kmp/e2e-integration-plan/seed-rebuild-spec/`. A
// diferencia de L0..L3 (que viven completamente en
// seeds/system/layers/<lN>_*.go), L4 separa la implementación de la
// capa (seeds/system/layers/l4_full.go) de los datos por dominio
// (este paquete).
//
// Estructura:
//   - resources.go        — ApplyResources: recursos del menú
//     (excluye announcements en L0 y materials
//     en L3).
//   - roles_permissions.go — ApplyRolesPermissions: roles + permisos
//   - role_permissions de los 5 roles que
//     faltan (student, teacher, guardian,
//     admin, school_admin). super_admin está
//     en L0; announcement_viewer en L1.
//   - screen_templates.go — ApplyScreenTemplates: templates
//     adicionales (dashboard-basic-v1 y
//     especiales). También refactoriza la
//     `definition` JSON de los 3 templates
//     base de L0 que se sembraron con `{}`
//     (zones requeridas por el SDUI engine
//     del KMP — bug detectado pre-Fase-6).
//   - screen_instances.go — ApplyScreenInstances: ~73 instances
//   - 14 screen_key_phantom legítimas
//     reportadas por el cross-checker.
//   - resource_screens.go — ApplyResourceScreens: mappings
//     recurso↔pantalla recalculados desde los
//     recursos y pantallas sobrevivientes.
//   - concept_types.go    — ApplyConceptTypes: tipos + definiciones
//     del catálogo conceptual.
//
// Cada Apply* se invoca desde `layers.l4Layer.Apply` en
// seeds/system/layers/l4_full.go.
//
// Política anti-copy-paste (F6-REQ-2.5): ningún sub-archivo puede
// ser una copia literal de un slice de [archivado pre-Fase-6] data.go.
// El legacy es inventario/guía; cada fila L4 debe haber sido revisada
// y documentada (renombre, descarte, refactor) en
// phase-6-layer-l4/decisions-log.md.
package l4
