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
//   - 3.73.0: F6 plan-024 (representante) — tipo de evaluacion practica/final. (1)
//     assessment.assessment gana la columna `kind` varchar(20) NOT NULL DEFAULT
//     'final' CHECK IN ('practice','final') (CHECK inline en el tag GORM, mismo
//     patron que status / source_type); toda evaluacion existente queda 'final'.
//     (2) NUEVA academic.practice_result: ESPEJO de academic.grade_item para
//     evaluaciones 'practice' (resultado FUERA del expediente, solo estadisticas).
//     Columnas id (PK uuid), membership_id/subject_id/period_id (FK academic CASCADE,
//     grain no-unico via idx_practice_result_grain), title, value decimal(5,2)
//     nullable (% 0–100), source varchar(20) con el mismo CHECK que grade_item
//     ('auto_scored','manual','auto_llm'), source_attempt_id (FK→assessment.
//     assessment_attempt SET NULL), source_assessment_id (FK→assessment.assessment
//     SET NULL), created_by_membership_id (FK→memberships RESTRICT), created_at/
//     updated_at. UNIQUE PARCIAL uq_practice_result_attempt (membership_id,
//     subject_id, period_id, source_attempt_id) WHERE source_attempt_id IS NOT NULL
//     (espejo de uq_grade_item_attempt; defensa del upsert por id determinista del
//     worker). Las FKs cross-schema, el trigger set_updated_at y el UNIQUE parcial
//     viven en post_gorm.sql (GORM no los materializa sin campo de relacion). El
//     worker ramifica por `kind` ('final'→grade_item, 'practice'→practice_result).
//     Cambio en entity (kind) + nueva entity + post_gorm.sql → ComputeFilesHash()
//     CAMBIA. Recrear BD, sin ALTER. L*_SEED_VERSION intacto (solo DDL, sin datos).
//   - 3.75.0 (2026-06-17): ADR 0024 sub-deuda "herencia del landing" — los 6 roles
//     alias (school_director/coordinator/assistant → dashboard-schooladmin;
//     assistant_teacher/observer/readonly_auditor → dashboard-teacher) reciben
//     landing_screen_key explícito (antes NULL → caían a school.default
//     "dashboard-home", shell sin contrato resoluble). Solo datos de seed L4
//     (L4_SEED_VERSION 1.69.0 → 1.70.0); sin DDL. Recrear BD para reseeding.
//   - 3.76.0 (2026-06-17): mismo frente — el rol de contrato L1 announcement_viewer
//     (scope school) recibe landing_screen_key=dashboard-schooladmin (antes NULL).
//     Cierra el 7º rol secundario que quedaba sin landing. Solo dato de seed L1
//     (L1_SEED_VERSION 1.3.0 → 1.4.0); sin DDL. Recrear BD para reseeding.
//   - 3.77.0 (2026-06-17): dashboard-home deja de ser "shell muerta" y pasa
//     a ser el dashboard básico por defecto (home genérico para roles sin
//     landing_screen_key propio; school.default_landing_screen_key sigue
//     apuntando aquí). El FE le dará un render real self-contained (otro
//     frente), así que su slot_data L4 deja de declarar el api_prefix
//     "learning" (inerte; el dashboard no consume endpoint) y queda solo
//     {"title":"Inicio"}; el description del instance se actualiza a la
//     nueva semántica. Solo dato de seed L4 (L4_SEED_VERSION 1.70.0 →
//     1.71.0) → cambia el hash de seeds → bump obligatorio. Sin DDL ni
//     cambios de permisos. Recrear BD para reseeding.
//   - 3.78.0 (2026-06-17): saneo de over-grants + Escuelas read-only
//     (bugs 0064/0065/0054). (1) bug 0064: se quita `admin.users.*` de
//     teacher y guardian (over-grant del panel Usuarios). (2) bug 0065: se
//     quita `admin.system_settings.*` del alumno y los switches
//     push/email del template settings-basic-v1 ganan
//     `permission:"admin.system_settings.update"` (dark_mode/theme siguen
//     sin permission; el alumno conserva `notifications.*`). (3) bug 0054:
//     schools-list pasa a actions_removed ["create","edit","delete"]
//     (read-only; gestión en admin-tool de Go). Solo dato de seed L4
//     (L4_SEED_VERSION 1.71.0 → 1.72.0) → cambia el hash de seeds → bump
//     obligatorio. Sin DDL. Recrear BD para reseeding.
//   - 3.79.0 (2026-06-17): bug 0048 — se quita el over-grant
//     `content.assessments.*` de studentPatterns (el alumno NO usa ningún
//     `content.assessments.*`: su flujo completo —ver asignadas, tomar, ver
//     resultados— corre sobre `content.assessments_student.*`, que se
//     CONSERVA). El wildcard docente le otorgaba publish/delete/update/assign/
//     create/grade/review → veía los botones de gestión en assessments-form.
//     Solo dato de seed L4 (L4_SEED_VERSION 1.72.0 → 1.73.0) → cambia el hash
//     de seeds → bump obligatorio. Sin DDL. Recrear BD para reseeding.
//   - 3.80.0 (2026-06-18): lote triage alpha Grupo 2 (seed SDUI mal
//     configurado). 0055 auditoría: columnas remapeadas a los campos reales
//     del DTO (action/actor_email/resource_type) + chip único "Solo críticos"
//     (filter_all "Todos" + filter_processing severity=critical; se retira el
//     chip de info). 0062 grades-list pasa a read-only (actions_removed
//     ["create","edit","delete"]) porque create/edit navegaban a grades-form
//     ELIMINADA → 404. 0057 form de anuncios expone toggle is_pinned (L2).
//     Solo datos de seed (L2_SEED_VERSION 1.3.0 → 1.4.0, L4_SEED_VERSION
//     1.73.0 → 1.74.0) → cambia el hash → bump obligatorio. Sin DDL. Recrear
//     BD para reseeding.
//   - 3.81.0 (2026-06-19): bug 0069 — flag `is_system` en iam.roles e
//     iam.permissions (Opción A). Nueva columna `is_system bool NOT NULL
//     DEFAULT false` (índices parciales idx_roles_system / idx_permissions_system,
//     vía tag GORM en las entities; AutoMigrate la crea). Los seeds L0–L4
//     marcan `IsSystem: true` en TODO rol/permiso del contrato (apply +
//     accessors espejo, para que el cross-checker coincida). Habilita el guard
//     de runtime en edugo-api-identity (otra tarea) que rechaza delete/update
//     de entradas del contrato. Cambio en entities + datos de seed L0/L1/L3/L4
//     (L0_SEED_VERSION 1.5.0 → 1.5.1, L1_SEED_VERSION 1.4.0 → 1.4.1,
//     L3_SEED_VERSION 1.4.0 → 1.4.1, L4_SEED_VERSION 1.74.0 → 1.74.1; L2 sin
//     cambio). Recrear BD, sin ALTER.
//   - 3.82.0 (2026-06-19): over-grant de permisos — school_admin pasa de
//     `context.*` a `context.browse_units` (`context.browse_schools`, scope
//     system, queda solo en super_admin). Corrige el "Cambiar escuela" que se
//     encendía y fallaba con 403 para un coordinador/admin de una sola escuela.
//     Solo datos de seed (L4_SEED_VERSION 1.74.1 → 1.74.2) → cambia el hash →
//     bump obligatorio. Sin DDL. Recrear BD para reseeding.
//   - 3.83.0 (2026-06-19): UI de invitaciones — nueva screen_instance
//     `invitations-detail` (detalle de solo lectura, form-basic-v1 en modo
//     lectura) + row-action `copy-code` (copiar código) en `invitations-list`.
//     Backend gana GET /schools/invitations/{id} (otra tarea, academic). Solo
//     datos de seed (L4_SEED_VERSION 1.74.2 → 1.74.3) → cambia el hash → bump
//     obligatorio. Sin DDL. Recrear BD para reseeding.
//   - 3.84.0 (2026-06-19): detalle de invitación gana campo `Unidad`
//     (academic_unit_name) — academic lo resuelve por JOIN a academic_units.
//     Solo datos de seed (L4_SEED_VERSION 1.74.3 → 1.74.4). Sin DDL.
//   - 3.85.0 (2026-06-20): plan 025 F1.1 — system `messaging` (WhatsApp). Seed
//     L4: nueva fila iam.systems (key `messaging`, "EduGo Mensajería") + 8 filas
//     iam.system_roles (super_admin + árbol teacher {teacher, assistant_teacher,
//     observer} + árbol school_admin {school_admin, school_director,
//     school_coordinator, school_assistant}); nuevo recurso API-only `messaging`
//     (b4000000-…-d0, IsMenuVisible=false, sin pantalla SDUI) + 3 permisos
//     messaging.{session.pair, message.send, device.link}; grant WILDCARD
//     `messaging.*` (allow) a teacher y school_admin (los aliases lo heredan vía
//     ADR-6; super_admin ya cubre por `*`). student/guardian/readonly_auditor/
//     announcement_viewer NO lo reciben (las familias son destinatarias, no
//     emisoras). La API edugo-api-messaging autoriza por los grants del JWT
//     (no consulta IAM); el system_roles existe para que la web pública/admin
//     reconozcan el system. Solo datos de seed L4 (L4_SEED_VERSION 1.74.4 →
//     1.75.0) → cambia el hash de seeds → bump obligatorio. Sin DDL. Recrear
//     BD para reseeding.
//   - 3.86.0 (2026-06-20): plan 025 F5 — `messaging` NAVEGABLE en el menú.
//     El recurso messaging pasa a IsMenuVisible=true (item de menú raíz, scope
//     system, accesible sin escuela activa); nuevo permiso `messaging.view`
//     (action `view`, slot.permission de la pantalla + gate del item — cubierto
//     por el wildcard `messaging.*` de school_admin/teacher, no se enumera por
//     rol); nueva screen_instance `messaging` (pantalla NATIVA Compose, slot_data
//     list-basic-v1 mínimo por higiene, sin api_prefix) + resource_screen default
//     (list) que liga el recurso al screen_key `messaging`. El FE mapeará el
//     screen_key `messaging` a Route.Messaging (pantalla nativa F5). Solo datos
//     de seed L4 (L4_SEED_VERSION 1.75.0 → 1.76.0) → cambia el hash de seeds →
//     bump obligatorio. Sin DDL. Recrear BD para reseeding.
//   - 3.87.0 (2026-06-21): `units-form` migra el campo `parent_unit_id` de
//     `remote_select` (combo) a `entity-picker` (lupa) — estándar del proyecto.
//     Búsqueda server-side + paginación contra academic:/api/v1/units. Solo datos
//     de seed L4 (L4_SEED_VERSION 1.76.0 → 1.77.0) → cambia el hash de seeds →
//     bump obligatorio. Sin DDL. Recrear BD para reseeding.
//   - 3.88.0 (2026-06-21): Plan 026 (overflow de navegación) — priority/pin
//     ADITIVOS al contrato del menú. La entity iam.resources gana dos columnas
//     nullable/default (priority *int → NULL, pin bool → false), materializadas
//     por GORM AutoMigrate igual que sort_order/is_menu_visible (índices parciales
//     idx_resources_priority / idx_resources_pin vía tag GORM). El seed L4
//     (l4ResourceRow + UPSERT) propaga ambas; los 31 recursos quedan en MODO
//     LEGACY (priority NULL, pin false) → cero regresión (el front cae a
//     sort_order). Cambio en entity (DDL aditiva) + datos de seed L4
//     (L4_SEED_VERSION 1.77.0 → 1.78.0) → cambia el hash de seeds → bump
//     obligatorio. Recrear BD, sin ALTER.
//   - 3.89.0 (2026-06-24): Plan 027 (permisología por proceso). Seed-only, sin
//     DDL. ADITIVO: 2 recursos my_* (academic.my_teaching → profesor ve solo lo
//     que dicta; academic.my_attendance → alumno lee solo su asistencia) con
//     permiso read:own, nodo de menú, screen-instance readonly (my-teaching-list
//     / my-attendance-list) y resource_screens. SUSTRACTIVO (cierre de fugas de
//     ESCRITURA): poda de grants en studentPatterns / guardianPatterns /
//     teacherPatterns — el alumno y el representante dejan de poder registrar
//     asistencia (POST /attendance/batch) y crear/editar/publicar anuncios,
//     materiales y evaluaciones; el profesor pierde capabilities de
//     administración (units/periods/invitations/join_requests/admisión) que
//     asume school_admin. reports.* del teacher se mantiene como deuda (acotar a
//     stats.unit rompería stats-dashboard, que apunta a /stats/global). Solo
//     datos de seed L4 (L4_SEED_VERSION 1.78.0 → 1.79.0) → cambia el hash de
//     seeds → bump obligatorio. Recrear BD, sin ALTER.
//   - 3.90.0 (plan 027 F4.8): deny en school_admin para sanear la amplitud del
//     arquetipo "Administra" — `academic.*.read:own` (quita ruido de menú my_* del
//     admin) y `admin.roles.{create,update,delete}` (no define roles IAM del
//     sistema). Solo datos de seed L4 (L4_SEED_VERSION 1.79.0 → 1.80.0). Recrear BD.
//   - 3.91.0 (bug 0074): el campo `subject_id` del `assessments-form` se repunta de
//     `academic:/api/v1/subjects` a `academic:/api/v1/me/subjects` (el endpoint admin
//     exigía `academic.subjects.read`, podado del teacher en F3 → 403; el nuevo usa
//     `academic.my_teaching.read:own` y devuelve solo las materias que el docente dicta).
//     Solo datos de seed L4 (L4_SEED_VERSION 1.80.0 → 1.81.0). Recrear BD, sin ALTER.
//   - 3.92.0 (plan 032 B1b): nueva tabla content.material_assignment — puente
//     material → oferta (subject_offering), calca assessment.assessment_assignment.
//     Entity content/material_assignment.go + registro en gorm_migrator.go; FKs
//     cross/same-schema (material/offering/membership), UNIQUE
//     (material_id, subject_offering_id) y trigger set_updated_at en post_gorm.sql.
//     DDL aditiva (cambia el hash de post_gorm.sql) → bump obligatorio. Recrear BD.
//   - 3.93.0 (bug 0081): el recurso `assessments_student` ("Tomar Evaluación") tenía
//     `assessments-list` como pantalla default, que pega al endpoint del docente
//     (GET /assessments, `content.assessments.read`) → 403 al alumno. Se elimina ese
//     mapping del recurso del estudiante y `assigned-assessments-list`
//     (GET /me/assigned-assessments) queda como única pantalla (isDefault). Solo datos
//     de seed L4 (L4_SEED_VERSION 1.81.0 → 1.82.0) → cambia el hash de seeds → bump
//     obligatorio. Recrear BD, sin ALTER.
//   - 3.94.0 (plan 033 B2a): nueva tabla content.user_material_tags — etiquetas
//     personales por usuario (D-B2.6). Entity content/user_material_tag.go +
//     registro en gorm_migrator.go. Sin FK dura (user_id vive en auth,
//     material_id en content); UNIQUE (user_id, material_id, tag) e indice
//     idx_user_material_tags_user via tags GORM (no toca pre/post_gorm.sql).
//     Cambio solo en entities → el hash no cambia, pero la regla 1 exige el bump.
//     Aplicada DIRECTO a Neon de forma aditiva (decision del dueno 2026-07-12).
//   - 3.95.0 (plan 032 B2a): assessment.is_public (catalogo de evaluaciones) —
//     nueva columna boolean not null default false en assessment.assessment, que
//     distingue disponibilidad/catalogo vs distribucion por grupos. Campo IsPublic
//     en entities/assessment.go via tag GORM (no toca pre/post_gorm.sql → el hash
//     no cambia, pero la regla 1 exige el bump). Aplicada DIRECTO a Neon de forma
//     aditiva (AutoMigrate focalizado, sin FORCE).
//   - 3.96.0 (plan 032 B2 — toggle catálogo por SDUI): dos acciones de toolbar
//     nuevas en el screen_instance `assessments-form` (`publicar-catalogo` /
//     `quitar-catalogo`, event_id publish-catalog / unpublish-catalog) que togglean
//     assessment.is_public por SDUI. Extiende el shape de `visible_when` de forma
//     aditiva (objeto simple → o LISTA de objetos en AND). Solo datos de seed L4
//     (L4_SEED_VERSION 1.82.0 → 1.83.0) → cambia el hash de seeds → bump
//     obligatorio. El slot_data del row se aplica DIRECTO a Neon de forma aditiva
//     (UPDATE del row, sin recrear, sin FORCE).
//   - 3.97.0 (plan 036 D-036.4): nueva columna teacher_feedback text nullable en
//     assessment.assessment_attempt — comentario global del profesor al finalizar
//     la revision. Campo TeacherFeedback en entities/assessment_attempt.go via tag
//     GORM (no toca pre/post_gorm.sql → el hash no cambia, pero la regla 1 exige
//     el bump). Solo LOCAL en esta tarea; Neon lo aplica otro paso del plan.
//   - 3.98.0 (plan 035 F1a — capa de práctica): diferenciación explícita
//     purpose + capa de trazas de práctica. En entities/assessment.go se ELIMINA
//     la columna `kind` (D-035.2: seed = verdad, sin columnas muertas) y se
//     agregan `purpose` varchar(20) not null default 'exam' CHECK in
//     (practice,exam,both) — D-035.1 — y `passing_score` smallint not null
//     default 60 CHECK 0..100 — D-035.8, umbral que bloquea el reintento del
//     examen. 3 tablas nuevas en el schema assessment (plano de práctica,
//     D-035.4): `practice_session` (cabecera del log, assessment_id nullable ON
//     DELETE SET NULL), `practice_session_answer` (detalle append-only, session_id
//     CASCADE / question_id SET NULL) y `user_question_stat` (acumulador acotado
//     del alumno: UNIQUE (membership_id, question_id), índices (school_id,
//     membership_id, subject_id) y (membership_id, next_review_at) para el SRS de
//     F2, question_id SET NULL). Las 3 entities registradas en el AutoMigrate; las
//     FKs cross-schema/SET NULL y los triggers set_updated_at de practice_session
//     y user_question_stat en post_gorm.sql → ComputeFilesHash() CAMBIA. Solo
//     LOCAL en esta tarea; Neon lo aplica aditivo (F1b) otro paso del plan.
//   - 3.99.0 (plan 037 F1g — worker a dieta): se ELIMINA la tabla
//     academic.practice_result, deprecada por el plan 036 (D-036.3). Era el
//     espejo de academic.grade_item para evaluaciones de práctica (resultado
//     fuera del expediente); ya no se materializa como tabla propia — la
//     trazabilidad de práctica vive en el plano assessment.practice_session /
//     practice_session_answer / user_question_stat (3.98.0). Se borra
//     entities/practice_result.go, su registro en el AutoMigrate (gorm_migrator.go)
//     y sus 3 bloques en post_gorm.sql (6 FKs, trigger set_updated_at, índice
//     parcial uq_practice_result_attempt) → ComputeFilesHash() CAMBIA → bump
//     obligatorio. Los índices GORM (idx_practice_result_grain/_attempt) se van
//     con la entity. Solo LOCAL: la BD nace sin la tabla al recrear; ninguna
//     dependencia viva la lee (el único uso era el worker por `kind`, ya retirado).
//   - 3.100.0 (plan 038, import evaluaciones): source_type admite 'imported'
//     (import externo de evaluaciones vía JSON, Riel 0); valor aditivo, default
//     sigue 'manual'. Cambio SOLO en el tag GORM de entities/assessment.go (CHECK
//     assessment_source_type_check + validate oneof); no toca pre/post_gorm.sql →
//     ComputeFilesHash() NO cambia. Bump obligatorio por la regla 1 (cambio en
//     entity). AutoMigrate no altera CHECKs existentes: la ampliación efectiva del
//     CHECK llega al recrear la BD (paso posterior coordinado, no en esta tarea).
//   - 3.101.0 (plan 039 F1 — terreno LLM, config por escuela): nueva tabla
//     academic.school_settings (entities/school_setting.go), configuración
//     clave/valor POR ESCUELA (D-039.1). PK compuesta (school_id, key); value
//     varchar not null; created_at/updated_at. Registrada en el AutoMigrate
//     (gorm_migrator.go, tras School por la FK). La FK school_id →
//     academic.schools ON DELETE CASCADE y el trigger set_updated_at viven en
//     post_gorm.sql (GORM no materializa la FK sin campo de relación) →
//     ComputeFilesHash() CAMBIA. El catálogo de claves válidas vive en código
//     (entities/school_setting_catalog.go): llm.generation.mode, llm.review.mode,
//     llm.review.flow (enums local|api|off / direct|teacher, defaults off/off/
//     teacher) e import.max_questions / import.max_json_bytes (int, defaults 100 /
//     1 MiB, env EDUGO_IMPORT_MAX_*). Seeds: SchoolSpec gana Settings opcional y
//     el fixture base siembra settings explícitos a San Ignacio (llm.review.mode=
//     api, llm.review.flow=teacher); la otra escuela queda sin filas (prueba la
//     resolución por default). Recrear BD, sin ALTER.
//   - 3.102.0 (plan 040 F0 — corrección IA prevalidada): dos cambios de contrato
//     en el schema assessment. (1) En entities/assessment_attempt_answer.go el
//     CHECK assessment_attempt_answer_review_status_check gana el valor
//     'ai_reviewed' (review_status IN pending,auto_graded,reviewed,ai_reviewed +
//     validate oneof); estado que marca la respuesta corregida por IA. (2) En
//     entities/attempt_review.go nace la columna review_source varchar(20) not
//     null default 'teacher' CHECK attempt_review_source_check IN (teacher,llm)
//     — materializa ADR 0033, distingue revisión manual vs corrección IA; filas
//     históricas quedan 'teacher' por default. Ambos cambios viven en tags GORM
//     (CHECKs inline) → no tocan pre/post_gorm.sql → ComputeFilesHash() NO cambia;
//     bump obligatorio por la regla 1 (cambio en entity). AutoMigrate no altera
//     CHECKs existentes: la ampliación efectiva del CHECK de review_status llega
//     al recrear la BD local. En Neon el CHECK requiere ALTER manual (drop+add
//     constraint) en el paso de despliegue coordinado, no en esta tarea.
//   - 3.103.0 (plan 040 F3 — template dedicado review-dashboard-v1): solo datos
//     de seed L4 (L4_SEED_VERSION 1.84.0 → 1.85.0) → cambia el hash de seeds →
//     bump obligatorio, sin DDL. Nace el template SDUI dedicado
//     review-dashboard-v1 (pattern "list", UUID a4000000-…-007) con zona de
//     filtros de 4 CHIP slots custom (filter_all / filter_pending_review /
//     filter_ai_reviewed / filter_completed) + lista student_name/status; la
//     instancia assessment-review-dashboard se repunta de dashboard-basic-v1 a
//     este template y suma los labels de chip (incl. «Prevalidado IA»). Habilita
//     que el contrato KMP renderice el chip de intentos prevalidados por IA.
//     Recrear BD para reseeding.
//   - 3.104.0 (plan 040 T3c — navegación a review-dashboard + fix detalle
//     grades-list): solo datos de seed L4 (L4_SEED_VERSION 1.85.0 → 1.86.0) →
//     cambia el hash de seeds → bump obligatorio, sin DDL. (P2a) La instancia
//     assessments-management-list gana la row-action `review-results` (scope
//     row, permiso content.assessments.read) que expone el handler KMP ya
//     existente y da el único punto de entrada a assessment-review-dashboard,
//     antes inalcanzable. (P2b) El detalle roto de grades-list (el contrato
//     KMP navegaba a grades-form, ELIMINADA 2026-06-09 → 404) se corrige en
//     KMP (GradesListContract: select-item → NoOp; grades-list es read-only y
//     no hay pantalla de detalle docente), fuera de este seed. Recrear BD para
//     reseeding.
//   - 3.105.0 (T4-1 — candado «en revisión por IA» con vencimiento): nueva
//     columna assessment.assessment_attempt.ai_review_claimed_at (timestamptz
//     NULL, entities/assessment_attempt.go, campo AIReviewClaimedAt *time.Time).
//     Marca cuándo un proceso de revisión por IA tomó el candado sobre el intento
//     (NULL = sin candado activo). Cambio SOLO en tag GORM del entity (columna
//     aditiva NULL) → no toca pre/post_gorm.sql → ComputeFilesHash() NO cambia;
//     bump obligatorio por la regla 1 (cambio en entity). AutoMigrate agrega la
//     columna nueva al recrear/migrar la BD local; en Neon el ADD COLUMN es
//     aditivo y llega en el paso de despliegue coordinado, no en esta tarea.
//   - 3.106.0 (plan 042 F0 — artefacto «llm_prep» por pregunta): cinco columnas
//     aditivas en assessment.question (entities/question.go):
//     llm_prep (jsonb NULL, artefacto de preparacion, contrato versionado v1),
//     llm_prep_status (varchar(20) NOT NULL DEFAULT 'none', CHECK inline
//     none|pending|processing|ready|failed|stale), llm_prep_source_hash
//     (varchar(64) NULL, sha256 de question_type+question_text+correct_answer+
//     explanation), llm_prep_feedback (text NULL, comentario del profesor que se
//     consume 1 vez) y llm_prep_updated_at (timestamptz NULL). Cambio SOLO en el
//     entity (columnas aditivas, CHECK en tag GORM) → no toca pre/post_gorm.sql →
//     ComputeFilesHash() NO cambia; bump obligatorio por la regla 1 (cambio en
//     entity). AutoMigrate agrega las columnas al recrear/migrar la BD local; en
//     Neon los ADD COLUMN son aditivos y llegan en el paso de despliegue
//     coordinado, no en esta tarea.
const SchemaVersion = "3.106.0"

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
