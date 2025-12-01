# Planner - FASE 1 UI Database

> **Plan detallado de fases y pasos atómicos**

---

## FASE 1: Preparación y Setup (15 min)

### Objetivo
Preparar el entorno y validar que tenemos todo lo necesario para implementar las migraciones.

### Paso 1.1: Verificar estructura de migraciones existentes
**Acciones**:
- [ ] Listar archivos en `postgres/migrations/structure/`
- [ ] Verificar última migración (debe ser 010_create_login_attempts.sql)
- [ ] Listar archivos en `postgres/migrations/constraints/`
- [ ] Verificar convenciones de nomenclatura

**Comando**:
```bash
ls -la postgres/migrations/structure/
ls -la postgres/migrations/constraints/
```

**Resultado esperado**: Confirmar que la última migración es 010

---

### Paso 1.2: Verificar función de trigger update_updated_at_column
**Acciones**:
- [ ] Buscar dónde está definida la función `update_updated_at_column()`
- [ ] Verificar que existe en migraciones previas
- [ ] Confirmar que podemos reutilizarla

**Comando**:
```bash
grep -r "update_updated_at_column" postgres/migrations/
```

**Resultado esperado**: Encontrar la función definida (probablemente en 000_initial_setup.sql o similar)

---

### Paso 1.3: Verificar conexión a base de datos local
**Acciones**:
- [ ] Conectarse a PostgreSQL local
- [ ] Verificar que la BD edugo_db existe
- [ ] Confirmar que migraciones actuales están aplicadas

**Comando**:
```bash
psql -U postgres -l | grep edugo
psql -U postgres -d edugo_db -c "\dt"
```

**Resultado esperado**: BD existe y tiene tablas actuales (users, schools, materials, etc.)

---

**COMMIT FASE 1**: No aplica (solo validación)

---

## FASE 2: Crear migración user_active_context (30 min)

### Objetivo
Implementar la tabla `user_active_context` con estructura completa, constraints, índices y trigger.

### Paso 2.1: Crear archivo de estructura
**Acciones**:
- [ ] Crear archivo `postgres/migrations/structure/011_create_user_active_context.sql`
- [ ] Implementar CREATE TABLE con todas las columnas
- [ ] Agregar constraints (UNIQUE, FKs)
- [ ] Crear índices (idx_user_active_context_user, idx_user_active_context_school)
- [ ] Agregar trigger para updated_at
- [ ] Incluir comentarios (COMMENT ON TABLE/COLUMN)

**Contenido**:
```sql
-- Tabla para almacenar el contexto/escuela activa del usuario
-- Permite filtrar datos en UI según la escuela seleccionada

CREATE TABLE IF NOT EXISTS user_active_context (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    school_id UUID NOT NULL,
    unit_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT uq_user_active_context_user UNIQUE(user_id),
    CONSTRAINT fk_user_active_context_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_user_active_context_school
        FOREIGN KEY (school_id)
        REFERENCES schools(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_user_active_context_unit
        FOREIGN KEY (unit_id)
        REFERENCES academic_units(id)
        ON DELETE SET NULL
);

-- Índices para performance
CREATE INDEX idx_user_active_context_user ON user_active_context(user_id);
CREATE INDEX idx_user_active_context_school ON user_active_context(school_id);

-- Trigger para updated_at automático
CREATE TRIGGER set_updated_at_user_active_context
    BEFORE UPDATE ON user_active_context
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comentarios
COMMENT ON TABLE user_active_context IS 'Almacena el contexto/escuela activa del usuario para filtrar datos en UI';
COMMENT ON COLUMN user_active_context.user_id IS 'Usuario propietario del contexto (UNIQUE: solo un contexto por usuario)';
COMMENT ON COLUMN user_active_context.school_id IS 'Escuela actualmente seleccionada por el usuario';
COMMENT ON COLUMN user_active_context.unit_id IS 'Unidad académica activa (opcional, puede ser NULL)';
```

**Validación**:
- [ ] Sintaxis SQL correcta
- [ ] Nombres de tablas/columnas consistentes con convención
- [ ] FKs apuntan a tablas correctas

---

### Paso 2.2: Crear archivo de constraints (si aplica)
**Acciones**:
- [ ] Crear archivo `postgres/migrations/constraints/011_create_user_active_context.sql`
- [ ] Evaluar si hay constraints adicionales (si todos están en structure, este archivo puede ser vacío)

