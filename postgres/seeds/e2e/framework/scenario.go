package framework

// Scenario es la receta de alto nivel que combina fixtures y les
// asigna un namespace propio. Cada scenario se registra una sola vez
// (RegisterScenario) y se aplica desde un test o desde el binario
// seed_e2e (composer.Apply).
//
// Reglas:
//
//   - BuildFixtures se llama una sola vez por aplicación del scenario;
//     debe construir instancias frescas (no compartir estado mutable
//     entre scenarios distintos).
//   - El scenario NO crea filas en tablas del production seed
//     (resources, permissions, roles del catálogo, screen_templates,
//     screen_instances, role_permissions del catálogo). Sólo combina
//     fixtures que respetan esa restricción.
//   - Los UUIDs y códigos visibles que el scenario produce derivan del
//     TenantPrefix/SchemaPrefix calculado por namespace.Derive a
//     partir del nombre del scenario (excepción: legacy_e2e fuerza
//     hash 00000000 para mantener paridad).
type Scenario interface {
	// Manifest devuelve el descriptor estático del scenario.
	Manifest() ScenarioManifest

	// BuildFixtures parametriza las fixtures que componen el scenario.
	// Recibe ApplyContext para que el scenario pueda inyectar valores
	// dinámicos (ej. el roleCode para fixture_role_only).
	BuildFixtures(ctx *ApplyContext) []Fixture
}

// ScenarioManifest describe estáticamente un scenario.
type ScenarioManifest struct {
	// Name identifica unívocamente al scenario dentro del registry.
	// Convención: snake_case sin prefijo "scenario_"
	// (ej. "teacher_grades_only", "legacy_e2e").
	Name string

	// Description es un docstring corto utilizado por la documentación
	// generada y el output de seed_e2e --list.
	Description string

	// Params permite parametrizar el scenario sin recompilar el binario.
	// Ej. roleCode = "teacher" vs "assistant_teacher". Las fixtures lo
	// leen vía ApplyContext.RawParams.
	Params map[string]string

	// FixtureNames es un orden declarativo. El composer aún resuelve
	// dependencias por Provides/Requires; este orden sirve como hint
	// cuando hay múltiples órdenes válidos y para fines documentales.
	FixtureNames []string

	// Tags es la lista de etiquetas usadas por la suite Fase D para
	// agrupar scenarios (ej. "rbac", "menu", "screen-config").
	Tags []string
}
