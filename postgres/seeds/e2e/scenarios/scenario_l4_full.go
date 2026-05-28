package scenarios

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l4"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
)

// L4Full valida end-to-end la capa L4 del seed system (Fase 6 del
// rebuild). A diferencia de l0_minimal..l3_isolation —que validan una
// capa puntual y se enfocan en la presencia de sus filas— l4_full
// asume que la BD bajo test tiene L0..L4 aplicadas (testdb.StartPostgres
// las aplica vía system.ApplySystem sin filtro de capa), y comprueba
// que la matriz "screens accesibles por rol" derivada de
// iam.role_permissions + iam.permissions + ui_config.resource_screens
// se mantenga estable para los 3 roles representativos del sistema:
//
//   - super_admin (L0): rol del bootstrap del sistema. Tiene asignados
//     SÓLO los 4 permisos `announcements:*` de L0 (no recibe permisos
//     adicionales en L4 — L4 siembra los otros 5 roles).
//   - teacher (L4): rol con CRUD sobre lo que produce en su unidad
//     (clase) — assessments, materials, grades, attendance, etc.
//   - student (L4): rol con sólo lectura sobre lo que tiene asignado
//     más la acción `assessments:attempt`.
//
// Ver `phase-6-layer-l4/{requirements,design}.md` (F6-REQ-6.x y
// F6-REQ-7.x).
//
// SQL puro (Opción A heredada de F5): no se levantan API server ni
// KMP runtime. La resolución de `GET /menu` y
// `screen-config/resolve/key/:key` se simula leyendo las mismas tablas
// que consumen los respectivos use-cases (`navigation.Repository`,
// `ResolveByKeyUseCase`):
//
//   - GET /menu              → iam.resources (is_menu_visible=true) +
//                              ui_config.resource_screens, filtrando
//                              por las resource_keys que asoma alguno
//                              de los permisos del rol; los padres del
//                              árbol se marcan visibles por
//                              propagación (ver
//                              `core/usecase/menu/get_user_menu.go::propagateVisibility`).
//   - screen-config/resolve  → ui_config.screen_instances filtrado por
//                              screen_key + parse del slot_data JSON
//                              (mismo flujo que
//                              `ResolveByKeyUseCase.Execute`).
//
// Detección de drift: la matriz expected = computeAccessMatrix() se
// construye PROGRAMÁTICAMENTE al inicio del test desde la BD; los
// asserts comparan el conteo de screens accesibles contra esa matriz
// derivada. Eso evita hardcodear números que cambien al sumar
// pantallas a L4. Lo que sí se hardcodea (porque es regresión real)
// es:
//
//   - los 3 roles deben tener > 0 screens (no romper acceso al menú).
//   - super_admin no debe tener TODOS los screens (no debe heredar
//     accesos de L4 por accidente — sólo announcements:*).
//   - teacher y student deben tener subsets explícitamente distintos
//     (sanidad de aislamiento por rol).
type L4Full struct{}

// Manifest implementa framework.Scenario.
func (s *L4Full) Manifest() framework.ScenarioManifest {
	return framework.ScenarioManifest{
		Name:         "l4_full",
		Description:  "Valida la capa L4 (sistema completo) del seed system y la matriz screens-por-rol derivada de role_permissions + resource_screens. Exporta las constantes L4 al JSON E2E.",
		FixtureNames: []string{"l4_full_export"},
		Tags:         []string{"l4", "system", "rbac", "menu", "screen-config", "phase-6"},
	}
}

// BuildFixtures implementa framework.Scenario. Devuelve siempre
// instancias frescas para no compartir estado mutable entre
// aplicaciones.
func (s *L4Full) BuildFixtures(ctx *framework.ApplyContext) []framework.Fixture {
	return []framework.Fixture{
		&l4FullExport{},
	}
}

