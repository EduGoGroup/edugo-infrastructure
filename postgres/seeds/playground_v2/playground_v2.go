// Package playground_v2 es la línea vigente de playgrounds focalizados
// de EduGo. Corren sobre el sistema completo (L0..L4) y se limitan a
// sembrar el envoltorio multi-tenant + roles/grants/usuarios para validar
// el CRUD sobre los recursos meta que L4 ya trae.
//
// Tiene registry propio, flag CLI propio (`--playground-v2`) y rangos
// UUID dedicados (62000000-..., 12000000-...). No participa de
// `ComputeFilesHash()` — cambiar un v2 no requiere bump de SchemaVersion.
//
// Tres clases de fixture (plan 053 T5.2 y HU-DEMO.5):
//   - Autónomos (`fixtures`): base, mp10_material, f6_representante. Siembran
//     todo lo que necesitan; base (y los que lo componen) trunca y resiembra
//     el mundo, incluido auth.users. Son los que expande "all".
//   - Autónomos fuera de "all" (`standalone`): demo, el mundo de la base de
//     demo (MVP-DEMO). Siembra todo, pero no trunca: se niega si encuentra otra
//     escuela, así que solo va sobre una base recién recreada.
//   - On-top (`onTop`): onboarding, n1_inscripcion, n17_secciones,
//     n4_evaluacion, multi_unidad, n0n1_escuelas. Se aplican ENCIMA de una BD
//     ya migrada con base, NO crean usuarios: sus usuarios (ScenarioUsers) los
//     crea antes la API de identidad (con su hash de contraseña) y ApplyOnTop
//     los resuelve por email y encola su sincronización con identity-core.
//
// Para agregar un nuevo playground v2:
//  1. Crear paquete `seeds/playground_v2/<name>/<name>.go` con una función
//     Apply idempotente: Apply(tx *gorm.DB) error si es autónomo, o
//     Users() []common.ScenarioUser + Apply(tx *gorm.DB, users common.UserIDs)
//     error si es on-top.
//  2. Agregar UNA línea al slice `fixtures`, `standalone` u `onTop`.
package playground_v2

