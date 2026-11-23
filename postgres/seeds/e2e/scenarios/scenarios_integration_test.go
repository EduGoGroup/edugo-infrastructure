//go:build integration
// +build integration

package scenarios_test

import (
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/internal/testdb"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/scenarios"
)

// TestScenariosIntegration_ApplyAndCleanup ejercita el scenario
// l0_minimal contra una BD efímera (testcontainers postgres:15-alpine
// con migrations + production seed). Por cada scenario:
//
//  1. Aplica vía composer.Apply.
//  2. Verifica que el ApplyContext devuelto tiene el TenantPrefix y
//     SchemaPrefix derivados por framework.Derive.
//  3. Verifica que ApplyContext.Constants no está vacío (cada fixture
//     debería haber agregado al menos una constante).
//  4. Aplica una segunda vez (idempotencia: no debe fallar).
//  5. Cleanup vía Cleaner.Cleanup; verifica que no falla.
//
// l0_minimal solo compone la fixture `L0ConstantsExport`.
//
// Ejecución:
//
//	ENABLE_INTEGRATION_TESTS=true go test -tags=integration \
//	    -run TestScenariosIntegration -count=1 \
//	    ./seeds/e2e/scenarios/...
func TestScenariosIntegration_ApplyAndCleanup(t *testing.T) {
	if !testdb.IntegrationGate() {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}

	gdb := testdb.StartPostgres(t)

	// R-008: los casos observer_audits, teacher_grades_only y
	// guardian_views_child llevaban un t.Skip fijo y se borraron; quedan
	// para el MVP de tests.
	cases := []struct {
		name     string
		scenario framework.Scenario
	}{
		{"l0_minimal", &scenarios.L0Minimal{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reg := framework.NewRegistry()
			if err := reg.RegisterScenario(c.scenario); err != nil {
				t.Fatalf("RegisterScenario: %v", err)
			}
			composer := framework.NewComposer(reg, framework.NewNopLogger())

			// Apply 1ª vez.
			ctx, err := composer.Apply(gdb, c.name)
			if err != nil {
				t.Fatalf("Apply 1: %v", err)
			}
			if ctx.ScenarioName != c.name {
				t.Errorf("ScenarioName=%q, want %q", ctx.ScenarioName, c.name)
			}
			wantTenant, wantSchema := framework.Derive(c.name)
			if ctx.TenantPrefix != wantTenant {
				t.Errorf("TenantPrefix=%q, want %q", ctx.TenantPrefix, wantTenant)
			}
			if ctx.SchemaPrefix != wantSchema {
				t.Errorf("SchemaPrefix=%q, want %q", ctx.SchemaPrefix, wantSchema)
			}
			if len(ctx.Constants) == 0 {
				t.Error("ctx.Constants vacío tras Apply — alguna fixture no llamó SetConstant")
			}

			// Apply 2ª vez (idempotencia, C-REQ-1.4).
			if _, err := composer.Apply(gdb, c.name); err != nil {
				t.Fatalf("Apply 2 (idempotencia rota): %v", err)
			}

			// Cleanup.
			cleaner := framework.NewCleaner(reg, framework.NewNopLogger())
			if err := cleaner.Cleanup(gdb, c.name); err != nil {
				t.Fatalf("Cleanup: %v", err)
			}
		})
	}
}

// TestScenarios_ResolveCorrectly_NoIntegration es la versión sin BD —
// ya existe en scenarios_test.go (build sin tag); este archivo se
// limita a tests integration con BD.
