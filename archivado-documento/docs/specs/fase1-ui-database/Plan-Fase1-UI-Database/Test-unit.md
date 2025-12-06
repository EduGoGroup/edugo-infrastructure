# Tests Unitarios y de Integración - FASE 1 UI Database

> **Especificación completa de tests a implementar**

---

## Filosofía de Testing

### Principios
1. **Cobertura completa**: Validar estructura, funcionalidad y performance
2. **Automatizable**: Scripts SQL ejecutables sin intervención manual
3. **Repetible**: Tests deben pasar múltiples veces con mismos resultados
4. **Aislado**: Cada test es independiente (usa transacciones o limpia datos)
5. **Documentado**: Cada test explica QUÉ valida y POR QUÉ es importante

### Tipos de Tests
- **Estructura**: Validar que tablas, columnas, constraints e índices existen correctamente
- **Funcionalidad**: Validar que constraints, triggers y FKs funcionan como se espera
- **Performance**: Validar que queries son eficientes con datos realistas
- **Integridad**: Validar CASCADE, SET NULL y datos históricos

---

## TEST SUITE 1: Estructura de Tablas

**Archivo**: `postgres/tests/test_fase1_structure.sql`  
**Propósito**: Validar que las 3 tablas fueron creadas con estructura correcta  
**Duración estimada**: < 1 segundo

### Test 1.1: Tablas existen

```sql
-- Test: Verificar que las 3 tablas nuevas existen
DO $$
DECLARE
    table_count INT;
BEGIN
    SELECT COUNT(*) INTO table_count
    FROM information_schema.tables
    WHERE table_schema = 'public'
      AND table_name IN ('user_active_context', 'user_favorites', 'user_activity_log');
    
    IF table_count = 3 THEN
        RAISE NOTICE '✅ PASS: Las 3 tablas existen';
    ELSE
        RAISE EXCEPTION '❌ FAIL: Solo % de 3 tablas existen', table_count;
    END IF;
END $$;
```

**Criterio de éxito**: NOTICE con "PASS"  
**Criterio de fallo**: EXCEPTION

---

### Test 1.2: Estructura de user_active_context

```sql
-- Test: Verificar columnas de user_active_context
DO $$
DECLARE
    expected_columns TEXT[] := ARRAY['id', 'user_id', 'school_id', 'unit_id', 'created_at', 'updated_at'];
    actual_columns TEXT[];
    missing_columns TEXT[];
BEGIN
    SELECT ARRAY_AGG(column_name::TEXT ORDER BY column_name)
    INTO actual_columns
    FROM information_schema.columns
    WHERE table_name = 'user_active_context';
    
    -- Verificar que todas las columnas esperadas existen
    SELECT ARRAY_AGG(col)
    INTO missing_columns
    FROM UNNEST(expected_columns) AS col
    WHERE col NOT IN (SELECT UNNEST(actual_columns));
    
    IF missing_columns IS NULL THEN
        RAISE NOTICE '✅ PASS: user_active_context tiene todas las columnas';
    ELSE
        RAISE EXCEPTION '❌ FAIL: Faltan columnas: %', missing_columns;
    END IF;
END $$;
```

**Valida**: 
- id, user_id, school_id, unit_id, created_at, updated_at existen

---

### Test 1.3: Tipos de datos correctos

```sql
-- Test: Verificar tipos de datos de user_active_context
DO $$
DECLARE
    wrong_types TEXT;
BEGIN
    SELECT STRING_AGG(column_name || ' (' || data_type || ')', ', ')
    INTO wrong_types
    FROM information_schema.columns
    WHERE table_name = 'user_active_context'
      AND (
          (column_name = 'id' AND data_type != 'uuid') OR
          (column_name = 'user_id' AND data_type != 'uuid') OR
          (column_name = 'school_id' AND data_type != 'uuid') OR
          (column_name = 'unit_id' AND data_type != 'uuid') OR
          (column_name LIKE '%_at' AND data_type != 'timestamp with time zone')
      );
    
    IF wrong_types IS NULL THEN
        RAISE NOTICE '✅ PASS: Tipos de datos correctos';
    ELSE
        RAISE EXCEPTION '❌ FAIL: Tipos incorrectos: %', wrong_types;
    END IF;
END $$;
```

