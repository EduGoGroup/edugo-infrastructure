package scenarios

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/fixtures"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
)

// L3Isolation valida el aislamiento de menús/permisos entre recursos
// y el soporte de CRUD parcial (materials sin :delete).
//
// Validaciones SQL (en L3IsolationConstants.Apply):
//   - F5-REQ-1.1: resource materials existe con scope=unit.
//   - F5-REQ-2.1: 3 permisos materials:{read,create,update}; ausencia explícita de :delete.
//   - F5-REQ-2.2: 3 role_permissions super_admin × materials.
//   - F5-REQ-3.x: 2 ScreenInstances + 2 resource_screens (list default, form no-default).
//   - No-regresión viewer L1: permisos siguen siendo {announcements:read}.
//
// Diferidos por Opción A (requieren API server / KMP runtime):
//   - F5-REQ-2.3: viewer GET /menu no devuelve materials.
//   - F5-REQ-4.1/4.2: GET /menu por rol.
//   - F5-REQ-4.3: viewer DELETE /materials retorna 404.
//   - F5-REQ-6.2: super_admin POST /materials retorna 201, DELETE retorna 405/403.
//   - F5-REQ-6.3: viewer GET /materials retorna 403.
//
// El scenario integration test (Wave 2) marcará los puntos diferidos
// con t.Skip("HTTP/UI deferred per Opción A — requires API server").
//
// Refs: phase-5-layer-l3/{requirements,design}.md.
type L3Isolation struct{}

// Manifest implementa framework.Scenario.
func (s *L3Isolation) Manifest() framework.ScenarioManifest {
	return framework.ScenarioManifest{
		Name:         "l3_isolation",
		Description:  "Valida la capa L3 (resource materials con CRUD parcial sin :delete + 2 ScreenInstances + 2 ResourceScreens) del seed system y exporta sus identificadores.",
		FixtureNames: []string{"l3_constants_export"},
		Tags:         []string{"l3", "system", "rbac", "menu", "screen-config"},
	}
}

// BuildFixtures implementa framework.Scenario. Devuelve siempre
// instancias frescas para no compartir estado mutable entre
// aplicaciones.
func (s *L3Isolation) BuildFixtures(ctx *framework.ApplyContext) []framework.Fixture {
	return []framework.Fixture{
		&fixtures.L3IsolationConstants{},
	}
}
