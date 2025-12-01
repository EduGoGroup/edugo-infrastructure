# Análisis Técnico Detallado - FASE 1 UI Database

> **Análisis profundo de las 3 nuevas tablas y sus implicaciones técnicas**

---

## 1. Tabla `user_active_context`

### 1.1 Análisis de Requisitos

**Problema a resolver**:
- Un usuario puede pertenecer a múltiples escuelas (via `memberships`)
- La UI necesita saber qué escuela está "activa" para filtrar datos
- Sin esta tabla, la UI no sabría qué datos mostrar

**Casos de uso**:
1. Usuario con múltiples escuelas selecciona una
2. Al iniciar sesión, se carga la última escuela activa
3. Cambiar de escuela durante la sesión
4. Filtrar materiales, cursos, etc. por escuela activa

**Relaciones**:
```
users (1) ──── (1) user_active_context
schools (1) ──── (N) user_active_context
academic_units (1) ──── (N) user_active_context [opcional]
```

### 1.2 Decisiones de Diseño

#### ¿Por qué UNIQUE en user_id?
- Un usuario solo puede tener UNA escuela activa a la vez
- Evita ambigüedad en la UI
- Simplifica queries: `SELECT * FROM user_active_context WHERE user_id = ?`

#### ¿Por qué unit_id es nullable?
- No todas las escuelas usan unidades académicas (campus, sedes)
- Escuelas pequeñas pueden tener una sola unidad
- Permite flexibilidad sin forzar el modelo

#### ¿Por qué CASCADE en user_id?
- Si un usuario es eliminado, su contexto ya no es válido
- Evita registros huérfanos
- Consistencia de datos

#### ¿Por qué CASCADE en school_id?
- Si una escuela es eliminada, el contexto debe ser eliminado
- El usuario necesitará seleccionar otra escuela
- No tiene sentido mantener contextos de escuelas inexistentes

#### ¿Por qué SET NULL en unit_id?
- Si una unidad es eliminada, el usuario puede seguir en la escuela
- El sistema puede asignar automáticamente otra unidad o dejar NULL
- Menos disruptivo que CASCADE

### 1.3 Índices y Performance

```sql
-- Índice en user_id (CRÍTICO)
CREATE INDEX idx_user_active_context_user ON user_active_context(user_id);
```
**Razón**: 
- Query más frecuente: "¿cuál es la escuela activa del usuario X?"
- Se ejecuta en CADA request de API que necesita filtrar por escuela
- UNIQUE constraint ya crea índice implícito, pero lo hacemos explícito

```sql
-- Índice en school_id
CREATE INDEX idx_user_active_context_school ON user_active_context(school_id);
```
**Razón**:
- Permite queries inversos: "¿qué usuarios tienen esta escuela activa?"
- Útil para analytics y reportes
- Performance en JOINs

### 1.4 Trigger para updated_at

```sql
CREATE TRIGGER set_updated_at_user_active_context
    BEFORE UPDATE ON user_active_context
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

**Razón**:
- Rastrea cuándo el usuario cambió de escuela
- Útil para auditoría
- Analytics: "frecuencia de cambio de contexto"

### 1.5 Queries Esperados

```sql
-- Obtener contexto activo de un usuario (API más frecuente)
SELECT * FROM user_active_context WHERE user_id = $1;

-- Cambiar escuela activa (UPSERT)
INSERT INTO user_active_context (user_id, school_id, unit_id)
VALUES ($1, $2, $3)
ON CONFLICT (user_id) 
DO UPDATE SET 
    school_id = EXCLUDED.school_id,
    unit_id = EXCLUDED.unit_id,
    updated_at = NOW();

