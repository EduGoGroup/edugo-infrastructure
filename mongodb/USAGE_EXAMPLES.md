# Ejemplos de Uso - MongoDB Module

Esta guía muestra cómo los clientes (api-admin, api-mobile, worker) deben usar el módulo MongoDB de edugo-infrastructure.

## 📦 Instalación

### Paso 1: Agregar Dependencia

```bash
# En tu proyecto Go (ejemplo: edugo-worker)
go get github.com/EduGoGroup/edugo-infrastructure/mongodb@v0.9.0
```

### Paso 2: Actualizar go.mod

```go
module github.com/EduGoGroup/edugo-worker

go 1.24

require (
    github.com/EduGoGroup/edugo-infrastructure/mongodb v0.9.0
    go.mongodb.org/mongo-driver v1.17.3
)
```

---

## 🎯 Uso 1: Tests de Integración (NUEVA API - Recomendado)

### Ejemplo Completo: Setup de Tests con Embed

```go
package tests

import (
    "context"
    "os"
    "testing"
    "time"
    
    "github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations"
    "github.com/stretchr/testify/suite"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type IntegrationTestSuite struct {
    suite.Suite
    client *mongo.Client
    db     *mongo.Database
    ctx    context.Context
}

// SetupSuite - Se ejecuta UNA vez antes de todos los tests
func (s *IntegrationTestSuite) SetupSuite() {
    s.ctx = context.Background()
    
    // 1. Conectar a MongoDB
    uri := os.Getenv("MONGO_TEST_URL")
    if uri == "" {
        uri = "mongodb://localhost:27017"
    }
    
    client, err := mongo.Connect(s.ctx, options.Client().ApplyURI(uri))
    if err != nil {
        s.T().Fatalf("Error conectando a MongoDB: %v", err)
    }
    
    s.client = client
    s.db = client.Database("edugo_test")
    
    // 2. Aplicar migraciones (estructura + constraints)
    if err := migrations.ApplyAll(s.ctx, s.db); err != nil {
        s.T().Fatalf("Error aplicando migraciones: %v", err)
    }
}

// SetupTest - Se ejecuta ANTES de cada test individual
func (s *IntegrationTestSuite) SetupTest() {
    // Limpiar collections entre tests
    collections, _ := s.db.ListCollectionNames(s.ctx, map[string]interface{}{})
    for _, collection := range collections {
        s.db.Collection(collection).Drop(s.ctx)
    }
    
    // Recrear estructura
    migrations.ApplyAll(s.ctx, s.db)
    
    // (Opcional) Aplicar datos mock para cada test
    if err := migrations.ApplyMockData(s.ctx, s.db); err != nil {
        s.T().Fatalf("Error aplicando mock data: %v", err)
    }
}

// TearDownSuite - Se ejecuta UNA vez después de todos los tests
func (s *IntegrationTestSuite) TearDownSuite() {
    if s.client != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        s.client.Disconnect(ctx)
    }
}

// Test de ejemplo
func (s *IntegrationTestSuite) TestCreateAssessment() {
    // Arrange
    assessment := map[string]interface{}{
        "material_id": "mat_test_123",
        "questions":   []interface{}{},
        "metadata": map[string]interface{}{
            "subject":    "Mathematics",
            "difficulty": "easy",
        },
        "created_at": time.Now(),
        "updated_at": time.Now(),
    }
    
    // Act
    result, err := s.db.Collection("material_assessment").InsertOne(s.ctx, assessment)
    
    // Assert
    s.NoError(err)
    s.NotNil(result.InsertedID)
}

func TestIntegrationSuite(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration tests")
    }
    suite.Run(t, new(IntegrationTestSuite))
}
```

---

## 🎯 Uso 2: Aplicación en Producción

### Inicialización al Arrancar

