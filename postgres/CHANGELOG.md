# Changelog

Este changelog comienza la nueva serie documental del modulo `postgres`.

Los tags historicos del modulo siguen existiendo en Git. El ultimo tag observado en esta fase es `postgres/v0.61.0`, pero el detalle narrativo de versiones anteriores no fue reconstruido aqui.

## [Unreleased]

## [0.900.4] - 2026-06-08

Material maestro-detalle (material = tema + N archivos), tipo de pregunta `multiple_select`,
seeds de evaluaciones en playground `n4_evaluacion` y poda de row-action SDUI heredada.
`SchemaVersion` 3.50.0 → 3.54.0; `L4_SEED_VERSION` 1.47.0 → 1.50.0.

### Added

- Tabla `content.material_file`: relación 1:N entre `content.materials` (ahora «tema») y sus
  archivos adjuntos (`original_filename`, `file_key`, `file_url`, `mime_type`, `file_size_bytes`,
  `position`, `uploaded_by_membership_id`). Permite que un material agrupe N recursos descargables.
- Columna `content.materials.summary` (`text`, nullable): resumen en Markdown del material/tema,
  editable a mano por el docente.
- Eliminada la tabla `content.material_version` (reemplazada por `content.material_file` en el
  modelo maestro-detalle; EduGo no está en producción).
- Tipo de pregunta `multiple_select` en `assessment.question`: CHECK extendido
  `question_type IN ('multiple_choice','true_false','short_answer','essay','multiple_select')`
  con validación que `correct_answer` es un array JSON cuando `question_type='multiple_select'`.
  Solo disponible en authoring (toma pendiente, deuda 009/010).
- Playground v2 `n4_evaluacion` gana 2 evaluaciones «Sistema Solar» (draft + published) con los
  5 tipos de pregunta (`multiple_choice`, `true_false`, `short_answer`, `essay`,
  `multiple_select`): fixture completa para validar authoring y toma de evaluaciones.
  `L4_SEED_VERSION` 1.47.0 → 1.50.0.

### Changed

- Poda de la row-action SDUI `edit` heredada en `assessment-questions-list`: la acción era un
  fantasma (no mapea a ninguna pantalla nativa) y generaba un botón inoperante en la UI.
  `L4_SEED_VERSION` bump incluido.
- `SchemaVersion` 3.50.0 → 3.54.0 (material maestro-detalle F1: → 3.52.0;
  `multiple_select` F2: → 3.54.0).

## [0.900.3] - 2026-06-07

Cierre de N4 (evaluación/contenido sobre la sesión + notas con procedencia). Planes 015 /
ADR 0019 / ADR 0020. `SchemaVersion` 3.48.0 → 3.50.0; `L4_SEED_VERSION` 1.45.0 → 1.47.0.

### Added

- Tabla `academic.grade_item` (N4 / ADR 0020 — modo detallado): componente de nota anclado a
  `(membership_id, subject_id, period_id)` con autor (`created_by_membership_id`) y trazabilidad al
  origen automático (`source_attempt_id` → `assessment.assessment_attempt`, `source_assessment_id` →
  `assessment.assessment`, ambas FK `ON DELETE SET NULL`). FKs académicas a `memberships`/`subjects`/
  `academic_periods` y UNIQUE parcial `uq_grade_item_attempt (membership_id, subject_id, period_id,
  source_attempt_id) WHERE source_attempt_id IS NOT NULL` (un solo componente `auto_scored` por intento
  de origen; los manuales quedan fuera del índice). Trigger `set_updated_at`.
- Tabla `academic.grade_history` (N4 / ADR 0020): auditoría append-only de override de notas, con FKs a
  `grades`/`grade_item` (CASCADE) y al `changed_by_membership_id` (RESTRICT), y CHECK XOR
  `grade_history_target_xor_check` (cada fila apunta a EXACTAMENTE uno de `grade_id`/`grade_item_id`).
  Sin `updated_at`/trigger (es append-only; `changed_at` lo fija el insert).
- Columna `academic.grades.source` (`varchar(20)` NOT NULL, default `'manual'`, CHECK
  `source IN ('auto_scored','manual','auto_llm')`): procedencia de la nota unificada (N4 / ADR 0020).
