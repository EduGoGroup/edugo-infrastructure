# Gap Analysis: Assessment Schema
# Infrastructure vs Isolated Design

**Fecha:** 17-Nov-2025  
**Proyecto:** edugo-infrastructure (Sprint-01)  
**Objetivo:** Identificar diferencias entre migraciones actuales (006-008) y diseño de `isolated/`

---

## RESUMEN EJECUTIVO

**Migraciones Actuales:** 006, 007, 008 (assessment, assessment_attempt, assessment_attempt_answer)  
**Diseño Isolated:** `edugo-api-mobile/docs/isolated/03-Design/DATA_MODEL.md`

**Resultado:** ⚠️ **DIFERENCIAS SIGNIFICATIVAS** - Se requieren migraciones adicionales (009-011)

---

## TABLA 1: assessment

### Campos Actuales (Infrastructure 006)

```sql
CREATE TABLE assessment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    material_id UUID NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    mongo_document_id VARCHAR(24) NOT NULL UNIQUE,
    questions_count INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'generated' 
        CHECK (status IN ('generated', 'published', 'archived')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

### Campos Esperados (Isolated Design)

```sql
CREATE TABLE assessment (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),              -- ⚠️ DIFERENTE
    material_id UUID NOT NULL,
    mongo_document_id VARCHAR(24) NOT NULL,
    title VARCHAR(255) NOT NULL,                            -- ❌ FALTA
    total_questions INTEGER NOT NULL,                       -- ⚠️ Nombrado como questions_count
    pass_threshold INTEGER NOT NULL DEFAULT 70,             -- ❌ FALTA
    max_attempts INTEGER DEFAULT NULL,                      -- ❌ FALTA
    time_limit_minutes INTEGER DEFAULT NULL,                -- ❌ FALTA
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Análisis de Diferencias

| Campo | Actual | Esperado | Estado | Impacto | Acción |
|-------|--------|----------|--------|---------|--------|
| **id** | `gen_random_uuid()` | `gen_uuid_v7()` | ⚠️ DIFERENTE | BAJO | Mantener actual (retrocompatibilidad) |
| **title** | ❌ No existe | `VARCHAR(255) NOT NULL` | ❌ FALTA | ALTO | AGREGAR con ALTER TABLE |
| **questions_count** | ✅ Existe | Nombrado `total_questions` | ⚠️ NOMBRE | MEDIO | Agregar alias con trigger |
| **total_questions** | ❌ No existe | `INTEGER NOT NULL` | ❌ FALTA | ALTO | AGREGAR y sincronizar |
| **pass_threshold** | ❌ No existe | `INTEGER DEFAULT 70` | ❌ FALTA | ALTO | AGREGAR con ALTER TABLE |
| **max_attempts** | ❌ No existe | `INTEGER DEFAULT NULL` | ❌ FALTA | MEDIO | AGREGAR con ALTER TABLE |
| **time_limit_minutes** | ❌ No existe | `INTEGER DEFAULT NULL` | ❌ FALTA | MEDIO | AGREGAR con ALTER TABLE |
| **status** | `('generated','published','archived')` | `('draft','published','closed')` | ⚠️ VALORES | MEDIO | EXTENDER valores (no reemplazar) |
| **deleted_at** | ✅ Existe | ❌ No especificado | ⚠️ EXTRA | BAJO | Mantener (soft deletes útil) |

### Recomendación: Migración 009

```sql
-- Migration 009: Extend assessment schema

ALTER TABLE assessment
    ADD COLUMN IF NOT EXISTS title VARCHAR(255),
    ADD COLUMN IF NOT EXISTS total_questions INTEGER,
    ADD COLUMN IF NOT EXISTS pass_threshold INTEGER DEFAULT 70 
        CHECK (pass_threshold >= 0 AND pass_threshold <= 100),
    ADD COLUMN IF NOT EXISTS max_attempts INTEGER,
    ADD COLUMN IF NOT EXISTS time_limit_minutes INTEGER;

-- Sincronizar questions_count <-> total_questions
UPDATE assessment SET total_questions = questions_count WHERE total_questions IS NULL;

CREATE TRIGGER trg_sync_questions_count
    BEFORE INSERT OR UPDATE ON assessment
    FOR EACH ROW
    EXECUTE FUNCTION sync_questions_count();

-- Extender status values (sin eliminar existentes)
ALTER TABLE assessment
    DROP CONSTRAINT IF EXISTS assessment_status_check,
    ADD CONSTRAINT assessment_status_check 
        CHECK (status IN ('draft', 'generated', 'published', 'archived', 'closed'));
```

---

## TABLA 2: assessment_attempt

### Campos Actuales (Infrastructure 007)

