package l4_test

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l4"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
	"github.com/EduGoGroup/edugo-shared/screenconfig"
)

// Harness de regresión F3.1 (plan 004-permisologia-mvp).
//
// Objetivo: garantizar que la migración de las screen_instances CRUD
// genéricas al "patrón delta" del SDUI (heredar default_actions del
// template + actions_removed) NO altera el CONJUNTO de acciones
// efectivas de ninguna pantalla. El label se normaliza al template a
// propósito (decisión del operador); por eso NO entra en la clave de
// invariante.
//
// Clave invariante por acción compuesta = (event_id_efectivo,
// permission, scope), donde event_id_efectivo se calcula como:
//   - el campo "event_id" explícito si está presente, o
//   - la inferencia del frontend: ids save/save_new/save_existing →
//     "submit-form"; el resto → el propio id.
//
// El test computa el composed real con
// screenconfig.ComposeActionsForResolve (el mismo que corre el backend
// en /screen-config/resolve) y compara contra un golden hardcodeado del
// estado POST-migración. El golden se generó ejecutando este mismo
// harness contra el estado PRE-migración (ver dumpGolden abajo); si una
// migración cambia el conjunto invariante de cualquier pantalla, el
// test falla con el diff exacto.

// invariantKey es la clave SEMÁNTICA estable de una acción compuesta.
//
// F3.1 (2ª pasada): la clave protege SOLO el contrato semántico de cada
// botón — qué evento dispara (event_id_efectivo) y con qué permiso se
// gatea (permission). Deliberadamente NO incluye los campos de
// PRESENTACIÓN (scope, label, condition, order, icon, style), porque la
// decisión del operador es "normalizar al template": esos campos ADOPTAN
// el canónico declarado en default_actions del template.
//
// Por qué scope salió del guard: scope es presentación, no semántica. El
// front KMP la trata como tal:
//   - normalizeScope convierte "form" → "form-submit" (los forms eran
//     render-equivalentes; el cambio inline→default es transparente).
//   - list-basic-v1 usa zonas de expansión dinámica (list_actions
//     scope="header", slots vacíos). Con `create` SIN scope inline, el
//     FAB de crear no se materializa hoy en el header (solo el botón del
//     empty-state); adoptar el scope canónico "header" del template hace
//     aparecer el FAB donde debe — es un ARREGLO (ausente→presente),
//     nunca una regresión.
//
// Lo único que NO puede cambiar al migrar es esta clave semántica: si el
// $resource$ expandido produjera un permission distinto al inline para
// algún id, la instancia NO es migrable por normalización pura.
type invariantKey struct {
	eventID    string
	permission string
}

func (k invariantKey) String() string {
	return fmt.Sprintf("%s|%s", k.eventID, k.permission)
}

// effectiveEventID replica la inferencia del frontend: las acciones de
// guardado colapsan al evento submit-form; el resto enrutan por su id.
func effectiveEventID(action map[string]any) string {
	if ev, ok := action["event_id"].(string); ok && ev != "" {
		return ev
	}
	id, _ := action["id"].(string)
	switch id {
	case "save", "save_new", "save_existing":
		return "submit-form"
	default:
		return id
	}
}

// collectAllInstances reúne las instancias de TODAS las capas
// (L0+L1+L2+L3+L4). L1 no siembra screen_instances.
func collectAllInstances(t *testing.T) []entities.ScreenInstance {
	t.Helper()
	var all []entities.ScreenInstance

	add := func(name string, fn func() ([]entities.ScreenInstance, error)) {
		rows, err := fn()
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		all = append(all, rows...)
	}
	add("L0ScreenInstances", layers.L0ScreenInstances)
	add("L2ScreenInstances", layers.L2ScreenInstances)
	add("L3ScreenInstances", layers.L3ScreenInstances)
	add("l4.ScreenInstances", l4.ScreenInstances)
	return all
}

// collectTemplatesByID indexa los templates de L0+L4 por su UUID string.
func collectTemplatesByID(t *testing.T) map[string]entities.ScreenTemplate {
	t.Helper()
	idx := make(map[string]entities.ScreenTemplate)

	add := func(name string, fn func() ([]entities.ScreenTemplate, error)) {
		rows, err := fn()
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		for _, tpl := range rows {
			idx[tpl.ID.String()] = tpl
		}
	}
	add("L0ScreenTemplates", layers.L0ScreenTemplates)
	add("l4.ScreenTemplates", l4.ScreenTemplates)
	return idx
}

