package scenarios

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/fixtures"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
)

// L1ReadOnly valida la capa L1 del seed system (Fase 3 del rebuild).
// El scenario no inserta filas en el sistema: confía en que
// system.ApplySystem aplicó el system seed (testdb.StartPostgres lo
// hace) y se limita a verificar la presencia del rol de contrato y
// exportar las constantes a ApplyContext.Constants →
// fixtures-constants.json.
//
// MP-09 F4: L1 quedó como CONTRATO PURO. El scenario se REDUCE al rol
// de contrato announcement_viewer; ya no valida DATO DE TENANT (escuela
// demo, usuario viewer, user_role, membership), que L1 dejó de sembrar.
// El nombre `l1_readonly` se conserva: el harness KMP valida escenarios
// por nombre en fixtures-constants.json.
//
// Refs: phase-3-layer-l1/{requirements,design}.md (F3-REQ-5.*), ADR-7.
type L1ReadOnly struct{}

// Manifest implementa framework.Scenario.
func (s *L1ReadOnly) Manifest() framework.ScenarioManifest {
	return framework.ScenarioManifest{
		Name:         "l1_readonly",
		Description:  "Valida la capa L1 (rol de contrato announcement_viewer) del seed system y exporta su identificador.",
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