-- Obtener usuarios activos en una escuela
SELECT u.* 
FROM users u
JOIN user_active_context uac ON u.id = uac.user_id
WHERE uac.school_id = $1;
```

### 1.6 Validaciones en API (NO en BD)

Las siguientes validaciones deben hacerse en la capa de API:

1. **Usuario tiene membership en la escuela**:
   ```sql
   SELECT 1 FROM memberships 
   WHERE user_id = $1 AND school_id = $2 AND is_active = true;
   ```

2. **School exists y está activa**:
   ```sql
   SELECT 1 FROM schools WHERE id = $1 AND is_active = true;
   ```

3. **Unit pertenece a la school** (si unit_id no es NULL):
   ```sql
   SELECT 1 FROM academic_units 
   WHERE id = $1 AND school_id = $2;
   ```

**¿Por qué en API y no en BD?**
- Mensajes de error más claros al usuario
- Lógica de negocio compleja (ej: "¿puede un admin cambiar contexto sin membership?")
- BD solo garantiza integridad referencial

---

## 2. Tabla `user_favorites`

### 2.1 Análisis de Requisitos

**Problema a resolver**:
- Usuarios quieren marcar materiales como favoritos
- Acceso rápido a materiales frecuentemente usados
- Lista personalizada de "lo más importante"

**Casos de uso**:
1. Usuario marca material como favorito
2. Usuario desmarca material de favoritos
3. Listar favoritos del usuario
4. Verificar si un material específico es favorito

**Relaciones**:
```
users (1) ──── (N) user_favorites
materials (1) ──── (N) user_favorites
```

### 2.2 Decisiones de Diseño

#### ¿Por qué UNIQUE(user_id, material_id)?
- Evita duplicados: un usuario no puede marcar el mismo material 2 veces
- Simplifica lógica de toggle: "si existe → eliminar, si no → insertar"
- Index implícito para queries frecuentes

#### ¿Por qué NO incluir school_id?
- Un material puede ser favorito independientemente de la escuela
- Simplifica el modelo
- Si se necesita filtrar por escuela, se hace via JOIN:
  ```sql
  SELECT m.* FROM materials m
  JOIN user_favorites uf ON m.id = uf.material_id
  WHERE uf.user_id = $1 AND m.school_id = $2;
  ```

#### ¿Por qué NO incluir updated_at?
- Un favorito no se "actualiza", solo se crea o elimina
- `created_at` es suficiente para ordenar por "agregado recientemente"

#### ¿Por qué CASCADE en ambos FKs?
- Si usuario es eliminado → favoritos ya no son válidos
- Si material es eliminado → favorito queda huérfano y sin sentido
- Limpieza automática

### 2.3 Índices y Performance

```sql
-- Índice en user_id (CRÍTICO)
CREATE INDEX idx_user_favorites_user ON user_favorites(user_id);
```
**Razón**:
- Query más frecuente: "listar favoritos del usuario X"
- UNIQUE(user_id, material_id) ya crea índice compuesto, pero queremos uno solo en user_id

```sql
-- Índice en material_id
CREATE INDEX idx_user_favorites_material ON user_favorites(material_id);
```
**Razón**:
- Query inverso: "¿cuántos usuarios tienen este material como favorito?"
- Útil para analytics: "materiales más favoriteados"
- Performance en DELETE cuando material es eliminado

```sql
-- Índice en created_at DESC
CREATE INDEX idx_user_favorites_created ON user_favorites(created_at DESC);
```
**Razón**:
- Ordenar favoritos por "agregado recientemente"
- UI puede mostrar: "Favoritos recientes primero"
- DESC para ORDER BY eficiente

### 2.4 Queries Esperados

```sql
-- Listar favoritos del usuario (ordenados por reciente)
SELECT m.* 
FROM materials m
JOIN user_favorites uf ON m.id = uf.material_id
WHERE uf.user_id = $1
ORDER BY uf.created_at DESC
LIMIT 20 OFFSET 0;

-- Agregar favorito (idempotente con ON CONFLICT)
INSERT INTO user_favorites (user_id, material_id)
VALUES ($1, $2)
ON CONFLICT (user_id, material_id) DO NOTHING;

-- Eliminar favorito
DELETE FROM user_favorites 
WHERE user_id = $1 AND material_id = $2;

-- Verificar si es favorito (para mostrar corazón lleno/vacío en UI)
SELECT EXISTS(
    SELECT 1 FROM user_favorites 
    WHERE user_id = $1 AND material_id = $2
) AS is_favorite;