// composeInvariants devuelve el conjunto ordenado de claves invariantes
// de una instancia, usando el composer real del backend.
func composeInvariants(t *testing.T, inst entities.ScreenInstance, templates map[string]entities.ScreenTemplate) []string {
	t.Helper()

	tpl, ok := templates[inst.TemplateID.String()]
	if !ok {
		t.Fatalf("instancia %q referencia template %s ausente del índice", inst.ScreenKey, inst.TemplateID)
	}

	reqPerm := ""
	if inst.RequiredPermission != nil {
		reqPerm = *inst.RequiredPermission
	}

	composedRaw, _ := screenconfig.ComposeActionsForResolve(inst.SlotData, tpl.Definition, reqPerm)

	var composed map[string]any
	if err := json.Unmarshal(composedRaw, &composed); err != nil {
		t.Fatalf("instancia %q: composed no es JSON válido: %v", inst.ScreenKey, err)
	}

	rawActions, _ := composed["actions"].([]any)
	// Conjunto (dedup) de claves semánticas: si dos acciones colapsan a
	// la misma (event_id_efectivo, permission) cuentan como una sola.
	seen := make(map[string]struct{}, len(rawActions))
	for _, a := range rawActions {
		action, ok := a.(map[string]any)
		if !ok {
			continue
		}
		perm, _ := action["permission"].(string)
		k := invariantKey{
			eventID:    effectiveEventID(action),
			permission: perm,
		}
		seen[k.String()] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// TestScreenActionsInvariantRoundTrip es el GATE de F3.1: para cada
// screen_instance del sistema, el conjunto de claves SEMÁNTICAS
// {event_id_efectivo, permission} debe coincidir con el golden.
func TestScreenActionsInvariantRoundTrip(t *testing.T) {
	templates := collectTemplatesByID(t)
	instances := collectAllInstances(t)

	if len(instances) == 0 {
		t.Fatal("no se recolectó ninguna screen_instance — el harness no validaría nada")
	}

	for _, inst := range instances {
		got := composeInvariants(t, inst, templates)
		want, known := goldenInvariants[inst.ScreenKey]
		if !known {
			t.Errorf("screen_key %q no está en el golden — agrega su set invariante a goldenInvariants (o revisa si es una instancia nueva)", inst.ScreenKey)
			continue
		}
		if !equalStringSets(got, want) {
			t.Errorf("screen_key %q: conjunto invariante cambió.\n  golden: %s\n  actual: %s",
				inst.ScreenKey, strings.Join(want, ", "), strings.Join(got, ", "))
		}
	}

	// Detecta golden huérfanos (keys que el golden declara pero ya no
	// existen como instancia) — protege contra eliminaciones silenciosas.
	present := make(map[string]bool, len(instances))
	for _, inst := range instances {
		present[inst.ScreenKey] = true
	}
	for key := range goldenInvariants {
		if !present[key] {
			t.Errorf("golden declara screen_key %q que ya no existe entre las instancias", key)
		}
	}
}

func equalStringSets(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestDumpGolden regenera el golden e imprime el literal Go listo para
// pegar en screen_actions_golden_test.go. Se skip-ea por defecto; para
// regenerar tras un cambio INTENCIONAL de acciones/permisos ejecutar:
//
//	DUMP_SCREEN_GOLDEN=1 go test ./seeds/system/l4/ -run TestDumpGolden -v
//
// Mantenerlo como Test (en vez de helper suelto) evita el falso positivo
// del linter `unused` y deja el procedimiento de regeneración ejecutable
// sin editar código.
func TestDumpGolden(t *testing.T) {
	if os.Getenv("DUMP_SCREEN_GOLDEN") == "" {
		t.Skip("set DUMP_SCREEN_GOLDEN=1 para regenerar el golden semántico")
	}
	templates := collectTemplatesByID(t)
	instances := collectAllInstances(t)

	keys := make([]string, 0, len(instances))
	byKey := make(map[string][]string, len(instances))
	for _, inst := range instances {
		keys = append(keys, inst.ScreenKey)
		byKey[inst.ScreenKey] = composeInvariants(t, inst, templates)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString("var goldenInvariants = map[string][]string{\n")
	for _, k := range keys {
		fmt.Fprintf(&b, "\t%q: {", k)
		quoted := make([]string, 0, len(byKey[k]))
		for _, v := range byKey[k] {
			quoted = append(quoted, fmt.Sprintf("%q", v))
		}
		b.WriteString(strings.Join(quoted, ", "))
		b.WriteString("},\n")
	}
	b.WriteString("}\n")
	t.Log("\n" + b.String())
}
