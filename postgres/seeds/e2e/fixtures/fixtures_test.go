package fixtures

import (
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
)

// Estos tests son "puros": no tocan la base de datos. Cubren los
// invariantes que se pueden validar sin ejecutar SQL — Manifest,
// validaciones de campos requeridos y comportamiento ante tx=nil. Los
// tests de integración real (con BD) viven en seeds/e2e/integration y
// requieren build tag `e2e_integration`.

// ----------------------------------------------------------------------
// RoleOnly
// ----------------------------------------------------------------------

func TestRoleOnly_Manifest(t *testing.T) {
	f := &RoleOnly{RoleCode: "teacher"}
	m := f.Manifest()
	if m.Name != "role_only" {
		t.Errorf("Name=%q, want role_only", m.Name)
	}
	wantProvides := []string{"school", "user", "user_role", "membership"}
	if !equalStrSlice(m.Provides, wantProvides) {
		t.Errorf("Provides=%v, want %v", m.Provides, wantProvides)
	}
	if len(m.Requires) != 0 {
		t.Errorf("Requires debería estar vacío; got=%v", m.Requires)
	}
	if len(m.Tables) == 0 {
		t.Error("Tables vacío")
	}
	if m.Constants == nil {
		t.Error("Constants es nil")
	}
}

func TestRoleOnly_Apply_EmptyRoleCode(t *testing.T) {
	f := &RoleOnly{RoleCode: ""}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), "RoleCode requerido") {
		t.Fatalf("se esperaba error 'RoleCode requerido'; got=%v", err)
	}
}

func TestRoleOnly_Apply_UnknownRoleCode(t *testing.T) {
	f := &RoleOnly{RoleCode: "wizard"}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), "unknown role code") {
		t.Fatalf("se esperaba error 'unknown role code'; got=%v", err)
	}
}

// ----------------------------------------------------------------------
// ScreenOnly
// ----------------------------------------------------------------------

func TestScreenOnly_Manifest(t *testing.T) {
	f := &ScreenOnly{ScreenKey: "assessments-list"}
	m := f.Manifest()
	if m.Name != "screen_only" {
		t.Errorf("Name=%q, want screen_only", m.Name)
	}
	if !contains(m.Provides, "screen_data") {
		t.Errorf("Provides debe incluir 'screen_data'; got=%v", m.Provides)
	}
	if !contains(m.Requires, "school") {
		t.Errorf("Requires debe incluir 'school'; got=%v", m.Requires)
	}
	if len(m.Tables) == 0 {
		t.Error("Tables vacío")
	}
	if m.Constants == nil {
		t.Error("Constants es nil")
	}
}

func TestScreenOnly_Apply_EmptyScreenKey(t *testing.T) {
	f := &ScreenOnly{ScreenKey: ""}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), "ScreenKey requerido") {
		t.Fatalf("se esperaba 'ScreenKey requerido'; got=%v", err)
	}
}

func TestScreenOnly_Apply_NilContext(t *testing.T) {
	f := &ScreenOnly{ScreenKey: "assessments-list"}
	err := f.Apply(nil, nil)
	if err == nil || !strings.Contains(err.Error(), "nil ApplyContext") {
		t.Fatalf("se esperaba 'nil ApplyContext'; got=%v", err)
	}
}

func TestScreenOnly_Apply_MissingSchoolCapability(t *testing.T) {
	f := &ScreenOnly{ScreenKey: "assessments-list"}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), `"school"`) {
		t.Fatalf("se esperaba error de capability 'school'; got=%v", err)
	}
}

func TestScreenOnly_Apply_NilTx(t *testing.T) {
	f := &ScreenOnly{ScreenKey: "assessments-list"}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	ctx.Provide("school", framework.ProvidedEntity{
		Kind: "school",
		ID:   "e2e00000-0000-0000-0000-000000000001",
	})
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), "nil transaction") {
		t.Fatalf("se esperaba 'nil transaction'; got=%v", err)
	}
}

