//go:build integration
// +build integration

package scenarios_test

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"

	"gorm.io/gorm"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/internal/testdb"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/scenarios"
)

// TestScenarioL4Full_Integration cierra el bloque E de Fase 6 — valida
// la capa L4 completa contra una BD efímera con L0..L4 aplicadas
// (testdb.StartPostgres corre system.ApplySystem sin filtro de capa).
// El test:
//
//  1. Aplica el scenario `l4_full` vía composer.Apply (el fixture
//     `l4_full_export` chequea cuentas y exporta constantes).
//
//  2. Valida idempotencia (segundo Apply no falla — C-REQ-1.4).
//
//  3. Genera la matriz `roleKey → []resource_keys accesibles` y
//     `roleKey → []screen_keys accesibles` PROGRAMÁTICAMENTE desde:
//
//     - iam.role_permissions JOIN iam.permissions p ON … →
//     permission_name → resource_key (prefijo antes del primer `:`).
//     - iam.resources WHERE is_menu_visible=true AND is_active=true
//     filtrados por las resource_keys, propagando visibilidad a
//     ancestros (mismo algoritmo que
//     core/usecase/menu/get_user_menu.go::propagateVisibility).
//     - ui_config.resource_screens JOIN iam.resources →
//     screen_keys accesibles por rol.
//
//  4. Para los 3 roles representativos (super_admin, teacher, student):
//
//     - GET /menu (simulado): asserts que el set de resource_keys del
//     menú es exactamente el esperado por la matriz computada.
//     - screen-config/resolve/key/:key (simulado): para CADA screen
//     accesible del rol, verifica que ui_config.screen_instances
//     tiene una fila por screen_key y que su slot_data parsea como
//     JSON (mismo flujo del use-case ResolveByKeyUseCase). Cualquier
//     screen_key del menú sin instancia en ui_config rompe el test
//     con un drift report (rol, screen_key, estado).
//     - Conteo total de screens accesibles coincide con la matriz
//     documentada en este comentario (valor calculado, no
//     hardcodeado — la matriz se imprime al inicio del test para
//     futura comparación).
//
//  5. Sanity checks de aislamiento por rol:
//     - cada rol tiene > 0 screens (acceso mínimo al menú).
//     - super_admin NO debe tener TODOS los screens del sistema —
//     L0 sólo le da announcements:*. Si hereda screens de L4
//     accidentalmente, hay drift en el seed.
//     - teacher y student tienen subsets distintos: simétrica !=
//     identidad (sanidad contra collapse del modelo RBAC).
//     - student ⊆ teacher en assessments/grades NO se asume — son
//     subsets independientes con shape común sólo en
//     dashboard/announcements/screens/menu.
//
// Ref F6-REQ-6.x y F6-REQ-7.x (`phase-6-layer-l4/{requirements,
// design}.md`).
//
// Ejecución:
//
//	ENABLE_INTEGRATION_TESTS=true go test -tags=integration \
//	    -run TestScenarioL4Full -count=1 -timeout=300s \
//	    ./seeds/e2e/scenarios/...
func TestScenarioL4Full_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("scenario_l4_full: skip en modo -short")
	}
	if !testdb.IntegrationGate() {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}

	gdb := testdb.StartPostgres(t)

	// 1. Aplicar el scenario `l4_full`. system.ApplySystem ya corrió en
	//    testdb.StartPostgres aplicando L0..L4; la fixture
	//    l4_full_export verifica cuentas y exporta constantes L4.
	reg := framework.NewRegistry()
	scenario := &scenarios.L4Full{}
	if err := reg.RegisterScenario(scenario); err != nil {
		t.Fatalf("RegisterScenario: %v", err)
	}
	composer := framework.NewComposer(reg, framework.NewNopLogger())

	ctx, err := composer.Apply(gdb, "l4_full")
	if err != nil {
		t.Fatalf("composer.Apply(l4_full): %v", err)
	}
	if ctx.ScenarioName != "l4_full" {
		t.Errorf("ScenarioName=%q, want %q", ctx.ScenarioName, "l4_full")
	}
	if len(ctx.Constants) == 0 {
		t.Error("ctx.Constants vacío tras Apply — la fixture l4_full_export no llamó SetConstant")
	}

	// 2. Idempotencia (C-REQ-1.4) — Apply 2ª vez no debe fallar.
	if _, err := composer.Apply(gdb, "l4_full"); err != nil {
		t.Fatalf("Apply 2 (idempotencia rota): %v", err)
	}

	// 3. Construir la matriz expected PROGRAMÁTICAMENTE.
	matrix := computeAccessMatrix(t, gdb)

	// Imprimir la matriz como ayuda al lector del log de CI cuando algo
	// rompe — no es output del test, sólo contexto.
	t.Logf("matriz screens-por-rol (calculada desde BD):")
	roleNames := make([]string, 0, len(matrix))
	for k := range matrix {
		roleNames = append(roleNames, k)
	}
	sort.Strings(roleNames)
	for _, role := range roleNames {
		t.Logf("  %s: %d screens accesibles (%d resources)",
			role, len(matrix[role].Screens), len(matrix[role].Resources))
	}

	// Roles representativos a validar (F6-REQ-7.1: 3 roles).
	representative := []string{"super_admin", "teacher", "student"}
	for _, role := range representative {
		if _, ok := matrix[role]; !ok {
			t.Fatalf("rol representativo %q ausente de la matriz computada — drift fatal en el seed", role)
		}
	}

	// 4. Para cada rol representativo, validar GET /menu y resolve/key.
	for _, role := range representative {
		role := role
		t.Run(fmt.Sprintf("F6-REQ-7.1_role_%s_menu_and_screens", role), func(t *testing.T) {
			access := matrix[role]
			// 4a. GET /menu simulado: cada resource_key del menú accesible
			//     debe existir en iam.resources con is_menu_visible=true
			//     y is_active=true. La construcción de access.Resources ya
			//     filtró por esos predicados — aquí re-validamos como
			//     defensa en profundidad.
			for _, rkey := range access.Resources {
				var n int64
				if err := gdb.Raw(
					`SELECT COUNT(*) FROM iam.resources
					 WHERE key = ? AND is_menu_visible = true AND is_active = true`,
					rkey,
				).Scan(&n).Error; err != nil {
					t.Fatalf("role=%s: query iam.resources[%s]: %v", role, rkey, err)
				}
				if n != 1 {
					t.Errorf("role=%s: resource %q no visible/activo (count=%d) — drift menu", role, rkey, n)
				}
			}

			// 4b. Por cada screen_key accesible, resolve/key simulado:
			//     leer ui_config.screen_instances + parse JSON de
			//     slot_data. Cualquier screen_key sin instance, o con
			//     slot_data inválido, es drift.
			for _, skey := range access.Screens {
				var slotJSON, templateID, instID string
				err := gdb.Raw(
					`SELECT id::text, COALESCE(slot_data::text, '{}'), template_id::text
					 FROM ui_config.screen_instances
					 WHERE screen_key = ?`,
					skey,
				).Row().Scan(&instID, &slotJSON, &templateID)
				if err != nil {
					t.Errorf("role=%s screen_key=%s: ui_config.screen_instances no devuelve fila — drift resolve/key: %v", role, skey, err)
					continue
				}
				var sd map[string]any
				if jerr := json.Unmarshal([]byte(slotJSON), &sd); jerr != nil {
					t.Errorf("role=%s screen_key=%s instance=%s: slot_data no parsea JSON — drift contract: %v", role, skey, instID, jerr)
					continue
				}
				// Defensa en profundidad: el template_id debe existir.
				var hasTpl int64
				if err := gdb.Raw(
					`SELECT COUNT(*) FROM ui_config.screen_templates WHERE id = ?::uuid`,
					templateID,
				).Scan(&hasTpl).Error; err != nil {
					t.Errorf("role=%s screen_key=%s: query template: %v", role, skey, err)
					continue
				}
				if hasTpl != 1 {
					t.Errorf("role=%s screen_key=%s instance=%s: template_id=%s ausente — drift FK", role, skey, instID, templateID)
				}
			}

			// 4c. El conteo total de screens accesibles debe coincidir
			//     con lo computado en `matrix[role].Screens`. Trivial por
			//     construcción, pero el assert deja constancia explícita
			//     del número en el log de CI.
			t.Logf("role=%s: %d screens accesibles", role, len(access.Screens))
			if len(access.Screens) == 0 {
				t.Errorf("role=%s: 0 screens accesibles — el rol no puede operar (¿BD vacía?)", role)
			}
		})
	}

	// 5. Sanity checks de aislamiento por rol.
	t.Run("F6-REQ-6.1_super_admin_minimal_inheritance", func(t *testing.T) {
		// super_admin recibe `announcements:*` en L0 y `materials:*` en
		// L3 (L3 siembra explícitamente super_admin × materials para
		// preservar el flujo administrativo). NO debe heredar otros
		// recursos del catálogo L4 (que sólo van a los 5 roles nuevos).
		// El set esperado es exactamente {announcements, materials}.
		sa := matrix["super_admin"]
		allowed := map[string]bool{"announcements": true, "materials": true}
		extras := []string{}
		for _, rkey := range sa.Resources {
			if !allowed[rkey] {
				extras = append(extras, rkey)
			}
		}
		if len(extras) != 0 {
			t.Errorf("super_admin tiene resources fuera de {announcements, materials}: %v — drift de herencia L4", extras)
		}
		// Defensa positiva: si falta alguno de los 2 esperados, también
		// hay drift (los seeds previos no se ejecutaron).
		got := setOf(sa.Resources)
		for k := range allowed {
			if !got[k] {
				t.Errorf("super_admin no ve resource %q — drift L0/L3", k)
			}
		}
	})

	t.Run("F6-REQ-6.2_teacher_student_distinct_sets", func(t *testing.T) {
		// teacher y student deben tener subsets de screens distintos
		// (sanidad contra collapse del modelo RBAC). El número exacto se
		// imprime al inicio del test; aquí sólo validamos
		// no-identidad.
		//
		// Nota: la matriz computada agrupa screens por recurso, no por
		// acción. teacher y student comparten varias screens (ej.
		// dashboard, announcements-list, materials-list,
		// assessments-list/form) porque ambos ven el recurso. La
		// diferencia real vive en las ACTIONS habilitadas dentro de
		// cada slot_data (los botones SAVE/DELETE/PUBLISH se ocultan
		// según permission_name del usuario, no según resource_screens).
		// Por eso aquí solo aserto !=, no subset estricto.
		ts := setOf(matrix["teacher"].Screens)
		ss := setOf(matrix["student"].Screens)
		if equalSets(ts, ss) {
			t.Error("teacher y student tienen sets de screens idénticos — colapso del modelo RBAC")
		}
		// Asimetría positiva: teacher tiene recursos administrativos
		// (memberships/periods/users/subjects/units/stats/reports) que
		// student NO tiene — el filtrado de acciones dentro de slot_data
		// es lo que diferencia el comportamiento real, pero la presencia
		// de estos screens ya debería ser mayor en teacher.
		if len(ts) <= len(ss) {
			t.Errorf("teacher tiene %d screens, student %d — esperado teacher > student", len(ts), len(ss))
		}
		// Poda F2 (plan 004-permisologia-mvp): report-card se retiró del
		// MVP, por lo que el recurso `reports` ya no resuelve a ninguna
		// pantalla y el student NO debe verlo. Antes de la poda esta
		// aserción era `!ss["report-card"]` (presencia obligatoria);
		// ahora valida la ausencia.
		if ss["report-card"] {
			t.Error("student ve `report-card` — debió retirarse en la poda F2 (plan 004)")
		}
	})

	t.Run("F6-REQ-6.3_total_role_count_in_db", func(t *testing.T) {
		// 12 roles totales (post PRE-4 — platform_admin eliminado):
		//   4 L4 canon (student, teacher, guardian, school_admin)
		// + 6 L4 alias (school_director, school_coordinator,
		//     school_assistant, assistant_teacher, observer,
		//     readonly_auditor)
		// + 1 L0 (super_admin)
		// + 1 L1 (announcement_viewer)
		var n int64
		if err := gdb.Raw(`SELECT COUNT(*) FROM iam.roles WHERE is_active = true`).Scan(&n).Error; err != nil {
			t.Fatalf("count iam.roles: %v", err)
		}
		const wantRoles = 12
		if n != wantRoles {
			t.Errorf("iam.roles activos=%d, want %d (4 L4 canon + 6 L4 alias + 1 L0 + 1 L1)", n, wantRoles)
		}
	})

	t.Run("F6-REQ-6.4_constants_exported_to_context", func(t *testing.T) {
		// Constantes mínimas que tests Kotlin (KMP) consumen vía
		// fixtures-constants.json. Si alguna falta tras Apply, el JSON
		// quedaría incompleto y los tests KMP romperían.
		// PRE-4: E2EFixtureL4RoleAdminName removida — platform_admin
		// eliminado del catálogo L4.
		wantKeys := []string{
			"E2EFixtureL4SeedVersion",
			"E2EFixtureL4LayerName",
			"E2EFixtureL4RoleStudentName",
			"E2EFixtureL4RoleTeacherName",
			"E2EFixtureL4RoleGuardianName",
			"E2EFixtureL4RoleSchoolAdminName",
		}
		for _, k := range wantKeys {
			if v, ok := ctx.Constants[k]; !ok || strings.TrimSpace(v) == "" {
				t.Errorf("constante %q ausente o vacía en ctx.Constants — la fixture no exportó el inventario L4", k)
			}
		}
	})
}

