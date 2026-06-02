// Package playground_v2 es la segunda línea de playgrounds focalizados
// de EduGo, paralela a `seeds/playground/`. A diferencia de los v1 que
// se aplicaban sobre L0 con recursos+pantallas sembrados ad-hoc, los
// v2 corren sobre el sistema completo (L0..L4) y se limitan a sembrar
// el envoltorio multi-tenant + roles/grants/usuarios para validar el
// CRUD sobre los recursos meta que L4 ya trae.
//
// Convive con `playground/` sin pisarlo: registry propio, flag CLI
// propio (`--playground-v2`) y rangos UUID dedicados (62000000-...,
// 12000000-...). No participa de `ComputeFilesHash()` — cambiar un v2
// no requiere bump de SchemaVersion.
//
// Para agregar un nuevo playground v2:
//  1. Crear paquete `seeds/playground_v2/<name>/<name>.go` con
//     una función Apply(tx *gorm.DB) error idempotente.
//  2. Agregar UNA línea al slice `fixtures` con su nombre y la
//     referencia a la función Apply.
package playground_v2

import (
	"fmt"

	multi_unidad "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/multi_unidad"
	n0n1_escuelas "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/n0n1_escuelas"
	n17_secciones "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/n17_secciones"
	n1_inscripcion "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/n1_inscripcion"
	onboarding "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/onboarding"
	"gorm.io/gorm"
)

// SeedVersion declara la versión semántica del paquete playground_v2.
// Bumpear al cambiar la forma de un fixture existente.
//
// v0.2.0 (2026-06-02, ADR 0016): materia = catálogo de ESCUELA. Los fixtures
// n1_inscripcion, n17_secciones y multi_unidad dejan de anclar sus materias a
// una unidad (AcademicUnitID = NULL); n0n1_escuelas ya estaba alineado. Cumplen
// UNIQUE(school_id, name) sin deduplicar (nombres distintos por escuela).
const SeedVersion = "v0.2.0"

// ApplyFunc es la firma estable de la función Apply de cada playground v2.
type ApplyFunc func(*gorm.DB) error

// Fixture representa un playground v2 registrado.
type Fixture struct {
	Name  string
	Apply ApplyFunc
}

// fixtures es el registry declarativo de playgrounds v2 disponibles.
// El orden se respeta cuando se aplica "all".
var fixtures = []Fixture{
	{Name: "onboarding", Apply: onboarding.Apply},
	{Name: "n1_inscripcion", Apply: n1_inscripcion.Apply},
	{Name: "n17_secciones", Apply: n17_secciones.Apply},
	{Name: "multi_unidad", Apply: multi_unidad.Apply},
	{Name: "n0n1_escuelas", Apply: n0n1_escuelas.Apply},
}

// Available retorna los nombres de playgrounds v2 disponibles.
func Available() []string {
	names := make([]string, 0, len(fixtures))
	for _, f := range fixtures {
		names = append(names, f.Name)
	}
	return names
}

// Apply ejecuta el playground v2 identificado por name. "all" expande a
// todos los registrados, en el orden del registry. Idempotente.
func Apply(gdb *gorm.DB, name string) error {
	if name == "all" {
		for _, f := range fixtures {
			if err := f.Apply(gdb); err != nil {
				return fmt.Errorf("playground_v2 %q: %w", f.Name, err)
			}
		}
		return nil
	}
	for _, f := range fixtures {
		if f.Name == name {
			return f.Apply(gdb)
		}
	}
	return fmt.Errorf("playground_v2: nombre desconocido %q (disponibles: %v, o \"all\")", name, Available())
}
