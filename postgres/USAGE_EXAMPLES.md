# Ejemplos de Uso - PostgreSQL Module

Esta guía muestra cómo los clientes (api-admin, api-mobile, worker) deben usar el módulo PostgreSQL de edugo-infrastructure.

## 📦 Instalación

### Paso 1: Agregar Dependencia

```bash
# En tu proyecto (ejemplo: edugo-api-mobile)
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.9.0
```

### Paso 2: Actualizar go.mod

```go
module github.com/EduGoGroup/edugo-api-mobile

go 1.24

require (
    github.com/EduGoGroup/edugo-infrastructure/postgres v0.9.0
    github.com/lib/pq v1.10.9
)
```

---

## 🎯 Uso 1: Tests de Integración (NUEVA API - Recomendado)

### Ejemplo Completo: Setup de Tests con Embed

```go
package tests

import (
    "database/sql"
    "os"
    "testing"
    
    "github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
    _ "github.com/lib/pq"
    "github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
    suite.Suite
    db *sql.DB
}

// SetupSuite - Se ejecuta UNA vez antes de todos los tests
func (s *IntegrationTestSuite) SetupSuite() {
    // 1. Conectar a PostgreSQL (puede ser testcontainer o BD local)
    connStr := os.Getenv("POSTGRES_TEST_URL")
    if connStr == "" {
        connStr = "host=localhost port=5432 user=edugo password=test dbname=edugo_test sslmode=disable"
    }
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        s.T().Fatalf("Error conectando a PostgreSQL: %v", err)
    }
    
    s.db = db
    
    // 2. Aplicar migraciones (estructura + constraints)
    if err := migrations.ApplyAll(s.db); err != nil {
        s.T().Fatalf("Error aplicando migraciones: %v", err)
    }
}

// SetupTest - Se ejecuta ANTES de cada test individual
func (s *IntegrationTestSuite) SetupTest() {
    // (Opcional) Aplicar datos mock para cada test
    if err := migrations.ApplyMockData(s.db); err != nil {
        s.T().Fatalf("Error aplicando mock data: %v", err)
    }
}

// TearDownSuite - Se ejecuta UNA vez después de todos los tests
func (s *IntegrationTestSuite) TearDownSuite() {
    if s.db != nil {
        s.db.Close()
    }
}

// Test de ejemplo
func (s *IntegrationTestSuite) TestCreateUser() {
    // Arrange
    query := `INSERT INTO users (id, email, role) VALUES ($1, $2, $3)`
    
    // Act
    _, err := s.db.Exec(query, "usr_123", "test@edugo.com", "student")
    
    // Assert
    s.NoError(err)
    
    var count int
    s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", "test@edugo.com").Scan(&count)
    s.Equal(1, count)
}

func TestIntegrationSuite(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration tests")
    }
    suite.Run(t, new(IntegrationTestSuite))
}
```

---

## 🎯 Uso 2: Migraciones en Aplicación (Producción)

### Opción A: CLI de Migraciones (Recomendado)

**Uso en CI/CD o scripts de deployment:**

```bash
#!/bin/bash
# scripts/migrate.sh

set -e

echo "Aplicando migraciones PostgreSQL..."

# Usar el CLI de infrastructure directamente
cd vendor/github.com/EduGoGroup/edugo-infrastructure/postgres
go run cmd/migrate/migrate.go up

echo "✅ Migraciones aplicadas exitosamente"
```

**Variables de entorno necesarias:**

```bash
export POSTGRES_HOST=db.production.com
export POSTGRES_PORT=5432
export POSTGRES_USER=edugo
export POSTGRES_PASSWORD=secret_password
export POSTGRES_DB=edugo_prod
export POSTGRES_SSLMODE=require
```

### Opción B: Programáticamente (Aplicación)