- Columna `academic.schools.grade_profile` (`varchar(20)` NOT NULL, default `'basic'`, CHECK
  `grade_profile IN ('basic','detailed')`): perfil de notas de la escuela (básico/detallado).
- Recurso L4 `grades_detail` ("Desglose de Notas", no menú-visible, scope `unit`) + 4 permisos
  `academic.grades_detail.{create,read,update,delete}`: catálogo del modo detallado de notas. Recurso
  propio (no comparte `resource_id` con `grades` por el unique `(resource_id, action)`); el grant es
  condicional por `schools.grade_profile` y lo inyecta identity en runtime (no se otorga vía
  `roleGrantPatterns`).
- Playground v2 `n4_evaluacion` (N4 F5.1): escuela "Colegio N4 Evaluacion" con
  `grade_profile='detailed'`, oferta de Ciencias Naturales, docente, alumnas (Ana/Bruno/Caro) y período
  activo; registrado en `playground_v2.go`. Soporta el E2E de cierre de N4.

### Changed

- Reconstrucción del subsistema de evaluación/contenido sobre el modelo de sesión (N4 F1, ADR 0019).
  Demolición + recreación sin backfill (EduGo no está en producción). Las tablas viejas llaveadas a
  `auth.users` y a subject/grade texto-libre se reescriben ancladas a
  `memberships`/`subjects`/`subject_offerings`:
  - `assessment.assessment`: `created_by_user_id` → `created_by_membership_id` (RESTRICT), subject/grade
    texto → `subject_id` (→ `academic.subjects` RESTRICT), `school_id` NOT NULL (CASCADE), `status IN
    (draft,published,archived)`.
  - `assessment.question` / `question_option` (renombradas a singular): la opción correcta vive en
    `question.correct_answer` (sin `is_correct` en la opción).
  - `assessment.assessment_material`: N:N con PK compuesta `(assessment_id, material_id)` → `content.materials`.
  - `assessment.assessment_assignment`: puente a la sesión — se elimina `student_id` XOR
    `academic_unit_id`; target = `subject_offering_id` (→ `academic.subject_offerings` CASCADE) + UNIQUE
    `(assessment_id, subject_offering_id)`.
  - `assessment.assessment_attempt`: `student_id` → `student_membership_id`; UNIQUE parcial
    `(assessment_id, student_membership_id) WHERE status='in_progress'` (un solo intento activo).
  - `assessment.attempt_review`: `reviewer_id` → `reviewer_membership_id`.
  - `content.materials`: subject/grade texto → `subject_id` (SET NULL, nullable), `uploaded_by_teacher_id`
    → `uploaded_by_membership_id` (RESTRICT); `content.material_version`: `changed_by` →
    `changed_by_membership_id`; `content.progress`: PK `(material_id, user_id)` →
    `(material_id, student_membership_id)`.
  - ELIMINADAS las tablas analíticas viejas `assessment.attempt_analytics` y `assessment.assessment_stats`
    (llaveadas a `auth.users`; analítica diferida) y los índices de assignment por
    `student_id`/`academic_unit_id`. Todas las FKs cross-schema y los UNIQUE viven en `post_gorm.sql`
    (GORM no los materializa sin campo de relación). `content.courses` queda fuera de alcance (intacto).
  - Se podan los playgrounds muertos `focal_evaluacion` / `focal_evaluacion_v2` / `focal_botonera` /
    `focal_static_screens` y se sanea el seed demo de evaluación.
- Seed SDUI de evaluación alineado al contrato nuevo (N4 F2.6): `assessment-question-form` gana el field
  `option-list` (`correct_answer` por opción); `assessments-form` / assignment / listas alineadas a
  `content.assessments.*` y al esquema nuevo (`subject_id`, `subject_offering`); `assessment-modality`
  eliminada (concepto muerto). `L4_SEED_VERSION` 1.45.0 → 1.46.0.
- `SchemaVersion` 3.48.0 → 3.50.0 (F1: 3.48.0 → 3.49.0; F4a: 3.49.0 → 3.50.0).
- `L4_SEED_VERSION` 1.45.0 → 1.47.0 (F2.6: → 1.46.0; F4a catálogo modo detallado: → 1.47.0).

