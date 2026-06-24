package l4

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ApplyResourceScreens siembra los mappings recurso↔pantalla
// (`ui_config.resource_screens`) que cierran el sistema. La tabla
// resuelve "dada una resource_key + screen_type, ¿qué screen_key
// ejecuto?" — la usa el resolver para responder al FE cuando éste
// pide la pantalla "default" de un recurso del menú.
//
// Excluidos (ya sembrados por capas anteriores):
//   - announcements:list  → applyL0ResourceScreens (L0)
//   - announcements:form  → applyL2ResourceScreens (L2)
//   - materials:list      → applyL3ResourceScreens (L3) — pantalla nativa
//     sin ScreenInstance. (materials:form fue podado junto con su
//     ScreenInstance: poda SDUI material 2026-06-07; ver L3.)
//
// Recursos sembrados por L4 (B1) sin pantalla asociada — intencional:
//   - admin, academic, content       → contenedores puros del árbol
//     del menú; el menú expande los hijos, no resuelve un screen_key
//     propio.
//   - screens, context, menu         → recursos "API-only"
//     (IsMenuVisible=false). El FE no pide screen_config para ellos.
//   - reports                        → contenedor del subárbol
//     reports/{stats}. report-card mapea a `reports` como
//     pantalla legacy de boleta, manteniéndose como default.
//
// Tabla UNIQUE: (resource_id, screen_type). Por eso varias pantallas
// del mismo recurso usan screen_types distintos (ver dashboard,
// assessments). is_default=true EXACTAMENTE en una fila por recurso
// (la pantalla que abre al hacer click en el item de menú).
//
// Cambios vs `legacy/data.go::resourceScreenSeedRows` (65 filas
// originales):
//   - +6 mappings nuevos para las screen_keys phantom seedadas por B4:
//     user-roles, membership-add, school-concepts-{list,form},
//     assessment-review-dashboard, attempt-review-detail,
//     assigned-assessments-list, notifications-list.
//     (assessment-modality se eliminó en plan 015 y assessment-assignment
//     se reemplazó por un modal nativo.)
//     (Cada phantom legítima del baseline cross-checker queda
//     accesible vía menú o sub-flujo.)
//   - -3 mappings descartados por dead-screens (cross-checker
//     `screen_key_dead`): material-detail, children-list,
//     child-progress (sus screen_instances no existen en B4 porque
//     el FE no las implementa).
//   - assessments-management-list: legacy lo declaraba 3 veces con
//     screen_types role-acoplados ("superadmin", "teacher",
//     "schoolcoordinator"). L4 consolida a una sola fila con
//     screen_type="management-list" — la visibilidad por rol se
//     controla por permisos (assessments:read), no duplicando la
//     fila en resource_screens. Cambio documentado en
//     phase-6-layer-l4/decisions-log.md.
//   - guardian_relations-form: alias underscore del FE legacy. Mismo
//     razonamiento — `guardian-relations-form` ya ocupa
//     (guardian_relations, form). NO se mapea el alias.
//   - app-login / app-settings: pantallas shell del sistema. No están
//     atadas a un recurso del menú — el FE las resuelve por
//     screen_key. NO se mapean.
//
// is_default por recurso (sanidad): exactamente 1 fila con
// is_default=true por cada recurso que tiene al menos una pantalla.
// Recursos contenedores y API-only sin pantalla: sin default
// (correcto, F6 acepta).
//
// UUIDs propios bajo prefijo b4500000-* (ADR-6 §6): el FE no
// referencia los UUIDs de filas de resource_screens (resuelve por
// resource_key+screen_type), por lo que se generan de cero.
//
// Idempotencia: UPSERT con conflict target (resource_id, screen_type)
// — UNIQUE definido en entities/resource_screen.go. DoNothing para
// no pisar customizaciones en environments live (mismo criterio que
// applyL2ResourceScreens / applyL3ResourceScreens).
func ApplyResourceScreens(tx *gorm.DB) error {
	mappings := buildL4ResourceScreens()

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "resource_id"}, {Name: "screen_type"}},
		DoNothing: true,
	}).CreateInBatches(&mappings, 50).Error; err != nil {
		return fmt.Errorf("ApplyResourceScreens: upsert resource_screens: %w", err)
	}
	return nil
}

