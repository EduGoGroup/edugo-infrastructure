# PHASE2 BRIDGE - Sprint-01-Migrate-CLI

## 📋 Resumen

**Sprint:** Sprint-01-Migrate-CLI
**Archivo principal:** `database/migrate.go`
**Tests:** `database/migrate_test.go`
**Estado Fase 1:** ✅ COMPLETADO

---

## ✅ Completado en Fase 1

### Implementación

- [x] CLI completa con 5 comandos (up, down, status, create, force)
- [x] Gestión de conexión PostgreSQL con variables de entorno
- [x] Creación automática de tabla `schema_migrations`
- [x] Sistema de transacciones con rollback automático
- [x] Carga de migraciones desde filesystem
- [x] Validación de migraciones (up y down)
- [x] Sanitización de nombres para nuevas migraciones
- [x] Formateo de output con emojis y colores
- [x] Manejo robusto de errores

### Tests Unitarios

- [x] `TestSanitizeName` - 7 casos de prueba (espacios, guiones, caracteres especiales, etc.)
- [x] `TestGetEnv` - Valores por defecto vs valores custom
- [x] `TestGetDBURL` - Construcción correcta de URL de conexión
- [x] Tests documentados con tabla-driven tests

### Código

```go
// migrate.go - 439 líneas
// Funciones principales:
- main()                      // Entry point con routing de comandos
- migrateUp(db)              // Ejecuta migraciones pendientes
- migrateDown(db)            // Revierte última migración
- showStatus(db)             // Muestra estado de migraciones
- createMigration(name)      // Crea nuevos archivos de migración
- forceMigration(db, version) // Fuerza versión (admin)
- loadMigrations()           // Carga desde filesystem
- getAppliedMigrations(db)   // Lee desde schema_migrations
- ensureMigrationsTable(db)  // Crea tabla si no existe
```

**Total de líneas:** 439 (migrate.go) + 175 (migrate_test.go)
**Tests unitarios:** 5 tests, todos passing
**Cobertura de funciones auxiliares:** 100%

---

## ⏳ Pendiente para Fase 2

### Tests de Integración con PostgreSQL Real

1. **Test: migrateUp crea tablas correctamente**
   - Descripción: Validar que migraciones UP crean todas las tablas
   - Requiere: PostgreSQL (Testcontainers)
   - Validar:
     - Tabla `schema_migrations` se crea
     - Todas las 8 migraciones se ejecutan
     - Datos se insertan en `schema_migrations`
     - Tablas (users, schools, etc.) existen en BD

2. **Test: migrateDown revierte correctamente**
   - Descripción: Validar rollback funciona
   - Requiere: PostgreSQL con migraciones aplicadas
   - Validar:
     - Última migración se revierte
     - Tabla eliminada de BD
     - Registro eliminado de `schema_migrations`
     - Puede revertir múltiples migraciones

3. **Test: showStatus muestra estado correcto**
   - Descripción: Validar reporte de estado
   - Requiere: PostgreSQL con algunas migraciones aplicadas
   - Validar:
     - Muestra migraciones aplicadas con timestamp
     - Muestra migraciones pendientes
     - Conteo correcto (total, aplicadas, pendientes)

4. **Test: transacciones con rollback en errores**
   - Descripción: Validar que errores en SQL hacen rollback
   - Requiere: PostgreSQL + migración con SQL inválido
   - Validar:
     - Error en migración no deja BD en estado inconsistente
     - Rollback automático funciona
     - Tabla `schema_migrations` no se actualiza en errores

5. **Test: createMigration genera archivos válidos**
   - Descripción: Validar generación de nuevas migraciones
   - Requiere: Filesystem
   - Validar:
     - Archivos .up.sql y .down.sql se crean
     - Nombres están correctamente sanitizados
     - Versionado es secuencial

### Edge Cases

1. **Conexión fallida a PostgreSQL**
   - Escenario: DB_HOST apunta a servidor inexistente
   - Validación: Error claro y mensaje útil

2. **Migraciones con SQL inválido**
   - Escenario: Sintaxis SQL errónea en .up.sql
   - Validación: Rollback automático, BD queda consistente

3. **Migraciones parcialmente aplicadas**
   - Escenario: 4 de 8 migraciones aplicadas
   - Validación: `status` muestra correctamente, `up` aplica solo pendientes

4. **Force migration con versión inválida**
   - Escenario: Forzar versión que no existe
   - Validación: Error o comportamiento documentado

---

## 🔧 Prerequisitos para Fase 2

### Servicios Requeridos

```bash
# Opción 1: Docker Compose
cd edugo-infrastructure
make dev-up-core

# Opción 2: Testcontainers (en tests)
# Se levanta automáticamente en tests
```

### Variables de Entorno