**Contenido**:
```sql
-- No additional constraints needed (already in structure file)
```

**Nota**: Seguir convención del proyecto (si 009 y 010 tienen archivos vacíos, mantener consistencia)

---

### Paso 2.3: Ejecutar migración en local
**Acciones**:
- [ ] Conectarse a BD local
- [ ] Ejecutar archivo de estructura
- [ ] Ejecutar archivo de constraints (si aplica)
- [ ] Verificar con `\d user_active_context`
- [ ] Verificar índices con `\di`

**Comando**:
```bash
psql -U postgres -d edugo_db -f postgres/migrations/structure/011_create_user_active_context.sql
psql -U postgres -d edugo_db -f postgres/migrations/constraints/011_create_user_active_context.sql

psql -U postgres -d edugo_db -c "\d user_active_context"
psql -U postgres -d edugo_db -c "\di idx_user_active_context*"
```

**Resultado esperado**: Tabla creada sin errores, estructura correcta

---

### Paso 2.4: Probar constraints y trigger
**Acciones**:
- [ ] Insertar registro de prueba
- [ ] Verificar UNIQUE constraint (intentar duplicado)
- [ ] Verificar FK constraints
- [ ] Verificar que trigger actualiza updated_at

**Comandos de prueba**:
```sql
-- Insertar válido
INSERT INTO user_active_context (user_id, school_id)
VALUES (
    (SELECT id FROM users LIMIT 1),
    (SELECT id FROM schools LIMIT 1)
);

-- Intentar duplicado (debe fallar)
INSERT INTO user_active_context (user_id, school_id)
VALUES (
    (SELECT id FROM users LIMIT 1),
    (SELECT id FROM schools LIMIT 1)
);

-- Verificar trigger
UPDATE user_active_context 
SET school_id = (SELECT id FROM schools LIMIT 1 OFFSET 1)
WHERE user_id = (SELECT id FROM users LIMIT 1);

SELECT created_at, updated_at FROM user_active_context 
WHERE user_id = (SELECT id FROM users LIMIT 1);
-- updated_at debe ser > created_at

-- Limpiar datos de prueba
DELETE FROM user_active_context;
```

**Resultado esperado**: Todas las validaciones funcionan correctamente

---

**COMMIT FASE 2**: 
```
feat(database): agregar tabla user_active_context para contexto de usuario

- Crear migración 011_create_user_active_context.sql
- Tabla almacena escuela activa del usuario para filtrado en UI
- Incluye constraints UNIQUE, FKs con CASCADE/SET NULL
- Índices en user_id y school_id para performance
- Trigger automático para updated_at

Parte de FASE 1 UI Roadmap - Bloquea APIs y UI

Relacionado: #[número-de-issue]
```

**ID del commit**: [Se llenará después de hacer el commit]

---

## FASE 3: Crear migración user_favorites (25 min)

### Objetivo
Implementar la tabla `user_favorites` para materiales favoritos del usuario.

### Paso 3.1: Crear archivo de estructura
**Acciones**:
- [ ] Crear archivo `postgres/migrations/structure/012_create_user_favorites.sql`
- [ ] Implementar CREATE TABLE con columnas necesarias
- [ ] Agregar constraints (UNIQUE compuesto, FKs)
- [ ] Crear índices (user_id, material_id, created_at)
- [ ] Incluir comentarios

**Contenido**:
```sql
-- Tabla para almacenar materiales marcados como favoritos por usuarios
-- Permite acceso rápido a contenido frecuentemente usado

CREATE TABLE IF NOT EXISTS user_favorites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    material_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT uq_user_favorites_user_material UNIQUE(user_id, material_id),
    CONSTRAINT fk_user_favorites_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_user_favorites_material
        FOREIGN KEY (material_id)
        REFERENCES materials(id)
        ON DELETE CASCADE
);

-- Índices para queries frecuentes
CREATE INDEX idx_user_favorites_user ON user_favorites(user_id);
CREATE INDEX idx_user_favorites_material ON user_favorites(material_id);
CREATE INDEX idx_user_favorites_created ON user_favorites(created_at DESC);

-- Comentarios
COMMENT ON TABLE user_favorites IS 'Materiales marcados como favoritos por usuarios';
COMMENT ON COLUMN user_favorites.user_id IS 'Usuario que marcó el favorito';
COMMENT ON COLUMN user_favorites.material_id IS 'Material marcado como favorito';
COMMENT ON COLUMN user_favorites.created_at IS 'Fecha cuando fue agregado a favoritos (para ordenar por reciente)';
```

