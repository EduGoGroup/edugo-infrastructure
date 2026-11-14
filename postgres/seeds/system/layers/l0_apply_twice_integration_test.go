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

// TestL0_ApplyTwice_Idempotent verifica que aplicar la capa L0 dos
// veces consecutivas produce exactamente el mismo dataset que aplicarla
// una sola vez: 13 filas en 8 tablas, sin duplicados ni errores.
//
// Justificación: la idempotencia de L0 es contrato (todos los inserts
// usan ON CONFLICT). Si una de las cuatro funciones applyL0_* perdiera
// el OnConflict, este test lo detectaría inmediatamente.
//
// Ejecución:
//
//	ENABLE_INTEGRATION_TESTS=true \
//	    go test -tags=integration -run TestL0_ApplyTwice_Idempotent -count=1 \
//	        ./seeds/system/layers/...
//
// El gate `ENABLE_INTEGRATION_TESTS=true` evita que el test corra en
// CI vanilla. El build tag `integration` añade un segundo cierre.
//
// El helper canónico (seeds/e2e/internal/testdb) no es importable
// desde este paquete (está bajo internal/), por lo que replicamos
// inline el patrón documentado allí: contenedor postgres:15-alpine +
// migrations.Migrate(Force=true) sobre BD efímera.
func TestL0_ApplyTwice_Idempotent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}

	gdb := startPostgresForL0Test(t)

	// migrations.Migrate(Force=true) ya aplica system.ApplySystem,
	// que invoca NewL0().Apply una vez. Por lo tanto los conteos
	// esperados ya deben cumplirse antes de re-aplicar.
	assertL0Counts(t, gdb, "after initial Migrate (Apply #1)")

	// Apply #2 — primera reaplicación explícita.
	if err := layers.NewL0().Apply(gdb); err != nil {
		t.Fatalf("NewL0().Apply (#2): %v", err)
	}
	assertL0Counts(t, gdb, "after Apply #2")

	// Apply #3 — segunda reaplicación. El spec pide "aplicar dos
	// veces"; este tercer ciclo refuerza el contrato de idempotencia
	// sin coste adicional relevante.
	if err := layers.NewL0().Apply(gdb); err != nil {
		t.Fatalf("NewL0().Apply (#3): %v", err)
	}
	assertL0Counts(t, gdb, "after Apply #3")
}

