package scenarios

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/fixtures"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
)

// SuperAdminGlobalFlow es el scenario que materializa el caso "global
// super_admin sin membership" — un usuario con rol L0 super_admin
// (school_id=NULL) que debe poder navegar el flujo SchoolSelector →
// switchContext → UnitSelector → Dashboard sin tener ninguna fila en
// `academic.memberships`.
//
// El scenario se usa principalmente desde el test cross-API en
// `edugo-dev-environment/test/integration/superadmin_flow/`. Reproduce
// los 4 bugs detectados en sesión 2026-05-12:
//
//  1. L4 super_admin sin context:browse_schools/context:browse_units.
//  2. academic.routes_school/routes_unit faltantes de RequireAnyPermission.
//  3. SchoolModels.kt usaba `data` en vez de `schools`.
//  4. switch_context.go bloqueaba al super_admin pidiendo membership.
//
// Composición:
//
//   - `l0_constants_export`: exporta los UUIDs del rol super_admin de L0
//     al JSON para que el test los referencie sin hardcodear.
//   - `role_only(school_admin)`: crea 1 escuela (necesaria para que el
//     SchoolSelector tenga al menos 1 fila listable) y un user "normal"
//     con membership (NO se usa en el test, pero garantiza que la BD
//     tiene un escenario de tenant baseline).
//   - `global_user_no_membership`: crea el user "global-super@..." con
//     rol super_admin global (school_id=NULL) y SIN membership.
//
// La unidad académica (academic_unit) la provee el demo seed que el
// harness aplica antes (demo.ApplyDemo siembra ~14 unidades).
//
// Refs: sesión-2026-05-12, ADR-6, F6-REQ-… (capa L4 context:browse_*).
type SuperAdminGlobalFlow struct{}

// Manifest implementa framework.Scenario.
func (s *SuperAdminGlobalFlow) Manifest() framework.ScenarioManifest {
	return framework.ScenarioManifest{
		Name:        "super_admin_global_flow",
		Description: "Global super_admin (school_id=NULL, sin membership) + escuela baseline para validar SchoolSelector → switchContext → UnitSelector → Dashboard.",
		FixtureNames: []string{
			"l0_constants_export",
			"role_only",
			"global_user_no_membership",
		},
		Tags: []string{"super_admin", "cross_api", "school_selector"},
	}
}

// BuildFixtures implementa framework.Scenario. Devuelve siempre
// instancias frescas para no compartir estado mutable entre
// aplicaciones.
func (s *SuperAdminGlobalFlow) BuildFixtures(ctx *framework.ApplyContext) []framework.Fixture {
	return []framework.Fixture{
		&fixtures.L0ConstantsExport{},
		&fixtures.RoleOnly{RoleCode: "school_admin"},
		&fixtures.GlobalUserNoMembership{},
	}
}
