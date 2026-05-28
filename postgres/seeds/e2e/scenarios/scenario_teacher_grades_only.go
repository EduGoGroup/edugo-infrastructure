package scenarios

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/fixtures"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
)

// TeacherGradesOnly reproduce el caso "docente con permisos parciales
// sobre la pantalla de calificaciones". Combina:
//
//   - fixture_role_only(teacher): docente con escuela/usuario/membership.
//   - fixture_partial_crud([grades]): rol overlay con (create, read)
//     sobre el recurso `grades` — modela "puede listar y registrar pero
//     no editar/borrar" (C-REQ-9.4).
//   - fixture_screen_only(grades-list): contenido mínimo para que la
//     pantalla cargue. La fixture crea 1 subject + 1 academic_period +
//     1 alumno (user + membership) + 1 grade que apunta al teacher
//     provisto por role_only para que la pantalla muestre datos reales.
//
// Tags: rbac, screen-config.
type TeacherGradesOnly struct{}

// Manifest implementa framework.Scenario.
func (s *TeacherGradesOnly) Manifest() framework.ScenarioManifest {
	return framework.ScenarioManifest{
		Name:        "teacher_grades_only",
		Description: "Docente con (create, read) parciales sobre grades + pantalla grades-list.",
		FixtureNames: []string{
			"role_only",
			"partial_crud",
			"screen_only",
		},
		Tags: []string{"rbac", "screen-config"},
	}
}

// BuildFixtures implementa framework.Scenario.
func (s *TeacherGradesOnly) BuildFixtures(ctx *framework.ApplyContext) []framework.Fixture {
	return []framework.Fixture{
		&fixtures.RoleOnly{RoleCode: "teacher"},
		&fixtures.PartialCrud{Resources: []string{"grades"}},
		&fixtures.ScreenOnly{ScreenKey: "grades-list"},
	}
}