**Valida**:
- UUIDs para id, user_id, school_id, unit_id
- TIMESTAMP WITH TIME ZONE para created_at, updated_at

---

### Test 1.4: Primary Key y Unique Constraints

```sql
-- Test: Verificar PK y UNIQUE en user_active_context
DO $$
DECLARE
    pk_count INT;
    unique_count INT;
BEGIN
    -- Verificar PK en id
    SELECT COUNT(*)
    INTO pk_count
    FROM information_schema.table_constraints
    WHERE table_name = 'user_active_context'
      AND constraint_type = 'PRIMARY KEY'
      AND constraint_name LIKE '%_pkey';
    
    -- Verificar UNIQUE en user_id
    SELECT COUNT(*)
    INTO unique_count
    FROM information_schema.table_constraints
    WHERE table_name = 'user_active_context'
      AND constraint_type = 'UNIQUE'
      AND constraint_name LIKE '%user%';
    
    IF pk_count >= 1 AND unique_count >= 1 THEN
        RAISE NOTICE '✅ PASS: PK y UNIQUE constraints existen';
    ELSE
        RAISE EXCEPTION '❌ FAIL: PK: %, UNIQUE: %', pk_count, unique_count;
    END IF;
END $$;
```

**Valida**:
- PRIMARY KEY en id
- UNIQUE constraint en user_id

---

### Test 1.5: Foreign Keys

```sql
-- Test: Verificar FKs de user_active_context
DO $$
DECLARE
    fk_count INT;
    expected_fks INT := 3; -- user_id, school_id, unit_id
BEGIN
    SELECT COUNT(*)
    INTO fk_count
    FROM information_schema.table_constraints
    WHERE table_name = 'user_active_context'
      AND constraint_type = 'FOREIGN KEY';
    
    IF fk_count >= expected_fks THEN
        RAISE NOTICE '✅ PASS: % Foreign Keys encontradas', fk_count;
    ELSE
        RAISE EXCEPTION '❌ FAIL: Solo % de % FKs', fk_count, expected_fks;
    END IF;
END $$;
```

**Valida**:
- FK a users (user_id)
- FK a schools (school_id)
- FK a academic_units (unit_id)

---

### Test 1.6: Índices

```sql
-- Test: Verificar índices de user_active_context
DO $$
DECLARE
    idx_user INT;
    idx_school INT;
BEGIN
    -- Índice en user_id
    SELECT COUNT(*)
    INTO idx_user
    FROM pg_indexes
    WHERE tablename = 'user_active_context'
      AND indexname LIKE '%user%';
    
    -- Índice en school_id
    SELECT COUNT(*)
    INTO idx_school
    FROM pg_indexes
    WHERE tablename = 'user_active_context'
      AND indexname LIKE '%school%';
    
    IF idx_user >= 1 AND idx_school >= 1 THEN
        RAISE NOTICE '✅ PASS: Índices en user_id y school_id';
    ELSE
        RAISE EXCEPTION '❌ FAIL: idx_user: %, idx_school: %', idx_user, idx_school;
    END IF;
END $$;
```

**Valida**:
- idx_user_active_context_user existe
- idx_user_active_context_school existe

---

### Test 1.7: Trigger para updated_at

```sql
-- Test: Verificar trigger en user_active_context
DO $$
DECLARE
    trigger_count INT;
BEGIN
    SELECT COUNT(*)
    INTO trigger_count
    FROM information_schema.triggers
    WHERE event_object_table = 'user_active_context'
      AND trigger_name LIKE '%updated_at%';
    
    IF trigger_count >= 1 THEN
        RAISE NOTICE '✅ PASS: Trigger updated_at existe';
    ELSE
        RAISE EXCEPTION '❌ FAIL: Trigger no encontrado';
    END IF;
END $$;
```

**Valida**:
- Trigger set_updated_at_user_active_context existe

---

### Tests similares para user_favorites

```sql
-- Test 1.8 a 1.12: Repetir estructura para user_favorites
-- Verificar: columnas, tipos, UNIQUE(user_id, material_id), FKs, índices
```

---

### Tests similares para user_activity_log

```sql
-- Test 1.13 a 1.18: Repetir estructura para user_activity_log
-- Adicional: Verificar ENUM activity_type
```

