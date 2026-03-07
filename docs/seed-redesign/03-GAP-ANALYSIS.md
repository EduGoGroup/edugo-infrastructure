# Gap Analysis — Flujo Estudiante y Evaluaciones

**Fecha:** 2026-03-07

---

## GAPS IDENTIFICADOS

### GAP-1: No hay endpoint para que estudiantes listen evaluaciones disponibles

**Estado actual:**
- `GET /api/v1/assessments` requiere `assessments:read` (solo docentes/admins)
- Los estudiantes acceden a evaluaciones SOLO via material: `GET /api/v1/materials/{id}/assessment`
- No existe `GET /api/v1/student/assessments` ni similar

**Impacto:** Un estudiante no puede ver una lista de "evaluaciones pendientes" o "evaluaciones disponibles".

**Solucion propuesta (fuera de scope seeds, requiere codigo):**
- Opcion A: Agregar endpoint `GET /api/v1/assessments/available` con permiso `assessments:attempt`
  que filtre por school_id del estudiante, status=published, y ventana de disponibilidad
- Opcion B: Reusar `GET /api/v1/assessments` pero que el middleware permita `assessments:read`
  a estudiantes (ya lo tiene en 004_role_permissions.sql) y filtrar por status=published

**Nota:** El permiso `assessments:read` YA esta asignado al rol `student` en los seeds de produccion.
Esto significa que `GET /api/v1/assessments` DEBERIA funcionar para estudiantes, solo que
filtraria por school_id. El unico gap real es que no hay una pantalla SDUI dedicada.

---

### GAP-2: assessment-take no tiene navegacion desde la lista

**Estado actual:**
- Existe screen_instance `assessment-take` (screen_key, template_id, handler)
- Existe contrato `AssessmentTakeContract` en KMP
- PERO: no hay un custom event handler en `AssessmentsListContract` que navegue a `assessment-take`
- El contrato actual navega a `assessments-form` (vista de docente, no de estudiante)

**Solucion propuesta (requiere codigo):**
- En `AssessmentsListContract`: el SELECT_ITEM deberia detectar si el usuario es estudiante
  y navegar a `assessment-take` en vez de `assessments-form`
- Alternativa: crear un contrato separado `StudentAssessmentsListContract` con screenKey diferente

---

### GAP-3: MongoDB vacio para assessments de seed

**Estado actual:**
- Los 3 assessments del seed (008_assessments.sql) tienen `mongo_document_id` apuntando a
  documentos que NO existen en MongoDB
- Cuando se hace GET /assessments/{id}, el servicio intenta buscar en MongoDB y falla con 500

**Solucion (cubierta por estos seeds):**
- Crear seed de MongoDB con los documentos correspondientes
- Cada assessment publicado debe tener su documento en `material_assessment_worker` con
  preguntas, opciones, respuestas correctas, y puntos

---

### GAP-4: Redis credentials invalidas

**Estado actual:**
- `CACHE_REDIS_URL` en debug.json de Mobile API tiene password incorrecta
- Error: `WRONGPASS invalid username-password pair or user is disabled`

**Impacto:** La API funciona sin cache (degraded mode), pero performance es menor.

**Solucion:** Actualizar credenciales en `.zed/debug.json` de mobile API (fuera de scope seeds).

---

## NO SON GAPS (confirmados como funcionales)

1. El permiso `assessments:read` ya esta asignado al rol student -- OK
2. El permiso `assessments:attempt` ya esta asignado al rol student -- OK
3. El permiso `assessments:view_results` ya esta asignado al rol student -- OK
4. Los endpoints para tomar evaluacion existen:
   - `POST /api/v1/materials/{id}/assessment/attempts` (crea intento + scoring)
   - `GET /api/v1/attempts/{id}/results` (ver resultados)
   - `GET /api/v1/users/me/attempts` (mis intentos)
5. Las screen_instances de assessment estan correctamente configuradas en produccion
6. Los resource_screens mapean correctamente assessments a las 5 pantallas

---

## PLAN DE ACCION

### Fase 1: Seeds (AHORA — este documento)
- [x] Disenar usuarios multi-escuela
- [x] Disenar terminologia por tipo
- [ ] Crear scripts PostgreSQL desde cero (development seeds)
- [ ] Crear scripts MongoDB (assessment documents)
- [ ] Actualizar production seeds si es necesario (screen instances nuevas)

### Fase 2: Codigo (DESPUES — no en este scope)
- [ ] Agregar navegacion de estudiante a assessment-take
- [ ] Posiblemente crear StudentAssessmentsListContract
- [ ] Revisar si assessments-list necesita filtrar por status=published para estudiantes
- [ ] Corregir credenciales de Redis