-- Top materiales más favoriteados (analytics)
SELECT m.*, COUNT(uf.id) as favorite_count
FROM materials m
JOIN user_favorites uf ON m.id = uf.material_id
GROUP BY m.id
ORDER BY favorite_count DESC
LIMIT 10;
```

### 2.5 Consideraciones de Escala

**¿Cuántos favoritos puede tener un usuario?**
- Estimación razonable: 10-50 favoritos por usuario activo
- Poder agregador: 1M usuarios × 30 favoritos = 30M registros

**¿Necesitamos particionamiento?**
- NO inicialmente
- Tabla pequeña comparada con `user_activity_log`
- Índices son suficientes para escala media

**¿Considerar soft deletes?**
- NO necesario
- Favoritos son preferencias temporales, no datos históricos
- Hard delete es apropiado

---

## 3. Tabla `user_activity_log`

### 3.1 Análisis de Requisitos

**Problema a resolver**:
- Mostrar "Actividad reciente" en Home
- Rastrear qué hace el usuario para analytics
- Calcular estadísticas (tiempo de estudio, materiales completados, etc.)

**Casos de uso**:
1. Usuario inicia un material → log `material_started`
2. Usuario completa una página → log `material_progress`
3. Usuario termina material → log `material_completed`
4. Usuario ve resumen → log `summary_viewed`
5. Usuario inicia quiz → log `quiz_started`
6. Usuario completa quiz → log `quiz_completed` + `quiz_passed`/`quiz_failed`

**Relaciones**:
```
users (1) ──── (N) user_activity_log
materials (1) ──── (N) user_activity_log [opcional]
schools (1) ──── (N) user_activity_log [opcional]
```

### 3.2 Decisiones de Diseño

#### ¿Por qué usar ENUM para activity_type?

```sql
CREATE TYPE activity_type AS ENUM (
    'material_started',
    'material_progress',
    'material_completed',
    'summary_viewed',
    'quiz_started',
    'quiz_completed',
    'quiz_passed',
    'quiz_failed'
);
```

**Ventajas**:
- Valores limitados y conocidos → previene typos
- Más eficiente en storage que VARCHAR
- Validación automática a nivel de BD
- Facilita agregaciones: `GROUP BY activity_type`

**Desventajas**:
- Agregar nuevos valores requiere `ALTER TYPE` (migración)
- No es tan flexible como VARCHAR

**Decisión**: ENUM es apropiado porque:
- Los tipos de actividad son estables (no cambian frecuentemente)
- Agregación y analytics son prioridad
- Mejor performance

#### ¿Por qué material_id y school_id son nullable?

- **material_id NULL**: Actividades que no están ligadas a un material específico
  - Ejemplo futuro: `settings_changed`, `profile_updated`
  
- **school_id NULL**: Actividades globales del usuario
  - Ejemplo futuro: `logged_in`, `logged_out`

**Decisión**: Nullable permite flexibilidad sin complicar el modelo

#### ¿Por qué JSONB para metadata?

```sql
metadata JSONB DEFAULT '{}'
```

**Casos de uso**:
```json
// Para material_progress
{
  "page": 5,
  "total_pages": 10,
  "time_spent_seconds": 120
}

// Para quiz_completed
{
  "score": 85,
  "total_questions": 20,
  "correct_answers": 17,
  "time_taken_seconds": 300
}