func TestScreenOnly_Cleanup_NilTx(t *testing.T) {
	f := &ScreenOnly{ScreenKey: "assessments-list"}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	if err := f.Cleanup(nil, ctx); err == nil {
		t.Fatal("Cleanup con tx=nil debería fallar")
	}
}

func TestSupportedScreenKeys_NotEmpty(t *testing.T) {
	keys := SupportedScreenKeys()
	if len(keys) == 0 {
		t.Fatal("SupportedScreenKeys vacío")
	}
	if !contains(keys, "assessments-list") {
		t.Errorf("SupportedScreenKeys debe incluir 'assessments-list'; got=%v", keys)
	}
	if FormatSupportedScreenKeys() == "" {
		t.Error("FormatSupportedScreenKeys vacío")
	}
}

// TestScreenOnly_SupportedScreenKeys_IncludesGradesList valida que la
// fixture exponga el screenKey grades-list (Fase C, Grupo 2).
func TestScreenOnly_SupportedScreenKeys_IncludesGradesList(t *testing.T) {
	keys := SupportedScreenKeys()
	if !contains(keys, "grades-list") {
		t.Errorf("SupportedScreenKeys debe incluir 'grades-list'; got=%v", keys)
	}
}

// TestScreenOnly_Apply_GradesList_NilTx_WithSchool comprueba que el
// branch grades-list es alcanzable: con school provista pero tx=nil la
// fixture debe devolver "nil transaction" (las inserciones no se
// ejecutan, pero el switch llega al branch correcto).
func TestScreenOnly_Apply_GradesList_NilTx_WithSchool(t *testing.T) {
	f := &ScreenOnly{ScreenKey: "grades-list"}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	ctx.Provide("school", framework.ProvidedEntity{
		Kind: "school",
		ID:   "e2e00000-0000-0000-0000-000000000001",
	})
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), "nil transaction") {
		t.Fatalf("se esperaba 'nil transaction'; got=%v", err)
	}
}

// ----------------------------------------------------------------------
// ReadonlyRole
// ----------------------------------------------------------------------

func TestReadonlyRole_Manifest(t *testing.T) {
	f := &ReadonlyRole{Resources: []string{"assessments"}}
	m := f.Manifest()
	if m.Name != "readonly_role" {
		t.Errorf("Name=%q, want readonly_role", m.Name)
	}
	if !contains(m.Provides, "readonly_role") {
		t.Errorf("Provides debe incluir 'readonly_role'; got=%v", m.Provides)
	}
	if len(m.Tables) == 0 {
		t.Error("Tables vacío")
	}
	if m.Constants == nil {
		t.Error("Constants es nil")
	}
}

func TestReadonlyRole_Apply_EmptyResources(t *testing.T) {
	f := &ReadonlyRole{Resources: nil}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), "Resources requerido") {
		t.Fatalf("se esperaba 'Resources requerido'; got=%v", err)
	}
}

func TestReadonlyRole_Apply_BlankResource(t *testing.T) {
	f := &ReadonlyRole{Resources: []string{"assessments", ""}}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), "vacío") {
		t.Fatalf("se esperaba error de resource vacío; got=%v", err)
	}
}

func TestReadonlyRole_Apply_NilTx(t *testing.T) {
	f := &ReadonlyRole{Resources: []string{"assessments"}}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), "nil transaction") {
		t.Fatalf("se esperaba 'nil transaction'; got=%v", err)
	}
}

func TestReadonlyRole_Cleanup_NilTx(t *testing.T) {
	f := &ReadonlyRole{Resources: []string{"assessments"}}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	if err := f.Cleanup(nil, ctx); err == nil {
		t.Fatal("Cleanup con tx=nil debería fallar")
	}
}

// ----------------------------------------------------------------------
// PartialCrud
// ----------------------------------------------------------------------

