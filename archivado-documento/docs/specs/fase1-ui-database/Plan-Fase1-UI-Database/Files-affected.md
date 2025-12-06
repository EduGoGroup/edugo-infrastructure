# Archivos Afectados - FASE 1 UI Database

> **Lista completa de archivos a crear, modificar y eliminar**

---

## Nuevos Archivos a Crear

### Migraciones de Base de Datos

#### 1. Migración 011: user_active_context

**Archivo**: `postgres/migrations/structure/011_create_user_active_context.sql`
- **Tipo**: Nuevo
- **Propósito**: Crear tabla user_active_context
- **Contenido**: CREATE TABLE, índices, trigger, comentarios
- **Líneas estimadas**: ~50

**Archivo**: `postgres/migrations/constraints/011_create_user_active_context.sql`
- **Tipo**: Nuevo
- **Propósito**: Constraints adicionales (si aplica)
- **Contenido**: Vacío o comentario (seguir convención)
- **Líneas estimadas**: ~1

---

#### 2. Migración 012: user_favorites

**Archivo**: `postgres/migrations/structure/012_create_user_favorites.sql`
- **Tipo**: Nuevo
- **Propósito**: Crear tabla user_favorites
- **Contenido**: CREATE TABLE, índices, comentarios
- **Líneas estimadas**: ~40

**Archivo**: `postgres/migrations/constraints/012_create_user_favorites.sql`
- **Tipo**: Nuevo
- **Propósito**: Constraints adicionales (si aplica)
- **Contenido**: Vacío o comentario
- **Líneas estimadas**: ~1

---

#### 3. Migración 013: user_activity_log

**Archivo**: `postgres/migrations/structure/013_create_user_activity_log.sql`
- **Tipo**: Nuevo
- **Propósito**: Crear ENUM activity_type y tabla user_activity_log
- **Contenido**: CREATE TYPE, CREATE TABLE, índices, comentarios
- **Líneas estimadas**: ~70

**Archivo**: `postgres/migrations/constraints/013_create_user_activity_log.sql`
- **Tipo**: Nuevo
- **Propósito**: Constraints adicionales (si aplica)
- **Contenido**: Vacío o comentario
- **Líneas estimadas**: ~1

---

### Tests

#### 1. Test de Estructura

**Archivo**: `postgres/tests/test_fase1_structure.sql`
- **Tipo**: Nuevo
- **Propósito**: Validar estructura de las 3 tablas
- **Contenido**: Queries de validación de columnas, tipos, constraints
- **Líneas estimadas**: ~100

---

#### 2. Test de Performance

**Archivo**: `postgres/tests/test_fase1_performance.sql`
- **Tipo**: Nuevo
- **Propósito**: Validar performance de queries con datos de prueba
- **Contenido**: Inserts masivos, EXPLAIN ANALYZE de queries frecuentes
- **Líneas estimadas**: ~80

---

#### 3. Test de Integridad

**Archivo**: `postgres/tests/test_fase1_integrity.sql`
- **Tipo**: Nuevo
- **Propósito**: Validar constraints, CASCADE, SET NULL
- **Contenido**: Tests de FK constraints, UNIQUE, triggers
- **Líneas estimadas**: ~120

---

### Documentación del Plan

**Nota**: Estos archivos ya fueron creados durante la planificación

**Archivo**: `docs/specs/fase1-ui-database/README.md`
- **Tipo**: Nuevo
- **Propósito**: Resumen ejecutivo de la FASE 1
- **Estado**: ✅ Creado

**Archivo**: `docs/specs/fase1-ui-database/ANALISIS-TECNICO.md`
- **Tipo**: Nuevo
- **Propósito**: Análisis técnico detallado de las 3 tablas
- **Estado**: ✅ Creado

**Archivo**: `docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/README.md`
- **Tipo**: Nuevo
- **Propósito**: Índice del plan de trabajo
- **Estado**: ✅ Creado

**Archivo**: `docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/Planner.md`
- **Tipo**: Nuevo
- **Propósito**: Fases y pasos detallados
- **Estado**: ✅ Creado

**Archivo**: `docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/Planner-commit.md`
- **Tipo**: Nuevo
- **Propósito**: Estrategia de commits
- **Estado**: ✅ Creado

**Archivo**: `docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/Files-affected.md`
- **Tipo**: Nuevo
- **Propósito**: Este archivo
- **Estado**: 🔄 En creación