## [0.900.2] - 2026-06-06

### Added

- Campo `PeriodID uuid.UUID` (`period_id`, NOT NULL, indexado) en la entity exportada
  `entities.SubjectOfferingEnrollment`: copia denormalizada del período de la oferta que habilita
  queries por materia/período sin JOIN a `subject_offerings` y completa el invariante de inscripción.
  Lo consume `edugo-api-academic`.
- FK `grades_teacher_fkey` en `academic.grades.teacher_id` → `academic.memberships(id)`
  ON DELETE SET NULL (materializada en `post_gorm.sql`; `teacher_id` nullable, la nota persiste sin
  docente al expirar su membresía). Notas por sesión de materia (N3, plan 013).
- Seeds L4 de notas/asistencia por sesión (N3 / N3.5, planes 012/013/014): entry-points
  `take-attendance`, `view-attendance`, `view-attendance-summary` y `put-grades` reubicados como
  row-actions en la card de cada sesión (`sessions-by-subject-list`, scope `row`); acción
  `view-grades-summary` + instancia `grades-subject-summary` (resumen de notas por sesión del docente);
  instancia `grades-batch` (pantalla "Poner notas"). Recurso `my_grades` ("Mis Notas") + permiso
  `academic.my_grades.read:own` con grant al rol student + instancia `my-grades-list`. `L4_SEED_VERSION`
  1.42.9 → 1.45.0.

### Changed

- Invariante de inscripción ampliado a período: el uniqueIndex `uq_enrollment_student_subject` pasa de
  `(student_membership_id, subject_id)` a `(student_membership_id, subject_id, period_id)`. Un alumno no
  puede inscribirse dos veces en la misma materia dentro del mismo período (bug 0036), pero sí puede
  cursarla en períodos distintos.
- `SchemaVersion` 3.47.3 → 3.48.0.
- SDUI de notas (`grades-*`): `api_prefix` `learning → academic` — las calificaciones viven en la API
  academic, no learning (plan 012, bug 0034).
- Tightening de privacidad: se quita el wildcard `academic.grades.*` del rol student (recibía notas
  ajenas vía `GET /grades`); el alumno queda solo con `academic.my_grades.read:own` ("Mis Notas"). El rol
  guardian conserva `academic.grades.*` a propósito.

## [0.900.1] - 2026-06-05

### Added

- Invariante "una oferta por materia": `subject_offering_enrollments` gana `subject_id`
  (copia denormalizada e inmutable del `subject_id` de la oferta) + uniqueIndex
  `uq_enrollment_student_subject (student_membership_id, subject_id)` + FK. Sostiene el guard
  del usecase de inscripción que rechaza doble inscripción de un alumno en la misma materia
  (bug 0036, PRE 1b).

### Changed

- Asistencia token-scoped (tenant→JWT, plan 012, PRE 1a): el form `attendance-batch` pierde el
  campo tenant `unit_id` (la unidad se deriva del JWT) y se elimina el screen huérfano
  `attendance-form` (cierra 2 de los latentes del bug 0034).
- Seeds L4 de asistencia por sesión de materia (N2, plan 008): `api_prefix` `learning → academic`
  en las instancias `attendance-*`; entry-point `take-attendance` (pasar lista) y acciones
  `submit-batch` / `view-attendance` / `view-attendance-summary` en `subjects-form`.
- `SchemaVersion` 3.45.0 → 3.47.3; `L4_SEED_VERSION` 1.42.4 → 1.42.9.

## [0.900.0] - 2026-06-05

### Changed

- El form `sessions-by-subject-form` limita el campo `section_label` a 10 caracteres
  (atributo `max_length: 10` en su `slot_data`), alineado con la validación del backend
  (`section_label max=10`): el SDUI ahora previene la entrada de más de 10 caracteres en
  lugar de fallar el guardado con 400. Soporte de `max_length` añadido al SDUI del KMP
  (modelo `Slot.maxLength` + `FormFieldsResolver` + `SlotRenderer.applyMaxLength`).
- El form `units-form` gana el campo `type` (select required, options
  school/grade/class/section/club/department): el DTO `CreateUnitRequest` lo exige
  `binding:"required"` y sin él el backend respondía 400 al crear una unidad.
