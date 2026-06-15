package migrations

import (
	"crypto/sha256"
	"fmt"
)

// SchemaVersion es la version actual de los scripts de migracion y seeds.
//
// OBLIGATORIO: Incrementar este valor cada vez que se modifique
// cualquier archivo en sql/*.sql o en seeds/.
// El migrador valida que esta version coincida con la registrada en BD.
//
// Historial reciente:
//   - 3.4.0: cierre Fase 1 rebuild (framework Layer + rename
//     production/development → system/demo).
//   - 3.5.0: cierre Fase 2 rebuild (capa L0 mínima + desactivación
//     Layer_Legacy + limpieza artefactos E2E legacy; ADR-6).
//   - 3.6.0: cierre Fase 3 rebuild (capa L1 readonly: rol
//     announcement_viewer + escuela mínima; ADR-7).
//   - 3.7.0: cierre Fase 4 rebuild (capa L2: segunda pantalla
//     announcement-form + mapping resource_screen tipo form).
//   - 3.8.0: cierre Fase 5 rebuild (capa L3: recurso materials con CRUD
//     parcial sin delete + 2 pantallas + aislamiento de menú).
//   - 3.8.1: fix Opción A (validación HTTP/UI): L1 ahora puebla
//     academic.memberships para destrabar switch-context del viewer.
//     L1_SEED_VERSION bumped a 1.1.0.
//   - 3.9.0: cierre Fase 6 rebuild — capa L4 completa (sistema completo
//     reorganizado por dominio: 31 recursos, 5 roles nuevos, ~89
//     permisos, 178 role_permissions, 5 templates nuevos + fix `zones`
//     en los 3 templates base L0, 68 screen_instances, 65
//     resource_screens, 5 concept_types + 50 definiciones), borrado
//     físico de `seeds/system/legacy/`, accessors públicos L4 para el
//     cross-checker, scenario `l4_full` con matriz role→screens
//     programática, baselines post-L4 archivados. Tickets de tooling
//     pendientes (TC-1..TC-5) documentados en decisions-log.
//   - 3.10.0: rescate Fase 6 — agregados 6 roles alias L4 (school_director,
//     school_coordinator, school_assistant, assistant_teacher, observer,
//     readonly_auditor con filtro readonly) + 7 permisos faltantes
//     (roles:create/update/delete, permissions_mgmt:create,
//     concept_types:create/update, attendance:update) +
//     screen_templates:* asignados a platform_admin + dashboard-home
//     shell + accessors L0..L3 (TC-5 cerrado) + demo seed y fixtures E2E
//     refactorizados con constantes L0/L1/L4. make seed-audit-strict
//     pasa con exit 0.
//   - 3.11.0: Fase 7 R-F7-1 — agregados 2 permisos system_settings:read
//     y system_settings:update al seed L4, con grants explícitos a 10
//     roles (read) y 6 roles (update). L4_SEED_VERSION bumped a 1.2.0.
//   - 3.12.0: super_admin gana context:browse_schools y context:browse_units
//     en L4 para destrabar SchoolSelector cuando login.schools[] viene
//     vacío. L4_SEED_VERSION bumped a 1.3.0.
//   - 3.13.0: nueva fixture E2E `global_user_no_membership` + scenario
//     `super_admin_global_flow` para el test cross-API en
//     edugo-dev-environment que valida SchoolSelector → switchContext →
//     UnitSelector → Dashboard del super_admin sin membership. Cambio
//     bajo seeds/e2e/* — bump obligatorio por regla CLAUDE.md.
//   - 3.14.0: PRE-4 (permissions-redesign) — agregadas 3 constantes al
//     enum `edugo-shared/common/types/enum/permission.go`
//     (PermissionAttendanceUpdate, PermissionSystemSettingsRead,
//     PermissionSystemSettingsUpdate) y fix de 4 strings buggy en KMP
//     (academic_units:* → units:*, user_roles:* → memberships:*). Las
//     filas de catálogo en el seed L4 ya existían; este bump cubre el
//     enum BE para que `IsValid()` reconozca los strings nuevos.
//   - 3.15.0: P1-1 (permissions-redesign Pass 1) — schema migration:
//     nuevas tablas iam.role_grants e iam.user_grants, funciones SQL
//     iam.permission_matches() e iam.scope_covers(), columna
//     auth.users.token_version (default 1) y columna nullable
//     iam.user_roles.scope_pattern. Sin datos aún; backfill en P1-2.
//   - 3.16.0: P1-2 (permissions-redesign Pass 1) — backfill 1:1 de
//     iam.role_grants desde iam.role_permissions (effect='allow', IDs
//     determinísticos vía SHA1) + backfill de iam.user_roles.scope_pattern
//     desde school_id/academic_unit_id. L4_SEED_VERSION bumped a 1.5.0.
//   - 3.17.0: P1-2 (permissions-redesign Pass 1) — ajuste de regex en
//     CHECK constraints role_grants_pattern_format y
//     user_grants_permission_format para aceptar ambos formatos: legacy
//     (recurso:accion[:own]) y path-based (recurso.accion[.*][:own]).
//     Necesario para mirror 1:1 desde role_permissions cuyo catálogo
//     legacy usa `:` como separador. El rename al formato nuevo está
//     planificado para Pass 2.
//   - 3.18.0: SUB-2 (post Pass 1 cleanup) — colapso de readonly_auditor
//     en seed L4: eliminados verbos mutativos (create/update/delete/publish/
//     finalize/grade/attempt/activate/approve/assign/review) del rol.
//     Adelanto del Pass 3 (D4) a Pass 1 por decisión del usuario, dado que
//     EduGo no está en producción. El mirror role_grants se actualiza
//     automáticamente vía applyL4RoleGrantsMirror. L4_SEED_VERSION → 1.6.0.
//   - 3.19.0: P2-4b (permissions-redesign Pass 2) — rename masivo del
//     catálogo iam.permissions.name del formato legacy `resource:action`
//     a path-based `<dominio_menu>.<recurso>.<accion>[:own]` (D3).
//     Afecta L0 (announcements: 4), L3 (materials: 3) y L4 (101 entradas).
//     El mirror iam.role_grants.pattern se regenera automáticamente vía
//     applyL4RoleGrantsMirror leyendo el `name` nuevo. slot_data.actions[]
//     .permission y screen_instances.required_permission también
//     renombrados. La tabla autoritativa de mapeo vive en
//     edugo-shared/common/types/enum/permission_path.go (legacyToPathTable).
//     L4_SEED_VERSION → 1.7.0.
//   - 3.20.0: extensión wildcard-first del PermissionMatcher para soportar
//     patterns `*.suffix` y `prefix.*.suffix`. Cambios:
//   - `iam.permission_matches()` extendida con dos ramas nuevas
//     (`*.suffix` y `prefix.*.suffix`).
//   - CHECK constraints `role_grants_pattern_format` y
//     `user_grants_permission_format` aceptan los nuevos patterns.
//   - Mirror cross-language: Go (`auth.PermissionMatches`), Kotlin
//     (`PermissionMatcher`) y el regex `enum.PathPermissionRegex`
//     extendidos en paralelo. Golden vector actualizado con casos
//     M33-M44 y G29-G35 para cubrir la nueva semántica.
//   - 3.21.0: P4-1 (plan B) — eliminación greenfield de iam.role_permissions.
//     La tabla legacy de asignaciones 1:1 rol×permiso se quita del schema
//     (entity, AutoMigrate, seeds, accessors, fixtures e2e, validators y
//     handlers HTTP). Los permisos efectivos por rol se resuelven sólo
//     vía iam.role_grants (patterns wildcard) + iam.user_grants.
//     Cambios principales:
//   - postgres/entities/role_permission.go: borrado.
//   - migrations/gorm_migrator.go: removido del AutoMigrate.
//   - sql/post_gorm.sql: borradas iam.get_user_permissions,
//     iam.get_user_resources, iam.user_has_permission (3 funciones
//     plpgsql sin consumers Go).
//   - seeds/system/layers/{l0_roles,l1_role_permission,l3_role_permissions}
//     y sus accessors: refactor para no sembrar role_permissions.
//   - seeds/system/l4/roles_permissions.go: eliminado bloque B2
//     legacy (rolePermissionGrants, applyL4RolePermissions y
//     helpers de filtrado readonly_auditor); applyL4RoleGrants
//     intacto.
//   - seeds/e2e/fixtures/{readonly_role,menu_subtree,partial_crud}:
//     re-cableados para sembrar iam.role_grants con patterns en vez
//     de iam.role_permissions.
//   - edugo-api-identity: removidos handlers HTTP de /roles/:id/
//     permissions (GET/POST/PUT/DELETE), use cases, mappers,
//     entities, ports.
//   - edugo-dev-environment/migrator: removido el validator
//     seedaudit/role_permissions; loader.SeedSnapshot pierde el
//     slice RolePermissions; production_loader.Snapshot.RolePermissions
//     queda en nil hasta P4-2.
//     L4_SEED_VERSION → 1.10.0; L0_SEED_VERSION sin bump (sólo se quita
//     una función auxiliar, no cambian datos visibles sembrados).
//   - 3.22.0: P4-2 (permissions-redesign Pass 4) — demo seed inicia
//     iam.user_grants productivo. Dos overrides ejemplo: deny puntual a
//     est.carlos sobre academic.grades.read (demo deny > allow) y allow
//     temporal a prof.martinez sobre admin.users.create con
//     expires_at=2026-06-11 (demo TTL). truncateDevelopmentData incorpora
//     iam.user_grants antes de iam.user_roles. demo.SeedVersion →
//     development-gorm-v2.
//   - 3.23.0: announcement-form gana el field `scope` (select
//     school|unit, required, default=school) en su slot_data L2. El form
//     ya enviaba title/body/published_at pero el DTO de backend exige
//     scope (oneof=school unit), por lo que el POST /api/v1/announcements
//     devolvía 400 al guardar desde el emulador. Cambios:
//   - seeds/system/layers/l2_screens.go: 4 fields en lugar de 3.
//   - seeds/e2e/fixtures/l2_constants_export.go: validador
//     actualizado (want 4 + key scope).
//     L2_SEED_VERSION → 1.1.0.
//   - 3.24.0: announcements-list slot_data agrega filter_ready_label
//     ("Fijados") y filter_processing_label ("No fijados") para
//     overridear los defaults "Activos"/"Otros" del template
//     list-basic-v1. La entidad Announcement no tiene is_active; los
//     labels genéricos eran engañosos en este contexto. L0_SEED_VERSION
//     → 1.2.0.
//   - 3.25.0: pulido de UI base — announcements-list gana
//     page_title="Anuncios" (antes el TopBar quedaba sin título por usar
//     "title" en vez de "page_title"); el template form-basic-v1 elimina
//     la zona form_header redundante con el TopBar. L0_SEED_VERSION →
//     1.3.0. Cambios complementarios en KMP: la acción `create` de las
//     listas se renderiza como FAB en lugar de IconButton en el TopBar.
//   - 3.26.0: announcement-form desdobla el botón Guardar en dos slots:
//     save_new (create-only, permission=create) y save (edit-only,
//     permission=update). Antes el único slot pedía `update` siempre,
//     lo que ocultaba el botón a usuarios con solo `create` (caso
//     focal-author). L2_SEED_VERSION → 1.2.0.
//   - 3.27.0: 4 slot_data de assessments (assessments-form,
//     assessments-management-list, assessment-questions-list,
//     assessment-question-form) migrados al nuevo estándar SDUI
//     consolidado en anuncios: zonas vacías con scope expandidas
//     por el resolver desde slot_data.actions[], page_title/edit_title
//     para el TopBar, save_new+save desdoblados con condition
//     create-only/edit-only y permisos diferenciados (create vs
//     update), destructive=true en eliminar. assessment-question-form
//     reemplaza los fields legacy (statement/kind/score/options) por
//     los del DTO real (question_text/question_type/points/
//     correct_answer/explanation/difficulty). assessments-form gana
//     los fields completos de CreateAssessmentRequest (pass_threshold,
//     max_attempts, time_limit_minutes, is_timed, shuffle_questions,
//     show_correct_answers, available_from/until). Además, en
//     resource_screens.go el recurso `assessments` cambia su default:
//     ahora `assessments-management-list` es is_default=true y
//     `assessments-list` queda en false. `assessments_student` queda
//     intacto (sigue con `assessments-list` como default).
//     L4_SEED_VERSION → 1.12.0.
//   - 3.28.0: fix de routing del menú docente para evaluaciones.
//     El bundle de sync expone `screens` indexadas por screen_type,
//     y el KMP toma `screens["list"]` para navegar — ignora
//     is_default. Mi cambio anterior a is_default no surtía efecto;
//     el menú seguía abriendo `assessments-list` (student-take) en
//     lugar de `assessments-management-list` (master-detail CRUD).
//     Fix en resource_screens.go: bajo el recurso `assessments`,
//     `screen_type="list"` ahora apunta directamente a
//     `assessments-management-list`. La fila legacy que mapeaba
//     `assessments-list` al recurso docente se eliminó —
//     `assessments-list` queda solo bajo `assessments_student`.
//     L4_SEED_VERSION → 1.13.0.
//   - 3.29.0: reincorporar actions perdidas en la migración del seed
//     legacy al rebuild L0-L4. assessments-form solo tenía
//     save_new/save/delete tras la migración Go; el botón "Preguntas"
//     (commit 15b3edc, marzo 2026) y los flujos publish/archive/assign
//     habían quedado fuera del slot_data, aunque los handlers KMP
//     (AssessmentFormContract) siguen vivos. Re-sembrados como actions
//     con scope=form y condition=edit-only.
//     L4_SEED_VERSION → 1.14.0.
//   - 3.30.0: rollback de 3.29.0 (Fase 0 del plan arquitectónico de
//     actions/iconos). Las 4 actions extra (view_questions/publish/
//     archive/assign) y los parches visuales SDUI en KMP
//     (SlotBindingResolver scope=form → ICON_BUTTON, SlotRenderer
//     tint destructive hardcoded, DSIconButton tint param) se
//     revierten porque eran parches: contrato icon-name sin
//     validación, sin overflow strategy en zonas ACTION_GROUP, style
//     hardcoded por if, mezcla semántica form-submit vs
//     resource-toolbar. assessments-form vuelve a save_new + save +
//     delete. Reincorporación de las acciones queda bloqueada hasta
//     definir tabla style→token, overflow strategy declarativa por
//     zona, contrato icon-name validado en build, y separación
//     scope=form-submit vs scope=resource-toolbar. Snapshot 002
//     anota el plan completo. L4_SEED_VERSION → 1.15.0.
//   - 3.31.0: composer SDUI defaults+added/removed + master-detail-v1
//     template + scope split form-submit/resource-toolbar + actions
//     reincorporadas a assessments-form sobre base arquitectónica del
//     snapshot 002. assessmentsForm pasa de form-basic-v1 a
//     master-detail-v1 con detail_config apuntando a
//     assessment-questions-list + actions_added [detail, publish,
//     archive] con scope=resource-toolbar. Defaults del template
//     aplican save_new/save/delete con scope=form-submit. Composer en
//     api-platform resuelve $resource$ placeholder, hace add/remove
//     sobre defaults, y reinjecta como slot_data.actions para el FE.
//     Retrocompat: instancias con "actions:[...]" legacy sin
//     added/removed se tratan como override total (announcement-form,
//     users-form, etc. siguen iguales). L4_SEED_VERSION → 1.16.0,
//     L0_SEED_VERSION → 1.4.0.
//   - 3.32.0: Fase 3 SDUI (F3-REQ-4 / bloque 7a) — nueva tabla
//     `academic.colors` (recurso CRUD plano de demostración usado por la
//     pantalla colors-list / colors-form que se resuelve vía
//     GenericListContract / GenericFormContract sin Kotlin nuevo).
//     AutoMigrate del entity Color; post_gorm.sql agrega trigger
//     set_updated_at, CHECK regex `^#[0-9A-Fa-f]{6}$` sobre `hex` y
//     UNIQUE compuesto (school_id, name). Las filas seed
//     (screen_instances, menu, permisos) son del bloque 7b.
//   - 3.33.0: nuevo registry `seeds/playground_v2/` paralelo a
//     `seeds/playground/`. MigrateOptions gana el campo `PlaygroundV2`
//     (mutuamente excluyente con `Playground`). El flujo Migrate
//     despacha a playground_v2.Apply cuando viene seteado.
//     Primer fixture: `v2_screens_catalog` (escuela+unidad+3 usuarios con
//     grants específicos a los 4 recursos meta del SDUI:
//     screen_templates, screen_instances, permissions_mgmt, roles).
//     Asume L4 completo: no siembra recursos/permisos/pantallas, solo
//     el envoltorio multi-tenant y la matriz de roles.
//   - 3.34.0: F1 (permisología MVP, ADR-6) — herencia de roles. Nueva
//     columna nullable `iam.roles.parent_role_id` (FK self-referencial a
//     iam.roles(id), ON DELETE SET NULL, índice idx_roles_parent) vía el
//     entity Role (AutoMigrate la crea). En el seed L4 los 5 alias
//     school_director/coordinator/assistant (→ school_admin) y
//     assistant_teacher/observer (→ teacher) dejan de declarar grants
//     propios y apuntan a su canónico; la herencia se resuelve y aplana
//     en el login (api-identity) sin tocar el formato del JWT ni el
//     matcher. readonly_auditor permanece standalone (no es superset
//     exacto de teacher). L4_SEED_VERSION → 1.17.3.
//   - 3.35.0: F2 (plan 004-permisologia-mvp) — poda del seed SDUI. Se
//     retiran 13 screen_instances y sus mappings en resource_screens
//     (guardian-relations-list/form, guardian_relations-form alias,
//     guardian-requests-list, schedules-list/form, calendar-list/form,
//     colors-list/form, stats-detail, progress-detail, report-card) más
//     el template L4 master-detail-basic-v1 (0 instancias). Se conservan
//     los dashboards progress-dashboard / stats-dashboard, el flujo de
//     examen completo, school-concepts y audit. Recursos huérfanos
//     resultantes (guardian_relations, schedules, calendar, colors,
//     reports) quedan como prune-later: NO se tocan iam.resources ni
//     iam.role_grants en esta pasada. L4_SEED_VERSION → 1.18.0.
//   - 3.36.0: F3 (plan 004-permisologia-mvp) — estándar de pantallas SDUI.
//     Las screen_instances migran al patrón delta (template define
//     default_actions; la instancia solo declara actions_added/removed).
//     Bumps L0→1.5.0, L2→1.3.0, L3→1.2.0, L4→1.19.0. Sin cambio semántico
//     (harness de round-trip verde). resource_screens intacta.
//   - 3.37.0: N0.0 (plan 005, capa de datos del onboarding) — 2 tablas
//     nuevas en el schema academic: `school_invitations` (códigos de
//     invitación con rol predefinido) y `school_join_requests`
//     (solicitudes de ingreso con doble gate de aprobación
//     school/unit + status pending/approved/rejected). Ambas entities
//     en AutoMigrate (SchoolInvitation antes de SchoolJoinRequest por la
//     FK invitation_id). post_gorm.sql agrega triggers set_updated_at,
//     las FKs cross-schema/cross-tabla (GORM no las materializa desde el
//     tag constraint sin campo de relación, mismo caso que
//     guardian_relations) y el índice UNIQUE parcial
//     idx_join_requests_pending_unique (una solicitud pendiente por
//     user/school/unit). Seeds L4: 3 recursos nuevos (invitations,
//     join_requests visibles bajo academic; join_request_approvals
//     API-only) + 9 permisos (invitations.{create,read,revoke},
//     join_requests.{create,read,reject}, join_request_approvals.{student,
//     teacher,guardian} — la acción ES el rol que se admite) + grants a
//     teacher (invitations.*, join_requests.*, approvals.student) y deny
//     a readonly_auditor (*.revoke, *.reject, approvals.*).
//     L4_SEED_VERSION → 1.20.0.
//   - 3.38.0: N0.4-A (plan 005, pantalla SDUI "gestionar invitaciones")
//     — 2 screen_instances nuevas en L4 sobre el recurso academic
//     `invitations`: `invitations-list` (list-basic-v1, scope school,
//     required_permission academic.invitations.read; patrón delta:
//     actions_removed [edit,delete] + actions_added [revoke] scope row
//     con permiso academic.invitations.revoke; create header heredado →
//     academic.invitations.create) e `invitations-form` (form-basic-v1,
//     create-only: actions_removed [save,delete]; fields academic_unit_id
//     remote_select, role select student/teacher/guardian, label,
//     expires_at, max_uses — code lo genera el backend). 2 filas en
//     resource_screens (invitations→list default, invitations→form).
//     Sin permisos/recursos nuevos (ya sembrados en N0.0).
//     L4_SEED_VERSION → 1.21.0.
//   - 3.39.0: N0.4-B (plan 005, bandeja de solicitudes pendientes) — 1
//     fila nueva en resource_screens: join_requests→list default
//     (screen_key `join-requests-inbox`). El FE la pinta con una pantalla
//     Compose NATIVA (no SDUI); NO se siembra screen_instance. Sin
//     permisos/recursos nuevos (ya sembrados en N0.0). L4_SEED_VERSION →
//     1.22.0.
//   - 3.40.0: N1.7 F0a etapa 1 (plan 010, ADR 0009) — capa de esquema
//     PURAMENTE ADITIVA de "sesiones de materia". 2 tablas nuevas en el
//     schema academic: `subject_offerings` (materia + seccion + periodo +
//     docente como unidad de enseñanza/inscripcion; unique compuesto
//     uq_subject_offerings_natural sobre school/subject/unit/section/period;
//     capacity reservado sin uso N1.7) y `subject_offering_enrollments`
//     (inscripcion del alumno a una sesion, PK compuesta offering+student).
//     Ambas entities en AutoMigrate (SubjectOffering antes de
//     SubjectOfferingEnrollment por la FK offering_id; ambas tras subjects/
//     academic_units/memberships/academic_periods). post_gorm.sql agrega las
//     7 FKs (GORM no las materializa sin campo de relacion, mismo caso que
//     subjects/school_invitations): docente=SET NULL, el resto CASCADE; mas
//     trigger set_updated_at en subject_offerings. En esta etapa la tabla
//     legacy `membership_subjects` aun coexistia con las nuevas; se elimina
//     en 3.41.0 (F0b).
//     Seed L4: recurso nuevo `subject_offerings` (b4000000-…-23, bajo
//     academic) + 5 permisos academic.subject_offerings.{create,read,update,
//     delete,enroll}. school_admin los cubre via wildcard academic.* (sin
//     enumerar); teacher gana literal academic.subject_offerings.read (paridad
//     con academic.subjects.read). Enum PermissionSubjectOfferings* agregado
//     en edugo-shared. L4_SEED_VERSION → 1.29.0.
//   - 3.41.0: N1.7 F0b (plan 010, ADR 0009) — eliminacion de la tabla legacy
//     `membership_subjects`. Se borra su entity y se quita de AutoMigrate, por
//     lo que la tabla ya no se materializa al recrear la BD (esquema
//     declarativo: sin DROP). El sentido "alumno-cursa-materia" / "docente-
//     dicta-materia" vive ahora en subject_offerings + subject_offering_
//     enrollments (3.40.0). Seeds demo y playground migrados a sesiones.
//     Catalogo L4 (Opcion A): se retira el feature "Mis Materias" del alumno
//     (instancia my-memberships-list + recurso/mapping/permiso/grant) y se
//     desembebe "alumnos por materia" del form de materia (subjects-form
//     vuelve a form-basic-v1, sin detail_config); se quita el campo
//     subject_ids del memberships-form. L4_SEED_VERSION → 1.30.0.
//   - 3.42.0: ADR 0016 (materia = catalogo de ESCUELA) — la entity
//     academic.subjects gana un unique compuesto uq_subjects_school_name sobre
//     (school_id, name). GORM lo materializa via tag uniqueIndex en ambas
//     columnas (mismo patron que uq_subject_offerings_natural; no requiere
//     post_gorm.sql). Respalda a nivel BD la validacion logica
//     ExistsByNameInSchool de la API academic e impide materias duplicadas por
//     nombre dentro de una escuela. Seeds reconciliados a materia=escuela
//     (academic_unit_id = NULL) en demo + playgrounds v2 n1_inscripcion/
//     n17_secciones/multi_unidad, deduplicando nombres repetidos por escuela.
//     L4_SEED_VERSION → 1.31.0 (scope sessions-by-subject-list school → unit).
//   - 3.43.0: el detalle de materia (subjects-form) queda SOLO con la pestaña
//     "Sesiones". Se retira la entrada "Alumnos" (students-by-subject-list) del
//     detail_configs y se ELIMINA esa screen_instance por completo (constructor,
//     registro en el slice y constante L4_SCREEN_INST_STUDENTS_BY_SUBJECT_ID):
//     era SOLO ese panel embebido, sin otra referencia (no estaba en menú ni en
//     resource_screens). Además `sessions-by-subject-form` corrige su scope
//     school → unit (form unidad-scoped, selector de docente requiere unidad
//     activa). Sin cambios de esquema/migraciones. L4_SEED_VERSION → 1.41.0.
//   - 3.44.0: se retira el camino de CREACIÓN DIRECTA de membresías (redundante
//     con invitación→solicitud→doble-gate→aprobación). Se ELIMINAN las
//     screen_instances `memberships-form` y `membership-add` (constructores,
//     registros y constantes …53/…d2) y sus 2 mappings en resource_screens;
//     `memberships-list` gana actions_removed:["create"] (conserva edit/delete/
//     expire). Leer/editar/expirar/borrar membresías intacto. Sin cambios de
//     esquema/migraciones. L4_SEED_VERSION → 1.42.0.
//   - 3.45.0: el período académico se ata además a la unidad. La entity
//     academic.academic_periods gana la columna `academic_unit_id` (uuid,
//     NOT NULL, index, FK a academic.academic_units(id) ON DELETE CASCADE,
//     espejo de school_id). El índice único parcial idx_academic_periods_active
//     pasa de (school_id) a (school_id, academic_unit_id) WHERE is_active=true,
//     por lo que la exclusividad del período activo es por unidad. Seeds que
//     insertan períodos (demo + playgrounds v2 n1_inscripcion/n17_secciones/
//     multi_unidad + fixture e2e screen_only) propagan academic_unit_id.
//   - 3.46.0: invariante "una oferta por materia por alumno" (bug 0036). La
//     entity academic.subject_offering_enrollments gana la columna `subject_id`
//     (uuid, NOT NULL, index; copia denormalizada e INMUTABLE del subject_id de
//     la oferta) con uniqueIndex compuesto uq_enrollment_student_subject
//     (student_membership_id, subject_id) que impide a un alumno inscribirse en
//     dos ofertas de la MISMA materia. post_gorm.sql agrega la FK
//     subject_offering_enrollments_subject_fkey → academic.subjects(id) ON
//     DELETE CASCADE (GORM no la materializa sin campo de relacion). Los seeds
//     que insertan enrollments (demo + playgrounds v2 n1_inscripcion/
//     n17_secciones + fixture integration academic_seed) propagan subject_id.
//   - 3.47.0: PRE 1a tenant→JWT de asistencia (cambio seed-only, sin DDL, igual
//     que 3.43.0/3.44.0). El form `attendance-batch` pierde el campo tenant
//     `unit_id` (la unidad sale del JWT vía RequireActiveContext, nunca del
//     form/query) y se ELIMINA el screen huérfano `attendance-form` (no mapeado
//     en resource_screens; solo lo respaldaba el contrato KMP, también
//     eliminado) — cierre del latente bug 0034. L4_SEED_VERSION → 1.42.6.
//   - 3.47.1: N2 feature de asistencia (plan 008, cambio seed-only, sin DDL,
//     igual que 3.47.0). (1) Las 3 instancias `attendance-*` corrigen
//     `api_prefix` de "learning" a "academic" (D5). (2) Entry-point "Pasar
//     lista": action `take-attendance` en `subjects-form` que navega a
//     `attendance-batch` con subjectId, gateada por `academic.attendance.create`
//     (D2). L4_SEED_VERSION → 1.42.7. Sin cambios de permisos.
//   - 3.47.2: N2.S2 cierre (plan 008 D5, cambio seed-only, sin DDL, igual que
//     3.47.1). El form `attendance-batch` (override nativo "pasar lista") declara
//     la action de submit `submit-batch` (scope header, permission
//     academic.attendance.create, event_id submit-batch) en su slot_data: es el
//     permiso del botón del override nativo (ADR 0003), espejo de la action
//     `enroll` de batch-enroll, y activa el gate cliente del botón (antes quedaba
//     null por falta de action de submit). L4_SEED_VERSION → 1.42.8. Sin cambios
//     de permisos.
//   - 3.47.3: N2.S3 (plan 008, cambio seed-only, sin DDL, igual que 3.47.2). El
//     form `subjects-form` suma dos entry-points de consulta de asistencia espejo
//     de "take-attendance": `view-attendance` ("Historial", event_id
//     view-attendance, order 21) y `view-attendance-summary` ("Resumen", event_id
//     view-attendance-summary, order 22), ambos scope resource-toolbar,
//     condition edit-only, permission academic.attendance.read. Navegan a las
//     pantallas SDUI genéricas attendance-list / attendance-summary pasando
//     subjectId; el destino del evento vive en SubjectsFormContract del KMP.
//     L4_SEED_VERSION → 1.42.9. Sin cambios de permisos.
//   - 3.48.0: plan 013 F1 — esquema de notas + invariante multi-período. (1)
//     academic.grades.teacher_id pasa a FK real → academic.memberships(id) ON
//     DELETE SET NULL: GORM no la materializa desde el tag
//     `constraint:grades_teacher_fkey` sin campo de relación, así que se declara
//     en post_gorm.sql (mismo patrón/política que subject_offerings_teacher_fkey;
//     teacher_id es nullable, el docente se desvincula sin borrar la nota). (2)
//     academic.subject_offering_enrollments gana la columna `period_id` (uuid,
//     NOT NULL, index; copia denormalizada e INMUTABLE del period_id de la
//     oferta, FK ya cubierta por CASCADE del propio offering_id). El uniqueIndex
//     uq_enrollment_student_subject pasa de 2 a 3 columnas
//     (student_membership_id, subject_id, period_id), por lo que el invariante
//     cambia de "una oferta por materia (ever)" a "una oferta por materia POR
//     PERÍODO" (D4): habilita Matemática-2025 + Matemática-2026, sigue
//     prohibiendo 2 secciones del mismo período (bug 0036). El guard de enroll
//     en academic (FindConflictingSubjectEnrollments) y el insert pasan a
//     considerar period_id. Seeds que insertan enrollments (demo + playgrounds v2
//     n1_inscripcion/n17_secciones) propagan period_id desde la oferta. EduGo no
//     está en producción → sin backfill.
//   - 3.48.0 (seed-only, sin DDL → SchemaVersion sin bump): N3.5 F1 (plan 014 /
//     ADR 0018). Reubicación de los entry-points de asistencia/notas de la materia
//     a la card de la sesión. Las 4 acciones del docente (take-attendance,
//     put-grades, view-attendance, view-attendance-summary) se BORRAN de
//     subjects-form (donde eran scope resource-toolbar y mezclaban el roster de un
//     docente con dos secciones de la misma materia) y se AÑADEN a
//     sessions-by-subject-list como row-actions (scope row, condition always): el
//     id de la fila es el offering_id, así cada acción opera sobre la sección
//     concreta. Mismos permisos (academic.attendance.create/read,
//     academic.grades.create; ya sembrados, cubiertos por el wildcard academic.*
//     de teacher). Además se reordenan las columnas de sessions-by-subject-list:
//     section_label pasa primero (headline que distingue A/B) y se quita
//     subject_name (redundante dentro del detalle de la materia). Reubicación, no
//     convivencia. L4_SEED_VERSION → 1.43.0.
//   - 3.49.0: N4 F1 (plan 015 / ADR 0019) — DEMOLICIÓN + RECONSTRUCCIÓN del
//     esquema de evaluación/contenido, anclado al modelo de sesión. EduGo no
//     está en producción → recrear BD sin backfill.
//     DEMOLIDO: el esquema viejo llaveado a auth.users + subject/grade texto-libre.
//   - entities borradas y reescritas: assessment, question, question_option,
//     assessment_material, assessment_assignment, assessment_attempt,
//     assessment_attempt_answer, attempt_review (schema assessment); material,
//     material_version, progress (schema content).
//   - post_gorm.sql: ELIMINADAS las tablas analíticas viejas
//     assessment.attempt_analytics y assessment.assessment_stats (llaveadas a
//     auth.users; analítica DIFERIDA en N4) y los índices de assignment por
//     student_id/academic_unit_id (modelo global muerto).
//     NUEVO (anclado a sesión):
//   - assessment.assessment: created_by_user_id → created_by_membership_id
//     (→academic.memberships RESTRICT), subject/grade texto → subject_id
//     (→academic.subjects RESTRICT), school_id NOT NULL (CASCADE), status
//     in (draft,published,archived), mongo_document_id reservado para V2.
//   - assessment.question / question_option: renombradas a singular; la opción
//     correcta vive en question.correct_answer (sin is_correct en la opción).
//   - assessment.assessment_material: N:N con PK compuesta (assessment_id,
//     material_id) → content.materials (arregla A4: lector deja de asumir 1:1).
//   - assessment.assessment_assignment: el PUENTE a la sesión. Se elimina
//     student_id XOR academic_unit_id + CHECK; target = subject_offering_id
//     (→academic.subject_offerings CASCADE) + UNIQUE (assessment_id,
//     subject_offering_id). Destinatarios se resuelven de
//     subject_offering_enrollments (arregla A2).
//   - assessment.assessment_attempt: student_id → student_membership_id
//     (→academic.memberships); UNIQUE parcial (assessment_id,
//     student_membership_id) WHERE status='in_progress' (un solo intento activo).
//   - assessment.attempt_review: reviewer_id → reviewer_membership_id.
//   - content.materials: subject/grade texto → subject_id (→academic.subjects
//     SET NULL, nullable), uploaded_by_teacher_id → uploaded_by_membership_id
//     (→academic.memberships RESTRICT).
//   - content.material_version: changed_by → changed_by_membership_id.
//   - content.progress: PK (material_id, user_id) → (material_id,
//     student_membership_id).
//     Todas las FKs cross-schema y el UNIQUE de assignment en post_gorm.sql
//     (GORM no las materializa sin campo de relación). content.courses queda
//     FUERA de alcance (intacto). Seeds de evaluación (demo + playground
//     focal_evaluacion*) y SDUI viejos de evaluación NO migrados aún: son F2/F4.
//   - 3.50.0: N4 F4.1 (plan 015 / ADR 0020) — esquema de notas con procedencia,
//     componentes, auditoría y perfil de escuela. EduGo no está en producción →
//     recrear BD sin backfill. (1) academic.grades gana la columna `source`
//     varchar(20) NOT NULL DEFAULT 'manual' CHECK IN ('auto_scored','manual',
//     'auto_llm') — procedencia de la nota unificada (CHECK inline en tag GORM,
//     mismo patrón que schools.subscription_tier). (2) NUEVA academic.grade_item
//     (componentes de nota): grain no-único (membership_id, subject_id, period_id)
//     vía idx_grade_item_grain; value/weight decimal(5,2) nullable (weight
//     informativo gen 1); source con el mismo CHECK; trazabilidad opcional al
//     origen auto vía source_attempt_id (FK→assessment.assessment_attempt SET NULL)
//   - source_assessment_id (FK→assessment.assessment SET NULL); created_by_
//     membership_id (FK→memberships RESTRICT); UNIQUE PARCIAL uq_grade_item_attempt
//     (membership_id, subject_id, period_id, source_attempt_id) WHERE
//     source_attempt_id IS NOT NULL (no duplicar el auto_scored por intento). (3)
//     NUEVA academic.grade_history (auditoría de override, append-only sin
//     updated_at): apunta a EXACTAMENTE UNO de grade_id (FK→grades CASCADE) /
//     grade_item_id (FK→grade_item CASCADE) vía CHECK XOR
//     grade_history_target_xor_check (((grade_id IS NOT NULL)::int + (grade_item_id
//     IS NOT NULL)::int) = 1); old_value/new_value decimal(5,2); changed_by_
//     membership_id (FK→memberships RESTRICT); changed_at default now(); reason
//     text. Índices idx_grade_history_grade / idx_grade_history_item. (4)
//     academic.schools gana la columna `grade_profile` varchar(20) NOT NULL
//     DEFAULT 'basic' CHECK IN ('basic','detailed') — perfil de notas básico/
//     detallado, gate por permisos en FE (CHECK inline en tag GORM, mismo patrón
//     que subscription_tier, misma tabla). Las FKs cross-schema (a assessment.*),
//     el CHECK XOR y el UNIQUE parcial viven en post_gorm.sql (GORM no los
//     materializa sin campo de relación). Sin tocar seeds (F4.6) ni APIs.
//   - 3.51.0 (seed-only, sin DDL): poda SDUI de material. L3 deja de
//     sembrar las 2 ScreenInstances `materials-list` / `material-form`
//     (+ slot_data) y el mapping resource_screen `material:form`. Eran
//     código muerto: las pantallas de material en la app son NATIVAS
//     (Compose) y no consumen esos seeds SDUI. El recurso materials sigue
//     en el menú vía el mapping `materials:list` (is_default, SIN
//     ScreenInstance — mismo patrón que material-detail / pantallas
//     nativas). L3_SEED_VERSION 1.2.0→1.3.0. Bump de SchemaVersion para
//     que el migrator recree el dataset (cambia el conteo de filas L3).
//   - 3.52.0: F2 (plan 018 / f2-spec) — rediseño de material a maestro-detalle.
//     EduGo no está en producción → recrear BD sin backfill.
//     MAESTRO content.materials: SE QUITAN las columnas inline de archivo
//     (file_url, file_type, file_size_bytes) que bajan al hijo; SE AGREGA
//     `summary` text nullable (markdown a mano del docente, DEC-2; distinto del
//     material_summary IA de MongoDB). Se conservan status (informativo, del
//     tema, DEC-4), description, processing_*, is_public, FKs y el índice parcial
//     idx_materials_status_active.
//     NUEVA content.material_file (DETALLE, N archivos por tema): id, material_id
//     (FK→content.materials CASCADE same-schema, la materializa GORM), file_url,
//     file_name (DEC-1, not null), file_type, file_size_bytes, created_at (DEC-3:
//     el orden sale de aquí). SIN status (DEC-4), SIN sort_order (DEC-3).
//     ELIMINADA content.material_version (entity material_version.go, su registro
//     en gorm_migrator y sus FKs material_version_{material,membership}_fkey en
//     post_gorm.sql): versionaba el único archivo inline, queda huérfana con N
//     archivos distintos (Hallazgo 1 — "no deprecar: eliminar"). El truncate del
//     demo seed sustituye content.material_version por content.material_file.
//     assessment.assessment_material intacto (el examen sigue apuntando al tema).
//     Sin tocar permisos: materials.delete sigue solo en L4 (ver nota infra).
//     FIX seed L3 (mismo bump, la BD nunca se aplicó a 3.52.0): la poda SDUI de
//     3.51.0 eliminó AMBAS screen_instances L3 (materials-list, material-form),
//     pero `materials-list` tiene mapping resource_screen (menú) y la FK
//     fk_resource_screens_screen_key exige su screen_instance → un recreate
//     limpio fallaba en L3 con violación 23503. Se RESTAURA la screen_instance
//     MÍNIMA `materials-list` (no renderizada; pantalla NATIVA Compose), patrón
//     batch-enroll/join-requests-inbox de L4. `material-form` SIGUE PODADO (sin
//     mapping → sin FK). L3_SEED_VERSION 1.3.0→1.4.0; +1 fila screen_instances.
//     Test l3_apply_twice y fixture e2e l3_constants_export ajustados (materials-
//     list: aserción negativa→positiva; material-form: sigue negativa).
//   - 3.53.0: Fase 2 — nuevo tipo de pregunta `multiple_select` (opción
//     múltiple con varias respuestas correctas, solo autoría). El CHECK
//     inline del entity Question (question_type_check) suma 'multiple_select'
//     a la lista permitida (de 4 a 5 tipos), igual que su tag `validate`
//     oneof. Cambia el output del migrador GORM → bump obligatorio. Contrato
//     de datos: para este tipo, assessment.question.correct_answer guarda un
//     ARRAY JSON de textos (["Texto A","Texto C"]); NO se añade is_correct a
//     question_option (los demás tipos no cambian). Acompaña el seed L4 del
//     form de pregunta (nuevo slot `options_multi` con selection_mode
//     multiple, visible_when question_type in [multiple_select]).
//     L4_SEED_VERSION → 1.49.0.
//   - 3.54.0: seed-only (sin DDL). assessment-questions-list elimina la
//     row-action SDUI heredada `edit` (default de list-basic-v1) vía
//     "actions_removed": ["edit"]: en el detalle de preguntas la edición la
//     cubre el botón nativo "Editar" del bottom-sheet; la acción SDUI no tenía
//     handler. L4_SEED_VERSION → 1.50.0.
//   - 3.55.0: seed-only (sin DDL). Dos ajustes de evaluación: (1) la action
//     "Publicar" de assessments-form alinea su slot.permission a
//     content.assessments.publish (antes .update) para igualar el gate del botón
//     con la ruta POST /assessments/:id/publish; (2) se ELIMINA la pantalla SDUI
//     assessment-assignment (reemplazada por modal nativo), conservando el
//     recurso assessments y el permiso content.assessments.assign. Cambia el
//     slot_data del seed L4 → bump obligatorio para invalidar la caché SDUI por
//     contenido. L4_SEED_VERSION → 1.51.0.
//   - 3.56.0: seed-only (sin DDL). plan 017 F2: assessments-form migra el campo
//     subject_id de remote_select a entity-picker (modal con búsqueda server-side
//   - paginación contra academic:/api/v1/subjects). Cambia el slot_data del
//     seed L4 → bump obligatorio para invalidar la caché SDUI por contenido.
//     L4_SEED_VERSION → 1.52.0.
//   - 3.57.0: seed-only (sin DDL). ADR 0022: assessments-form declara view_when a
//     nivel slot_data → el front pone el form read-only total cuando la evaluación
//     no es borrador. Acompaña backend learning (subject_id editable solo en
//     borrador; update fuera de borrador → 400 BUSINESS_ASSESSMENT_NOT_DRAFT).
//     Cambia el slot_data del seed L4 → bump para invalidar la caché SDUI por
//     contenido. L4_SEED_VERSION → 1.53.0.
//   - 3.58.0: seed-only (sin DDL). Poda de dos pantallas SDUI legacy huérfanas:
//     (1) grades-form (reemplazada por nativas my-grade-detail/grades-batch) y
//     (2) user-roles (huérfana, sin reemplazo ni entry-point). Ambas tenían
//     controles remote_select MUERTOS (student_id/subject_id, user_id). Se
//     eliminan sus screen_instances + mappings en resource_screens + constantes.
//     Cambia el set de screens del seed L4 → bump para invalidar la caché SDUI
//     por contenido. Sin cambios de roles ni permisos. L4_SEED_VERSION → 1.54.0.
//   - 3.59.0: plan 020 N5 F1.1/F1.2 — push M2M + device tokens. NUEVAS
//     notifications.device_tokens (entity DeviceToken, FK auth.users CASCADE,
//     UNIQUE user_id+device_token, índice parcial idx_device_tokens_user_active)
//     y auth.service_clients (entity ServiceClient, scopes text[], índice parcial
//     idx_service_clients_active). Seed L5-m2m: edugo-worker y edugo-api-learning
//     con scope notifications.dispatch y secret_hash bcrypt del dev secret de
//     push-secrets.env. L5_SEED_VERSION → 1.0.0.
//   - 3.60.0: plan 020 F4.6.8 — persistir tenant en la notificación in-app. La
//     entity Notification suma school_id/unit_id (uuid nullable) a
//     notifications.notifications para que la lista in-app pueda resolver el
//     context-switch multi-tenant al tocar (antes solo viajaban en el push, V1).
//     Solo DDL aditivo vía AutoMigrate (2 columnas nullable, sin índice); sin
//     cambios de seeds ni SQL post_gorm. Bump por la regla 1 (cambio en entity).
//   - 3.60.1: eliminación de la feature muerta `content.courses` — ningún código
//     vivo la lee (la API learning ya fue limpiada). Se borra la entity Course,
//     su registro en el AutoMigrate de gorm_migrator.go y el seedCourses + su
//     truncate en demo/development.go. AutoMigrate nunca dropea, así que un
//     recreate fresco simplemente ya no crea la tabla (sin DROP, sin SQL
//     post_gorm). Bump por la regla 1 (cambio en migrations/ + seeds/), aunque
//     ComputeFilesHash() no cambia (solo hashea pre/post_gorm.sql). En este
//     paso content.progress / entities.Progress aun seguian vivas; se
//     eliminan despues en 3.60.3.
//   - 3.60.2: seed-only (sin DDL). Entry-point "Gestionar Conceptos" en el
//     form `schools-form`: nueva action de navegación `manage-concepts`
//     (scope form-submit, condition edit-only, permission
//     admin.concept_types.read, event_id manage-concepts) que abre la
//     pantalla ya sembrada `school-concepts-list` (el wiring KMP en
//     SchoolsFormContract ya existía). Cambia el slot_data del seed L4 →
//     bump para invalidar la caché SDUI por contenido. L4_SEED_VERSION
//     1.54.0 → 1.55.0. Sin cambios de esquema ni de permisos.
//   - 3.60.3 (seed/DDL): elimina la tabla content.progress huerfana —
//     productor y lector removidos en paralelo (MP-04). Se borra la entity
//     Progress, su registro en el AutoMigrate de gorm_migrator.go, las 2 FKs
//     (progress_material_fkey / progress_student_fkey) y el trigger
//     set_updated_at de content.progress en post_gorm.sql, y su truncate en
//     demo/development.go. El schema content NO se borra (content.materials
//     sigue viva). AutoMigrate nunca dropea: un recreate fresco simplemente ya
//     no crea la tabla (sin DROP). Bump por la regla 1; ComputeFilesHash()
//     CAMBIA esta vez (se editó post_gorm.sql).
//   - 3.60.4: seed-only (sin DDL, igual que 3.60.2). Arregla la pantalla
//     `audit-detail` (detalle de evento de auditoría) que pintaba campos de
//     material/archivo ("Tamaño/Subido/Estado/Descripción" + botón
//     "Descargar"): el renderer de detalle del KMP está atado a las `zones`
//     del template y el slot_data del instance no puede cambiar los `field`
//     ni los slots, solo los labels. Se mina un template propio L4
//     `audit-detail-v1` (a4000000-...006, pattern detail) con los campos
//     REALES del evento (actor_email/role/ip/user_agent, service_name,
//     action, resource_type/id, request_method/path, status_code, severity,
//     category, created_at) en solo lectura, labels en español, ícono "list"
//     y sin descarga; `auditDetail()` se reapunta a él (antes
//     detail-basic-v1 de L0). Endpoint (identity:/api/v1/audit/events/:id) y
//     permiso (admin.audit.read) intactos. Cambia slot_data + se agrega un
//     template → bump para invalidar la caché SDUI por contenido.
//     L4_SEED_VERSION 1.55.0 → 1.56.0. Sin cambios de esquema ni de permisos.
//   - 3.60.5: seed-only (sin DDL, igual que 3.60.4). Fix de render de
//     `audit-detail-v1`: las filas de detalle usaban controlType "list-item"
//     (DSListRow), que pinta el valor como headline + un chevron de navegación
//     y deja el label vacío (DSListRow solo toma el atributo estático `label`
//     como supporting, ignora bind/default). Ahora cada campo es una sub-zona
//     container con DOS slots controlType "label" (uno de etiqueta, texto
//     español estático en `value`, style "caption"; otro de valor con `field`,
//     style "body"), espejo de cómo detail-basic-v1 pinta sus filas de valor:
//     sin chevron y se ve "Etiqueta / valor" en solo lectura. Cambia el
//     definition del template L4 → bump para invalidar la caché SDUI por
//     contenido. L4_SEED_VERSION 1.56.0 → 1.57.0. Sin cambios de esquema ni de
//     permisos.
//   - 3.60.6: seed-only (sin DDL, igual que 3.60.5). MP-03 F3: rangos numéricos
//     declarativos en los forms SDUI. Cada campo `"type": "number"` ahora lleva
//     `min`/`max` en su slot_data para que el FE KMP valide antes de enviar,
//     espejando el binding real del backend donde existe (assessments-form:
//     pass_threshold 0–100, max_attempts/time_limit_minutes min=1;
//     assessment-question-form: points min=0) y con un mínimo conservador donde
//     el backend no declara rango (period-form academic_year 1900–2100,
//     invitations-form max_uses min=1). Cambia slot_data de instances L4 → bump
//     para invalidar la caché SDUI por contenido. L4_SEED_VERSION 1.57.0 →
//     1.58.0. Sin cambios de esquema ni de permisos.
//   - 3.60.7 — MP-01 F3: poda de tablas muertas academic.{schedule,calendar_event,colors} y playground focal_colors_demo
//   - 3.60.8 — ADR 0024 F0: landing_screen_key en roles + default_landing_screen_key en schools
//   - 3.61.0 — MP-08 F0 (aditivo, solo esquema): 4 tablas nuevas modelando en
//     datos el acceso por sistema y la equivalencia tipo-de-invitacion->rol
//     (todo por FK de id, nunca por nombre). iam.systems (catalogo de apps),
//     iam.system_roles (puente sistema<->rol), academic.invitation_types
//     (catalogo global de tipos de invitacion) y academic.school_invitation_roles
//     (equivalencia (escuela, tipo) -> rol IAM, FK cross-schema academic->iam).
//     Las 4 entities entran en AutoMigrate; post_gorm.sql agrega sus FKs (GORM
//     no las materializa sin campo de relacion) y los triggers set_updated_at.
//     Puramente ADITIVO: no toca tablas existentes. SIN seeds (los valores de
//     los catalogos los siembra F1); L4_SEED_VERSION intacto.
//   - 3.62.0 — MP-08 F3 (swap de columna, NO aditivo): la columna role (varchar
//   - CHECK inline del enum de roles) se reemplaza por invitation_type_id (uuid,
//     FK -> academic.invitation_types(id)) en las 3 tablas academic.{memberships,
//     school_invitations,school_join_requests}. El CHECK ..._role_check se elimina
//     (la validez del tipo la garantiza la FK). Las 3 FKs nuevas y el reexpresado
//     del indice parcial idx_memberships_unit_role_active ->
//     idx_memberships_unit_invitation_type_active viven en post_gorm.sql. Los
//     seeds resuelven la key del tipo a su id via catalog.ResolveInvitationTypeID
//     (data-driven, sin hardcodear UUIDs); L1 adelanta ApplyInvitationTypes para
//     que su membresia pueda resolver el FK antes de L4. L4_SEED_VERSION intacto
//     (no cambian filas de catalogo). Requiere recrear BD (sin ALTER).
//   - 3.63.0 — MP-08 F4 (seed-only, sin DDL): dos ajustes de slot_data SDUI en
//     L4. (1) P5: el form `invitations-form` cambia el campo `role` (select
//     estatico, enum legacy muerto) por `invitation_type` (remote_select contra
//     GET /api/v1/schools/invitation-types; value_field=key, display_field=label),
//     alineado a CreateInvitationRequest.InvitationType. (2) P4 (DEC-D):
//     `schools-list` retira la accion `create` del header (actions_removed
//     ["create"]); el alta de escuelas pasa al admin-tool de Go. Se conserva
//     schools-form + manage-concepts y la edicion de escuelas existentes.
//     L4_SEED_VERSION 1.60.0 -> 1.61.0. Bump de SchemaVersion para invalidar la
//     cache SDUI por contenido (recrear BD, sin ALTER).
//   - 3.64.0 — aprobacion de ingreso: SELLO x TIPO (seed-only, sin DDL). El
//     permiso unico academic.join_request_approvals.<tipo> se separa en
//     academic.join_request_approvals.{school,unit}.{student,teacher,guardian}
//     (catalogo: 3 filas -> 6). El doble gate (colegio->unidad) gana un permiso
//     por sello; approve.go (academic) evalua el permiso del sello concreto que
//     firma. El grant de teacher pasa de ...student a ...unit.student (admite
//     alumnos a SU clase = sello de unidad; ya no firma el sello de colegio).
//     school_admin (academic.*) / super_admin (*) cubren ambos sub-namespaces
//     por subarbol; readonly_auditor sigue denegado por su deny de prefijo
//     academic.join_request_approvals.*. L4_SEED_VERSION 1.61.0 -> 1.62.0. Bump
//     de SchemaVersion por cambio de catalogo+grant (recrear BD, sin ALTER).
//   - 3.65.0 — ADR 0024 DEC-4: elimina la columna decorativa scope_pattern de
//     iam.user_grants (el motor de auth nunca la evaluaba; el scope efectivo
//     vive en el JWT, no en el grant). Cambios: entity UserGrant pierde el campo
//     ScopePattern y el indice unico uq_user_grants_user_scope_perm_effect se
//     reescribe a uq_user_grants_user_perm_effect sobre (user_id,
//     permission_pattern, effect); post_gorm.sql elimina el CHECK
//     user_grants_scope_format; el demo seed (seedUserGrants) deja de sembrar
//     ScopePattern en sus 2 filas. NO toca iam.role_grants (ya limpio) ni
//     iam.user_roles.scope_pattern (sigue en uso). Cambio en entity + SQL +
//     demo seed (no L4) -> L4_SEED_VERSION intacto. Requiere recrear BD (sin
//     ALTER). ComputeFilesHash() CAMBIA (se editó post_gorm.sql).
//   - 3.66.0 — plan 022 / ADR 0024 foco 3: poda del recurso grades_detail
//     (seed-only, sin DDL). Se eliminan del catálogo L4 el recurso
//     `grades_detail` (…37) y sus 4 permisos academic.grades_detail.{create,
//     read,update,delete}. El modo detallado de notas ya no se gobierna con un
//     permiso: academic lo decide leyendo `grade_profile` de la escuela (el
//     permiso era un mensajero eliminable). Se retira también el grant condicional
//     por perfil que vivía en identity. L4_SEED_VERSION 1.62.0 -> 1.63.0. Bump de
//     SchemaVersion por cambio de catálogo de recursos+permisos (recrear BD, sin
//     ALTER).
//   - 3.67.0 — MP-09 F2-A: eliminación del paquete seeds/demo. El dataset de
//     desarrollo ya lo provee seeds/playground_v2/base (default del migrador
//     desde F1). Se borra el paquete seeds/demo (development.go + su test de
//     integración), se repuntan los consumidores no-test (cmd/runner, cmd/seed,
//     tools/mock-generator) a base.Apply, se elimina el branch SeedDemo de
//     migrate.go junto con el campo MigrateOptions.SeedDemo, y se retira el flag
//     `--seed-demo` del migrador. ComputeFilesHash() de seeds/ deja de incluir
//     la línea "demo:<SeedVersion>" → cambia el hash de seeds → bump obligatorio.
//   - 3.68.0 — MP-09 F4: las capas system L1..L4 quedan como CONTRATO PURO, sin
//     DATO DE TENANT. (1) L1 deja de sembrar escuela demo, usuario viewer,
//     user_role y membership; sólo conserva el rol de contrato
//     announcement_viewer (se borran l1_{school,user,user_role,membership}.go y
//     las constantes de tenant; L1_SEED_VERSION 1.2.0 -> 1.3.0). (2) L4 deja de
//     sembrar las equivalencias tipo->rol de la escuela demo (se elimina
//     ApplyDemoSchoolInvitationRoles y su paso 10; el helper genérico
//     SeedDefaultSchoolInvitationRoles se conserva; L4_SEED_VERSION 1.64.0 ->
//     1.65.0). (3) playground_v2/base ahora siembra school_invitation_roles para
//     sus 2 escuelas vía l4.SeedDefaultSchoolInvitationRoles (antes ninguna las
//     tenía). (4) tests/fixtures/scenarios L1..L3 actualizados a la nueva
//     realidad (sin viewer; scenario l1_readonly reducido a rol, nombre
//     conservado). Bump de seeds (L1/L4 SeedVersion) -> cambia el hash de seeds
//     -> bump obligatorio de SchemaVersion.
//   - 3.69.0: F1 plan-024 (representante): guardian_relations.school_id (NOT NULL,
//     índice único +school_id) + academic_unit_id; school_guardian_policy (política
//     por escuela); school_invitations.student_id (FK auth.users SET NULL). Recrear
//     BD, sin ALTER.
//   - 3.70.0: F4·S3·M0 plan-024 (representante): academic.memberships gana ESTADO
//     EXPLÍCITO `status` varchar(12) NOT NULL DEFAULT 'active' CHECK IN
//     ('pending','active','withdrawn') como ÚNICA fuente de verdad del estado; se
//     ELIMINA la columna `is_active` (era derivable: is_active=true ⟺
//     status='active'). `withdrawn_at` se conserva como timestamp informativo. El
//     índice parcial idx_memberships_unit_invitation_type_active pasa de
//     `WHERE is_active = true` a `WHERE status = 'active'`. CHECK inline en el tag
//     GORM del entity (mismo patrón que assessment.status / schools.grade_profile);
//     post_gorm.sql cambia el WHERE del índice → ComputeFilesHash() CAMBIA. Sin
//     cambio de comportamiento (default active equivale al is_active=true de hoy).
//     Seeds playground_v2 (common helper + base) migrados a status='active'; no
//     son parte del hash (MP-09). Recrear BD, sin ALTER. academic/identity migran
//     su lectura del estado después (otra tarea). L*_SEED_VERSION intacto (no
//     cambia ningún dato de las capas system L0–L4).
//   - 3.71.0 (2026-06-15): se ELIMINA el recurso/pantalla `progress`
//     (progress-dashboard) del seed L4 — su screen SDUI apuntaba a
//     /api/v1/stats/student (inexistente → 404) y era redundante con el
//     dashboard nativo del alumno. Cambia el catálogo de recursos/permisos
//     (resource `progress`, permisos `reports.progress.*` + grants, la
//     screen_instance/mapping `progress-dashboard`). L4_SEED_VERSION
//     1.66.0 → 1.67.0 → cambia el hash de seeds → bump obligatorio de
//     SchemaVersion. Recrear BD, sin ALTER. `stats`/`reports` intactos.
//   - 3.72.0 (2026-06-15): M4 plan-024 (representante) — higiene del seed L4: se
//     quita el campo `api_prefix:"learning"` INERTE del slot_data del
//     screen_instance `dashboard-guardian` (el dashboard del representante es
//     NATIVO y ya no carga por el pipe SDUI; nadie consume ese campo).
//     L4_SEED_VERSION 1.67.0 → 1.68.0 → cambia el hash de seeds → bump
//     obligatorio de SchemaVersion. Recrear BD, sin ALTER.
const SchemaVersion = "3.72.0"

// ComputeFilesHash calcula un SHA256 de los archivos SQL embebidos
// en el paquete migrations (pre_gorm.sql y post_gorm.sql).
func ComputeFilesHash() string {
	h := sha256.New()
	for _, name := range []string{"sql/pre_gorm.sql", "sql/post_gorm.sql"} {
		content, err := sqlFiles.ReadFile(name)
		if err != nil {
			continue
		}
		h.Write([]byte(name))
		h.Write(content)
	}
	return fmt.Sprintf("%x", h.Sum(nil))[:16]
}