// Para summary_viewed
{
  "summary_length_chars": 500,
  "read_time_seconds": 45
}
```

**Ventajas**:
- Flexibilidad: cada activity_type puede tener metadata diferente
- No necesitamos columnas para cada atributo posible
- Queries JSONB en PostgreSQL son eficientes
- Evita crear tablas separadas para cada tipo

**Desventajas**:
- Validación de estructura debe hacerse en API
- Queries más complejas

**Decisión**: JSONB es ideal para metadata heterogénea

#### ¿Por qué SET NULL en FKs?

```sql
material_id UUID REFERENCES materials(id) ON DELETE SET NULL
school_id UUID REFERENCES schools(id) ON DELETE SET NULL
```

**Razón**:
- Logs son **históricos** → no deben eliminarse
- Si material/school es eliminado, el log sigue siendo valioso
- Analytics: "usuario completó un material (ahora eliminado) el 15/12/2024"
- Mantiene integridad de datos históricos

**Alternativa considerada**: CASCADE
- **Rechazada** porque perderíamos datos históricos valiosos

### 3.3 Índices y Performance

Esta tabla crecerá **RÁPIDO**. Estimación:
- 100K usuarios activos
- 10 actividades/día por usuario
- = 1M registros/día = 365M registros/año

**Índices críticos**:

```sql
-- Query más frecuente: "actividad reciente del usuario"
CREATE INDEX idx_user_activity_user_created 
ON user_activity_log(user_id, created_at DESC);
```
**Uso**:
```sql
SELECT * FROM user_activity_log
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT 10;
```

```sql
-- Analytics por escuela
CREATE INDEX idx_user_activity_school 
ON user_activity_log(school_id, created_at DESC);
```
**Uso**:
```sql
SELECT * FROM user_activity_log
WHERE school_id = $1 AND created_at > NOW() - INTERVAL '30 days'
ORDER BY created_at DESC;
```

```sql
-- Filtrar por tipo de actividad
CREATE INDEX idx_user_activity_type 
ON user_activity_log(activity_type);
```
**Uso**:
```sql
SELECT COUNT(*) FROM user_activity_log
WHERE activity_type = 'material_completed'
  AND created_at > NOW() - INTERVAL '7 days';
```

```sql
-- Índice parcial para rate limiting
CREATE INDEX idx_user_activity_rate_limit
ON user_activity_log(user_id, activity_type, created_at)
WHERE created_at > NOW() - INTERVAL '1 hour';
```
**Uso**:
- Detectar usuarios que generan spam
- Rate limiting: "máximo 100 material_progress por hora"
- Índice parcial → más pequeño → más rápido

### 3.4 Particionamiento (FUTURO)

Con 365M registros/año, **necesitaremos particionamiento**.

**Estrategia sugerida**: Particionamiento por rango de fecha (mensual)

```sql
-- Ejemplo (NO implementar ahora, solo referencia futura)
CREATE TABLE user_activity_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    activity_type activity_type NOT NULL,
    material_id UUID REFERENCES materials(id) ON DELETE SET NULL,
    school_id UUID REFERENCES schools(id) ON DELETE SET NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- Crear partición por mes
CREATE TABLE user_activity_log_y2025m12 PARTITION OF user_activity_log
    FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');

CREATE TABLE user_activity_log_y2026m01 PARTITION OF user_activity_log
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
```

**Beneficios**:
- Queries más rápidos (solo escanea partición relevante)
- Fácil eliminar datos viejos: `DROP TABLE user_activity_log_y2023m01`
- Mejor performance en INSERT

**Cuándo implementar**:
- Cuando alcancemos ~50M registros
- O cuando queries empiecen a degradarse
- Requiere migración compleja (recrear tabla)

### 3.5 Queries Esperados

```sql
-- Actividad reciente del usuario (Home Screen)
SELECT 
    id,
    activity_type,
    material_id,
    m.title as material_title,
    metadata,
    created_at
FROM user_activity_log ual
LEFT JOIN materials m ON ual.material_id = m.id
WHERE ual.user_id = $1
ORDER BY created_at DESC
LIMIT 10;

-- Estadísticas del usuario (últimos 30 días)
SELECT 
    activity_type,
    COUNT(*) as count,
    DATE_TRUNC('day', created_at) as date
FROM user_activity_log
WHERE user_id = $1 
  AND created_at > NOW() - INTERVAL '30 days'
GROUP BY activity_type, DATE_TRUNC('day', created_at)
ORDER BY date DESC;

-- Materiales más activos (analytics)
SELECT 
    m.id,
    m.title,
    COUNT(*) as activity_count,
    COUNT(DISTINCT user_id) as unique_users
