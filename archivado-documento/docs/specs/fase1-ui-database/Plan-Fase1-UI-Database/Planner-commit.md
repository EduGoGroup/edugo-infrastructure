# Estrategia de Commits - FASE 1 UI Database

> **Commits atómicos y convencionales para trazabilidad**

---

## Filosofía de Commits

### Principios
1. **Atómico**: Cada commit representa una unidad completa de trabajo
2. **Convencional**: Seguir Conventional Commits (feat, fix, docs, test, etc.)
3. **Descriptivo**: Mensaje claro que explica QUÉ y POR QUÉ
4. **Trazable**: Referencia a issue/PR para contexto completo

### Estructura de Mensaje de Commit

```
<tipo>(<alcance>): <descripción corta>

<cuerpo opcional: explicación detallada>

<footer opcional: referencias>
```

**Ejemplo**:
```
feat(database): agregar tabla user_active_context para contexto de usuario

- Crear migración 011_create_user_active_context.sql
- Tabla almacena escuela activa del usuario para filtrado en UI
- Incluye constraints UNIQUE, FKs con CASCADE/SET NULL
- Índices en user_id y school_id para performance
- Trigger automático para updated_at

Parte de FASE 1 UI Roadmap - Bloquea APIs y UI

Relacionado: #45
```

---

## Tipos de Commit (Conventional Commits)

| Tipo | Uso | Ejemplo |
|------|-----|---------|
| `feat` | Nueva funcionalidad | `feat(database): agregar tabla user_favorites` |
| `fix` | Corrección de bug | `fix(migration): corregir FK en user_active_context` |
| `docs` | Solo documentación | `docs(database): actualizar README con nuevas tablas` |
| `test` | Tests (sin cambio de código) | `test(database): agregar tests de validación FASE 1` |
| `refactor` | Refactorización sin cambio de funcionalidad | `refactor(migration): optimizar índices` |
| `perf` | Mejora de performance | `perf(database): agregar índice parcial en activity_log` |
| `chore` | Tareas de mantenimiento | `chore(database): actualizar CHANGELOG` |

---

## Alcance (Scope)

Para este proyecto, usamos:
- `database`: Cambios en schema/migraciones
- `migration`: Específico de archivos de migración
- `test`: Tests
- `docs`: Documentación

---

## Plan de Commits para FASE 1

### Commit 1: FASE 2 - user_active_context

**Tipo**: `feat(database)`  
**Mensaje corto**: agregar tabla user_active_context para contexto de usuario

**Cuerpo**:
```
- Crear migración 011_create_user_active_context.sql (structure)
- Crear migración 011_create_user_active_context.sql (constraints)
- Tabla almacena escuela activa del usuario para filtrado en UI
- UNIQUE constraint en user_id (solo un contexto por usuario)
- Foreign keys con CASCADE a users/schools, SET NULL a academic_units
- Índices en user_id y school_id para performance
- Trigger automático para updated_at

Bloquea: APIs de contexto (FASE 2) y UI selector de escuela (FASE 4)

Parte de FASE 1 UI Roadmap
```

**Footer**:
```
Relacionado: #[número-de-issue]
```

**Archivos incluidos**:
- `postgres/migrations/structure/011_create_user_active_context.sql`
- `postgres/migrations/constraints/011_create_user_active_context.sql`

**Comando**:
```bash
git add postgres/migrations/structure/011_create_user_active_context.sql
git add postgres/migrations/constraints/011_create_user_active_context.sql
git commit -m "feat(database): agregar tabla user_active_context para contexto de usuario

- Crear migración 011_create_user_active_context.sql (structure)
- Crear migración 011_create_user_active_context.sql (constraints)
- Tabla almacena escuela activa del usuario para filtrado en UI
- UNIQUE constraint en user_id (solo un contexto por usuario)
- Foreign keys con CASCADE a users/schools, SET NULL a academic_units
- Índices en user_id y school_id para performance
- Trigger automático para updated_at

Bloquea: APIs de contexto (FASE 2) y UI selector de escuela (FASE 4)

Parte de FASE 1 UI Roadmap

Relacionado: #[número-de-issue]"
```

---

### Commit 2: FASE 3 - user_favorites

**Tipo**: `feat(database)`  
**Mensaje corto**: agregar tabla user_favorites para materiales favoritos