**Validación**:
- [ ] Sintaxis SQL correcta
- [ ] UNIQUE constraint compuesto correcto
- [ ] Índices apropiados

---

### Paso 3.2: Crear archivo de constraints
**Acciones**:
- [ ] Crear archivo `postgres/migrations/constraints/012_create_user_favorites.sql`

**Contenido**:
```sql
-- No additional constraints needed (already in structure file)
```

---

### Paso 3.3: Ejecutar migración en local
**Acciones**:
- [ ] Ejecutar archivo de estructura
- [ ] Ejecutar archivo de constraints
- [ ] Verificar estructura con `\d user_favorites`
- [ ] Verificar índices

**Comando**:
```bash
psql -U postgres -d edugo_db -f postgres/migrations/structure/012_create_user_favorites.sql
psql -U postgres -d edugo_db -f postgres/migrations/constraints/012_create_user_favorites.sql

psql -U postgres -d edugo_db -c "\d user_favorites"
psql -U postgres -d edugo_db -c "\di idx_user_favorites*"
```

**Resultado esperado**: Tabla creada sin errores

---

### Paso 3.4: Probar constraints
**Acciones**:
- [ ] Insertar favorito válido
- [ ] Intentar duplicado (debe fallar)
- [ ] Verificar CASCADE en FKs

**Comandos de prueba**:
```sql
-- Insertar favorito
INSERT INTO user_favorites (user_id, material_id)
VALUES (
    (SELECT id FROM users LIMIT 1),
    (SELECT id FROM materials LIMIT 1)
);

-- Intentar duplicado (debe fallar)
INSERT INTO user_favorites (user_id, material_id)
VALUES (
    (SELECT id FROM users LIMIT 1),
    (SELECT id FROM materials LIMIT 1)
);

-- Query de favoritos del usuario
SELECT uf.*, m.title 
FROM user_favorites uf
JOIN materials m ON uf.material_id = m.id
WHERE uf.user_id = (SELECT id FROM users LIMIT 1)
ORDER BY uf.created_at DESC;

-- Limpiar
DELETE FROM user_favorites;
```

**Resultado esperado**: Constraints funcionan correctamente

---

**COMMIT FASE 3**:
```
feat(database): agregar tabla user_favorites para materiales favoritos

- Crear migración 012_create_user_favorites.sql
- Tabla almacena materiales marcados como favoritos por usuarios
- UNIQUE(user_id, material_id) evita duplicados
- Índices en user_id, material_id y created_at
- CASCADE en ambos FKs para limpieza automática

Parte de FASE 1 UI Roadmap - Funcionalidad de favoritos en UI

Relacionado: #[número-de-issue]
```

**ID del commit**: [Se llenará después de hacer el commit]

---

## FASE 4: Crear migración user_activity_log (35 min)

### Objetivo
Implementar la tabla `user_activity_log` para rastrear actividades del usuario.

### Paso 4.1: Crear ENUM activity_type
**Acciones**:
- [ ] Crear archivo `postgres/migrations/structure/013_create_user_activity_log.sql`
- [ ] Definir ENUM con tipos de actividad
- [ ] Implementar CREATE TABLE

**Contenido (Parte 1 - ENUM)**:
```sql
-- Tipo ENUM para clasificar actividades del usuario
CREATE TYPE activity_type AS ENUM (
    'material_started',      -- Usuario inició un material
    'material_progress',     -- Usuario avanzó en lectura
    'material_completed',    -- Usuario completó material
    'summary_viewed',        -- Usuario vio resumen generado
    'quiz_started',          -- Usuario inició quiz
    'quiz_completed',        -- Usuario completó quiz
    'quiz_passed',          -- Usuario aprobó quiz
    'quiz_failed'           -- Usuario reprobó quiz
);
```

---

### Paso 4.2: Crear tabla con JSONB metadata
**Acciones**:
- [ ] Implementar CREATE TABLE con columnas
- [ ] Agregar FKs con SET NULL (datos históricos)
- [ ] Crear índices estratégicos

