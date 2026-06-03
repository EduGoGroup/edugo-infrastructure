package l4

import (
	"encoding/json"
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ApplyScreenInstances siembra las screen_instances restantes del
// sistema (70 filas) en `ui_config.screen_instances`. Excluye lo ya
// sembrado por capas anteriores:
//   - L0: `announcements-list`
//   - L2: `announcement-form`
//   - L3: `materials-list`, `material-form`
//
// Cobertura (F6-REQ-2.1, F6-REQ-3.1):
//   - 61 instances derivadas del inventario legacy
//     `[archivado pre-Fase-6] data.go::screenInstanceSeedRows[]`
//     (lineas 615-691). Las definiciones se redefinen con criterio
//     (NO copy-paste — F6-REQ-2.5).
//   - 9 instances NUEVAS exigidas por el FE pero ausentes del legacy
//     (parte del set de 14 `screen_key_phantom` reportadas por el
//     cross-checker baseline):
//   - assessment-assignment, assessment-modality,
//     assessment-review-dashboard, assigned-assessments-list,
//     attempt-review-detail (5 pantallas de assessments).
//   - membership-add, user-roles (admin).
//   - school-concepts-list, school-concepts-form (concept_types
//     a nivel escuela).
//   - Las 5 phantom restantes (app-login, app-settings,
//     attendance-form, guardian-relations-list, guardian_relations-form)
//     SI existian en legacy y figuran como phantom solo porque
//     Layer_Legacy esta desactivado desde Fase 2 (ADR-6); B4 las
//     siembra como parte del set de 61.
//
// Descartes vs legacy (12 keys), todos documentados:
//   - `announcements-list`, `materials-list` → ya en L0/L3.
//   - `announcements-form` (linea 686) → duplicado plural del key L2
//     `announcement-form`. La forma canonica es singular (FE usa
//     `announcement-form` en L2 y `announcements-form` no aparece en
//     ningun contrato KMP).
//   - `material-detail`, `child-progress`, `children-list` → reportados
//     como `screen_key_dead` por el cross-checker (FE no implementa
//     composable). Aceptados como dead — F6-REQ-3 no exige seedar
//     dead-screens.
//   - `announcement-detail` → DEAD (no hay contrato KMP). Coherente
//     con la decision de tratar material-detail como dead.
//   - `calendar-events-list`, `calendar-event-detail`,
//     `calendar-event-form` → duplicados semanticos del par
//     canonico `calendar-list` / `calendar-form` (ambos en legacy con
//     contratos KMP). Descartados los `-events` / `-event-` por no
//     tener contrato KMP.
//   - `schedule-form`, `schedule-detail` → duplicados singulares de
//     `schedules-form` (canonico, con contrato KMP). Descartados.
//
// api_prefix por servicio (estable, alineado al routing real del
// backend; coherente con L0 platform / L3 academic):
//   - `platform`  → announcements (L0), screens config, system-settings,
//     notifications.
//   - `identity`  → users, roles, permissions, audit, screen_templates,
//     screen_instances (la gestion del catalogo iam vive
//     en el servicio identity).
//   - `academic`  → units, memberships, subjects, periods, calendar,
//     schedules, guardian_relations, materials (L3),
//     directorios escolares, concept_types (school).
//   - `learning`  → assessments, grades, attendance, progress, stats,
//     report-card, dashboards (KPIs).
//
// Correcciones aplicadas vs los 9 `service_prefix_mismatch` del
// baseline (F6-REQ-3.2):
//   - `audit`              → seed declara `identity`. La canónica
//     del cross-checker pide `iam:` pero el routing HTTP real expone
//     auditoria bajo el servicio identity. Decision: SEED CORRECTO,
//     la tabla canonica del cross-checker queda desactualizada (no se
//     toca el FE, no se toca el seed; el warning del checker se
//     acepta como "tabla canonica TODO").
//   - `roles`, `users`     → idem que `audit` (identity).
//   - `screens`            → `platform`. Coincide con el FE actual.
//     Mismo razonamiento: tabla canonica obsoleta.
//   - `guardian_relations` → seed declara `academic` (no `learning`
//     como hace el FE). Esto CORRIGE el seed; cualquier slot_data
//     legacy que usaba `learning` se reescribe a `academic`. El FE
//     debera alinearse — queda como ticket FE.
//   - los 3 `info` (concept_types, permissions_mgmt, reports) son
//     simplemente recursos no mapeados en la tabla canonica del
//     cross-checker. Decision similar: seed correcto, tabla TODO.
//
// Idempotencia: UPSERT con conflict target `screen_key` (UNIQUE en el
// schema, ver entities/screen_instance.go). DoNothing para evitar
// pisar customizaciones manuales en environments live; B0..B6 deben
// poder re-correrse sin efectos.
func ApplyScreenInstances(tx *gorm.DB) error {
	instances := buildL4ScreenInstances()

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "screen_key"}},
		DoNothing: true,
	}).CreateInBatches(&instances, 50).Error; err != nil {
		return fmt.Errorf("ApplyScreenInstances: upsert screen_instances: %w", err)
	}
	return nil
}

