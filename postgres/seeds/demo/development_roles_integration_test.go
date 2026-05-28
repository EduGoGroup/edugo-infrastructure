//go:build integration
// +build integration

package demo_test

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

// TestDemo_UserRoles_AfterMigrate verifica que tras Migrate(SeedDemo=true)
// los 21 usuarios demo tienen al menos un user_role asignado y que los
// role_id referencian roles válidos (FK virtual contra iam.roles).
//
// Justificación: post-Fase-6 se eliminó el catálogo legacy de roles
// `10000000-...-0001..0012`. seedUserRoles asignaba esos UUIDs y por
// tanto generaba filas con FK rota. Este test es el guardrail de que
// el refactor a L0/L4 quedó consistente.
//
// Ejecución:
//
//	ENABLE_INTEGRATION_TESTS=true \
//	    go test -tags=integration -run TestDemo_UserRoles_AfterMigrate -count=1 \
//	        ./seeds/demo/...
func TestDemo_UserRoles_AfterMigrate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}

	gdb := startPostgresForDemoTest(t)

	// 1) Los 21 usuarios demo deben existir.
	var userCount int64
	if err := gdb.Raw(`SELECT COUNT(*) FROM auth.users WHERE email LIKE '%@edugo.test'`).Scan(&userCount).Error; err != nil {
		t.Fatalf("count demo users: %v", err)
	}
	if userCount != 21 {
		t.Errorf("demo users: got %d, want 21", userCount)
	}

	// 2) Cada uno de los 21 usuarios debe tener al menos un user_role.
	type emailRow struct {
		Email string
		Cnt   int64
	}
	var rows []emailRow
	if err := gdb.Raw(`
		SELECT u.email AS email, COUNT(ur.id) AS cnt
		FROM auth.users u
		LEFT JOIN iam.user_roles ur ON ur.user_id = u.id AND ur.is_active = true
		WHERE u.email LIKE '%@edugo.test'
		GROUP BY u.email
		ORDER BY u.email
	`).Scan(&rows).Error; err != nil {
		t.Fatalf("user_roles per user: %v", err)
	}
	for _, r := range rows {
		if r.Cnt < 1 {
			t.Errorf("user %s has no active user_role", r.Email)
		}
	}

	// 3) Mapeo email → role_name esperado (al menos los principales).
	expectedRoleByEmail := map[string]string{
		"super@edugo.test":            "super_admin",
		"admin.sanignacio@edugo.test": "school_admin",
		"admin.crearte@edugo.test":    "school_admin",
		"coord.academico@edugo.test":  "school_coordinator",
		"prof.martinez@edugo.test":    "teacher",
		"prof.gonzalez@edugo.test":    "teacher",
		"facilitador.ruiz@edugo.test": "teacher",
		"est.carlos@edugo.test":       "student",
		"est.sofia@edugo.test":        "student",
		"est.diego@edugo.test":        "student",
		"est.valentina@edugo.test":    "student",
		"est.mateo@edugo.test":        "student",
		"tutor.mendoza@edugo.test":    "guardian",
		"tutora.herrera@edugo.test":   "guardian",
		// PRE-4: el rol `platform_admin` fue eliminado; este usuario
		// ahora hereda `super_admin` (L0), que cubre el mismo acceso global.
		"admin.plataforma@edugo.test":    "super_admin",
		"director.sanignacio@edugo.test": "school_director",
		"asist.admin@edugo.test":         "school_assistant",
		"asist.prof@edugo.test":          "assistant_teacher",
		"observador@edugo.test":          "observer",
		"guardian.pendiente@edugo.test":  "guardian",
		"readonly@edugo.test":            "readonly_auditor",
	}

	for email, wantRole := range expectedRoleByEmail {
		var got string
		err := gdb.Raw(`
			SELECT r.name
			FROM auth.users u
			JOIN iam.user_roles ur ON ur.user_id = u.id AND ur.is_active = true
			JOIN iam.roles r ON r.id = ur.role_id
			WHERE u.email = ?
			ORDER BY r.name
			LIMIT 1
		`, email).Scan(&got).Error
		if err != nil {
			t.Errorf("lookup role for %s: %v", email, err)
			continue
		}
		if got != wantRole {
			t.Errorf("user %s: got role %q, want %q", email, got, wantRole)
		}
	}

	// 4) FK virtual: NO debe haber user_role con role_id huérfano.
	var orphan int64
	if err := gdb.Raw(`
		SELECT COUNT(*)
		FROM iam.user_roles ur
		WHERE ur.role_id NOT IN (SELECT id FROM iam.roles)
	`).Scan(&orphan).Error; err != nil {
		t.Fatalf("orphan user_roles query: %v", err)
	}
	if orphan != 0 {
		t.Errorf("found %d user_roles with role_id NOT IN iam.roles (FK virtual rota)", orphan)
	}

	// 5) Sanity: los role_id usados deben ser los constants L0/L4.
	expectedRoleIDs := []string{
		layers.L0_ROLE_SUPER_ADMIN_ID,
		l4.L4_ROLE_STUDENT_ID,
		l4.L4_ROLE_TEACHER_ID,
		l4.L4_ROLE_GUARDIAN_ID,
		// PRE-4: l4.L4_ROLE_ADMIN_ID removido (platform_admin eliminado).
		l4.L4_ROLE_SCHOOL_ADMIN_ID,
		l4.L4_ROLE_SCHOOL_DIRECTOR_ID,
		l4.L4_ROLE_SCHOOL_COORDINATOR_ID,
		l4.L4_ROLE_SCHOOL_ASSISTANT_ID,
		l4.L4_ROLE_ASSISTANT_TEACHER_ID,
		l4.L4_ROLE_OBSERVER_ID,
		l4.L4_ROLE_READONLY_AUDITOR_ID,
	}
	for _, rid := range expectedRoleIDs {
		var exists int64
		if err := gdb.Raw(`SELECT COUNT(*) FROM iam.roles WHERE id = ?`, rid).Scan(&exists).Error; err != nil {
			t.Fatalf("verify role %s exists: %v", rid, err)
		}
		if exists != 1 {
			t.Errorf("role id %s not found in iam.roles (need 1 row, got %d)", rid, exists)
		}
	}

	// 6) Conteo total user_roles del demo (los 27 cc000000-...).
	var demoURCount int64
	if err := gdb.Raw(`
		SELECT COUNT(*) FROM iam.user_roles
		WHERE id::text LIKE 'cc000000-%'
	`).Scan(&demoURCount).Error; err != nil {
		t.Fatalf("count demo user_roles: %v", err)
	}
	if demoURCount != 27 {
		t.Errorf("demo user_roles (cc000000-*): got %d, want 27", demoURCount)
	}
}

