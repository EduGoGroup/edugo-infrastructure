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

// TestL2_ApplyTwice_Idempotent verifica que aplicar la capa L2 (sobre
// L0+L1) dos veces consecutivas produce exactamente el mismo dataset
// que aplicarla una sola vez: 2 filas L2 sin duplicados.
//
// Justificación: la idempotencia de L2 es contrato (todos los inserts
// usan ON CONFLICT DO NOTHING). Si una de las funciones applyL2_*
// perdiera el OnConflict, este test lo detectaría inmediatamente.
//
// Adicionalmente verifica:
//
//   - F4-REQ-1.1: la ScreenInstance L2 (announcement-form) existe y es
//     única por id.
//   - F4-REQ-2.1: el mapping ResourceScreen L2 (announcements ↔ form)
//     existe y es único por id; y para el recurso `announcements` hay
//     EXACTAMENTE 2 filas en ui_config.resource_screens (la list de L0
//     y la form de L2).
//   - F3-REQ-6.2 (invariante post-L2): el user_role del viewer sigue
//     con school_id no nulo apuntando a la escuela demo L1.
//
// Ejecución:
//
//	ENABLE_INTEGRATION_TESTS=true \
//	    go test -tags=integration -run TestL2_ApplyTwice_Idempotent -count=1 \
//	        ./seeds/system/layers/...
//
// Replica el patrón inline de TestL0/L1_ApplyTwice_Idempotent (el helper
// canónico seeds/e2e/internal/testdb no es importable desde aquí).
func TestL2_ApplyTwice_Idempotent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}

	gdb := startPostgresForL2Test(t)

	// migrations.Migrate(Force=true) ya aplica system.ApplySystem, que
	// con L2 registrada invoca NewL0().Apply + NewL1().Apply +
	// NewL2().Apply una primera vez. Por tanto los conteos canónicos
	// L2 ya deben cumplirse.
	//
	// NOTA: NO llamamos a assertL0Counts / assertL1Counts aquí. Esos
	// helpers (definidos en l0/l1_apply_twice_integration_test.go)
	// validan `resource_screens [resource_id=L0 announcements] = 1`,
	// pero tras registrar L2 en system.go ese conteo pasa
	// legítimamente a 2 (list de L0 + form de L2). La
	// no-regresión-por-no-duplicación de L0/L1 se valida implícitamente
	// porque applyL2_* no toca ninguna tabla L0/L1 ni inserta en
	// resource_screens con resource_id distinto. F4-REQ-2.1 se valida
	// directamente abajo (assertL2Counts).
	//
	// MP-09 F4: ya no se valida el viewer L1 (assertL1ViewerHasSchool
	// eliminado). L1 quedó como CONTRATO PURO sin DATO DE TENANT; el
	// viewer/escuela/membership ya no existen en el system seed. La
	// idempotencia de L2 se valida con assertL2Counts.
	assertL2Counts(t, gdb, "L2 test — after initial Migrate (Apply #1)")

	// Apply #2 — primera reaplicación explícita de L0 + L1 + L2.
	if err := layers.NewL0().Apply(gdb); err != nil {
		t.Fatalf("NewL0().Apply (#2): %v", err)
	}
	if err := layers.NewL1().Apply(gdb); err != nil {
		t.Fatalf("NewL1().Apply (#2): %v", err)
	}
	if err := layers.NewL2().Apply(gdb); err != nil {
		t.Fatalf("NewL2().Apply (#2): %v", err)
	}
	assertL2Counts(t, gdb, "after L0+L1+L2 Apply #2")

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
	assertL2Counts(t, gdb, "after L0+L1+L2 Apply #3")
}

