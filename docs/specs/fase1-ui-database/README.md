# FASE 1: Base de Datos - UI Roadmap

> **Proyecto**: edugo-infrastructure  
> **Responsable**: Base de Datos (PostgreSQL)  
> **Duración estimada**: 1-2 días  
> **Prioridad**: 🔴 CRÍTICA  
> **Bloquea**: APIs (Fase 2) y UI Estudiantes (Fase 4)

---

## Contexto

Este proyecto es parte del UI Roadmap de EduGo. La FASE 1 implementa las tablas de base de datos necesarias para soportar las nuevas funcionalidades de la interfaz de usuario para estudiantes.

**Dependencias**:
- Sin BD → APIs no pueden funcionar
- Sin APIs → Apps no tienen datos

**Orden del Roadmap**:
```
FASE 1: BASE DE DATOS (edugo-infrastructure) ← ESTAMOS AQUÍ
   ↓
FASE 2: APIs (api-mobile primero, luego api-admin)
   ↓
FASE 3: MÓDULOS CROSS (SPM compartidos)
   ↓
FASE 4: APP ESTUDIANTES (completa)
   ↓
FASE 5: APP ADMINISTRACIÓN (completa)
```

---

## Objetivo

Implementar 3 nuevas tablas PostgreSQL que soportan:

1. **Contexto de usuario** (`user_active_context`)
   - Almacenar qué escuela tiene seleccionada el usuario
   - Permitir filtrar datos por escuela activa
   - Bloquea: Selector de escuela en UI

2. **Favoritos** (`user_favorites`)
   - Materiales marcados como favoritos por usuarios
   - Bloquea: Funcionalidad de favoritos en UI

3. **Log de actividad** (`user_activity_log`)
   - Rastrear actividades del usuario
   - Bloquea: Sección "Actividad reciente" en Home

---

## Tareas a Implementar

### Tarea 1.1.1: Tabla `user_active_context`
**Prioridad**: 🔴 CRÍTICA  
**Archivo**: `postgres/migrations/structure/011_create_user_active_context.sql`

**Descripción**:
Almacena el contexto/escuela activa del usuario para filtrar datos en la UI.

**Columnas**:
- `id`: UUID (PK)
- `user_id`: UUID (FK a users, NOT NULL, UNIQUE)
- `school_id`: UUID (FK a schools, NOT NULL)
- `unit_id`: UUID (FK a academic_units, nullable)
- `created_at`: TIMESTAMP WITH TIME ZONE
- `updated_at`: TIMESTAMP WITH TIME ZONE

**Constraints**:
- UNIQUE(user_id) - Un usuario solo tiene un contexto activo
- FK a users ON DELETE CASCADE
- FK a schools ON DELETE CASCADE
- FK a academic_units ON DELETE SET NULL

**Índices**:
- `idx_user_active_context_user` en user_id
- `idx_user_active_context_school` en school_id

**Trigger**:
- `set_updated_at_user_active_context` para actualizar `updated_at`

---

### Tarea 1.1.2: Tabla `user_favorites`
**Prioridad**: 🟡 MEDIA  
**Archivo**: `postgres/migrations/structure/012_create_user_favorites.sql`

**Descripción**:
Almacena materiales marcados como favoritos por usuarios.

**Columnas**:
- `id`: UUID (PK)
- `user_id`: UUID (FK a users, NOT NULL)
- `material_id`: UUID (FK a materials, NOT NULL)
- `created_at`: TIMESTAMP WITH TIME ZONE

**Constraints**:
- UNIQUE(user_id, material_id) - Un usuario no puede duplicar favoritos
- FK a users ON DELETE CASCADE
- FK a materials ON DELETE CASCADE

**Índices**:
- `idx_user_favorites_user` en user_id
- `idx_user_favorites_material` en material_id
- `idx_user_favorites_created` en created_at DESC

---

### Tarea 1.1.3: Tabla `user_activity_log`
**Prioridad**: 🟡 MEDIA  
**Archivo**: `postgres/migrations/structure/013_create_user_activity_log.sql`

**Descripción**:
Log de actividades del usuario para historial y analytics.

**Tipo ENUM**: `activity_type`
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

**Columnas**:
- `id`: UUID (PK)
- `user_id`: UUID (FK a users, NOT NULL)
- `activity_type`: activity_type (NOT NULL)
- `material_id`: UUID (FK a materials, nullable)
- `school_id`: UUID (FK a schools, nullable)
- `metadata`: JSONB (default '{}')
- `created_at`: TIMESTAMP WITH TIME ZONE

**Constraints**:
- FK a users ON DELETE CASCADE
- FK a materials ON DELETE SET NULL
- FK a schools ON DELETE SET NULL

**Índices**:
- `idx_user_activity_user_created` en (user_id, created_at DESC)
- `idx_user_activity_school` en (school_id, created_at DESC)
- `idx_user_activity_type` en activity_type
- `idx_user_activity_rate_limit` en (user_id, activity_type, created_at) WHERE ...

**Nota**: Considerar particionamiento por fecha para escalabilidad futura.

---

## Archivos a Crear/Modificar

### Nuevos archivos:

1. **Structure (Tablas)**:
   - `postgres/migrations/structure/011_create_user_active_context.sql`
   - `postgres/migrations/structure/012_create_user_favorites.sql`
   - `postgres/migrations/structure/013_create_user_activity_log.sql`