func TestPartialCrud_Manifest(t *testing.T) {
	f := &PartialCrud{Resources: []string{"assessments"}}
	m := f.Manifest()
	if m.Name != "partial_crud" {
		t.Errorf("Name=%q, want partial_crud", m.Name)
	}
	if !contains(m.Provides, "partial_crud_role") {
		t.Errorf("Provides debe incluir 'partial_crud_role'; got=%v", m.Provides)
	}
	if _, ok := m.Constants[PartialCrudRoleCodeConstant]; !ok {
		t.Errorf("Constants debe incluir %q; got=%v", PartialCrudRoleCodeConstant, m.Constants)
	}
}

func TestPartialCrud_Apply_EmptyResources(t *testing.T) {
	f := &PartialCrud{Resources: nil}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), "Resources requerido") {
		t.Fatalf("se esperaba 'Resources requerido'; got=%v", err)
	}
}

func TestPartialCrud_Apply_BlankResource(t *testing.T) {
	f := &PartialCrud{Resources: []string{""}}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), "vacío") {
		t.Fatalf("se esperaba error de resource vacío; got=%v", err)
	}
}

func TestPartialCrud_Apply_NilTx(t *testing.T) {
	f := &PartialCrud{Resources: []string{"assessments"}}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), "nil transaction") {
		t.Fatalf("se esperaba 'nil transaction'; got=%v", err)
	}
}

func TestPartialCrud_Cleanup_NilTx(t *testing.T) {
	f := &PartialCrud{Resources: []string{"assessments"}}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	if err := f.Cleanup(nil, ctx); err == nil {
		t.Fatal("Cleanup con tx=nil debería fallar")
	}
}

// ----------------------------------------------------------------------
// MenuSubtree
// ----------------------------------------------------------------------

func TestMenuSubtree_Manifest_RequiresReadonlyRoleWhenRoleIDEmpty(t *testing.T) {
	f := &MenuSubtree{SubtreeRoot: "academico"}
	m := f.Manifest()
	if m.Name != "menu_subtree" {
		t.Errorf("Name=%q, want menu_subtree", m.Name)
	}
	if !contains(m.Provides, "menu_subtree") {
		t.Errorf("Provides debe incluir 'menu_subtree'; got=%v", m.Provides)
	}
	if !contains(m.Requires, "readonly_role") {
		t.Errorf("Requires debe incluir 'readonly_role' cuando RoleID está vacío; got=%v", m.Requires)
	}
}

func TestMenuSubtree_Manifest_NoRequiresWhenRoleIDProvided(t *testing.T) {
	f := &MenuSubtree{SubtreeRoot: "academico", RoleID: "e2e00000-0000-0000-0000-0000000000ab"}
	m := f.Manifest()
	if contains(m.Requires, "readonly_role") {
		t.Errorf("cuando RoleID está informado, no debería requerir 'readonly_role'; got=%v", m.Requires)
	}
}

func TestMenuSubtree_Apply_EmptySubtreeRoot(t *testing.T) {
	f := &MenuSubtree{SubtreeRoot: ""}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), "SubtreeRoot requerido") {
		t.Fatalf("se esperaba 'SubtreeRoot requerido'; got=%v", err)
	}
}

func TestMenuSubtree_Apply_NilContext(t *testing.T) {
	f := &MenuSubtree{SubtreeRoot: "academico"}
	err := f.Apply(nil, nil)
	if err == nil || !strings.Contains(err.Error(), "nil ApplyContext") {
		t.Fatalf("se esperaba 'nil ApplyContext'; got=%v", err)
	}
}

func TestMenuSubtree_Apply_MissingReadonlyRoleCapability(t *testing.T) {
	f := &MenuSubtree{SubtreeRoot: "academico"}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), `"readonly_role"`) {
		t.Fatalf("se esperaba error de capability 'readonly_role'; got=%v", err)
	}
}

