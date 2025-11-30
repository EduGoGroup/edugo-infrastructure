package types

// TableToEntity mapea nombre de tabla SQL a nombre de entity Go
var TableToEntity = map[string]string{
	"users":              "User",
	"schools":            "School",
	"academic_units":     "AcademicUnit",
	"memberships":        "Membership",
	"materials":          "Material",
	"subjects":           "Subject",
	"units":              "Unit",
	"guardian_relations": "GuardianRelation",
}

// TableToCamel convierte nombre de tabla a CamelCase para estructuras
var TableToCamel = map[string]string{
	"users":              "Users",
	"schools":            "Schools",
	"academic_units":     "AcademicUnits",
	"memberships":        "Memberships",
	"materials":          "Materials",
	"subjects":           "Subjects",
	"units":              "Units",
	"guardian_relations": "GuardianRelations",
}

// GetEntityName retorna el nombre de la entity para una tabla
func GetEntityName(table string) string {
	if entity, ok := TableToEntity[table]; ok {
		return entity
	}
	return table
}

// GetTableCamel retorna el nombre CamelCase para una tabla
func GetTableCamel(table string) string {
	if camel, ok := TableToCamel[table]; ok {
		return camel
	}
	return table
}