FROM user_activity_log ual
JOIN materials m ON ual.material_id = m.id
WHERE ual.school_id = $1
  AND ual.created_at > NOW() - INTERVAL '7 days'
GROUP BY m.id, m.title
ORDER BY activity_count DESC
LIMIT 10;

-- Insertar actividad (desde API)
INSERT INTO user_activity_log (
    user_id, 
    activity_type, 
    material_id, 
    school_id, 
    metadata
) VALUES ($1, $2, $3, $4, $5);
```

### 3.6 Consideraciones de Privacidad

**Pregunta**: ¿Estos logs contienen información sensible?

**Respuesta**: Depende del metadata

**Datos seguros**:
- `activity_type`: público
- `material_id`: referencia a material (público dentro de escuela)
- `created_at`: timestamp de actividad

**Datos potencialmente sensibles en metadata**:
- Scores de quizzes → pueden revelar rendimiento académico
- Tiempo de lectura → patrones de comportamiento
- IPs o device info (si se almacenan)

**Recomendaciones**:
1. **NO almacenar PII** en metadata (nombres, emails, etc.)
2. **NO almacenar IPs** (usar tabla separada si se necesita)
3. **Considerar anonimización** después de X meses
4. **GDPR**: Permitir eliminación de logs al eliminar usuario (CASCADE en user_id)

### 3.7 Archivado y Retención

**Pregunta**: ¿Cuánto tiempo mantener logs?

**Sugerencias**:
- **Logs recientes** (0-3 meses): Tabla activa, acceso frecuente
- **Logs históricos** (3-12 meses): Partición archivada, acceso ocasional
- **Logs antiguos** (12+ meses): Exportar a data warehouse, eliminar de BD transaccional

**Estrategia de archivado** (futuro):
```sql
-- Mover logs viejos a tabla de archivo (trimestral)
CREATE TABLE user_activity_log_archive_2025q1 AS
SELECT * FROM user_activity_log
WHERE created_at >= '2025-01-01' AND created_at < '2025-04-01';

-- Eliminar de tabla activa
DELETE FROM user_activity_log
WHERE created_at >= '2025-01-01' AND created_at < '2025-04-01';
```

---

## 4. Impacto en Performance del Sistema

### 4.1 Writes por Segundo (WPS)

**Estimación de carga**:

| Tabla | Estimación WPS | Pico WPS | Notas |
|-------|----------------|----------|-------|
| `user_active_context` | 0.1 | 10 | Solo al cambiar escuela |
| `user_favorites` | 1 | 50 | Toggle favoritos ocasional |
| `user_activity_log` | 100 | 1000 | ALTA frecuencia, cada acción |

**user_activity_log** es la tabla crítica para performance.

### 4.2 Estrategias de Optimización

#### Opción 1: Batch Inserts desde API
```go
// En vez de 1 INSERT por actividad
for _, activity := range activities {
    db.Insert(activity) // N queries
}

// Hacer batch insert
db.BatchInsert(activities) // 1 query
```

#### Opción 2: Async Logging
```go
// API no espera confirmación de INSERT
activityLogger.LogAsync(activity)

// Background worker procesa cola
go func() {
    for activity := range activityQueue {
        db.Insert(activity)
    }
}()
```

#### Opción 3: Buffer en Memoria + Flush Periódico
```go
// Acumular en memoria
activityBuffer.Add(activity)

// Flush cada 10 segundos o cada 100 items
if buffer.ShouldFlush() {
    db.BatchInsert(buffer.Drain())
}
```

**Recomendación**: Opción 2 (Async) + Opción 3 (Buffering)
- API responde rápido
- BD no se satura
- Pérdida mínima si falla (solo buffer en memoria)

### 4.3 Monitoring y Alertas

**Métricas a monitorear**:

```sql
-- Tamaño de tabla user_activity_log
SELECT pg_size_pretty(pg_total_relation_size('user_activity_log'));

-- Tasa de crecimiento
SELECT 
    DATE_TRUNC('hour', created_at) as hour,
    COUNT(*) as records_per_hour