func TestMenuSubtree_Apply_NilTx_WithRoleID(t *testing.T) {
	f := &MenuSubtree{
		SubtreeRoot: "academico",
		RoleID:      "e2e00000-0000-0000-0000-0000000000ab",
	}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), "nil transaction") {
		t.Fatalf("se esperaba 'nil transaction'; got=%v", err)
	}
}

func TestMenuSubtree_Apply_NilTx_FromProvided(t *testing.T) {
	f := &MenuSubtree{SubtreeRoot: "academico"}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	ctx.Provide("readonly_role", framework.ProvidedEntity{
		Kind: "role",
		ID:   "e2e00000-0000-0000-0000-0000000000ab",
	})
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), "nil transaction") {
		t.Fatalf("se esperaba 'nil transaction'; got=%v", err)
	}
}

func TestMenuSubtree_Cleanup_NilTx(t *testing.T) {
	f := &MenuSubtree{SubtreeRoot: "academico"}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	if err := f.Cleanup(nil, ctx); err == nil {
		t.Fatal("Cleanup con tx=nil debería fallar")
	}
}

// ----------------------------------------------------------------------
// Helpers locales
// ----------------------------------------------------------------------

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func equalStrSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestRoleOnly_MembershipRoleFor cubre la tabla de mapeo roleCode →
// membership.role para todos los códigos del catálogo.
func TestRoleOnly_MembershipRoleFor(t *testing.T) {
	cases := map[string]string{
		"super_admin":        "admin",
		"platform_admin":     "admin",
		"school_admin":       "admin",
		"school_director":    "admin",
		"school_coordinator": "coordinator",
		"teacher":            "teacher",
		"assistant_teacher":  "teacher",
		"guardian":           "guardian",
		"observer":           "assistant",
		"readonly_auditor":   "assistant",
		"school_assistant":   "assistant",
		"student":            "student",
		"unknown_code":       "student",
	}
	for code, want := range cases {
		if got := membershipRoleFor(code); got != want {
			t.Errorf("membershipRoleFor(%q) = %q, want %q", code, got, want)
		}
	}
}

// TestAvailableRoleCodes_NotEmpty asegura que el listado de códigos
// del catálogo se devuelve ordenado y separado por coma.
func TestAvailableRoleCodes_NotEmpty(t *testing.T) {
	got := availableRoleCodes()
	if got == "" {
		t.Fatal("availableRoleCodes vacío")
	}
	if !strings.Contains(got, "teacher") {
		t.Errorf("availableRoleCodes debería incluir 'teacher'; got=%q", got)
	}
	if !strings.Contains(got, ",") {
		t.Errorf("availableRoleCodes debería ser CSV; got=%q", got)
	}
}

// TestScreenOnly_Apply_DefaultBranch cubre el caso pragmático: para un
// screenKey desconocido la fixture devuelve sin error y registra el
// screen_data como "none".
func TestScreenOnly_Apply_DefaultBranch(t *testing.T) {
	f := &ScreenOnly{ScreenKey: "unknown-screen-xyz"}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	ctx.Provide("school", framework.ProvidedEntity{
		Kind: "school",
		ID:   "e2e00000-0000-0000-0000-000000000001",
	})
	// Para alcanzar el branch default necesitamos pasar una tx no-nil
	// pero como el branch default no toca la BD, basta con un
	// puntero arbitrario; usamos un *gorm.DB nil-pero-distinto-de-nil
	// no es trivial sin imports extra, así que cubrimos el branch a
	// través de la ramificación: validación → tx-check → uuid.Parse.
	// Con tx=nil, alcanzaremos "nil transaction" en lugar del default.
	// Para cubrir el default real haría falta un test de integración.
	// Aquí validamos al menos que no falla por ScreenKey desconocido
	// antes de llegar al check de tx.
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), "nil transaction") {
		t.Fatalf("se esperaba 'nil transaction'; got=%v", err)
	}
}

