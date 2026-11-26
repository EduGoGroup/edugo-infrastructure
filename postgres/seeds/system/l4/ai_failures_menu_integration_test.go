//go:build integration
// +build integration

package l4_test

import (
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l4"
)

// TestL4_AIFailuresMenu_Seeded (HU-IA.7) verifica, tras sembrar L4, que existen
// el recurso ai_failures bajo content, su permiso, el grant literal de la
// docente, la pantalla ai-failures-list y su resource_screen por defecto; y que
// re-aplicar las cuatro piezas no duplica nada.
//
// Reusa startPostgresForL4DashboardHomeTest (mismo paquete de test).
func TestL4_AIFailuresMenu_Seeded(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}
	gdb := startPostgresForL4DashboardHomeTest(t)

	check := func(label string) {
		t.Helper()
		expectCount(t, gdb, label+": recurso", 1,
			`SELECT count(*) FROM iam.resources r JOIN iam.resources p ON p.id = r.parent_id
			 WHERE r.id = ? AND r.key = 'ai_failures' AND p.key = 'content' AND r.is_menu_visible`,
			l4.L4_RESOURCE_AI_FAILURES_ID)
		expectCount(t, gdb, label+": permiso", 1,
			`SELECT count(*) FROM iam.permissions WHERE id = ? AND name = 'content.ai_failures.review' AND resource_id = ?`,
			l4.L4_PERM_AI_FAILURES_REVIEW_ID, l4.L4_RESOURCE_AI_FAILURES_ID)
		expectCount(t, gdb, label+": grant de la docente", 1,
			`SELECT count(*) FROM iam.role_grants WHERE role_id = ? AND pattern = 'content.ai_failures.review' AND effect = 'allow'`,
			l4.L4_ROLE_TEACHER_ID)
		expectCount(t, gdb, label+": pantalla", 1,
			`SELECT count(*) FROM ui_config.screen_instances WHERE id = ? AND screen_key = 'ai-failures-list'`,
			l4.L4_SCREEN_INST_AI_FAILURES_LIST_ID)
		expectCount(t, gdb, label+": resource_screen", 1,
			`SELECT count(*) FROM ui_config.resource_screens WHERE resource_id = ? AND screen_key = 'ai-failures-list' AND is_default`,
			l4.L4_RESOURCE_AI_FAILURES_ID)
	}
	check("primera siembra")

	for name, apply := range map[string]func(*gorm.DB) error{
		"ApplyResources":        l4.ApplyResources,
		"ApplyRolesPermissions": l4.ApplyRolesPermissions,
		"ApplyScreenInstances":  l4.ApplyScreenInstances,
		"ApplyResourceScreens":  l4.ApplyResourceScreens,
	} {
		if err := apply(gdb); err != nil {
			t.Fatalf("re-apply %s: %v", name, err)
		}
	}
	check("segunda siembra")
}

func expectCount(t *testing.T, gdb *gorm.DB, label string, want int64, query string, args ...any) {
	t.Helper()
	var n int64
	if err := gdb.Raw(query, args...).Scan(&n).Error; err != nil {
		t.Fatalf("%s: %v", label, err)
	}
	if n != want {
		t.Errorf("%s: %d filas, se esperan %d", label, n, want)
	}
}
