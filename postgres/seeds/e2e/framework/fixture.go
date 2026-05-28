package framework

import "gorm.io/gorm"

// Fixture es la pieza atómica componible del seed E2E. Cada fixture
// declara estáticamente lo que necesita y lo que aporta (Manifest), y
// se aplica/limpia siempre en el contexto de un Scenario que ya derivó
// los prefijos de aislamiento.
//
// Reglas que toda implementación debe respetar (validadas por tests):
//
//   - Manifest debe ser puro: el mismo valor cada vez que se invoca.
//   - Apply debe ser idempotente — aplicar la fixture N veces produce
//     el mismo estado final que aplicarla una sola vez (C-REQ-1.4,
//     C-REQ-4.4).
//   - Apply jamás escribe filas cuyo UUID caiga en el namespace del
//     production seed (10000000-..., c1000000-...). Para detectar el
//     bug a tiempo, antes de cada INSERT se debe llamar a
//     AssertNotProductionNamespace (C-REQ-10.2).
//   - Cleanup borra exclusivamente filas cuyo identificador o código
//     contenga TenantPrefix/SchemaPrefix del scenario (C-REQ-3.1, 3.2).
//   - Cleanup debe ser un no-op seguro si Apply nunca corrió antes
//     (C-REQ-3.3).
type Fixture interface {
	// Manifest devuelve el descriptor estático de la fixture (orden,
	// dependencias, tablas tocadas, constantes exportadas).
	Manifest() FixtureManifest

	// Apply ejecuta la fixture sobre tx, usando ctx para resolver
	// prefijos y reusar entidades ya creadas por fixtures previas
	// dentro del mismo scenario.
	Apply(tx *gorm.DB, ctx *ApplyContext) error

	// Cleanup borra (siempre por prefijo del scenario) lo que esta
	// fixture creó. Debe ser segura aún si Apply nunca corrió.
	Cleanup(tx *gorm.DB, ctx *ApplyContext) error
}

// FixtureManifest describe estáticamente una fixture. El composer lo
// usa para resolver dependencias (Provides/Requires), construir el
// orden de cleanup (Tables) y poblar el JSON de constantes exportadas.
type FixtureManifest struct {
	// Name identifica unívocamente a la fixture dentro del registry.
	// Convención: snake_case sin prefijo "fixture_" (ej. "role_only").
	Name string

	// Provides enumera las capacidades que esta fixture aporta a la
	// composición. Otras fixtures pueden referenciarlas en Requires.
	// Ejemplos: ["school", "user", "user_role", "membership"].
	Provides []string

	// Requires enumera las capacidades que necesitan ser provistas por
	// otra fixture en la misma composición. Si al resolver no se
	// satisface, el composer falla con `unsatisfied requirement`.
	Requires []string

	// Tables lista las tablas (schema.table) que la fixture toca.
	// El cleanup borra siguiendo el orden inverso de declaración.
	Tables []string

	// Constants son pares clave→valor que se exportarán al JSON de
	// constantes (seeds/e2e/exports/fixtures-constants.json). El
	// framework enriquece automáticamente las plantillas con
	// TenantPrefix y SchemaPrefix. Las claves siguen el patrón
	// E2E<Fixture><EntityKind> (C-REQ-6.2).
	Constants map[string]string

	// Description es un docstring corto utilizado por el catálogo y la
	// documentación generada.
	Description string
}