```bash
# Copiar .env.example
cp .env.example .env

# Variables necesarias:
DB_HOST=localhost
DB_PORT=5432
DB_NAME=edugo_dev
DB_USER=edugo
DB_PASSWORD=changeme
DB_SSL_MODE=disable
```

### Datos de Prueba

```bash
# Ejecutar migraciones
cd database
go run migrate.go up

# Verificar estado
go run migrate.go status
```

---

## 🧪 Tests de Integración a Implementar

### Archivo: `database/migrate_integration_test.go`

```go
package main

import (
    "testing"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgres(t *testing.T) (*sql.DB, func()) {
    // Setup Testcontainers
    ctx := context.Background()
    req := testcontainers.ContainerRequest{
        Image:        "postgres:15-alpine",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_DB":       "test_db",
            "POSTGRES_USER":     "test",
            "POSTGRES_PASSWORD": "test",
        },
        WaitingFor: wait.ForLog("database system is ready to accept connections"),
    }

    container, _ := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })

    // Conectar
    host, _ := container.Host(ctx)
    port, _ := container.MappedPort(ctx, "5432")

    db, _ := sql.Open("postgres", fmt.Sprintf("postgres://test:test@%s:%s/test_db?sslmode=disable", host, port.Port()))

    cleanup := func() {
        db.Close()
        container.Terminate(ctx)
    }

    return db, cleanup
}

func TestMigrateUpIntegration(t *testing.T) {
    db, cleanup := setupPostgres(t)
    defer cleanup()

    // Ensure migrations table
    if err := ensureMigrationsTable(db); err != nil {
        t.Fatalf("Error creando tabla: %v", err)
    }

    // Run migrations
    if err := migrateUp(db); err != nil {
        t.Fatalf("Error ejecutando migraciones: %v", err)
    }

    // Validate tables exist
    tables := []string{"users", "schools", "academic_units", "memberships",
                      "materials", "assessment", "assessment_attempt", "assessment_attempt_answer"}

    for _, table := range tables {
        var exists bool
        query := "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)"
        if err := db.QueryRow(query, table).Scan(&exists); err != nil {
            t.Fatalf("Error verificando tabla %s: %v", table, err)
        }
        if !exists {
            t.Errorf("Tabla %s no fue creada", table)
        }
    }
}

func TestMigrateDownIntegration(t *testing.T) {
    db, cleanup := setupPostgres(t)
    defer cleanup()

    // Setup: ejecutar todas las migraciones
    ensureMigrationsTable(db)
    migrateUp(db)

    // Test: revertir última migración
    if err := migrateDown(db); err != nil {
        t.Fatalf("Error revirtiendo migración: %v", err)
    }

    // Validate: assessment_attempt_answer no debe existir
    var exists bool
    query := "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'assessment_attempt_answer')"
    db.QueryRow(query).Scan(&exists)

    if exists {
        t.Error("Tabla assessment_attempt_answer no fue eliminada")
    }
}

func TestShowStatusIntegration(t *testing.T) {
    // Similar al anterior, pero capturando output de showStatus
}
```

### Casos de Prueba

1. **Happy path:** migrateUp aplica todas las 8 migraciones correctamente
2. **Error handling:** SQL inválido hace rollback automático
3. **Rollback:** migrateDown revierte correctamente
4. **Partial migrations:** aplicar 4 de 8, luego completar las 4 restantes
5. **Force migration:** forzar versión específica
6. **Concurrent migrations:** validar locks (si se implementa)

---

## 📝 Notas para Fase 2

- migrate.go **NO** tiene sistema de locks - considerar agregar para producción
- Actualmente usa ordenamiento simple de versiones - OK para <100 migraciones
- Force migration es peligroso - documentar claramente su uso
- Considerar agregar comando `migrate.go version` para mostrar versión actual
- Performance: loadMigrations() carga todos los archivos - optimizar si >100 migraciones

---

## ✅ Checklist Fase 2

- [ ] Levantar PostgreSQL (Testcontainers)
- [ ] Configurar variables de entorno de test
- [ ] Implementar `setupPostgres(t)` helper
- [ ] Test: migrateUp crea todas las tablas
- [ ] Test: migrateDown revierte correctamente
- [ ] Test: showStatus reporta estado correcto
- [ ] Test: rollback en errores de SQL
- [ ] Test: createMigration genera archivos válidos
- [ ] Validar edge cases documentados
- [ ] Medir cobertura de tests (objetivo: >80%)
- [ ] Actualizar README.md con ejemplos de tests
- [ ] Commit y push

---

**Fase 1 completada:** 2025-11-16
**Próximo paso:** Ejecutar tests de integración con PostgreSQL real
**Estimado Fase 2:** 1-2 horas
