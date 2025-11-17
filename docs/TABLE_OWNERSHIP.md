# 🗄️ Ownership de Tablas - EduGo

**Owner:** edugo-infrastructure  
**Fecha:** 15 de Noviembre, 2025

---

## 🎯 Propósito

Este documento define QUIÉN crea y mantiene cada tabla en PostgreSQL.

**Regla de oro:** Solo el owner puede crear/modificar/eliminar la tabla.

---

## 📊 Tabla de Ownership

| Tabla | Owner | Creada en migración | Usada por | Puede modificar |
|-------|-------|-------------------|-----------|-----------------|
| **users** | infrastructure | 001 | api-admin, api-mobile, worker | Solo infrastructure |
| **schools** | infrastructure | 002 | api-admin, api-mobile | Solo infrastructure |
| **academic_units** | infrastructure | 003 | api-admin, api-mobile | Solo infrastructure |
| **memberships** | infrastructure | 004 | api-admin, api-mobile | Solo infrastructure |
| **materials** | infrastructure | 005 | api-mobile, worker | Solo infrastructure |
| **assessment** | infrastructure | 006 | api-mobile, worker | Solo infrastructure |
| **assessment_attempt** | infrastructure | 007 | api-mobile | Solo infrastructure |
| **assessment_attempt_answer** | infrastructure | 008 | api-mobile | Solo infrastructure |

---

## ✅ Reglas de Ownership

### 1. Solo infrastructure crea tablas

```sql
-- ✅ CORRECTO: En edugo-infrastructure/database/migrations/
CREATE TABLE users (...);

-- ❌ INCORRECTO: En api-admin/migrations/ o api-mobile/migrations/
CREATE TABLE users (...);  -- NO HACER ESTO
```

### 2. Proyectos solo USAN las tablas

```go
// ✅ CORRECTO: En api-admin o api-mobile
db.Query("SELECT * FROM users WHERE email = $1", email)
db.Exec("INSERT INTO users (...) VALUES (...)")

// ❌ INCORRECTO: En api-admin o api-mobile
db.Exec("CREATE TABLE users (...)")  // NO HACER ESTO
db.Exec("ALTER TABLE users ADD COLUMN avatar TEXT")  // NO HACER ESTO
```

### 3. Cambios de esquema = Nueva migración en infrastructure

Si necesitas agregar columna `avatar` a `users`:

```bash
# En edugo-infrastructure/
cd database
go run migrate.go create "add_avatar_to_users"

# Editar: migrations/postgres/009_add_avatar_to_users.up.sql
ALTER TABLE users ADD COLUMN avatar TEXT;

# Editar: migrations/postgres/009_add_avatar_to_users.down.sql
ALTER TABLE users DROP COLUMN avatar;

# Ejecutar migración
go run migrate.go up
```

---

## 🔄 Orden de Ejecución de Migraciones

### En Desarrollo Local

```bash
# 1. Levantar PostgreSQL
cd edugo-infrastructure
make dev-up-core

# 2. Ejecutar migraciones de infrastructure
cd database
go run migrate.go up

# ✅ Listo! Todas las tablas creadas
# Ahora puedes correr api-admin, api-mobile, worker
```

### En CI/CD

```yaml
# .github/workflows/ci.yml
jobs:
  test:
    steps:
      - name: Setup infrastructure
        run: |
          cd edugo-infrastructure
          docker-compose up -d postgres
          cd database && go run migrate.go up
      
      - name: Test api-admin
        run: cd edugo-api-admin && go test ./...
      
      - name: Test api-mobile
        run: cd edugo-api-mobile && go test ./...
```

---

## 🚫 Anti-Patrones (NO HACER)

### ❌ Anti-Patrón 1: Duplicar definiciones

```sql
-- ❌ MAL: En api-admin/migrations/001.sql
CREATE TABLE IF NOT EXISTS users (...);

-- ❌ MAL: En api-mobile/migrations/001.sql
CREATE TABLE IF NOT EXISTS users (...);

-- Problema: Definiciones pueden divergir
```

### ❌ Anti-Patrón 2: Migraciones en proyectos individuales

```bash
# ❌ MAL
api-admin/
  migrations/
    001_create_users.sql  # NO!

api-mobile/
  migrations/
    001_create_materials.sql  # NO!
```

### ✅ Correcto: Todo en infrastructure

```bash
# ✅ BIEN
edugo-infrastructure/
  database/
    migrations/
      postgres/
        001_create_users.sql
        002_create_schools.sql
        003_create_academic_units.sql
        004_create_memberships.sql
        005_create_materials.sql
        006_create_assessments.sql
        007_create_assessment_attempts.sql
        008_create_assessment_answers.sql
```

---

## 🔍 Validación en Proyectos

Cada proyecto puede validar que tablas existen antes de correr:

```makefile
# api-admin/Makefile
.PHONY: validate-schema
validate-schema:
	@psql $(DB_URL) -c "SELECT 1 FROM users LIMIT 1" > /dev/null || \
	  (echo "ERROR: Tabla users no existe. Ejecutar migraciones de infrastructure primero" && exit 1)
	@echo "✅ Schema validado"

.PHONY: run
run: validate-schema
	go run cmd/api/main.go
```

---

## 📋 Preguntas Frecuentes

### P: ¿Qué pasa si necesito una tabla solo para api-mobile?

**R:** Igual se crea en `edugo-infrastructure`.

Ejemplo: Si necesitas tabla `student_progress` solo para api-mobile:
- Crear migración en `edugo-infrastructure/database/migrations/postgres/009_create_student_progress.sql`
- Documentar en este archivo que es usada solo por api-mobile
- api-mobile la usa normalmente

### P: ¿Puedo tener migrations/ en api-admin para cosas específicas?

**R:** NO. Todas las migraciones van en infrastructure.

Si tienes dudas si algo es "específico", pregúntate:
- ¿Modifica el esquema de PostgreSQL? → infrastructure
- ¿Es configuración de la app? → .env del proyecto
- ¿Es seed de datos? → infrastructure/seeds/

### P: ¿Cómo manejo datos iniciales (seeds)?

**R:** Seeds también en infrastructure:

```bash
edugo-infrastructure/
  seeds/
    postgres/
      users.sql      # 3 usuarios de prueba
      schools.sql    # 2 escuelas de prueba
```

---

## ✅ Checklist para Nuevas Tablas

- [ ] Crear migración UP en `database/migrations/postgres/00X_create_*.up.sql`
- [ ] Crear migración DOWN en `database/migrations/postgres/00X_create_*.down.sql`
- [ ] Agregar tabla a este documento (TABLE_OWNERSHIP.md)
- [ ] Especificar quién usa la tabla
- [ ] Agregar comentarios SQL (`COMMENT ON TABLE ...`)
- [ ] Crear índices necesarios
- [ ] Testear migración UP y DOWN localmente
- [ ] Commit en rama `dev` de infrastructure
- [ ] PR y merge

---

**Última actualización:** 15 de Noviembre, 2025  
**Mantenedor:** Equipo EduGo