// l4FullExport es una fixture pasiva inline del scenario `l4_full`.
// Verifica la presencia (no la forma) de las filas sembradas por L4
// vía los conteos públicos expuestos por `l4.Resources()` /
// `l4.Permissions()` / etc. — y exporta el inventario base (cuenta de
// recursos L4, cuenta de roles L4, version semántica de la capa) al
// ApplyContext para que tests downstream (Kotlin) puedan referenciarlo.
//
// La fixture es deliberadamente mínima: la validación funcional fuerte
// vive en TestScenarioL4Full_Integration (matriz por rol). Aquí solo
// se garantiza que system.ApplySystem haya corrido y dejado las filas
// esperadas en BD — si la cuenta no coincide con el accessor estático
// de l4.*, hay un drift entre el código y lo que llegó a la BD.
type l4FullExport struct{}

// Manifest implementa framework.Fixture.
func (f *l4FullExport) Manifest() framework.FixtureManifest {
	return framework.FixtureManifest{
		Name:        "l4_full_export",
		Description: "Verifica la presencia (count match) de las filas L4 sembradas por system.ApplySystem y exporta el inventario base al JSON.",
		// PRE-4: claves E2EFixtureL4RoleAdmin{ID,Name} removidas —
		// el rol `platform_admin` fue eliminado del catálogo L4 y los
		// tests Kotlin del KMP no las consumen.
		Constants: map[string]string{
			"E2EFixtureL4SeedVersion":         layers.L4_SEED_VERSION,
			"E2EFixtureL4LayerName":           layers.L4_LAYER_NAME,
			"E2EFixtureL4RoleStudentID":       l4.L4_ROLE_STUDENT_ID,
			"E2EFixtureL4RoleStudentName":     l4.L4_ROLE_STUDENT_NAME,
			"E2EFixtureL4RoleTeacherID":       l4.L4_ROLE_TEACHER_ID,
			"E2EFixtureL4RoleTeacherName":     l4.L4_ROLE_TEACHER_NAME,
			"E2EFixtureL4RoleGuardianID":      l4.L4_ROLE_GUARDIAN_ID,
			"E2EFixtureL4RoleGuardianName":    l4.L4_ROLE_GUARDIAN_NAME,
			"E2EFixtureL4RoleSchoolAdminID":   l4.L4_ROLE_SCHOOL_ADMIN_ID,
			"E2EFixtureL4RoleSchoolAdminName": l4.L4_ROLE_SCHOOL_ADMIN_NAME,
		},
	}
}

