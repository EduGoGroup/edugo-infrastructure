// Package testdb provee un helper compartido para los tests integration
// del seed E2E (paquete `seeds/e2e`). Encapsula el setup de Postgres
// vía testcontainers-go con dos modos:
//
//   - Local (default): arranca un contenedor postgres:15-alpine,
//     aplica migrations completas y el production seed. Útil para CI
//     local y nightly contra Docker.
//   - Cloud (override): si la variable POSTGRES_URI está definida en el
//     environment, abre conexión directa a esa URI (Neon u otro) sin
//     levantar contenedor. Útil para los benchmarks que miden la cota
//     real contra cloud.
//
// El cleanup del contenedor y de la conexión se registra con
// tb.Cleanup() — no es necesario llamarlo manualmente desde el test.
package testdb

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
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system"
)

// StartPostgres prepara una BD lista para los tests integration del
// seed E2E. Si POSTGRES_URI está definido, abre conexión directa contra
// esa URI y asume que ya tiene migrations + production seed aplicados
// (modo cloud / Neon). En caso contrario, arranca un contenedor
// postgres:15-alpine, aplica migrations.Migrate(Force=true) y
// system.ApplySystem sobre la BD efímera.
//
// Devuelve el *gorm.DB listo para usar. La conexión y el contenedor
// se cierran automáticamente con tb.Cleanup().
func StartPostgres(tb testing.TB) *gorm.DB {
	tb.Helper()
	if uri := os.Getenv("POSTGRES_URI"); uri != "" {
		return openExisting(tb, uri)
	}
	return startContainer(tb, "")
}

// StartPostgresUpTo es igual que StartPostgres pero limita el system
// seed a las capas hasta `upToLayer` inclusive (vacío = todas). Útil
// para scenarios cuyo assert depende de aislamiento de capa (p.ej.
// l3_isolation, que validan ausencia de filas introducidas por L4).
//
// En modo cloud (POSTGRES_URI definido) `upToLayer` se ignora — la BD
// remota se asume ya seedeada.
func StartPostgresUpTo(tb testing.TB, upToLayer string) *gorm.DB {
	tb.Helper()
	if uri := os.Getenv("POSTGRES_URI"); uri != "" {
		return openExisting(tb, uri)
	}
	return startContainer(tb, upToLayer)
}

// openExisting abre una conexión directa contra una URI existente
// (modo cloud). No aplica migrations ni seeds: asume que la BD ya está
// preparada.
func openExisting(tb testing.TB, uri string) *gorm.DB {
	tb.Helper()
	sqlDB, err := sql.Open("postgres", uri)
	if err != nil {
		tb.Fatalf("testdb: sql.Open POSTGRES_URI: %v", err)
	}
	if err := sqlDB.PingContext(context.Background()); err != nil {
		tb.Fatalf("testdb: ping POSTGRES_URI: %v", err)
	}
	tb.Cleanup(func() { _ = sqlDB.Close() })
	return openGORM(tb, sqlDB)
}

// startContainer levanta un contenedor postgres:15-alpine, aplica las
// migraciones completas (Force=true) y el production seed mínimo que
// las fixtures E2E necesitan (resources, roles, permisos, ui_config).
// Mantiene paridad con el patrón canónico de
// migrations/migrations_integration_test.go.
func startContainer(tb testing.TB, upToLayer string) *gorm.DB {
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
		tb.Fatalf("testdb: testcontainers GenericContainer: %v", err)
	}
	tb.Cleanup(func() { _ = container.Terminate(ctx) })

	host, err := container.Host(ctx)
	if err != nil {
		tb.Fatalf("testdb: container.Host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		tb.Fatalf("testdb: container.MappedPort: %v", err)
	}
	connStr := fmt.Sprintf(
		"host=%s port=%s user=test password=test dbname=testdb sslmode=disable",
		host, port.Port(),
	)

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		tb.Fatalf("testdb: sql.Open container: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		tb.Fatalf("testdb: ping container: %v", err)
	}
	tb.Cleanup(func() { _ = sqlDB.Close() })

	if _, err := migrations.Migrate(sqlDB, migrations.MigrateOptions{
		Force:         true,
		DBUser:        "test",
		SeedUpToLayer: upToLayer,
	}); err != nil {
		tb.Fatalf("testdb: migrations.Migrate: %v", err)
	}
	// migrations.Migrate ya aplica el system seed cuando opts.Force
	// dispara el flujo completo; lo invocamos de nuevo para cubrir el
	// caso (poco probable en este helper) de un no-op donde la BD haya
	// quedado con tablas pero sin seeds. Es idempotente vía OnConflict.
	if err := system.ApplySystem(sqlDB, upToLayer); err != nil {
		tb.Fatalf("testdb: system.ApplySystem: %v", err)
	}
	return openGORM(tb, sqlDB)
}

// openGORM envuelve un *sql.DB en gorm.DB con logger silencioso.
func openGORM(tb testing.TB, sqlDB *sql.DB) *gorm.DB {
	tb.Helper()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		tb.Fatalf("testdb: gorm.Open: %v", err)
	}
	return gdb
}

// IntegrationGate decide si los tests integration deben correr. Devuelve
// true si ENABLE_INTEGRATION_TESTS=true. El caller suele invocarla así:
//
//	if !testdb.IntegrationGate() {
//	    t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
//	}
func IntegrationGate() bool {
	return os.Getenv("ENABLE_INTEGRATION_TESTS") == "true"
}
