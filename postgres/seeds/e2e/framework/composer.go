package framework

import (
	"fmt"
	"maps"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Composer aplica scenarios (o composiciones ad-hoc de fixtures)
// dentro de una transacción única, resolviendo dependencias por orden
// topológico y emitiendo logs estructurados.
//
// Uso típico:
//
//	composer := framework.NewComposer(framework.DefaultRegistry, nil)
//	ctx, err := composer.Apply(db, "teacher_grades_only")
type Composer struct {
	Registry *Registry
	Logger   Logger
	Now      func() time.Time
}

// NewComposer construye un composer con valores por defecto: logger
// JSON a stdout, reloj real. Si reg es nil se usa DefaultRegistry.
func NewComposer(reg *Registry, log Logger) *Composer {
	if reg == nil {
		reg = DefaultRegistry
	}
	if log == nil {
		log = NewJSONLogger()
	}
	return &Composer{Registry: reg, Logger: log, Now: time.Now}
}

// Apply ejecuta un scenario completo (lookup en registry + Compose).
// Toda la operación se envuelve en db.Transaction para que un fallo
// en cualquier fixture revierta el scenario por completo.
func (c *Composer) Apply(db *gorm.DB, scenarioName string) (*ApplyContext, error) {
	scenario, err := c.Registry.LookupScenario(scenarioName)
	if err != nil {
		return nil, err
	}

	release, err := c.Registry.AcquireApplyLock(scenarioName)
	if err != nil {
		return nil, err
	}
	defer release()

	tenantPrefix, schemaPrefix := Derive(scenarioName)
	ctx := NewApplyContext(scenarioName, tenantPrefix, schemaPrefix)

	manifest := scenario.Manifest()
	maps.Copy(ctx.RawParams, manifest.Params)

	fixtures := scenario.BuildFixtures(ctx)
	if len(fixtures) == 0 {
		// Scenario vacío es válido (C-REQ-1.5): no toca BD pero
		// devuelve ApplyContext válido para que el caller pueda
		// inspeccionar prefijos y constantes.
		c.Logger.Emit(LogEntry{
			Event:        EventScenarioApply,
			Scenario:     scenarioName,
			TenantPrefix: tenantPrefix,
		})
		c.Logger.Emit(LogEntry{
			Event:        EventScenarioDone,
			Scenario:     scenarioName,
			TenantPrefix: tenantPrefix,
		})
		return ctx, nil
	}

	resolved, err := resolve(fixtures)
	if err != nil {
		return nil, err
	}

	c.Logger.Emit(LogEntry{
		Event:        EventScenarioApply,
		Scenario:     scenarioName,
		TenantPrefix: tenantPrefix,
	})
	start := c.Now()
	if db == nil {
		return nil, fmt.Errorf("composer.Apply: nil db")
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		for _, f := range resolved {
			if err := c.applyFixture(tx, ctx, f); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	c.Logger.Emit(LogEntry{
		Event:        EventScenarioDone,
		Scenario:     scenarioName,
		TenantPrefix: tenantPrefix,
		DurationMs:   c.Now().Sub(start).Milliseconds(),
	})
	return ctx, nil
}

// Compose aplica una lista arbitraria de fixtures con un prefijo
// explícito. Pensado para tests que quieren combinar fixtures sin
// registrar un scenario completo.
func (c *Composer) Compose(db *gorm.DB, scenarioName string, fixtures []Fixture) (*ApplyContext, error) {
	if scenarioName == "" {
		return nil, fmt.Errorf("composer.Compose: empty scenarioName")
	}
	tenantPrefix, schemaPrefix := Derive(scenarioName)
	ctx := NewApplyContext(scenarioName, tenantPrefix, schemaPrefix)
	if len(fixtures) == 0 {
		return ctx, nil
	}
	resolved, err := resolve(fixtures)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, fmt.Errorf("composer.Compose: nil db")
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		for _, f := range resolved {
			if err := c.applyFixture(tx, ctx, f); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return ctx, nil
}

// applyFixture invoca Apply de una fixture, emite los eventos y captura
// errores con stage/last_table para diagnóstico.
func (c *Composer) applyFixture(tx *gorm.DB, ctx *ApplyContext, f Fixture) error {
	manifest := f.Manifest()
	start := c.Now()
	if err := f.Apply(tx, ctx); err != nil {
		c.Logger.Emit(LogEntry{
			Event:        EventFixtureError,
			Scenario:     ctx.ScenarioName,
			Fixture:      manifest.Name,
			TenantPrefix: ctx.TenantPrefix,
			Stage:        "apply",
			Error:        err.Error(),
		})
		return fmt.Errorf("fixture %q apply: %w", manifest.Name, err)
	}
	c.Logger.Emit(LogEntry{
		Event:        EventFixtureApply,
		Scenario:     ctx.ScenarioName,
		Fixture:      manifest.Name,
		TenantPrefix: ctx.TenantPrefix,
		Tables:       manifest.Tables,
		DurationMs:   c.Now().Sub(start).Milliseconds(),
	})
	return nil
}

// resolve calcula el orden topológico de la composición, detectando:
//
//   - provider conflict: dos fixtures declaran el mismo Provides con
//     manifests distintos (C-REQ-10.1).
//   - unsatisfied requirement: una fixture pide algo que ninguna otra
//     en la composición provee (C-REQ-1.3).
//   - dependency cycle: A requiere B y B requiere A.
//
// Si dos fixtures aportan capacidades disjuntas, el orden resultante
// respeta el orden de declaración (estable, salvo por dependencias).
func resolve(fixtures []Fixture) ([]Fixture, error) {
	if len(fixtures) == 0 {
		return nil, nil
	}

	// providerOf: capability -> índice de la fixture que la provee.
	providerOf := map[string]int{}
	manifests := make([]FixtureManifest, len(fixtures))
	for i, f := range fixtures {
		manifests[i] = f.Manifest()
		if manifests[i].Name == "" {
			return nil, fmt.Errorf("resolve: fixture #%d has empty manifest Name", i)
		}
		for _, cap := range manifests[i].Provides {
			if prev, dup := providerOf[cap]; dup {
				return nil, fmt.Errorf("provider conflict: %q claimed by [%s, %s]",
					cap, manifests[prev].Name, manifests[i].Name)
			}
			providerOf[cap] = i
		}
	}

	// Construir grafo de dependencias: arista i -> j si i depende de j.
	deps := make([][]int, len(fixtures))
	for i, m := range manifests {
		for _, req := range m.Requires {
			pIdx, ok := providerOf[req]
			if !ok {
				return nil, fmt.Errorf("unsatisfied requirement: %s requires %q (no fixture in composition provides it)", m.Name, req)
			}
			if pIdx == i {
				continue
			}
			deps[i] = append(deps[i], pIdx)
		}
	}

	// Topological sort con detección de ciclos. Usa visit recursivo
	// con estados unvisited/temporary/permanent.
	const (
		stUnvisited = 0
		stTemporary = 1
		stPermanent = 2
	)
	state := make([]int, len(fixtures))
	order := make([]int, 0, len(fixtures))

	var visit func(i int, path []int) error
	visit = func(i int, path []int) error {
		switch state[i] {
		case stPermanent:
			return nil
		case stTemporary:
			cycle := append([]int{}, path...)
			cycle = append(cycle, i)
			names := make([]string, len(cycle))
			for k, idx := range cycle {
				names[k] = manifests[idx].Name
			}
			return fmt.Errorf("dependency cycle: %s", joinPath(names))
		}
		state[i] = stTemporary
		// Orden estable: visitar primero por índice de fixture proveedora.
		sortedDeps := append([]int{}, deps[i]...)
		sort.Ints(sortedDeps)
		for _, d := range sortedDeps {
			if err := visit(d, append(path, i)); err != nil {
				return err
			}
		}
		state[i] = stPermanent
		order = append(order, i)
		return nil
	}

	for i := range fixtures {
		if err := visit(i, nil); err != nil {
			return nil, err
		}
	}

	out := make([]Fixture, len(order))
	for k, idx := range order {
		out[k] = fixtures[idx]
	}
	return out, nil
}

// joinPath formatea un ciclo como "A -> B -> C -> A".
func joinPath(names []string) string {
	return strings.Join(names, " -> ")
}