```go
package main

import (
    "context"
    "log"
    "os"
    
    "github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    ctx := context.Background()
    
    // 1. Conectar a MongoDB
    uri := os.Getenv("MONGO_URL")
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatalf("Error conectando: %v", err)
    }
    defer client.Disconnect(ctx)
    
    db := client.Database("edugo_prod")
    
    // 2. Aplicar migraciones si es necesario
    if os.Getenv("RUN_MIGRATIONS") == "true" {
        log.Println("📊 Aplicando migraciones...")
        if err := migrations.ApplyAll(ctx, db); err != nil {
            log.Fatalf("Error: %v", err)
        }
        log.Println("✅ Migraciones aplicadas")
    }
    
    // 3. Aplicar seeds en desarrollo
    if os.Getenv("ENV") == "development" {
        if err := migrations.ApplySeeds(ctx, db); err != nil {
            log.Printf("Warning: %v", err)
        }
    }
    
    // 4. Iniciar aplicación
    startServer(db)
}
```

---

## 🎯 Uso 3: API Flexible

### Aplicar por Capas

```go
import "github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations"

ctx := context.Background()

// Opción 1: Todo de una vez
if err := migrations.ApplyAll(ctx, db); err != nil {
    log.Fatal(err)
}

// Opción 2: Capa por capa
if err := migrations.ApplyStructure(ctx, db); err != nil {
    log.Fatal(err)
}
if err := migrations.ApplyConstraints(ctx, db); err != nil {
    log.Fatal(err)
}
if err := migrations.ApplySeeds(ctx, db); err != nil {
    log.Fatal(err)
}
if err := migrations.ApplyMockData(ctx, db); err != nil {
    log.Fatal(err)
}
```

### Listar Funciones Disponibles

```go
functions := migrations.ListFunctions()

for layer, funcs := range functions {
    fmt.Printf("%s:\n", layer)
    for _, funcName := range funcs {
        fmt.Printf("  - %s\n", funcName)
    }
}

// Output:
// structure:
//   - CreateMaterialAssessment
//   - CreateMaterialContent
//   - CreateAssessmentAttemptResult
//   ...
// constraints:
//   - CreateMaterialAssessmentIndexes
//   - CreateMaterialContentIndexes
//   ...
```

### Usar Funciones Individuales

```go
import (
    "github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations/structure"
    "github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations/constraints"
)

// Crear una collection específica
if err := structure.CreateMaterialAssessment(ctx, db); err != nil {
    log.Fatal(err)
}

// Crear índices específicos
if err := constraints.CreateMaterialAssessmentIndexes(ctx, db); err != nil {
    log.Fatal(err)
}
```

---

## 🎯 Uso 4: Con Testcontainers (Recomendado para CI)

```go
package tests

import (
    "context"
    "testing"
    
    "github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func setupMongoContainer(t *testing.T) (*mongo.Database, func()) {
    ctx := context.Background()
    
    // 1. Crear container MongoDB
    req := testcontainers.ContainerRequest{
        Image:        "mongo:7",
        ExposedPorts: []string{"27017/tcp"},
        WaitingFor:   wait.ForLog("Waiting for connections"),
    }
    
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        t.Fatalf("Error creando container: %v", err)
    }
    
    // 2. Obtener endpoint
    host, _ := container.Host(ctx)
    port, _ := container.MappedPort(ctx, "27017")
    uri := fmt.Sprintf("mongodb://%s:%s", host, port.Port())
    
    // 3. Conectar
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        t.Fatalf("Error conectando: %v", err)
    }
    
    db := client.Database("edugo_test")
    
    // 4. Aplicar migraciones
    if err := migrations.ApplyAll(ctx, db); err != nil {
        t.Fatalf("Error aplicando migraciones: %v", err)
    }
    
    // 5. Retornar DB y función de cleanup
    cleanup := func() {
        client.Disconnect(ctx)
        container.Terminate(ctx)
    }
    
    return db, cleanup
}

// Uso en test
func TestWithContainer(t *testing.T) {
    db, cleanup := setupMongoContainer(t)
    defer cleanup()
    
    // Tu test aquí
    ctx := context.Background()
    count, _ := db.Collection("material_assessment").CountDocuments(ctx, map[string]interface{}{})
    
    if count != 0 {
        t.Errorf("Expected 0 documents, got %d", count)
    }
}
```