FROM user_activity_log
WHERE created_at > NOW() - INTERVAL '24 hours'
GROUP BY hour
ORDER BY hour DESC;

-- Índices más usados
SELECT 
    schemaname, 
    tablename, 
    indexname, 
    idx_scan as scans,
    idx_tup_read as tuples_read
FROM pg_stat_user_indexes
WHERE tablename = 'user_activity_log'
ORDER BY idx_scan DESC;

-- Índices no usados (considerar eliminar)
SELECT 
    schemaname, 
    tablename, 
    indexname
FROM pg_stat_user_indexes
WHERE tablename = 'user_activity_log' 
  AND idx_scan = 0
  AND indexrelname NOT LIKE '%_pkey';
```

**Alertas sugeridas**:
- ⚠️ Si tabla > 100GB → considerar particionamiento
- ⚠️ Si inserts/sec > 1000 → considerar buffering
- ⚠️ Si query latency > 100ms → revisar índices

---

## 5. Plan de Rollback

### 5.1 ¿Qué pasa si necesitamos revertir?

**Escenario 1**: Migración falla a mitad
- PostgreSQL usa transacciones en DDL
- Si falla, se revierte automáticamente
- **Acción**: Revisar logs, corregir SQL, reintentar

**Escenario 2**: Migración exitosa pero APIs tienen bugs
- Tablas ya están creadas
- APIs usan tablas nuevas y fallan
- **Acción**: Revertir código de API (no eliminar tablas)
- **Razón**: Eliminar tablas pierde datos si ya hay registros

**Escenario 3**: Necesitamos eliminar tablas en producción
```sql
-- Reversal de 013_create_user_activity_log.sql
DROP INDEX IF EXISTS idx_user_activity_rate_limit;
DROP INDEX IF EXISTS idx_user_activity_type;
DROP INDEX IF EXISTS idx_user_activity_school;
DROP INDEX IF EXISTS idx_user_activity_user_created;
DROP TRIGGER IF EXISTS set_updated_at_user_activity_log ON user_activity_log;
DROP TABLE IF EXISTS user_activity_log CASCADE;
DROP TYPE IF EXISTS activity_type CASCADE;

-- Reversal de 012_create_user_favorites.sql
DROP INDEX IF EXISTS idx_user_favorites_created;
DROP INDEX IF EXISTS idx_user_favorites_material;
DROP INDEX IF EXISTS idx_user_favorites_user;
DROP TABLE IF EXISTS user_favorites CASCADE;

-- Reversal de 011_create_user_active_context.sql
DROP INDEX IF EXISTS idx_user_active_context_school;
DROP INDEX IF EXISTS idx_user_active_context_user;
DROP TRIGGER IF EXISTS set_updated_at_user_active_context ON user_active_context;
DROP TABLE IF EXISTS user_active_context CASCADE;
```

**⚠️ ADVERTENCIA**: Solo ejecutar en desarrollo/staging
- En producción, coordinar con equipo
- Hacer backup antes: `pg_dump edugo_db > backup_before_rollback.sql`

### 5.2 Migraciones Down (Opcional)

Podemos crear migraciones "down" para reversar:

```
postgres/migrations/
├── structure/
│   ├── 011_create_user_active_context.sql       (UP)
│   ├── 011_create_user_active_context_down.sql  (DOWN)
│   ├── 012_create_user_favorites.sql
│   ├── 012_create_user_favorites_down.sql
│   ├── 013_create_user_activity_log.sql
│   └── 013_create_user_activity_log_down.sql
```

**Debate**: ¿Vale la pena?
- ✅ **Pro**: Fácil revertir en dev/staging
- ❌ **Contra**: Doble mantenimiento, nunca usamos en prod
- **Decisión**: Crear solo si tenemos herramienta de migración que lo soporte

---

## 6. Testing de Migraciones

### 6.1 Tests a Ejecutar

**Test 1: Migración ejecuta sin errores**
```bash
psql -U postgres -d edugo_db_test < 011_create_user_active_context.sql
psql -U postgres -d edugo_db_test < 012_create_user_favorites.sql
psql -U postgres -d edugo_db_test < 013_create_user_activity_log.sql