---

### Test 1.19: ENUM activity_type

```sql
-- Test: Verificar que ENUM activity_type existe y tiene valores correctos
DO $$
DECLARE
    enum_values TEXT[];
    expected_count INT := 8;
    actual_count INT;
BEGIN
    SELECT ARRAY_AGG(enumlabel ORDER BY enumsortorder)
    INTO enum_values
    FROM pg_enum
    JOIN pg_type ON pg_enum.enumtypid = pg_type.oid
    WHERE pg_type.typname = 'activity_type';
    
    actual_count := ARRAY_LENGTH(enum_values, 1);
    
    IF actual_count = expected_count THEN
        RAISE NOTICE '✅ PASS: ENUM activity_type tiene % valores: %', actual_count, enum_values;
    ELSE
        RAISE EXCEPTION '❌ FAIL: ENUM tiene % valores, esperados %', actual_count, expected_count;
    END IF;
    
    -- Verificar valores específicos
    IF 'material_started' = ANY(enum_values) AND
       'quiz_passed' = ANY(enum_values) THEN
        RAISE NOTICE '✅ PASS: Valores clave del ENUM presentes';
    ELSE
        RAISE EXCEPTION '❌ FAIL: Faltan valores en ENUM';
    END IF;
END $$;
```

**Valida**:
- ENUM tiene 8 valores
- Contiene: material_started, material_progress, material_completed, summary_viewed, quiz_started, quiz_completed, quiz_passed, quiz_failed

---

## TEST SUITE 2: Funcionalidad

**Archivo**: `postgres/tests/test_fase1_integrity.sql`  
**Propósito**: Validar que constraints, triggers y FKs funcionan  
**Duración estimada**: 2-3 segundos

### Test 2.1: UNIQUE constraint en user_active_context

```sql
-- Test: UNIQUE en user_id debe prevenir duplicados
DO $$
DECLARE
    test_user_id UUID;
    test_school_id UUID;
BEGIN
    -- Obtener IDs de prueba
    SELECT id INTO test_user_id FROM users LIMIT 1;
    SELECT id INTO test_school_id FROM schools LIMIT 1;
    
    BEGIN
        -- Insertar primer registro (debe funcionar)
        INSERT INTO user_active_context (user_id, school_id)
        VALUES (test_user_id, test_school_id);
        
        -- Intentar duplicado (debe fallar)
        INSERT INTO user_active_context (user_id, school_id)
        VALUES (test_user_id, test_school_id);
        
        -- Si llegamos aquí, el test falló
        RAISE EXCEPTION '❌ FAIL: UNIQUE constraint no funcionó, duplicado permitido';
        
    EXCEPTION
        WHEN unique_violation THEN
            -- Comportamiento esperado
            RAISE NOTICE '✅ PASS: UNIQUE constraint funciona correctamente';
            ROLLBACK;
    END;
END $$;
```

**Valida**:
- No se pueden insertar 2 contextos para el mismo user_id
- Lanza unique_violation error

---

### Test 2.2: Trigger updated_at funciona

```sql
-- Test: Trigger actualiza updated_at automáticamente
DO $$
DECLARE
    test_user_id UUID;
    test_school_id UUID;
    test_school_id_2 UUID;
    created TIMESTAMP WITH TIME ZONE;
    updated TIMESTAMP WITH TIME ZONE;
BEGIN
    SELECT id INTO test_user_id FROM users LIMIT 1;
    SELECT id INTO test_school_id FROM schools LIMIT 1;
    SELECT id INTO test_school_id_2 FROM schools LIMIT 1 OFFSET 1;
    
    -- Insertar
    INSERT INTO user_active_context (user_id, school_id)
    VALUES (test_user_id, test_school_id)
    RETURNING created_at INTO created;
    
    -- Esperar 1 segundo
    PERFORM pg_sleep(1);
    
    -- Actualizar
    UPDATE user_active_context
    SET school_id = test_school_id_2
    WHERE user_id = test_user_id
    RETURNING updated_at INTO updated;
    
    -- Verificar que updated_at > created_at
    IF updated > created THEN
        RAISE NOTICE '✅ PASS: Trigger updated_at funciona (created: %, updated: %)', created, updated;
    ELSE
        RAISE EXCEPTION '❌ FAIL: updated_at no se actualizó';
    END IF;
    
    -- Limpiar
    DELETE FROM user_active_context WHERE user_id = test_user_id;
END $$;
```

