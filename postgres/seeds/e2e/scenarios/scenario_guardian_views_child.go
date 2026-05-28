package scenarios

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/fixtures"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
)

// GuardianViewsChild reproduce el caso "apoderado consulta el progreso
// del hijo". Combina:
//
//   - fixture_role_only(guardian): apoderado con escuela/usuario.
//   - fixture_guardian_relation: vínculo guardian↔student real en
//     academic.guardian_relations (más el student.user y su membership)
//     para que la pantalla `child-progress` resuelva contra datos reales.
//   - fixture_screen_only(child-progress): contenido mínimo
//     para que la pantalla cargue (cae al branch default de screen_only
//     porque la pantalla `child-progress` no requiere inserts
//     específicos en assessment/etc.: el chrome lo provee el production
//     seed y el vínculo lo provee guardian_relation).
//
// El screenKey es "child-progress" (catálogo ui_config del production
// seed). NO es "guardian-child-progress".
//
// Tags: rbac, screen-config.
type GuardianViewsChild struct{}

// Manifest implementa framework.Scenario.
func (s *GuardianViewsChild) Manifest() framework.ScenarioManifest {
	return framework.ScenarioManifest{
		Name:        "guardian_views_child",
		Description: "Apoderado con vínculo a un student y pantalla child-progress.",
		FixtureNames: []string{
			"role_only",
			"guardian_relation",
			"screen_only",
		},
		Tags: []string{"rbac", "screen-config"},
	}
}

// BuildFixtures implementa framework.Scenario.
func (s *GuardianViewsChild) BuildFixtures(ctx *framework.ApplyContext) []framework.Fixture {
	return []framework.Fixture{
		&fixtures.RoleOnly{RoleCode: "guardian"},
		&fixtures.GuardianRelation{},
		&fixtures.ScreenOnly{ScreenKey: "child-progress"},
	}
}