# Verificar sin errores
echo $?  # Debe ser 0
```

**Test 2: Estructura de tablas correcta**
```sql
\d user_active_context
\d user_favorites
\d user_activity_log
```

Verificar:
- Columnas correctas
- Tipos de datos correctos
- Constraints (PK, FK, UNIQUE)
- Defaults

**Test 3: Índices creados**
```sql
\di idx_user_active_context_user
\di idx_user_active_context_school
\di idx_user_favorites_user
\di idx_user_favorites_material
\di idx_user_favorites_created
\di idx_user_activity_user_created
\di idx_user_activity_school
\di idx_user_activity_type
\di idx_user_activity_rate_limit
```

**Test 4: Triggers funcionan**
```sql
-- Insertar registro en user_active_context
INSERT INTO user_active_context (user_id, school_id)
VALUES (
    (SELECT id FROM users LIMIT 1),
    (SELECT id FROM schools LIMIT 1)
);

-- Verificar updated_at
SELECT updated_at FROM user_active_context ORDER BY created_at DESC LIMIT 1;

-- Actualizar registro
UPDATE user_active_context 
SET school_id = (SELECT id FROM schools LIMIT 1 OFFSET 1)
WHERE id = (SELECT id FROM user_active_context ORDER BY created_at DESC LIMIT 1);

-- Verificar que updated_at cambió
SELECT updated_at, created_at FROM user_active_context 
ORDER BY created_at DESC LIMIT 1;
-- updated_at debe ser mayor que created_at
```

**Test 5: Constraints funcionan**
```sql
-- Test UNIQUE constraint en user_active_context
INSERT INTO user_active_context (user_id, school_id)
VALUES (
    (SELECT id FROM users LIMIT 1),
    (SELECT id FROM schools LIMIT 1)
);
-- Debe fallar con: duplicate key value violates unique constraint

-- Test UNIQUE constraint en user_favorites
INSERT INTO user_favorites (user_id, material_id)
VALUES (
    (SELECT id FROM users LIMIT 1),
    (SELECT id FROM materials LIMIT 1)
);
INSERT INTO user_favorites (user_id, material_id)
VALUES (
    (SELECT id FROM users LIMIT 1),
    (SELECT id FROM materials LIMIT 1)
);
-- Segundo debe fallar

-- Test FK constraints
INSERT INTO user_active_context (user_id, school_id)
VALUES ('00000000-0000-0000-0000-000000000000', 
        '00000000-0000-0000-0000-000000000000');
-- Debe fallar con: violates foreign key constraint
```

**Test 6: CASCADE y SET NULL funcionan**
```sql
-- Crear datos de prueba
INSERT INTO user_active_context (user_id, school_id)
VALUES (
    (SELECT id FROM users WHERE email = 'student1@edugo.test'),
    (SELECT id FROM schools LIMIT 1)
) RETURNING id;

-- Eliminar usuario
DELETE FROM users WHERE email = 'student1@edugo.test';

-- Verificar que user_active_context fue eliminado (CASCADE)
SELECT COUNT(*) FROM user_active_context 
WHERE user_id = (SELECT id FROM users WHERE email = 'student1@edugo.test');
-- Debe retornar 0

-- Test SET NULL en user_activity_log
INSERT INTO user_activity_log (user_id, activity_type, material_id, school_id)
VALUES (
    (SELECT id FROM users WHERE email = 'student2@edugo.test'),
    'material_started',
    (SELECT id FROM materials LIMIT 1),
    (SELECT id FROM schools LIMIT 1)
);

-- Eliminar material
DELETE FROM materials WHERE id = (SELECT material_id FROM user_activity_log LIMIT 1);

