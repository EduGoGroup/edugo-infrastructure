package framework

// ApplyContext es el contexto compartido durante la aplicación de un
// scenario. Las fixtures lo consumen para resolver prefijos, reusar
// entidades ya creadas por fixtures previas y declarar las constantes
// que se exportarán al JSON consumible por los tests Kotlin.
//
// El composer crea un ApplyContext fresco al inicio de cada Apply y lo
// pasa a las fixtures en orden topológico. Un Cleanup utiliza el mismo
// tipo (con Provided vacío) para localizar las filas a borrar a partir
// del nombre del scenario.
type ApplyContext struct {
	// ScenarioName es el nombre del scenario al que pertenece esta
	// aplicación. Si la composición se invoca directamente (Composer
	// .Compose, sin scenario registrado) el campo lleva un nombre ad-hoc.
	ScenarioName string

	// TenantPrefix se aplica a los códigos visibles
	// (ej. "E2E-A1B2C3D4-SCHOOL-01").
	TenantPrefix string

	// SchemaPrefix se aplica a todos los UUIDs generados
	// (ej. "e2ea1b2c3d4-0000-0000-0000-000000000001").
	SchemaPrefix string

	// Provided acumula las entidades creadas hasta el momento.
	// Sirve para que una fixture detecte que otra ya creó la escuela
	// y reutilice su ID en vez de duplicarla (C-REQ-1.2).
	// La clave es la capacidad declarada en FixtureManifest.Provides.
	Provided map[string]ProvidedEntity

	// Constants acumula los pares clave→valor que se materializarán en
	// el JSON exportado al final de la aplicación (C-REQ-6.3).
	Constants map[string]string

	// RawParams es la copia de ScenarioManifest.Params para que las
	// fixtures puedan leer parámetros sin acoplarse al scenario.
	RawParams map[string]string
}

// NewApplyContext construye un ApplyContext con maps inicializados.
// Es la forma recomendada de crear un contexto desde tests.
func NewApplyContext(scenarioName, tenantPrefix, schemaPrefix string) *ApplyContext {
	return &ApplyContext{
		ScenarioName: scenarioName,
		TenantPrefix: tenantPrefix,
		SchemaPrefix: schemaPrefix,
		Provided:     map[string]ProvidedEntity{},
		Constants:    map[string]string{},
		RawParams:    map[string]string{},
	}
}

// Provide registra una entidad como provista por una fixture y la deja
// disponible para fixtures posteriores que la requieran. Si la
// capacidad ya estaba registrada por otra fixture con un Code distinto
// el composer detectará el conflicto durante resolve(); aquí se
// permite la sobreescritura silenciosa porque la deduplicación ya
// ocurrió antes de tocar la BD.
func (c *ApplyContext) Provide(capability string, entity ProvidedEntity) {
	if c.Provided == nil {
		c.Provided = map[string]ProvidedEntity{}
	}
	c.Provided[capability] = entity
}

// SetConstant agrega una clave/valor al mapa de constantes que se
// exportará al JSON. Se usa por las fixtures dentro de Apply.
func (c *ApplyContext) SetConstant(key, value string) {
	if c.Constants == nil {
		c.Constants = map[string]string{}
	}
	c.Constants[key] = value
}

// ProvidedEntity es la metadata mínima que una fixture comparte con el
// resto de la composición sobre una entidad que creó (o reusó).
type ProvidedEntity struct {
	// Kind clasifica la entidad: "school", "user", "role",
	// "academic_unit", etc.
	Kind string

	// ID es el UUID con SchemaPrefix aplicado.
	ID string

	// Code es el código visible con TenantPrefix aplicado, cuando
	// aplique. Vacío para entidades sin código (ej. role_permissions).
	Code string

	// Extra permite adjuntar metadata específica de la entidad
	// (ej. email del usuario, password en claro). Las claves son
	// libres pero se recomienda documentarlas en el manifest.
	Extra map[string]string
}