2. **Constraints**:
   - `postgres/migrations/constraints/011_create_user_active_context.sql`
   - `postgres/migrations/constraints/012_create_user_favorites.sql`
   - `postgres/migrations/constraints/013_create_user_activity_log.sql`

3. **Indexes** (si aplica):
   - `postgres/migrations/indexes/011_create_user_active_context.sql`
   - `postgres/migrations/indexes/012_create_user_favorites.sql`
   - `postgres/migrations/indexes/013_create_user_activity_log.sql`

### Archivos a actualizar:

- `postgres/README.md` - Documentar nuevas tablas
- `CHANGELOG.md` - Agregar entrada de versión

---

## Criterios de Aceptación

✅ **Migraciones creadas correctamente**:
- Archivos SQL con sintaxis válida
- Constraints definidos correctamente
- Índices para performance

✅ **Migraciones ejecutadas sin errores**:
- En ambiente local (desarrollo)
- En ambiente dev (testing)

✅ **Verificaciones**:
- `\d user_active_context` muestra estructura correcta
- `\d user_favorites` muestra estructura correcta
- `\d user_activity_log` muestra estructura correcta
- Índices creados: `\di` muestra todos los índices
- Triggers funcionando correctamente

✅ **Documentación actualizada**:
- README con descripción de nuevas tablas
- CHANGELOG con nueva versión

---

## Migraciones Futuras (NO para esta fase)

Estas se implementarán cuando se trabaje la App de Administración (FASE 5):

| Tabla | Prioridad | Sprint | Descripción |
|-------|-----------|--------|-------------|
| `academic_cycles` | 🟡 Media | Admin Sprint 2 | Ciclos escolares/años académicos |
| `academic_periods` | 🟡 Media | Admin Sprint 2 | Trimestres, bimestres, semestres |
| `schedules` | 🟡 Media | Admin Sprint 3 | Horarios de clases |
| `schedule_blocks` | 🟡 Media | Admin Sprint 3 | Bloques de tiempo del horario |
| `classrooms` | 🟢 Baja | Admin Sprint 4 | Aulas físicas |
| `school_events` | 🟢 Baja | Admin Sprint 4 | Eventos escolares |
| `import_jobs` | 🟡 Media | Admin Sprint 3 | Trabajos de importación masiva |
| `grading_scales` | 🟢 Baja | Admin Sprint 5 | Escalas de calificación |
| `custom_roles` | 🟢 Baja | Admin Sprint 5 | Roles personalizados |
| `certificates` | 🟢 Baja | Admin Sprint 6 | Certificados de cursos |
| `fee_concepts` | 🟢 Baja | Admin Sprint 6 | Conceptos de cobro |
| `payments` | 🟢 Baja | Admin Sprint 6 | Pagos de estudiantes |

---

## Referencias

- **Plan de trabajo completo**: `/Users/jhoanmedina/source/EduGo/Analisys/docs/specs/ui-roadmap/PLAN-TRABAJO-ORDENADO.md`
- **Endpoints backend requeridos**: `/Users/jhoanmedina/source/EduGo/Analisys/docs/specs/ui-roadmap/ENDPOINTS-BACKEND-REQUERIDOS.md`
- **Convenciones de migraciones**: `postgres/README.md`

---

## Notas Técnicas

### Numeración de migraciones
- Última migración existente: `010_create_login_attempts.sql`
- Próximas migraciones: `011`, `012`, `013`

### Convenciones
- Usar `UUID` para PKs (gen_random_uuid())
- Usar `TIMESTAMP WITH TIME ZONE` para fechas
- Siempre incluir `created_at`
- Incluir `updated_at` si la tabla se actualiza
- Triggers para `updated_at` automático
- Comentarios con `COMMENT ON TABLE/COLUMN`

### Performance
- Índices en FKs para joins rápidos
- Índices en columnas de filtrado frecuente
- Considerar índices parciales (WHERE clause) para queries específicos
- JSONB para metadata flexible

### Validaciones
- Constraints a nivel de BD (UNIQUE, FK, NOT NULL)
- Validaciones de negocio en capa de API
- Considerar triggers para validaciones complejas si es necesario

---

## Checklist de Implementación

```
□ Crear migración 011_create_user_active_context.sql (structure)
□ Crear migración 011_create_user_active_context.sql (constraints)
□ Crear migración 012_create_user_favorites.sql (structure)
□ Crear migración 012_create_user_favorites.sql (constraints)
□ Crear migración 013_create_user_activity_log.sql (structure)
□ Crear migración 013_create_user_activity_log.sql (constraints)
□ Ejecutar migraciones en ambiente local
□ Verificar estructura de tablas con \d
□ Verificar índices con \di
□ Ejecutar migraciones en ambiente dev
□ Actualizar postgres/README.md
□ Actualizar CHANGELOG.md
□ Crear tag de versión (postgres/v0.11.0)
```

---

## Próximos Pasos (FASE 2)

Una vez completada esta fase, el siguiente paso es:

**FASE 2: API-Mobile**
- Implementar endpoints que consuman estas nuevas tablas
- Ubicación: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile`
- Endpoints: `/v1/users/me/schools`, `/v1/users/me/active-school`, etc.