- El form `units-form` se sanea (plan 011 Eje C): se QUITAN los campos `level` y
  `period_id`, que el DTO `CreateUnitRequest` no acepta y el contrato KMP
  `UnitsFormContract` descartaba en silencio. El form queda con `name` + `type`. El
  `UnitsFormContract` del KMP migra además de `/api/v1/schools/{schoolId}/units` al
  endpoint JWT-only `/api/v1/units` (school_id del token, no del path; bug 0015).
- `L4_SEED_VERSION` 1.42.1 → 1.42.4. Sin cambios de esquema (`SchemaVersion` intacto) ni de permisos.

## [0.5.0] - 2026-06-03

### Changed

- `academic.academic_periods` gana la columna `academic_unit_id` (uuid, NOT NULL, con índice
  y FK a `academic.academic_units(id)` ON DELETE CASCADE), espejo de `school_id`: el período
  queda atado además a la unidad académica.
- El índice único parcial `idx_academic_periods_active` pasa de `(school_id)` a
  `(school_id, academic_unit_id) WHERE is_active = true`; la exclusividad del período activo
  ahora es por unidad, no por escuela.
- `SchemaVersion` 3.44.0 → 3.45.0. Seeds que insertan períodos (demo y playgrounds v2
  `n1_inscripcion` / `n17_secciones` / `multi_unidad`, más la fixture e2e `screen_only`)
  propagan `academic_unit_id`.

## [0.4.0] - 2026-06-02

### Added

- Playground v2 `n0n1_escuelas`: 3 escuelas, 13 docentes y una solicitud N0 pendiente
  (InglesLab), para validar el flujo de aprobación N0/N1.

### Changed

- `sessions-by-subject-form` pasa de scope `school` a `unit`; se elimina la pestaña *Alumnos*
  del detalle de materia y la instancia `studentsBySubjectList` asociada.
- Ajustes en los playgrounds v2 `multi_unidad`, `n17_secciones` y `n1_inscripcion`, y en la
  entidad `subject`.
- `L4_SEED_VERSION` 1.37.0 → 1.41.0; `SchemaVersion` 3.41.0 → 3.43.0; demo `SeedVersion`
  `development-gorm-v3` → `development-gorm-v4`.

### Removed

- `v2_screens_catalog`: el CRUD de los 4 recursos meta SDUI se migró al admin-tool de Go.

## [0.3.0] - 2026-05-30

### Added

- Playground v2 `multi_unidad` para validar el selector de unidad (multi-unidad).

### Changed

- Scope `unit` en los recursos `memberships`, `subjects` y `subject_offerings`
  (contexto requerido derivado del scope del recurso).
- `L4_SEED_VERSION` 1.35.0 → 1.37.0.

### Removed

- Poda del menú SDUI: se eliminan 6 recursos sin pantalla KMP asociada.

## [0.2.0] - 2026-05-28

### Added
- Seeds SDUI de N1.7 (`L4_SEED_VERSION` 1.35.0):
  - Pantallas nativas `batch-enroll` (inscripción por lote), `enroll-one` (inscripción 1-a-1) y `sessions-by-subject-list` (listado de sesiones por materia).
  - Entrada de menú "Sesiones de Materia".
- Master-detail generalizado a N detalles vía `detail_configs` (antes limitado a un único detalle).
- Playground v2 `n17_secciones`: secciones A/B y un docente con dos sesiones.

### Changed
- `L4_SEED_VERSION` bumpeado a `1.35.0`.

## [0.1.0] - 2026-05-27

### Added
- Reinicio de la versión a `v0.1.0` (borrón y cuenta nueva).
- Estructura limpia y alineación del esquema de base de datos relacional y seeds para PostgreSQL.

## [0.77.2] - 2026-03-30
### Changed
- fix seeds structure

## [0.77.1] - 2026-03-29
### Changed
- fix seeds structure

## [0.77.0] - 2026-03-29
### Changed
- fix seeds structure

## [0.76.0] - 2026-03-28
### Changed
- fix seeds production

## [0.75.0] - 2026-03-28
### Changed
- fix seeds UI