```go
package main

import (
    "database/sql"
    "log"
    
    pgtesting "github.com/EduGoGroup/edugo-infrastructure/postgres/testing"
    _ "github.com/lib/pq"
)

func main() {
    // 1. Conectar a PostgreSQL
    db, err := sql.Open("postgres", "host=localhost port=5432 user=edugo dbname=edugo_prod sslmode=disable")
    if err != nil {
        log.Fatalf("Error conectando: %v", err)
    }
    defer db.Close()
    
    // 2. Aplicar migraciones
    migrationsPath := "./vendor/github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
    if err := pgtesting.ApplyMigrations(db, migrationsPath); err != nil {
        log.Fatalf("Error aplicando migraciones: %v", err)
    }
    
    log.Println("✅ Migraciones aplicadas")
    
    // 3. Iniciar aplicación
    startServer(db)
}
```

---

## 🎯 Uso 3: Testcontainers (Recomendado para CI)

```go
package tests

import (
    "context"
    "database/sql"
    "fmt"
    "testing"
    "time"
    
    pgtesting "github.com/EduGoGroup/edugo-infrastructure/postgres/testing"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgresContainer(t *testing.T) (*sql.DB, func()) {
    ctx := context.Background()
    
    // 1. Crear container PostgreSQL
    req := testcontainers.ContainerRequest{
        Image:        "postgres:15-alpine",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_USER":     "edugo",
            "POSTGRES_PASSWORD": "test123",
            "POSTGRES_DB":       "edugo_test",
        },
        WaitingFor: wait.ForLog("database system is ready to accept connections").
            WithOccurrence(2).
            WithStartupTimeout(30 * time.Second),
    }
    
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        t.Fatalf("Error creando container: %v", err)
    }
    
    // 2. Obtener puerto mapeado
    host, _ := container.Host(ctx)
    port, _ := container.MappedPort(ctx, "5432")
    
    // 3. Conectar a PostgreSQL
    connStr := fmt.Sprintf("host=%s port=%s user=edugo password=test123 dbname=edugo_test sslmode=disable",
        host, port.Port())
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        t.Fatalf("Error conectando: %v", err)
    }
    
    // 4. Aplicar migraciones de infrastructure
    migrationsPath := "../vendor/github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
    if err := pgtesting.ApplyMigrations(db, migrationsPath); err != nil {
        t.Fatalf("Error aplicando migraciones: %v", err)
    }
    
    // 5. Retornar DB y función de cleanup
    cleanup := func() {
        db.Close()
        container.Terminate(ctx)
    }
    
    return db, cleanup
}

// Uso en test
func TestWithContainer(t *testing.T) {
    db, cleanup := setupPostgresContainer(t)
    defer cleanup()
    
    // Tu test aquí
    var result int
    db.QueryRow("SELECT COUNT(*) FROM users").Scan(&result)
    if result != 0 {
        t.Errorf("Expected 0 users, got %d", result)
    }
}
```

---

## 🎯 Uso 4: Seeds Personalizados

### Crear Seeds en tu Proyecto

```sql
-- edugo-api-mobile/testdata/seeds/demo_users.sql

INSERT INTO users (id, email, role, created_at, updated_at) VALUES
('usr_student_1', 'student1@test.com', 'student', NOW(), NOW()),
('usr_teacher_1', 'teacher1@test.com', 'teacher', NOW(), NOW()),
('usr_admin_1', 'admin1@test.com', 'admin', NOW(), NOW());

INSERT INTO schools (id, name, address, created_at, updated_at) VALUES
('sch_test_1', 'Escuela Demo', 'Calle Falsa 123', NOW(), NOW());
```

### Aplicar Seeds en Tests

```go
func (s *IntegrationTestSuite) SetupTest() {
    // 1. Limpiar BD
    pgtesting.CleanDatabase(s.db)
    
    // 2. Aplicar seeds de infrastructure (si existen)
    infraSeeds := "../vendor/github.com/EduGoGroup/edugo-infrastructure/postgres/seeds"
    pgtesting.ApplySeeds(s.db, infraSeeds)
    
    // 3. Aplicar seeds de tu proyecto
    projectSeeds := "../testdata/seeds"
    pgtesting.ApplySeeds(s.db, projectSeeds)
}
```

