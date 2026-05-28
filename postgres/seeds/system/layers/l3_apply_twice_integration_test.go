//go:build integration
// +build integration

package layers_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
)

// TestL3_ApplyTwice_Idempotent verifica que aplicar la capa L3 (sobre
// L0+L1+L2) dos veces consecutivas produce exactamente el mismo dataset
// que aplicarla una sola vez: 11 filas L3 sin duplicados.
//
// Justificación: la idempotencia de L3 es contrato (todos los inserts
// usan ON CONFLICT DO NOTHING). Si una de las funciones applyL3_*
// perdiera el OnConflict, este test lo detectaría inmediatamente.
//
// Adicionalmente verifica:
//
//   - F5-REQ-1.1: el recurso `materials` existe (1 fila en iam.resources).
//   - F5-REQ-2.1: 3 permisos materials:{read,create,update}; NO existe
//     `materials:delete` en iam.permissions (assertion explícita).
//   - F5-REQ-2.2: 3 role_permissions super_admin × cada permiso L3.
//   - F5-REQ-3.x: 2 ScreenInstances + 2 ResourceScreens para materials
//     (list default + form no-default).
//   - Invariantes macro post-L3:
//   - iam.resources total ≥ 2 (announcements + materials).
//   - ui_config.resource_screens para announcements sigue siendo 2
//     (list L0 + form L2) — L3 no debe contaminar ese conteo.
//   - Cadena L1 viewer→permissions sigue sin contener materials:*.
//
// NOTA Fase 6: el seed se limita a `L3_LAYER_NAME` para aislar este
// test del contenido de L4 (que sí siembra `materials:delete/download
// /publish` por decisión documentada de F6 B2 — completa la
// superficie del recurso materials para los roles de plataforma).
// La assertion negativa F5-REQ-2.1 sobre `materials:delete` se
// preserva como invariante de "L3 por sí solo".
//
// Ejecución:
//
//	ENABLE_INTEGRATION_TESTS=true \
//	    go test -tags=integration -run TestL3_ApplyTwice_Idempotent -count=1 \
//	        ./seeds/system/layers/...
//
// Replica el patrón inline de TestL0/L1/L2_ApplyTwice_Idempotent (el
// helper canónico seeds/e2e/internal/testdb no es importable desde aquí).
func TestL3_ApplyTwice_Idempotent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}

	gdb := startPostgresForL3Test(t)

	// migrations.Migrate(Force=true) ya aplica system.ApplySystem, que
	// con L3 registrada invoca NewL0().Apply + NewL1().Apply +
	// NewL2().Apply + NewL3().Apply una primera vez. Por tanto los
	// conteos canónicos L3 ya deben cumplirse.
	//
	// NOTA: NO llamamos a assertL0Counts / assertL1Counts / assertL2Counts
	// aquí (lección de Fase 4): esos helpers contienen invariantes
	// específicas de su capa que pueden desactualizarse al sumar L3 (el
	// recurso `materials` infla resources total, p.ej.). Validamos
	// directamente los conteos L3-aislados + invariantes macro
	// específicos abajo en assertL3Counts.
	assertL3Counts(t, gdb, "L3 test — after initial Migrate (Apply #1)")

	// Apply #2 — primera reaplicación explícita de L0 + L1 + L2 + L3.
	if err := layers.NewL0().Apply(gdb); err != nil {
		t.Fatalf("NewL0().Apply (#2): %v", err)
	}
	if err := layers.NewL1().Apply(gdb); err != nil {
		t.Fatalf("NewL1().Apply (#2): %v", err)
	}
	if err := layers.NewL2().Apply(gdb); err != nil {
		t.Fatalf("NewL2().Apply (#2): %v", err)
	}
	if err := layers.NewL3().Apply(gdb); err != nil {
		t.Fatalf("NewL3().Apply (#2): %v", err)
	}
	assertL3Counts(t, gdb, "after L0+L1+L2+L3 Apply #2")

	// Apply #3 — segunda reaplicación. Refuerza idempotencia.
	if err := layers.NewL0().Apply(gdb); err != nil {
		t.Fatalf("NewL0().Apply (#3): %v", err)
	}
	if err := layers.NewL1().Apply(gdb); err != nil {
		t.Fatalf("NewL1().Apply (#3): %v", err)
	}
	if err := layers.NewL2().Apply(gdb); err != nil {
		t.Fatalf("NewL2().Apply (#3): %v", err)
	}
	if err := layers.NewL3().Apply(gdb); err != nil {
		t.Fatalf("NewL3().Apply (#3): %v", err)
	}
	assertL3Counts(t, gdb, "after L0+L1+L2+L3 Apply #3")
}

