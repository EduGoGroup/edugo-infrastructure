package framework

import (
	"strings"
	"testing"

	"gorm.io/gorm"
)

// fakeFixture es una fixture mínima que no toca la base de datos. Sirve
// para los tests de resolución de dependencias.
type fakeFixture struct {
	manifest FixtureManifest
	apply    func(tx *gorm.DB, ctx *ApplyContext) error
}

func (f *fakeFixture) Manifest() FixtureManifest { return f.manifest }
func (f *fakeFixture) Apply(tx *gorm.DB, ctx *ApplyContext) error {
	if f.apply != nil {
		return f.apply(tx, ctx)
	}
	return nil
}
func (f *fakeFixture) Cleanup(tx *gorm.DB, ctx *ApplyContext) error { return nil }

func newFake(name string, provides, requires []string) *fakeFixture {
	return &fakeFixture{manifest: FixtureManifest{Name: name, Provides: provides, Requires: requires}}
}

func TestResolve_TopologicalOrder(t *testing.T) {
	a := newFake("A", []string{"school"}, nil)
	b := newFake("B", []string{"user"}, []string{"school"})
	c := newFake("C", nil, []string{"user", "school"})

	ordered, err := resolve([]Fixture{c, b, a})
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}
	if len(ordered) != 3 {
		t.Fatalf("len(ordered) = %d, want 3", len(ordered))
	}
	indexOf := map[string]int{}
	for i, f := range ordered {
		indexOf[f.Manifest().Name] = i
	}
	if indexOf["A"] >= indexOf["B"] {
		t.Errorf("A debería ir antes que B")
	}
	if indexOf["B"] >= indexOf["C"] {
		t.Errorf("B debería ir antes que C")
	}
}

func TestResolve_ProviderConflict(t *testing.T) {
	a := newFake("A", []string{"school"}, nil)
	b := newFake("B", []string{"school"}, nil)
	_, err := resolve([]Fixture{a, b})
	if err == nil || !strings.Contains(err.Error(), "provider conflict") {
		t.Fatalf("se esperaba provider conflict; got=%v", err)
	}
}

func TestResolve_UnsatisfiedRequirement(t *testing.T) {
	a := newFake("A", nil, []string{"school"})
	_, err := resolve([]Fixture{a})
	if err == nil || !strings.Contains(err.Error(), "unsatisfied requirement") {
		t.Fatalf("se esperaba unsatisfied requirement; got=%v", err)
	}
}

func TestResolve_DependencyCycle(t *testing.T) {
	a := newFake("A", []string{"x"}, []string{"y"})
	b := newFake("B", []string{"y"}, []string{"x"})
	_, err := resolve([]Fixture{a, b})
	if err == nil || !strings.Contains(err.Error(), "dependency cycle") {
		t.Fatalf("se esperaba dependency cycle; got=%v", err)
	}
}

func TestResolve_EmptyComposition(t *testing.T) {
	out, err := resolve(nil)
	if err != nil {
		t.Fatalf("resolve(nil) error: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("resolve(nil) len = %d, want 0", len(out))
	}
}

func TestResolve_EmptyManifestName(t *testing.T) {
	bad := &fakeFixture{manifest: FixtureManifest{}}
	_, err := resolve([]Fixture{bad})
	if err == nil || !strings.Contains(err.Error(), "empty manifest Name") {
		t.Fatalf("se esperaba error de manifest vacío; got=%v", err)
	}
}

func TestResolve_StableOrderForIndependent(t *testing.T) {
	a := newFake("A", []string{"x"}, nil)
	b := newFake("B", []string{"y"}, nil)
	out, err := resolve([]Fixture{a, b})
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}
	if out[0].Manifest().Name != "A" || out[1].Manifest().Name != "B" {
		t.Errorf("orden no estable: %s,%s", out[0].Manifest().Name, out[1].Manifest().Name)
	}
}

// TestComposer_Apply_NilDB asegura que Apply no intenta tocar la BD si
// db es nil (evita que un test mal configurado entre a la transacción).
func TestComposer_Apply_NilDB(t *testing.T) {
	reg := NewRegistry()
	if err := reg.RegisterScenario(&fakeScenario{name: "noop"}); err != nil {
		t.Fatalf("register: %v", err)
	}
	c := NewComposer(reg, NewNopLogger())
	_, err := c.Apply(nil, "noop")
	if err == nil {
		t.Fatal("se esperaba error con db nil")
	}
}

// fakeScenario implementa Scenario con una fixture trivial que sí
// requiere DB para Apply (forzamos el path nil-db).
type fakeScenario struct {
	name string
}

func (s *fakeScenario) Manifest() ScenarioManifest {
	return ScenarioManifest{Name: s.name, Description: "fake"}
}
func (s *fakeScenario) BuildFixtures(_ *ApplyContext) []Fixture {
	return []Fixture{newFake("only", []string{"x"}, nil)}
}

// TestComposer_Apply_EmptyScenario_NoDB confirma C-REQ-1.5: un scenario
// que no produce fixtures es no-op y no toca la BD.
func TestComposer_Apply_EmptyScenario_NoDB(t *testing.T) {
	reg := NewRegistry()
	if err := reg.RegisterScenario(&emptyScenario{}); err != nil {
		t.Fatalf("register: %v", err)
	}
	c := NewComposer(reg, NewMemoryLogger())
	ctx, err := c.Apply(nil, "empty")
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if ctx == nil {
		t.Fatal("ctx debería ser no-nil")
	}
	if ctx.ScenarioName != "empty" {
		t.Errorf("ScenarioName=%q, want %q", ctx.ScenarioName, "empty")
	}
}

type emptyScenario struct{}

func (emptyScenario) Manifest() ScenarioManifest         { return ScenarioManifest{Name: "empty"} }
func (emptyScenario) BuildFixtures(*ApplyContext) []Fixture { return nil }

func TestComposer_Apply_UnknownScenario(t *testing.T) {
	c := NewComposer(NewRegistry(), NewNopLogger())
	_, err := c.Apply(nil, "missing")
	if err == nil || !strings.Contains(err.Error(), "unregistered scenario") {
		t.Fatalf("se esperaba unregistered scenario; got=%v", err)
	}
}

func TestComposer_Compose_EmptyName(t *testing.T) {
	c := NewComposer(NewRegistry(), NewNopLogger())
	_, err := c.Compose(nil, "", nil)
	if err == nil {
		t.Fatal("se esperaba error con scenarioName vacío")
	}
}

func TestComposer_Compose_EmptyFixtures_NoDB(t *testing.T) {
	c := NewComposer(NewRegistry(), NewNopLogger())
	ctx, err := c.Compose(nil, "ad_hoc_test", nil)
	if err != nil {
		t.Fatalf("Compose error: %v", err)
	}
	if ctx.ScenarioName != "ad_hoc_test" {
		t.Errorf("ScenarioName=%q", ctx.ScenarioName)
	}
}

func TestJoinPath(t *testing.T) {
	got := joinPath([]string{"A", "B", "A"})
	if got != "A -> B -> A" {
		t.Errorf("joinPath got=%q", got)
	}
	if joinPath(nil) != "" {
		t.Error("joinPath(nil) should be empty")
	}
}