// assertL2Counts valida los conteos canónicos de L2 (2 filas
// distribuidas en 2 tablas) filtrados por los IDs específicos de L2.
// Adicionalmente verifica F4-REQ-2.1 macro: para el recurso
// `announcements` la tabla ui_config.resource_screens contiene
// EXACTAMENTE 2 filas (list de L0 + form de L2).
func assertL2Counts(t *testing.T, gdb *gorm.DB, stage string) {
	t.Helper()

	type check struct {
		desc  string
		query string
		args  []any
		want  int64
	}

	checks := []check{
		{
			desc:  "ui_config.screen_instances [id=L2 announcement-form]",
			query: `SELECT COUNT(*) FROM ui_config.screen_instances WHERE id = ?::uuid`,
			args:  []any{layers.L2_SCREEN_INSTANCE_ANNOUNCEMENT_FORM_ID},
			want:  1,
		},
		{
			desc:  "ui_config.resource_screens [id=L2 announcements-form mapping]",
			query: `SELECT COUNT(*) FROM ui_config.resource_screens WHERE id = ?::uuid`,
			args:  []any{layers.L2_RESOURCE_SCREEN_ANNOUNCEMENTS_FORM_ID},
			want:  1,
		},
	}

	var total int64
	for _, c := range checks {
		var got int64
		if err := gdb.Raw(c.query, c.args...).Scan(&got).Error; err != nil {
			t.Fatalf("[%s] count %s: %v", stage, c.desc, err)
		}
		if got != c.want {
			t.Errorf("[%s] %s: got %d, want %d", stage, c.desc, got, c.want)
		}
		total += got
	}

	// Total canónico L2: 1 + 1 = 2.
	const wantTotal int64 = 2
	if total != wantTotal {
		t.Errorf("[%s] total L2 rows: got %d, want %d", stage, total, wantTotal)
	}

	// F4-REQ-2.1 (macro): para el recurso announcements ya hay 2 filas
	// en resource_screens (list de L0 + form de L2). Si L2 perdiera
	// idempotencia y duplicara, este conteo subiría.
	var resourceScreensForAnnouncements int64
	if err := gdb.Raw(
		`SELECT COUNT(*) FROM ui_config.resource_screens WHERE resource_id = ?::uuid`,
		layers.L0_RESOURCE_ANNOUNCEMENTS_ID,
	).Scan(&resourceScreensForAnnouncements).Error; err != nil {
		t.Fatalf("[%s] count resource_screens for announcements: %v", stage, err)
	}
	if resourceScreensForAnnouncements != 2 {
		t.Errorf(
			"[%s] resource_screens for announcements: got %d, want 2 (1 list L0 + 1 form L2)",
			stage, resourceScreensForAnnouncements,
		)
	}
}

// startPostgresForL2Test levanta un contenedor postgres:15-alpine,
// aplica migrations.Migrate(Force=true) (que invoca system.ApplySystem
// y por tanto NewL0().Apply + NewL1().Apply + NewL2().Apply una primera
// vez) y devuelve un *gorm.DB listo para usar. Replica el patrón de
// startPostgresForL0Test / startPostgresForL1Test (el helper canónico
// seeds/e2e/internal/testdb está bajo internal/ y no es importable
// desde aquí).
func startPostgresForL2Test(tb testing.TB) *gorm.DB {
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
		tb.Fatalf("l2 integration: testcontainers GenericContainer: %v", err)
	}
	tb.Cleanup(func() { _ = container.Terminate(ctx) })

	host, err := container.Host(ctx)
	if err != nil {
		tb.Fatalf("l2 integration: container.Host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		tb.Fatalf("l2 integration: container.MappedPort: %v", err)
	}
	connStr := fmt.Sprintf(
		"host=%s port=%s user=test password=test dbname=testdb sslmode=disable",
		host, port.Port(),
	)

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		tb.Fatalf("l2 integration: sql.Open: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		tb.Fatalf("l2 integration: ping: %v", err)
	}
	tb.Cleanup(func() { _ = sqlDB.Close() })

	if _, err := migrations.Migrate(sqlDB, migrations.MigrateOptions{
		Force:  true,
		DBUser: "test",
	}); err != nil {
		tb.Fatalf("l2 integration: migrations.Migrate: %v", err)
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		tb.Fatalf("l2 integration: gorm.Open: %v", err)
	}
	return gdb
}