// Apply implementa framework.Fixture. Sólo lee + exporta constantes.
// Idempotente.
func (f *l4FullExport) Apply(tx *gorm.DB, ctx *framework.ApplyContext) error {
	// Cuentas declaradas estáticamente vía accessors de l4.*. Si la BD
	// no las refleja, system.ApplySystem (en testdb.StartPostgres) no
	// terminó o falló silenciosamente — fallamos temprano.
	staticResources, err := l4.Resources()
	if err != nil {
		return fmt.Errorf("l4_full_export: l4.Resources(): %w", err)
	}
	staticRoles, err := l4.Roles()
	if err != nil {
		return fmt.Errorf("l4_full_export: l4.Roles(): %w", err)
	}
	staticPerms, err := l4.Permissions()
	if err != nil {
		return fmt.Errorf("l4_full_export: l4.Permissions(): %w", err)
	}
	// P4-1 (plan B): la asignación 1:1 rol×permiso fue eliminada.
	// El conteo de role_grants se verifica a partir de un mínimo
	// declarado por L4 (12 roles × patterns wildcard, ~141 filas).
	staticRS := l4.ResourceScreens
	rs, err := staticRS()
	if err != nil {
		return fmt.Errorf("l4_full_export: l4.ResourceScreens(): %w", err)
	}

	// resources: L0+L3+L4 viven en iam.resources. Validamos
	// >= count(L4) — los de L0/L3 sumarán encima.
	var totalResources int64
	if err := tx.Raw(`SELECT COUNT(*) FROM iam.resources WHERE is_active = true`).Scan(&totalResources).Error; err != nil {
		return fmt.Errorf("l4_full_export: count iam.resources: %w", err)
	}
	if int(totalResources) < len(staticResources) {
		return fmt.Errorf("l4_full_export: iam.resources count=%d < accessor L4 count=%d (system.ApplySystem incompleto?)", totalResources, len(staticResources))
	}

	// roles L4 (5): student, teacher, guardian, admin, school_admin.
	// + super_admin (L0) + announcement_viewer (L1) → 7 totales.
	var totalRoles int64
	if err := tx.Raw(`SELECT COUNT(*) FROM iam.roles WHERE is_active = true`).Scan(&totalRoles).Error; err != nil {
		return fmt.Errorf("l4_full_export: count iam.roles: %w", err)
	}
	if int(totalRoles) < len(staticRoles) {
		return fmt.Errorf("l4_full_export: iam.roles count=%d < accessor L4 count=%d", totalRoles, len(staticRoles))
	}

	var totalPerms int64
	if err := tx.Raw(`SELECT COUNT(*) FROM iam.permissions WHERE is_active = true`).Scan(&totalPerms).Error; err != nil {
		return fmt.Errorf("l4_full_export: count iam.permissions: %w", err)
	}
	if int(totalPerms) < len(staticPerms) {
		return fmt.Errorf("l4_full_export: iam.permissions count=%d < accessor L4 count=%d", totalPerms, len(staticPerms))
	}

	// P4-1 (plan B): se valida que iam.role_grants haya quedado
	// sembrado (al menos 1 fila por rol activo). El modelo nuevo no
	// expone un accessor estático con count fijo — el mínimo se deriva
	// de la cardinalidad de iam.roles activos.
	var totalRoleGrants int64
	if err := tx.Raw(`SELECT COUNT(*) FROM iam.role_grants`).Scan(&totalRoleGrants).Error; err != nil {
		return fmt.Errorf("l4_full_export: count iam.role_grants: %w", err)
	}
	if int(totalRoleGrants) < int(totalRoles) {
		return fmt.Errorf("l4_full_export: iam.role_grants count=%d < total roles=%d (¿applyL4RoleGrants se ejecutó?)", totalRoleGrants, totalRoles)
	}

	var totalRS int64
	if err := tx.Raw(`SELECT COUNT(*) FROM ui_config.resource_screens WHERE is_active = true`).Scan(&totalRS).Error; err != nil {
		return fmt.Errorf("l4_full_export: count ui_config.resource_screens: %w", err)
	}
	if int(totalRS) < len(rs) {
		return fmt.Errorf("l4_full_export: ui_config.resource_screens count=%d < accessor L4 count=%d", totalRS, len(rs))
	}

	// Exportamos sólo strings (Constants es map[string]string).
	ctx.SetConstant("E2EFixtureL4SeedVersion", layers.L4_SEED_VERSION)
	ctx.SetConstant("E2EFixtureL4LayerName", layers.L4_LAYER_NAME)
	ctx.SetConstant("E2EFixtureL4RoleStudentID", l4.L4_ROLE_STUDENT_ID)
	ctx.SetConstant("E2EFixtureL4RoleStudentName", l4.L4_ROLE_STUDENT_NAME)
	ctx.SetConstant("E2EFixtureL4RoleTeacherID", l4.L4_ROLE_TEACHER_ID)
	ctx.SetConstant("E2EFixtureL4RoleTeacherName", l4.L4_ROLE_TEACHER_NAME)
	ctx.SetConstant("E2EFixtureL4RoleGuardianID", l4.L4_ROLE_GUARDIAN_ID)
	ctx.SetConstant("E2EFixtureL4RoleGuardianName", l4.L4_ROLE_GUARDIAN_NAME)
	// PRE-4: E2EFixtureL4RoleAdmin{ID,Name} eliminadas — platform_admin
	// fue removido del catálogo L4.
	ctx.SetConstant("E2EFixtureL4RoleSchoolAdminID", l4.L4_ROLE_SCHOOL_ADMIN_ID)
	ctx.SetConstant("E2EFixtureL4RoleSchoolAdminName", l4.L4_ROLE_SCHOOL_ADMIN_NAME)
	return nil
}

// Cleanup implementa framework.Fixture. La fixture es pasiva (no
// inserta filas — sólo lee + exporta constantes), así que Cleanup es
// un no-op seguro (C-REQ-3.3). Mismo patrón que L0ConstantsExport y
// L3IsolationConstants.
func (f *l4FullExport) Cleanup(tx *gorm.DB, ctx *framework.ApplyContext) error {
	return nil
}