**Valida**:
- updated_at se actualiza automáticamente
- updated_at > created_at después de UPDATE

---

### Test 2.3: CASCADE en user_favorites al eliminar usuario

```sql
-- Test: Eliminar usuario elimina sus favoritos (CASCADE)
DO $$
DECLARE
    temp_user_id UUID;
    temp_material_id UUID;
    favorite_count INT;
BEGIN
    -- Crear usuario temporal
    INSERT INTO users (id, email, password_hash, role, first_name, last_name)
    VALUES (gen_random_uuid(), 'test_cascade@example.com', 'hash', 'student', 'Test', 'User')
    RETURNING id INTO temp_user_id;
    
    SELECT id INTO temp_material_id FROM materials LIMIT 1;
    
    -- Agregar favorito
    INSERT INTO user_favorites (user_id, material_id)
    VALUES (temp_user_id, temp_material_id);
    
    -- Eliminar usuario
    DELETE FROM users WHERE id = temp_user_id;
    
    -- Verificar que favorito fue eliminado
    SELECT COUNT(*) INTO favorite_count
    FROM user_favorites
    WHERE user_id = temp_user_id;
    
    IF favorite_count = 0 THEN
        RAISE NOTICE '✅ PASS: CASCADE funciona, favorito eliminado';
    ELSE
        RAISE EXCEPTION '❌ FAIL: Favorito no fue eliminado, CASCADE no funciona';
    END IF;
END $$;
```

**Valida**:
- ON DELETE CASCADE en FK user_id funciona
- Eliminar usuario elimina automáticamente sus favoritos

---

### Test 2.4: SET NULL en user_activity_log al eliminar material

```sql
-- Test: Eliminar material pone NULL en activity_log (SET NULL)
DO $$
DECLARE
    test_user_id UUID;
    temp_material_id UUID;
    log_material_id UUID;
BEGIN
    SELECT id INTO test_user_id FROM users LIMIT 1;
    
    -- Crear material temporal
    INSERT INTO materials (id, title, school_id, subject_id, created_by)
    VALUES (
        gen_random_uuid(),
        'Test Material for Activity Log',
        (SELECT id FROM schools LIMIT 1),
        (SELECT id FROM subjects LIMIT 1),
        test_user_id
    )
    RETURNING id INTO temp_material_id;
    
    -- Insertar actividad
    INSERT INTO user_activity_log (user_id, activity_type, material_id)
    VALUES (test_user_id, 'material_started', temp_material_id);
    
    -- Eliminar material
    DELETE FROM materials WHERE id = temp_material_id;
    
    -- Verificar que log existe pero material_id es NULL
    SELECT material_id INTO log_material_id
    FROM user_activity_log
    WHERE user_id = test_user_id
      AND activity_type = 'material_started'
    ORDER BY created_at DESC
    LIMIT 1;
    
    IF log_material_id IS NULL THEN
        RAISE NOTICE '✅ PASS: SET NULL funciona, log preservado con material_id = NULL';
    ELSE
        RAISE EXCEPTION '❌ FAIL: material_id no es NULL: %', log_material_id;
    END IF;
    
    -- Limpiar
    DELETE FROM user_activity_log WHERE user_id = test_user_id;
END $$;
```

**Valida**:
- ON DELETE SET NULL en FK material_id funciona
- Log histórico se preserva con material_id = NULL

---

### Test 2.5: ENUM valida valores en user_activity_log

```sql
-- Test: ENUM rechaza valores inválidos
DO $$
DECLARE
    test_user_id UUID;
BEGIN
    SELECT id INTO test_user_id FROM users LIMIT 1;
    
    BEGIN
        -- Intentar insertar valor inválido
        INSERT INTO user_activity_log (user_id, activity_type)
        VALUES (test_user_id, 'invalid_activity'::activity_type);
        
        -- Si llegamos aquí, falló
        RAISE EXCEPTION '❌ FAIL: ENUM aceptó valor inválido';
        
    EXCEPTION
        WHEN invalid_text_representation THEN
            RAISE NOTICE '✅ PASS: ENUM rechaza valores inválidos correctamente';
    END;
END $$;
```

