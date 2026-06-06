package l4

// UUIDs hardcodeados de las screen_instances sembradas por L4 (B4).
// Prefijo `b4400000-*` por convención de fase + sub-bloque (B4 = 4-ésimo
// bloque). Diferenciado del bloque b4000000-* usado en B1 (resources)
// y de b4...000003 reservado en B2 para role_permissions.
//
// El FE NO hardcodea UUIDs de screen_instances (resuelve por
// `screen_key`), así que se pueden regenerar libremente respecto al
// legacy (ADR-6 §6). Verificación: `grep -rn "[0-9a-f]\{8\}-...\
// screen_instance" /EduUI/edugo-ui-kmp/` no produce hits.
//
// Bloques mnemotécnicos del último segmento (legibilidad humana, sin
// dependencia semántica con el legacy):
//
//	0001..0009 → auth & shell (app-login, app-settings)
//	0010..001F → dashboards por rol
//	0020..002F → admin: users, schools, roles, permissions
//	0030..003F → admin: screens, system_settings, concept_types
//	0040..004F → admin: audit
//	0050..005F → academic: units, memberships, subjects, periods
//	0060..006F → academic: guardian_relations, calendar, schedules
//	0070..007F → academic: grades, attendance
//	0080..009F → content: assessments + take/result/review/assignment
//	00A0..00AF → reports: progress, stats, report-card
//	00B0..00BF → guardian: children/requests
//	00C0..00CF → notifications & directories
//	00D0..00DF → phantom-nuevas (B4: school-concepts, user-roles)
const (
	// auth & shell
	L4_SCREEN_INST_APP_LOGIN_ID      = "b4400000-0000-0000-0000-000000000001"
	L4_SCREEN_INST_APP_SETTINGS_ID   = "b4400000-0000-0000-0000-000000000002"
	L4_SCREEN_INST_DASHBOARD_HOME_ID = "b4400000-0000-0000-0000-000000000003"

	// dashboards (5 roles + 2 specialized)
	L4_SCREEN_INST_DASH_TEACHER_ID    = "b4400000-0000-0000-0000-000000000010"
	L4_SCREEN_INST_DASH_STUDENT_ID    = "b4400000-0000-0000-0000-000000000011"
	L4_SCREEN_INST_DASH_SUPERADMIN_ID = "b4400000-0000-0000-0000-000000000012"
	L4_SCREEN_INST_DASH_SCHOOLADM_ID  = "b4400000-0000-0000-0000-000000000013"
	L4_SCREEN_INST_DASH_GUARDIAN_ID   = "b4400000-0000-0000-0000-000000000014"
	L4_SCREEN_INST_PROGRESS_DASH_ID   = "b4400000-0000-0000-0000-000000000015"
	L4_SCREEN_INST_STATS_DASH_ID      = "b4400000-0000-0000-0000-000000000016"

	// admin: users / schools
	// Poda menú (2026-05-29): L4_SCREEN_INST_ROLES_*/PERMISSIONS_* (…24..…27)
	// eliminadas junto con los recursos `roles`/`permissions_mgmt`; UUIDs libres.
	L4_SCREEN_INST_USERS_LIST_ID   = "b4400000-0000-0000-0000-000000000020"
	L4_SCREEN_INST_USERS_FORM_ID   = "b4400000-0000-0000-0000-000000000021"
	L4_SCREEN_INST_SCHOOLS_LIST_ID = "b4400000-0000-0000-0000-000000000022"
	L4_SCREEN_INST_SCHOOLS_FORM_ID = "b4400000-0000-0000-0000-000000000023"

	// admin: system-settings, concept-types
	// Poda menú (2026-06-01): las instancias …30..…33 (screen-templates-list,
	// screen-instances-list/form, screens-form) se eliminaron junto con los
	// recursos `screen_templates`/`screen_instances` — ese CRUD de
	// configuración SDUI se reimplementó en el admin-tool de Go. UUIDs libres.
	L4_SCREEN_INST_SYSTEM_SETTINGS_ID    = "b4400000-0000-0000-0000-000000000034"
	L4_SCREEN_INST_CONCEPT_TYPES_LIST_ID = "b4400000-0000-0000-0000-000000000035"
	L4_SCREEN_INST_CONCEPT_TYPES_FORM_ID = "b4400000-0000-0000-0000-000000000036"

	// admin: audit
	L4_SCREEN_INST_AUDIT_LIST_ID   = "b4400000-0000-0000-0000-000000000040"
	L4_SCREEN_INST_AUDIT_DETAIL_ID = "b4400000-0000-0000-0000-000000000041"

	// academic: structure (units, memberships, subjects, periods)
	L4_SCREEN_INST_UNITS_LIST_ID       = "b4400000-0000-0000-0000-000000000050"
	L4_SCREEN_INST_UNITS_FORM_ID       = "b4400000-0000-0000-0000-000000000051"
	L4_SCREEN_INST_MEMBERSHIPS_LIST_ID = "b4400000-0000-0000-0000-000000000052"
	// memberships-form queda como pantalla de SOLO EDICIÓN (editar una membresía
	// existente). La creación directa se retiró (sin FAB de crear, sin POST, sin
	// membership-add), pero el form de edición se conserva — lo alcanza la acción
	// "editar" de memberships-list (LOAD_DATA por id + submit PUT).
	L4_SCREEN_INST_MEMBERSHIPS_FORM_ID = "b4400000-0000-0000-0000-000000000053"
	// "Mis materias" del alumno (plan 006, N1.C ETAPA 2). Reintroducida en
	// N1.7 F1 sobre el modelo de sesiones: el contrato KMP consume
	// GET /api/v1/me/subject-offerings (lector A); el seed solo declara
	// columnas/título/permiso. Permiso propio academic.my_memberships.read:own.
	L4_SCREEN_INST_MY_MEMBERSHIPS_LIST_ID = "b4400000-0000-0000-0000-00000000005b"
	// "Mis notas" del alumno (N3 F4, consulta de notas). Lista readonly de
	// las notas propias por sesión de materia. El contrato KMP consume
	// GET /api/v1/me/grades; el seed solo declara columnas/título/permiso.
	// Permiso propio academic.my_grades.read:own. Espejo de my-memberships-list.
	// Sufijo …0068 (slot libre del bloque academic structure; los …0060..0067
	// quedaron libres tras la poda de guardian).
	L4_SCREEN_INST_MY_GRADES_LIST_ID = "b4400000-0000-0000-0000-000000000068"
	L4_SCREEN_INST_SUBJECTS_LIST_ID  = "b4400000-0000-0000-0000-000000000054"
	L4_SCREEN_INST_SUBJECTS_FORM_ID       = "b4400000-0000-0000-0000-000000000055"
	L4_SCREEN_INST_PERIODS_LIST_ID        = "b4400000-0000-0000-0000-000000000056"
	L4_SCREEN_INST_PERIODS_FORM_ID        = "b4400000-0000-0000-0000-000000000057"
	L4_SCREEN_INST_INVITATIONS_LIST_ID    = "b4400000-0000-0000-0000-000000000058"
	L4_SCREEN_INST_INVITATIONS_FORM_ID    = "b4400000-0000-0000-0000-000000000059"
	// join-requests-inbox (N0.4-B): instancia mínima requerida por la FK
	// resource_screens.screen_key → screen_instances.screen_key. La
	// pantalla se renderiza NATIVA en el FE (no SDUI), así que el
	// slot_data nunca se usa para render; existe solo para satisfacer la
	// FK y permitir que el menú resuelva el screen_key.
	L4_SCREEN_INST_JOIN_REQUESTS_INBOX_ID = "b4400000-0000-0000-0000-00000000005a"
	// batch-enroll (N1.7 F1): pantalla NATIVA de "inscripción por lote" de
	// alumnos en una sesión de materia (subject_offering). Igual que
	// join-requests-inbox, esta screen_instance existe SOLO para satisfacer la
	// FK resource_screens.screen_key → screen_instances.screen_key y para que el
	// menú resuelva el screen_key `batch-enroll`. El slot_data NO se renderiza
	// por el SDUI engine: el FE intercepta el screen_key y pinta la pantalla
	// Compose nativa. Sufijo …5c, adyacente a my_memberships …5b.
	L4_SCREEN_INST_SUBJECT_OFFERINGS_BATCH_ENROLL_ID = "b4400000-0000-0000-0000-00000000005c"
	// enroll-one (N1.7 F2): pantalla NATIVA de "inscripción individual" de un
	// alumno en una sesión de materia (subject_offering). Igual que batch-enroll,
	// existe SOLO para satisfacer la FK resource_screens.screen_key →
	// screen_instances.screen_key y para que el menú/handler resuelva el
	// screen_key `enroll-one`. El slot_data NO se renderiza por el SDUI engine: el
	// FE intercepta el screen_key y pinta la pantalla Compose nativa. Sufijo …5d,
	// adyacente a batch-enroll …5c.
	L4_SCREEN_INST_ENROLL_ONE_ID = "b4400000-0000-0000-0000-00000000005d"
	// sessions-by-subject-list (N1.7 F2; reubicada en F2.2): lista hija de
	// "sesiones de la materia". Se alcanza embebida como pestaña "Sesiones" del
	// master-detail subjects-form (detail_configs[]); el contenedor inyecta
	// subjectId y la pantalla consume GET /api/v1/subject-offerings?subject_id={id}.
	// Pantalla SDUI list estándar (columns subject/section/period/teacher).
	// Sufijo …5e, adyacente a enroll-one …5d.
	L4_SCREEN_INST_SESSIONS_BY_SUBJECT_ID = "b4400000-0000-0000-0000-00000000005e"
	// sessions-by-subject-form (N1.7 F2.3): formulario crear/editar de "sesión de
	// materia". Se renderiza como modal del master-detail subjects-form (la
	// pestaña "Sesiones" lo declara en detail_configs[].modal_screen_key). El
	// contenedor inyecta subjectId (parent) y, en edición, id; el contrato KMP
	// inyecta subject_id al body en create y consume POST/PUT
	// /api/v1/subject-offerings. Pantalla SDUI form estándar. Sufijo …5f,
	// adyacente a sessions-by-subject-list …5e.
	L4_SCREEN_INST_SESSIONS_BY_SUBJECT_FORM_ID = "b4400000-0000-0000-0000-00000000005f"

	// academic: guardian / calendar / schedules
	// Poda F2 (plan 004-permisologia-mvp): retiradas las constantes de
	// guardian (rel-list/form, alias, req-list) 0060..0063, calendar
	// (list/form) 0064..0065 y schedules (list/form) 0066..0067 junto
	// con sus constructores y filas en resource_screens.go. UUIDs
	// 0060..0067 quedan libres para reuso futuro.

	// academic: grades / attendance
	L4_SCREEN_INST_GRADES_LIST_ID        = "b4400000-0000-0000-0000-000000000070"
	L4_SCREEN_INST_GRADES_FORM_ID        = "b4400000-0000-0000-0000-000000000071"
	L4_SCREEN_INST_ATTENDANCE_LIST_ID    = "b4400000-0000-0000-0000-000000000072"
	L4_SCREEN_INST_ATTENDANCE_BATCH_ID   = "b4400000-0000-0000-0000-000000000073"
	L4_SCREEN_INST_GRADES_BATCH_ID       = "b4400000-0000-0000-0000-000000000074"
	L4_SCREEN_INST_ATTENDANCE_SUMMARY_ID = "b4400000-0000-0000-0000-000000000075"
	// grades-subject-summary (N3 F4, consulta de notas): resumen de notas por
	// sesión (vista del docente). Espejo de attendance-summary. Sufijo …0076,
	// adyacente a attendance-summary …0075.
	L4_SCREEN_INST_GRADES_SUBJECT_SUMMARY_ID = "b4400000-0000-0000-0000-000000000076"

	// content: assessments (gestión + estudiante)
	L4_SCREEN_INST_ASSESS_LIST_ID           = "b4400000-0000-0000-0000-000000000080"
	L4_SCREEN_INST_ASSESS_FORM_ID           = "b4400000-0000-0000-0000-000000000081"
	L4_SCREEN_INST_ASSESS_MGMT_LIST_ID      = "b4400000-0000-0000-0000-000000000082"
	L4_SCREEN_INST_ASSESS_TAKE_ID           = "b4400000-0000-0000-0000-000000000083"
	L4_SCREEN_INST_ASSESS_RESULT_ID         = "b4400000-0000-0000-0000-000000000084"
	L4_SCREEN_INST_ASSESS_QUESTIONS_LIST_ID = "b4400000-0000-0000-0000-000000000085"
	L4_SCREEN_INST_ASSESS_QUESTION_FORM_ID  = "b4400000-0000-0000-0000-000000000086"
	L4_SCREEN_INST_ASSESS_ASSIGNMENT_ID     = "b4400000-0000-0000-0000-000000000087"
	L4_SCREEN_INST_ASSESS_MODALITY_ID       = "b4400000-0000-0000-0000-000000000088"
	L4_SCREEN_INST_ASSESS_REVIEW_DASH_ID    = "b4400000-0000-0000-0000-000000000089"
	L4_SCREEN_INST_ASSESS_ASSIGNED_LIST_ID  = "b4400000-0000-0000-0000-00000000008a"
	L4_SCREEN_INST_ATTEMPT_REVIEW_DETAIL_ID = "b4400000-0000-0000-0000-00000000008b"

	// reports
	// Poda F2 (plan 004-permisologia-mvp): retiradas las constantes de
	// progress-detail 00a0, stats-detail 00a1 y report-card 00a2 junto
	// con sus constructores y filas en resource_screens.go. UUIDs
	// 00a0..00a2 quedan libres para reuso futuro.

	// directories & misc
	L4_SCREEN_INST_UNIT_DIRECTORY_ID = "b4400000-0000-0000-0000-0000000000c0"
	// L4_SCREEN_INST_STUDENTS_BY_SUBJECT_ID (UUID …c1) eliminado (2026-06-02):
	// su screen_instance `students-by-subject-list` era SOLO el panel detalle
	// "Alumnos" embebido en subjects-form, retirado porque un alumno se inscribe
	// en una SESIÓN, no en la materia. UUID …c1 queda libre para reuso futuro.
	// L4_SCREEN_INST_NOTIFICATIONS_LIST_ID retirado en B7-fix: FE no
	// implementa NotificationsListContract.kt aún. Re-sembrar al
	// agregar el Contract correspondiente.

	// phantom-nuevas (B4): no presentes en legacy, exigidas por el FE
	L4_SCREEN_INST_SCHOOL_CONCEPTS_LIST_ID = "b4400000-0000-0000-0000-0000000000d0"
	L4_SCREEN_INST_SCHOOL_CONCEPTS_FORM_ID = "b4400000-0000-0000-0000-0000000000d1"
	// L4_SCREEN_INST_MEMBERSHIP_ADD_ID (…d2) retirado: la creación directa de
	// membresías se eliminó (redundante con invitación→solicitud→aprobación).
	// UUID …d2 queda libre para reuso futuro.
	L4_SCREEN_INST_USER_ROLES_ID = "b4400000-0000-0000-0000-0000000000d3"

	// Fase 3 (B7b) — demo CRUD data-driven sin Kotlin.
	// Poda F2 (plan 004-permisologia-mvp): retiradas las constantes de
	// la pareja demo colors-list 00e0 / colors-form 00e1 junto con sus
	// constructores y filas en resource_screens.go. UUIDs 00e0..00e1
	// quedan libres para reuso futuro.
)

// UUIDs L0 que el bloque B4 necesita referenciar como template_id.
// Duplicados aquí para evitar el ciclo `seeds/system/l4 →
// seeds/system/layers`. Origen autoritativo:
// `seeds/system/layers/l0_constants.go`. Mismo patrón usado por B2 en
// roles_permissions_constants.go.
const (
	L0_SCREEN_TPL_LIST_ID_REF          = "30000000-0000-0000-0000-000000000001"
	L0_SCREEN_TPL_DETAIL_ID_REF        = "30000000-0000-0000-0000-000000000002"
	L0_SCREEN_TPL_FORM_ID_REF          = "30000000-0000-0000-0000-000000000003"
	L0_SCREEN_TPL_MASTER_DETAIL_ID_REF = "30000000-0000-0000-0000-000000000004"
)