---

## 📚 Referencia de Funciones

### `pgtesting.ApplyMigrations(db, path)`

- **Qué hace:** Aplica todas las migraciones pendientes desde un directorio
- **Cuándo usar:** SetupSuite (una vez antes de todos los tests)
- **Idempotente:** Sí (solo aplica migraciones nuevas)

### `pgtesting.CleanDatabase(db)`

- **Qué hace:** Trunca todas las tablas excepto schema_migrations
- **Cuándo usar:** SetupTest (antes de cada test individual)
- **Seguro:** Sí (desactiva triggers temporalmente)

### `pgtesting.ApplySeeds(db, path)`

- **Qué hace:** Ejecuta archivos .sql de un directorio en orden correcto
- **Cuándo usar:** SetupTest (después de CleanDatabase)
- **Orden:** Automático por dependencias de FK

---

## ⚠️ Notas Importantes

### 1. Path a Migraciones

El path cambia según tu estructura:

```go
// Opción 1: Vendored
"../vendor/github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"

// Opción 2: Go modules cache
// Go lo encuentra automáticamente si usas:
import _ "github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
// Pero necesitas el path físico para ApplyMigrations
```

### 2. Tests Cortos vs Integración

```go
func TestIntegrationSuite(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration tests")
    }
    suite.Run(t, new(IntegrationTestSuite))
}
```

```bash
# Solo tests unitarios (rápido)
go test -short ./...

# Todos los tests incluido integración
go test ./...
```

### 3. Variables de Entorno

El CLI de migraciones usa estas variables:

- `POSTGRES_HOST` (default: localhost)
- `POSTGRES_PORT` (default: 5432)
- `POSTGRES_USER` (default: edugo)
- `POSTGRES_PASSWORD` (default: edugo_dev_2024)
- `POSTGRES_DB` (default: edugo_db)
- `POSTGRES_SSLMODE` (default: disable)

---

## 🚀 Makefile de Ejemplo

```makefile
# edugo-api-mobile/Makefile

.PHONY: test-integration
test-integration: ## Tests de integración con PostgreSQL
	@echo "🧪 Ejecutando tests de integración..."
	@go test -v ./tests/integration/...

.PHONY: test-all
test-all: ## Todos los tests
	@go test -v ./...

.PHONY: migrate-up
migrate-up: ## Aplicar migraciones
	@cd vendor/github.com/EduGoGroup/edugo-infrastructure/postgres && \
	go run cmd/migrate/migrate.go up

.PHONY: migrate-status
migrate-status: ## Ver estado de migraciones
	@cd vendor/github.com/EduGoGroup/edugo-infrastructure/postgres && \
	go run cmd/migrate/migrate.go status
```

---

## ❓ FAQ

### ¿Debo usar ApplyMigrations en producción?

**No recomendado.** Usa el CLI de migraciones (`cmd/migrate/migrate.go`) en scripts de deployment.

### ¿Puedo modificar las migraciones de infrastructure?

**No.** Las migraciones de infrastructure son inmutables. Si necesitas cambios, crea migraciones locales en tu proyecto.

### ¿Cómo agrego migraciones específicas de mi proyecto?

Crea tu propio sistema de migraciones local o usa herramientas como `golang-migrate/migrate`.

---

**Generado por:** edugo-infrastructure  
**Versión:** v0.9.0  
**Fecha:** 2025-11-18

## 📝 Notas de Arquitectura

**PostgreSQL usa archivos SQL embebidos:**
- `embed.FS` con archivos `.sql`
- API consistente con MongoDB
- Ambos módulos: `ApplyAll()`, `ApplyStructure()`, `ApplyConstraints()`, etc.
