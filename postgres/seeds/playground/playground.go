// Package playground contiene datasets pequeños y focalizados para iterar
// sobre EduGo carpeta a carpeta. Cada subdirectorio (admin/, focal_pantalla/,
// crud_full/, ...) expone un Fixture que se aplica encima de L0 (el piso del
// sistema).
//
// El paquete NO se incluye en seeds.ComputeFilesHash(): cambiar un playground
// no requiere bump de SchemaVersion. Sólo la versión de L0 (o cualquier capa
// de system/) participa del hash que valida el migrator.
//
// Para agregar un nuevo playground:
//  1. Crear paquete en seeds/playground/<name>/<name>.go con una función
//     Apply(tx *gorm.DB) error idempotente.
//  2. Agregar UNA línea al slice `fixtures` de este archivo con su nombre
//     y la referencia a la función Apply. No hace falta tocar nada más —
//     el resto del migrator descubre el playground por enumerar el registry.
//
// "all" como nombre de playground es alias para aplicar todos los registrados
// en orden, alineado con la convención de fotos acumulativas
// (project_edugo_playgrounds_convention en memoria).
package playground

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground/admin"
	focal_pantalla "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground/focal_pantalla"
	"gorm.io/gorm"
)

// SeedVersion declara la versión semántica del paquete playground.
// Bumpear al cambiar la forma de un fixture existente.
const SeedVersion = "v0.4.0"

// ApplyFunc es la firma estable de la función Apply de cada playground.
type ApplyFunc func(*gorm.DB) error

// Fixture representa un playground registrado: nombre + función Apply.
type Fixture struct {
	Name  string
	Apply ApplyFunc
}

// fixtures es el registry declarativo de playgrounds disponibles.
// Para agregar uno nuevo: agregar una línea acá apuntando a tu paquete.
// El orden se respeta cuando se aplica "all".
var fixtures = []Fixture{
	{Name: "admin", Apply: admin.Apply},
	{Name: "focal-pantalla", Apply: focal_pantalla.Apply},
	// ELIMINADOS en N4 F1 (plan 015) — cadena anclada al contrato viejo de
	// evaluación (entities.Assessment con created_by_user_id / subject texto-libre,
	// ahora demolido):
	//   - focal-evaluacion / focal-evaluacion-v2: sembraban la evaluación vieja.
	//   - focal-botonera (deprecado, plan 004) chainea focal_evaluacion_v2 para su
	//     tenant; focal-static-screens chainea focal_botonera. Toda la cadena era
	//     dead-by-transitivity y rompía el build. Se reconstruye en F2 sobre el
	//     esquema nuevo (assessment.* por sesión) si se necesita.
}

// Available retorna los nombres de playgrounds disponibles, en el orden
// del registry. Útil para CLIs y help messages.
func Available() []string {
	names := make([]string, 0, len(fixtures))
	for _, f := range fixtures {
		names = append(names, f.Name)
	}
	return names
}

// Apply ejecuta el playground identificado por name. "all" expande a todos
// los registrados, en el orden del registry. Devuelve error si el nombre es
// desconocido. Idempotente — cada Apply reaplica sin duplicar.
func Apply(gdb *gorm.DB, name string) error {
	if name == "all" {
		for _, f := range fixtures {
			if err := f.Apply(gdb); err != nil {
				return fmt.Errorf("playground %q: %w", f.Name, err)
			}
		}
		return nil
	}
	for _, f := range fixtures {
		if f.Name == name {
			return f.Apply(gdb)
		}
	}
	return fmt.Errorf("playground: nombre desconocido %q (disponibles: %v, o \"all\")", name, Available())
}
