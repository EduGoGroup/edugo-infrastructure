# Changelog

Este changelog comienza la nueva serie documental del modulo `postgres`.

Los tags historicos del modulo siguen existiendo en Git. El ultimo tag observado en esta fase es `postgres/v0.61.0`, pero el detalle narrativo de versiones anteriores no fue reconstruido aqui.

## [Unreleased]

## [0.900.23] - 2026-07-15

Semilla SDUI: acción **«Revisiones»** en la toolbar del form de evaluaciones. `SchemaVersion` sin
cambios (**3.104.0**): no hay DDL, solo datos de `screenconfig`.

### Added

- **Acción «Revisiones»** en `assessments-form` (`system/l4/screen_instances_rows.go`): abre el
  tablero de revisión (`assessment-review-dashboard`, `event_id: view-reviews`) desde la evaluación,
  que hasta ahora solo se alcanzaba por el tablero «por calificar». Visible sobre evaluaciones
  `published`, permiso `content.assessments.read`, orden 25.

### Notas de despliegue

- **El `ExpectedContentHash` NO cambia.** Aunque combina migraciones y semillas, la parte de semillas
  se calcula sobre `nombre de capa + SeedVersion()` declarada (`seeds/version.go`), **no sobre el
  contenido de los archivos**. Este cambio no toca `L4_SEED_VERSION`, así que el hash queda en
  `9c21a3421fcf67ea` y `cloud-status` sigue verde. Verificado contra Neon con el migrator compilado
  contra esta versión.
- `ApplyScreenInstances` es `ON CONFLICT (screen_key) DO NOTHING` **a propósito** (no pisar
  customizaciones manuales en live): re-correr la semilla **no** actualiza una pantalla existente. En
  un environment ya sembrado, un cambio de `slot_data` se aplica con un UPDATE focalizado, o
  recreando. En un environment nuevo la semilla ya trae la acción.
- **Neon (staging) ya tenía esta acción aplicada a mano** (verificada campo por campo el 2026-07-15):
  este cambio es el código poniéndose al día para que un recreate desde cero la reproduzca. No hizo
  falta ningún UPDATE.

Consumidor real de la semilla: `edugo-dev-environment`.

## [0.900.22] - 2026-07-15

Planes 039 (terreno LLM) y 040 (corrección IA prevalidada). `SchemaVersion` 3.100.0 → **3.104.0**,
aplicado aditivamente (3.101.0 → 3.104.0). Introduce los cimientos de la revisión asistida por IA:
política por escuela, estado `ai_reviewed`, procedencia de la revisión, el evento
`attempt.review_requested` y el tablero de revisión de intentos en SDUI.

### Added

- **`identity.school_settings` + catálogo de claves** (3.101.0, plan 039): tabla de configuración por
  escuela con su catálogo de claves y seeds base; soporta la política por carril del terreno LLM
  (credenciales LLM viven por env del ecosistema; la escuela solo elige política).
- **`assessment.review_status = 'ai_reviewed'`** (3.102.0, plan 040 F0): nuevo estado de revisión para
  intentos corregidos por la IA (prevalidación), junto a los estados existentes.
- **`assessment.attempt_review.review_source`** (3.102.0, plan 040 F1): procedencia de la revisión
  (docente vs IA), para distinguir quién produjo la corrección.
- **Schema del evento `attempt.review_requested`** (3.102.0, plan 040 F1): contrato del evento que
  dispara la solicitud de revisión asistida (empareja con `edugo-shared/messaging/events`).
- **Template SDUI `review-dashboard-v1`** (3.103.0, plan 040 F3): plantilla dedicada al
  `assessment-review-dashboard` con chips de filtro (incl. «Prevalidado IA»).
- **Row-action «Revisar intentos»** (3.104.0, plan 040 T3c P2a): acción de fila en
  `assessments-management-list` que navega hacia `assessment-review-dashboard`.

## [0.900.21] - 2026-07-14

Planes 037 F1 (worker a dieta) y 038 Riel 0 (import externo de evaluaciones). `SchemaVersion`
3.98.0 → **3.100.0**, aplicado aditivamente (3.99.0 y 3.100.0). MongoDB sale del ecosistema (D-037.11).

### Added

