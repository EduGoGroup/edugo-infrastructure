package scenarios

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/fixtures"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
)

// L2Form valida end-to-end la capa L2 del seed system (Fase 4 del
// rebuild). El scenario no inserta filas en el sistema: confía en que
// system.ApplySystem aplicó L0 + L1 + L2 (testdb.StartPostgres lo
// hace) y se limita, vía la fixture l2_constants_export, a:
//
//   - Verificar presencia y forma de la ScreenInstance
//     "announcement-form" (template_id, slot_data con 4 fields /
//     3 actions / api_prefix="platform").
//   - Verificar el mapping resource_screens
//     (announcements, screen_type=form, is_default=false).
//   - Verificar la derivación de permisos por evento:
//     SAVE_NEW → announcements:create, SAVE_EXISTING → announcements:update.
//   - No-regresión sobre la cadena L1 viewer→permisos:
//     viewer@edugo.demo sigue teniendo SOLO announcements:read.
//   - Exportar las constantes L2 al ApplyContext → fixtures-constants.json.
//
// Cobertura por F4-REQ:
//
//	F4-REQ-1.1  cubierto por SQL en l2_constants_export.verifyScreenInstance
//	F4-REQ-1.2  cubierto por parse JSON en l2_constants_export.verifyScreenInstance
//	F4-REQ-2.1  cubierto por SQL en l2_constants_export.verifyResourceScreen
//	F4-REQ-3.1  cubierto por parse JSON en l2_constants_export.verifyScreenInstance
//	F4-REQ-3.2  cubierto por parse JSON en l2_constants_export.verifyScreenInstance
//	(no-regresión L1) cubierto por l2_constants_export.verifyViewerPermissionsNoRegression
//
// Diferidos (Opción A — HTTP/UI deferred per spec, requires API server):
//
//	F4-REQ-1.3  HTTP resolve de pantalla — diferido hasta cierre del plan.
//	F4-REQ-3.3  KMP oculta botón Guardar para viewer — UI, fuera de scope SQL.
//	F4-REQ-5.2  POST 201 / GET menu con super_admin — requiere API server.
//	F4-REQ-5.3  POST 403 con viewer — requiere API server.
//
// El scenario integration test (Wave 2) marcará los puntos diferidos
// con t.Skip("HTTP/UI deferred per Opción A — requires API server").
//
// Refs: phase-4-layer-l2/{requirements,design}.md.
type L2Form struct{}

// Manifest implementa framework.Scenario.
func (s *L2Form) Manifest() framework.ScenarioManifest {
	return framework.ScenarioManifest{
		Name:         "l2_form",
		Description:  "Valida la capa L2 (ScreenInstance announcement-form + ResourceScreen form) del seed system y exporta sus identificadores.",
		FixtureNames: []string{"l2_constants_export"},
		Tags:         []string{"l2", "system", "screen-config"},
	}
}

// BuildFixtures implementa framework.Scenario. Devuelve siempre
// instancias frescas para no compartir estado mutable entre
// aplicaciones.
func (s *L2Form) BuildFixtures(ctx *framework.ApplyContext) []framework.Fixture {
	return []framework.Fixture{
		&fixtures.L2ConstantsExport{},
	}
}