// l4ScreenInstanceRow describe una fila lista para convertir a
// entities.ScreenInstance. Concentrar los campos primitivos aqui
// facilita revisar la tabla completa sin ruido de tipos
// (uuid.UUID, *string, json.RawMessage).
//
//   - templateID: literal UUID; debe ser uno de los 3 templates base
//     L0 (L0_SCREEN_TPL_*_ID_REF) o uno de los 5 L4
//     (l4TplLoginV1ID, etc.).
//   - slotData: JSON literal RAW. Se inyecta sin re-serializar para
//     preservar formato (mismo patron que announcementsListSlotData
//     en seeds/system/layers/l0_screens.go).
//   - requiredPermission/handlerKey: opcional. "" se traduce a NULL
//     en el upsert.
//   - scope: "system" | "school" | "unit". Refleja el alcance del
//     dato, NO del rol que la consume.
type l4ScreenInstanceRow struct {
	id                 string
	screenKey          string
	templateID         string
	name               string
	description        string
	slotData           string
	scope              string
	requiredPermission string
	handlerKey         string
}

func (r l4ScreenInstanceRow) toEntity() entities.ScreenInstance {
	desc := r.description
	var descPtr *string
	if desc != "" {
		descPtr = &desc
	}
	rp := r.requiredPermission
	var rpPtr *string
	if rp != "" {
		rpPtr = &rp
	}
	hk := r.handlerKey
	var hkPtr *string
	if hk != "" {
		hkPtr = &hk
	}
	return entities.ScreenInstance{
		ID:                 mustParseL4UUID(r.id, "screen_instance:"+r.screenKey),
		ScreenKey:          r.screenKey,
		TemplateID:         mustParseL4UUID(r.templateID, "screen_instance.template_id:"+r.screenKey),
		Name:               r.name,
		Description:        descPtr,
		SlotData:           json.RawMessage([]byte(r.slotData)),
		Scope:              r.scope,
		RequiredPermission: rpPtr,
		HandlerKey:         hkPtr,
		IsActive:           true,
	}
}

func buildL4ScreenInstances() []entities.ScreenInstance {
	rows := l4ScreenInstanceRows()
	out := make([]entities.ScreenInstance, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.toEntity())
	}
	return out
}