**Contenido (Parte 2 - Tabla)**:
```sql
-- Tabla para log de actividades del usuario
-- Uso: historial, analytics, actividad reciente en Home
CREATE TABLE IF NOT EXISTS user_activity_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    activity_type activity_type NOT NULL,
    material_id UUID,
    school_id UUID,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT fk_user_activity_log_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_user_activity_log_material
        FOREIGN KEY (material_id)
        REFERENCES materials(id)
        ON DELETE SET NULL,
    CONSTRAINT fk_user_activity_log_school
        FOREIGN KEY (school_id)
        REFERENCES schools(id)
        ON DELETE SET NULL
);

-- Índices para queries frecuentes
CREATE INDEX idx_user_activity_user_created 
    ON user_activity_log(user_id, created_at DESC);

CREATE INDEX idx_user_activity_school 
    ON user_activity_log(school_id, created_at DESC);

CREATE INDEX idx_user_activity_type 
    ON user_activity_log(activity_type);

-- Índice parcial para rate limiting (solo última hora)
CREATE INDEX idx_user_activity_rate_limit
    ON user_activity_log(user_id, activity_type, created_at)
    WHERE created_at > NOW() - INTERVAL '1 hour';

-- Comentarios
COMMENT ON TABLE user_activity_log IS 'Log de actividades del usuario para historial y analytics';
COMMENT ON COLUMN user_activity_log.activity_type IS 'Tipo de actividad realizada';
COMMENT ON COLUMN user_activity_log.metadata IS 'Datos adicionales en JSON (ej: score, pages, time_spent)';
COMMENT ON COLUMN user_activity_log.created_at IS 'Timestamp de la actividad';
```

**Validación**:
- [ ] ENUM definido correctamente
- [ ] Tabla usa el ENUM
- [ ] JSONB para metadata flexible
- [ ] SET NULL en FKs para preservar históricos

---

### Paso 4.3: Crear archivo de constraints
**Acciones**:
- [ ] Crear archivo `postgres/migrations/constraints/013_create_user_activity_log.sql`

**Contenido**:
```sql
-- No additional constraints needed (already in structure file)
```

---

### Paso 4.4: Ejecutar migración en local
**Acciones**:
- [ ] Ejecutar archivo de estructura
- [ ] Ejecutar archivo de constraints
- [ ] Verificar ENUM: `\dT activity_type`
- [ ] Verificar tabla: `\d user_activity_log`
- [ ] Verificar índices

**Comando**:
```bash
psql -U postgres -d edugo_db -f postgres/migrations/structure/013_create_user_activity_log.sql
psql -U postgres -d edugo_db -f postgres/migrations/constraints/013_create_user_activity_log.sql

psql -U postgres -d edugo_db -c "\dT activity_type"
psql -U postgres -d edugo_db -c "\d user_activity_log"
psql -U postgres -d edugo_db -c "\di idx_user_activity*"
```

**Resultado esperado**: ENUM y tabla creados sin errores

---

### Paso 4.5: Probar ENUM y JSONB
**Acciones**:
- [ ] Insertar actividad con metadata JSON
- [ ] Verificar que ENUM valida valores
- [ ] Probar queries con JSONB

**Comandos de prueba**:
```sql
-- Insertar actividad válida con metadata
INSERT INTO user_activity_log (user_id, activity_type, material_id, school_id, metadata)
VALUES (
    (SELECT id FROM users LIMIT 1),
    'material_progress',
    (SELECT id FROM materials LIMIT 1),
    (SELECT id FROM schools LIMIT 1),
    '{"page": 5, "total_pages": 10, "time_spent_seconds": 120}'::jsonb
);

-- Insertar quiz completado
INSERT INTO user_activity_log (user_id, activity_type, material_id, metadata)
VALUES (
    (SELECT id FROM users LIMIT 1),
    'quiz_passed',
    (SELECT id FROM materials LIMIT 1),
    '{"score": 85, "total_questions": 20, "correct_answers": 17}'::jsonb
);

-- Intentar valor inválido en ENUM (debe fallar)
INSERT INTO user_activity_log (user_id, activity_type, material_id)
VALUES (
    (SELECT id FROM users LIMIT 1),
    'invalid_activity',  -- Error: no es parte del ENUM
    (SELECT id FROM materials LIMIT 1)
);

-- Query con filtro en JSONB
SELECT * FROM user_activity_log
WHERE metadata->>'score' IS NOT NULL
  AND (metadata->>'score')::int > 80;

-- Query de actividad reciente
SELECT 
    activity_type,
    metadata,
    created_at
FROM user_activity_log
WHERE user_id = (SELECT id FROM users LIMIT 1)
ORDER BY created_at DESC
LIMIT 5;

-- Limpiar
DELETE FROM user_activity_log;
```

