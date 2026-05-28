package framework

import (
	"testing"

	"gorm.io/gorm"
)

// nopFakeFixture es una fixture que no toca tx — ideal para verificar
// que el composer emita los eventos `fixture.apply` correctamente sin
// requerir una BD.
type nopFakeFixture struct {
	manifest FixtureManifest
}

func (f *nopFakeFixture) Manifest() FixtureManifest { return f.manifest }
func (f *nopFakeFixture) Apply(tx *gorm.DB, ctx *ApplyContext) error {
	if ctx != nil {
		ctx.SetConstant("touched", "yes")
	}
	return nil
}
func (f *nopFakeFixture) Cleanup(tx *gorm.DB, ctx *ApplyContext) error { return nil }

// TestComposer_ApplyFixture_EmitsApplyEvent cubre C-REQ-8.1: cada
// fixture aplicada produce un log estructurado con event=fixture.apply
// que incluye scenario, fixture, tenant_prefix y la lista de tablas
// declaradas en el Manifest.
func TestComposer_ApplyFixture_EmitsApplyEvent(t *testing.T) {
	mem := NewMemoryLogger()
	c := NewComposer(NewRegistry(), mem)

	ctx := NewApplyContext("ad_hoc", "E2E-XXXXX-", "e2eXXXXX-")
	f := &nopFakeFixture{manifest: FixtureManifest{
		Name:   "ad_hoc",
		Tables: []string{"academic.schools", "auth.users"},
	}}

	if err := c.applyFixture(nil, ctx, f); err != nil {
		t.Fatalf("applyFixture: %v", err)
	}

	got := mem.Captured()
	if len(got) != 1 {
		t.Fatalf("se esperaba 1 evento; got=%d (%+v)", len(got), got)
	}
	ev := got[0]
	if ev.Event != EventFixtureApply {
		t.Errorf("event=%v, want %v", ev.Event, EventFixtureApply)
	}
	if ev.Fixture != "ad_hoc" {
		t.Errorf("fixture=%q", ev.Fixture)
	}
	if ev.Scenario != "ad_hoc" {
		t.Errorf("scenario=%q", ev.Scenario)
	}
	if ev.TenantPrefix != "E2E-XXXXX-" {
		t.Errorf("tenant_prefix=%q", ev.TenantPrefix)
	}
	if len(ev.Tables) != 2 {
		t.Errorf("tables=%v, esperaban 2", ev.Tables)
	}
	if ev.Time.IsZero() {
		t.Error("time debería poblarse")
	}
}

// errFakeFixture devuelve siempre un error desde Apply, ideal para
// validar el path `fixture.error` (C-REQ-8.3).
type errFakeFixture struct {
	manifest FixtureManifest
	err      error
}

func (f *errFakeFixture) Manifest() FixtureManifest                    { return f.manifest }
func (f *errFakeFixture) Apply(tx *gorm.DB, ctx *ApplyContext) error   { return f.err }
func (f *errFakeFixture) Cleanup(tx *gorm.DB, ctx *ApplyContext) error { return nil }

// TestComposer_ApplyFixture_EmitsErrorEvent cubre C-REQ-8.3: cuando
// una fixture falla, el composer emite `fixture.error` con stage=apply
// antes de propagar el error.
func TestComposer_ApplyFixture_EmitsErrorEvent(t *testing.T) {
	mem := NewMemoryLogger()
	c := NewComposer(NewRegistry(), mem)
	ctx := NewApplyContext("err", "E2E-", "e2e00000-")
	f := &errFakeFixture{
		manifest: FixtureManifest{Name: "boom"},
		err:      errMockApply,
	}
	err := c.applyFixture(nil, ctx, f)
	if err == nil {
		t.Fatal("se esperaba error")
	}
	got := mem.Captured()
	if len(got) != 1 || got[0].Event != EventFixtureError {
		t.Fatalf("esperaba un evento fixture.error; got=%+v", got)
	}
	if got[0].Stage != "apply" {
		t.Errorf("stage=%q, want apply", got[0].Stage)
	}
	if got[0].Error == "" {
		t.Error("error debería estar poblado")
	}
}

var errMockApply = mockError("simulated apply failure")

type mockError string

func (e mockError) Error() string { return string(e) }
