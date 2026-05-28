package scenarios

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/fixtures"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
)

// L1ReadOnly valida end-to-end la capa L1 del seed system (Fase 3
// del rebuild). El scenario no inserta filas en el sistema: confía
// en que system.ApplySystem aplicó L0 + L1 (testdb.StartPostgres lo
// hace) y se limita a verificar su presencia y exportar las
// constantes a ApplyContext.Constants → fixtures-constants.json.
//
// Refs: phase-3-layer-l1/{requirements,design}.md (F3-REQ-5.*,
// F3-REQ-6.2), ADR-7 (escuela mínima L1 para scope=school).
type L1ReadOnly struct{}

// Manifest implementa framework.Scenario.
func (s *L1ReadOnly) Manifest() framework.ScenarioManifest {
	return framework.ScenarioManifest{
		Name:         "l1_readonly",
		Description:  "Valida la capa L1 (rol announcement_viewer + escuela demo + usuario viewer) del seed system y exporta sus identificadores.",
		FixtureNames: []string{"l1_constants_export"},
		Tags:         []string{"l1", "system", "rbac"},
	}
}

// BuildFixtures implementa framework.Scenario. Devuelve siempre
// instancias frescas para no compartir estado mutable entre
// aplicaciones.
func (s *L1ReadOnly) BuildFixtures(ctx *framework.ApplyContext) []framework.Fixture {
	return []framework.Fixture{
		&fixtures.L1ConstantsExport{},
	}
}
