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

## Resolución en Fase 2 ✅

**Fecha de Resolución:** 22 de Noviembre, 2025  
**Opción Elegida:** Opción 1 - Crear Migraciones SQL

### Acciones Realizadas

1. ✅ **Análisis de proyectos hermanos:**
   - api-mobile: MaterialVersion, Progress
   - api-admin: Subject, Unit, GuardianRelation
   - Descubrimiento: AssessmentQuestion y AssessmentAnswer están en MongoDB

2. ✅ **Migraciones SQL creadas (5):**
   - 012_create_material_versions.up.sql + down.sql
   - 013_create_subjects.up.sql + down.sql
   - 014_create_units.up.sql + down.sql
   - 015_create_guardian_relations.up.sql + down.sql
   - 016_create_progress.up.sql + down.sql

3. ✅ **Entities PostgreSQL creadas (5):**
   - material_version.go
   - subject.go
   - unit.go
   - guardian_relation.go
   - progress.go

4. ✅ **Corrección de alcance:**
   - AssessmentQuestion y AssessmentAnswer NO son PostgreSQL
   - Están en MongoDB como parte de MaterialAssessment (ya creadas en Fase 1)
   - Alcance real: 13 entities PostgreSQL (no 14)

5. ✅ **Validación:**
   - Compilación postgres: exitosa
   - Compilación mongodb: exitosa
   - Tests: exitosos

### Resultado Final

- ✅ **13/13 entities PostgreSQL** (100%)
- ✅ **3/3 entities MongoDB** (100%)
- ✅ **16/16 entities totales** (100% del alcance real)
- ✅ **Documentado en:** `tracking/FASE-2-COMPLETE.md`
- ✅ **Commits:** 20564c7, d40c563

---

## Estado Final: RESUELTO ✅

Todas las entities bloqueadas han sido resueltas en Fase 2.

El sprint puede continuar a Fase 3 (Validación y PR).

---

**Generado por:** Claude Code - Fase 1
**Siguiente acción:** Verificar compilación de entities PostgreSQL creadas