**Valida**:
- ENUM activity_type solo acepta valores definidos
- Rechaza valores no en el ENUM

---

### Test 2.6: JSONB metadata en user_activity_log

```sql
-- Test: JSONB acepta JSON válido y permite queries
DO $$
DECLARE
    test_user_id UUID;
    log_id UUID;
    score INT;
BEGIN
    SELECT id INTO test_user_id FROM users LIMIT 1;
    
    -- Insertar actividad con metadata
    INSERT INTO user_activity_log (user_id, activity_type, metadata)
    VALUES (
        test_user_id,
        'quiz_passed',
        '{"score": 95, "total_questions": 20, "time_seconds": 300}'::jsonb
    )
    RETURNING id INTO log_id;
    
    -- Extraer score de JSONB
    SELECT (metadata->>'score')::INT INTO score
    FROM user_activity_log
    WHERE id = log_id;
    
    IF score = 95 THEN
        RAISE NOTICE '✅ PASS: JSONB funciona, score extraído: %', score;
    ELSE
        RAISE EXCEPTION '❌ FAIL: Score incorrecto: %', score;
    END IF;
    
    -- Limpiar
    DELETE FROM user_activity_log WHERE id = log_id;
END $$;
```

**Valida**:
- Columna metadata acepta JSONB
- Queries JSONB funcionan (`->`, `->>`)

---

## TEST SUITE 3: Performance

**Archivo**: `postgres/tests/test_fase1_performance.sql`  
**Propósito**: Validar que queries son eficientes con volumen realista  
**Duración estimada**: 5-10 segundos

### Test 3.1: Inserción masiva en user_activity_log

```sql
-- Test: Insertar 10K registros rápidamente
DO $$
DECLARE
    test_user_id UUID;
    test_material_id UUID;
    test_school_id UUID;
    start_time TIMESTAMP;
    end_time TIMESTAMP;
    duration INTERVAL;
    insert_count INT;
BEGIN
    SELECT id INTO test_user_id FROM users LIMIT 1;
    SELECT id INTO test_material_id FROM materials LIMIT 1;
    SELECT id INTO test_school_id FROM schools LIMIT 1;
    
    start_time := clock_timestamp();
    
    -- Insertar 10K registros
    INSERT INTO user_activity_log (user_id, activity_type, material_id, school_id, metadata)
    SELECT 
        test_user_id,
        'material_progress',
        test_material_id,
        test_school_id,
        jsonb_build_object('page', gs.n, 'total_pages', 100)
    FROM generate_series(1, 10000) AS gs(n);
    
    GET DIAGNOSTICS insert_count = ROW_COUNT;
    end_time := clock_timestamp();
    duration := end_time - start_time;
    
    RAISE NOTICE '✅ Insertados % registros en %', insert_count, duration;
    
    IF duration < INTERVAL '3 seconds' THEN
        RAISE NOTICE '✅ PASS: Inserción rápida (< 3s)';
    ELSE
        RAISE WARNING '⚠️  WARNING: Inserción lenta: %', duration;
    END IF;
    
    -- Limpiar
    DELETE FROM user_activity_log WHERE user_id = test_user_id;
END $$;
```

**Valida**:
- 10K inserts en < 3 segundos
- Columna JSONB no degrada performance significativamente

---

### Test 3.2: Query de actividad reciente usa índice

```sql
-- Test: Query de actividad reciente usa idx_user_activity_user_created
DO $$
DECLARE
    test_user_id UUID;
    query_plan TEXT;
BEGIN
    SELECT id INTO test_user_id FROM users LIMIT 1;
    
    -- Poblar con datos
    INSERT INTO user_activity_log (user_id, activity_type)
    SELECT test_user_id, 'material_progress'
    FROM generate_series(1, 1000);
    
    -- Obtener query plan
    SELECT query_plan INTO query_plan
    FROM (
        EXPLAIN (FORMAT TEXT)
        SELECT * FROM user_activity_log
        WHERE user_id = test_user_id
        ORDER BY created_at DESC
        LIMIT 10
    ) AS plan_output(query_plan);
    
    IF query_plan LIKE '%idx_user_activity_user_created%' OR query_plan LIKE '%Index Scan%' THEN
        RAISE NOTICE '✅ PASS: Query usa índice';
        RAISE NOTICE 'Plan: %', query_plan;
    ELSE
        RAISE WARNING '⚠️  WARNING: Query no usa índice';
        RAISE NOTICE 'Plan: %', query_plan;
    END IF;
    
    -- Limpiar
    DELETE FROM user_activity_log WHERE user_id = test_user_id;
END $$;
```