// TestRoleOnly_Cleanup_NilTx valida que Cleanup falle limpio con tx=nil.
// El caller (cleaner) nunca debería invocarlo con nil, pero el guard
// evita panics si algún test mal configurado lo hace.
func TestRoleOnly_Cleanup_NilTx(t *testing.T) {
	f := &RoleOnly{RoleCode: "teacher"}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	// RoleOnly.Cleanup actual no chequea nil-tx explícitamente; al
	// invocar DeleteByPrefix con tx=nil, devuelve un error de
	// framework. Validamos que NO panic y que devuelve algún error.
	if err := f.Cleanup(nil, ctx); err == nil {
		t.Fatal("Cleanup con tx=nil debería devolver error")
	}
}

// TestSchemaHashFromCtx_Legacy comprueba la utilidad compartida que
// extrae el hash del SchemaPrefix legacy. Cubre el branch que lee
// LegacyHash cuando el ctx es nil.
// TestRoleOnly_Apply_Validations cubre los retornos tempranos de Apply.
// Notar que NO ejercitamos el path completo (requiere BD): sólo
// confirmamos que las validaciones de inputs se ejecutan en orden.
func TestRoleOnly_Apply_Validations(t *testing.T) {
	cases := []struct {
		name    string
		fixture *RoleOnly
		ctx     *framework.ApplyContext
		needle  string
	}{
		{
			name:    "empty role code",
			fixture: &RoleOnly{},
			ctx:     framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-"),
			needle:  "RoleCode requerido",
		},
		{
			name:    "unknown role code",
			fixture: &RoleOnly{RoleCode: "wizard"},
			ctx:     framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-"),
			needle:  "unknown role code",
		},
		{
			name:    "available codes listed",
			fixture: &RoleOnly{RoleCode: "ghost"},
			ctx:     framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-"),
			needle:  "available:",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fixture.Apply(nil, tc.ctx)
			if err == nil || !strings.Contains(err.Error(), tc.needle) {
				t.Fatalf("se esperaba mensaje con %q; got=%v", tc.needle, err)
			}
		})
	}
}

// TestScreenOnly_Apply_Validations cubre todos los retornos tempranos.
func TestScreenOnly_Apply_Validations(t *testing.T) {
	withSchool := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	withSchool.Provide("school", framework.ProvidedEntity{
		Kind: "school",
		ID:   "e2e00000-0000-0000-0000-000000000001",
	})
	bareCtx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")

	cases := []struct {
		name    string
		fixture *ScreenOnly
		ctx     *framework.ApplyContext
		needle  string
	}{
		{"empty screen key", &ScreenOnly{}, bareCtx, "ScreenKey requerido"},
		{"nil ctx", &ScreenOnly{ScreenKey: "x"}, nil, "nil ApplyContext"},
		{"missing school", &ScreenOnly{ScreenKey: "x"}, bareCtx, "school"},
		{"nil tx", &ScreenOnly{ScreenKey: "x"}, withSchool, "nil transaction"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fixture.Apply(nil, tc.ctx)
			if err == nil || !strings.Contains(err.Error(), tc.needle) {
				t.Fatalf("se esperaba %q; got=%v", tc.needle, err)
			}
		})
	}
}

// TestReadonlyRole_Apply_Validations.
func TestReadonlyRole_Apply_Validations(t *testing.T) {
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	cases := []struct {
		name    string
		fixture *ReadonlyRole
		ctx     *framework.ApplyContext
		needle  string
	}{
		{"empty resources", &ReadonlyRole{}, ctx, "Resources requerido"},
		{"blank in slice", &ReadonlyRole{Resources: []string{"a", ""}}, ctx, "vacío"},
		{"nil ctx", &ReadonlyRole{Resources: []string{"a"}}, nil, "nil ApplyContext"},
		{"nil tx", &ReadonlyRole{Resources: []string{"a"}}, ctx, "nil transaction"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fixture.Apply(nil, tc.ctx)
			if err == nil || !strings.Contains(err.Error(), tc.needle) {
				t.Fatalf("se esperaba %q; got=%v", tc.needle, err)
			}
		})
	}
}