// l4ScreenInstanceRows retorna la tabla literal de las 70 screen
// instances sembradas por L4. Agrupadas por dominio y comentadas con
// la decision aplicada al refactorizar el slot_data legacy.
//
// Convencion del slot_data:
//   - "title": label en el topBar (consumido por list/form/detail
//     templates v1).
//   - "actions": botones del topBar/row; cada accion declara
//     {id, label, icon, permission, [scope, event]}.
//   - "columns" (lists) o "fields" (forms): definicion mínima del
//     contenido. NO se redefine la geometria UI (eso vive en
//     screen_templates.go).
//   - "api_prefix": servicio HTTP donde el FE hace la llamada de
//     datos. Reemplaza el `data_endpoint` literal del legacy
//     (`learning:/api/v1/...`); el FE arma el path por convencion
//     `{prefix}:/api/v1/{resource}`.
func l4ScreenInstanceRows() []l4ScreenInstanceRow {
	return []l4ScreenInstanceRow{
		// ===========================================================
		// AUTH & SHELL
		// ===========================================================
		// app-login: handler_key="login" → el FE delega la accion de
		// auth al handler nativo, no a la maquina SDUI.
		{
			id:          L4_SCREEN_INST_APP_LOGIN_ID,
			screenKey:   "app-login",
			templateID:  l4TplLoginV1ID,
			name:        "Inicio de Sesión",
			description: "Pantalla de login con autenticacion local y social",
			slotData: `{
  "app_logo": "edugo_logo",
  "app_name": "EduGo",
  "app_tagline": "Aprender es facil",
  "email_label": "Email",
  "password_label": "Contraseña",
  "remember_label": "Recordarme",
  "login_btn_label": "Ingresar",
  "forgot_password_label": "Olvidaste tu contraseña?",
  "divider_text": "o continuar con",
  "google_btn_label": "Google",
  "api_prefix": "identity"
}`,
			scope:      "system",
			handlerKey: "login",
		},
		// app-settings: settings del usuario; logout via accion.
		{
			id:          L4_SCREEN_INST_APP_SETTINGS_ID,
			screenKey:   "app-settings",
			templateID:  l4TplSettingsV1ID,
			name:        "Configuración",
			description: "Configuracion de cuenta y preferencias del usuario",
			slotData: `{
  "title": "Configuración",
  "appearance_title": "Apariencia",
  "dark_mode_label": "Modo oscuro",
  "theme_label": "Tema",
  "notifications_title": "Notificaciones",
  "push_label": "Notificaciones push",
  "email_label": "Notificaciones por email",
  "logout_label": "Cerrar sesión",
  "api_prefix": "identity"
}`,
			scope: "system",
		},
		// dashboard-home: pantalla shell de routing. El KMP la mapea a
		// DynamicDashboardScreen, que delega al dashboard especifico del
		// rol activo (dashboard-teacher/student/superadmin/...). NO se
		// mapea en resource_screens — el FE la resuelve por screen_key
		// (mismo patron que app-login / app-settings). El cross-checker
		// la reportaba como phantom solo porque el seed previo no la
		// declaraba; se siembra aqui como shell sin mapping.
		{
			id:          L4_SCREEN_INST_DASHBOARD_HOME_ID,
			screenKey:   "dashboard-home",
			templateID:  l4TplDashboardV1ID,
			name:        "Inicio",
			description: "Pantalla shell de inicio que delega al dashboard del rol activo",
			slotData: `{
  "title": "Inicio",
  "api_prefix": "learning"
}`,
			scope: "system",
		},

		// ===========================================================
		// DASHBOARDS POR ROL (5 roles implementados + 2 dashboards
		// agregados: progress / stats)
		// ===========================================================
		// Todos los dashboards usan dashboard-basic-v1 (template L4).
		// api_prefix=learning porque los KPIs vienen del servicio
		// learning (stats/progress endpoint).
		{
			id:          L4_SCREEN_INST_DASH_TEACHER_ID,
			screenKey:   "dashboard-teacher",
			templateID:  l4TplDashboardV1ID,
			name:        "Dashboard Profesor",
			description: "Panel principal del profesor",
			slotData: `{
  "title": "Inicio",
  "greeting_text": "Buenos días",
  "kpi_students_label": "Estudiantes",
  "kpi_materials_label": "Materiales",
  "kpi_avg_score_label": "Nota promedio",
  "kpi_completion_label": "Avance",
  "activity_title": "Actividad reciente",
  "upload_label": "Subir material",
  "progress_label": "Ver progreso",
  "api_prefix": "learning"
}`,
			scope: "school",
		},
		{
			id:          L4_SCREEN_INST_DASH_STUDENT_ID,
			screenKey:   "dashboard-student",
			templateID:  l4TplDashboardV1ID,
			name:        "Dashboard Estudiante",
			description: "Panel principal del estudiante",
			slotData: `{
  "title": "Inicio",
  "greeting_text": "Hola",
  "kpi_students_label": "Cursos",
  "kpi_materials_label": "Materiales",
  "kpi_avg_score_label": "Mi nota",
  "kpi_completion_label": "Mi progreso",
  "activity_title": "Actividad reciente",
  "upload_label": "Mis materiales",
  "progress_label": "Mi progreso",
  "api_prefix": "learning"
}`,
			scope: "unit",
		},
		{
			id:          L4_SCREEN_INST_DASH_SUPERADMIN_ID,
			screenKey:   "dashboard-superadmin",
			templateID:  l4TplDashboardV1ID,
			name:        "Dashboard Superadmin",
			description: "Panel principal del superadministrador",
			slotData: `{
  "title": "Inicio",
  "greeting_text": "Hola, admin",
  "kpi_students_label": "Usuarios",
  "kpi_materials_label": "Escuelas",
  "kpi_avg_score_label": "Roles",
  "kpi_completion_label": "Permisos",
  "activity_title": "Actividad reciente",
  "upload_label": "Crear escuela",
  "progress_label": "Ver estadísticas",
  "api_prefix": "learning"
}`,
			scope: "system",
		},
		{
			id:          L4_SCREEN_INST_DASH_SCHOOLADM_ID,
			screenKey:   "dashboard-schooladmin",
			templateID:  l4TplDashboardV1ID,
			name:        "Dashboard Administrador Escolar",
			description: "Panel principal del administrador de la escuela",
			slotData: `{
  "title": "Inicio",
  "greeting_text": "Hola, equipo",
  "kpi_students_label": "Estudiantes",
  "kpi_materials_label": "Materiales",
  "kpi_avg_score_label": "Nota promedio",
  "kpi_completion_label": "Avance",
  "activity_title": "Actividad reciente",
  "upload_label": "Nueva clase",
  "progress_label": "Reporte semanal",
  "api_prefix": "learning"
}`,
			scope: "school",
		},
		{
			id:          L4_SCREEN_INST_DASH_GUARDIAN_ID,
			screenKey:   "dashboard-guardian",
			templateID:  l4TplDashboardV1ID,
			name:        "Dashboard Padres/Tutores",
			description: "Panel principal del guardian",
			slotData: `{
  "title": "Inicio",
  "greeting_text": "Hola",
  "kpi_students_label": "Hijos",
  "kpi_materials_label": "Materiales",
  "kpi_avg_score_label": "Promedio",
  "kpi_completion_label": "Asistencia",
  "activity_title": "Últimas novedades",
  "upload_label": "Vincular hijo",
  "progress_label": "Progreso",
  "api_prefix": "learning"
}`,
			scope: "school",
		},
		// progress-dashboard / stats-dashboard: reutilizan
		// dashboard-basic-v1. Conservados del legacy porque el FE
		// (ProgressDashboardContract.kt, StatsDashboardContract.kt)
		// declara estos screenKeys explicitamente.
		{
			id:          L4_SCREEN_INST_PROGRESS_DASH_ID,
			screenKey:   "progress-dashboard",
			templateID:  l4TplDashboardV1ID,
			name:        "Progreso Académico",
			description: "Dashboard de progreso académico",
			slotData: `{
  "title": "Progreso",
  "greeting_text": "Progreso general",
  "kpi_students_label": "Estudiantes",
  "kpi_materials_label": "Materiales",
  "kpi_avg_score_label": "Nota media",
  "kpi_completion_label": "% Avance",
  "activity_title": "Hitos recientes",
  "api_prefix": "learning"
}`,
			scope: "unit",
		},
		{
			id:          L4_SCREEN_INST_STATS_DASH_ID,
			screenKey:   "stats-dashboard",
			templateID:  l4TplDashboardV1ID,
			name:        "Estadísticas",
			description: "Dashboard de estadísticas del sistema",
			slotData: `{
  "title": "Estadísticas",
  "greeting_text": "Resumen",
  "kpi_students_label": "Estudiantes",
  "kpi_materials_label": "Materiales",
  "kpi_avg_score_label": "Calif.",
  "kpi_completion_label": "Asistencia",
  "activity_title": "Tendencias",
  "api_prefix": "learning"
}`,
			scope: "school",
		},

		// ===========================================================
		// ADMIN: USERS / SCHOOLS / ROLES / PERMISSIONS
		// ===========================================================
		// api_prefix=identity. La canónica del cross-checker pide
		// "iam:" pero el routing HTTP real expone el catálogo bajo
		// identity (ver doc principal de ApplyScreenInstances).
		usersList(),
		usersForm(),
		schoolsList(),
		schoolsForm(),
		// Poda menú (2026-05-29): rolesList/rolesForm/permissionsList/permissionsForm
		// eliminadas — el FE KMP no implementa esas pantallas y los recursos
		// `roles`/`permissions_mgmt` fueron retirados del menú.

		// Poda menú (2026-06-01): screenTplList/screenInstList/screenInstForm/
		// screensForm eliminadas — las pantallas de configuración SDUI se
		// reimplementaron en el admin-tool de Go y los recursos
		// `screen_templates`/`screen_instances` se retiraron del menú.

		// ===========================================================
		// ADMIN: SYSTEM SETTINGS + CONCEPT TYPES
		// ===========================================================
		systemSettings(),
		conceptTypesList(),
		conceptTypesForm(),

		// ===========================================================
		// ADMIN: AUDIT
		// ===========================================================
		// audit-events-list: list; audit-detail: detail.
		// api_prefix=identity (auditoria vive bajo identity en el
		// routing real).
		auditEventsList(),
		auditDetail(),

		// ===========================================================
		// ACADEMIC: UNITS / MEMBERSHIPS / SUBJECTS / PERIODS
		// ===========================================================
		unitsList(),
		unitsForm(),
		membershipsList(),
		membershipsForm(),   // SOLO EDICIÓN (sin save_new): editar una membresía existente; la creación directa se retiró
		myMembershipsList(), // "Mis materias" (alumno), reintroducida en N1.7 F1 sobre sesiones
		subjectsList(),
		subjectsForm(),
		periodsList(),
		periodsForm(),
		invitationsList(),
		invitationsForm(),
		joinRequestsInbox(),
		subjectOfferingsBatchEnroll(), // inscripción por lote (pantalla nativa), N1.7 F1
		enrollOne(),                   // inscripción individual (pantalla nativa), N1.7 F2
		sessionsBySubjectList(),       // sesiones por materia (lista hija SDUI), N1.7 F2
		sessionsBySubjectForm(),       // crear/editar sesión de materia (modal SDUI), N1.7 F2.3

		// ===========================================================
		// ACADEMIC: GUARDIAN / CALENDAR / SCHEDULES
		// ===========================================================
		// Poda F2 (plan 004-permisologia-mvp): guardian-relations-list,
		// guardian-relations-form, guardian_relations-form (alias),
		// guardian-requests-list, calendar-list, calendar-form,
		// schedules-list y schedules-form se retiraron del MVP. Sus
		// constructores y constantes también se eliminaron. El recurso
		// academic.guardian_relations queda huérfano (prune-later, ver
		// docs/plans/004-permisologia-mvp/diferido.md); no se toca en
		// esta pasada por riesgo del proxy de dashboards.

		// ===========================================================
		// ACADEMIC: GRADES / ATTENDANCE
		// ===========================================================
		gradesList(),
		gradesForm(),
		attendanceList(),
		attendanceBatch(),
		attendanceForm(),
		attendanceSummary(),

		// ===========================================================
		// CONTENT: ASSESSMENTS (gestion + estudiante + nuevas)
		// ===========================================================
		// Las 5 instancias "phantom-nuevas" (assessment-assignment,
		// assessment-modality, assessment-review-dashboard,
		// assigned-assessments-list, attempt-review-detail) son
		// completamente nuevas en B4 — el legacy no las declara y
		// el FE las exige. F6-REQ-3.1.
		assessmentsList(),
		assessmentsForm(),
		assessmentsMgmtList(),
		assessmentTake(),
		assessmentResult(),
		assessmentQuestionsList(),
		assessmentQuestionForm(),
		assessmentAssignment(),      // phantom-nueva
		assessmentModality(),        // phantom-nueva
		assessmentReviewDashboard(), // phantom-nueva
		assignedAssessmentsList(),   // phantom-nueva
		attemptReviewDetail(),       // phantom-nueva

		// ===========================================================
		// REPORTS (detalles + report-card)
		// ===========================================================
		// Poda F2 (plan 004-permisologia-mvp): progress-detail,
		// stats-detail y report-card se retiraron del MVP (los
		// dashboards progress-dashboard / stats-dashboard SÍ se
		// conservan). Sus constructores y constantes se eliminaron.

		// ===========================================================
		// DIRECTORIES & MISC (unit-directory)
		// ===========================================================
		// notifications-list retirado en B7-fix: el FE no implementa
		// ningún Contract.kt para esta screen_key — el cross-checker
		// la reportaba como screen_key_dead. Re-sembrar cuando el FE
		// agregue NotificationsListContract.kt.
		unitDirectory(),
		// students-by-subject-list eliminada (2026-06-02): era SOLO el panel
		// detalle "Alumnos" embebido en subjects-form, retirado porque un alumno
		// se inscribe en una SESIÓN, no en la materia. Sin otras referencias
		// (no estaba en menú ni en resource_screens). Su constructor y constante
		// se eliminaron.

		// ===========================================================
		// PHANTOM-NUEVAS NO-ASSESSMENT (3 adicionales)
		// ===========================================================
		// school-concepts-list / school-concepts-form: variante del
		// CRUD de concept_types con scope=school (overrides locales).
		// El FE las trata como pantallas separadas del concept-types-*
		// de scope=system.
		// membership-add se retiró: la creación directa de membresías se
		// eliminó (redundante con invitación→solicitud→aprobación).
		// user-roles: edicion de roles asignados a un usuario.
		// Decisiones de permisos documentadas en contract-check-NOTES.md
		// (TC-A): user-roles usa users:update; no se crean permisos
		// nuevos.
		schoolConceptsList(),
		schoolConceptsForm(),
		userRoles(),

		// ===========================================================
		// FASE 3 (B7b) — DEMO CRUD DATA-DRIVEN SIN KOTLIN
		// ===========================================================
		// Poda F2 (plan 004-permisologia-mvp): colors-list / colors-form
		// (pareja demo) se retiraron del MVP. Sus constructores y
		// constantes se eliminaron. El recurso platform.colors queda
		// huérfano (prune-later); no se toca en esta pasada.
	}
}
