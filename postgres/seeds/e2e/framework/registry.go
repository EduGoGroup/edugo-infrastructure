package framework

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Registry es el catálogo de fixtures y scenarios disponibles para el
// composer. Es seguro para uso concurrente; en la práctica el binario
// seed_e2e lo puebla una sola vez al arrancar (RegisterAll), y los
// tests crean instancias propias mediante NewRegistry().
type Registry struct {
	mu        sync.RWMutex
	fixtures  map[string]Fixture
	scenarios map[string]Scenario

	// applying lleva el control de qué scenarios están en pleno Apply
	// para detectar reentradas concurrentes en el mismo proceso
	// (C-REQ-2.5, C-REQ-10.5).
	applying map[string]bool
}

// NewRegistry construye un registry vacío.
func NewRegistry() *Registry {
	return &Registry{
		fixtures:  map[string]Fixture{},
		scenarios: map[string]Scenario{},
		applying:  map[string]bool{},
	}
}

// DefaultRegistry es el registry global usado por el binario seed_e2e.
// Los tests deben preferir NewRegistry() para evitar contaminación
// cruzada.
var DefaultRegistry = NewRegistry()

// RegisterFixture añade una fixture al registry. Falla si ya existe
// otra con el mismo nombre (C-REQ-10.3 aplicado al catálogo).
func (r *Registry) RegisterFixture(f Fixture) error {
	if f == nil {
		return fmt.Errorf("registry: cannot register nil fixture")
	}
	manifest := f.Manifest()
	if manifest.Name == "" {
		return fmt.Errorf("registry: fixture has empty Name")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, dup := r.fixtures[manifest.Name]; dup {
		return fmt.Errorf("registry: duplicate fixture %q", manifest.Name)
	}
	r.fixtures[manifest.Name] = f
	return nil
}

// RegisterScenario añade un scenario al registry.
func (r *Registry) RegisterScenario(s Scenario) error {
	if s == nil {
		return fmt.Errorf("registry: cannot register nil scenario")
	}
	manifest := s.Manifest()
	if manifest.Name == "" {
		return fmt.Errorf("registry: scenario has empty Name")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, dup := r.scenarios[manifest.Name]; dup {
		return fmt.Errorf("registry: duplicate scenario %q", manifest.Name)
	}
	r.scenarios[manifest.Name] = s
	return nil
}

// LookupFixture devuelve la fixture o un error con la lista de
// fixtures disponibles (mensaje accionable, C-REQ-10.4 aplicado al
// catálogo).
func (r *Registry) LookupFixture(name string) (Fixture, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if f, ok := r.fixtures[name]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("registry: unknown fixture %q (available: %s)", name, joinSorted(r.fixtures))
}

// LookupScenario devuelve el scenario o un error
// `unregistered scenario: <name>` (C-REQ-10.3).
func (r *Registry) LookupScenario(name string) (Scenario, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if s, ok := r.scenarios[name]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("unregistered scenario: %s (available: %s)", name, joinSorted(r.scenarios))
}

// FixtureNames devuelve el listado ordenado de fixtures registradas
// (utilizado por catálogos y mensajes de error).
func (r *Registry) FixtureNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.fixtures))
	for n := range r.fixtures {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// ScenarioNames devuelve el listado ordenado de scenarios registrados.
func (r *Registry) ScenarioNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.scenarios))
	for n := range r.scenarios {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// AcquireApplyLock marca un scenario como "en aplicación" y devuelve
// un release que el composer debe invocar al terminar (típicamente con
// defer). Si otro Apply del mismo scenario está en curso, falla con
// `cleanup already in progress` style error (C-REQ-10.5).
func (r *Registry) AcquireApplyLock(scenarioName string) (release func(), err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.applying[scenarioName] {
		return nil, fmt.Errorf("scenario already in progress: %s", scenarioName)
	}
	r.applying[scenarioName] = true
	return func() {
		r.mu.Lock()
		delete(r.applying, scenarioName)
		r.mu.Unlock()
	}, nil
}

// joinSorted produce una lista ordenada y separada por comas con las
// claves del map; usado por los mensajes de error.
func joinSorted[V any](m map[string]V) string {
	if len(m) == 0 {
		return "<empty>"
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}