```sql
CREATE TABLE assessment_attempt (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES assessment(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    score DECIMAL(5,2),
    max_score DECIMAL(5,2),
    percentage DECIMAL(5,2),
    status VARCHAR(50) NOT NULL DEFAULT 'in_progress' 
        CHECK (status IN ('in_progress', 'completed', 'abandoned')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

### Campos Esperados (Isolated Design)

```sql
CREATE TABLE assessment_attempt (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),              -- ⚠️ DIFERENTE
    assessment_id UUID NOT NULL,
    student_id UUID NOT NULL,
    score INTEGER NOT NULL,                                 -- ⚠️ INTEGER vs DECIMAL
    max_score INTEGER NOT NULL DEFAULT 100,                 -- ⚠️ INTEGER vs DECIMAL
    time_spent_seconds INTEGER NOT NULL,                    -- ❌ FALTA
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    idempotency_key VARCHAR(64) DEFAULT NULL,               -- ❌ FALTA
    
    CHECK (completed_at > started_at),                      -- ❌ FALTA
    CHECK (time_spent_seconds = EXTRACT(EPOCH FROM ...))    -- ❌ FALTA
);
```

### Análisis de Diferencias

| Campo | Actual | Esperado | Estado | Impacto | Acción |
|-------|--------|----------|--------|---------|--------|
| **id** | `gen_random_uuid()` | `gen_uuid_v7()` | ⚠️ DIFERENTE | BAJO | Mantener actual |
| **score** | `DECIMAL(5,2)` | `INTEGER` | ⚠️ TIPO | BAJO | Mantener DECIMAL (más flexible) |
| **max_score** | `DECIMAL(5,2)` | `INTEGER` | ⚠️ TIPO | BAJO | Mantener DECIMAL |
| **percentage** | ✅ Existe | ❌ No especificado | ⚠️ EXTRA | BAJO | Mantener (útil para reportes) |
| **time_spent_seconds** | ❌ No existe | `INTEGER NOT NULL` | ❌ FALTA | ALTO | AGREGAR con ALTER TABLE |
| **idempotency_key** | ❌ No existe | `VARCHAR(64) UNIQUE` | ❌ FALTA | MEDIO | AGREGAR con ALTER TABLE |
| **status** | ✅ Existe | ❌ No especificado | ⚠️ EXTRA | BAJO | Mantener (útil para tracking) |
| **CHECK (completed > started)** | ❌ No existe | ✅ Requerido | ❌ FALTA | ALTO | AGREGAR constraint |

### Recomendación: Migración 010

```sql
-- Migration 010: Extend assessment_attempt

ALTER TABLE assessment_attempt
    ADD COLUMN IF NOT EXISTS time_spent_seconds INTEGER 
        CHECK (time_spent_seconds > 0 AND time_spent_seconds <= 7200),
    ADD COLUMN IF NOT EXISTS idempotency_key VARCHAR(64) UNIQUE;

-- Agregar CHECK constraint
ALTER TABLE assessment_attempt
    ADD CONSTRAINT check_attempt_time_logical 
        CHECK (completed_at IS NULL OR completed_at > started_at);

-- Índice parcial para idempotency_key
CREATE INDEX IF NOT EXISTS idx_attempt_idempotency_key 
    ON assessment_attempt(idempotency_key) 
    WHERE idempotency_key IS NOT NULL;
```

---

## TABLA 3: assessment_attempt_answer

### Campos Actuales (Infrastructure 008)

```sql
CREATE TABLE assessment_attempt_answer (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_id UUID NOT NULL REFERENCES assessment_attempt(id) ON DELETE CASCADE,
    question_index INTEGER NOT NULL,
    student_answer TEXT,
    is_correct BOOLEAN,
    points_earned DECIMAL(5,2),
    max_points DECIMAL(5,2),
    answered_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(attempt_id, question_index)
);
```

### Campos Esperados (Isolated Design)

```sql
CREATE TABLE assessment_attempt_answer (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),              -- ⚠️ DIFERENTE
    attempt_id UUID NOT NULL,
    question_id VARCHAR(50) NOT NULL,                       -- ⚠️ NOMBRE DIFERENTE
    selected_answer_id VARCHAR(50) NOT NULL,                -- ⚠️ NOMBRE DIFERENTE
    is_correct BOOLEAN NOT NULL,
    time_spent_seconds INTEGER NOT NULL,                    -- ❌ FALTA
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    UNIQUE(attempt_id, question_id)                         -- ⚠️ question_id vs question_index
);
```

### Análisis de Diferencias

| Campo | Actual | Esperado | Estado | Impacto | Acción |
|-------|--------|----------|--------|---------|--------|
| **id** | `gen_random_uuid()` | `gen_uuid_v7()` | ⚠️ DIFERENTE | BAJO | Mantener actual |
| **question_index** | `INTEGER` | `question_id VARCHAR(50)` | ⚠️ DIFERENTE | MEDIO | Mantener index (más eficiente) |
| **student_answer** | `TEXT` | `selected_answer_id VARCHAR(50)` | ⚠️ DIFERENTE | MEDIO | Mantener TEXT (más flexible) |
| **points_earned** | ✅ Existe | ❌ No especificado | ⚠️ EXTRA | BAJO | Mantener (útil para scoring) |
| **max_points** | ✅ Existe | ❌ No especificado | ⚠️ EXTRA | BAJO | Mantener |
| **time_spent_seconds** | ❌ No existe | `INTEGER NOT NULL` | ❌ FALTA | ALTO | AGREGAR con ALTER TABLE |
| **answered_at** | ✅ Existe | ❌ No especificado | ⚠️ EXTRA | BAJO | Mantener (timestamp útil) |
| **updated_at** | ✅ Existe | ❌ No especificado | ⚠️ EXTRA | BAJO | Mantener |

### Nota Arquitectónica: question_index vs question_id

**Decisión Infrastructure (actual):**
- `question_index INTEGER` - Índice numérico (0-based) de la pregunta en el assessment
- Más eficiente para storage y queries
- Mapeo se hace en capa de aplicación

**Diseño Isolated (esperado):**
- `question_id VARCHAR(50)` - ID de pregunta desde MongoDB
- Más explícito pero menos eficiente

**Recomendación:** **Mantener `question_index`** por eficiencia. Las APIs mapearán `question_index` ↔ `question_id` según sea necesario.

### Recomendación: Migración 011

```sql
-- Migration 011: Extend assessment_attempt_answer