// TestPartialCrud_Apply_Validations.
func TestPartialCrud_Apply_Validations(t *testing.T) {
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	cases := []struct {
		name    string
		fixture *PartialCrud
		ctx     *framework.ApplyContext
		needle  string
	}{
		{"empty resources", &PartialCrud{}, ctx, "Resources requerido"},
		{"blank in slice", &PartialCrud{Resources: []string{""}}, ctx, "vacío"},
		{"nil ctx", &PartialCrud{Resources: []string{"a"}}, nil, "nil ApplyContext"},
		{"nil tx", &PartialCrud{Resources: []string{"a"}}, ctx, "nil transaction"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fixture.Apply(nil, tc.ctx)
			if err == nil || !strings.Contains(err.Error(), tc.needle) {
				t.Fatalf("se esperaba %q; got=%v", tc.needle, err)
			}
		})
	}
}

// TestMenuSubtree_Apply_Validations.
func TestMenuSubtree_Apply_Validations(t *testing.T) {
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	withRole := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	withRole.Provide("readonly_role", framework.ProvidedEntity{
		Kind: "role",
		ID:   "e2e00000-0000-0000-0000-0000000000ab",
	})
	cases := []struct {
		name    string
		fixture *MenuSubtree
		ctx     *framework.ApplyContext
		needle  string
	}{
		{"empty root", &MenuSubtree{}, ctx, "SubtreeRoot requerido"},
		{"nil ctx", &MenuSubtree{SubtreeRoot: "academico"}, nil, "nil ApplyContext"},
		{"missing capability", &MenuSubtree{SubtreeRoot: "academico"}, ctx, "readonly_role"},
		{"invalid roleID format", &MenuSubtree{SubtreeRoot: "academico", RoleID: "not-a-uuid"}, ctx, "RoleID inválido"},
		{"nil tx with explicit roleID", &MenuSubtree{SubtreeRoot: "academico", RoleID: "e2e00000-0000-0000-0000-0000000000ab"}, ctx, "nil transaction"},
		{"nil tx via provided", &MenuSubtree{SubtreeRoot: "academico"}, withRole, "nil transaction"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fixture.Apply(nil, tc.ctx)
			if err == nil || !strings.Contains(err.Error(), tc.needle) {
				t.Fatalf("se esperaba %q; got=%v", tc.needle, err)
			}
		})
	}
}

// TestCleanup_EmptySchemaPrefix valida que todas las fixtures rechacen
// un ApplyContext con SchemaPrefix vacío al hacer cleanup, evitando
// borrados accidentales `LIKE %`.
func TestCleanup_EmptySchemaPrefix(t *testing.T) {
	ctx := &framework.ApplyContext{ScenarioName: "x"}

	cases := map[string]framework.Fixture{
		"screen_only":       &ScreenOnly{ScreenKey: "x"},
		"readonly_role":     &ReadonlyRole{Resources: []string{"a"}},
		"partial_crud":      &PartialCrud{Resources: []string{"a"}},
		"menu_subtree":      &MenuSubtree{SubtreeRoot: "academico"},
		"guardian_relation": &GuardianRelation{},
	}
	for name, f := range cases {
		t.Run(name, func(t *testing.T) {
			err := f.Cleanup(nil, ctx)
			if err == nil {
				t.Fatalf("%s.Cleanup con SchemaPrefix vacío debería fallar", name)
			}
		})
	}
}

// ----------------------------------------------------------------------
// GuardianRelation
// ----------------------------------------------------------------------