## [0.74.0] - 2026-03-28
### Changed
- fix seeds

## [0.73.0] - 2026-03-27
### Added
- fix seeds

## [0.72.0] - 2026-03-27
### Added
- Agregar assessments:assign y assessments:review a permisos del sistema
- Asignar ambos a super_admin, school_admin, school_director, school_coordinator, teacher
- Agregar 2 assessments manuales de ejemplo con preguntas PG, opciones y asignaciones
- Bump SchemaVersion 1.2.0 → 1.3.0

### Changed
- 003_permissions.sql — 2 permisos nuevos (68 total)
- 004_role_permissions.sql — 10 asignaciones nuevas (5 roles x 2 permisos)
- 008_assessments.sql — 2 assessments manuales + 6 preguntas PG + opciones + 3 asignaciones
- version.go — SchemaVersion 1.3.0

## [0.71.0] - 2026-03-27

### Added
- Tablas para el sistema de evaluaciones: `assessment.questions`, `assessment.question_options`, `assessment.assessment_assignments`.
- Tablas para el sistema de revisiones: `assessment.attempt_reviews`.
- Tablas para el sistema de notificaciones: `notifications.notifications`.
- Entidades GORM para todas las nuevas tablas: `Question`, `QuestionOption`, `AssessmentAssignment`, `AttemptReview`, `Notification`.

### Changed
- Modificada la tabla `assessment.assessment` para soportar `source_type` (manual/ai_generated) y hacer `mongo_document_id` opcional.
- Actualizada la tabla `assessment.assessment_attempt_answer` con referencia a `question_id` y `review_status`.

## [0.70.0] - 2026-03-25

### Added
- Nuevo rol `readonly_auditor` (scope: `school`) en seeds de produccion.
- 27 permisos de solo lectura asignados a `readonly_auditor` (todos los recursos del ecosistema).
- Permisos `context:browse_schools` y `context:browse_units` para `school_coordinator` (total: 43).
- Permisos `context:browse_schools` y `context:browse_units` para `school_assistant` (total: 15).
- Usuario de prueba `readonly@edugo.test` (U-21) en seeds de desarrollo.
- Membership `m028` para usuario ReadOnly en San Ignacio (rol: `readonly_auditor`).
- Asignacion de rol `ur27` para usuario ReadOnly en seeds de desarrollo.

### Changed
- `SchemaVersion` bumpeado de `1.1.0` a `1.1.4`.

## [0.66.0] - 2026-03-23

### Changed
- Renombrada la clave de permiso `system:settings` a `system_settings:settings` en seeds de produccion.
- Actualizada la instancia de pantalla `system-settings` para usar la nueva clave de permiso y la plantilla `settings-system-v1`.

### Added
- Nueva plantilla de pantalla `settings-system-v1` en seeds de produccion.

## [0.65.0] - 2026-03-20

### Added
- Nueva documentacion fase 1 del modulo.
- Documentacion fase 2 de integracion ecosistemica del modulo.
- Indice local en `docs/` con procesos, arquitectura e integracion.
- `Makefile` uniforme con `release-check` y wrappers de release.
- CLI `cmd/seed` para ejecutar seeds embebidos sin scripts externos.
- `internal/dbutil`: paquete interno con `BuildDBURL` y `EnvFirst`, compartido por `cmd/runner` y `cmd/seed`.
- `internal/sqlutil`: paquete interno con `IsEmptyOrComment`, compartido por `migrations` y `seeds`.

### Changed
- `README.md` reescrito para representar el estado actual del modulo y no la documentacion heredada.
- `cmd/runner` ahora usa migraciones y seeds embebidos en lugar de rutas obsoletas del filesystem.
- `migrations/embed.go` y `seeds/embed.go` simplifican su API publica: se eliminaron `GetScript`, `ListScripts` y `GetScriptsByLayer` (sin callers confirmados en todo el ecosistema).

### Removed
- `cmd/migrate/`: CLI legacy de migraciones incrementales (`up/down/status/create/force`). Era incompatible con el modelo actual de recreacion completa de schema y no tenia callers en ningun proyecto del ecosistema.
- Targets de Makefile: `migrate-up`, `migrate-down`, `migrate-status`, `migrate-create`.