ALTER TABLE assessment_attempt_answer
    ADD COLUMN IF NOT EXISTS time_spent_seconds INTEGER 
        CHECK (time_spent_seconds >= 0);

COMMENT ON COLUMN assessment_attempt_answer.time_spent_seconds IS 
    'Tiempo que tomó responder esta pregunta en segundos';

-- No creamos question_id ni selected_answer_id
-- Mantenemos question_index y student_answer (más flexibles)
-- Mapeo se hace en capa de aplicación
```

---

## RESUMEN DE ACCIONES REQUERIDAS

### Migraciones Nuevas a Crear

| Migración | Nombre | Objetivo | Prioridad |
|-----------|--------|----------|-----------|
| **009** | `extend_assessment_schema` | Agregar title, pass_threshold, max_attempts, time_limit | CRÍTICA |
| **010** | `extend_assessment_attempt` | Agregar time_spent_seconds, idempotency_key, CHECK | CRÍTICA |
| **011** | `extend_assessment_answer` | Agregar time_spent_seconds | ALTA |

### Campos a Agregar

**Tabla `assessment`:**
- ✅ `title VARCHAR(255)`
- ✅ `total_questions INTEGER` (alias de questions_count)
- ✅ `pass_threshold INTEGER DEFAULT 70`
- ✅ `max_attempts INTEGER`
- ✅ `time_limit_minutes INTEGER`

**Tabla `assessment_attempt`:**
- ✅ `time_spent_seconds INTEGER`
- ✅ `idempotency_key VARCHAR(64) UNIQUE`
- ✅ CHECK constraint: `completed_at > started_at`

**Tabla `assessment_attempt_answer`:**
- ✅ `time_spent_seconds INTEGER`

### Decisiones Arquitectónicas

| Decisión | Razón | Impacto |
|----------|-------|---------|
| **Mantener gen_random_uuid()** | Retrocompatibilidad | Bajo (UUIDv7 puede agregarse después) |
| **Mantener DECIMAL para scores** | Mayor precisión para cálculos | Bajo (más flexible que INTEGER) |
| **Mantener question_index** | Eficiencia de storage | Medio (mapeo en aplicación) |
| **Mantener student_answer TEXT** | Flexibilidad para tipos de respuesta | Bajo (soporta JSON, texto libre, etc.) |
| **Agregar campos con ALTER** | Retrocompatibilidad | Alto (no rompe código existente) |
| **Trigger para sincronización** | Transición gradual questions_count ↔ total_questions | Medio (evita migración de datos) |

### Impacto en Retrocompatibilidad

✅ **100% RETROCOMPATIBLE**

- Todos los campos nuevos son opcionales (DEFAULT o NULL)
- `IF NOT EXISTS` en todos los ALTER TABLE
- Sin DROP COLUMN ni RENAME TABLE
- Status values extendidos (no reemplazados)
- Trigger mantiene sincronización gradual

### Próximos Pasos

1. ✅ Crear migraciones 009, 010, 011 (UP y DOWN)
2. ✅ Validar en BD de prueba
3. ✅ Crear seeds de datos con nuevos campos
4. ✅ Documentar en CHANGELOG
5. ✅ PR y release tag

---

## CONCLUSIÓN

**Estado Actual:** Migraciones base (006-008) funcionan pero están **incompletas** según diseño isolated.

**Acción Recomendada:** Crear migraciones **009-011** para extender schema manteniendo **retrocompatibilidad 100%**.

**Impacto:** BAJO riesgo, ALTA ganancia (sincronización con ecosistema).

---

**Fecha:** 17-Nov-2025  
**Responsable:** Claude Code + Jhoan Medina  
**Proyecto:** edugo-infrastructure Sprint-01