- **`assessment.source_type = 'imported'`** (3.100.0, plan 038 Riel 0): el CHECK
  `assessment_source_type_check` admite el nuevo origen `imported` — evaluación creada desde un JSON
  externo (`assessment_import` v1) — junto a `manual` (UI) y `ai_generated` (LLM interno). El default
  sigue `manual`; el validador `oneof` se actualiza en paralelo.

### Removed

- **`academic.practice_result`** (3.99.0, plan 037 F1g): se elimina la tabla y su entity
  (`practice_result.go`), deprecada por el plan 036 (D-036.3). Era el espejo de `academic.grade_item`
  para evaluaciones de práctica; la trazabilidad de práctica vive ahora en el plano
  `assessment.practice_session` / `practice_session_answer` / `user_question_stat` (3.98.0). Se retira
  del AutoMigrate (`gorm_migrator.go`) y sus 3 bloques en `post_gorm.sql` (6 FKs, trigger
  `set_updated_at`, índice parcial `uq_practice_result_attempt`) → `ComputeFilesHash()` cambia.

### Changed

- **MongoDB fuera del ecosistema** (plan 037 D-037.11): limpieza de las referencias documentales en los
  comentarios de las entities. `assessment.mongo_document_id` queda como **columna legada reservada**
  (sin backing store documental); el `summary` de `material` se describe como resumen manual del docente,
  sin mención a Mongo. Sin cambio de schema.

## [0.900.20] - 2026-07-13

Planes 036 (plano de examen sano) y 035 F1 (capa de práctica). `SchemaVersion` 3.96.0 → **3.98.0**;
`L4_SEED_VERSION` con el field `purpose` + `visible_when` en `assessments-form`. Todo aplicado
aditivamente a Neon (3.97.0 y 3.98.0 estampados con hash combinado).

### Added

- **`assessment_attempt.teacher_feedback`** (3.97.0): comentario global del docente al finalizar la
  revisión (plan 036 D-036.4).
- **Capa de práctica** (3.98.0, plan 035 F1a): `assessment.purpose` (practice|exam|both, default exam)
  + `assessment.passing_score`; entities nuevas `practice_session`, `practice_session_answer`,
  `user_question_stat` (acumulador acotado, UNIQUE membership+question); backfill `purpose` ← `kind`.
- **Seed L4 `assessments-form`**: select `purpose` con `visible_when` en `max_attempts`/`passing_score`.

### Removed

- **`Assessment.Kind`** fuera de la entity (el seed/BD local nace sin `kind`; el DROP en Neon es la
  tarea 035-F1k, tras el deploy de learning sin referencias).

## [0.900.19] - 2026-07-13

Plan 032 (ola 2) — catálogo de evaluaciones. `SchemaVersion` 3.94.0 → **3.96.0**;
`L4_SEED_VERSION` 1.82.0 → **1.83.0**. Ambas aplicadas aditivamente a Neon.

### Added

- **`assessment.is_public`** (3.95.0): columna que habilita el catálogo de evaluaciones (visibilidad
  pública, autor mantiene propiedad).
- **Seed L4 — recurso `assessments-form`** (3.96.0): acciones `publish-catalog`/`unpublish-catalog`,
  con `visible_when` en lista combinado en AND. `L4_SEED_VERSION` **1.83.0**.

## [0.900.18] - 2026-07-12

Plan 033 Bloque B2a — etiquetas personales de materiales. `SchemaVersion` 3.93.0 → **3.94.0**;
`L4_SEED_VERSION` sin cambios (solo estructura, sin tocar seeds de datos).

### Added

- **`content.user_material_tags`**: etiquetas personales de materiales por usuario. Nueva entity
  `UserMaterialTag` con `UNIQUE(user_id, material_id, tag)` e índice por `user_id`, registrada en
  `AutoMigrate`. Aplicada aditivamente a Neon.

## [0.900.17] - 2026-07-03

Fix del bug 0081 (403 al abrir "Tomar Evaluación" como alumno). `SchemaVersion` 3.92.0 → **3.93.0**;
`L4_SEED_VERSION` 1.81.0 → **1.82.0**.

### Fixed

- **Seed L4 — recurso `assessments_student` ("Tomar Evaluación")**: la pantalla default era
  `assessments-list` (lista del docente, `GET /assessments`, permiso `content.assessments.read`) → 403
  al alumno. Se elimina ese mapping del recurso del estudiante; queda `assigned-assessments-list`
  (`GET /me/assigned-assessments`, permiso `content.assessments_student.read`) como única pantalla
  (`isDefault`). Bug 0081.

