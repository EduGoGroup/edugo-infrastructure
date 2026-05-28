package scenarios

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/fixtures"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
)

// ObserverAudits reproduce el caso "auditor read-only mira el sub-árbol
// de auditoría". Combina:
//
//   - fixture_role_only(readonly_auditor): aporta school/user/user_role
//     /membership ligados al rol del catálogo "readonly_auditor".
//   - fixture_readonly_role([audit]): crea el rol overlay con
//     scope=unit que aporta la capacidad "readonly_role" requerida por
//     fixture_menu_subtree (decisión de diseño: role_only NO provee
//     "readonly_role"; necesitamos un rol overlay del scenario para
//     ligarle los permisos del subtree sin tocar el catálogo). El
//     resource key es "audit" (catálogo iam.resources del production
//     seed), no "audit-events" (ese es el screenKey de la pantalla).
//   - fixture_menu_subtree(audit): asigna <resource>:read a todos
//     los recursos del subtree para iluminar la rama del menú.
//
// Tags: rbac, menu, audit.
type ObserverAudits struct{}

// Manifest implementa framework.Scenario.
func (s *ObserverAudits) Manifest() framework.ScenarioManifest {
	return framework.ScenarioManifest{
		Name:        "observer_audits",
		Description: "Auditor read-only con visibilidad sobre el subtree de auditoría (rol overlay + menú).",
		FixtureNames: []string{
			"role_only",
			"readonly_role",
			"menu_subtree",
		},
		Tags: []string{"rbac", "menu", "audit"},
	}
}

// BuildFixtures implementa framework.Scenario. Devuelve siempre
// instancias frescas para no compartir estado mutable entre
// aplicaciones.
func (s *ObserverAudits) BuildFixtures(ctx *framework.ApplyContext) []framework.Fixture {
	return []framework.Fixture{
		&fixtures.RoleOnly{RoleCode: "readonly_auditor"},
		&fixtures.ReadonlyRole{Resources: []string{"audit"}},
		&fixtures.MenuSubtree{SubtreeRoot: "audit"},
	}
}
