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
// migrations + system seed que incluye L0 + L1).
//
// Cubre:
//
//   - F3-REQ-5.1: scenario se aplica sobre BD con L0 + L1.
//   - F3-REQ-5.3 (parte positiva): el viewer tiene exactamente el
//     permiso `announcements:read` vía la cadena
//     `auth.users → iam.user_roles → iam.roles → iam.role_permissions
//     → iam.permissions`.
//   - F3-REQ-5.3 (parte negativa) / F3-REQ-5.4 (parte SQL): el viewer
//     NO tiene `announcements:create|update|delete`.
//   - F3-REQ-6.2: el `iam.user_roles` del viewer tiene `school_id` no
//     nulo apuntando a la escuela demo L1.
//
// El sub-test http_403_diferido marca con t.Skip la parte HTTP de
// F3-REQ-5.4 (Opción A confirmada con el usuario: validación HTTP
// equivalente queda diferida hasta el cierre del plan, cuando se
// levante el API server).
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
	//    en testdb.StartPostgres aplicando L0 + L1; la fixture verifica
	//    presencia y exporta constantes.
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

	// 2. F3-REQ-5.3 (parte positiva) — la cadena
	//    user → user_role → role → role_permission → permission
	//    devuelve EXACTAMENTE 1 fila con `announcements:read` para el
	//    viewer.
	t.Run("F3-REQ-5.3_positive_read_permission", func(t *testing.T) {
		const q = `
SELECT COUNT(*)
FROM auth.users u
JOIN iam.user_roles ur ON ur.user_id = u.id AND ur.is_active = TRUE
JOIN iam.roles r ON r.id = ur.role_id AND r.is_active = TRUE
JOIN iam.role_permissions rp ON rp.role_id = r.id
JOIN iam.permissions p ON p.id = rp.permission_id AND p.is_active = TRUE
WHERE u.email = ? AND p.name = ?
`
		var count int64
		if err := gdb.Raw(q, layers.L1_VIEWER_EMAIL, "academic.announcements.read").Scan(&count).Error; err != nil {
			t.Fatalf("query announcements:read: %v", err)
		}
		if count != 1 {
			t.Errorf("viewer %q tiene %d filas con announcements:read; want 1", layers.L1_VIEWER_EMAIL, count)
		}
	})

	// 3. F3-REQ-5.3 (parte negativa) + F3-REQ-5.4 (parte SQL) — el viewer
	//    NO tiene `announcements:create|update|delete`.
	t.Run("F3-REQ-5.3_negative_no_write_permissions", func(t *testing.T) {
		const q = `
SELECT COUNT(*)
FROM auth.users u
JOIN iam.user_roles ur ON ur.user_id = u.id AND ur.is_active = TRUE
JOIN iam.roles r ON r.id = ur.role_id AND r.is_active = TRUE
JOIN iam.role_permissions rp ON rp.role_id = r.id
JOIN iam.permissions p ON p.id = rp.permission_id AND p.is_active = TRUE
WHERE u.email = ? AND p.name = ?
`
		for _, action := range []string{"academic.announcements.create", "academic.announcements.update", "academic.announcements.delete"} {
			var count int64
			if err := gdb.Raw(q, layers.L1_VIEWER_EMAIL, action).Scan(&count).Error; err != nil {
				t.Fatalf("query %s: %v", action, err)
			}
			if count != 0 {
				t.Errorf("viewer %q tiene %d filas con %s; want 0 (gating read-only)", layers.L1_VIEWER_EMAIL, count, action)
			}
		}
	})

	// 4. F3-REQ-6.2 — el user_role del viewer tiene school_id no nulo
	//    apuntando a la escuela demo L1.
	t.Run("F3-REQ-6.2_user_role_school_scoped", func(t *testing.T) {
		type row struct {
			SchoolID *string
		}
		const q = `
SELECT ur.school_id::text AS school_id
FROM iam.user_roles ur
JOIN iam.roles r ON r.id = ur.role_id
WHERE ur.user_id = ?::uuid AND r.scope = 'school'
`
		var rows []row
		if err := gdb.Raw(q, layers.L1_USER_VIEWER_ID).Scan(&rows).Error; err != nil {
			t.Fatalf("query user_roles del viewer: %v", err)
		}
		if len(rows) != 1 {
			t.Fatalf("user_roles para viewer con scope=school: got=%d, want 1", len(rows))
		}
		if rows[0].SchoolID == nil {
			t.Fatal("school_id IS NULL — violación del contrato scope=school (post_gorm.sql:~311)")
		}
		if *rows[0].SchoolID != layers.L1_SCHOOL_DEMO_ID {
			t.Errorf("school_id=%q, want %q (L1 demo school)", *rows[0].SchoolID, layers.L1_SCHOOL_DEMO_ID)
		}
	})

	// 5. F3-REQ-5.4 (parte HTTP) — diferida.
	t.Run("F3-REQ-5.4_http_403_deferred", func(t *testing.T) {
		t.Skip("Diferido hasta cierre del plan: requiere API server levantado (Opción A confirmada con usuario). " +
			"La parte SQL de F3-REQ-5.4 se valida en F3-REQ-5.3_negative_no_write_permissions.")
	})
}