**Valida**:
- Query de actividad reciente usa Index Scan
- No hace Sequential Scan (lento)

---

### Test 3.3: Aggregation por activity_type es eficiente

```sql
-- Test: Aggregation usa índice idx_user_activity_type
EXPLAIN ANALYZE
SELECT activity_type, COUNT(*), AVG((metadata->>'score')::INT)
FROM user_activity_log
WHERE user_id = (SELECT id FROM users LIMIT 1)
  AND created_at > NOW() - INTERVAL '30 days'
GROUP BY activity_type;
```

**Criterio de éxito**: 
- Execution time < 100ms con 10K registros
- Usa idx_user_activity_user_created para filtro
- Usa idx_user_activity_type para GROUP BY (idealmente)

---

### Test 3.4: JSONB queries son eficientes

```sql
-- Test: Query JSONB con filtro
EXPLAIN ANALYZE
SELECT * FROM user_activity_log
WHERE metadata->>'score' IS NOT NULL
  AND (metadata->>'score')::INT > 80
LIMIT 100;
```

**Criterio de éxito**:
- Execution time razonable (< 200ms con 10K registros)
- **Nota**: Sin índice GIN en metadata, será Sequential Scan (aceptable para MVP)
- **Futuro**: Considerar `CREATE INDEX idx_activity_metadata_gin ON user_activity_log USING GIN (metadata);`

---

## TEST SUITE 4: Integración

**Propósito**: Validar escenarios de uso realistas  
**Duración estimada**: 3-5 segundos

### Test 4.1: Flujo completo de cambio de contexto

```sql
-- Test: Usuario cambia de escuela activa (UPSERT)
DO $$
DECLARE
    test_user_id UUID;
    school1_id UUID;
    school2_id UUID;
    context_id UUID;
    active_school UUID;
BEGIN
    SELECT id INTO test_user_id FROM users LIMIT 1;
    SELECT id INTO school1_id FROM schools LIMIT 1;
    SELECT id INTO school2_id FROM schools LIMIT 1 OFFSET 1;
    
    -- Primera vez: INSERT
    INSERT INTO user_active_context (user_id, school_id)
    VALUES (test_user_id, school1_id)
    ON CONFLICT (user_id) DO UPDATE SET school_id = EXCLUDED.school_id
    RETURNING id INTO context_id;
    
    -- Verificar
    SELECT school_id INTO active_school
    FROM user_active_context
    WHERE user_id = test_user_id;
    
    IF active_school = school1_id THEN
        RAISE NOTICE '✅ PASS: Primera inserción correcta';
    ELSE
        RAISE EXCEPTION '❌ FAIL: School incorrecta';
    END IF;
    
    -- Segunda vez: UPDATE via UPSERT
    INSERT INTO user_active_context (user_id, school_id)
    VALUES (test_user_id, school2_id)
    ON CONFLICT (user_id) DO UPDATE SET school_id = EXCLUDED.school_id;
    
    -- Verificar
    SELECT school_id INTO active_school
    FROM user_active_context
    WHERE user_id = test_user_id;
    
    IF active_school = school2_id THEN
        RAISE NOTICE '✅ PASS: UPSERT actualización correcta';
    ELSE
        RAISE EXCEPTION '❌ FAIL: UPSERT no actualizó';
    END IF;
    
    -- Limpiar
    DELETE FROM user_active_context WHERE user_id = test_user_id;
END $$;
```

**Valida**:
- INSERT ... ON CONFLICT funciona para cambiar contexto
- UNIQUE constraint permite UPSERT pattern

---

### Test 4.2: Flujo completo de toggle favorito