// buildL4ResourceScreens materializa las filas literales en
// entities.ResourceScreen. Helper compartido por ApplyResourceScreens y
// por el accessor público l4.ResourceScreens().
func buildL4ResourceScreens() []entities.ResourceScreen {
	rows := l4ResourceScreens()
	mappings := make([]entities.ResourceScreen, 0, len(rows))
	for _, r := range rows {
		mappings = append(mappings, entities.ResourceScreen{
			ID:          mustParseL4UUID(r.id, "resource_screen:"+r.resourceKey+":"+r.screenType),
			ResourceID:  mustParseL4UUID(r.resourceID, "resource_screen.resource_id:"+r.resourceKey),
			ResourceKey: r.resourceKey,
			ScreenKey:   r.screenKey,
			ScreenType:  r.screenType,
			IsDefault:   r.isDefault,
			SortOrder:   r.sortOrder,
			IsActive:    true,
		})
	}
	return mappings
}

// l4ResourceScreenRow describe una fila de ui_config.resource_screens
// con primitivos para simplificar la inspección de la tabla literal.
// La conversión a entities.ResourceScreen ocurre en ApplyResourceScreens.
type l4ResourceScreenRow struct {
	id          string
	resourceID  string
	resourceKey string
	screenKey   string
	screenType  string
	isDefault   bool
	sortOrder   int
}