**Resultado esperado**: ENUM valida correctamente, JSONB funciona

---

**COMMIT FASE 4**:
```
feat(database): agregar tabla user_activity_log para tracking de actividades

- Crear ENUM activity_type con 8 tipos de actividad
- Crear migración 013_create_user_activity_log.sql
- Tabla almacena log de actividades para historial y analytics
- JSONB metadata para datos flexibles por tipo de actividad
- SET NULL en FKs para preservar datos históricos
- Índices: (user_id, created_at), (school_id, created_at), activity_type
- Índice parcial para rate limiting (última hora)

Parte de FASE 1 UI Roadmap - Actividad reciente en Home

Relacionado: #[número-de-issue]
```

**ID del commit**: [Se llenará después de hacer el commit]

---

## FASE 5: Ejecutar Tests de Validación (30 min)

### Objetivo
Validar que las 3 migraciones funcionan correctamente con tests exhaustivos.

### Paso 5.1: Test de estructura de tablas
**Acciones**:
- [ ] Verificar que las 3 tablas existen
- [ ] Verificar columnas y tipos de datos
- [ ] Verificar constraints (PK, FK, UNIQUE)
- [ ] Verificar defaults

**Script de test** (crear en `postgres/tests/test_fase1_structure.sql`):
```sql
-- Test: Verificar que tablas existen
SELECT 
    CASE 
        WHEN COUNT(*) = 3 THEN 'PASS: Las 3 tablas existen'
        ELSE 'FAIL: Faltan tablas'
    END as test_result
FROM information_schema.tables
WHERE table_name IN ('user_active_context', 'user_favorites', 'user_activity_log');

-- Test: Verificar columnas de user_active_context
SELECT 
    CASE 
        WHEN COUNT(*) = 6 THEN 'PASS: user_active_context tiene 6 columnas'
        ELSE 'FAIL: Columnas incorrectas'
    END as test_result
FROM information_schema.columns
WHERE table_name = 'user_active_context';

-- Test: Verificar UNIQUE constraint
SELECT 
    CASE 
        WHEN COUNT(*) >= 1 THEN 'PASS: UNIQUE constraint existe'
        ELSE 'FAIL: Falta UNIQUE constraint'
    END as test_result
FROM information_schema.table_constraints
WHERE table_name = 'user_active_context' 
  AND constraint_type = 'UNIQUE';

-- Similares para user_favorites y user_activity_log...
```

**Ejecutar**:
```bash
psql -U postgres -d edugo_db -f postgres/tests/test_fase1_structure.sql
```

---

### Paso 5.2: Test de performance
**Acciones**:
- [ ] Insertar 10K registros en user_activity_log
- [ ] Medir tiempo de INSERT
- [ ] Medir tiempo de queries frecuentes
- [ ] Verificar que índices se usan

**Script de test**:
```sql
-- Poblar con datos de prueba
INSERT INTO user_activity_log (user_id, activity_type, material_id, school_id, metadata)
SELECT 
    (SELECT id FROM users LIMIT 1),
    'material_progress',
    (SELECT id FROM materials LIMIT 1),
    (SELECT id FROM schools LIMIT 1),
    jsonb_build_object('page', gs.n, 'total_pages', 100)
FROM generate_series(1, 10000) AS gs(n);

-- Test: Query de actividad reciente (debe usar índice)
EXPLAIN ANALYZE
SELECT * FROM user_activity_log
WHERE user_id = (SELECT id FROM users LIMIT 1)
ORDER BY created_at DESC
LIMIT 10;

-- Verificar que usa idx_user_activity_user_created

-- Test: Aggregation
EXPLAIN ANALYZE
SELECT activity_type, COUNT(*)
FROM user_activity_log
WHERE user_id = (SELECT id FROM users LIMIT 1)
GROUP BY activity_type;

-- Limpiar
DELETE FROM user_activity_log;
```