**Archivo**: `docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/Test-unit.md`
- **Tipo**: Nuevo
- **Propósito**: Especificación de tests unitarios
- **Estado**: ⏳ Pendiente

**Archivo**: `docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/error.md`
- **Tipo**: Nuevo
- **Propósito**: Template para tracking de errores
- **Estado**: ⏳ Pendiente

---

## Archivos a Modificar

### 1. postgres/README.md

**Ubicación**: `postgres/README.md`
- **Tipo**: Modificación
- **Sección a modificar**: Agregar sección de "Tablas" o buscar donde documentan schema
- **Cambios**:
  - Agregar descripción de `user_active_context`
  - Agregar descripción de `user_favorites`
  - Agregar descripción de `user_activity_log`
  - Incluir propósito, relaciones, queries comunes
  - Consideraciones de escala
- **Líneas a agregar**: ~150

**Contenido a agregar**:
```markdown
### Nuevas Tablas - FASE 1 UI Roadmap (v0.11.0)

#### user_active_context
[Descripción completa según Planner.md FASE 6]

#### user_favorites
[Descripción completa según Planner.md FASE 6]

#### user_activity_log
[Descripción completa según Planner.md FASE 6]
```

---

### 2. CHANGELOG.md

**Ubicación**: `CHANGELOG.md` (raíz del proyecto)
- **Tipo**: Modificación
- **Sección a modificar**: Agregar nueva versión al inicio
- **Cambios**:
  - Agregar sección `## [postgres/v0.11.0] - 2025-12-01`
  - Listar las 3 tablas nuevas bajo `### Added`
  - Listar archivos de migración
  - Mencionar testing y documentación
  - Incluir referencias (issue, PR)
- **Líneas a agregar**: ~80

**Ubicación en el archivo**:
```markdown
# Changelog

## [postgres/v0.11.0] - 2025-12-01    ← AGREGAR AQUÍ

### Added - FASE 1 UI Roadmap
...

## [postgres/v0.10.1] - [fecha anterior]   ← Esto ya existe
...
```

---

### 3. README.md (raíz del proyecto)

**Ubicación**: `README.md` (raíz del proyecto)
- **Tipo**: Modificación (opcional, según instrucción del comando)
- **Sección a modificar**: Agregar pequeño comentario al final
- **Cambios**:
  - Agregar link al plan de trabajo actual
  - Indicar último trabajo realizado
- **Líneas a agregar**: ~5

**Contenido a agregar** (al final del README):
```markdown
---

## 📋 Último Plan de Trabajo

**FASE 1: UI Database Infrastructure** - [Ver plan completo](./docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/README.md)

Implementación de 3 nuevas tablas PostgreSQL para soportar UI Roadmap:
- `user_active_context` - Contexto/escuela activa del usuario
- `user_favorites` - Materiales favoritos
- `user_activity_log` - Log de actividades

**Estado**: 🔄 En progreso  
**Fecha**: 1 de Diciembre, 2025
```

---

## Archivos a Eliminar

**Ninguno** - Esta fase solo agrega archivos nuevos.

---

## Estructura de Directorios Resultante

```
edugo-infrastructure/
├── postgres/
│   ├── migrations/
│   │   ├── structure/
│   │   │   ├── 001_create_users.sql
│   │   │   ├── ...
│   │   │   ├── 010_create_login_attempts.sql
│   │   │   ├── 011_create_user_active_context.sql    ← NUEVO
│   │   │   ├── 012_create_user_favorites.sql          ← NUEVO
│   │   │   └── 013_create_user_activity_log.sql       ← NUEVO
│   │   └── constraints/
│   │       ├── 001_create_users.sql
│   │       ├── ...
│   │       ├── 010_create_login_attempts.sql
│   │       ├── 011_create_user_active_context.sql    ← NUEVO
│   │       ├── 012_create_user_favorites.sql          ← NUEVO
│   │       └── 013_create_user_activity_log.sql       ← NUEVO
│   ├── tests/
│   │   ├── test_fase1_structure.sql                   ← NUEVO
│   │   ├── test_fase1_performance.sql                 ← NUEVO
│   │   └── test_fase1_integrity.sql                   ← NUEVO
│   └── README.md                                       ← MODIFICADO
├── docs/
│   └── specs/
│       └── fase1-ui-database/                         ← NUEVO DIRECTORIO
│           ├── README.md                              ← NUEVO
│           ├── ANALISIS-TECNICO.md                    ← NUEVO
│           └── Plan-Fase1-UI-Database/                ← NUEVO DIRECTORIO
│               ├── README.md                          ← NUEVO
│               ├── Planner.md                         ← NUEVO
│               ├── Planner-commit.md                  ← NUEVO
│               ├── Files-affected.md                  ← NUEVO (este archivo)
│               ├── Test-unit.md                       ← NUEVO
│               └── error.md                           ← NUEVO
├── CHANGELOG.md                                       ← MODIFICADO
└── README.md                                          ← MODIFICADO (opcional)
```

