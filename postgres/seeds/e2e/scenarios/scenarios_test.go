package scenarios

import (
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
)

// Tests "puros": no tocan la base de datos. Validan que los Manifest
// son consistentes, que BuildFixtures devuelve el orden esperado y que
// los Provides/Requires de la composición resuelven sin
// "unsatisfied requirement". El test de Apply se detiene en
// `composer.Apply: nil db`, comprobando que la resolución previa fue
// exitosa.

// ----------------------------------------------------------------------
// Manifests
// ----------------------------------------------------------------------

func TestObserverAudits_Manifest(t *testing.T) {
	s := &ObserverAudits{}
	m := s.Manifest()
	if m.Name != "observer_audits" {
		t.Errorf("Name=%q, want observer_audits", m.Name)
	}
	if strings.TrimSpace(m.Description) == "" {
		t.Error("Description vacía")
	}
	wantTags := []string{"rbac", "menu", "audit"}
	if !equalStrSlice(m.Tags, wantTags) {
		t.Errorf("Tags=%v, want %v", m.Tags, wantTags)
	}
}

func TestTeacherGradesOnly_Manifest(t *testing.T) {
	s := &TeacherGradesOnly{}
	m := s.Manifest()
	if m.Name != "teacher_grades_only" {
		t.Errorf("Name=%q, want teacher_grades_only", m.Name)
	}
	if strings.TrimSpace(m.Description) == "" {
		t.Error("Description vacía")
	}
	wantTags := []string{"rbac", "screen-config"}
	if !equalStrSlice(m.Tags, wantTags) {
		t.Errorf("Tags=%v, want %v", m.Tags, wantTags)
	}
}

func TestGuardianViewsChild_Manifest(t *testing.T) {
	s := &GuardianViewsChild{}
	m := s.Manifest()
	if m.Name != "guardian_views_child" {
		t.Errorf("Name=%q, want guardian_views_child", m.Name)
	}
	if strings.TrimSpace(m.Description) == "" {
		t.Error("Description vacía")
	}
	wantTags := []string{"rbac", "screen-config"}
	if !equalStrSlice(m.Tags, wantTags) {
		t.Errorf("Tags=%v, want %v", m.Tags, wantTags)
	}
}

// ----------------------------------------------------------------------
// BuildFixtures: orden y composición
// ----------------------------------------------------------------------

func TestObserverAudits_BuildFixtures_Order(t *testing.T) {
	ctx := framework.NewApplyContext("observer_audits", "E2E-TEST-", "e2etest00-")
	fs := (&ObserverAudits{}).BuildFixtures(ctx)
	want := []string{"role_only", "readonly_role", "menu_subtree"}
	assertFixtureOrder(t, fs, want)
}

func TestTeacherGradesOnly_BuildFixtures_Order(t *testing.T) {
	ctx := framework.NewApplyContext("teacher_grades_only", "E2E-TEST-", "e2etest00-")
	fs := (&TeacherGradesOnly{}).BuildFixtures(ctx)
	want := []string{"role_only", "partial_crud", "screen_only"}
	assertFixtureOrder(t, fs, want)
}

func TestGuardianViewsChild_BuildFixtures_Order(t *testing.T) {
	ctx := framework.NewApplyContext("guardian_views_child", "E2E-TEST-", "e2etest00-")
	fs := (&GuardianViewsChild{}).BuildFixtures(ctx)
	want := []string{"role_only", "guardian_relation", "screen_only"}
	assertFixtureOrder(t, fs, want)
}

// ----------------------------------------------------------------------
// RegisterAll
// ----------------------------------------------------------------------

func TestRegisterAll_PopulatesRegistry(t *testing.T) {
	reg := framework.NewRegistry()
	if err := RegisterAll(reg); err != nil {
		t.Fatalf("RegisterAll error inesperado: %v", err)
	}
	wantNames := []string{
		"observer_audits",
		"teacher_grades_only",
		"guardian_views_child",
	}
	for _, name := range wantNames {
		if _, err := reg.LookupScenario(name); err != nil {
			t.Errorf("scenario %q no registrado: %v", name, err)
		}
	}
}

func TestRegisterAll_Idempotente(t *testing.T) {
	reg := framework.NewRegistry()
	if err := RegisterAll(reg); err != nil {
		t.Fatalf("primer RegisterAll error: %v", err)
	}
	err := RegisterAll(reg)
	if err == nil {
		t.Fatal("se esperaba error en segundo RegisterAll, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate scenario") {
		t.Errorf("error=%v, se esperaba contener 'duplicate scenario'", err)
	}
}

// ----------------------------------------------------------------------
// Composer.Apply: validar resolución de Provides/Requires
// ----------------------------------------------------------------------

func TestScenarios_ResolveCorrectly(t *testing.T) {
	cases := []string{
		"observer_audits",
		"teacher_grades_only",
		"guardian_views_child",
	}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			reg := framework.NewRegistry()
			if err := RegisterAll(reg); err != nil {
				t.Fatalf("RegisterAll: %v", err)
			}
			composer := framework.NewComposer(reg, framework.NewNopLogger())
			_, err := composer.Apply(nil, name)
			if err == nil {
				t.Fatal("se esperaba error 'nil db', got nil")
			}
			if strings.Contains(err.Error(), "unsatisfied requirement") {
				t.Fatalf("composición inválida (Provides/Requires inconsistentes): %v", err)
			}
			if strings.Contains(err.Error(), "provider conflict") {
				t.Fatalf("composición inválida (provider conflict): %v", err)
			}
			if strings.Contains(err.Error(), "dependency cycle") {
				t.Fatalf("composición inválida (dependency cycle): %v", err)
			}
			if !strings.Contains(err.Error(), "nil db") {
				t.Fatalf("se esperaba error 'nil db'; got=%v", err)
			}
		})
	}
}

// ----------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------

func assertFixtureOrder(t *testing.T, fs []framework.Fixture, want []string) {
	t.Helper()
	if len(fs) != len(want) {
		t.Fatalf("len(fixtures)=%d, want %d (got=%v)", len(fs), len(want), fixtureNames(fs))
	}
	for i, f := range fs {
		if got := f.Manifest().Name; got != want[i] {
			t.Errorf("fixtures[%d].Name=%q, want %q", i, got, want[i])
		}
	}
}

func fixtureNames(fs []framework.Fixture) []string {
	out := make([]string, len(fs))
	for i, f := range fs {
		out[i] = f.Manifest().Name
	}
	return out
}

func equalStrSlice(a, b []string) bool {
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