```sql
-- Test: Usuario marca/desmarca favorito
DO $$
DECLARE
    test_user_id UUID;
    test_material_id UUID;
    is_favorite BOOLEAN;
BEGIN
    SELECT id INTO test_user_id FROM users LIMIT 1;
    SELECT id INTO test_material_id FROM materials LIMIT 1;
    
    -- Agregar favorito
    INSERT INTO user_favorites (user_id, material_id)
    VALUES (test_user_id, test_material_id)
    ON CONFLICT DO NOTHING;
    
    -- Verificar que existe
    SELECT EXISTS(
        SELECT 1 FROM user_favorites
        WHERE user_id = test_user_id AND material_id = test_material_id
    ) INTO is_favorite;
    
    IF is_favorite THEN
        RAISE NOTICE '✅ PASS: Favorito agregado';
    ELSE
        RAISE EXCEPTION '❌ FAIL: Favorito no se agregó';
    END IF;
    
    -- Eliminar favorito
    DELETE FROM user_favorites
    WHERE user_id = test_user_id AND material_id = test_material_id;
    
    -- Verificar que no existe
    SELECT EXISTS(
        SELECT 1 FROM user_favorites
        WHERE user_id = test_user_id AND material_id = test_material_id
    ) INTO is_favorite;
    
    IF NOT is_favorite THEN
        RAISE NOTICE '✅ PASS: Favorito eliminado';
    ELSE
        RAISE EXCEPTION '❌ FAIL: Favorito no se eliminó';
    END IF;
END $$;
```

**Valida**:
- Pattern de toggle favorito funciona
- INSERT ... ON CONFLICT DO NOTHING para idempotencia

---

### Test 4.3: Flujo de logging de actividad con metadata

```sql
-- Test: Simular actividad completa de usuario en material
DO $$
DECLARE
    test_user_id UUID;
    test_material_id UUID;
    test_school_id UUID;
    activity_count INT;
BEGIN
    SELECT id INTO test_user_id FROM users LIMIT 1;
    SELECT id INTO test_material_id FROM materials LIMIT 1;
    SELECT id INTO test_school_id FROM schools LIMIT 1;
    
    -- 1. Iniciar material
    INSERT INTO user_activity_log (user_id, activity_type, material_id, school_id, metadata)
    VALUES (test_user_id, 'material_started', test_material_id, test_school_id, '{"page": 1}'::jsonb);
    
    -- 2. Progreso en múltiples páginas
    INSERT INTO user_activity_log (user_id, activity_type, material_id, school_id, metadata)
    SELECT 
        test_user_id,
        'material_progress',
        test_material_id,
        test_school_id,
        jsonb_build_object('page', n, 'time_spent_seconds', n * 30)
    FROM generate_series(2, 5) AS n;
    
    -- 3. Completar material
    INSERT INTO user_activity_log (user_id, activity_type, material_id, school_id, metadata)
    VALUES (test_user_id, 'material_completed', test_material_id, test_school_id, '{"total_time_seconds": 600}'::jsonb);
    
    -- 4. Ver resumen
    INSERT INTO user_activity_log (user_id, activity_type, material_id, school_id, metadata)
    VALUES (test_user_id, 'summary_viewed', test_material_id, test_school_id, '{"summary_length": 500}'::jsonb);
    
    -- 5. Hacer quiz
    INSERT INTO user_activity_log (user_id, activity_type, material_id, school_id, metadata)
    VALUES (test_user_id, 'quiz_started', test_material_id, test_school_id, '{"question_count": 10}'::jsonb);
    
    INSERT INTO user_activity_log (user_id, activity_type, material_id, school_id, metadata)
    VALUES (test_user_id, 'quiz_passed', test_material_id, test_school_id, '{"score": 90, "correct": 9, "total": 10}'::jsonb);
    
    -- Verificar que todo se registró
    SELECT COUNT(*) INTO activity_count
    FROM user_activity_log
    WHERE user_id = test_user_id
      AND material_id = test_material_id;
    
    IF activity_count = 9 THEN  -- 1 start + 4 progress + 1 complete + 1 summary + 2 quiz
        RAISE NOTICE '✅ PASS: Flujo completo registrado (% actividades)', activity_count;
    ELSE
        RAISE EXCEPTION '❌ FAIL: Solo % actividades registradas', activity_count;
    END IF;
    
    -- Probar query de "actividad reciente"
    PERFORM * FROM user_activity_log
    WHERE user_id = test_user_id
    ORDER BY created_at DESC
    LIMIT 5;
    
    RAISE NOTICE '✅ PASS: Query de actividad reciente funciona';
    
    -- Limpiar
    DELETE FROM user_activity_log WHERE user_id = test_user_id;
END $$;
```

