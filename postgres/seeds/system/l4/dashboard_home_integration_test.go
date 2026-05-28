//go:build integration
// +build integration

package l4_test

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
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l4"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
)

// TestL4_DashboardHome_Seeded verifica que tras aplicar L4:
//
//  1. Existe una fila con screen_key="dashboard-home" en
//     ui_config.screen_instances.
//  2. NO existe ninguna fila en ui_config.resource_screens que mapee
//     este screen_key — dashboard-home es shell sin mapping (mismo
//     patron que app-login / app-settings).
//  3. La aplicación es idempotente: re-invocar ApplyScreenInstances
//     no produce filas duplicadas (UPSERT por screen_key).
//
// Ejecución:
//
//	ENABLE_INTEGRATION_TESTS=true \
//	    go test -tags=integration -run TestL4_DashboardHome -count=1 \
//	        ./postgres/seeds/system/l4/...
//
// Requiere docker corriendo (testcontainers).
func TestL4_DashboardHome_Seeded(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}

	gdb := startPostgresForL4DashboardHomeTest(t)

	// (1) dashboard-home existe en screen_instances con el UUID esperado.
	var cnt int64
	if err := gdb.Raw(
		`SELECT COUNT(*) FROM ui_config.screen_instances WHERE screen_key = ?`,
		"dashboard-home",
	).Scan(&cnt).Error; err != nil {
		t.Fatalf("count dashboard-home in screen_instances: %v", err)
	}
	if cnt != 1 {
		t.Fatalf("dashboard-home: got %d rows in screen_instances, want 1", cnt)
	}

	var id string
	if err := gdb.Raw(
		`SELECT id::text FROM ui_config.screen_instances WHERE screen_key = ?`,
		"dashboard-home",
	).Scan(&id).Error; err != nil {
		t.Fatalf("load dashboard-home id: %v", err)
	}
	if id != l4.L4_SCREEN_INST_DASHBOARD_HOME_ID {
		t.Errorf("dashboard-home id: got %s, want %s", id, l4.L4_SCREEN_INST_DASHBOARD_HOME_ID)
	}

	// (2) Sin mapping en resource_screens — es shell, igual que
	// app-login / app-settings.
	var rsCount int64
	if err := gdb.Raw(
		`SELECT COUNT(*) FROM ui_config.resource_screens WHERE screen_key = ?`,
		"dashboard-home",
	).Scan(&rsCount).Error; err != nil {
		t.Fatalf("count dashboard-home in resource_screens: %v", err)
	}
	if rsCount != 0 {
		t.Errorf(
			"dashboard-home: got %d rows in resource_screens, want 0 (shell sin mapping)",
			rsCount,
		)
	}

	// (3) Idempotencia: re-aplicar ApplyScreenInstances no duplica.
	if err := l4.ApplyScreenInstances(gdb); err != nil {
		t.Fatalf("re-apply ApplyScreenInstances: %v", err)
	}
	if err := gdb.Raw(
		`SELECT COUNT(*) FROM ui_config.screen_instances WHERE screen_key = ?`,
		"dashboard-home",
	).Scan(&cnt).Error; err != nil {
		t.Fatalf("count dashboard-home after re-apply: %v", err)
	}
	if cnt != 1 {
		t.Errorf("dashboard-home no idempotente: got %d rows after re-apply, want 1", cnt)
	}
}

// startPostgresForL4DashboardHomeTest levanta postgres:15-alpine y
// ejecuta migrations.Migrate(Force=true, SeedUpToLayer=L4_LAYER_NAME).
// Replica el patrón de startPostgresForL4AliasTest.
func startPostgresForL4DashboardHomeTest(tb testing.TB) *gorm.DB {
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
		tb.Fatalf("l4 dashboard-home integration: testcontainers GenericContainer: %v", err)
	}
	tb.Cleanup(func() { _ = container.Terminate(ctx) })

	host, err := container.Host(ctx)
	if err != nil {
		tb.Fatalf("l4 dashboard-home integration: container.Host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		tb.Fatalf("l4 dashboard-home integration: container.MappedPort: %v", err)
	}
	connStr := fmt.Sprintf(
		"host=%s port=%s user=test password=test dbname=testdb sslmode=disable",
		host, port.Port(),
	)

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		tb.Fatalf("l4 dashboard-home integration: sql.Open: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		tb.Fatalf("l4 dashboard-home integration: ping: %v", err)
	}
	tb.Cleanup(func() { _ = sqlDB.Close() })

	if _, err := migrations.Migrate(sqlDB, migrations.MigrateOptions{
		Force:         true,
		DBUser:        "test",
		SeedUpToLayer: layers.L4_LAYER_NAME,
	}); err != nil {
		tb.Fatalf("l4 dashboard-home integration: migrations.Migrate: %v", err)
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		tb.Fatalf("l4 dashboard-home integration: gorm.Open: %v", err)
	}
	return gdb
}
