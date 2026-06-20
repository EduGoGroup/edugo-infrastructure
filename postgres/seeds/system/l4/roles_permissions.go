package l4

import (
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ApplyRolesPermissions siembra 11 roles (5 canónicos + 6 alias)
// implementados en KMP — student, teacher, guardian, school_admin más
// los 6 alias school_director, school_coordinator, school_assistant,
// assistant_teacher, observer, readonly_auditor — más los permisos del
// catálogo, y los grants en iam.role_grants con patterns wildcard.
//
// P4-1 (plan B): la asignación 1:1 rol×permiso en iam.role_permissions
// fue eliminada. Los permisos efectivos por rol se resuelven
// exclusivamente vía iam.role_grants (allow + deny con patterns) y
// iam.user_grants (overrides por usuario). Cualquier UI de gestión de
// permisos se reconstruye en P4-2 sobre el modelo nuevo.
//
// Composición del bloque B2:
//   - 11 roles            → applyL4Roles (5 canónicos + 6 alias)
//   - 80 permisos NUEVOS  → applyL4Permissions (excluye `announcements:*`
//     ya en L0 y `materials:read/create/update`
//     ya en L3) — el catálogo de permisos sigue
//     existiendo como referencia semántica, pero
//     las asignaciones a roles se hacen vía
//     patterns en iam.role_grants.
//   - role_grants         → applyL4RoleGrants (patterns wildcard-first;
//     readonly_auditor combina allow amplio +
//     deny mutativos).
//
// Refs: phase-6-layer-l4/{requirements,design,tasks}.md (B2).
func ApplyRolesPermissions(tx *gorm.DB) error {
	if err := applyL4Roles(tx); err != nil {
		return fmt.Errorf("ApplyRolesPermissions: roles: %w", err)
	}
	if err := applyL4Permissions(tx); err != nil {
		return fmt.Errorf("ApplyRolesPermissions: permissions: %w", err)
	}
	if err := applyL4RoleGrants(tx); err != nil {
		return fmt.Errorf("applyL4RoleGrants: %w", err)
	}
	return nil
}

// -----------------------------------------------------------------------
// Roles
// -----------------------------------------------------------------------

// l4RoleSpec describe una fila lista para convertir a entities.Role.
type l4RoleSpec struct {
	idStr       string
	name        string
	displayName string
	description string
	scope       string
	// parentIDStr es el rol canónico del que hereda grants (ADR-6).
	// Vacío para roles canónicos (parent_role_id NULL). Los alias
	// apuntan aquí a su canónico y dejan de declarar grants propios:
	// la herencia se resuelve y aplana en el login.
	parentIDStr string
	// landingScreenKey es el screen_key del dashboard de inicio de este rol
	// (ADR 0024 F0 DEC-2). Vacío → NULL → cae al default de la escuela o al
	// fallback de sistema. Los 5 canónicos y los 6 alias lo declaran
	// explícitamente (el landing NO se hereda vía parent_role_id: la cascada
	// del backend mira solo el campo propio del rol).
	landingScreenKey string
}

// l4RoleSpecs retorna las specs declarativas de los 11 roles que L4
// siembra (5 canónicos + 6 alias). Helper compartido por applyL4Roles
// y por el accessor público l4.Roles() — la lógica de construcción del
// slice de entities vive una sola vez (buildL4Roles).
func l4RoleSpecs() []l4RoleSpec {
	return []l4RoleSpec{
		{
			idStr:            L4_ROLE_STUDENT_ID,
			name:             L4_ROLE_STUDENT_NAME,
			displayName:      "Estudiante",
			description:      "Alumno inscrito en una unidad académica.",
			scope:            "unit",
			landingScreenKey: "dashboard-student",
		},
		{
			idStr:            L4_ROLE_TEACHER_ID,
			name:             L4_ROLE_TEACHER_NAME,
			displayName:      "Profesor",
			description:      "Docente con permisos de gestión de clase (asistencia, calificaciones, evaluaciones, materiales).",
			scope:            "unit",
			landingScreenKey: "dashboard-teacher",
		},
		{
			idStr:            L4_ROLE_GUARDIAN_ID,
			name:             L4_ROLE_GUARDIAN_NAME,
			displayName:      "Apoderado",
			description:      "Tutor legal o apoderado vinculado a uno o más estudiantes.",
			scope:            "unit",
			landingScreenKey: "dashboard-guardian",
		},
		// PRE-4: el rol `platform_admin` (L4_ROLE_ADMIN_*) fue
		// eliminado. Sus capacidades quedan cubiertas por `super_admin`
		// (L0) que ya tiene acceso global.
		{
			idStr:            L4_ROLE_SCHOOL_ADMIN_ID,
			name:             L4_ROLE_SCHOOL_ADMIN_NAME,
			displayName:      "Administrador de Escuela",
			description:      "Administrador con control total dentro de una institución educativa.",
			scope:            "school",
			landingScreenKey: "dashboard-schooladmin",
		},
		// --- Alias roles (heredan grants del canónico) ---
		// landing_screen_key de los alias (ADR 0024 sub-deuda "herencia del
		// landing"): los 6 alias reciben EXPLÍCITAMENTE el dashboard de su rol
		// canónico. La cascada del backend (rol ?? escuela ?? "dashboard-home")
		// solo mira el campo PROPIO del rol —no resuelve la herencia de grants
		// (ADR-6) para el landing—, así que un alias con NULL caía al default de
		// la escuela, que es "dashboard-home" (el dashboard básico genérico).
		// Sembrar el landing aquí hace que coordinador/director/asistente
		// aterricen en su dashboard real en vez del home genérico.
		{
			idStr:            L4_ROLE_SCHOOL_DIRECTOR_ID,
			name:             L4_ROLE_SCHOOL_DIRECTOR_NAME,
			displayName:      "Director de Escuela",
			description:      "Director de la institución educativa. Alias de school_admin (hereda todos sus permisos).",
			scope:            "school",
			parentIDStr:      L4_ROLE_SCHOOL_ADMIN_ID,
			landingScreenKey: "dashboard-schooladmin",
		},
		{
			idStr:            L4_ROLE_SCHOOL_COORDINATOR_ID,
			name:             L4_ROLE_SCHOOL_COORDINATOR_NAME,
			displayName:      "Coordinador de Escuela",
			description:      "Coordinador académico de la institución. Alias de school_admin (hereda todos sus permisos).",
			scope:            "school",
			parentIDStr:      L4_ROLE_SCHOOL_ADMIN_ID,
			landingScreenKey: "dashboard-schooladmin",
		},
		{
			idStr:            L4_ROLE_SCHOOL_ASSISTANT_ID,
			name:             L4_ROLE_SCHOOL_ASSISTANT_NAME,
			displayName:      "Asistente de Escuela",
			description:      "Personal de apoyo administrativo de la institución. Alias de school_admin (hereda todos sus permisos).",
			scope:            "school",
			parentIDStr:      L4_ROLE_SCHOOL_ADMIN_ID,
			landingScreenKey: "dashboard-schooladmin",
		},
		{
			idStr:            L4_ROLE_ASSISTANT_TEACHER_ID,
			name:             L4_ROLE_ASSISTANT_TEACHER_NAME,
			displayName:      "Profesor Asistente",
			description:      "Docente auxiliar. Alias de teacher (hereda todos sus permisos).",
			scope:            "unit",
			parentIDStr:      L4_ROLE_TEACHER_ID,
			landingScreenKey: "dashboard-teacher",
		},
		{
			idStr:            L4_ROLE_OBSERVER_ID,
			name:             L4_ROLE_OBSERVER_NAME,
			displayName:      "Observador",
			description:      "Observador con visibilidad sobre la clase. Alias de teacher (hereda todos sus permisos).",
			scope:            "unit",
			parentIDStr:      L4_ROLE_TEACHER_ID,
			landingScreenKey: "dashboard-teacher",
		},
		{
			idStr:       L4_ROLE_READONLY_AUDITOR_ID,
			name:        L4_ROLE_READONLY_AUDITOR_NAME,
			displayName: "Auditor de Solo Lectura",
			// readonly_auditor NO hereda de ningún canónico (allow read-only
			// propio; ver nota abajo). Aterriza en dashboard-teacher: es scope
			// unit y su acceso es la vista de clase en solo lectura, el dashboard
			// más cercano a su superficie. Sin landing caería al home genérico
			// "dashboard-home" en vez del dashboard de su superficie.
			landingScreenKey: "dashboard-teacher",
			// NO hereda: su allow read-only no coincide con el de teacher
			// (teacher carece de academic.guardian_relations/memberships y
			// content.assessments_student, y a la vez aporta
			// admin.system_settings.settings y admin.users.update:own que
			// readonly nunca tuvo). Heredar de teacher cambiaría el set
			// efectivo aplanado → se mantiene standalone para preservar el
			// baseline exacto (ver reporte F1 / nota al ADR-6).
			description: "Auditor con permisos exclusivamente de lectura. Allow read-only amplio + deny de todas las acciones de mutación (create/update/delete/publish/finalize/activate/approve/grade/attempt/assign/review/manage/request).",
			scope:       "unit",
		},
	}
}

// buildL4Roles materializa las specs en entities.Role. Devuelve error
// si algún UUID está corrupto en las constantes — el panic se diferiría
// a Apply, pero el accessor también lo necesita para reportar a su caller.
func buildL4Roles() ([]entities.Role, error) {
	specs := l4RoleSpecs()
	roles := make([]entities.Role, 0, len(specs))
	for _, s := range specs {
		id, err := uuid.Parse(s.idStr)
		if err != nil {
			return nil, fmt.Errorf("parse role id %s: %w", s.idStr, err)
		}
		var parentID *uuid.UUID
		if s.parentIDStr != "" {
			pid, perr := uuid.Parse(s.parentIDStr)
			if perr != nil {
				return nil, fmt.Errorf("parse parent role id %s for %s: %w", s.parentIDStr, s.name, perr)
			}
			parentID = &pid
		}
		desc := s.description
		var landing *string
		if s.landingScreenKey != "" {
			lk := s.landingScreenKey
			landing = &lk
		}
		roles = append(roles, entities.Role{
			ID:               id,
			Name:             s.name,
			DisplayName:      s.displayName,
			Description:      &desc,
			Scope:            s.scope,
			ParentRoleID:     parentID,
			LandingScreenKey: landing,
			IsActive:         true,
			IsSystem:         true,
		})
	}
	return roles, nil
}

// applyL4Roles inserta los 11 roles (5 canónicos + 6 alias) con UUIDs
// propios L4. Patrón idéntico a applyL1Role / upsertL0Role:
// OnConflict (id) DO NOTHING.
func applyL4Roles(tx *gorm.DB) error {
	roles, err := buildL4Roles()
	if err != nil {
		return err
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).CreateInBatches(&roles, 10).Error
}

// -----------------------------------------------------------------------
// Permisos
// -----------------------------------------------------------------------

// l4PermissionSpec describe una fila de iam.permissions a sembrar.
// Mantiene las constantes legacy de UUID para minimizar reescritura;
// los `name` cumplen el CHECK del esquema (regex camelcase con `:`).
type l4PermissionSpec struct {
	idStr       string
	resourceID  string
	name        string
	displayName string
	description string
	action      string
	scope       string
}

// l4Permissions retorna el catálogo de permisos sembrados por L4.
//
// NO incluye permisos ya sembrados por capas previas:
//   - `announcements:read/create/update/delete` → L0 (layers.L0_PERM_ANNOUNCEMENTS_*)
//   - `materials:read/create/update`           → L3 (layers.L3_PERM_MATERIALS_*)
//
// SÍ incluye `materials:delete/download/publish` (3 acciones que L3
// dejó fuera por design F5-REQ-2.1) usando como FK el recurso materials
// de L3 (b3000000-0000-0000-0000-000000000001).
//
// El UUID del recurso materials de L3 difiere del UUID del legacy
// (20000000-…-30 vs b3000000-…-01); B2 usa el UUID L3 para esos 3
// permisos extra para que la FK resuelva.
func l4Permissions() []l4PermissionSpec {
	const materialsResourceID = "b3000000-0000-0000-0000-000000000001" // = layers.L3_RESOURCE_MATERIALS_ID
	return []l4PermissionSpec{
		// --- assessments (resource 20000000-…-31) ---
		{"8c8d7a5b-2688-4646-9888-bc53600dbbc0", L4_RESOURCE_ASSESSMENTS_ID, "content.assessments.attempt", "Rendir Evaluaciones", "Intentar evaluaciones como estudiante", "attempt", "unit"},
		{"0e72f5de-e2df-4bc6-b6ec-68266855d1e8", L4_RESOURCE_ASSESSMENTS_ID, "content.assessments.create", "Crear Evaluaciones", "Crear evaluaciones y exámenes", "create", "unit"},
		{"8b39f083-c630-4d1c-92ec-58afabb3376b", L4_RESOURCE_ASSESSMENTS_ID, "content.assessments.delete", "Eliminar Evaluaciones", "Eliminar evaluaciones", "delete", "unit"},
		{"d477a2fd-f996-41b5-9dd2-d66cda3460d6", L4_RESOURCE_ASSESSMENTS_ID, "content.assessments.grade", "Calificar Evaluaciones", "Calificar respuestas de estudiantes", "grade", "unit"},
		{"1b69ea9a-1e15-4c38-a5e1-68c408bb1b97", L4_RESOURCE_ASSESSMENTS_ID, "content.assessments.publish", "Publicar Evaluaciones", "Publicar evaluaciones para estudiantes", "publish", "unit"},
		{"2606efba-8615-48ef-bc2a-6bf576a9158c", L4_RESOURCE_ASSESSMENTS_ID, "content.assessments.read", "Ver Evaluaciones", "Ver evaluaciones", "read", "unit"},
		{"91c8ae21-a955-4d59-b52e-6cf74bd6532b", L4_RESOURCE_ASSESSMENTS_ID, "content.assessments.update", "Editar Evaluaciones", "Modificar evaluaciones", "update", "unit"},
		{"a5000000-0000-0000-0000-000000000001", L4_RESOURCE_ASSESSMENTS_ID, "content.assessments.assign", "Asignar Evaluaciones", "Asignar evaluaciones a estudiantes o unidades", "assign", "unit"},
		{"a6000000-0000-0000-0000-000000000001", L4_RESOURCE_ASSESSMENTS_ID, "content.assessments.review", "Revisar Evaluaciones", "Revisar intentos de evaluaciones de estudiantes", "review", "unit"},
		{"b457a385-29bc-4e5e-a79d-06c74f81d23e", L4_RESOURCE_ASSESSMENTS_ID, "content.assessments.view_results", "Ver Resultados", "Ver resultados propios", "view_results", "unit"},

		// --- assessments_student (resource 20000000-…-33) ---
		{"30000000-0000-0000-0000-000000000033", L4_RESOURCE_ASSESSMENTS_STUDENT_ID, "content.assessments_student.read", "Ver evaluaciones como estudiante", "Ver y tomar evaluaciones desde perspectiva del estudiante", "read", "unit"},

		// --- materials gap (L3 sembró read/create/update; faltan 3) ---
		{"6358de3d-11ef-49c0-be42-da51ffcdfbc1", materialsResourceID, "content.materials.delete", "Eliminar Materiales", "Eliminar materiales", "delete", "unit"},
		{"9b0c10e0-0a7b-4e73-af9a-2d5adf99790f", materialsResourceID, "content.materials.download", "Descargar Materiales", "Descargar materiales educativos", "download", "unit"},
		{"9aba2d20-ca23-403e-b127-6bc967eec751", materialsResourceID, "content.materials.publish", "Publicar Materiales", "Publicar materiales para estudiantes", "publish", "unit"},

		// --- memberships (resource 20000000-…-21) ---
		{"989bbeb5-9884-4728-8c98-d87d9d27f088", L4_RESOURCE_MEMBERSHIPS_ID, "academic.memberships.create", "Crear Membresías", "Asignar usuarios a unidades académicas", "create", "school"},
		{"c78ca0cf-0b6d-4e70-94d8-5ee0d021233f", L4_RESOURCE_MEMBERSHIPS_ID, "academic.memberships.delete", "Eliminar Membresías", "Eliminar membresías de unidades", "delete", "school"},
		{"28dfb6b5-c680-4442-8530-73b67199fbcb", L4_RESOURCE_MEMBERSHIPS_ID, "academic.memberships.read", "Ver Membresías", "Ver membresías de unidades académicas", "read", "school"},
		{"0f53cce3-0133-4f93-9b8c-7c62b3a8eb3c", L4_RESOURCE_MEMBERSHIPS_ID, "academic.memberships.update", "Editar Membresías", "Modificar membresías", "update", "school"},

		// --- my_memberships (resource 20000000-…-22) ---
		// Permiso único del feature "mis materias" del alumno. Vive bajo un path
		// PROPIO (academic.my_memberships.*) — NO bajo academic.memberships.* —
		// para que el gate de menú por path-prefix NO le filtre el item admin
		// "memberships" (roster de unidad). Cubre las tres caras del feature self
		// del student: visibilidad del item de menú "Mis materias",
		// slot.permission de la pantalla my-memberships-list y route gate del
		// dato. El contrato KMP consume el lector A
		// (GET /api/v1/me/subject-offerings). El teacher (literal
		// `academic.memberships.read`) NO toca este path, así que no ve "Mis
		// materias"; el admin sí lo ve vía `academic.*`. Reintroducido en N1.7 F1
		// sobre el modelo de sesiones.
		{L4_PERM_MY_MEMBERSHIPS_READ_OWN_ID, L4_RESOURCE_MY_MEMBERSHIPS_ID, "academic.my_memberships.read:own", "Ver Mis Materias", "Ver el item de menú de materias propias del alumno", "read:own", "unit"},

		// --- my_grades (resource 20000000-…-24) ---
		// Permiso único del feature "mis notas" del alumno (N3 F4 — consulta de
		// notas). Vive bajo un path PROPIO (academic.my_grades.*) — NO bajo
		// academic.grades.* — para que el gate de menú por path-prefix NO le filtre
		// el item admin "grades" (CRUD docente). Cubre las tres caras del feature
		// self del student: visibilidad del item de menú "Mis notas",
		// slot.permission de la pantalla my-grades-list y route gate del dato. El
		// contrato KMP consume el lector self GET /api/v1/me/grades. Espejo de
		// academic.my_memberships.read:own.
		{L4_PERM_MY_GRADES_READ_OWN_ID, L4_RESOURCE_MY_GRADES_ID, "academic.my_grades.read:own", "Ver Mis Notas", "Ver el item de menú de notas propias del alumno", "read:own", "unit"},

		// --- my_wards (resources b4000000-…-25/26/27/28) — plan 024 F1 ---
		// Vistas `:own` del acudido para el rol guardián. Cada una bajo su path
		// propio academic.my_wards_* (gate de menú por prefijo distinto del CRUD
		// docente). El lector real que las sirve llega en F3; aquí solo se declaran.
		{L4_PERM_MY_WARDS_GRADES_READ_OWN_ID, L4_RESOURCE_MY_WARDS_GRADES_ID, "academic.my_wards_grades.read:own", "Ver Notas de Acudidos", "Ver notas de los alumnos a cargo del representante", "read:own", "unit"},
		{L4_PERM_MY_WARDS_ATTENDANCE_READ_OWN_ID, L4_RESOURCE_MY_WARDS_ATTENDANCE_ID, "academic.my_wards_attendance.read:own", "Ver Asistencia de Acudidos", "Ver asistencia de los alumnos a cargo del representante", "read:own", "unit"},
		{L4_PERM_MY_WARDS_ANNOUNCEMENTS_READ_OWN_ID, L4_RESOURCE_MY_WARDS_ANNOUNCEMENTS_ID, "academic.my_wards_announcements.read:own", "Ver Anuncios de Acudidos", "Ver anuncios dirigidos a los alumnos a cargo del representante", "read:own", "unit"},
		{L4_PERM_MY_WARDS_MATERIALS_READ_OWN_ID, L4_RESOURCE_MY_WARDS_MATERIALS_ID, "academic.my_wards_materials.read:own", "Ver Materiales de Acudidos", "Ver materiales de los alumnos a cargo del representante", "read:own", "unit"},
		{L4_PERM_MY_WARDS_ASSESSMENTS_READ_OWN_ID, L4_RESOURCE_MY_WARDS_ASSESSMENTS_ID, "academic.my_wards_assessments.read:own", "Ver Evaluaciones de Acudidos", "Ver evaluaciones de los alumnos a cargo del representante", "read:own", "unit"},

		// Poda menú (2026-05-29): permisos `admin.permissions_mgmt.*`
		// eliminados junto con el recurso `permissions_mgmt`.

		// Eliminado (2026-06-15): permisos `reports.progress.*` junto con el
		// recurso `progress`. Su pantalla SDUI apuntaba a /api/v1/stats/student
		// (inexistente → 404) y era redundante con el dashboard nativo del alumno.

		// Poda menú (2026-05-29): permisos `admin.roles.*` eliminados junto
		// con el recurso `roles`.

		// --- schools (resource 20000000-…-11) ---
		{"611df7ce-b4cd-474f-901d-9bfd8873a9c1", L4_RESOURCE_SCHOOLS_ID, "admin.schools.create", "Crear Escuelas", "Crear nuevas instituciones educativas", "create", "system"},
		{"5bd8088b-1506-4b22-aa7e-9e4eb50de24e", L4_RESOURCE_SCHOOLS_ID, "admin.schools.delete", "Eliminar Escuelas", "Eliminar escuelas del sistema", "delete", "system"},
		{"8545c3be-3117-40a1-b1fb-da78d6233ae1", L4_RESOURCE_SCHOOLS_ID, "admin.schools.manage", "Gestionar Escuela", "Control total de la escuela", "manage", "school"},
		{"bc15c7a1-f203-46e0-80be-2850fad94b0e", L4_RESOURCE_SCHOOLS_ID, "admin.schools.read", "Ver Escuelas", "Ver información de escuelas", "read", "system"},
		{"2b823ad1-d875-4951-9c85-3baafa3f1f65", L4_RESOURCE_SCHOOLS_ID, "admin.schools.update", "Editar Escuelas", "Modificar datos de escuelas", "update", "school"},

		// Poda menú (2026-06-01): permisos admin.screen_instances.* y
		// admin.screen_templates.* eliminados junto con sus recursos — el CRUD
		// de configuración SDUI se reimplementó en el admin-tool de Go. El
		// permiso screens.read (lectura de pantallas desde mobile) se conserva.

		// --- screens mobile (resource 20000000-…-52) ---
		{"2b31df13-4c54-43fc-8bcd-8a9265fba1a0", L4_RESOURCE_SCREENS_ID, "screens.read", "Leer Pantallas (Mobile)", "Leer configuración de pantallas desde mobile", "read", "system"},

		// --- subjects (resource 20000000-…-32) ---
		{"30000000-0000-0000-0000-000000000001", L4_RESOURCE_SUBJECTS_ID, "academic.subjects.create", "Crear Materia", "Crear materias en el plan de estudios", "create", "school"},
		{"30000000-0000-0000-0000-000000000002", L4_RESOURCE_SUBJECTS_ID, "academic.subjects.read", "Ver Materias", "Ver materias del plan de estudios", "read", "school"},
		{"30000000-0000-0000-0000-000000000003", L4_RESOURCE_SUBJECTS_ID, "academic.subjects.update", "Editar Materia", "Modificar datos de materias", "update", "school"},
		{"30000000-0000-0000-0000-000000000004", L4_RESOURCE_SUBJECTS_ID, "academic.subjects.delete", "Eliminar Materia", "Eliminar materias del plan de estudios", "delete", "school"},

		// --- subject_offerings (resource b4000000-…-23): plan 010 N1.7 / ADR 0009 ---
		// `enroll` cubre alta y baja de matrícula de la sesión (inscripción por lote).
		{"30000000-0000-0000-0000-000000000011", L4_RESOURCE_SUBJECT_OFFERINGS_ID, "academic.subject_offerings.create", "Crear Sesión de Materia", "Crear una sesión (materia + sección + período + docente)", "create", "school"},
		{"30000000-0000-0000-0000-000000000012", L4_RESOURCE_SUBJECT_OFFERINGS_ID, "academic.subject_offerings.read", "Ver Sesiones de Materia", "Ver las sesiones de materia y sus inscritos", "read", "school"},
		{"30000000-0000-0000-0000-000000000013", L4_RESOURCE_SUBJECT_OFFERINGS_ID, "academic.subject_offerings.update", "Editar Sesión de Materia", "Modificar una sesión (p.ej. reasignar docente)", "update", "school"},
		{"30000000-0000-0000-0000-000000000014", L4_RESOURCE_SUBJECT_OFFERINGS_ID, "academic.subject_offerings.delete", "Eliminar Sesión de Materia", "Eliminar una sesión y sus inscripciones (cascade)", "delete", "school"},
		{"30000000-0000-0000-0000-000000000015", L4_RESOURCE_SUBJECT_OFFERINGS_ID, "academic.subject_offerings.enroll", "Inscribir en Sesión", "Inscribir/desinscribir alumnos en una sesión de materia", "enroll", "school"},

		// --- stats (resource 20000000-…-41) ---
		{"8a9fbae4-1b64-4870-ad14-41c436348bcc", L4_RESOURCE_STATS_ID, "reports.stats.global", "Estadísticas Globales", "Ver estadísticas de toda la plataforma", "global", "system"},
		{"f35d45c1-9539-422d-974f-5075d8f9b296", L4_RESOURCE_STATS_ID, "reports.stats.school", "Estadísticas de Escuela", "Ver estadísticas de la institución", "school", "school"},
		{"f47983b3-a721-461e-8de5-05fea4eda3fe", L4_RESOURCE_STATS_ID, "reports.stats.unit", "Estadísticas de Unidad", "Ver estadísticas de la clase", "unit", "unit"},

		// --- units (resource 20000000-…-20) ---
		{"619f6f66-6806-4894-a965-3c266a483be3", L4_RESOURCE_UNITS_ID, "academic.units.create", "Crear Unidades", "Crear unidades académicas (clases, grados)", "create", "school"},
		{"8d3a079d-5b6c-452f-ab49-b725547a052c", L4_RESOURCE_UNITS_ID, "academic.units.delete", "Eliminar Unidades", "Eliminar unidades académicas", "delete", "school"},
		{"61633d6c-aa56-40c1-a048-8e21f2893058", L4_RESOURCE_UNITS_ID, "academic.units.read", "Ver Unidades", "Ver unidades académicas", "read", "school"},
		{"4809f4d8-16dc-4222-9e12-1fca5f3c7ab7", L4_RESOURCE_UNITS_ID, "academic.units.update", "Editar Unidades", "Modificar unidades académicas", "update", "school"},

		// --- users (resource 20000000-…-10) ---
		{"eff25f87-711d-43a5-b8d3-1e3fb6be6a19", L4_RESOURCE_USERS_ID, "admin.users.create", "Crear Usuarios", "Crear nuevos usuarios en el sistema", "create", "system"},
		{"4129c4b5-89a1-4908-8b21-b6289e1ad095", L4_RESOURCE_USERS_ID, "admin.users.delete", "Eliminar Usuarios", "Eliminar usuarios del sistema", "delete", "system"},
		{"1ae1ad50-857c-4601-8378-5fd25128f11d", L4_RESOURCE_USERS_ID, "admin.users.read", "Ver Usuarios", "Ver información de usuarios", "read", "school"},
		{"813077d4-4624-4817-b3f2-69d60f6cb7a9", L4_RESOURCE_USERS_ID, "admin.users.update", "Editar Usuarios", "Modificar datos de usuarios", "update", "school"},
		{"8098577f-e5d8-4e07-aee7-4c1521cbe88b", L4_RESOURCE_USERS_ID, "admin.users.read:own", "Ver Perfil Propio", "Ver propio perfil de usuario", "read:own", "system"},
		{"668b0d86-f4cb-45cd-a8c4-afa8f6cfe9b6", L4_RESOURCE_USERS_ID, "admin.users.update:own", "Editar Perfil Propio", "Modificar propio perfil", "update:own", "system"},

		// Poda menú (2026-05-29): permisos `academic.guardian_relations.*`
		// eliminados junto con el recurso `guardian_relations`.

		// --- invitations (resource …61) ---
		{"77474519-5e37-4cf4-b594-af6b0b5cd56d", L4_RESOURCE_INVITATIONS_ID, "academic.invitations.create", "Crear Invitación", "Generar códigos de invitación a colegio/unidad", "create", "school"},
		{"4873dd88-4268-4075-80ce-8e361802ae42", L4_RESOURCE_INVITATIONS_ID, "academic.invitations.read", "Ver Invitaciones", "Listar códigos de invitación", "read", "school"},
		{"d872dc59-2b6e-405f-acf3-d2c144d99d0f", L4_RESOURCE_INVITATIONS_ID, "academic.invitations.revoke", "Revocar Invitación", "Desactivar un código de invitación", "revoke", "school"},

		// --- join_requests (resource …62) ---
		{"c183939b-16af-4393-b52b-df915503b952", L4_RESOURCE_JOIN_REQUESTS_ID, "academic.join_requests.create", "Redimir Invitación", "Crear una solicitud de ingreso al redimir un código", "create", "school"},
		{"aca9a4cc-d572-4a21-8e58-ff3ccceb7daf", L4_RESOURCE_JOIN_REQUESTS_ID, "academic.join_requests.read", "Ver Solicitudes de Ingreso", "Ver la bandeja de solicitudes pendientes", "read", "school"},
		{"fed74f16-24ce-4247-8dea-abaaf927fbd6", L4_RESOURCE_JOIN_REQUESTS_ID, "academic.join_requests.reject", "Rechazar Solicitud", "Rechazar una solicitud de ingreso", "reject", "school"},

		// --- join_request_approvals (resource …63): SELLO × TIPO ---
		// El doble gate de ingreso (colegio→unidad) tiene UN permiso POR SELLO y
		// POR TIPO de invitación. El path es
		// `academic.join_request_approvals.<sello>.<tipo>` con sello ∈
		// {school, unit} y tipo ∈ {student, teacher, guardian}. La key de acción
		// (`<sello>.<tipo>`) es única dentro del recurso (UNIQUE resource_id+action).
		// approve.go evalúa el permiso del sello concreto que va a firmar, así un
		// rol puede firmar el sello de unidad de su clase (teacher → unit.student)
		// sin firmar el de colegio. Los wildcards amplios (school_admin
		// `academic.*`, super_admin `*`) cubren ambos sub-namespaces por subárbol.
		{"437fd7c4-bbcc-4359-87b4-b3444c8f2abe", L4_RESOURCE_JOIN_REQUEST_APPROVALS_ID, "academic.join_request_approvals.school.student", "Aprobar Alumnos (colegio)", "Firmar el sello de COLEGIO de solicitudes con rol student", "school.student", "school"},
		{"2637d38f-6ae2-4e65-bb5e-c28b8acc35d8", L4_RESOURCE_JOIN_REQUEST_APPROVALS_ID, "academic.join_request_approvals.school.teacher", "Aprobar Profesores (colegio)", "Firmar el sello de COLEGIO de solicitudes con rol teacher", "school.teacher", "school"},
		{"d369b2ac-b84a-4f8d-8099-ffb27128bc10", L4_RESOURCE_JOIN_REQUEST_APPROVALS_ID, "academic.join_request_approvals.school.guardian", "Aprobar Apoderados (colegio)", "Firmar el sello de COLEGIO de solicitudes con rol guardian", "school.guardian", "school"},
		{"d75e2c2a-7b9b-4472-8fcd-98d7ca74ce9e", L4_RESOURCE_JOIN_REQUEST_APPROVALS_ID, "academic.join_request_approvals.unit.student", "Aprobar Alumnos (unidad)", "Firmar el sello de UNIDAD de solicitudes con rol student", "unit.student", "school"},
		{"ac9b949a-896d-4b62-a005-ab97b00b96a6", L4_RESOURCE_JOIN_REQUEST_APPROVALS_ID, "academic.join_request_approvals.unit.teacher", "Aprobar Profesores (unidad)", "Firmar el sello de UNIDAD de solicitudes con rol teacher", "unit.teacher", "school"},
		{"63551147-a030-420f-a60b-a4e05e731040", L4_RESOURCE_JOIN_REQUEST_APPROVALS_ID, "academic.join_request_approvals.unit.guardian", "Aprobar Apoderados (unidad)", "Firmar el sello de UNIDAD de solicitudes con rol guardian", "unit.guardian", "school"},

		// --- audit (resource 20000000-…-70) ---
		{"a1000000-0000-0000-0000-000000000001", L4_RESOURCE_AUDIT_ID, "admin.audit.read", "Ver Auditoría", "Ver registros de auditoría del sistema", "read", "system"},
		{"a1000000-0000-0000-0000-000000000002", L4_RESOURCE_AUDIT_ID, "admin.audit.export", "Exportar Auditoría", "Exportar registros de auditoría", "export", "system"},

		// --- concept_types (resource 20000000-…-80) ---
		{L4_PERM_CONCEPT_TYPES_CREATE_ID, L4_RESOURCE_CONCEPT_TYPES_ID, "admin.concept_types.create", "Crear Tipo de Concepto", "Crear tipo de concepto", "create", "system"},
		{"c2000000-0000-0000-0000-000000000002", L4_RESOURCE_CONCEPT_TYPES_ID, "admin.concept_types.read", "Ver Tipos de Concepto", "Listar tipos de institucion", "read", "system"},
		{L4_PERM_CONCEPT_TYPES_UPDATE_ID, L4_RESOURCE_CONCEPT_TYPES_ID, "admin.concept_types.update", "Actualizar Tipo de Concepto", "Actualizar tipo de concepto", "update", "system"},

		// --- dashboard (resource 20000000-…-01) ---
		{"d0000000-0000-0000-0000-000000000001", L4_RESOURCE_DASHBOARD_ID, "dashboard.view", "Ver Dashboard", "Ver panel principal según rol del usuario", "view", "system"},

		// --- periods (resource 20000000-…-34) ---
		{"e1000000-0000-0000-0000-000000000001", L4_RESOURCE_PERIODS_ID, "academic.periods.read", "Ver Periodos", "Ver periodos académicos", "read", "school"},
		{"e1000000-0000-0000-0000-000000000002", L4_RESOURCE_PERIODS_ID, "academic.periods.create", "Crear Periodo", "Crear periodos académicos", "create", "school"},
		{"e1000000-0000-0000-0000-000000000003", L4_RESOURCE_PERIODS_ID, "academic.periods.update", "Editar Periodo", "Modificar periodos académicos", "update", "school"},
		{"e1000000-0000-0000-0000-000000000004", L4_RESOURCE_PERIODS_ID, "academic.periods.delete", "Eliminar Periodo", "Eliminar periodos académicos", "delete", "school"},
		{"e1000000-0000-0000-0000-000000000005", L4_RESOURCE_PERIODS_ID, "academic.periods.activate", "Activar Periodo", "Activar/desactivar periodos académicos", "activate", "school"},

		// --- grades (resource 20000000-…-35) ---
		{"e2000000-0000-0000-0000-000000000001", L4_RESOURCE_GRADES_ID, "academic.grades.read", "Ver Calificaciones", "Ver calificaciones de estudiantes", "read", "unit"},
		{"e2000000-0000-0000-0000-000000000002", L4_RESOURCE_GRADES_ID, "academic.grades.create", "Crear Calificación", "Registrar calificaciones", "create", "unit"},
		{"e2000000-0000-0000-0000-000000000003", L4_RESOURCE_GRADES_ID, "academic.grades.update", "Editar Calificación", "Modificar calificaciones", "update", "unit"},
		{"e2000000-0000-0000-0000-000000000004", L4_RESOURCE_GRADES_ID, "academic.grades.finalize", "Finalizar Calificación", "Finalizar y cerrar calificaciones", "finalize", "unit"},

		// --- attendance (resource 20000000-…-36) ---
		{"e3000000-0000-0000-0000-000000000001", L4_RESOURCE_ATTENDANCE_ID, "academic.attendance.read", "Ver Asistencia", "Ver registros de asistencia", "read", "unit"},
		{"e3000000-0000-0000-0000-000000000002", L4_RESOURCE_ATTENDANCE_ID, "academic.attendance.create", "Registrar Asistencia", "Registrar asistencia de estudiantes", "create", "unit"},
		{L4_PERM_ATTENDANCE_UPDATE_ID, L4_RESOURCE_ATTENDANCE_ID, "academic.attendance.update", "Actualizar Asistencia", "Actualizar asistencia", "update", "unit"},

		// Poda menú (2026-05-29): permisos `academic.schedules.*` y
		// `academic.calendar.*` eliminados junto con los recursos `schedules`
		// y `calendar`.

		// --- reports (resource 20000000-…-05) ---
		{"e7000000-0000-0000-0000-000000000001", L4_RESOURCE_REPORTS_ID, "reports.read", "Ver Reportes", "Ver reportes y estadísticas generales", "read", "school"},

		// --- context (resource 20000000-…-A0) ---
		{"f0000000-0000-0000-0000-000000000001", L4_RESOURCE_CONTEXT_ID, "context.browse_schools", "Explorar Escuelas", "Listar todas las escuelas para seleccion de contexto", "browse_schools", "system"},
		{"f0000000-0000-0000-0000-000000000002", L4_RESOURCE_CONTEXT_ID, "context.browse_units", "Explorar Unidades", "Listar todas las unidades de una escuela para seleccion de contexto", "browse_units", "system"},

		// --- notifications (resource 20000000-…-B0) ---
		{"f1000000-0000-0000-0000-000000000001", L4_RESOURCE_NOTIFICATIONS_ID, "notifications.read", "Ver Notificaciones", "Ver y gestionar notificaciones propias", "read", "system"},

		// --- menu (resource 20000000-…-C0) ---
		{"f2000000-0000-0000-0000-000000000001", L4_RESOURCE_MENU_ID, "menu.read", "Ver Menu", "Ver menu de navegación según permisos del rol", "read", "system"},
		{"f2000000-0000-0000-0000-000000000002", L4_RESOURCE_MENU_ID, "menu.full_read", "Ver Menu Completo", "Ver menu completo sin filtrar por permisos", "full_read", "system"},

		// --- system_settings (resource 20000000-…-90) ---
		{"d1000000-0000-0000-0000-000000000001", L4_RESOURCE_SYSTEM_SETTINGS_ID, "admin.system_settings.settings", "Configuración del Sistema", "Acceder a la configuración y mantenimiento del sistema", "settings", "system"},
		{"d1000000-0000-0000-0000-000000000002", L4_RESOURCE_SYSTEM_SETTINGS_ID, "admin.system_settings.read", "Leer Configuración del Sistema", "Leer configuración del sistema", "read", "system"},
		{"d1000000-0000-0000-0000-000000000003", L4_RESOURCE_SYSTEM_SETTINGS_ID, "admin.system_settings.update", "Actualizar Configuración del Sistema", "Actualizar configuración del sistema", "update", "system"},

		// --- messaging (resource b4000000-…-d0): plan 025 (WhatsApp) ---
		// Permisos del staff para operar la mensajería a las familias. Path de 3
		// segmentos (messaging.<sub>.<accion>); la API messaging los lee de los
		// grants del JWT (no consulta IAM). La acción es el path tras el recurso
		// `messaging` (UNIQUE resource_id+action, las 3 son distintas). Scope
		// system: la capability no se ata a una escuela; el alcance lo da el JWT.
		{L4_PERM_MESSAGING_SESSION_PAIR_ID, L4_RESOURCE_MESSAGING_ID, "messaging.session.pair", "Vincular Sesión de WhatsApp", "Vincular (parear) la sesión de WhatsApp de la escuela", "session.pair", "system"},
		{L4_PERM_MESSAGING_MESSAGE_SEND_ID, L4_RESOURCE_MESSAGING_ID, "messaging.message.send", "Enviar Mensaje de WhatsApp", "Enviar mensajes de WhatsApp a las familias", "message.send", "system"},
		{L4_PERM_MESSAGING_DEVICE_LINK_ID, L4_RESOURCE_MESSAGING_ID, "messaging.device.link", "Enlazar Dispositivo de WhatsApp", "Enlazar un dispositivo emisor de WhatsApp", "device.link", "system"},
	}
}

// buildL4Permissions materializa las specs declarativas en
// entities.Permission. Helper compartido por applyL4Permissions y por
// el accessor público l4.Permissions().
func buildL4Permissions() ([]entities.Permission, error) {
	specs := l4Permissions()
	perms := make([]entities.Permission, 0, len(specs))
	for _, s := range specs {
		pid, err := uuid.Parse(s.idStr)
		if err != nil {
			return nil, fmt.Errorf("parse permission id %s (%s): %w", s.idStr, s.name, err)
		}
		rid, err := uuid.Parse(s.resourceID)
		if err != nil {
			return nil, fmt.Errorf("parse resource_id %s for %s: %w", s.resourceID, s.name, err)
		}
		desc := s.description
		perms = append(perms, entities.Permission{
			ID:          pid,
			Name:        s.name,
			DisplayName: s.displayName,
			Description: &desc,
			ResourceID:  rid,
			Action:      s.action,
			Scope:       s.scope,
			IsActive:    true,
			IsSystem:    true,
		})
	}
	return perms, nil
}

// applyL4Permissions inserta las filas de iam.permissions definidas
// por L4. Idempotente vía OnConflict por id (alineado con
// applyL3Permissions). El UNIQUE compuesto (resource_id, action)
// también se respeta por design: no hay duplicados en el slice.
func applyL4Permissions(tx *gorm.DB) error {
	perms, err := buildL4Permissions()
	if err != nil {
		return err
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).CreateInBatches(&perms, 20).Error
}

// -----------------------------------------------------------------------
// P3-1: role_grants con patterns wildcard
// -----------------------------------------------------------------------

// roleGrantPatterns devuelve los patterns de iam.role_grants para cada
// rol. Sustituye al espejo 1:1 anterior por wildcards explícitos —
// principio "wildcard-first": si un rol cubre semánticamente todo el
// subárbol `prefix.*`, el grant es ese pattern (no la enumeración).
//
// Compatibilidad del matcher (edugo-shared/auth/permission_matcher.go):
// soporta `*` global, literales exactos, `prefix.*` (cubre `prefix`,
// `prefix.X`, `prefix.X.Y`), `*.suffix` y `prefix.*.suffix`.
//
// readonly_auditor: allow amplio + denies `*.suffix` sobre verbos
// mutativos. Cualquier permiso mutativo nuevo del catálogo cae
// automáticamente bajo el deny sin tocar el seed.
//
// ADR-6 (herencia de roles): los 5 alias que SÍ heredan
// (school_director/coordinator/assistant → school_admin;
// assistant_teacher/observer → teacher) ya NO aparecen aquí — sus grants
// se resuelven desde el canónico vía roles.parent_role_id y se aplanan en
// el login. readonly_auditor permanece standalone por no ser un superset
// exacto de teacher (ver l4RoleSpecs).
func roleGrantPatterns() map[string][]string {
	// Espejo de role IDs definidos en capas previas (L0, L1).
	const (
		l0RoleSuperAdmin         = "10000000-0000-0000-0000-000000000001" // layers.L0_ROLE_SUPER_ADMIN_ID
		l1RoleAnnouncementViewer = "b1000000-0000-0000-0000-000000000001" // layers.L1_ROLE_ANNOUNCEMENT_VIEWER_ID
	)

	// Bloques reusables.
	schoolAdminPatterns := []string{
		"academic.*",
		"admin.*",
		"content.*",
		// Grant LITERAL a `context.browse_units` (ver/elegir las unidades de SU escuela
		// para el selector de contexto), NO el wildcard `context.*`: ese también otorga
		// `context.browse_schools` (scope `system`: "listar TODAS las escuelas del
		// sistema"), una capability de super_admin. Heredada por school_admin/coordinator
		// encendía el "Cambiar escuela" a un admin de UNA sola escuela y fallaba con 403
		// al elegir una escuela ajena. Mismo criterio que los literales de teacher/student
		// abajo: literal donde el wildcard sobre-otorgaría.
		"context.browse_units",
		"reports.*",
		"dashboard.*",
		"menu.*",
		"notifications.*",
		"screens.*",
		// Plan 025: el admin de escuela comunica a las familias por WhatsApp.
		// Wildcard del subárbol messaging.* (session.pair/message.send/device.link).
		// Lo heredan school_director/coordinator/assistant (ADR-6).
		"messaging.*",
	}
	teacherPatterns := []string{
		"academic.announcements.*",
		"academic.attendance.*",
		"academic.grades.*",
		// El docente lee membresías para el roster/directorio de su unidad
		// (unit-directory). Grant LITERAL a `.read`, NO el wildcard
		// `academic.memberships.*`: el profesor no crea, edita ni elimina
		// membresías (eso es de school_admin vía academic.*). Mismo criterio
		// que `academic.join_request_approvals.unit.student` — literal donde el
		// wildcard sobre-otorgaría. (DIFERIBLE en F0b: se conserva porque tiene
		// usos vivos roster/unit-directory; la vista "alumnos por materia" que
		// también lo usaba se retiró, ver plan 010 N1.7.)
		"academic.memberships.read",
		"academic.periods.*",
		// Plan 006 (Trozo A): el docente NO gestiona materias por defecto;
		// solo las ve (decisión del dueño). Grant LITERAL a `.read`, NO el
		// wildcard `academic.subjects.*`: crear/editar/eliminar materias es
		// de school_admin (vía academic.*). Mismo criterio que
		// `academic.memberships.read` arriba.
		"academic.subjects.read",
		// Plan 010 (N1.7): el docente VE sus sesiones de materia (lectura),
		// igual que ve materias. Grant LITERAL a `.read`, NO el wildcard
		// `academic.subject_offerings.*`: crear/editar/eliminar sesiones e
		// inscribir es de school_admin (vía academic.*). Si un docente
		// concreto debe inscribir, se le da `...enroll` vía user_grant.
		"academic.subject_offerings.read",
		"academic.units.*",
		// Onboarding (plan 005, sello × tipo): el profesor gestiona invitaciones
		// y solicitudes de su clase, pero SOLO firma el sello de UNIDAD de
		// alumnos (admite alumnos a SU clase) → grant literal a
		// `.unit.student`. NO firma el sello de COLEGIO (eso es de school_admin
		// vía academic.*) ni aprobaciones de profesores/apoderados; nunca el
		// wildcard `.*` sobre approvals.
		"academic.invitations.*",
		"academic.join_requests.*",
		"academic.join_request_approvals.unit.student",
		"content.assessments.*",
		"content.materials.*",
		"admin.system_settings.*",
		"reports.*",
		"dashboard.*",
		"menu.*",
		"notifications.*",
		"screens.*",
		// Plan 025: el docente comunica a las familias de su clase por WhatsApp.
		// Wildcard del subárbol messaging.*. Lo heredan assistant_teacher/observer
		// (ADR-6). El alumno NO lo recibe (familias = destinatarias, no emisoras).
		"messaging.*",
	}
	studentPatterns := []string{
		"academic.announcements.*",
		"academic.attendance.*",
		// El alumno NO recibe el wildcard `academic.grades.*` (CRUD docente):
		// ese grant lo dejaba ver/crear/editar notas ajenas vía GET/POST /grades
		// y ver el menú "Calificaciones" (grades-list) — fuga de privacidad
		// (N3 F4.1, decisión del dueño 2026-06-06). Su única lectura de notas es
		// el feature self `academic.my_grades.read:own` (abajo), que sirve solo
		// sus propias notas vía GET /api/v1/me/grades.
		// Permiso ÚNICO del feature "mis materias" del alumno (visibilidad de
		// menú + slot.permission + route gate del dato). Grant LITERAL al
		// recurso my_memberships, bajo un path PROPIO (academic.my_memberships.*)
		// — NUNCA el wildcard `academic.memberships.*` (eso listaría/mutaría la
		// unidad) ni un `:own` bajo `academic.memberships.` (eso le filtraría el
		// item de menú admin "memberships" por el gate path-prefix). No matchea
		// `academic.memberships.read` amplio, así que el alumno sigue recibiendo
		// 403 en GET /memberships (listar unidad). El dato propio llega vía el
		// lector A GET /api/v1/me/subject-offerings. Reintroducido en N1.7 F1.
		"academic.my_memberships.read:own",
		// Permiso ÚNICO del feature "mis notas" del alumno (N3 F4 — consulta de
		// notas). Grant LITERAL al recurso my_grades, bajo un path PROPIO
		// (academic.my_grades.*) — NO el wildcard `academic.grades.*` (eso
		// permanece como CRUD docente, no es el lector self del alumno). Cubre la
		// visibilidad del item de menú "Mis notas", slot.permission de
		// my-grades-list y route gate del dato propio (GET /api/v1/me/grades).
		// Espejo de academic.my_memberships.read:own.
		"academic.my_grades.read:own",
		"content.assessments_student.*",
		"content.materials.*",
		"dashboard.*",
		"menu.*",
		"notifications.*",
		"screens.*",
	}
	guardianPatterns := []string{
		"academic.announcements.*",
		"academic.attendance.*",
		// academic.grades.* ELIMINADO en F1 (privacidad): el guardián ve solo a su
		// hijo vía academic.my_wards_*.read:own (lector real en F3).
		"academic.guardian_relations.*", // revertir poda 2026-05-29 (revive rutas backend)
		"academic.my_wards_grades.read:own",
		"academic.my_wards_attendance.read:own",
		"academic.my_wards_announcements.read:own",
		"academic.my_wards_materials.read:own",
		"academic.my_wards_assessments.read:own",
		"content.assessments.*",
		"content.materials.*",
		"admin.system_settings.*",
		"reports.read",
		"dashboard.*",
		"menu.*",
		"notifications.*",
		"screens.*",
	}

	out := map[string][]string{
		// L0 super_admin: acceso total al sistema, un solo pattern.
		l0RoleSuperAdmin: {"*"},

		// L1 announcement_viewer: literales (sin cambio respecto al
		// estado previo — sólo lee anuncios y system_settings).
		l1RoleAnnouncementViewer: {
			"academic.announcements.read",
			"admin.system_settings.read",
		},

		// Canónicos L4.
		L4_ROLE_STUDENT_ID:      studentPatterns,
		L4_ROLE_TEACHER_ID:      teacherPatterns,
		L4_ROLE_GUARDIAN_ID:     guardianPatterns,
		L4_ROLE_SCHOOL_ADMIN_ID: schoolAdminPatterns,

		// ADR-6 (herencia de roles): los alias school_director,
		// school_coordinator, school_assistant (→ school_admin) y
		// assistant_teacher, observer (→ teacher) NO declaran allow
		// propios; heredan el set completo del canónico vía
		// roles.parent_role_id. La herencia se resuelve y aplana en el
		// login.

		// readonly_auditor NO hereda (ver l4RoleSpecs): conserva su allow
		// read-only amplio. Su deny mutativo está en roleGrantDenyPatterns.
		L4_ROLE_READONLY_AUDITOR_ID: {
			"academic.*",
			"content.*",
			"reports.*",
			"dashboard.*",
			"menu.*",
			"notifications.*",
			"screens.*",
			"admin.users.read",
			"admin.users.read:own",
			"admin.system_settings.read",
		},
	}
	return out
}

// roleGrantDenyPatterns devuelve los patterns con effect=deny por rol.
// Para readonly_auditor: bloqueo wildcard `*.suffix` de cualquier verbo
// mutativo. Esto opera sobre el allow amplio del rol y deja como neto
// sólo accesos read-only — cualquier permiso mutativo nuevo del catálogo
// cae automáticamente aquí.
func roleGrantDenyPatterns() map[string][]string {
	return map[string][]string{
		L4_ROLE_READONLY_AUDITOR_ID: {
			"*.create",
			"*.update",
			"*.delete",
			"*.publish",
			"*.finalize",
			"*.activate",
			"*.approve",
			"*.grade",
			"*.attempt",
			"*.assign",
			"*.review",
			"*.manage",
			"*.request",
			// Onboarding (plan 005): higiene deny-wins. readonly_auditor
			// tiene allow `academic.*`; sin estos deny podría revocar
			// invitaciones, rechazar solicitudes o aprobar ingresos. El
			// namespace de aprobación no es un verbo de mutación clásico, así
			// que se deniega completo (la acción ES el rol).
			"*.revoke",
			"*.reject",
			"academic.join_request_approvals.*",
		},
	}
}

// applyL4RoleGrants siembra iam.role_grants con patterns wildcard por
// rol según roleGrantPatterns(), más los grants literales filtrados
// para readonly_auditor (sin cambios respecto al pase anterior).
//
// ID determinístico: SHA1 sobre (roleID + ":" + pattern + ":allow") — el
// mismo par produce el mismo UUID, idempotente.
//
// OnConflict por la clave natural (role_id, pattern, effect) — coincide
// con el UNIQUE compuesto `uq_role_grants_role_pattern_effect`.
func applyL4RoleGrants(tx *gorm.DB) error {
	now := time.Now().UTC()
	allowPatterns := roleGrantPatterns()
	denyPatterns := roleGrantDenyPatterns()

	total := 0
	for _, ps := range allowPatterns {
		total += len(ps)
	}
	for _, ps := range denyPatterns {
		total += len(ps)
	}
	grants := make([]entities.RoleGrant, 0, total)

	// Orden estable de inserción.
	roleOrder := []string{
		"10000000-0000-0000-0000-000000000001", // L0 super_admin
		"b1000000-0000-0000-0000-000000000001", // L1 announcement_viewer
		L4_ROLE_STUDENT_ID,
		L4_ROLE_TEACHER_ID,
		L4_ROLE_GUARDIAN_ID,
		L4_ROLE_SCHOOL_ADMIN_ID,
		L4_ROLE_SCHOOL_DIRECTOR_ID,
		L4_ROLE_SCHOOL_COORDINATOR_ID,
		L4_ROLE_SCHOOL_ASSISTANT_ID,
		L4_ROLE_ASSISTANT_TEACHER_ID,
		L4_ROLE_OBSERVER_ID,
		L4_ROLE_READONLY_AUDITOR_ID,
	}
	for _, ridStr := range roleOrder {
		roleID, err := uuid.Parse(ridStr)
		if err != nil {
			return fmt.Errorf("parse role_id %s: %w", ridStr, err)
		}
		for _, pattern := range allowPatterns[ridStr] {
			derivedID := uuid.NewSHA1(uuid.NameSpaceOID, []byte(roleID.String()+":"+pattern+":allow"))
			grants = append(grants, entities.RoleGrant{
				ID:        derivedID,
				RoleID:    roleID,
				Pattern:   pattern,
				Effect:    "allow",
				CreatedAt: now,
			})
		}
		for _, pattern := range denyPatterns[ridStr] {
			derivedID := uuid.NewSHA1(uuid.NameSpaceOID, []byte(roleID.String()+":"+pattern+":deny"))
			grants = append(grants, entities.RoleGrant{
				ID:        derivedID,
				RoleID:    roleID,
				Pattern:   pattern,
				Effect:    "deny",
				CreatedAt: now,
			})
		}
	}

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "role_id"}, {Name: "pattern"}, {Name: "effect"}},
		DoNothing: true,
	}).CreateInBatches(&grants, 50).Error
}