// assertL0Counts valida los ocho conteos canónicos de L0 (13 filas
// distribuidas en 8 tablas) filtrados por los IDs/keys específicos
// de L0 para no contaminar con dataset legado u otras capas.
//
// P4-1 (plan B, SchemaVersion 3.21.0): iam.role_permissions ya no
// existe y L0 dejó de sembrar los 4 role_permissions del super_admin
// (ver l0_minimal.go / l0_roles.go). El permiso efectivo del rol se
// otorga vía iam.role_grants (pattern `*`) desde L4, y se valida aparte
// en assertSuperAdminWildcardGrant porque no es una fila de L0.
func assertL0Counts(t *testing.T, gdb *gorm.DB, stage string) {
	t.Helper()

	type check struct {
		desc  string
		query string
		args  []any
		want  int64
	}

	checks := []check{
		{
			desc:  "iam.resources [key=announcements]",
			query: `SELECT COUNT(*) FROM iam.resources WHERE key = ?`,
			args:  []any{layers.L0_RESOURCE_ANNOUNCEMENTS_KEY},
			want:  1,
		},
		{
			desc:  "iam.roles [name=super_admin]",
			query: `SELECT COUNT(*) FROM iam.roles WHERE name = ?`,
			args:  []any{layers.L0_ROLE_SUPER_ADMIN_NAME},
			want:  1,
		},
		{
			// Los 4 permisos CRUD de L0 se llaman
			// academic.announcements.{read,create,update,delete}
			// (formato path; el formato legacy era `announcements:<accion>`).
			// Se filtra por los 4 IDs de L0 para que el check siga siendo
			// invariante L0 aunque otra capa añada permisos al recurso.
			desc: "iam.permissions [id IN L0 announcements CRUD, name LIKE academic.announcements.%]",
			query: `SELECT COUNT(*) FROM iam.permissions ` +
				`WHERE id IN (?::uuid, ?::uuid, ?::uuid, ?::uuid) AND name LIKE ?`,
			args: []any{
				layers.L0_PERM_ANNOUNCEMENTS_READ,
				layers.L0_PERM_ANNOUNCEMENTS_CREATE,
				layers.L0_PERM_ANNOUNCEMENTS_UPDATE,
				layers.L0_PERM_ANNOUNCEMENTS_DELETE,
				"academic.announcements.%",
			},
			want: 4,
		},
		{
			desc:  "ui_config.screen_templates [name in L0 trio]",
			query: `SELECT COUNT(*) FROM ui_config.screen_templates WHERE name IN (?, ?, ?)`,
			args:  []any{"list-basic-v1", "detail-basic-v1", "form-basic-v1"},
			want:  3,
		},
		{
			desc:  "ui_config.screen_instances [screen_key=announcements-list]",
			query: `SELECT COUNT(*) FROM ui_config.screen_instances WHERE screen_key = ?`,
			args:  []any{layers.L0_SCREEN_KEY_ANNOUNCEMENTS_LIST},
			want:  1,
		},
		{
			// Filtrado por screen_type='list' para que el check siga
			// siendo invariante L0 incluso cuando capas superiores
			// (L2 form, futuras detail/etc.) agreguen más mappings
			// para el mismo resource_id.
			desc:  "ui_config.resource_screens [resource_id=L0 announcements, type=list]",
			query: `SELECT COUNT(*) FROM ui_config.resource_screens WHERE resource_id = ? AND screen_type = ?`,
			args:  []any{layers.L0_RESOURCE_ANNOUNCEMENTS_ID, "list"},
			want:  1,
		},
		{
			desc:  "auth.users [email=super_admin@edugo.system]",
			query: `SELECT COUNT(*) FROM auth.users WHERE email = ?`,
			args:  []any{layers.L0_SUPER_ADMIN_EMAIL},
			want:  1,
		},
		{
			desc:  "iam.user_roles [user_id=L0 super_admin AND role_id=L0 super_admin]",
			query: `SELECT COUNT(*) FROM iam.user_roles WHERE user_id = ? AND role_id = ?`,
			args:  []any{layers.L0_USER_SUPER_ADMIN_ID, layers.L0_ROLE_SUPER_ADMIN_ID},
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

	// Total canónico L0: 1 + 1 + 4 + 3 + 1 + 1 + 1 + 1 = 13.
	const wantTotal int64 = 13
	if total != wantTotal {
		t.Errorf("[%s] total L0 rows: got %d, want %d", stage, total, wantTotal)
	}

	assertSuperAdminWildcardGrant(t, gdb, stage)
}

// assertSuperAdminWildcardGrant valida que el super_admin de L0 conserva
// su permiso efectivo: exactamente 1 grant `*` con effect=allow en
// iam.role_grants. Sustituye al conteo de los 4 iam.role_permissions
// super_admin × announcements, tabla eliminada en P4-1. El grant lo
// siembra L4 (applyL4RoleGrants), no L0: por eso no entra en el total
// de filas de L0, y reaplicar L0 no debe alterarlo.
func assertSuperAdminWildcardGrant(t *testing.T, gdb *gorm.DB, stage string) {
	t.Helper()

	var got int64
	if err := gdb.Raw(
		`SELECT COUNT(*) FROM iam.role_grants WHERE role_id = ?::uuid AND pattern = ? AND effect = ?`,
		layers.L0_ROLE_SUPER_ADMIN_ID, "*", "allow",
	).Scan(&got).Error; err != nil {
		t.Fatalf("[%s] count iam.role_grants [role_id=L0 super_admin, pattern=*]: %v", stage, err)
	}
	if got != 1 {
		t.Errorf("[%s] iam.role_grants [role_id=L0 super_admin, pattern=*, effect=allow]: got %d, want 1", stage, got)
	}
}

// startPostgresForL0Test levanta un contenedor postgres:15-alpine,
// aplica migrations.Migrate(Force=true) (que incluye system.ApplySystem
// y por tanto NewL0().Apply una primera vez) y devuelve un *gorm.DB
// listo para usar. Replica el patrón de seeds/e2e/internal/testdb
// porque ese paquete está bajo internal/ y no es importable desde aquí.
func startPostgresForL0Test(tb testing.TB) *gorm.DB {
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
		tb.Fatalf("l0 integration: testcontainers GenericContainer: %v", err)
	}
	tb.Cleanup(func() { _ = container.Terminate(ctx) })

	host, err := container.Host(ctx)
	if err != nil {
		tb.Fatalf("l0 integration: container.Host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		tb.Fatalf("l0 integration: container.MappedPort: %v", err)
	}
	connStr := fmt.Sprintf(
		"host=%s port=%s user=test password=test dbname=testdb sslmode=disable",
		host, port.Port(),
	)

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		tb.Fatalf("l0 integration: sql.Open: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		tb.Fatalf("l0 integration: ping: %v", err)
	}
	tb.Cleanup(func() { _ = sqlDB.Close() })

	if _, err := migrations.Migrate(sqlDB, migrations.MigrateOptions{
		Force:  true,
		DBUser: "test",
	}); err != nil {
		tb.Fatalf("l0 integration: migrations.Migrate: %v", err)
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		tb.Fatalf("l0 integration: gorm.Open: %v", err)
	}
	return gdb
}
