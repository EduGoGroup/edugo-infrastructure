# Entities Bloqueadas - Fase 1

**Fecha:** 2025-11-22
**Sprint:** Sprint Entities - Fase 1
**Razón:** Migraciones SQL no existen

---

## Entities Bloqueadas

Las siguientes entities mencionadas en `docs/cicd/sprints/SPRINT-ENTITIES.md` **NO pueden ser creadas** en Fase 1 porque **no existen migraciones SQL** correspondientes:

### PostgreSQL Entities sin Migraciones

| # | Entity | Tabla Esperada | Estado | Razón |
|---|--------|----------------|--------|-------|
| 1 | `MaterialVersion` | `material_versions` | ❌ Bloqueada | No existe migración `create_material_versions.up.sql` |
| 2 | `Subject` | `subjects` | ❌ Bloqueada | No existe migración `create_subjects.up.sql` |
| 3 | `Unit` | `units` | ❌ Bloqueada | No existe migración `create_units.up.sql` |
| 4 | `GuardianRelation` | `guardian_relations` | ❌ Bloqueada | No existe migración `create_guardian_relations.up.sql` |
| 5 | `AssessmentQuestion` | `assessment_questions` | ❌ Bloqueada | No existe migración `create_assessment_questions.up.sql` |
| 6 | `AssessmentAnswer` | `assessment_answers` | ❌ Bloqueada | No existe migración `create_assessment_answers.up.sql` |
| 7 | `Progress` | `progress` | ❌ Bloqueada | No existe migración `create_progress.up.sql` |

---

## Entities Creadas Exitosamente (8 de 14)

Las siguientes entities **SÍ fueron creadas** porque tienen migraciones SQL existentes:

| # | Entity | Archivo | Migración Base | Estado |
|---|--------|---------|----------------|--------|
| 1 | `User` | `postgres/entities/user.go` | `001_create_users.up.sql` | ✅ Creada |
| 2 | `School` | `postgres/entities/school.go` | `002_create_schools.up.sql` | ✅ Creada |
| 3 | `AcademicUnit` | `postgres/entities/academic_unit.go` | `003_create_academic_units.up.sql` | ✅ Creada |
| 4 | `Membership` | `postgres/entities/membership.go` | `004_create_memberships.up.sql` | ✅ Creada |
| 5 | `Material` | `postgres/entities/material.go` | `005_create_materials.up.sql` | ✅ Creada |
| 6 | `Assessment` | `postgres/entities/assessment.go` | `006_create_assessments.up.sql` + `009_extend_assessment_schema.up.sql` | ✅ Creada |
| 7 | `AssessmentAttempt` | `postgres/entities/assessment_attempt.go` | `007_create_assessment_attempts.up.sql` + `010_extend_assessment_attempt.up.sql` | ✅ Creada |
| 8 | `AssessmentAttemptAnswer` | `postgres/entities/assessment_attempt_answer.go` | `008_create_assessment_answers.up.sql` + `011_extend_assessment_answer.up.sql` | ✅ Creada |

---

## Decisión: Continuar sin Stubs

**Opción elegida:** No crear stubs para entities sin migraciones

**Razón:**
- Las entities son reflejos exactos de tablas SQL
- Sin migraciones, no podemos conocer el schema exacto
- Crear stubs inventaría estructura incorrecta
- Mejor esperar a que se creen las migraciones en Fase 2 o sprints futuros

**Alternativas consideradas:**
1. ❌ Crear stubs basados en nombres → No confiable sin schema real
2. ❌ Inferir estructura de APIs → Acoplamiento innecesario
3. ✅ **Documentar como bloqueadas y continuar** → Más seguro y transparente

---

## Impacto en el Sprint

### Objetivo Original
- Crear 14 entities PostgreSQL (según `SPRINT-ENTITIES.md`)

### Resultado Fase 1
- ✅ 8 entities creadas (57% del total)
- ❌ 6 entities bloqueadas (43% del total)

### Proyectos Afectados

**Proyectos que pueden usar entities creadas:**
- ✅ api-mobile: Puede usar User, School, AcademicUnit, Membership, Material, Assessment*
- ✅ api-administracion: Puede usar User, School, AcademicUnit, Membership
- ✅ worker: Puede usar todas las entities creadas

**Funcionalidades bloqueadas por entities faltantes:**
- ❌ api-mobile: MaterialVersion (versionado de materiales)
- ❌ api-administracion: Subject, Unit (organización curricular)
- ❌ api-administracion: GuardianRelation (gestión de apoderados)
- ❌ api-mobile: Progress (seguimiento de progreso estudiantil)
- ❌ api-mobile: AssessmentQuestion, AssessmentAnswer (modelo isolated de assessments)

---

## Plan para Fase 2

### Opción 1: Crear Migraciones SQL (Recomendado)

Si en Fase 2 tenemos acceso a un ambiente con PostgreSQL:

1. **Crear migraciones SQL** para las 6 entities faltantes basándose en:
   - Necesidades de api-mobile, api-administracion, worker
   - Diseños existentes en documentación de esos proyectos
   - Esquemas inferidos de código actual

2. **Ejecutar migraciones** en entorno de desarrollo/test

3. **Crear entities Go** reflejando las nuevas tablas

4. **Validar** con tests de integración

### Opción 2: Actualizar Sprint (Alternativa)

Si las migraciones no son prioritarias:

1. **Actualizar `SPRINT-ENTITIES.md`** para reflejar solo 8 entities
2. **Marcar como completado** con scope reducido
3. **Crear nuevo sprint** para las 6 entities faltantes cuando sea necesario

---

## Próximos Pasos

1. ✅ Continuar con MongoDB entities (si tienen migraciones/schemas)
2. ✅ Verificar compilación de entities PostgreSQL creadas
3. ✅ Crear tests para las 8 entities existentes
4. ✅ Documentar uso en READMEs
5. ⏳ **Fase 2:** Decidir entre Opción 1 o Opción 2 según prioridades del proyecto

---

**Generado por:** Claude Code - Fase 1
**Siguiente acción:** Verificar compilación de entities PostgreSQL creadas