---

### Paso 5.3: Test de integridad referencial
**Acciones**:
- [ ] Test CASCADE en user_favorites
- [ ] Test SET NULL en user_activity_log
- [ ] Test UNIQUE constraints

**Script de test**:
```sql
-- Test CASCADE: eliminar usuario elimina favoritos
BEGIN;
    INSERT INTO user_favorites (user_id, material_id)
    VALUES (
        (SELECT id FROM users WHERE email = 'student1@edugo.test'),
        (SELECT id FROM materials LIMIT 1)
    );
    
    DELETE FROM users WHERE email = 'student1@edugo.test';
    
    SELECT 
        CASE 
            WHEN COUNT(*) = 0 THEN 'PASS: CASCADE funciona'
            ELSE 'FAIL: Favorito no fue eliminado'
        END as test_result
    FROM user_favorites
    WHERE user_id = (SELECT id FROM users WHERE email = 'student1@edugo.test');
ROLLBACK;

-- Test SET NULL: eliminar material deja NULL en activity_log
BEGIN;
    -- Crear material temporal
    INSERT INTO materials (id, title, school_id)
    VALUES (
        gen_random_uuid(),
        'Test Material',
        (SELECT id FROM schools LIMIT 1)
    ) RETURNING id AS temp_material_id;
    
    -- Insertar actividad
    INSERT INTO user_activity_log (user_id, activity_type, material_id)
    VALUES (
        (SELECT id FROM users LIMIT 1),
        'material_started',
        (SELECT id FROM materials WHERE title = 'Test Material')
    );
    
    -- Eliminar material
    DELETE FROM materials WHERE title = 'Test Material';
    
    -- Verificar que material_id es NULL pero log existe
    SELECT 
        CASE 
            WHEN COUNT(*) = 1 AND material_id IS NULL 
            THEN 'PASS: SET NULL funciona'
            ELSE 'FAIL: Log fue eliminado o material_id no es NULL'
        END as test_result
    FROM user_activity_log
    WHERE activity_type = 'material_started'
    ORDER BY created_at DESC LIMIT 1;
ROLLBACK;
```

---

**COMMIT FASE 5**:
```
test(database): agregar tests de validación para FASE 1 UI Database

- Tests de estructura de tablas (columnas, tipos, constraints)
- Tests de performance con 10K registros
- Tests de integridad referencial (CASCADE, SET NULL)
- Tests de índices y query plans
- Scripts en postgres/tests/test_fase1_*.sql

Todas las validaciones pasan correctamente

Relacionado: #[número-de-issue]
```

**ID del commit**: [Se llenará después de hacer el commit]

---

## FASE 6: Actualizar Documentación (20 min)

### Objetivo
Documentar las nuevas tablas en README y CHANGELOG del proyecto.

### Paso 6.1: Actualizar postgres/README.md
**Acciones**:
- [ ] Abrir `postgres/README.md`
- [ ] Buscar sección de "Tablas" o "Schema"
- [ ] Agregar descripción de las 3 nuevas tablas
- [ ] Incluir casos de uso y relaciones

**Contenido a agregar**:
```markdown
### user_active_context

Almacena el contexto/escuela activa del usuario para filtrar datos en la UI.

**Propósito**: 
- Permite que usuarios con membresías en múltiples escuelas seleccionen cuál está activa
- Filtra materiales, cursos y datos según la escuela seleccionada
- Mejora UX evitando mostrar datos mezclados de múltiples escuelas

**Relaciones**:
- `user_id` → `users.id` (1:1, CASCADE)
- `school_id` → `schools.id` (N:1, CASCADE)
- `unit_id` → `academic_units.id` (N:1, SET NULL, opcional)

**Queries comunes**:
```sql
-- Obtener contexto activo del usuario
SELECT * FROM user_active_context WHERE user_id = ?;