**Cuerpo**:
```
- Crear migración 012_create_user_favorites.sql (structure)
- Crear migración 012_create_user_favorites.sql (constraints)
- Tabla almacena materiales marcados como favoritos por usuarios
- UNIQUE constraint compuesto (user_id, material_id) evita duplicados
- CASCADE en ambos FKs para limpieza automática
- Índices en user_id, material_id y created_at
- Ordenamiento por created_at DESC para "favoritos recientes"

Bloquea: Funcionalidad de favoritos en UI (FASE 4)

Parte de FASE 1 UI Roadmap
```

**Footer**:
```
Relacionado: #[número-de-issue]
```

**Archivos incluidos**:
- `postgres/migrations/structure/012_create_user_favorites.sql`
- `postgres/migrations/constraints/012_create_user_favorites.sql`

**Comando**:
```bash
git add postgres/migrations/structure/012_create_user_favorites.sql
git add postgres/migrations/constraints/012_create_user_favorites.sql
git commit -m "feat(database): agregar tabla user_favorites para materiales favoritos

- Crear migración 012_create_user_favorites.sql (structure)
- Crear migración 012_create_user_favorites.sql (constraints)
- Tabla almacena materiales marcados como favoritos por usuarios
- UNIQUE constraint compuesto (user_id, material_id) evita duplicados
- CASCADE en ambos FKs para limpieza automática
- Índices en user_id, material_id y created_at
- Ordenamiento por created_at DESC para favoritos recientes

Bloquea: Funcionalidad de favoritos en UI (FASE 4)

Parte de FASE 1 UI Roadmap

Relacionado: #[número-de-issue]"
```

---

### Commit 3: FASE 4 - user_activity_log

**Tipo**: `feat(database)`  
**Mensaje corto**: agregar tabla user_activity_log para tracking de actividades

**Cuerpo**:
```
- Crear ENUM activity_type con 8 tipos de actividad
- Crear migración 013_create_user_activity_log.sql (structure)
- Crear migración 013_create_user_activity_log.sql (constraints)
- Tabla almacena log de actividades para historial y analytics
- JSONB metadata para datos flexibles por tipo de actividad
- SET NULL en FKs para preservar datos históricos
- Índices estratégicos:
  * (user_id, created_at DESC) para actividad reciente
  * (school_id, created_at DESC) para analytics por escuela
  * activity_type para agregaciones
  * Índice parcial para rate limiting (última hora)

Consideraciones de escala:
- Tabla de alto volumen (estimado 1M+ registros/día)
- Índice parcial reduce tamaño manteniendo funcionalidad

Bloquea: Actividad reciente en Home (FASE 4)

Parte de FASE 1 UI Roadmap
```

**Footer**:
```
Relacionado: #[número-de-issue]
```

**Archivos incluidos**:
- `postgres/migrations/structure/013_create_user_activity_log.sql`
- `postgres/migrations/constraints/013_create_user_activity_log.sql`

**Comando**:
```bash
git add postgres/migrations/structure/013_create_user_activity_log.sql
git add postgres/migrations/constraints/013_create_user_activity_log.sql
git commit -m "feat(database): agregar tabla user_activity_log para tracking de actividades

- Crear ENUM activity_type con 8 tipos de actividad
- Crear migración 013_create_user_activity_log.sql (structure)
- Crear migración 013_create_user_activity_log.sql (constraints)
- Tabla almacena log de actividades para historial y analytics
- JSONB metadata para datos flexibles por tipo de actividad
- SET NULL en FKs para preservar datos históricos
- Índices estratégicos:
  * (user_id, created_at DESC) para actividad reciente
  * (school_id, created_at DESC) para analytics por escuela
  * activity_type para agregaciones
  * Índice parcial para rate limiting (última hora)

Consideraciones de escala:
- Tabla de alto volumen (estimado 1M+ registros/día)
- Índice parcial reduce tamaño manteniendo funcionalidad

Bloquea: Actividad reciente en Home (FASE 4)

Parte de FASE 1 UI Roadmap

Relacionado: #[número-de-issue]"
```

---

### Commit 4: FASE 5 - Tests de validación

**Tipo**: `test(database)`  
**Mensaje corto**: agregar tests de validación para FASE 1 UI Database

**Cuerpo**:
```
- Tests de estructura de tablas (columnas, tipos, constraints)
- Tests de performance con 10K inserts en user_activity_log
- Tests de integridad referencial:
  * CASCADE en user_active_context y user_favorites
  * SET NULL en user_activity_log
- Tests de índices y query plans
- Validación de ENUM activity_type
- Validación de JSONB metadata
- Scripts ejecutables en postgres/tests/test_fase1_*.sql

Resultados: Todas las validaciones pasan correctamente

Parte de FASE 1 UI Roadmap
```