// assertL3Counts valida los conteos canónicos de L3 (11 filas
// distribuidas en 5 tablas) filtrados por los IDs específicos de L3,
// más las invariantes macro post-L3 (cadena viewer, conteo de
// resource_screens para announcements).
func assertL3Counts(t *testing.T, gdb *gorm.DB, stage string) {
	t.Helper()

	type check struct {
		desc  string
		query string
		args  []any
		want  int64
	}

	checks := []check{
		// F5-REQ-1.1: resource materials existe.
		{
			desc:  "iam.resources [id=L3_RESOURCE_MATERIALS_ID]",
			query: `SELECT COUNT(*) FROM iam.resources WHERE id = ?::uuid`,
			args:  []any{layers.L3_RESOURCE_MATERIALS_ID},
			want:  1,
		},
		// F5-REQ-2.1: 3 permisos materials:{read,create,update}.
		{
			desc:  "iam.permissions [resource_id=L3_RESOURCE_MATERIALS_ID]",
			query: `SELECT COUNT(*) FROM iam.permissions WHERE resource_id = ?::uuid`,
			args:  []any{layers.L3_RESOURCE_MATERIALS_ID},
			want:  3,
		},
		// F5-REQ-2.1 (negativa): NO existe materials:delete.
		{
			desc:  "iam.permissions [name=materials:delete] (negativa)",
			query: `SELECT COUNT(*) FROM iam.permissions WHERE name = ?`,
			args:  []any{"content.materials.delete"},
			want:  0,
		},
		// F5-REQ-2.2: 3 role_permissions super_admin × materials.
		{
			desc:  "iam.role_permissions [role_id=L0 super_admin AND permission_id IN L3 perms]",
			query: `SELECT COUNT(*) FROM iam.role_permissions WHERE role_id = ?::uuid AND permission_id IN (?::uuid, ?::uuid, ?::uuid)`,
			args: []any{
				layers.L0_ROLE_SUPER_ADMIN_ID,
				layers.L3_PERM_MATERIALS_READ_ID,
				layers.L3_PERM_MATERIALS_CREATE_ID,
				layers.L3_PERM_MATERIALS_UPDATE_ID,
			},
			want: 3,
		},
		// F5-REQ-3.1/3.2: 2 ScreenInstances L3.
		{
			desc:  "ui_config.screen_instances [id IN (L3 materials-list, L3 material-form)]",
			query: `SELECT COUNT(*) FROM ui_config.screen_instances WHERE id IN (?::uuid, ?::uuid)`,
			args: []any{
				layers.L3_SCREEN_INSTANCE_MATERIALS_LIST_ID,
				layers.L3_SCREEN_INSTANCE_MATERIAL_FORM_ID,
			},
			want: 2,
		},
		// F5-REQ-3.3: total 2 resource_screens para materials.
		{
			desc:  "ui_config.resource_screens [resource_id=L3 materials]",
			query: `SELECT COUNT(*) FROM ui_config.resource_screens WHERE resource_id = ?::uuid`,
			args:  []any{layers.L3_RESOURCE_MATERIALS_ID},
			want:  2,
		},
		// F5-REQ-3.3 (list default).
		{
			desc:  "ui_config.resource_screens [resource_id=L3 materials, type=list, default=true]",
			query: `SELECT COUNT(*) FROM ui_config.resource_screens WHERE resource_id = ?::uuid AND screen_type = ? AND is_default = TRUE`,
			args:  []any{layers.L3_RESOURCE_MATERIALS_ID, "list"},
			want:  1,
		},
		// F5-REQ-3.3 (form no-default).
		{
			desc:  "ui_config.resource_screens [resource_id=L3 materials, type=form, default=false]",
			query: `SELECT COUNT(*) FROM ui_config.resource_screens WHERE resource_id = ?::uuid AND screen_type = ? AND is_default = FALSE`,
			args:  []any{layers.L3_RESOURCE_MATERIALS_ID, "form"},
			want:  1,
		},
	}

	for _, c := range checks {
		var got int64
		if err := gdb.Raw(c.query, c.args...).Scan(&got).Error; err != nil {
			t.Fatalf("[%s] count %s: %v", stage, c.desc, err)
		}
		if got != c.want {
			t.Errorf("[%s] %s: got %d, want %d", stage, c.desc, got, c.want)
		}
	}

	// Invariantes macro post-L3.

	// (1) iam.resources total ≥ 2 (announcements de L0 + materials de L3).
	var resourcesTotal int64
	if err := gdb.Raw(`SELECT COUNT(*) FROM iam.resources`).Scan(&resourcesTotal).Error; err != nil {
		t.Fatalf("[%s] count iam.resources total: %v", stage, err)
	}
	if resourcesTotal < 2 {
		t.Errorf("[%s] iam.resources total: got %d, want ≥ 2 (announcements + materials)", stage, resourcesTotal)
	}

	// (2) resource_screens para announcements sigue siendo 2 (list L0 + form L2).
	// Si L3 contaminara accidentalmente esa tabla con un resource_id mal
	// puesto, este conteo se rompería.
	var resourceScreensForAnnouncements int64
	if err := gdb.Raw(
		`SELECT COUNT(*) FROM ui_config.resource_screens WHERE resource_id = ?::uuid`,
		layers.L0_RESOURCE_ANNOUNCEMENTS_ID,
	).Scan(&resourceScreensForAnnouncements).Error; err != nil {
		t.Fatalf("[%s] count resource_screens for announcements: %v", stage, err)
	}
	if resourceScreensForAnnouncements != 2 {
		t.Errorf(
			"[%s] resource_screens for announcements: got %d, want 2 (1 list L0 + 1 form L2) — L3 no debe modificar ese set",
			stage, resourceScreensForAnnouncements,
		)
	}

	// (3) Cadena L1 viewer → permisos sigue sin contener materials:*.
	// No-regresión a nivel SQL de F5-REQ-2.3: tras L3 el viewer no
	// adquiere accidentalmente ningún permiso sobre materials.
	var viewerMaterialsCount int64
	if err := gdb.Raw(
		`SELECT COUNT(*) FROM iam.role_permissions rp `+
			`JOIN iam.permissions p ON rp.permission_id = p.id `+
			`WHERE rp.role_id = (`+
			`  SELECT ur.role_id FROM iam.user_roles ur `+
			`  JOIN auth.users u ON ur.user_id = u.id `+
			`  WHERE u.email = ? LIMIT 1`+
			`) AND p.name LIKE ?`,
		layers.L1_VIEWER_EMAIL, "materials:%",
	).Scan(&viewerMaterialsCount).Error; err != nil {
		t.Fatalf("[%s] count viewer materials:* permissions: %v", stage, err)
	}
	if viewerMaterialsCount != 0 {
		t.Errorf(
			"[%s] viewer materials:* permissions: got %d, want 0 (L3 no debe filtrar permisos materials:* al viewer)",
			stage, viewerMaterialsCount,
		)
	}
}