-- Cambiar escuela activa (UPSERT)
INSERT INTO user_active_context (user_id, school_id)
VALUES (?, ?)
ON CONFLICT (user_id) DO UPDATE SET school_id = EXCLUDED.school_id;
```

### user_favorites

Almacena materiales marcados como favoritos por usuarios.

**Propósito**:
- Acceso rápido a materiales frecuentemente usados
- Lista personalizada de "lo más importante" para cada usuario
- Toggle de favoritos en UI de materiales

**Relaciones**:
- `user_id` → `users.id` (N:1, CASCADE)
- `material_id` → `materials.id` (N:1, CASCADE)

**Queries comunes**:
```sql
-- Listar favoritos del usuario
SELECT m.* FROM materials m
JOIN user_favorites uf ON m.id = uf.material_id
WHERE uf.user_id = ?
ORDER BY uf.created_at DESC;

-- Toggle favorito
INSERT INTO user_favorites (user_id, material_id) VALUES (?, ?)
ON CONFLICT DO NOTHING;  -- O DELETE si ya existe
```

### user_activity_log

Log de actividades del usuario para historial y analytics.

**Propósito**:
- Mostrar "Actividad reciente" en Home de la app
- Analytics de uso (materiales más usados, tiempo de estudio, etc.)
- Tracking de progreso del usuario
- Rate limiting de acciones

**Relaciones**:
- `user_id` → `users.id` (N:1, CASCADE)
- `material_id` → `materials.id` (N:1, SET NULL - históricos)
- `school_id` → `schools.id` (N:1, SET NULL - históricos)

**activity_type ENUM**:
- `material_started`, `material_progress`, `material_completed`
- `summary_viewed`
- `quiz_started`, `quiz_completed`, `quiz_passed`, `quiz_failed`

**metadata JSONB ejemplos**:
```json
// material_progress
{"page": 5, "total_pages": 10, "time_spent_seconds": 120}

// quiz_passed
{"score": 85, "total_questions": 20, "correct_answers": 17}
```

**Queries comunes**:
```sql
-- Actividad reciente del usuario
SELECT * FROM user_activity_log
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT 10;

-- Estadísticas de actividad (últimos 30 días)
SELECT activity_type, COUNT(*)
FROM user_activity_log
WHERE user_id = ? AND created_at > NOW() - INTERVAL '30 days'
GROUP BY activity_type;
```

**Consideraciones de escala**:
- Tabla de alto volumen (estimado 1M+ registros/día en producción)
- Considerar particionamiento por fecha cuando alcance ~50M registros
- Índice parcial para rate limiting (solo última hora)
```

---

### Paso 6.2: Actualizar CHANGELOG.md
**Acciones**:
- [ ] Abrir `CHANGELOG.md` en la raíz del proyecto
- [ ] Agregar nueva sección de versión (ej: postgres/v0.11.0)
- [ ] Listar las 3 nuevas tablas como features

**Contenido a agregar**:
```markdown
## [postgres/v0.11.0] - 2025-12-01

### Added - FASE 1 UI Roadmap

#### Nueva tabla: `user_active_context`
- Almacena contexto/escuela activa del usuario para filtrado en UI
- UNIQUE constraint en `user_id` (solo un contexto por usuario)
- Foreign keys con CASCADE a users/schools, SET NULL a academic_units
- Índices en user_id y school_id para performance
- Trigger automático para updated_at
- **Bloquea**: APIs de contexto de usuario (FASE 2)

#### Nueva tabla: `user_favorites`
- Almacena materiales marcados como favoritos
- UNIQUE constraint compuesto (user_id, material_id)
- CASCADE en ambos FKs para limpieza automática
- Índices en user_id, material_id y created_at
- **Bloquea**: Funcionalidad de favoritos en UI (FASE 4)

#### Nueva tabla: `user_activity_log`
- Log de actividades del usuario para historial y analytics
- ENUM `activity_type` con 8 tipos de actividad
- JSONB `metadata` para datos flexibles
- SET NULL en FKs para preservar datos históricos
- Índices estratégicos: (user_id, created_at), (school_id, created_at), activity_type
- Índice parcial para rate limiting (última hora)
- **Bloquea**: Actividad reciente en Home (FASE 4)

### Migration Files
- `postgres/migrations/structure/011_create_user_active_context.sql`
- `postgres/migrations/constraints/011_create_user_active_context.sql`
- `postgres/migrations/structure/012_create_user_favorites.sql`
- `postgres/migrations/constraints/012_create_user_favorites.sql`
- `postgres/migrations/structure/013_create_user_activity_log.sql`
- `postgres/migrations/constraints/013_create_user_activity_log.sql`

### Testing
- Tests de estructura de tablas
- Tests de performance (10K inserts en user_activity_log)
- Tests de integridad referencial (CASCADE, SET NULL)
- Todos los tests pasan exitosamente

### Documentation
- Actualizado postgres/README.md con nuevas tablas
- Documentación técnica detallada en docs/specs/fase1-ui-database/

### Related
- Parte de UI Roadmap FASE 1
- Issue #[número]
- PR #[número]
```

