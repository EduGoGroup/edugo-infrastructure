package framework

import (
	"strings"
	"sync"
	"testing"
)

func TestRegistry_RegisterFixture_Duplicate(t *testing.T) {
	reg := NewRegistry()
	if err := reg.RegisterFixture(newFake("A", nil, nil)); err != nil {
		t.Fatalf("first register: %v", err)
	}
	err := reg.RegisterFixture(newFake("A", nil, nil))
	if err == nil || !strings.Contains(err.Error(), "duplicate fixture") {
		t.Fatalf("se esperaba duplicate fixture; got=%v", err)
	}
}

func TestRegistry_RegisterFixture_EmptyName(t *testing.T) {
	reg := NewRegistry()
	err := reg.RegisterFixture(&fakeFixture{})
	if err == nil || !strings.Contains(err.Error(), "empty Name") {
		t.Fatalf("se esperaba empty Name; got=%v", err)
	}
}

func TestRegistry_RegisterScenario_Duplicate(t *testing.T) {
	reg := NewRegistry()
	if err := reg.RegisterScenario(&fakeScenario{name: "x"}); err != nil {
		t.Fatalf("first register: %v", err)
	}
	err := reg.RegisterScenario(&fakeScenario{name: "x"})
	if err == nil || !strings.Contains(err.Error(), "duplicate scenario") {
		t.Fatalf("se esperaba duplicate scenario; got=%v", err)
	}
}

func TestRegistry_LookupFixture_Unknown(t *testing.T) {
	reg := NewRegistry()
	_, err := reg.LookupFixture("missing")
	if err == nil || !strings.Contains(err.Error(), "unknown fixture") {
		t.Fatalf("se esperaba unknown fixture; got=%v", err)
	}
}

func TestRegistry_LookupScenario_UnknownLists(t *testing.T) {
	reg := NewRegistry()
	if err := reg.RegisterScenario(&fakeScenario{name: "available_one"}); err != nil {
		t.Fatalf("register: %v", err)
	}
	_, err := reg.LookupScenario("missing")
	if err == nil || !strings.Contains(err.Error(), "unregistered scenario") {
		t.Fatalf("se esperaba unregistered scenario; got=%v", err)
	}
	if !strings.Contains(err.Error(), "available_one") {
		t.Errorf("error debe listar scenarios disponibles; got=%v", err)
	}
}

func TestRegistry_Names_Sorted(t *testing.T) {
	reg := NewRegistry()
	for _, n := range []string{"c", "a", "b"} {
		_ = reg.RegisterFixture(newFake(n, nil, nil))
		_ = reg.RegisterScenario(&fakeScenario{name: n})
	}
	got := reg.FixtureNames()
	if got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Errorf("FixtureNames no ordenado: %v", got)
	}
	got = reg.ScenarioNames()
	if got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Errorf("ScenarioNames no ordenado: %v", got)
	}
}

func TestRegistry_AcquireApplyLock_DoubleApply(t *testing.T) {
	reg := NewRegistry()
	rel, err := reg.AcquireApplyLock("x")
	if err != nil {
		t.Fatalf("first acquire: %v", err)
	}
	defer rel()
	_, err = reg.AcquireApplyLock("x")
	if err == nil || !strings.Contains(err.Error(), "already in progress") {
		t.Fatalf("se esperaba 'already in progress'; got=%v", err)
	}
}

func TestRegistry_AcquireApplyLock_ReleaseAllowsReentry(t *testing.T) {
	reg := NewRegistry()
	rel, err := reg.AcquireApplyLock("x")
	if err != nil {
		t.Fatalf("first acquire: %v", err)
	}
	rel()
	rel2, err := reg.AcquireApplyLock("x")
	if err != nil {
		t.Fatalf("second acquire después de release: %v", err)
	}
	rel2()
}

func TestRegistry_RegisterFixture_Nil(t *testing.T) {
	reg := NewRegistry()
	if err := reg.RegisterFixture(nil); err == nil {
		t.Fatal("se esperaba error con fixture nil")
	}
	if err := reg.RegisterScenario(nil); err == nil {
		t.Fatal("se esperaba error con scenario nil")
	}
}

func TestRegistry_ConcurrentRegister(t *testing.T) {
	reg := NewRegistry()
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = reg.RegisterFixture(newFake("dup", nil, nil))
		}()
	}
	wg.Wait()
	// Sólo uno debe haber ganado.
	if names := reg.FixtureNames(); len(names) != 1 {
		t.Fatalf("se esperaba 1 fixture registrada; got=%v", names)
	}
}
