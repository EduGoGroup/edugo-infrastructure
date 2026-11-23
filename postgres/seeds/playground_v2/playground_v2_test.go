package playground_v2_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/common"
	"github.com/google/uuid"
)

func TestOnTopScenarios(t *testing.T) {
	want := []string{"onboarding", "n1_inscripcion", "n17_secciones", "n4_evaluacion", "multi_unidad", "n0n1_escuelas"}
	if got := playground_v2.OnTopScenarios(); !slices.Equal(got, want) {
		t.Fatalf("OnTopScenarios() = %v, esperado %v", got, want)
	}
	avail := playground_v2.Available()
	for _, n := range append(want, "base", "mp10_material", "f6_representante", "demo") {
		if !slices.Contains(avail, n) {
			t.Errorf("Available() no incluye %q: %v", n, avail)
		}
	}
}

func TestScenarioUsers(t *testing.T) {
	wantPerScenario := map[string]int{
		"onboarding": 3, "n1_inscripcion": 5, "n17_secciones": 8,
		"n4_evaluacion": 5, "multi_unidad": 5, "n0n1_escuelas": 15,
	}
	total := 0
	seen := map[string]string{}
	for _, name := range playground_v2.OnTopScenarios() {
		users, err := playground_v2.ScenarioUsers(name)
		if err != nil {
			t.Fatalf("ScenarioUsers(%q): %v", name, err)
		}
		if len(users) != wantPerScenario[name] {
			t.Errorf("%s: %d usuarios, esperado %d", name, len(users), wantPerScenario[name])
		}
		for _, u := range users {
			if u.Email == "" || u.FirstName == "" || u.LastName == "" {
				t.Errorf("%s: usuario incompleto %+v", name, u)
			}
			if !strings.HasSuffix(u.Email, "edugo.local") {
				t.Errorf("%s: email fuera del dominio de playground: %s", name, u.Email)
			}
			if prev, dup := seen[u.Email]; dup {
				t.Errorf("email %s repetido en %s y %s", u.Email, prev, name)
			}
			seen[u.Email] = name
		}
		total += len(users)
	}
	if total != 41 {
		t.Errorf("total de usuarios on-top = %d, esperado 41", total)
	}

	for _, name := range []string{"base", "mp10_material", "demo", "all", "nope"} {
		if _, err := playground_v2.ScenarioUsers(name); err == nil {
			t.Errorf("ScenarioUsers(%q) debía fallar: no es on-top", name)
		}
	}
}

func TestUserIDsRequireListsAllMissing(t *testing.T) {
	users := []common.ScenarioUser{{Email: "b@x"}, {Email: "a@x"}, {Email: "c@x"}}
	ids := common.UserIDs{"c@x": uuid.New()}
	err := ids.Require(users)
	if err == nil || !strings.Contains(err.Error(), "a@x, b@x") || strings.Contains(err.Error(), "c@x") {
		t.Fatalf("Require() = %v; esperaba listar a@x, b@x ordenados y no c@x", err)
	}
	ids["a@x"], ids["b@x"] = uuid.New(), uuid.New()
	if err := ids.Require(users); err != nil {
		t.Fatalf("Require() con todos presentes = %v", err)
	}
}