-- Verificar que material_id ahora es NULL
SELECT material_id FROM user_activity_log 
ORDER BY created_at DESC LIMIT 1;
-- material_id debe ser NULL, registro debe existir
```

### 6.2 Tests de Performance

```sql
-- Test 1: Insertar 10K registros en user_activity_log
EXPLAIN ANALYZE
INSERT INTO user_activity_log (user_id, activity_type, material_id, school_id, metadata)
SELECT 
    (SELECT id FROM users LIMIT 1),
    'material_progress',
    (SELECT id FROM materials LIMIT 1),
    (SELECT id FROM schools LIMIT 1),
    '{"page": 5}'::jsonb
FROM generate_series(1, 10000);

-- Verificar tiempo < 1 segundo

-- Test 2: Query de actividad reciente
EXPLAIN ANALYZE
SELECT * FROM user_activity_log
WHERE user_id = (SELECT id FROM users LIMIT 1)
ORDER BY created_at DESC
LIMIT 10;

-- Verificar que usa idx_user_activity_user_created
-- Tiempo < 10ms

-- Test 3: Aggregation query
EXPLAIN ANALYZE
SELECT activity_type, COUNT(*)
FROM user_activity_log
WHERE user_id = (SELECT id FROM users LIMIT 1)
  AND created_at > NOW() - INTERVAL '30 days'
GROUP BY activity_type;

-- Verificar tiempo < 50ms
```

---

## 7. Checklist Final de Validación

```
Migración 011 (user_active_context):
□ Archivo SQL sin errores de sintaxis
□ Tabla creada correctamente
□ Constraint UNIQUE(user_id) existe
□ FKs a users, schools, academic_units configurados
□ Índices idx_user_active_context_user y school creados
□ Trigger set_updated_at_user_active_context funciona
□ CASCADE en user_id y school_id funcionan
□ SET NULL en unit_id funciona
□ Comentarios COMMENT ON TABLE/COLUMN agregados

Migración 012 (user_favorites):
□ Archivo SQL sin errores de sintaxis
□ Tabla creada correctamente
□ Constraint UNIQUE(user_id, material_id) existe
□ FKs a users, materials configurados
□ Índices en user_id, material_id, created_at creados
□ CASCADE en ambos FKs funciona
□ Test de inserción duplicada falla correctamente

Migración 013 (user_activity_log):
□ Archivo SQL sin errores de sintaxis
□ ENUM activity_type creado correctamente
□ Tabla creada correctamente
□ FKs a users, materials, schools configurados
□ Índices idx_user_activity_* creados (4 total)
□ SET NULL en material_id y school_id funciona
□ CASCADE en user_id funciona
□ JSONB metadata acepta JSON válido
□ Performance test con 10K inserts < 1 segundo

Documentación:
□ postgres/README.md actualizado con nuevas tablas
□ CHANGELOG.md con nueva versión
□ Comentarios en código SQL explicando decisiones
□ Este documento (ANALISIS-TECNICO.md) completo

Testing:
□ Migraciones ejecutadas en local sin errores
□ Migraciones ejecutadas en dev sin errores
□ Tests de constraints pasados
□ Tests de performance pasados
□ Queries esperados validados
```

---

## 8. Recursos y Referencias

### Documentación PostgreSQL
- [CREATE TABLE](https://www.postgresql.org/docs/current/sql-createtable.html)
- [CREATE INDEX](https://www.postgresql.org/docs/current/sql-createindex.html)
- [Table Partitioning](https://www.postgresql.org/docs/current/ddl-partitioning.html)
- [JSONB](https://www.postgresql.org/docs/current/datatype-json.html)
- [ENUM Types](https://www.postgresql.org/docs/current/datatype-enum.html)

### Best Practices
- [Indexing Best Practices](https://www.postgresql.org/docs/current/indexes-types.html)
- [Performance Tips](https://wiki.postgresql.org/wiki/Performance_Optimization)
- [Partitioning Strategies](https://www.postgresql.org/docs/current/ddl-partitioning.html#DDL-PARTITIONING-DECLARATIVE)

### Internal References
- Existing migrations: `postgres/migrations/`
- Migration conventions: `postgres/README.md`
- Schema documentation: `schemas/`

---

**Fin del Análisis Técnico**

Este documento debe servir como referencia durante la implementación y para futuras decisiones de arquitectura de BD.