**Valida**:
- Flujo completo de actividades se registra correctamente
- Diferentes tipos de metadata JSONB funcionan
- Query de actividad reciente retorna resultados ordenados

---

## Resumen de Tests

| Suite | Archivo | Tests | Propósito |
|-------|---------|-------|-----------|
| **1. Estructura** | test_fase1_structure.sql | ~19 tests | Tablas, columnas, tipos, constraints, índices |
| **2. Funcionalidad** | test_fase1_integrity.sql | ~6 tests | UNIQUE, CASCADE, SET NULL, ENUM, JSONB, triggers |
| **3. Performance** | test_fase1_performance.sql | ~4 tests | Inserts masivos, índices, query plans |
| **4. Integración** | test_fase1_integrity.sql | ~3 tests | Flujos completos de uso |
| **Total** | 3 archivos | **~32 tests** | Cobertura completa |

---

## Ejecución de Tests

### Ejecutar todos los tests

```bash
# Suite 1: Estructura
psql -U postgres -d edugo_db -f postgres/tests/test_fase1_structure.sql

# Suite 2 & 4: Funcionalidad e Integración
psql -U postgres -d edugo_db -f postgres/tests/test_fase1_integrity.sql

# Suite 3: Performance
psql -U postgres -d edugo_db -f postgres/tests/test_fase1_performance.sql
```

### Ejecutar con reporte

```bash
# Con salida a archivo
psql -U postgres -d edugo_db -f postgres/tests/test_fase1_structure.sql > test_results.txt 2>&1

# Verificar solo errores
grep -E "FAIL|ERROR" test_results.txt
```

### Script automatizado

```bash
#!/bin/bash
# run_fase1_tests.sh

echo "🧪 Ejecutando tests FASE 1 UI Database..."

RESULTS_DIR="postgres/tests/results"
mkdir -p $RESULTS_DIR

TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Suite 1
echo "📋 Test Suite 1: Estructura..."
psql -U postgres -d edugo_db -f postgres/tests/test_fase1_structure.sql \
    > $RESULTS_DIR/structure_$TIMESTAMP.log 2>&1

# Suite 2 & 4
echo "🔗 Test Suite 2 & 4: Integridad..."
psql -U postgres -d edugo_db -f postgres/tests/test_fase1_integrity.sql \
    > $RESULTS_DIR/integrity_$TIMESTAMP.log 2>&1

# Suite 3
echo "⚡ Test Suite 3: Performance..."
psql -U postgres -d edugo_db -f postgres/tests/test_fase1_performance.sql \
    > $RESULTS_DIR/performance_$TIMESTAMP.log 2>&1

# Resumen
echo ""
echo "📊 Resumen de Resultados:"
echo "========================"

for log in $RESULTS_DIR/*_$TIMESTAMP.log; do
    echo "Archivo: $(basename $log)"
    grep -E "✅ PASS|❌ FAIL" $log | wc -l | xargs echo "  Tests ejecutados:"
    grep "✅ PASS" $log | wc -l | xargs echo "  Pasados:"
    grep "❌ FAIL" $log | wc -l | xargs echo "  Fallados:"
    echo ""
done

echo "✅ Tests completados. Logs en: $RESULTS_DIR/"
```

---

## Criterios de Aceptación Global

Para que la FASE 1 sea considerada completa:

```
✅ Todos los tests de estructura pasan (19/19)
✅ Todos los tests de funcionalidad pasan (6/6)
✅ Tests de performance muestran tiempos aceptables (< 3s inserts, < 100ms queries)
✅ Tests de integración simulan flujos reales exitosamente (3/3)
✅ No hay errores en logs de PostgreSQL
✅ Query plans usan índices apropiados
```

**Total esperado**: 32 tests PASS, 0 FAIL

---

**Fin de Especificación de Tests**

Estos tests aseguran que las 3 tablas nuevas funcionan correctamente y están listas para ser consumidas por las APIs en FASE 2.
