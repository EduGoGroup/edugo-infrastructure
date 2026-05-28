// Package framework define el motor de fixtures compositivas que
// reemplaza la pila aditiva fase0..fase4 del seed E2E vigente.
//
// Conceptos clave:
//
//   - Fixture: pieza pequeña, idempotente, declara qué necesita
//     (Requires) y qué provee (Provides). Vive en seeds/e2e/fixtures/.
//   - Scenario: receta que combina fixtures, asigna un namespace único
//     (TenantPrefix + SchemaPrefix) y declara constantes públicas para
//     los tests. Vive en seeds/e2e/scenarios/.
//   - Composer: resuelve dependencias por orden topológico y aplica
//     todo el scenario dentro de una única transacción.
//   - Cleaner: borra exclusivamente las filas con el prefijo del
//     scenario; nunca toca el production seed (resources, permissions,
//     roles, screen_templates, screen_instances del catálogo).
//
// El framework hereda dos lecciones críticas:
//
//   - F2·H5: nunca depender del tag gorm:"default:..." para booleanos
//     críticos (IsActive, IsMenuVisible, IsDefault, IsPinned, IsPublic,
//     IsTimed, membership.IsActive). Usar UpsertBool o asignación
//     explícita.
//   - DA2 del plan E2E: namespace UUID e2e00000-... y prefijo de código
//     E2E- siguen reservados; los nuevos scenarios usan un sub-namespace
//     derivado por SHA-1 del nombre.
//
// La spec completa vive en
// EduUI/edugo-ui-kmp/e2e-integration-plan/system-data-quality-spec/phase-c-fixtures-refactor.
package framework