**Footer**:
```
Relacionado: #[número-de-issue]
```

**Archivos incluidos**:
- `postgres/tests/test_fase1_structure.sql`
- `postgres/tests/test_fase1_performance.sql`
- `postgres/tests/test_fase1_integrity.sql`

**Comando**:
```bash
git add postgres/tests/test_fase1_*.sql
git commit -m "test(database): agregar tests de validación para FASE 1 UI Database

- Tests de estructura de tablas (columnas, tipos, constraints)
- Tests de performance con 10K inserts en user_activity_log
- Tests de integridad referencial:
  * CASCADE en user_active_context y user_favorites
  * SET NULL en user_activity_log
- Tests de índices y query plans
- Validación de ENUM activity_type
- Validación de JSONB metadata
- Scripts ejecutables en postgres/tests/test_fase1_*.sql

Resultados: Todas las validaciones pasan correctamente

Parte de FASE 1 UI Roadmap

Relacionado: #[número-de-issue]"
```

---

### Commit 5: FASE 6 - Documentación

**Tipo**: `docs(database)`  
**Mensaje corto**: actualizar documentación para FASE 1 UI Database

**Cuerpo**:
```
- Actualizar postgres/README.md con 3 nuevas tablas:
  * user_active_context: propósito, relaciones, queries comunes
  * user_favorites: casos de uso, toggle de favoritos
  * user_activity_log: tipos de actividad, metadata JSONB, consideraciones de escala
- Actualizar CHANGELOG.md con versión postgres/v0.11.0
- Incluir archivos de migración, features, testing
- Documentar consideraciones de particionamiento futuro

Parte de FASE 1 UI Roadmap
```

**Footer**:
```
Relacionado: #[número-de-issue]
```

**Archivos incluidos**:
- `postgres/README.md`
- `CHANGELOG.md`

**Comando**:
```bash
git add postgres/README.md
git add CHANGELOG.md
git commit -m "docs(database): actualizar documentación para FASE 1 UI Database

- Actualizar postgres/README.md con 3 nuevas tablas:
  * user_active_context: propósito, relaciones, queries comunes
  * user_favorites: casos de uso, toggle de favoritos
  * user_activity_log: tipos de actividad, metadata JSONB, consideraciones de escala
- Actualizar CHANGELOG.md con versión postgres/v0.11.0
- Incluir archivos de migración, features, testing
- Documentar consideraciones de particionamiento futuro

Parte de FASE 1 UI Roadmap

Relacionado: #[número-de-issue]"
```

---

### Commit 6 (Opcional): Documentación del plan

**Tipo**: `docs(planning)`  
**Mensaje corto**: agregar documentación de planificación FASE 1

**Cuerpo**:
```
- Agregar docs/specs/fase1-ui-database/ con análisis completo
- README.md: resumen ejecutivo y criterios de aceptación
- ANALISIS-TECNICO.md: análisis profundo de 3 tablas
- Plan-Fase1-UI-Database/: planificación detallada
  * Planner.md: fases y pasos atómicos
  * Planner-commit.md: estrategia de commits
  * Files-affected.md: archivos modificados
  * Test-unit.md: tests a implementar
  * error.md: tracking de errores (template)

Documentación para referencia futura y onboarding

Parte de FASE 1 UI Roadmap
```

**Footer**:
```
Relacionado: #[número-de-issue]
```

**Archivos incluidos**:
- `docs/specs/fase1-ui-database/**`

**Comando**:
```bash
git add docs/specs/fase1-ui-database/
git commit -m "docs(planning): agregar documentación de planificación FASE 1

- Agregar docs/specs/fase1-ui-database/ con análisis completo
- README.md: resumen ejecutivo y criterios de aceptación
- ANALISIS-TECNICO.md: análisis profundo de 3 tablas
- Plan-Fase1-UI-Database/: planificación detallada
  * Planner.md: fases y pasos atómicos
  * Planner-commit.md: estrategia de commits
  * Files-affected.md: archivos modificados
  * Test-unit.md: tests a implementar
  * error.md: tracking de errores (template)

Documentación para referencia futura y onboarding

Parte de FASE 1 UI Roadmap

Relacionado: #[número-de-issue]"
```

---

## Tag de Versión

**Después de todos los commits**, crear tag anotado:

