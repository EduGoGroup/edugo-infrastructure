//go:build integration
// +build integration

package playground_v2_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// TestApplyOnTop_Scenarios valida el contrato on-top (plan 053 T5.2) para los 6
// escenarios, sobre una BD migrada con system + base:
//  1. con usuarios ausentes falla listando TODOS los emails que faltan y no
//     escribe nada;
//  2. con los usuarios presentes (creados como lo haría la API: UUID nuevo)
//     aplica, usa esos ids y encola una fila de outbox por usuario;
//  3. es re-ejecutable: los conteos no cambian y el outbox sube de generación.
func TestApplyOnTop_Scenarios(t *testing.T) {
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}
	gdb := startPostgres(t)

	baseUsers := count(t, gdb, "auth.users")

	for _, name := range playground_v2.OnTopScenarios() {
		t.Run(name, func(t *testing.T) {
			users, err := playground_v2.ScenarioUsers(name)
			if err != nil {
				t.Fatalf("ScenarioUsers: %v", err)
			}

			// (1) Solo la primera mitad presente: falla sin escribir.
			half := users[:len(users)/2]
			createUsers(t, gdb, half)
			before := snapshot(t, gdb)
			err = playground_v2.ApplyOnTop(gdb, name)
			if err == nil {
				t.Fatalf("ApplyOnTop con usuarios ausentes debía fallar")
			}
			for _, u := range users[len(users)/2:] {
				if !strings.Contains(err.Error(), u.Email) {
					t.Errorf("el error no lista el email ausente %s: %v", u.Email, err)
				}
			}
			for _, u := range half {
				if strings.Contains(err.Error(), u.Email) {
					t.Errorf("el error lista un email presente %s: %v", u.Email, err)
				}
			}
			if after := snapshot(t, gdb); after != before {
				t.Fatalf("ApplyOnTop fallido escribió: antes %+v, después %+v", before, after)
			}

			// (2) Todos presentes: aplica y encola.
			createUsers(t, gdb, users[len(users)/2:])
			if err := playground_v2.ApplyOnTop(gdb, name); err != nil {
				t.Fatalf("ApplyOnTop: %v", err)
			}
			first := snapshot(t, gdb)
			if first.memberships <= before.memberships {
				t.Errorf("ApplyOnTop no sembró memberships: %+v", first)
			}
			ids := userIDs(t, gdb, users)
			if got := countWhere(t, gdb, "auth.identity_sync_outbox", "user_id IN ?", ids); got != int64(len(users)) {
				t.Errorf("outbox: %d filas para los usuarios del escenario, esperado %d", got, len(users))
			}
			// Los user_roles del escenario cuelgan de los ids resueltos.
			if got := countWhere(t, gdb, "iam.user_roles", "user_id IN ?", ids); got == 0 {
				t.Errorf("ningún iam.user_roles apunta a los usuarios resueltos")
			}

			// (3) Re-ejecutable.
			if err := playground_v2.ApplyOnTop(gdb, name); err != nil {
				t.Fatalf("ApplyOnTop (2ª vez): %v", err)
			}
			second := snapshot(t, gdb)
			if second != first {
				t.Fatalf("la 2ª aplicación cambió conteos: %+v → %+v", first, second)
			}
			if got := countWhere(t, gdb, "auth.identity_sync_outbox", "user_id IN ? AND generation = 2", ids); got != int64(len(users)) {
				t.Errorf("outbox tras re-aplicar: %d filas en generación 2, esperado %d", got, len(users))
			}
		})
	}

	// Ningún escenario on-top crea ni borra usuarios: solo están los de base y
	// los creados por el test.
	created := int64(0)
	for _, name := range playground_v2.OnTopScenarios() {
		users, _ := playground_v2.ScenarioUsers(name)
		created += int64(len(users))
	}
	if got := count(t, gdb, "auth.users"); got != baseUsers+created {
		t.Errorf("auth.users = %d, esperado %d (base) + %d (test)", got, baseUsers, created)
	}

	// Apply(name) de un on-top delega en ApplyOnTop (aquí ya hay usuarios).
	if err := playground_v2.Apply(gdb, "onboarding"); err != nil {
		t.Errorf("Apply(onboarding) con usuarios presentes: %v", err)
	}
}

