//go:build integration
// +build integration

package framework_test

import (
	"os"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/fixtures"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/internal/testdb"
)

// Benchmarks del Bloque 8 de Fase C. Cumplen las cotas de C-REQ-7
// (5s/8s/30s) midiendo apply de fixtures atómicas y de una composición
// completa. El backend de BD se elige automáticamente:
//
//   - Local (default): testcontainers postgres:15-alpine + migrations
//     completas + production seed. Es el modo del CI/nightly.
//   - Cloud (override): si POSTGRES_URI está definido, se conecta
//     directamente a esa URI y NO levanta contenedor. Útil para medir
//     la cota real contra Neon sin instalar Docker.
//
// Ejecución típica:
//
//	# Contra contenedor (default):
//	ENABLE_INTEGRATION_TESTS=true go test -tags=integration -bench=. \
//	    -benchmem -count=1 -benchtime=1x -timeout=10m \
//	    ./seeds/e2e/framework/...
//
//	# Contra Neon:
//	source EduBack/edugo-dev-environment/migrator/.env.cloud
//	ENABLE_INTEGRATION_TESTS=true POSTGRES_URI=$POSTGRES_URI \
//	    go test -tags=integration -bench=. -benchmem -count=1 \
//	        -benchtime=1x -timeout=10m ./seeds/e2e/framework/...
//
// Cotas (C-REQ-7):
//
//   - BenchmarkFixtureRoleOnly             < 5 segundos
//   - BenchmarkRoleOnlyScreenOnlyCompose   < 8 segundos
//
// BenchmarkLegacyE2EFullApply fue removido en Fase 2 (ADR-6): el
// scenario legacy_e2e y sus fixtures dejaron de existir.
//
// Con testcontainers el setup de la primera iteración paga ~3-5 s
// extra por levantar el contenedor; ese costo NO se contabiliza en las
// cotas (queda fuera de b.N). Para medir la cota neta se recomienda
// `-benchtime=1x` y leer el `ns/op` del run individual.

// gateOrSkipBench salta el benchmark si ENABLE_INTEGRATION_TESTS no
// está activo. Centraliza el guard para que los 3 benchmarks compartan
// el mismo mensaje.
func gateOrSkipBench(b *testing.B) {
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		b.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration benchmarks")
	}
}

func BenchmarkFixtureRoleOnly(b *testing.B) {
	gateOrSkipBench(b)
	gdb := testdb.StartPostgres(b)
	c := framework.NewComposer(framework.NewRegistry(), framework.NewNopLogger())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.Compose(gdb, "bench_role_only", []framework.Fixture{
			&fixtures.RoleOnly{RoleCode: "teacher"},
		})
		if err != nil {
			b.Fatalf("Compose: %v", err)
		}
	}
}

func BenchmarkRoleOnlyScreenOnlyCompose(b *testing.B) {
	gateOrSkipBench(b)
	gdb := testdb.StartPostgres(b)
	c := framework.NewComposer(framework.NewRegistry(), framework.NewNopLogger())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.Compose(gdb, "bench_compose", []framework.Fixture{
			&fixtures.RoleOnly{RoleCode: "teacher"},
			&fixtures.ScreenOnly{ScreenKey: "assessments-list"},
		})
		if err != nil {
			b.Fatalf("Compose: %v", err)
		}
	}
}

