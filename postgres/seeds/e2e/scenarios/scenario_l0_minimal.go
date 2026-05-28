package scenarios

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/fixtures"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
)

// L0Minimal valida end-to-end la capa L0 del seed system (Fase 2
// del rebuild). El scenario no inserta filas en el sistema: confía
// en que system.ApplySystem aplicó L0 (testdb.StartPostgres lo hace)
// y se limita a verificar su presencia y exportar las constantes a
// ApplyContext.Constants → fixtures-constants.json.
//
// Refs: phase-2-layer-l0/{requirements,design}.md (F2-REQ-6),
// ADR-6 (no coexistencia con legacy).
type L0Minimal struct{}

// Manifest implementa framework.Scenario.
func (s *L0Minimal) Manifest() framework.ScenarioManifest {
	return framework.ScenarioManifest{
		Name:         "l0_minimal",
		Description:  "Valida la capa L0 (mínimo viable) del seed system y exporta sus identificadores.",
		FixtureNames: []string{"l0_constants_export"},
		Tags:         []string{"l0", "system"},
	}
}

// BuildFixtures implementa framework.Scenario. Devuelve siempre
// instancias frescas para no compartir estado mutable entre
// aplicaciones.
func (s *L0Minimal) BuildFixtures(ctx *framework.ApplyContext) []framework.Fixture {
	return []framework.Fixture{
		&fixtures.L0ConstantsExport{},
	}
}