## [0.900.16] - 2026-06-30

Plan 033 Bloque B1b — biblioteca de materiales por grupos. `SchemaVersion` 3.90.0 → **3.92.0**;
`L4_SEED_VERSION` → **1.81.0**.

### Added

- **`content.material_assignment`**: distribución de materiales por grupo (`subject_offering_id`), calco
  de `assessment_assignment` (`material_id` → offering, `assigned_by_membership_id`,
  `available_from/until`, `UNIQUE(material_id, offering_id)`, FKs CASCADE/CASCADE/RESTRICT, trigger
  `updated_at`). Solo estructura; sin tocar seeds de datos.

### Fixed

- **Seed L4 — form de evaluación**: el campo Materia apuntaba al catálogo admin
  (`academic:/api/v1/subjects`, permiso podado del profesor → 403); ahora usa
  `academic:/api/v1/me/subjects` con `subtitle_field=code` (bug 0074).

## [0.900.15] - 2026-06-24

Plan 027 — permisología por proceso (altitud y arquetipos de acceso). Cierra fugas de ESCRITURA en roles
de consumo y mueve "ver lo mío" a recursos `my_*` dedicados. `SchemaVersion` 3.88.0 → **3.90.0**;
`L4_SEED_VERSION` → **1.80.0**. (Agrupa los commits `7fc5620` poda + recursos `my_*`, y `1afba1a` deny F4.8;
el release previo `v0.900.14` = 3.88.0 quedó sin documentar.)

### Added

- **Recursos `my_teaching` (profesor) y `my_attendance` (alumno)** con `read:own` — patrón "consume-lo-propio"
  espejo de `my_grades`/`my_memberships`/`my_wards_*`. Habilitan que profesor vea solo sus materias y alumno
  solo su asistencia, sin allows anchos sobre recursos compartidos.

### Changed

- **Poda de grants de consumo** (`l4/roles_permissions.go`): los wildcards anchos de los roles de consumo
  (alumno/representante) se recortan a `.read` literal, eliminando las fugas de escritura auditadas en F0
  (alumno podía `POST /attendance/batch` y crear/borrar anuncios/materiales; representante creaba/publicaba
  evaluaciones).
- **`school_admin` — deny F4.8** (`1afba1a`): se añade `deny` de `academic.*.read:own` (quita el ruido de los
  recursos `my_*` en el menú del admin) y de `admin.roles.{create,update,delete}` (el admin de escuela no crea
  roles IAM del contrato). Requiere el motor de menú deny-aware de `edugo-api-platform` (release `1920dce`).

> Detalle del plan: `../../../docs/plans/027-permisologia-por-proceso/` · cierre de despliegue:
> `../../../docs/plans/027-permisologia-por-proceso/despliegue.md`.

## [0.900.13] - 2026-06-19

Flag `is_system` para proteger el contrato de roles/permisos del sistema en runtime (deuda 0069).
`SchemaVersion` 3.80.0 → 3.81.0.

### Added

- **Columna `is_system` (`bool NOT NULL DEFAULT false`) en `iam.roles` e `iam.permissions`**: marca las
  entradas del contrato del sistema (sembradas por L0–L4) como inmutables en runtime. Cuando es `true`,
  los usecases de mutación (delete/update) de edugo-api-identity rechazan la operación para impedir que el
  catálogo de roles/permisos se borre o edite vía API; las entradas creadas por usuarios en runtime quedan
  en `false` (default). Se añaden índices parciales `idx_roles_system` / `idx_permissions_system` (vía tag
  GORM en las entities; `AutoMigrate` crea columna e índices). Los seeds L0–L4 marcan `is_system=true` en
  TODO rol/permiso del contrato (apply + accessors espejo, para que el cross-checker coincida):
  `L0_SEED_VERSION` 1.5.0 → 1.5.1, `L1_SEED_VERSION` 1.4.0 → 1.4.1, `L3_SEED_VERSION` 1.4.0 → 1.4.1,
  `L4_SEED_VERSION` 1.74.0 → 1.74.1 (L2 sin cambio). Requiere recrear BD (sin ALTER).

## [0.900.12] - 2026-06-18

Release previo no documentado en su momento (sincronización de seeds/merge dev).

## [0.900.11] - 2026-06-15