// startPostgresForDemoTest levanta un contenedor postgres:15-alpine,
// ejecuta migrations.Migrate con SeedDemo=true y SeedUpToLayer=""
// (toda la pila system L0..L4 + demo) y devuelve un *gorm.DB.
func startPostgresForDemoTest(tb testing.TB) *gorm.DB {
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
		tb.Fatalf("demo integration: testcontainers GenericContainer: %v", err)
	}
	tb.Cleanup(func() { _ = container.Terminate(ctx) })

	host, err := container.Host(ctx)
	if err != nil {
		tb.Fatalf("demo integration: container.Host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		tb.Fatalf("demo integration: container.MappedPort: %v", err)
	}
	connStr := fmt.Sprintf(
		"host=%s port=%s user=test password=test dbname=testdb sslmode=disable",
		host, port.Port(),
	)

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		tb.Fatalf("demo integration: sql.Open: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		tb.Fatalf("demo integration: ping: %v", err)
	}
	tb.Cleanup(func() { _ = sqlDB.Close() })

	if _, err := migrations.Migrate(sqlDB, migrations.MigrateOptions{
		Force:         true,
		SeedDemo:      true,
		SeedUpToLayer: "",
		DBUser:        "test",
	}); err != nil {
		tb.Fatalf("demo integration: migrations.Migrate: %v", err)
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		tb.Fatalf("demo integration: gorm.Open: %v", err)
	}
	return gdb
}
