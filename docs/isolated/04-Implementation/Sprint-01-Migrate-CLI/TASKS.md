# TASKS - Sprint-01-Migrate-CLI

## ✅ Fase 1 - COMPLETADAS

### Implementación de migrate.go

- [x] **Estructura base del CLI**
  - [x] main() con routing de comandos
  - [x] printHelp() con documentación de comandos
  - [x] getEnv() y getDBURL() para configuración

- [x] **Gestión de conexión PostgreSQL**
  - [x] Conexión con variables de entorno
  - [x] Validación de conexión con db.Ping()
  - [x] Cierre correcto de conexiones

- [x] **Tabla schema_migrations**
  - [x] ensureMigrationsTable() crea tabla si no existe
  - [x] Campos: version (INT PRIMARY KEY), name (VARCHAR), applied_at (TIMESTAMP)

- [x] **Comando: migrate up**
  - [x] loadMigrations() carga archivos .up.sql y .down.sql
  - [x] getAppliedMigrations() lee desde schema_migrations
  - [x] migrateUp() ejecuta solo migraciones pendientes
  - [x] Transacciones con rollback automático en errores
  - [x] Registro en schema_migrations después de aplicar
  - [x] Output formateado con emojis ✅

- [x] **Comando: migrate down**
  - [x] migrateDown() encuentra última migración aplicada
  - [x] Ejecuta SQL de .down.sql
  - [x] Elimina registro de schema_migrations
  - [x] Transacción con rollback automático

- [x] **Comando: migrate status**
  - [x] showStatus() lista todas las migraciones
  - [x] Marca aplicadas con ✅ y timestamp
  - [x] Marca pendientes con ⬜
  - [x] Muestra conteo total/aplicadas/pendientes

- [x] **Comando: migrate create**
  - [x] createMigration(name) genera archivos .up.sql y .down.sql
  - [x] sanitizeName() limpia caracteres especiales
  - [x] Versionado secuencial automático (001, 002, etc.)
  - [x] Templates con comentarios de fecha

- [x] **Comando: migrate force**
  - [x] forceMigration() fuerza versión específica
  - [x] Limpia schema_migrations
  - [x] Advertencia en output ⚠️

### Tests Unitarios

- [x] **TestSanitizeName**
  - [x] Espacios → underscores
  - [x] Guiones → underscores
  - [x] Mayúsculas → minúsculas
  - [x] Caracteres especiales → eliminados
  - [x] Números → preservados

- [x] **TestGetEnv**
  - [x] Retorna valor de env var cuando está seteada
  - [x] Retorna default cuando env var no existe

- [x] **TestGetDBURL**
  - [x] Construye URL correcta con defaults
  - [x] Construye URL correcta con env vars custom

- [x] **TestLoadMigrations** (skipped)
  - [x] Marcado como skip (requiere refactoring para testing)

- [x] **TestCreateMigrationFiles**
  - [x] Smoke test de sanitización de nombres

### Documentación

- [x] Comentarios inline en código
- [x] README.md del sprint
- [x] PHASE2_BRIDGE.md con pendientes
- [x] Ejemplos de uso en README principal

---

## ⏳ Fase 2 - PENDIENTES

### Tests de Integración

- [ ] **TestMigrateUpIntegration**
  - [ ] Setup PostgreSQL con Testcontainers
  - [ ] Ejecutar todas las 8 migraciones
  - [ ] Validar que tablas existen (users, schools, etc.)
  - [ ] Validar registros en schema_migrations

- [ ] **TestMigrateDownIntegration**
  - [ ] Setup BD con migraciones aplicadas
  - [ ] Revertir última migración
  - [ ] Validar que tabla fue eliminada
  - [ ] Validar que registro fue eliminado de schema_migrations

- [ ] **TestShowStatusIntegration**
  - [ ] Setup BD con algunas migraciones aplicadas
  - [ ] Ejecutar showStatus
  - [ ] Validar output (migraciones aplicadas vs pendientes)

- [ ] **TestTransactionRollback**
  - [ ] Crear migración con SQL inválido
  - [ ] Intentar ejecutar migrateUp
  - [ ] Validar que rollback funcionó
  - [ ] Validar que BD quedó consistente

- [ ] **TestCreateMigrationIntegration**
  - [ ] Ejecutar createMigration con nombre de prueba
  - [ ] Validar que archivos .up.sql y .down.sql se crearon
  - [ ] Validar contenido de archivos

### Edge Cases

- [ ] **Conexión fallida a PostgreSQL**
  - [ ] DB_HOST apunta a servidor inexistente
  - [ ] Validar mensaje de error claro

- [ ] **SQL inválido en migración**
  - [ ] Migración con sintaxis SQL errónea
  - [ ] Validar rollback automático

- [ ] **Migraciones parcialmente aplicadas**
  - [ ] Aplicar solo 4 de 8 migraciones
  - [ ] Ejecutar status
  - [ ] Ejecutar up y validar que aplica solo pendientes

- [ ] **Force migration con versión inválida**
  - [ ] Intentar forzar versión que no existe
  - [ ] Documentar comportamiento

### Mejoras Futuras

- [ ] Sistema de locks para evitar ejecuciones concurrentes
- [ ] Comando `version` para mostrar versión actual de BD
- [ ] Rollback múltiple (down N)
- [ ] Dry-run mode (mostrar SQL sin ejecutar)
- [ ] Mejor manejo de errores con tipos custom
- [ ] Logging estructurado (JSON)

---

## 📊 Métricas

### Fase 1
- **Líneas de código:** 439 (migrate.go) + 175 (migrate_test.go) = 614 total
- **Tests unitarios:** 5 tests
- **Tests passing:** 4/5 (1 skipped)
- **Cobertura:** 100% de funciones auxiliares
- **Comandos implementados:** 5/5 (up, down, status, create, force)

### Fase 2 (objetivos)
- **Tests de integración:** 5+
- **Edge cases validados:** 4+
- **Cobertura total:** >80%
- **Performance:** <1s para 8 migraciones

---

## 🔗 Referencias

- Código: `database/migrate.go`
- Tests: `database/migrate_test.go`
- Docs: `README.md`, `PHASE2_BRIDGE.md`
- Migraciones SQL: `database/migrations/postgres/`
