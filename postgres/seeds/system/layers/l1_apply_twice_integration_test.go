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

// TestL1_ApplyTwice_Idempotent verifica que aplicar las capas L0 + L1
// dos veces consecutivas produce exactamente el mismo dataset que
// aplicarlas una sola vez: 18 filas (17 de L0 + 1 de L1) sin duplicados.
//
// Justificación: la idempotencia de L1 es contrato (el insert usa
// ON CONFLICT DO NOTHING). Si applyL1Role perdiera el OnConflict, este
// test lo detectaría inmediatamente.
//
// MP-09 F4: L1 dejó de sembrar DATO DE TENANT (escuela demo, usuario
// viewer, user_role, membership). system/ es CONTRATO PURO: L1 sólo
// siembra el rol de contrato announcement_viewer. Se retiró
// assertL1ViewerHasSchool (no hay viewer en el contrato).
//
// Ejecución:
//
//	ENABLE_INTEGRATION_TESTS=true \
//	    go test -tags=integration -run TestL1_ApplyTwice_Idempotent -count=1 \
//	        ./seeds/system/layers/...
//
// Replica el patrón inline de TestL0_ApplyTwice_Idempotent (el helper
// canónico seeds/e2e/internal/testdb no es importable desde aquí).
func TestL1_ApplyTwice_Idempotent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}

	gdb := startPostgresForL1Test(t)

	// migrations.Migrate(Force=true) ya aplica system.ApplySystem,
	// que con L1 registrado invoca NewL0().Apply + NewL1().Apply una
	// primera vez. Por tanto los conteos canónicos ya deben cumplirse.
	assertL0Counts(t, gdb, "L1 test — after initial Migrate (Apply #1)")
	assertL1Counts(t, gdb, "L1 test — after initial Migrate (Apply #1)")

	// Apply #2 — primera reaplicación explícita de L0 + L1.
	if err := layers.NewL0().Apply(gdb); err != nil {
		t.Fatalf("NewL0().Apply (#2): %v", err)
	}
	if err := layers.NewL1().Apply(gdb); err != nil {
		t.Fatalf("NewL1().Apply (#2): %v", err)
	}
	assertL0Counts(t, gdb, "after L0+L1 Apply #2")
	assertL1Counts(t, gdb, "after L0+L1 Apply #2")

	// Apply #3 — segunda reaplicación. Refuerza idempotencia.
	if err := layers.NewL0().Apply(gdb); err != nil {
		t.Fatalf("NewL0().Apply (#3): %v", err)
	}
	if err := layers.NewL1().Apply(gdb); err != nil {
		t.Fatalf("NewL1().Apply (#3): %v", err)
	}
	assertL0Counts(t, gdb, "after L0+L1 Apply #3")
	assertL1Counts(t, gdb, "after L0+L1 Apply #3")
}

// assertL1Counts valida el único conteo canónico de L1: el rol de
// contrato announcement_viewer (1 fila en iam.roles).
//
// MP-09 F4: L1 dejó de sembrar DATO DE TENANT. El permiso efectivo del
// rol (academic.announcements.read) se otorga vía iam.role_grants desde
// L4, así que NO se valida aquí (este test sólo aplica L0+L1).
func assertL1Counts(t *testing.T, gdb *gorm.DB, stage string) {
	t.Helper()

	var got int64
	if err := gdb.Raw(
		`SELECT COUNT(*) FROM iam.roles WHERE name = ?`,
		layers.L1_ROLE_ANNOUNCEMENT_VIEWER_NAME,
	).Scan(&got).Error; err != nil {
		t.Fatalf("[%s] count iam.roles [name=announcement_viewer]: %v", stage, err)
	}

	// Total canónico L1: 1 (sólo el rol de contrato).
	const wantTotal int64 = 1
	if got != wantTotal {
		t.Errorf("[%s] iam.roles [name=announcement_viewer]: got %d, want %d", stage, got, wantTotal)
	}
}

// startPostgresForL1Test levanta un contenedor postgres:15-alpine,
// aplica migrations.Migrate(Force=true) (que invoca system.ApplySystem
// y por tanto NewL0().Apply + NewL1().Apply una primera vez) y
// devuelve un *gorm.DB listo para usar. Replica el patrón de
// startPostgresForL0Test (el helper canónico seeds/e2e/internal/testdb
// está bajo internal/ y no es importable desde aquí).
func startPostgresForL1Test(tb testing.TB) *gorm.DB {
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
		tb.Fatalf("l1 integration: testcontainers GenericContainer: %v", err)
	}
	tb.Cleanup(func() { _ = container.Terminate(ctx) })

	host, err := container.Host(ctx)
	if err != nil {
		tb.Fatalf("l1 integration: container.Host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		tb.Fatalf("l1 integration: container.MappedPort: %v", err)
	}
	connStr := fmt.Sprintf(
		"host=%s port=%s user=test password=test dbname=testdb sslmode=disable",
		host, port.Port(),
	)

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		tb.Fatalf("l1 integration: sql.Open: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		tb.Fatalf("l1 integration: ping: %v", err)
	}
	tb.Cleanup(func() { _ = sqlDB.Close() })

	if _, err := migrations.Migrate(sqlDB, migrations.MigrateOptions{
		Force:  true,
		DBUser: "test",
	}); err != nil {
		tb.Fatalf("l1 integration: migrations.Migrate: %v", err)
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		tb.Fatalf("l1 integration: gorm.Open: %v", err)
	}
	return gdb
}