### Fixed
- **Mirror `L3ScreenInstances()`**: el accessor devolvía slice vacío desde la poda SDUI material
  (2026-06-07, v1.3.0), pero el seed real sí siembra la screen_instance mínima `materials-list`
  para satisfacer la FK `fk_resource_screens_screen_key`. El auditor `seed-audit` reportaba
  `RS_SCREEN_MISSING` en modo strict. Accessor actualizado para reflejar el seed real (v1.4.0).

## [0.900.10] - 2026-06-15

Eliminación completa del recurso/pantalla `progress` (`progress-dashboard`). `SchemaVersion`
3.70.0 → 3.71.0; `L4_SEED_VERSION` 1.66.0 → 1.67.0.

### Removed

- **Recurso L4 `progress` (…40) + su pantalla `progress-dashboard` (seed-only, sin DDL)**: la
  screen SDUI apuntaba a `/api/v1/stats/student` (endpoint inexistente → 404) y era redundante con
  el dashboard nativo del alumno; en el menú aparecía como "Reportes › Progreso" abriendo una
  pantalla vacía. Se eliminan: el recurso `progress` (`resources.go`, `resources_constants.go`), sus
  permisos `reports.progress.{read,read:own,update}` del catálogo y los grants `reports.progress.*` /
  `reports.progress.read:own` de los roles `student` y `guardian` (`roles_permissions.go`), la
  `screen_instance` `progress-dashboard` y su constante `L4_SCREEN_INST_PROGRESS_DASH_ID`
  (`screen_instances.go`, `screen_instances_constants.go`), y el mapping
  `progress → progress-dashboard` (`resource_screens.go`). UUIDs …40 (recurso) y …15
  (screen_instance) quedan libres. El recurso hermano `stats` (→ `stats-dashboard`,
  `/api/v1/stats/global`, vivo) y el padre `reports` se conservan intactos. Requiere recrear BD
  (sin ALTER).

## [0.900.9] - 2026-06-13

Poda del recurso/permisos `grades_detail` (plan 022 / ADR 0024 foco 3). `SchemaVersion` 3.65.0 → 3.66.0;
`L4_SEED_VERSION` 1.62.0 → 1.63.0.

### Removed

- **Recurso L4 `grades_detail` (…37) + sus 4 permisos `academic.grades_detail.{create,read,update,delete}`**
  (seed-only, sin DDL): se eliminan del catálogo L4 (`resources.go`, `resources_constants.go`,
  `roles_permissions.go`). El modo detallado de notas ya no se gobierna con un permiso: ahora lo decide
  `academic` leyendo `grade_profile` de la escuela. El permiso era un mensajero eliminable, así que se
  retira también el grant condicional por perfil que vivía en `identity`. UUID …37 queda libre. Requiere
  recrear BD (sin ALTER).

## [0.900.8] - 2026-06-13

### Fixed

- **Seed demo (`seedMemberships`)**: migrado de la columna `role` (eliminada por MP-08 F3) a
  `invitation_type_id` (FK → `academic.invitation_types`), alineando el demo con el esquema 3.64.0. El
  defecto hacía fallar la siembra demo de `postgres/v0.900.7` con `column "role" ... does not exist`.
  Las keys (`student`/`teacher`/`admin`/`coordinator`/`guardian`/`assistant`) se resuelven a su id vía
  `catalog.ResolveInvitationTypeID` (data-driven, sin UUIDs hardcodeados), igual que los playground_v2 y
  los playground legacy ya migrados. Sin cambio de esquema → `SchemaVersion` sigue 3.64.0.

## [0.900.7] - 2026-06-13

Cierre de varios micro-planes sobre el esquema: MP-08 (acceso por sistema + contexto de invitación
data-driven, swap `role`→`invitation_type_id`, split del permiso de aprobación), poda de features
muertas (MP-04 Track A / MP-01), rangos numéricos declarativos en forms SDUI (MP-03 F3), landing
data-driven (ADR 0024 F0) y persistencia de tenant en la notificación in-app. `SchemaVersion`
3.59.0 → 3.64.0; `L4_SEED_VERSION` 1.57.0 → 1.62.0.

### Added

