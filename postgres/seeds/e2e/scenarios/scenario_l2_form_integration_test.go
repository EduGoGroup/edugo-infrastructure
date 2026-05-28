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

// TestScenarioL2Form_Integration ejercita el scenario l2_form contra
// una BD efímera (testcontainers postgres:15-alpine con migrations +
// system seed que incluye L0 + L1 + L2).
//
// Cubre:
//
//   - F4-REQ-1.1: ScreenInstance announcement-form existe con
//     template_id = L0_SCREEN_TPL_FORM_ID y screen_key correcto.
//   - F4-REQ-1.2: slot_data es JSON válido con 3 fields (title, body,
//     published_at), 3 actions y api_prefix=platform.
//   - F4-REQ-2.1: ResourceScreen (resource=announcements,
//     screen_type=form, is_default=false) existe con el id canónico de
//     L2.
//   - F4-REQ-3.1: action SAVE_NEW lleva permission=announcements:create.
//   - F4-REQ-3.2: action SAVE_EXISTING lleva
//     permission=announcements:update.
//   - No-regresión F3-REQ-5.3: tras aplicar L2 el viewer L1 sigue
//     teniendo EXACTAMENTE el permiso announcements:read.
//
// Toda la validación SQL está implementada en la fixture
// l2_constants_export.Apply (invocada por composer.Apply). Si alguna
// assertion falla, composer.Apply retorna error y el test rompe en el
// primer Apply.
//
// Los sub-tests http_*_deferred / ui_*_deferred marcan con t.Skip las
// partes diferidas por Opción A (HTTP/UI requieren API server o KMP).
//
// Ejecución:
//
//	ENABLE_INTEGRATION_TESTS=true go test -tags=integration \
//	    -run TestScenarioL2Form -count=1 \
//	    ./seeds/e2e/scenarios/...
func TestScenarioL2Form_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("scenario_l2_form: skip en modo -short")
	}
	if !testdb.IntegrationGate() {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}

	gdb := testdb.StartPostgres(t)

	// 1. Aplicar el scenario `l2_form`. system.ApplySystem ya corrió en
	//    testdb.StartPostgres aplicando L0 + L1 + L2; la fixture
	//    l2_constants_export verifica presencia y forma de las filas L2,
	//    valida la no-regresión L1 y exporta constantes.
	reg := framework.NewRegistry()
	scenario := &scenarios.L2Form{}
	if err := reg.RegisterScenario(scenario); err != nil {
		t.Fatalf("RegisterScenario: %v", err)
	}
	composer := framework.NewComposer(reg, framework.NewNopLogger())

	ctx, err := composer.Apply(gdb, "l2_form")
	if err != nil {
		t.Fatalf("composer.Apply(l2_form): %v", err)
	}
	if ctx.ScenarioName != "l2_form" {
		t.Errorf("ScenarioName=%q, want %q", ctx.ScenarioName, "l2_form")
	}
	if len(ctx.Constants) == 0 {
		t.Error("ctx.Constants vacío tras Apply — la fixture l2_constants_export no llamó SetConstant")
	}

	// Idempotencia (C-REQ-1.4) — Apply 2ª vez no debe fallar.
	if _, err := composer.Apply(gdb, "l2_form"); err != nil {
		t.Fatalf("Apply 2 (idempotencia rota): %v", err)
	}

	// 2. F4-REQ-2.1 (macro) — para el recurso announcements la tabla
	//    ui_config.resource_screens contiene EXACTAMENTE 2 filas (la
	//    list de L0 + la form de L2). Si L2 perdiera idempotencia este
	//    conteo subiría.
	t.Run("F4-REQ-2.1_resource_screens_count_for_announcements", func(t *testing.T) {
		var count int64
		if err := gdb.Raw(
			`SELECT COUNT(*) FROM ui_config.resource_screens WHERE resource_id = ?::uuid`,
			layers.L0_RESOURCE_ANNOUNCEMENTS_ID,
		).Scan(&count).Error; err != nil {
			t.Fatalf("query resource_screens for announcements: %v", err)
		}
		if count != 2 {
			t.Errorf("resource_screens for announcements: got %d, want 2 (1 list L0 + 1 form L2)", count)
		}
	})

	// 3. Constantes L2 exportadas correctamente al ApplyContext (las
	//    consumirán los tests Kotlin del KMP vía fixtures-constants.json).
	t.Run("constants_exported_to_context", func(t *testing.T) {
		wantKeys := []string{
			"E2EFixtureL2ScreenInstanceAnnouncementFormID",
			"E2EFixtureL2ResourceScreenAnnouncementsFormID",
			"E2EFixtureL2ScreenKeyAnnouncementForm",
		}
		for _, k := range wantKeys {
			v, ok := ctx.Constants[k]
			if !ok {
				t.Errorf("constante %q ausente en ctx.Constants tras Apply", k)
				continue
			}
			if v == "" {
				t.Errorf("constante %q exportada con valor vacío", k)
			}
		}
	})

	// 4. Sub-tests diferidos por Opción A (HTTP/UI fuera de scope SQL).
	t.Run("F4-REQ-1.3_http_screen_resolve_deferred", func(t *testing.T) {
		t.Skip("HTTP/UI deferred per Opción A — requires API server")
	})
	t.Run("F4-REQ-3.3_kmp_hides_save_button_for_viewer_deferred", func(t *testing.T) {
		t.Skip("HTTP/UI deferred per Opción A — requires API server")
	})
	t.Run("F4-REQ-5.2_post_201_get_menu_super_admin_deferred", func(t *testing.T) {
		t.Skip("HTTP/UI deferred per Opción A — requires API server")
	})
	t.Run("F4-REQ-5.3_post_403_viewer_deferred", func(t *testing.T) {
		t.Skip("HTTP/UI deferred per Opción A — requires API server")
	})
}