```bash
git tag -a postgres/v0.11.0 -m "Release postgres/v0.11.0 - FASE 1 UI Database

Nuevas tablas para soportar UI Roadmap:
- user_active_context: contexto/escuela activa del usuario
- user_favorites: materiales favoritos
- user_activity_log: log de actividades para analytics

Incluye:
- 6 archivos de migración (structure + constraints)
- Tests de validación completos (estructura, performance, integridad)
- Documentación actualizada (README, CHANGELOG)
- Análisis técnico detallado

Bloquea FASE 2 (APIs) y FASE 4 (UI Estudiantes) del roadmap.

Relacionado: #[número-de-issue]
"
```

---

## Resumen de Commits

| # | Tipo | Alcance | Descripción | Archivos |
|---|------|---------|-------------|----------|
| 1 | feat | database | user_active_context | 2 migraciones |
| 2 | feat | database | user_favorites | 2 migraciones |
| 3 | feat | database | user_activity_log | 2 migraciones |
| 4 | test | database | tests de validación | 3 scripts de test |
| 5 | docs | database | actualizar README y CHANGELOG | 2 archivos |
| 6 | docs | planning | documentación de plan | múltiples |
| TAG | - | - | postgres/v0.11.0 | - |

---

## Convenciones Adicionales

### Mensajes en Español vs Inglés

**Decisión**: Commits en inglés (estándar de la industria)
- Tipo y alcance: inglés
- Descripción corta: inglés
- Cuerpo: puede incluir español para explicaciones detalladas si el equipo lo prefiere

**Razón**: 
- Conventional Commits es estándar internacional
- Facilita integración con herramientas (CHANGELOG automático, etc.)
- Permite colaboración con desarrolladores de otros países

### Longitud de Mensaje

- **Línea de asunto**: Máximo 72 caracteres
- **Cuerpo**: Envuelto a 72 caracteres por línea
- **Listas**: Usar guiones (-) o asteriscos (*)

### Referencias

- **Issues**: `#123`
- **PRs**: `!456` o `PR #456`
- **Otros commits**: `abc1234` (SHA corto)

### Breaking Changes

Si un commit introduce breaking changes (no aplica en este caso):

```
feat(database)!: cambiar schema de user_active_context

BREAKING CHANGE: Renombrar columna school_id a active_school_id
```

---

## Validación de Commits

Antes de cada commit:

```bash
# Verificar qué archivos están staged
git status

# Revisar diff
git diff --staged

# Verificar sintaxis de mensaje (si hay herramienta)
commitlint --edit

# Hacer commit
git commit -m "..."
```

Después de cada commit:

```bash
# Verificar que commit está en log
git log -1 --oneline

# Verificar archivos incluidos
git show --stat
```

---

## Squash vs Commits Atómicos

**Para este proyecto**: Mantener commits atómicos separados

**NO hacer squash** porque:
- Cada fase es independiente y completa
- Facilita revisión granular en PR
- Permite revert específico si hay problemas
- Mejor trazabilidad de cambios

**Cuándo SI hacer squash**:
- Si hay commits de "fix typo" o correcciones menores
- En feature branches con muchos commits WIP
- Al final, si el reviewers lo requiere

---

## Checklist Pre-Commit

Antes de cada commit, verificar:

```
□ Archivos correctos están staged
□ Sintaxis SQL validada (sin errores)
□ Tests pasan (si aplica)
□ Mensaje de commit sigue convención
□ Descripción clara del cambio
□ Referencia a issue incluida
□ NO incluir archivos temporales (.DS_Store, node_modules, etc.)
```

---

## Ejemplo de Historia de Commits Final

```bash
git log --oneline --graph

* a1b2c3d (HEAD -> feature/fase1-ui-database-infrastructure, tag: postgres/v0.11.0) docs(planning): agregar documentación de planificación FASE 1
* e4f5g6h docs(database): actualizar documentación para FASE 1 UI Database
* i7j8k9l test(database): agregar tests de validación para FASE 1 UI Database
* m1n2o3p feat(database): agregar tabla user_activity_log para tracking de actividades
* q4r5s6t feat(database): agregar tabla user_favorites para materiales favoritos
* u7v8w9x feat(database): agregar tabla user_active_context para contexto de usuario
* 9dc9add (origin/main, origin/dev, main, dev) chore: sync dev to main - mock-generator typed entities
```

---

## Herramientas Útiles

### Commitizen (Opcional)

Para ayudar con mensajes de commit:

```bash
npm install -g commitizen
git cz
```

### Commitlint (Opcional)

Para validar mensajes:

```bash
npm install -g @commitlint/cli @commitlint/config-conventional
echo "module.exports = {extends: ['@commitlint/config-conventional']}" > commitlint.config.js
```

---

**Fin de Estrategia de Commits**

Este documento sirve como guía durante la implementación para mantener consistencia en el historial de Git.