- **MP-08 F0 (3.61.0, aditivo, solo esquema)**: 4 entities nuevas que modelan en datos el acceso por
  sistema y la equivalencia tipo-de-invitación → rol (todo por FK de id, nunca por nombre):
  `iam.systems` (catálogo de apps), `iam.system_roles` (puente sistema↔rol),
  `academic.invitation_types` (catálogo global de tipos de invitación) y
  `academic.school_invitation_roles` (equivalencia `(escuela, tipo) → rol` IAM, FK cross-schema
  `academic→iam`). Las 4 entran en AutoMigrate; `post_gorm.sql` agrega sus FKs y triggers
  `set_updated_at`. Sin seeds (los catálogos los siembra F1).
- **MP-08 F1 (seeds L4, `L4_SEED_VERSION` → 1.60.0)**: siembra de `systems`/`system_roles`/
  `invitation_types` y las equivalencias por escuela en las nuevas tablas del catálogo.

### Changed

- **MP-08 F3 (3.62.0, swap de columna, NO aditivo)**: la columna `role` (varchar + CHECK inline del
  enum de roles) se reemplaza por `invitation_type_id` (uuid, FK → `academic.invitation_types(id)`) en
  `academic.{memberships,school_invitations,school_join_requests}`. Se elimina el CHECK `..._role_check`
  (la validez la garantiza la FK) y el índice parcial `idx_memberships_unit_role_active` se reexpresa a
  `idx_memberships_unit_invitation_type_active` (en `post_gorm.sql`). Los seeds resuelven la key del tipo
  a su id vía `catalog.ResolveInvitationTypeID` (data-driven, sin UUIDs hardcodeados). Requiere recrear
  BD (sin ALTER).
- **MP-08 F4 (3.63.0, seed-only)**: el form `invitations-form` cambia el campo `role` (select estático,
  enum legacy muerto) por `invitation_type` (remote_select contra
  `GET /api/v1/schools/invitation-types`); y `schools-list` retira la acción `create` del header
  (`actions_removed ["create"]`) — el alta de escuelas pasa al admin-tool de Go (se conserva
  `schools-form` + `manage-concepts` y la edición de escuelas existentes). `L4_SEED_VERSION` → 1.61.0.
- **MP-08 aprobación SELLO x TIPO (3.64.0, seed-only)**: el permiso único
  `academic.join_request_approvals.<tipo>` se separa en
  `academic.join_request_approvals.{school,unit}.{student,teacher,guardian}` (catálogo 3 → 6 filas). El
  CHECK de permisos/grants se amplía de `{0,2}` a `{0,3}` segmentos (alinea con
  `enum.PathPermissionRegex` del shared, ahora 4 segmentos) en `permission.go` + `post_gorm.sql`. El
  grant de `teacher` pasa de `...student` a `...unit.student` (admite alumnos a su clase = sello de
  unidad, ya no firma el sello de colegio); `school_admin`/`super_admin` cubren ambos sub-namespaces por
  subárbol y `readonly_auditor` sigue denegado por su deny de prefijo. `L4_SEED_VERSION` → 1.62.0.
  Requiere recrear BD (sin ALTER).
- **ADR 0024 F0 (3.60.8)**: landing data-driven — `landing_screen_key` en roles y
  `default_landing_screen_key` en schools (la pantalla de aterrizaje deja de estar hardcodeada).
- **MP-03 F3 (3.60.6, seed-only)**: rangos numéricos declarativos en los forms SDUI. Cada campo
  `"type": "number"` lleva ahora `min`/`max` en su slot_data para que el front KMP valide antes de
  enviar, espejando el binding real del backend (assessments-form: `pass_threshold` 0–100,
  `max_attempts`/`time_limit_minutes` min=1; assessment-question-form: `points` min=0) y con un mínimo
  conservador donde el backend no declara rango. `L4_SEED_VERSION` → 1.58.0.
- **Notificación in-app (3.60.0, plan 020 F4.6.8)**: la entity `Notification` suma `school_id`/`unit_id`
  (uuid nullable) a `notifications.notifications` para que la lista in-app resuelva el context-switch
  multi-tenant al tocar (antes solo viajaban en el push). DDL aditivo vía AutoMigrate (2 columnas
  nullable, sin índice).
- **Seeds audit-detail (3.60.2/3.60.4/3.60.5, seed-only — MP-01 Ola 2)**: entry-point
  "Gestionar Conceptos" (`manage-concepts`) en `schools-form`, y nuevo template L4 `audit-detail-v1`
  que pinta los campos reales del evento de auditoría (actor, acción, recurso, status, severidad…) en
  solo lectura, en vez de los campos heredados de material/archivo. `L4_SEED_VERSION` 1.54.0 → 1.57.0.

