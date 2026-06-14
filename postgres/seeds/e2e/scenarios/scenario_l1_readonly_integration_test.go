//go:build integration
// +build integration

package scenarios_test

import (
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/internal/testdb"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/scenarios"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
)

// TestScenarioL1Readonly_Integration ejercita el scenario l1_readonly
// contra una BD efímera (testcontainers postgres:15-alpine con
// migrations + system seed completo L0..L4).
//
// MP-09 F4: L1 quedó como CONTRATO PURO. Ya no siembra DATO DE TENANT
// (escuela demo, usuario viewer, user_role, membership); sólo siembra
// el rol de contrato announcement_viewer. El permiso efectivo del rol
// (academic.announcements.read) se modela hoy con iam.role_grants
// (patterns glob), sembrado en L4 — NO con la tabla iam.role_permissions
// (eliminada en P4-1). El test valida el contrato read-only a nivel de
// ROL, no del viewer (que ya no es parte del contrato; su dato vivo
// vive en playground_v2/base).
//
// Cubre:
//
//   - F3-REQ-5.1: scenario se aplica sobre BD con el system seed.
//   - F3-REQ-5.3 (parte positiva): el rol announcement_viewer tiene el
//     grant allow `academic.announcements.read` en iam.role_grants.
//   - F3-REQ-5.3 (parte negativa) / F3-REQ-5.4 (parte SQL): el rol NO
//     tiene grants de escritura sobre announcements.
//
// La parte HTTP de F3-REQ-5.4 (403) y la validación end-to-end del
// usuario viewer son ahora responsabilidad de los tests de playground/
// API (el viewer es dato vivo, no contrato).
//
// Ejecución:
//
//	ENABLE_INTEGRATION_TESTS=true go test -tags=integration \
//	    -run TestScenarioL1Readonly -count=1 \
//	    ./seeds/e2e/scenarios/...
func TestScenarioL1Readonly_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("scenario_l1_readonly: skip en modo -short")
	}
	if !testdb.IntegrationGate() {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}

	gdb := testdb.StartPostgres(t)

	// 1. Aplicar el scenario `l1_readonly`. system.ApplySystem ya corrió
	//    en testdb.StartPostgres aplicando L0..L4; la fixture verifica
	//    presencia del rol de contrato y exporta constantes.
	reg := framework.NewRegistry()
	scenario := &scenarios.L1ReadOnly{}
	if err := reg.RegisterScenario(scenario); err != nil {
		t.Fatalf("RegisterScenario: %v", err)
	}
	composer := framework.NewComposer(reg, framework.NewNopLogger())

	ctx, err := composer.Apply(gdb, "l1_readonly")
	if err != nil {
		t.Fatalf("composer.Apply(l1_readonly): %v", err)
	}
	if ctx.ScenarioName != "l1_readonly" {
		t.Errorf("ScenarioName=%q, want %q", ctx.ScenarioName, "l1_readonly")
	}
	if len(ctx.Constants) == 0 {
		t.Error("ctx.Constants vacío tras Apply — la fixture l1_constants_export no llamó SetConstant")
	}

	// Idempotencia (C-REQ-1.4) — Apply 2ª vez no debe fallar.
	if _, err := composer.Apply(gdb, "l1_readonly"); err != nil {
		t.Fatalf("Apply 2 (idempotencia rota): %v", err)
	}

	// 2. F3-REQ-5.3 (parte positiva) — el rol announcement_viewer tiene
	//    EXACTAMENTE 1 grant allow `academic.announcements.read` en
	//    iam.role_grants.
	t.Run("F3-REQ-5.3_positive_read_grant", func(t *testing.T) {
		const q = `
SELECT COUNT(*)
FROM iam.role_grants rg
WHERE rg.role_id = ?::uuid AND rg.pattern = ? AND rg.effect = 'allow'
`
		var count int64
		if err := gdb.Raw(q, layers.L1_ROLE_ANNOUNCEMENT_VIEWER_ID, "academic.announcements.read").Scan(&count).Error; err != nil {
			t.Fatalf("query announcements.read grant: %v", err)
		}
		if count != 1 {
			t.Errorf("rol announcement_viewer tiene %d grants allow con academic.announcements.read; want 1", count)
		}
	})

	// 3. F3-REQ-5.3 (parte negativa) + F3-REQ-5.4 (parte SQL) — el rol NO
	//    tiene grants de escritura sobre announcements.
	t.Run("F3-REQ-5.3_negative_no_write_grants", func(t *testing.T) {
		const q = `
SELECT COUNT(*)
FROM iam.role_grants rg
WHERE rg.role_id = ?::uuid AND rg.pattern = ? AND rg.effect = 'allow'
`
		for _, pattern := range []string{
			"academic.announcements.create",
			"academic.announcements.update",
			"academic.announcements.delete",
		} {
			var count int64
			if err := gdb.Raw(q, layers.L1_ROLE_ANNOUNCEMENT_VIEWER_ID, pattern).Scan(&count).Error; err != nil {
				t.Fatalf("query %s grant: %v", pattern, err)
			}
			if count != 0 {
				t.Errorf("rol announcement_viewer tiene %d grants allow con %s; want 0 (gating read-only)", count, pattern)
			}
		}
	})
}