import (
	"fmt"
	"strings"

	base "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/base"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/common"
	demo "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/demo"
	f6_representante "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/f6_representante"
	mp10_material "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/mp10_material"
	multi_unidad "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/multi_unidad"
	n0n1_escuelas "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/n0n1_escuelas"
	n17_secciones "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/n17_secciones"
	n1_inscripcion "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/n1_inscripcion"
	n4_evaluacion "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/n4_evaluacion"
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
//
// v0.3.0 (2026-11-16, plan 053 T5.2): los 6 escenarios con usuarios pasan a
// on-top. Dejan de insertar en auth.users (41 usuarios @edugo.local con UUID
// literal y bcrypt del seed); los resuelven por email y ApplyOnTop encola su
// sincronización. "all" deja de incluirlos.
//
// v0.4.0 (2026-11-17, plan 053 T6.2): base deja de sembrar password_hash (sus
// 10 usuarios y los 2 representantes) y de truncar auth.refresh_tokens; el
// login lo hace identity-core. Mismos usuarios, UUID y datos.
//
// v0.5.0 (2026-11-20, HU-DEMO.5): nace `demo`, el mundo de la base de demo, en
// un registro aparte (`standalone`) que "all" no expande. Los fixtures que ya
// existían no cambian.
const SeedVersion = "v0.5.0"

// ApplyFunc es la firma estable de la función Apply de cada playground v2
// autónomo.
type ApplyFunc func(*gorm.DB) error

// Fixture representa un playground v2 autónomo registrado.
type Fixture struct {
	Name  string
	Apply ApplyFunc
}

// ScenarioUser es un usuario que un escenario on-top necesita en auth.users.
// Quien aplica el escenario lo crea antes por la API de identidad.
type ScenarioUser struct{ Email, FirstName, LastName string }

// onTopFixture representa un escenario on-top registrado.
type onTopFixture struct {
	Name  string
	Users func() []common.ScenarioUser
	Apply func(*gorm.DB, common.UserIDs) error
}

// fixtures es el registry declarativo de playgrounds v2 autónomos. El orden se
// respeta cuando se aplica "all".
var fixtures = []Fixture{
	// base = mundo de datos por defecto (MP-09). Va primero en "all".
	{Name: "base", Apply: base.Apply},
	// mp10_material compone ENCIMA de base (material publicado para los hijos
	// del representante); valida la pantalla "Material del hijo" (plan 024 M3).
	// Invoca base.Apply: trunca el mundo, incluidos los escenarios on-top que
	// se hubieran aplicado antes. Aplicarlo ANTES de cualquier on-top.
	{Name: "mp10_material", Apply: mp10_material.Apply},
	// f6_representante compone ENCIMA de base (evaluaciones asignadas a la sesión
	// del hijo); valida la pantalla "Evaluaciones del hijo" (plan 024 / GET
	// /me/wards/assessments). Igual que mp10_material: invoca base.Apply, así que
	// va ANTES de cualquier on-top.
	{Name: "f6_representante", Apply: f6_representante.Apply},
}

// standalone es el registry de autónomos que "all" NO expande: se aplican solo
// por su nombre.
var standalone = []Fixture{
	// demo = el mundo de la demo en KMP (MVP-DEMO, HU-DEMO.5). Va solo en la base
	// de demo, recién recreada: no trunca y se niega si encuentra otra escuela.
	{Name: "demo", Apply: demo.Apply},
}

// onTop es el registry de escenarios on-top (plan 053 T5.2), en orden estable.
var onTop = []onTopFixture{
	{Name: "onboarding", Users: onboarding.Users, Apply: onboarding.Apply},
	{Name: "n1_inscripcion", Users: n1_inscripcion.Users, Apply: n1_inscripcion.Apply},
	{Name: "n17_secciones", Users: n17_secciones.Users, Apply: n17_secciones.Apply},
	{Name: "n4_evaluacion", Users: n4_evaluacion.Users, Apply: n4_evaluacion.Apply},
	{Name: "multi_unidad", Users: multi_unidad.Users, Apply: multi_unidad.Apply},
	{Name: "n0n1_escuelas", Users: n0n1_escuelas.Users, Apply: n0n1_escuelas.Apply},
}

// Available retorna los nombres de todos los playgrounds v2 (autónomos, los
// que quedan fuera de "all" y on-top).
func Available() []string {
	names := make([]string, 0, len(fixtures)+len(standalone)+len(onTop))
	for _, f := range fixtures {
		names = append(names, f.Name)
	}
	for _, f := range standalone {
		names = append(names, f.Name)
	}
	return append(names, OnTopScenarios()...)
}

// OnTopScenarios lista los escenarios que se aplican encima de base, sin force.
func OnTopScenarios() []string {
	names := make([]string, 0, len(onTop))
	for _, f := range onTop {
		names = append(names, f.Name)
	}
	return names
}

func findOnTop(name string) (onTopFixture, bool) {
	for _, f := range onTop {
		if f.Name == name {
			return f, true
		}
	}
	return onTopFixture{}, false
}

// ScenarioUsers devuelve los usuarios que el escenario necesita (error si el nombre no es on-top).
func ScenarioUsers(name string) ([]ScenarioUser, error) {
	f, ok := findOnTop(name)
	if !ok {
		return nil, fmt.Errorf("playground_v2: %q no es un escenario on-top (on-top: %v)", name, OnTopScenarios())
	}
	src := f.Users()
	out := make([]ScenarioUser, 0, len(src))
	for _, u := range src {
		out = append(out, ScenarioUser{Email: u.Email, FirstName: u.FirstName, LastName: u.LastName})
	}
	return out, nil
}

// enqueueIdentitySyncSQL es la SQL de Outbox.Enqueue de edugo-api-identity
// (internal/infra/postgres/identitysync/outbox.go), copiada tal cual: el
// escenario apunta el cambio de accesos de cada usuario igual que la API.
const enqueueIdentitySyncSQL = `INSERT INTO auth.identity_sync_outbox (user_id) VALUES (?::uuid)
		ON CONFLICT (user_id) DO UPDATE SET
			generation = auth.identity_sync_outbox.generation + 1,
			attempts = 0, next_attempt_at = now(), last_error = NULL, updated_at = now()`

// ApplyOnTop aplica el escenario sobre una BD ya migrada con base, sin DROP y sin
// borrar auth.users. Resuelve cada user_id por email en auth.users (deleted_at IS NULL);
// si falta alguno devuelve un error que lista TODOS los emails ausentes, sin escribir nada.
// Por cada usuario del escenario encola su sync en auth.identity_sync_outbox
// (misma SQL que Outbox.Enqueue de edugo-api-identity). Todo en una transacción.
//
// Re-ejecutable: las inserciones del escenario son OnConflict DoNothing por id
// y el outbox sube de generación. Límite conocido: las memberships e
// invitaciones llevan id fijo, así que si un usuario se borra y se vuelve a
// crear (UUID nuevo) sobre una BD ya sembrada, su membership antigua conserva
// el user_id viejo; recrear la BD (base) y volver a aplicar lo arregla.
func ApplyOnTop(gdb *gorm.DB, name string) error {
	f, ok := findOnTop(name)
	if !ok {
		return fmt.Errorf("playground_v2: %q no es un escenario on-top (on-top: %v)", name, OnTopScenarios())
	}
	users := f.Users()
	return gdb.Transaction(func(tx *gorm.DB) error {
		ids, missing, err := common.ResolveUserIDs(tx, users)
		if err != nil {
			return fmt.Errorf("playground_v2 %q: %w", name, err)
		}
		if len(missing) > 0 {
			return fmt.Errorf("playground_v2 %q (ApplyOnTop): faltan %d usuarios en auth.users; créalos antes con la API de identidad: %s",
				name, len(missing), strings.Join(missing, ", "))
		}
		if err := f.Apply(tx, ids); err != nil {
			return fmt.Errorf("playground_v2 %q: %w", name, err)
		}
		for _, u := range users {
			if err := tx.Exec(enqueueIdentitySyncSQL, ids[u.Email].String()).Error; err != nil {
				return fmt.Errorf("playground_v2 %q: encolando la sincronización de %s: %w", name, u.Email, err)
			}
		}
		return nil
	})
}

// Apply ejecuta el playground v2 identificado por name. Idempotente.
//
//   - "all" expande SOLO a los autónomos (base, mp10_material,
//     f6_representante), en el orden del registry. Los on-top quedan fuera:
//     "all" corre dentro de Migrate sobre una BD recién forzada, donde los
//     usuarios que crea la API no pueden existir todavía. demo (standalone)
//     también queda fuera: se aplica solo por su nombre.
//   - Un escenario on-top delega en ApplyOnTop. Dentro de Migrate (BD recién
//     recreada) falla siempre con la lista de usuarios ausentes: el camino es
//     migrar con base, crear los usuarios por la API y llamar a ApplyOnTop.
func Apply(gdb *gorm.DB, name string) error {
	if name == "all" {
		for _, f := range fixtures {
			if err := f.Apply(gdb); err != nil {
				return fmt.Errorf("playground_v2 %q: %w", f.Name, err)
			}
		}
		return nil
	}
	for _, list := range [][]Fixture{fixtures, standalone} {
		for _, f := range list {
			if f.Name == name {
				return f.Apply(gdb)
			}
		}
	}
	if _, ok := findOnTop(name); ok {
		return ApplyOnTop(gdb, name)
	}
	return fmt.Errorf("playground_v2: nombre desconocido %q (disponibles: %v, o \"all\")", name, Available())
}
