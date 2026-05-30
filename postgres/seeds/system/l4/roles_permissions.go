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
}

// l4RoleSpecs retorna las specs declarativas de los 11 roles que L4
// siembra (5 canónicos + 6 alias). Helper compartido por applyL4Roles
// y por el accessor público l4.Roles() — la lógica de construcción del
// slice de entities vive una sola vez (buildL4Roles).
func l4RoleSpecs() []l4RoleSpec {
	return []l4RoleSpec{
		{
			idStr:       L4_ROLE_STUDENT_ID,
			name:        L4_ROLE_STUDENT_NAME,
			displayName: "Estudiante",
			description: "Alumno inscrito en una unidad académica.",
			scope:       "unit",
		},
		{
			idStr:       L4_ROLE_TEACHER_ID,
			name:        L4_ROLE_TEACHER_NAME,
			displayName: "Profesor",
			description: "Docente con permisos de gestión de clase (asistencia, calificaciones, evaluaciones, materiales).",
			scope:       "unit",
		},
		{
			idStr:       L4_ROLE_GUARDIAN_ID,
			name:        L4_ROLE_GUARDIAN_NAME,
			displayName: "Apoderado",
			description: "Tutor legal o apoderado vinculado a uno o más estudiantes.",
			scope:       "unit",
		},
		// PRE-4: el rol `platform_admin` (L4_ROLE_ADMIN_*) fue
		// eliminado. Sus capacidades quedan cubiertas por `super_admin`
		// (L0) que ya tiene acceso global.
		{
			idStr:       L4_ROLE_SCHOOL_ADMIN_ID,
			name:        L4_ROLE_SCHOOL_ADMIN_NAME,
			displayName: "Administrador de Escuela",
			description: "Administrador con control total dentro de una institución educativa.",
			scope:       "school",
		},
		// --- Alias roles (heredan grants del canónico) ---
		{
			idStr:       L4_ROLE_SCHOOL_DIRECTOR_ID,
			name:        L4_ROLE_SCHOOL_DIRECTOR_NAME,
			displayName: "Director de Escuela",
			description: "Director de la institución educativa. Alias de school_admin (hereda todos sus permisos).",
			scope:       "school",
			parentIDStr: L4_ROLE_SCHOOL_ADMIN_ID,
		},
		{
			idStr:       L4_ROLE_SCHOOL_COORDINATOR_ID,
			name:        L4_ROLE_SCHOOL_COORDINATOR_NAME,
			displayName: "Coordinador de Escuela",
			description: "Coordinador académico de la institución. Alias de school_admin (hereda todos sus permisos).",
			scope:       "school",
			parentIDStr: L4_ROLE_SCHOOL_ADMIN_ID,
		},
		{
			idStr:       L4_ROLE_SCHOOL_ASSISTANT_ID,
			name:        L4_ROLE_SCHOOL_ASSISTANT_NAME,
			displayName: "Asistente de Escuela",
			description: "Personal de apoyo administrativo de la institución. Alias de school_admin (hereda todos sus permisos).",
			scope:       "school",
			parentIDStr: L4_ROLE_SCHOOL_ADMIN_ID,
		},
		{
			idStr:       L4_ROLE_ASSISTANT_TEACHER_ID,
			name:        L4_ROLE_ASSISTANT_TEACHER_NAME,
			displayName: "Profesor Asistente",
			description: "Docente auxiliar. Alias de teacher (hereda todos sus permisos).",
			scope:       "unit",
			parentIDStr: L4_ROLE_TEACHER_ID,
		},
		{
			idStr:       L4_ROLE_OBSERVER_ID,
			name:        L4_ROLE_OBSERVER_NAME,
			displayName: "Observador",
			description: "Observador con visibilidad sobre la clase. Alias de teacher (hereda todos sus permisos).",
			scope:       "unit",
			parentIDStr: L4_ROLE_TEACHER_ID,
		},
		{
			idStr:       L4_ROLE_READONLY_AUDITOR_ID,
			name:        L4_ROLE_READONLY_AUDITOR_NAME,
			displayName: "Auditor de Solo Lectura",
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
		roles = append(roles, entities.Role{
			ID:           id,
			Name:         s.name,
			DisplayName:  s.displayName,
			Description:  &desc,
			Scope:        s.scope,
			ParentRoleID: parentID,
			IsActive:     true,
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

		// Poda menú (2026-05-29): permisos `admin.permissions_mgmt.*`
		// eliminados junto con el recurso `permissions_mgmt`.

		// --- progress (resource 20000000-…-40) ---
		// NOTA: se descartó `progress:read:own` (zombie + sólo super_admin
		//       lo usaba; el backend no chequea ese permiso en runtime).
		{"89336e44-3636-4744-a056-aea878f57b18", L4_RESOURCE_PROGRESS_ID, "reports.progress.read", "Ver Progreso", "Ver progreso académico", "read", "unit"},
		{"d033dfd0-47a5-476b-b51a-f52f5fc66d7a", L4_RESOURCE_PROGRESS_ID, "reports.progress.read:own", "Ver Progreso Propio", "Ver propio progreso (usado por student)", "read:own", "unit"},
		{"19d017d1-ee5b-4fc3-828a-2d12056631b4", L4_RESOURCE_PROGRESS_ID, "reports.progress.update", "Actualizar Progreso", "Actualizar progreso de estudiantes", "update", "unit"},

		// Poda menú (2026-05-29): permisos `admin.roles.*` eliminados junto
		// con el recurso `roles`.

		// --- schools (resource 20000000-…-11) ---
		{"611df7ce-b4cd-474f-901d-9bfd8873a9c1", L4_RESOURCE_SCHOOLS_ID, "admin.schools.create", "Crear Escuelas", "Crear nuevas instituciones educativas", "create", "system"},
		{"5bd8088b-1506-4b22-aa7e-9e4eb50de24e", L4_RESOURCE_SCHOOLS_ID, "admin.schools.delete", "Eliminar Escuelas", "Eliminar escuelas del sistema", "delete", "system"},
		{"8545c3be-3117-40a1-b1fb-da78d6233ae1", L4_RESOURCE_SCHOOLS_ID, "admin.schools.manage", "Gestionar Escuela", "Control total de la escuela", "manage", "school"},
		{"bc15c7a1-f203-46e0-80be-2850fad94b0e", L4_RESOURCE_SCHOOLS_ID, "admin.schools.read", "Ver Escuelas", "Ver información de escuelas", "read", "system"},
		{"2b823ad1-d875-4951-9c85-3baafa3f1f65", L4_RESOURCE_SCHOOLS_ID, "admin.schools.update", "Editar Escuelas", "Modificar datos de escuelas", "update", "school"},

		// --- screen_instances (resource 20000000-…-51) ---
		{"fa35b956-665f-48f4-a51e-ad1393e72652", L4_RESOURCE_SCREEN_INSTANCES_ID, "admin.screen_instances.create", "Crear Instancias de Pantalla", "Crear nuevas instancias de pantalla", "create", "system"},
		{"4096f489-b3f8-49bd-8ecb-6e3588a85f84", L4_RESOURCE_SCREEN_INSTANCES_ID, "admin.screen_instances.delete", "Eliminar Instancias de Pantalla", "Eliminar instancias de pantalla configuradas", "delete", "system"},
		{"ebfd0911-43bf-42ef-9523-8dd93079db47", L4_RESOURCE_SCREEN_INSTANCES_ID, "admin.screen_instances.read", "Ver Instancias de Pantalla", "Ver instancias de pantalla configuradas", "read", "system"},
		{"1ad07392-4c86-4ec9-b249-b66be3f97ce8", L4_RESOURCE_SCREEN_INSTANCES_ID, "admin.screen_instances.update", "Actualizar Instancias de Pantalla", "Modificar instancias de pantalla existentes", "update", "system"},

		// --- screen_templates (resource 20000000-…-50) ---
		{"52011396-5981-4c59-a772-1f353d10a3e9", L4_RESOURCE_SCREEN_TEMPLATES_ID, "admin.screen_templates.create", "Crear Templates de Pantalla", "Crear nuevos templates de pantalla", "create", "system"},
		{"b6db1991-4a2c-429a-9c45-0ed177b6e3ed", L4_RESOURCE_SCREEN_TEMPLATES_ID, "admin.screen_templates.delete", "Eliminar Templates de Pantalla", "Eliminar templates de pantalla del sistema", "delete", "system"},
		{"3d89c941-cbe5-4c1b-8cf0-0b55b4aaa313", L4_RESOURCE_SCREEN_TEMPLATES_ID, "admin.screen_templates.read", "Ver Templates de Pantalla", "Ver templates de pantalla del sistema", "read", "system"},
		{"e5bf88e6-73ff-40d4-93a4-8c787d3930af", L4_RESOURCE_SCREEN_TEMPLATES_ID, "admin.screen_templates.update", "Actualizar Templates de Pantalla", "Modificar templates de pantalla existentes", "update", "system"},

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

		// --- join_request_approvals (resource …63): la acción ES el rol que se admite ---
		{"8881638e-6f57-4849-9238-bcf8b7af5a93", L4_RESOURCE_JOIN_REQUEST_APPROVALS_ID, "academic.join_request_approvals.student", "Aprobar Alumnos", "Firmar el visto bueno de solicitudes con rol student", "student", "school"},
		{"f423efb0-ec3f-4d2e-bb45-f4bc2262156c", L4_RESOURCE_JOIN_REQUEST_APPROVALS_ID, "academic.join_request_approvals.teacher", "Aprobar Profesores", "Firmar el visto bueno de solicitudes con rol teacher", "teacher", "school"},
		{"995a9f3d-c1d9-4e30-8bdf-feb7aee84178", L4_RESOURCE_JOIN_REQUEST_APPROVALS_ID, "academic.join_request_approvals.guardian", "Aprobar Apoderados", "Firmar el visto bueno de solicitudes con rol guardian", "guardian", "school"},

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
		"context.*",
		"reports.*",
		"dashboard.*",
		"menu.*",
		"notifications.*",
		"screens.*",
	}
	teacherPatterns := []string{
		"academic.announcements.*",
		"academic.attendance.*",
		"academic.grades.*",
		// El docente lee membresías para el roster/directorio de su unidad
		// (unit-directory). Grant LITERAL a `.read`, NO el wildcard
		// `academic.memberships.*`: el profesor no crea, edita ni elimina
		// membresías (eso es de school_admin vía academic.*). Mismo criterio
		// que `academic.join_request_approvals.student` — literal donde el
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
		// Onboarding (plan 005): el profesor gestiona invitaciones y
		// solicitudes de su clase, pero SOLO firma aprobaciones de alumnos
		// (no de profesores ni apoderados) → grant literal a `.student`,
		// nunca el wildcard `.*` sobre approvals.
		"academic.invitations.*",
		"academic.join_requests.*",
		"academic.join_request_approvals.student",
		"content.assessments.*",
		"content.materials.*",
		"admin.users.*",
		"admin.system_settings.*",
		"reports.*",
		"dashboard.*",
		"menu.*",
		"notifications.*",
		"screens.*",
	}
	studentPatterns := []string{
		"academic.announcements.*",
		"academic.attendance.*",
		"academic.grades.*",
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
		"content.assessments.*",
		"content.assessments_student.*",
		"content.materials.*",
		"admin.system_settings.*",
		"reports.progress.*",
		"dashboard.*",
		"menu.*",
		"notifications.*",
		"screens.*",
	}
	guardianPatterns := []string{
		"academic.announcements.*",
		"academic.attendance.*",
		"academic.grades.*",
		"content.assessments.*",
		"content.materials.*",
		"admin.users.*",
		"admin.system_settings.*",
		"reports.progress.*",
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
