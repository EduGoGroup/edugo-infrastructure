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
// aplicarlas una sola vez: 23 filas (17 de L0 + 6 de L1) sin duplicados.
//
// Justificación: la idempotencia de L1 es contrato (todos los inserts
// usan ON CONFLICT DO NOTHING). Si una de las cinco funciones
// applyL1_* perdiera el OnConflict, este test lo detectaría
// inmediatamente.
//
// Adicionalmente verifica F3-REQ-6.2: el user_role del viewer debe
// tener school_id NOT NULL (requisito del contrato scope=school).
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
	assertL1ViewerHasSchool(t, gdb, "L1 test — after initial Migrate (Apply #1)")

	// Apply #2 — primera reaplicación explícita de L0 + L1.
	if err := layers.NewL0().Apply(gdb); err != nil {
		t.Fatalf("NewL0().Apply (#2): %v", err)
	}
	if err := layers.NewL1().Apply(gdb); err != nil {
		t.Fatalf("NewL1().Apply (#2): %v", err)
	}
	assertL0Counts(t, gdb, "after L0+L1 Apply #2")
	assertL1Counts(t, gdb, "after L0+L1 Apply #2")
	assertL1ViewerHasSchool(t, gdb, "after L0+L1 Apply #2")

	// Apply #3 — segunda reaplicación. Refuerza idempotencia.
	if err := layers.NewL0().Apply(gdb); err != nil {
		t.Fatalf("NewL0().Apply (#3): %v", err)
	}
	if err := layers.NewL1().Apply(gdb); err != nil {
		t.Fatalf("NewL1().Apply (#3): %v", err)
	}
	assertL0Counts(t, gdb, "after L0+L1 Apply #3")
	assertL1Counts(t, gdb, "after L0+L1 Apply #3")
	assertL1ViewerHasSchool(t, gdb, "after L0+L1 Apply #3")
}

// assertL1Counts valida los seis conteos canónicos de L1 (6 filas
// distribuidas en 6 tablas) filtrados por los IDs específicos de L1
// para no contaminar con datos de L0.
func assertL1Counts(t *testing.T, gdb *gorm.DB, stage string) {
	t.Helper()

	type check struct {
		desc  string
		query string
		args  []any
		want  int64
	}

	checks := []check{
		{
			desc:  "academic.schools [code=L1-DEMO]",
			query: `SELECT COUNT(*) FROM academic.schools WHERE code = ?`,
			args:  []any{layers.L1_SCHOOL_DEMO_CODE},
			want:  1,
		},
		{
			desc:  "iam.roles [name=announcement_viewer]",
			query: `SELECT COUNT(*) FROM iam.roles WHERE name = ?`,
			args:  []any{layers.L1_ROLE_ANNOUNCEMENT_VIEWER_NAME},
			want:  1,
		},
		{
			desc:  "auth.users [email=viewer@edugo.demo]",
			query: `SELECT COUNT(*) FROM auth.users WHERE email = ?`,
			args:  []any{layers.L1_VIEWER_EMAIL},
			want:  1,
		},
		{
			desc:  "iam.role_permissions [role_id=L1 viewer]",
			query: `SELECT COUNT(*) FROM iam.role_permissions WHERE role_id = ?`,
			args:  []any{layers.L1_ROLE_ANNOUNCEMENT_VIEWER_ID},
			want:  1,
		},
		{
			desc:  "iam.user_roles [user_id=L1 viewer AND role_id=L1 viewer]",
			query: `SELECT COUNT(*) FROM iam.user_roles WHERE user_id = ? AND role_id = ?`,
			args:  []any{layers.L1_USER_VIEWER_ID, layers.L1_ROLE_ANNOUNCEMENT_VIEWER_ID},
			want:  1,
		},
		{
			desc:  "academic.memberships [user_id=L1 viewer AND school_id=L1 school]",
			query: `SELECT COUNT(*) FROM academic.memberships WHERE user_id = ? AND school_id = ? AND is_active = true`,
			args:  []any{layers.L1_USER_VIEWER_ID, layers.L1_SCHOOL_DEMO_ID},
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

	// Total canónico L1: 1 + 1 + 1 + 1 + 1 + 1 = 6.
	const wantTotal int64 = 6
	if total != wantTotal {
		t.Errorf("[%s] total L1 rows: got %d, want %d", stage, total, wantTotal)
	}
}

// assertL1ViewerHasSchool valida F3-REQ-6.2: el user_role del viewer
// debe tener school_id NOT NULL (post_gorm.sql:~311 — contrato
// scope=school exige school_id presente).
func assertL1ViewerHasSchool(t *testing.T, gdb *gorm.DB, stage string) {
	t.Helper()
	var nullCount int64
	if err := gdb.Raw(
		`SELECT COUNT(*) FROM iam.user_roles WHERE user_id = ? AND role_id = ? AND school_id IS NULL`,
		layers.L1_USER_VIEWER_ID, layers.L1_ROLE_ANNOUNCEMENT_VIEWER_ID,
	).Scan(&nullCount).Error; err != nil {
		t.Fatalf("[%s] check viewer school_id: %v", stage, err)
	}
	if nullCount != 0 {
		t.Errorf("[%s] F3-REQ-6.2: viewer user_role has NULL school_id (%d rows); expected school_id NOT NULL",
			stage, nullCount)
	}

	var matchCount int64
	if err := gdb.Raw(
		`SELECT COUNT(*) FROM iam.user_roles WHERE user_id = ? AND role_id = ? AND school_id = ?`,
		layers.L1_USER_VIEWER_ID, layers.L1_ROLE_ANNOUNCEMENT_VIEWER_ID, layers.L1_SCHOOL_DEMO_ID,
	).Scan(&matchCount).Error; err != nil {
		t.Fatalf("[%s] check viewer school_id match: %v", stage, err)
	}
	if matchCount != 1 {
		t.Errorf("[%s] F3-REQ-6.2: viewer user_role school_id != L1_SCHOOL_DEMO_ID (got %d matching rows; want 1)",
			stage, matchCount)
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
