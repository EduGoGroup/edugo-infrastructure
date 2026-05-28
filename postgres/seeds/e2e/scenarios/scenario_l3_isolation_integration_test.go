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

// TestScenarioL3Isolation_Integration ejercita el scenario l3_isolation
// contra una BD efímera (testcontainers postgres:15-alpine con
// migrations + system seed que incluye L0 + L1 + L2 + L3).
//
// Cubre por SQL (en L3IsolationConstants.Apply, invocada por
// composer.Apply):
//
//   - F5-REQ-1.1: resource materials existe con scope=unit.
//   - F5-REQ-2.1: 3 permisos materials:{read,create,update}; ausencia
//     explícita de materials:delete.
//   - F5-REQ-2.2: 3 role_permissions super_admin × materials.
//   - F5-REQ-3.1: ScreenInstance materials-list con slot_data correcto.
//   - F5-REQ-3.2: ScreenInstance material-form con slot_data correcto.
//   - F5-REQ-3.3: 2 resource_screens (list default + form no-default).
//   - No-regresión L1: viewer sigue con EXACTAMENTE {announcements:read}.
//
// Los sub-tests *_deferred marcan con t.Skip las partes diferidas por
// Opción A (HTTP/UI requieren API server o KMP runtime).
//
// Ejecución:
//
//	ENABLE_INTEGRATION_TESTS=true go test -tags=integration \
//	    -run TestScenarioL3Isolation -count=1 \
//	    ./seeds/e2e/scenarios/...
func TestScenarioL3Isolation_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("scenario_l3_isolation: skip en modo -short")
	}
	if !testdb.IntegrationGate() {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}

	// Limitamos el system seed a L3 (upTo=L3_LAYER_NAME) para preservar
	// el aislamiento que valida este scenario: L4 introduce
	// materials:{delete,download,publish} (decisión documentada de F6
	// B2) y rompería la assertion negativa F5-REQ-2.1 sobre
	// `materials:delete`.
	gdb := testdb.StartPostgresUpTo(t, layers.L3_LAYER_NAME)

	// 1. Aplicar el scenario `l3_isolation`. system.ApplySystem ya corrió
	//    en testdb.StartPostgres aplicando L0 + L1 + L2 + L3; la fixture
	//    l3_constants_export verifica presencia y forma de las filas L3,
	//    valida la no-regresión L1 viewer y exporta constantes.
	reg := framework.NewRegistry()
	scenario := &scenarios.L3Isolation{}
	if err := reg.RegisterScenario(scenario); err != nil {
		t.Fatalf("RegisterScenario: %v", err)
	}
	composer := framework.NewComposer(reg, framework.NewNopLogger())

	ctx, err := composer.Apply(gdb, "l3_isolation")
	if err != nil {
		t.Fatalf("composer.Apply(l3_isolation): %v", err)
	}
	if ctx.ScenarioName != "l3_isolation" {
		t.Errorf("ScenarioName=%q, want %q", ctx.ScenarioName, "l3_isolation")
	}
	if len(ctx.Constants) == 0 {
		t.Error("ctx.Constants vacío tras Apply — la fixture l3_constants_export no llamó SetConstant")
	}

	// Idempotencia (C-REQ-1.4) — Apply 2ª vez no debe fallar.
	if _, err := composer.Apply(gdb, "l3_isolation"); err != nil {
		t.Fatalf("Apply 2 (idempotencia rota): %v", err)
	}

	// F5-REQ-1.1 — resource materials existe.
	t.Run("F5-REQ-1.1_resource_materials_exists", func(t *testing.T) {
		var count int64
		if err := gdb.Raw(
			`SELECT COUNT(*) FROM iam.resources WHERE id = ?::uuid`,
			layers.L3_RESOURCE_MATERIALS_ID,
		).Scan(&count).Error; err != nil {
			t.Fatalf("query iam.resources: %v", err)
		}
		if count != 1 {
			t.Errorf("iam.resources [id=L3_RESOURCE_MATERIALS_ID]: got %d, want 1", count)
		}
	})

	// F5-REQ-2.1 (negativa) — materials:delete NO existe.
	t.Run("F5-REQ-2.1_no_materials_delete_permission", func(t *testing.T) {
		var count int64
		if err := gdb.Raw(
			`SELECT COUNT(*) FROM iam.permissions WHERE name = ?`,
			"content.materials.delete",
		).Scan(&count).Error; err != nil {
			t.Fatalf("query iam.permissions[materials:delete]: %v", err)
		}
		if count != 0 {
			t.Errorf("iam.permissions[name=materials:delete]: got %d, want 0 (L3 valida CRUD parcial sin :delete)", count)
		}
	})

	// F5-REQ-3.3 — total resource_screens para materials = 2.
	t.Run("F5-REQ-3.3_resource_screens_count_for_materials", func(t *testing.T) {
		var count int64
		if err := gdb.Raw(
			`SELECT COUNT(*) FROM ui_config.resource_screens WHERE resource_id = ?::uuid`,
			layers.L3_RESOURCE_MATERIALS_ID,
		).Scan(&count).Error; err != nil {
			t.Fatalf("query resource_screens for materials: %v", err)
		}
		if count != 2 {
			t.Errorf("resource_screens for materials: got %d, want 2 (1 list default + 1 form no-default)", count)
		}
	})

	// Constantes L3 exportadas correctamente al ApplyContext (las
	// consumirán los tests Kotlin del KMP vía fixtures-constants.json).
	t.Run("constants_exported_to_context", func(t *testing.T) {
		wantConstants := map[string]string{
			"E2EFixtureL3ResourceMaterialsID":           layers.L3_RESOURCE_MATERIALS_ID,
			"E2EFixtureL3ResourceMaterialsKey":          layers.L3_RESOURCE_MATERIALS_KEY,
			"E2EFixtureL3PermMaterialsReadID":           layers.L3_PERM_MATERIALS_READ_ID,
			"E2EFixtureL3PermMaterialsCreateID":         layers.L3_PERM_MATERIALS_CREATE_ID,
			"E2EFixtureL3PermMaterialsUpdateID":         layers.L3_PERM_MATERIALS_UPDATE_ID,
			"E2EFixtureL3ScreenInstanceMaterialsListID": layers.L3_SCREEN_INSTANCE_MATERIALS_LIST_ID,
			"E2EFixtureL3ScreenInstanceMaterialFormID":  layers.L3_SCREEN_INSTANCE_MATERIAL_FORM_ID,
			"E2EFixtureL3ScreenKeyMaterialsList":        layers.L3_SCREEN_KEY_MATERIALS_LIST,
			"E2EFixtureL3ScreenKeyMaterialForm":         layers.L3_SCREEN_KEY_MATERIAL_FORM,
		}
		for k, want := range wantConstants {
			got, ok := ctx.Constants[k]
			if !ok {
				t.Errorf("constante %q ausente en ctx.Constants tras Apply", k)
				continue
			}
			if got != want {
				t.Errorf("constante %q: got %q, want %q", k, got, want)
			}
		}
	})

	// Sub-tests diferidos por Opción A (HTTP/UI fuera de scope SQL).
	t.Run("F5-REQ-2.3_viewer_menu_no_materials_deferred", func(t *testing.T) {
		t.Skip("HTTP/UI deferred per Opción A — requires API server")
	})
	t.Run("F5-REQ-4.1_super_admin_menu_two_items_deferred", func(t *testing.T) {
		t.Skip("HTTP/UI deferred per Opción A — requires API server")
	})
	t.Run("F5-REQ-4.2_viewer_menu_one_item_deferred", func(t *testing.T) {
		t.Skip("HTTP/UI deferred per Opción A — requires API server")
	})
	t.Run("F5-REQ-4.3_viewer_delete_returns_404_deferred", func(t *testing.T) {
		t.Skip("HTTP/UI deferred per Opción A — requires API server")
	})
	t.Run("F5-REQ-6.2_super_admin_crud_deferred", func(t *testing.T) {
		t.Skip("HTTP/UI deferred per Opción A — requires API server")
	})
	t.Run("F5-REQ-6.3_viewer_get_forbidden_deferred", func(t *testing.T) {
		t.Skip("HTTP/UI deferred per Opción A — requires API server")
	})
}