---

**COMMIT FASE 6**:
```
docs(database): actualizar documentación para FASE 1 UI Database

- Actualizar postgres/README.md con 3 nuevas tablas
- Describir propósito, relaciones y queries comunes
- Actualizar CHANGELOG.md con versión postgres/v0.11.0
- Incluir consideraciones de escala para user_activity_log

Parte de FASE 1 UI Roadmap

Relacionado: #[número-de-issue]
```

**ID del commit**: [Se llenará después de hacer el commit]

---

## FASE 7: Crear Tag y Finalizar (10 min)

### Objetivo
Crear tag de versión y preparar para merge a dev.

### Paso 7.1: Verificar que todo está commiteado
**Acciones**:
- [ ] `git status` debe estar limpio
- [ ] Revisar que todos los commits de las fases anteriores están hechos

**Comando**:
```bash
git status
git log --oneline -n 10
```

---

### Paso 7.2: Crear tag de versión
**Acciones**:
- [ ] Crear tag anotado `postgres/v0.11.0`
- [ ] Incluir mensaje descriptivo

**Comando**:
```bash
git tag -a postgres/v0.11.0 -m "Release postgres/v0.11.0 - FASE 1 UI Database

Nuevas tablas para soportar UI Roadmap:
- user_active_context: contexto/escuela activa del usuario
- user_favorites: materiales favoritos
- user_activity_log: log de actividades para analytics

Incluye:
- 6 archivos de migración (structure + constraints)
- Tests de validación completos
- Documentación actualizada (README, CHANGELOG)

Bloquea FASE 2 (APIs) y FASE 4 (UI Estudiantes) del roadmap.
"

git tag -l postgres/v0.11.0
```

---

### Paso 7.3: Preparar para push
**Acciones**:
- [ ] Revisar todos los commits
- [ ] Confirmar que rama está lista para push
- [ ] Documentar próximos pasos

**Comando**:
```bash
git log --oneline origin/dev..HEAD
```

**Resultado esperado**: 6-7 commits listos para push

---

**COMMIT FASE 7**: No aplica (solo tag)

**TAG**: `postgres/v0.11.0`

---

## Resumen de Commits Esperados

Al finalizar, deberíamos tener estos commits en la rama:

1. **FASE 2**: feat(database): agregar tabla user_active_context
2. **FASE 3**: feat(database): agregar tabla user_favorites  
3. **FASE 4**: feat(database): agregar tabla user_activity_log
4. **FASE 5**: test(database): agregar tests de validación FASE 1
5. **FASE 6**: docs(database): actualizar documentación FASE 1
6. **TAG**: postgres/v0.11.0

---

## Próximos Pasos (Fuera de este plan)

1. **Push a remoto**: 
   ```bash
   git push origin feature/fase1-ui-database-infrastructure
   git push origin postgres/v0.11.0
   ```

2. **Crear Pull Request** a `dev`:
   - Título: "feat(database): FASE 1 UI Database - 3 nuevas tablas"
   - Descripción: Incluir resumen, checklist de validación
   - Reviewers: Asignar a equipo

3. **Ejecutar migraciones en dev**:
   - Coordinar con DevOps/Backend
   - Ejecutar scripts en ambiente dev
   - Validar que todo funciona

4. **Iniciar FASE 2**: APIs que consumen estas tablas
   - Proyecto: edugo-api-mobile
   - Endpoints: `/v1/users/me/schools`, `/v1/users/me/active-school`, etc.

---

## Notas Finales

- ✅ Plan sigue metodología TDD
- ✅ Commits atómicos (1 por fase)
- ✅ Tests antes de documentar
- ✅ Documentación continua
- ✅ Tag de versión al final

**Duración total estimada**: ~2 horas de trabajo enfocado