// l4ResourceScreens retorna las 65 filas de mappings recurso↔pantalla
// sembradas por L4. Agrupadas por recurso, con el default explícito al
// frente del bloque para revisión visual.
func l4ResourceScreens() []l4ResourceScreenRow {
	return []l4ResourceScreenRow{
		// =============================================================
		// dashboard (5 dashboards por rol; screen_types distintos para
		// satisfacer UNIQUE (resource_id, screen_type)).
		// Default: teacher (legacy preservado — el primer rol que ve
		// la app post-login es teacher en el flujo de onboarding del
		// MVP). Los demás dashboards resuelven cuando el usuario
		// loguea con el rol correspondiente.
		// =============================================================
		{id: "b4500000-0000-0000-0000-000000000001", resourceID: L4_RESOURCE_DASHBOARD_ID, resourceKey: L4_RESOURCE_DASHBOARD_KEY, screenKey: "dashboard-teacher", screenType: "dashboard", isDefault: true, sortOrder: 1},
		{id: "b4500000-0000-0000-0000-000000000002", resourceID: L4_RESOURCE_DASHBOARD_ID, resourceKey: L4_RESOURCE_DASHBOARD_KEY, screenKey: "dashboard-student", screenType: "dashboard-student", isDefault: false, sortOrder: 2},
		{id: "b4500000-0000-0000-0000-000000000003", resourceID: L4_RESOURCE_DASHBOARD_ID, resourceKey: L4_RESOURCE_DASHBOARD_KEY, screenKey: "dashboard-superadmin", screenType: "dashboard-superadmin", isDefault: false, sortOrder: 3},
		{id: "b4500000-0000-0000-0000-000000000004", resourceID: L4_RESOURCE_DASHBOARD_ID, resourceKey: L4_RESOURCE_DASHBOARD_KEY, screenKey: "dashboard-schooladmin", screenType: "dashboard-schooladmin", isDefault: false, sortOrder: 4},
		{id: "b4500000-0000-0000-0000-000000000005", resourceID: L4_RESOURCE_DASHBOARD_ID, resourceKey: L4_RESOURCE_DASHBOARD_KEY, screenKey: "dashboard-guardian", screenType: "dashboard-guardian", isDefault: false, sortOrder: 5},

		// =============================================================
		// admin → users / schools / roles / permissions_mgmt
		// =============================================================
		{id: "b4500000-0000-0000-0000-000000000010", resourceID: L4_RESOURCE_USERS_ID, resourceKey: L4_RESOURCE_USERS_KEY, screenKey: "users-list", screenType: "list", isDefault: true, sortOrder: 1},
		{id: "b4500000-0000-0000-0000-000000000011", resourceID: L4_RESOURCE_USERS_ID, resourceKey: L4_RESOURCE_USERS_KEY, screenKey: "users-form", screenType: "form", isDefault: false, sortOrder: 2},
		// user-roles ELIMINADA (2026-06-09): mapping (screen_type "roles", UUID
		// …0012) retirado junto con su screen_instance — pantalla SDUI legacy
		// huérfana sin entry-point. UUID …0012 queda libre para reuso futuro.

		{id: "b4500000-0000-0000-0000-000000000015", resourceID: L4_RESOURCE_SCHOOLS_ID, resourceKey: L4_RESOURCE_SCHOOLS_KEY, screenKey: "schools-list", screenType: "list", isDefault: true, sortOrder: 1},
		{id: "b4500000-0000-0000-0000-000000000016", resourceID: L4_RESOURCE_SCHOOLS_ID, resourceKey: L4_RESOURCE_SCHOOLS_KEY, screenKey: "schools-form", screenType: "form", isDefault: false, sortOrder: 2},

		// Poda menú (2026-05-29): mappings de `roles` (roles-list/form) y
		// `permissions_mgmt` (permissions-list/form) eliminados junto con sus
		// recursos y screen_instances.

		// Poda menú (2026-06-01): mappings de `screen_templates`
		// (screen-templates-list) y `screen_instances` (screen-instances-list/
		// form, screens-form) eliminados junto con sus recursos y
		// screen_instances — ese CRUD de configuración SDUI se reimplementó en
		// el admin-tool de Go.

		// =============================================================
		// admin → audit / concept_types / system_settings
		// =============================================================
		{id: "b4500000-0000-0000-0000-000000000040", resourceID: L4_RESOURCE_AUDIT_ID, resourceKey: L4_RESOURCE_AUDIT_KEY, screenKey: "audit-events-list", screenType: "list", isDefault: true, sortOrder: 1},
		{id: "b4500000-0000-0000-0000-000000000041", resourceID: L4_RESOURCE_AUDIT_ID, resourceKey: L4_RESOURCE_AUDIT_KEY, screenKey: "audit-detail", screenType: "detail", isDefault: false, sortOrder: 2},

		{id: "b4500000-0000-0000-0000-000000000045", resourceID: L4_RESOURCE_CONCEPT_TYPES_ID, resourceKey: L4_RESOURCE_CONCEPT_TYPES_KEY, screenKey: "concept-types-list", screenType: "list", isDefault: true, sortOrder: 1},
		{id: "b4500000-0000-0000-0000-000000000046", resourceID: L4_RESOURCE_CONCEPT_TYPES_ID, resourceKey: L4_RESOURCE_CONCEPT_TYPES_KEY, screenKey: "concept-types-form", screenType: "form", isDefault: false, sortOrder: 2},
		// school-concepts-{list,form} (phantom B4): variantes school-
		// scoped del CRUD de concept_types. Sub-flujo en el form de
		// escuela; no default.
		{id: "b4500000-0000-0000-0000-000000000047", resourceID: L4_RESOURCE_CONCEPT_TYPES_ID, resourceKey: L4_RESOURCE_CONCEPT_TYPES_KEY, screenKey: "school-concepts-list", screenType: "school-list", isDefault: false, sortOrder: 3},
		{id: "b4500000-0000-0000-0000-000000000048", resourceID: L4_RESOURCE_CONCEPT_TYPES_ID, resourceKey: L4_RESOURCE_CONCEPT_TYPES_KEY, screenKey: "school-concepts-form", screenType: "school-form", isDefault: false, sortOrder: 4},

		{id: "b4500000-0000-0000-0000-000000000050", resourceID: L4_RESOURCE_SYSTEM_SETTINGS_ID, resourceKey: L4_RESOURCE_SYSTEM_SETTINGS_KEY, screenKey: "system-settings", screenType: "settings", isDefault: true, sortOrder: 1},

		// =============================================================
		// academic → units / memberships / subjects
		// =============================================================
		{id: "b4500000-0000-0000-0000-000000000060", resourceID: L4_RESOURCE_UNITS_ID, resourceKey: L4_RESOURCE_UNITS_KEY, screenKey: "units-list", screenType: "list", isDefault: true, sortOrder: 1},
		{id: "b4500000-0000-0000-0000-000000000061", resourceID: L4_RESOURCE_UNITS_ID, resourceKey: L4_RESOURCE_UNITS_KEY, screenKey: "units-form", screenType: "form", isDefault: false, sortOrder: 2},

		{id: "b4500000-0000-0000-0000-000000000065", resourceID: L4_RESOURCE_MEMBERSHIPS_ID, resourceKey: L4_RESOURCE_MEMBERSHIPS_KEY, screenKey: "memberships-list", screenType: "list", isDefault: true, sortOrder: 1},
		// memberships-form queda como pantalla de SOLO EDICIÓN (editar una membresía
		// existente desde la acción "editar" de la lista). La creación directa se
		// retiró: sin FAB de crear en la lista (actions_removed:["create"]), sin POST
		// en el backend y membership-add eliminado.
		{id: "b4500000-0000-0000-0000-000000000066", resourceID: L4_RESOURCE_MEMBERSHIPS_ID, resourceKey: L4_RESOURCE_MEMBERSHIPS_KEY, screenKey: "memberships-form", screenType: "form", isDefault: false, sortOrder: 2},
		{id: "b4500000-0000-0000-0000-000000000067", resourceID: L4_RESOURCE_MEMBERSHIPS_ID, resourceKey: L4_RESOURCE_MEMBERSHIPS_KEY, screenKey: "unit-directory", screenType: "directory", isDefault: false, sortOrder: 3},

		// academic → my_memberships (plan 006, N1.C ETAPA 2): default del
		// item de menú "Mis materias" del alumno. Recurso separado de
		// memberships; abre la lista readonly de materias propias.
		// Reintroducido en N1.7 F1 sobre el modelo de sesiones.
		{id: "b4500000-0000-0000-0000-000000000069", resourceID: L4_RESOURCE_MY_MEMBERSHIPS_ID, resourceKey: L4_RESOURCE_MY_MEMBERSHIPS_KEY, screenKey: "my-memberships-list", screenType: "list", isDefault: true, sortOrder: 1},

		// academic → my_grades (N3 F4, consulta de notas): default del item de
		// menú "Mis notas" del alumno. Recurso separado de grades (CRUD docente);
		// abre la lista readonly de notas propias. Espejo de my_memberships.
		{id: "b4500000-0000-0000-0000-000000000068", resourceID: L4_RESOURCE_MY_GRADES_ID, resourceKey: L4_RESOURCE_MY_GRADES_KEY, screenKey: "my-grades-list", screenType: "list", isDefault: true, sortOrder: 1},

		// academic → my_teaching (plan 027 F3): default del item de menú "Mis
		// Materias" del profesor. Recurso separado de subjects/subject_offerings;
		// abre la lista readonly de sesiones que dicta. Espejo de my_grades.
		{id: "b4500000-0000-0000-0000-00000000006e", resourceID: L4_RESOURCE_MY_TEACHING_ID, resourceKey: L4_RESOURCE_MY_TEACHING_KEY, screenKey: "my-teaching-list", screenType: "list", isDefault: true, sortOrder: 1},

		// academic → my_attendance (plan 027 F2): default del item de menú "Mi
		// Asistencia" del alumno. Recurso separado de attendance (CRUD docente);
		// abre la lista readonly de su propia asistencia. Espejo de my_grades.
		{id: "b4500000-0000-0000-0000-00000000006f", resourceID: L4_RESOURCE_MY_ATTENDANCE_ID, resourceKey: L4_RESOURCE_MY_ATTENDANCE_KEY, screenKey: "my-attendance-list", screenType: "list", isDefault: true, sortOrder: 1},

		// academic → subject_offerings (sesiones de materia, plan 010 / ADR
		// 0009). batch-enroll = "inscripción por lote", pantalla NATIVA del FE
		// (Compose, NO SDUI). Único mapping del recurso (list, is_default=true):
		// el menú expone el screen_key `batch-enroll` y el FE lo intercepta. El
		// recurso subject_offerings es IsMenuVisible=false (contenedor de la
		// sesión); el ítem se alcanza vía el flujo de sesiones. N1.7 F1.
		{id: "b4500000-0000-0000-0000-00000000006a", resourceID: L4_RESOURCE_SUBJECT_OFFERINGS_ID, resourceKey: L4_RESOURCE_SUBJECT_OFFERINGS_KEY, screenKey: "batch-enroll", screenType: "list", isDefault: true, sortOrder: 1},
		// N1.7 F2: enroll-one (inscripción individual, pantalla NATIVA) y
		// sessions-by-subject-list (lista hija SDUI; F2.2 la reubica como pestaña
		// "Sesiones" embebida en subjects-form vía detail_configs[]). Ambas bajo el
		// mismo recurso subject_offerings para satisfacer la FK screen_key; el
		// default sigue siendo batch-enroll. No necesitan ser visibles en menú por
		// sí solas.
		{id: "b4500000-0000-0000-0000-00000000006b", resourceID: L4_RESOURCE_SUBJECT_OFFERINGS_ID, resourceKey: L4_RESOURCE_SUBJECT_OFFERINGS_KEY, screenKey: "enroll-one", screenType: "list", isDefault: false, sortOrder: 2},
		{id: "b4500000-0000-0000-0000-00000000006c", resourceID: L4_RESOURCE_SUBJECT_OFFERINGS_ID, resourceKey: L4_RESOURCE_SUBJECT_OFFERINGS_KEY, screenKey: "sessions-by-subject-list", screenType: "list", isDefault: false, sortOrder: 3},
		// N1.7 F2.3: sessions-by-subject-form, formulario crear/editar de sesión
		// (modal del master-detail subjects-form). Bajo el mismo recurso
		// subject_offerings para satisfacer la FK screen_key; no es default ni
		// visible en menú (se alcanza desde la pestaña "Sesiones").
		{id: "b4500000-0000-0000-0000-00000000006d", resourceID: L4_RESOURCE_SUBJECT_OFFERINGS_ID, resourceKey: L4_RESOURCE_SUBJECT_OFFERINGS_KEY, screenKey: "sessions-by-subject-form", screenType: "form", isDefault: false, sortOrder: 4},

		{id: "b4500000-0000-0000-0000-000000000070", resourceID: L4_RESOURCE_SUBJECTS_ID, resourceKey: L4_RESOURCE_SUBJECTS_KEY, screenKey: "subjects-list", screenType: "list", isDefault: true, sortOrder: 1},
		{id: "b4500000-0000-0000-0000-000000000071", resourceID: L4_RESOURCE_SUBJECTS_ID, resourceKey: L4_RESOURCE_SUBJECTS_KEY, screenKey: "subjects-form", screenType: "form", isDefault: false, sortOrder: 2},

		// =============================================================
		// academic → invitations (N0.4-A, plan 005)
		// Default: invitations-list (lo que abre el item de menú
		// `invitations`). El form es sub-flujo (create-only), no default.
		// =============================================================
		{id: "b4500000-0000-0000-0000-000000000075", resourceID: L4_RESOURCE_INVITATIONS_ID, resourceKey: L4_RESOURCE_INVITATIONS_KEY, screenKey: "invitations-list", screenType: "list", isDefault: true, sortOrder: 1},
		{id: "b4500000-0000-0000-0000-000000000076", resourceID: L4_RESOURCE_INVITATIONS_ID, resourceKey: L4_RESOURCE_INVITATIONS_KEY, screenKey: "invitations-form", screenType: "form", isDefault: false, sortOrder: 2},
		{id: "b4500000-0000-0000-0000-000000000078", resourceID: L4_RESOURCE_INVITATIONS_ID, resourceKey: L4_RESOURCE_INVITATIONS_KEY, screenKey: "invitations-detail", screenType: "detail", isDefault: false, sortOrder: 3},

		// =============================================================
		// academic → join_requests (N0.4-B, plan 005)
		// Bandeja NATIVA del doble visto bueno. Único mapping (list,
		// is_default=true): el screen_key `join-requests-inbox` lo
		// resuelve el FE con una pantalla Compose nativa (NO SDUI), por
		// eso no hay screen_instance — el resolver solo necesita que el
		// menú exponga el screen_key vía resource_screens. El item de
		// menú aparece para quien tenga `academic.join_requests.read`.
		// =============================================================
		{id: "b4500000-0000-0000-0000-000000000077", resourceID: L4_RESOURCE_JOIN_REQUESTS_ID, resourceKey: L4_RESOURCE_JOIN_REQUESTS_KEY, screenKey: "join-requests-inbox", screenType: "list", isDefault: true, sortOrder: 1},

		// =============================================================
		// academic → guardian_relations
		// =============================================================
		// Poda F2 (plan 004-permisologia-mvp): retiradas las pantallas
		// guardian-requests-list y guardian-relations-form. El recurso
		// academic.guardian_relations queda sin mappings (huérfano,
		// prune-later — ver docs/plans/004-permisologia-mvp/diferido.md).

		// =============================================================
		// academic → periods / grades / attendance / schedules /
		// calendar
		// =============================================================
		{id: "b4500000-0000-0000-0000-000000000080", resourceID: L4_RESOURCE_PERIODS_ID, resourceKey: L4_RESOURCE_PERIODS_KEY, screenKey: "periods-list", screenType: "list", isDefault: true, sortOrder: 1},
		{id: "b4500000-0000-0000-0000-000000000081", resourceID: L4_RESOURCE_PERIODS_ID, resourceKey: L4_RESOURCE_PERIODS_KEY, screenKey: "periods-form", screenType: "form", isDefault: false, sortOrder: 2},

		{id: "b4500000-0000-0000-0000-000000000085", resourceID: L4_RESOURCE_GRADES_ID, resourceKey: L4_RESOURCE_GRADES_KEY, screenKey: "grades-list", screenType: "list", isDefault: true, sortOrder: 1},
		// grades-form ELIMINADA (2026-06-09): mapping (screen_type "form", UUID
		// …0086) retirado junto con su screen_instance — form SDUI legacy
		// reemplazado por nativas (my-grade-detail / grades-batch). UUID …0086
		// queda libre para reuso futuro.
		// grades-subject-summary (N3 F4, consulta de notas): resumen de notas por
		// sesión (vista docente). Espejo del mapping attendance-summary; screen_type
		// distinto ("summary") para satisfacer UNIQUE (resource_id, screen_type). No
		// es default (se alcanza desde la card de sesión vía event view-grades-summary).
		{id: "b4500000-0000-0000-0000-000000000087", resourceID: L4_RESOURCE_GRADES_ID, resourceKey: L4_RESOURCE_GRADES_KEY, screenKey: "grades-subject-summary", screenType: "summary", isDefault: false, sortOrder: 3},

		{id: "b4500000-0000-0000-0000-000000000090", resourceID: L4_RESOURCE_ATTENDANCE_ID, resourceKey: L4_RESOURCE_ATTENDANCE_KEY, screenKey: "attendance-list", screenType: "list", isDefault: true, sortOrder: 1},
		// attendance-batch ocupa screen_type=form para el recurso attendance.
		{id: "b4500000-0000-0000-0000-000000000091", resourceID: L4_RESOURCE_ATTENDANCE_ID, resourceKey: L4_RESOURCE_ATTENDANCE_KEY, screenKey: "attendance-batch", screenType: "form", isDefault: false, sortOrder: 2},
		{id: "b4500000-0000-0000-0000-000000000092", resourceID: L4_RESOURCE_ATTENDANCE_ID, resourceKey: L4_RESOURCE_ATTENDANCE_KEY, screenKey: "attendance-summary", screenType: "summary", isDefault: false, sortOrder: 3},

		// Poda F2 (plan 004-permisologia-mvp): retiradas las pantallas
		// schedules-list/form (95,96) y calendar-list/form (a0,a1). Los
		// recursos academic.schedules y academic.calendar quedan sin
		// mappings (huérfanos, prune-later).

		// =============================================================
		// content → assessments (configuración + revisión por docente)
		// =============================================================
		// El docente entra al CRUD vía screen_type=list, que apunta a
		// `assessments-management-list` (master-detail con questions +
		// assignment). La pantalla `assessments-list` (flujo student-take)
		// vive SOLO bajo el recurso `assessments_student` — no se mapea
		// acá para evitar que el menú docente caiga en el flow de tomar
		// el examen al hacer tap en una fila.
		{id: "b4500000-0000-0000-0000-0000000000b1", resourceID: L4_RESOURCE_ASSESSMENTS_ID, resourceKey: L4_RESOURCE_ASSESSMENTS_KEY, screenKey: "assessments-form", screenType: "form", isDefault: false, sortOrder: 2},
		{id: "b4500000-0000-0000-0000-0000000000b2", resourceID: L4_RESOURCE_ASSESSMENTS_ID, resourceKey: L4_RESOURCE_ASSESSMENTS_KEY, screenKey: "assessments-management-list", screenType: "list", isDefault: true, sortOrder: 1},
		{id: "b4500000-0000-0000-0000-0000000000b3", resourceID: L4_RESOURCE_ASSESSMENTS_ID, resourceKey: L4_RESOURCE_ASSESSMENTS_KEY, screenKey: "assessment-take", screenType: "detail", isDefault: false, sortOrder: 4},
		{id: "b4500000-0000-0000-0000-0000000000b4", resourceID: L4_RESOURCE_ASSESSMENTS_ID, resourceKey: L4_RESOURCE_ASSESSMENTS_KEY, screenKey: "assessment-questions-list", screenType: "questions", isDefault: false, sortOrder: 5},
		{id: "b4500000-0000-0000-0000-0000000000b5", resourceID: L4_RESOURCE_ASSESSMENTS_ID, resourceKey: L4_RESOURCE_ASSESSMENTS_KEY, screenKey: "assessment-question-form", screenType: "question-form", isDefault: false, sortOrder: 6},
		{id: "b4500000-0000-0000-0000-0000000000b6", resourceID: L4_RESOURCE_ASSESSMENTS_ID, resourceKey: L4_RESOURCE_ASSESSMENTS_KEY, screenKey: "assessment-result", screenType: "result", isDefault: false, sortOrder: 7},
		// b7/assessment-assignment eliminado: reemplazado por modal nativo.
		// b8/assessment-modality eliminado en plan 015 (concepto muerto).
		{id: "b4500000-0000-0000-0000-0000000000b9", resourceID: L4_RESOURCE_ASSESSMENTS_ID, resourceKey: L4_RESOURCE_ASSESSMENTS_KEY, screenKey: "assessment-review-dashboard", screenType: "review-dashboard", isDefault: false, sortOrder: 10},
		{id: "b4500000-0000-0000-0000-0000000000ba", resourceID: L4_RESOURCE_ASSESSMENTS_ID, resourceKey: L4_RESOURCE_ASSESSMENTS_KEY, screenKey: "attempt-review-detail", screenType: "attempt-review", isDefault: false, sortOrder: 11},

		// =============================================================
		// content → assessments_student (flujo del estudiante)
		// Comparte screen_key `assessments-list` con `assessments`
		// (mismo screen_instance, distinto recurso). Resolución por
		// permisos: assessments_student:read.
		// =============================================================
		{id: "b4500000-0000-0000-0000-0000000000c0", resourceID: L4_RESOURCE_ASSESSMENTS_STUDENT_ID, resourceKey: L4_RESOURCE_ASSESSMENTS_STUDENT_KEY, screenKey: "assessments-list", screenType: "list", isDefault: true, sortOrder: 1},
		{id: "b4500000-0000-0000-0000-0000000000c1", resourceID: L4_RESOURCE_ASSESSMENTS_STUDENT_ID, resourceKey: L4_RESOURCE_ASSESSMENTS_STUDENT_KEY, screenKey: "assigned-assessments-list", screenType: "assigned-list", isDefault: false, sortOrder: 2},

		// =============================================================
		// reports → stats (report-card legacy)
		// =============================================================
		// Poda F2 (plan 004-permisologia-mvp): retiradas las pantallas
		// detalle progress-detail (d1), stats-detail (d6) y report-card
		// (da). Eliminado (2026-06-15): el mapping progress → progress-dashboard
		// (d0) junto con el recurso `progress` (apuntaba a un endpoint
		// inexistente). Se conserva el dashboard stats-dashboard. El recurso
		// `reports` raíz queda sin mapping (huérfano, prune-later — ver
		// docs/plans/004-permisologia-mvp/diferido.md).
		{id: "b4500000-0000-0000-0000-0000000000d5", resourceID: L4_RESOURCE_STATS_ID, resourceKey: L4_RESOURCE_STATS_KEY, screenKey: "stats-dashboard", screenType: "dashboard", isDefault: true, sortOrder: 1},

		// notifications: mapping retirado en B7-fix junto con la
		// screen_instance notifications-list (FE no implementa
		// NotificationsListContract.kt aún).

		// =============================================================
		// messaging (plan 025 F5 — WhatsApp del staff hacia familias)
		// Único mapping (list, is_default=true): el screen_key
		// `messaging` lo resuelve el FE con una pantalla Compose nativa
		// (NO SDUI). is_default=true hace que el item de menú abra esta
		// pantalla. Aparece para quien tenga `messaging.view` (cubierto
		// por el wildcard `messaging.*` de school_admin/teacher). Sufijo
		// …e0, espejo de la screen_instance L4_SCREEN_INST_MESSAGING_ID.
		// =============================================================
		{id: "b4500000-0000-0000-0000-0000000000e0", resourceID: L4_RESOURCE_MESSAGING_ID, resourceKey: L4_RESOURCE_MESSAGING_KEY, screenKey: "messaging", screenType: "list", isDefault: true, sortOrder: 1},
	}
}