---

## 📚 Referencia de Funciones

### `migrations.ApplyAll(ctx, db)`

- **Qué hace:** Aplica structure + constraints
- **Cuándo usar:** Inicializar BD completa
- **Retorna:** `error`

### `migrations.ApplyStructure(ctx, db)`

- **Qué hace:** Crea collections con validators
- **Cuándo usar:** Crear solo las collections
- **Retorna:** `error`

### `migrations.ApplyConstraints(ctx, db)`

- **Qué hace:** Crea índices
- **Cuándo usar:** Después de ApplyStructure()
- **Retorna:** `error`

### `migrations.ApplySeeds(ctx, db)`

- **Qué hace:** Inserta datos iniciales del ecosistema
- **Cuándo usar:** Producción/staging
- **Retorna:** `error`

### `migrations.ApplyMockData(ctx, db)`

- **Qué hace:** Inserta datos mock para testing
- **Cuándo usar:** Tests, desarrollo
- **Retorna:** `error`

### `migrations.ListFunctions()`

- **Retorna:** `map[string][]string` - Mapa de capa -> array de nombres de funciones
- **Uso:** Debugging y documentación
- **Ejemplo:**
  ```go
  functions := migrations.ListFunctions()
  fmt.Println(functions)
  // {
  //   "structure": ["CreateMaterialAssessment", "CreateMaterialContent", ...],
  //   "constraints": ["CreateMaterialAssessmentIndexes", ...],
  //   "seeds": [],
  //   "testing": []
  // }
  ```

---

## ⚠️ Notas Importantes

### 1. Orden de Ejecución

```go
// ✅ CORRECTO
migrations.ApplyStructure(ctx, db)
migrations.ApplyConstraints(ctx, db)

// ❌ INCORRECTO - constraints fallarán sin collections
migrations.ApplyConstraints(ctx, db)
migrations.ApplyStructure(ctx, db)
```

### 2. Idempotencia

Los scripts de structure **NO son idempotentes**:

```go
// ❌ Esto fallará si ya se ejecutó
migrations.ApplyStructure(ctx, db)
migrations.ApplyStructure(ctx, db) // Error: collection already exists
```

Para desarrollo, limpiar primero:

```go
// Limpiar BD completa
collections, _ := db.ListCollectionNames(ctx, map[string]interface{}{})
for _, collection := range collections {
    db.Collection(collection).Drop(ctx)
}

// Ahora sí aplicar
migrations.ApplyAll(ctx, db)
```

### 3. Contexto

Siempre usar context con timeout en producción:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if err := migrations.ApplyAll(ctx, db); err != nil {
    log.Fatal(err)
}
```

---

## 🚀 Makefile de Ejemplo

```makefile
# Makefile para proyecto Go

.PHONY: test
test: ## Tests unitarios
	@go test -short ./...

.PHONY: test-integration
test-integration: ## Tests de integración
	@MONGO_TEST_URL=mongodb://localhost:27017 go test -v ./tests/integration/...

.PHONY: migrate
migrate: ## Aplicar migraciones MongoDB
	@echo "📊 Aplicando migraciones..."
	@go run cmd/migrate/main.go

.PHONY: dev-setup
dev-setup: ## Setup desarrollo (migraciones + seeds + mock)
	@echo "🚀 Iniciando setup de desarrollo..."
	@go run scripts/dev-setup.go
```

---

**Generado por:** edugo-infrastructure  
**Versión:** v0.9.0  
**Fecha:** 2025-11-18  

## 📝 Notas de Arquitectura

**Diferencia clave con PostgreSQL:**
- PostgreSQL usa `embed.FS` con archivos SQL
- MongoDB usa funciones Go nativas (no puede embeder `.go`)
- Ambos exponen la misma API: `ApplyAll()`, `ApplyStructure()`, `ApplyConstraints()`, etc.
- Consistencia en la experiencia del desarrollador
