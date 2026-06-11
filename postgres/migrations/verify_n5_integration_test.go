package migrations_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestN5InfraGate(t *testing.T) {
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true")
	}
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "test",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)
	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432/tcp")
	dsn := "postgres://test:test@" + host + ":" + port.Port() + "/test?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	_, err = migrations.Migrate(db, migrations.MigrateOptions{Force: true, SeedDemo: false, DBUser: "test"})
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	var dtExists, scExists bool
	_ = db.QueryRow(`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema='notifications' AND table_name='device_tokens')`).Scan(&dtExists)
	_ = db.QueryRow(`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema='auth' AND table_name='service_clients')`).Scan(&scExists)
	if !dtExists {
		t.Fatal("missing notifications.device_tokens")
	}
	if !scExists {
		t.Fatal("missing auth.service_clients")
	}
	var count int
	_ = db.QueryRow(`SELECT COUNT(*) FROM auth.service_clients WHERE is_active = true AND client_id IN ('edugo-worker','edugo-api-learning')`).Scan(&count)
	if count != 2 {
		t.Fatalf("expected 2 active M2M clients, got %d", count)
	}
	var idx1, idx2 bool
	_ = db.QueryRow(`SELECT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname='idx_device_tokens_user_active')`).Scan(&idx1)
	_ = db.QueryRow(`SELECT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname='idx_service_clients_active')`).Scan(&idx2)
	if !idx1 || !idx2 {
		t.Fatalf("partial indexes missing: device=%v service=%v", idx1, idx2)
	}
}