---

## Resumen de Cambios por Tipo

| Tipo de Cambio | Cantidad | Archivos |
|----------------|----------|----------|
| **Nuevos** | 15 | 6 migraciones + 3 tests + 6 docs |
| **Modificados** | 2-3 | README.md (postgres), CHANGELOG.md, README.md (raíz, opcional) |
| **Eliminados** | 0 | Ninguno |
| **Total** | 17-18 | - |

---

## Tamaño Estimado de Cambios

| Categoría | Líneas de Código | Porcentaje |
|-----------|------------------|------------|
| **SQL (migraciones)** | ~160 líneas | 20% |
| **SQL (tests)** | ~300 líneas | 38% |
| **Documentación** | ~330 líneas | 42% |
| **Total** | ~790 líneas | 100% |

---

## Dependencias entre Archivos

```
Planner.md
    ↓
011_create_user_active_context.sql ──┐
012_create_user_favorites.sql ────────┼─→ test_fase1_structure.sql
013_create_user_activity_log.sql ────┘      test_fase1_performance.sql
                                             test_fase1_integrity.sql
                                                    ↓
                                            postgres/README.md
                                            CHANGELOG.md
```

**Orden de creación**:
1. Migraciones (011, 012, 013)
2. Tests (structure, performance, integrity)
3. Documentación (README, CHANGELOG)

---

## Validación de Archivos

### Pre-commit Checklist

Antes de commitear cada archivo:

**Migraciones**:
```bash
# Validar sintaxis SQL
psql -U postgres -d edugo_db --dry-run -f <archivo.sql>

# O usar linter SQL si está disponible
sqlfluff lint <archivo.sql>
```

**Tests**:
```bash
# Ejecutar test y verificar salida
psql -U postgres -d edugo_db -f <test.sql>
```

**Documentación**:
```bash
# Validar Markdown
markdownlint <archivo.md>

# O verificar links
markdown-link-check <archivo.md>
```

---

## Backup y Seguridad

### Antes de ejecutar migraciones

```bash
# Backup de BD local
pg_dump -U postgres edugo_db > backup_before_fase1_$(date +%Y%m%d_%H%M%S).sql

# Verificar backup
ls -lh backup_before_fase1_*.sql
```

### En caso de error

```bash
# Restaurar desde backup
psql -U postgres -d edugo_db < backup_before_fase1_YYYYMMDD_HHMMSS.sql
```

---

## Checklist de Archivos

```
Migraciones:
□ postgres/migrations/structure/011_create_user_active_context.sql
□ postgres/migrations/constraints/011_create_user_active_context.sql
□ postgres/migrations/structure/012_create_user_favorites.sql
□ postgres/migrations/constraints/012_create_user_favorites.sql
□ postgres/migrations/structure/013_create_user_activity_log.sql
□ postgres/migrations/constraints/013_create_user_activity_log.sql

Tests:
□ postgres/tests/test_fase1_structure.sql
□ postgres/tests/test_fase1_performance.sql
□ postgres/tests/test_fase1_integrity.sql

Documentación:
□ postgres/README.md (modificación)
□ CHANGELOG.md (modificación)
□ README.md raíz (modificación opcional)
□ docs/specs/fase1-ui-database/README.md
□ docs/specs/fase1-ui-database/ANALISIS-TECNICO.md
□ docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/README.md
□ docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/Planner.md
□ docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/Planner-commit.md
□ docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/Files-affected.md
□ docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/Test-unit.md
□ docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/error.md
```

---

**Total de archivos a gestionar**: 21 archivos (15 nuevos + 6 docs ya creados + 2-3 modificaciones)