// -------------------------------------------------------------------
// Helpers — matriz screens-por-rol calculada PROGRAMÁTICAMENTE desde
// la BD (NO hardcodeada). E.2 del entregable F6.
// -------------------------------------------------------------------

// roleAccess captura el set de resources/screens visibles para un rol.
// Las listas vienen ordenadas alfabéticamente para que el log de CI sea
// determinístico y los diffs entre corridas no marquen ruido.
type roleAccess struct {
	Resources []string // resource_keys visibles en el menú (con padres propagados)
	Screens   []string // screen_keys resolvibles via screen-config/resolve/key
}

// resourceRow es una proyección mínima de iam.resources usada para
// reconstruir el árbol del menú y propagar visibilidad de hijos a
// padres (igual algoritmo que core/usecase/menu/get_user_menu.go).
type resourceRow struct {
	ID            string
	Key           string
	ParentID      string
	IsMenuVisible bool
	IsActive      bool
}

// computeAccessMatrix produce el mapa role_name → roleAccess derivado
// de la BD. Es la fuente de verdad del test: si esta matriz cambia
// entre corridas para el mismo seed, ha habido drift.
//
// Implementación:
//
//  1. Cargar TODOS los iam.resources activos.
//  2. Cargar TODAS las filas (role_name, permission_name) de
//     iam.role_permissions JOIN iam.permissions p ON
//     rp.permission_id=p.id JOIN iam.roles ro ON rp.role_id=ro.id.
//  3. Cargar TODOS los ui_config.resource_screens activos (resource_key
//     → screen_key).
//  4. Para cada rol:
//     - Extraer resource_keys del rol = { permission_name antes del `:` }.
//     - Filtrar resources que coinciden con esos keys AND
//     is_menu_visible=true AND is_active=true.
//     - Propagar visibilidad a ancestros (parent_id) — sin esto, los
//     contenedores como `admin`/`academic`/`content` no aparecerían
//     aunque sus hijos sí.
//     - Para cada resource visible: agregar todos sus screen_keys de
//     resource_screens al set.
func computeAccessMatrix(t *testing.T, gdb *gorm.DB) map[string]*roleAccess {
	t.Helper()

	// 1. resources.
	var resourceRows []resourceRow
	rows, err := gdb.Raw(
		`SELECT id::text, key, COALESCE(parent_id::text, ''), is_menu_visible, is_active
		 FROM iam.resources
		 WHERE is_active = true`,
	).Rows()
	if err != nil {
		t.Fatalf("computeAccessMatrix: query iam.resources: %v", err)
	}
	for rows.Next() {
		var r resourceRow
		if err := rows.Scan(&r.ID, &r.Key, &r.ParentID, &r.IsMenuVisible, &r.IsActive); err != nil {
			rows.Close()
			t.Fatalf("computeAccessMatrix: scan iam.resources: %v", err)
		}
		resourceRows = append(resourceRows, r)
	}
	rows.Close()
	resByID := make(map[string]*resourceRow, len(resourceRows))
	resByKey := make(map[string]*resourceRow, len(resourceRows))
	for i := range resourceRows {
		r := &resourceRows[i]
		resByID[r.ID] = r
		resByKey[r.Key] = r
	}

	// 2. role → []permission_name. Filtramos por roles activos y
	//    permisos activos.
	type rolePermRow struct {
		RoleName string
		PermName string
	}
	var rpRows []rolePermRow
	rows, err = gdb.Raw(
		`SELECT ro.name, p.name
		 FROM iam.role_permissions rp
		   JOIN iam.permissions p ON rp.permission_id = p.id
		   JOIN iam.roles       ro ON rp.role_id      = ro.id
		 WHERE ro.is_active = true AND p.is_active = true`,
	).Rows()
	if err != nil {
		t.Fatalf("computeAccessMatrix: query role_permissions: %v", err)
	}
	for rows.Next() {
		var rp rolePermRow
		if err := rows.Scan(&rp.RoleName, &rp.PermName); err != nil {
			rows.Close()
			t.Fatalf("computeAccessMatrix: scan role_permissions: %v", err)
		}
		rpRows = append(rpRows, rp)
	}
	rows.Close()

	// permsByRole[role_name] = set(permission_name)
	permsByRole := make(map[string]map[string]bool)
	for _, rp := range rpRows {
		if permsByRole[rp.RoleName] == nil {
			permsByRole[rp.RoleName] = make(map[string]bool)
		}
		permsByRole[rp.RoleName][rp.PermName] = true
	}

	// 3. resource_screens activos.
	type resScrRow struct {
		ResourceKey string
		ScreenKey   string
	}
	var rsRows []resScrRow
	rows, err = gdb.Raw(
		`SELECT resource_key, screen_key
		 FROM ui_config.resource_screens
		 WHERE is_active = true`,
	).Rows()
	if err != nil {
		t.Fatalf("computeAccessMatrix: query resource_screens: %v", err)
	}
	for rows.Next() {
		var r resScrRow
		if err := rows.Scan(&r.ResourceKey, &r.ScreenKey); err != nil {
			rows.Close()
			t.Fatalf("computeAccessMatrix: scan resource_screens: %v", err)
		}
		rsRows = append(rsRows, r)
	}
	rows.Close()
	screensByResource := make(map[string][]string)
	for _, r := range rsRows {
		screensByResource[r.ResourceKey] = append(screensByResource[r.ResourceKey], r.ScreenKey)
	}

	// 4. construir la matriz por rol.
	matrix := make(map[string]*roleAccess, len(permsByRole))
	for roleName, perms := range permsByRole {
		// resource_keys = prefijo de permission_name antes del primer `:`.
		visibleKeys := make(map[string]bool)
		for permName := range perms {
			idx := strings.IndexByte(permName, ':')
			if idx <= 0 {
				continue
			}
			visibleKeys[permName[:idx]] = true
		}

		// Propagar a ancestros — sin esto los contenedores
		// `admin`/`academic`/`content`/`reports` no aparecen pese a
		// tener hijos visibles (algoritmo idéntico al del use-case
		// menu/get_user_menu.go).
		toPropagate := make([]string, 0, len(visibleKeys))
		for k := range visibleKeys {
			toPropagate = append(toPropagate, k)
		}
		for _, k := range toPropagate {
			r, ok := resByKey[k]
			if !ok {
				continue
			}
			propagateVisibility(r.ParentID, resByID, visibleKeys)
		}

		// Filtrar a is_menu_visible.
		resources := make([]string, 0, len(visibleKeys))
		for key := range visibleKeys {
			r, ok := resByKey[key]
			if !ok || !r.IsMenuVisible || !r.IsActive {
				continue
			}
			resources = append(resources, key)
		}
		sort.Strings(resources)

		// screens accesibles = union de screensByResource[rkey] para
		// cada rkey en `resources`. Deduplicar.
		screenSet := make(map[string]bool)
		for _, rkey := range resources {
			for _, sk := range screensByResource[rkey] {
				screenSet[sk] = true
			}
		}
		screens := make([]string, 0, len(screenSet))
		for sk := range screenSet {
			screens = append(screens, sk)
		}
		sort.Strings(screens)

		matrix[roleName] = &roleAccess{
			Resources: resources,
			Screens:   screens,
		}
	}
	return matrix
}

// propagateVisibility recorre la cadena de parent_id marcando ancestros
// como visibles. Es la traducción literal de
// core/usecase/menu/get_user_menu.go::propagateVisibility.
func propagateVisibility(parentID string, resByID map[string]*resourceRow, visibleKeys map[string]bool) {
	if parentID == "" {
		return
	}
	parent, ok := resByID[parentID]
	if !ok {
		return
	}
	if visibleKeys[parent.Key] {
		return
	}
	visibleKeys[parent.Key] = true
	propagateVisibility(parent.ParentID, resByID, visibleKeys)
}

func setOf(xs []string) map[string]bool {
	out := make(map[string]bool, len(xs))
	for _, x := range xs {
		out[x] = true
	}
	return out
}

func equalSets(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}