// startPostgresForL3Test levanta un contenedor postgres:15-alpine,
// aplica migrations.Migrate(Force=true) (que invoca system.ApplySystem
// y por tanto NewL0().Apply + NewL1().Apply + NewL2().Apply +
// NewL3().Apply una primera vez) y devuelve un *gorm.DB listo para
// usar. Replica el patrón de startPostgresForL{0,1,2}Test (el helper
// canónico seeds/e2e/internal/testdb está bajo internal/ y no es
// importable desde aquí).
func startPostgresForL3Test(tb testing.TB) *gorm.DB {
	tb.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		tb.Fatalf("l3 integration: testcontainers GenericContainer: %v", err)
	}
	tb.Cleanup(func() { _ = container.Terminate(ctx) })

	host, err := container.Host(ctx)
	if err != nil {
		tb.Fatalf("l3 integration: container.Host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		tb.Fatalf("l3 integration: container.MappedPort: %v", err)
	}
	connStr := fmt.Sprintf(
		"host=%s port=%s user=test password=test dbname=testdb sslmode=disable",
		host, port.Port(),
	)

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		tb.Fatalf("l3 integration: sql.Open: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		tb.Fatalf("l3 integration: ping: %v", err)
	}
	tb.Cleanup(func() { _ = sqlDB.Close() })

	if _, err := migrations.Migrate(sqlDB, migrations.MigrateOptions{
		Force:         true,
		DBUser:        "test",
		SeedUpToLayer: layers.L3_LAYER_NAME,
	}); err != nil {
		tb.Fatalf("l3 integration: migrations.Migrate: %v", err)
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		tb.Fatalf("l3 integration: gorm.Open: %v", err)
	}
	return gdb
}