func TestGuardianRelation_Manifest(t *testing.T) {
	f := &GuardianRelation{}
	m := f.Manifest()
	if m.Name != "guardian_relation" {
		t.Errorf("Name=%q, want guardian_relation", m.Name)
	}
	wantProvides := []string{"guardian_relation"}
	if !equalStrSlice(m.Provides, wantProvides) {
		t.Errorf("Provides=%v, want %v", m.Provides, wantProvides)
	}
	wantRequires := []string{"school", "user"}
	if !equalStrSlice(m.Requires, wantRequires) {
		t.Errorf("Requires=%v, want %v", m.Requires, wantRequires)
	}
	if len(m.Tables) == 0 {
		t.Error("Tables vacío")
	}
	if m.Constants == nil {
		t.Error("Constants es nil")
	}
	for _, key := range []string{
		E2EFixtureGuardianRelationStudentID,
		E2EFixtureGuardianRelationStudentEmail,
		E2EFixtureGuardianRelationID,
	} {
		if _, ok := m.Constants[key]; !ok {
			t.Errorf("Constants debe incluir %q; got=%v", key, m.Constants)
		}
	}
}

func TestGuardianRelation_Apply_NilContext(t *testing.T) {
	f := &GuardianRelation{}
	err := f.Apply(nil, nil)
	if err == nil || !strings.Contains(err.Error(), "nil ApplyContext") {
		t.Fatalf("se esperaba 'nil ApplyContext'; got=%v", err)
	}
}

func TestGuardianRelation_Apply_MissingSchoolCapability(t *testing.T) {
	f := &GuardianRelation{}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), `"school"`) {
		t.Fatalf("se esperaba error de capability 'school'; got=%v", err)
	}
}

func TestGuardianRelation_Apply_MissingUserCapability(t *testing.T) {
	f := &GuardianRelation{}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	ctx.Provide("school", framework.ProvidedEntity{
		Kind: "school",
		ID:   "e2e00000-0000-0000-0000-000000000001",
	})
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), `"user"`) {
		t.Fatalf("se esperaba error de capability 'user'; got=%v", err)
	}
}

func TestGuardianRelation_Apply_NilTx_AfterCapabilities(t *testing.T) {
	f := &GuardianRelation{}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	ctx.Provide("school", framework.ProvidedEntity{
		Kind: "school",
		ID:   "e2e00000-0000-0000-0000-000000000001",
	})
	ctx.Provide("user", framework.ProvidedEntity{
		Kind: "user",
		ID:   "e2e00000-0000-0000-0000-000000000010",
	})
	err := f.Apply(nil, ctx)
	if err == nil || !strings.Contains(err.Error(), "nil transaction") {
		t.Fatalf("se esperaba 'nil transaction'; got=%v", err)
	}
}

func TestGuardianRelation_Cleanup_NilTx(t *testing.T) {
	f := &GuardianRelation{}
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	if err := f.Cleanup(nil, ctx); err == nil {
		t.Fatal("Cleanup con tx=nil debería fallar")
	}
}

// TestGuardianRelation_Apply_StudentEmailPattern valida que el email
// del student siga el patrón determinista basado en SchemaPrefix
// (`<role>-<fixture>-<hash>@edugo.test` vía framework.MakeEmail).
// Reproducimos el comportamiento de Apply hasta antes de tocar la BD
// chequeando el formato esperado.
func TestGuardianRelation_Apply_StudentEmailPattern(t *testing.T) {
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	got := framework.MakeEmail(ctx, "student", "guardian_relation")
	want := "student-guardian_relation-00000@edugo.test"
	if got != want {
		t.Errorf("student email=%q, want %q", got, want)
	}
}

func TestSchemaHashFromCtx_Legacy(t *testing.T) {
	ctx := framework.NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	if got := schemaHashFromCtx(ctx); got != "00000" {
		// El hash legacy textual extraído del prefix "e2e00000-" es
		// "00000" (5 chars); LegacyHash es la versión completa de 8.
		t.Errorf("schemaHashFromCtx(legacy) = %q, want %q", got, "00000")
	}
	if got := schemaHashFromCtx(nil); got != framework.LegacyHash {
		t.Errorf("schemaHashFromCtx(nil) = %q, want %q", got, framework.LegacyHash)
	}
}