type counts struct {
	schools, memberships, userRoles, joinRequests, outbox, users int64
}

func snapshot(t *testing.T, gdb *gorm.DB) counts {
	t.Helper()
	return counts{
		schools:      count(t, gdb, "academic.schools"),
		memberships:  count(t, gdb, "academic.memberships"),
		userRoles:    count(t, gdb, "iam.user_roles"),
		joinRequests: count(t, gdb, "academic.school_join_requests"),
		outbox:       count(t, gdb, "auth.identity_sync_outbox"),
		users:        count(t, gdb, "auth.users"),
	}
}

func count(t *testing.T, gdb *gorm.DB, table string) int64 {
	t.Helper()
	var n int64
	if err := gdb.Table(table).Count(&n).Error; err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	return n
}

func countWhere(t *testing.T, gdb *gorm.DB, table, where string, args ...any) int64 {
	t.Helper()
	var n int64
	if err := gdb.Table(table).Where(where, args...).Count(&n).Error; err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	return n
}

// createUsers inserta los usuarios como lo haría la API de identidad: UUID
// nuevo (no el literal que llevaba el seed) y un hash que no sirve para entrar.
func createUsers(t *testing.T, gdb *gorm.DB, users []playground_v2.ScenarioUser) {
	t.Helper()
	for _, u := range users {
		row := entities.User{
			ID:           uuid.New(),
			Email:        u.Email,
			PasswordHash: "!unusable",
			FirstName:    u.FirstName,
			LastName:     u.LastName,
			IsActive:     true,
		}
		if err := gdb.Create(&row).Error; err != nil {
			t.Fatalf("creando usuario %s: %v", u.Email, err)
		}
	}
}

func userIDs(t *testing.T, gdb *gorm.DB, users []playground_v2.ScenarioUser) []string {
	t.Helper()
	emails := make([]string, 0, len(users))
	for _, u := range users {
		emails = append(emails, u.Email)
	}
	var ids []string
	if err := gdb.Raw(`SELECT id::text FROM auth.users WHERE email IN ?`, emails).Scan(&ids).Error; err != nil {
		t.Fatalf("leyendo ids: %v", err)
	}
	if len(ids) != len(users) {
		t.Fatalf("ids resueltos %d, esperado %d", len(ids), len(users))
	}
	return ids
}

// startPostgres arranca un Postgres efímero, migra con system + base y
// devuelve un *gorm.DB.
func startPostgres(tb testing.TB) *gorm.DB {
	tb.Helper()
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "test",
				"POSTGRES_PASSWORD": "test",
				"POSTGRES_DB":       "testdb",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		},
		Started: true,
	})
	if err != nil {
		tb.Fatalf("playground_v2 integration: GenericContainer: %v", err)
	}
	tb.Cleanup(func() { _ = container.Terminate(ctx) })

	host, err := container.Host(ctx)
	if err != nil {
		tb.Fatalf("playground_v2 integration: Host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		tb.Fatalf("playground_v2 integration: MappedPort: %v", err)
	}
	sqlDB, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=test password=test dbname=testdb sslmode=disable", host, port.Port()))
	if err != nil {
		tb.Fatalf("playground_v2 integration: sql.Open: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		tb.Fatalf("playground_v2 integration: ping: %v", err)
	}
	tb.Cleanup(func() { _ = sqlDB.Close() })

	if _, err := migrations.Migrate(sqlDB, migrations.MigrateOptions{
		Force:        true,
		DBUser:       "test",
		PlaygroundV2: "base",
	}); err != nil {
		tb.Fatalf("playground_v2 integration: Migrate: %v", err)
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		tb.Fatalf("playground_v2 integration: gorm.Open: %v", err)
	}
	return gdb
}