### Removed

- **MP-04 Track A / MP-01 — poda de features muertas** (sin DROP; AutoMigrate nunca dropea, un recreate
  fresco simplemente ya no crea las tablas):
  - `content.courses` (3.60.1): se borra la entity `Course`, su registro en AutoMigrate y el `seedCourses`.
  - `content.progress` (3.60.3): se borra la entity `Progress`, su registro en AutoMigrate, sus 2 FKs y
    el trigger `set_updated_at` en `post_gorm.sql`.
  - `academic.{schedule,calendar_event,colors}` y el playground `focal_colors_demo` (3.60.7, MP-01 F3).
  Acompaña dedup de helpers de playground v2 → `playground_v2/common` (MP-01 F2.1) y la baja de los
  permisos/recursos muertos asociados.

## [0.900.6] - 2026-06-11

### Added

- Plan 020 N5 F1.1/F1.2: tablas `notifications.device_tokens` y `auth.service_clients`
  (entities GORM + índices parciales). Capa de seed L5-m2m con clientes `edugo-worker` y
  `edugo-api-learning` (scope `notifications.dispatch`, `secret_hash` bcrypt del dev secret de
  `push-secrets.env`). `SchemaVersion` 3.58.0 → 3.59.0; `L5_SEED_VERSION` 1.0.0.

## [0.900.5] - 2026-06-09

Despliegue del plan 019: `entity-picker` en `assessments-form.subject_id`, `view_when` read-only
fuera de borrador (ADR 0022), baja de los forms SDUI legacy `grades-form`/`user-roles` y fixture de
conformidad del contrato entity-picker. `SchemaVersion` 3.54.0 → 3.58.0; `L4_SEED_VERSION`
1.50.0 → 1.54.0.

### Changed

- Plan 017 F2 (picker de entidad): `assessments-form` migra el campo `subject_id` de
  `remote_select` a `entity-picker` (control nuevo). El selector de materia abre un modal con
  búsqueda server-side + paginación contra `academic:/api/v1/subjects` (`search_param=search`,
  `page_size=20`) en lugar de cargar todas las opciones al montar. Se conservan
  `remote_endpoint`/`display_field`/`value_field` (claves legacy con fallback en el resolver KMP
  `FormFieldsResolver`). Seed-only (sin DDL). `SchemaVersion` 3.55.0 → 3.56.0; `L4_SEED_VERSION`
  1.51.0 → 1.52.0.

- ADR 0022 (campos estructurales solo en borrador): `assessments-form` declara `view_when`
  (`{"field":"status","in":["published","archived"]}`) a nivel `slot_data`; el front pone el form en
  read-only total fuera de borrador. Acompaña el backend `learning`: el update persiste `subject_id`
  solo en borrador y rechaza el update fuera de borrador con 400 `BUSINESS_ASSESSMENT_NOT_DRAFT`;
  `AssessmentResponse` añade `subject_name` (subquery a `academic.subjects`) para el label del picker.
  Seed-only (sin DDL). `SchemaVersion` 3.56.0 → 3.57.0; `L4_SEED_VERSION` 1.52.0 → 1.53.0.

### Removed

- Poda de dos pantallas SDUI legacy huérfanas: (1) `grades-form` (reemplazada por las nativas
  `my-grade-detail`/`grades-batch`) y (2) `user-roles` (huérfana, sin reemplazo ni entry-point).
  Ambas tenían controles `remote_select` MUERTOS (`student_id`/`subject_id`, `user_id`). Se eliminan
  sus `screen_instances` + mappings en `resource_screens` + constantes. Sin cambios de roles ni
  permisos. Seed-only (sin DDL). `SchemaVersion` 3.57.0 → 3.58.0; `L4_SEED_VERSION` 1.53.0 → 1.54.0.

### Tests

- Fixture de conformidad del contrato `entity-picker` (plan 019 WI-5): test de validación del
  `slot_data` sembrado para `assessments-form.subject_id`, asegurando que el control migrado cumple
  el contrato esperado por el front.
- `gofmt` sobre `migrations/version.go` para pasar el `fmt-check` de CI (sin cambios semánticos).

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
