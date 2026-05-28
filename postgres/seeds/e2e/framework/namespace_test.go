package framework

import (
	"strings"
	"testing"
)

func TestDerive_Stable(t *testing.T) {
	tenant1, schema1 := Derive("teacher_grades_only")
	tenant2, schema2 := Derive("teacher_grades_only")
	if tenant1 != tenant2 || schema1 != schema2 {
		t.Fatalf("Derive no es determinista: %q/%q vs %q/%q", tenant1, schema1, tenant2, schema2)
	}
	if !strings.HasPrefix(tenant1, "E2E-") || !strings.HasSuffix(tenant1, "-") {
		t.Errorf("TenantPrefix con forma inesperada: %q", tenant1)
	}
	if !strings.HasPrefix(schema1, "e2e") || !strings.HasSuffix(schema1, "-") {
		t.Errorf("SchemaPrefix con forma inesperada: %q", schema1)
	}
}

func TestDerive_Legacy(t *testing.T) {
	tenant, schema := Derive(LegacyScenarioName)
	if tenant != "E2E-" {
		t.Errorf("legacy TenantPrefix esperaba %q, got %q", "E2E-", tenant)
	}
	if schema != "e2e00000-" {
		t.Errorf("legacy SchemaPrefix esperaba %q, got %q", "e2e00000-", schema)
	}
}

func TestDerive_DifferentNames_DifferentPrefixes(t *testing.T) {
	t1, s1 := Derive("scenario_a")
	t2, s2 := Derive("scenario_b")
	if t1 == t2 || s1 == s2 {
		t.Fatalf("Derive devolvió el mismo prefijo para nombres distintos: %q/%q", t1, s1)
	}
}

func TestAssertNotProductionNamespace(t *testing.T) {
	cases := []struct {
		uuid    string
		wantErr bool
	}{
		{"e2ea1b2c3d4-0000-0000-0000-000000000001", false},
		{"e2e00000-0000-0000-0000-000000000001", false},
		{"10000000-0000-0000-0000-000000000003", true},
		{"c1000000-0000-0000-0000-000000000001", true},
		{"00000000-0000-0000-0000-000000000001", true},
		{"", true},
	}
	for _, tc := range cases {
		err := AssertNotProductionNamespace(tc.uuid)
		if tc.wantErr && err == nil {
			t.Errorf("uuid=%q: se esperaba error, got nil", tc.uuid)
		}
		if !tc.wantErr && err != nil {
			t.Errorf("uuid=%q: error inesperado: %v", tc.uuid, err)
		}
	}
}

func TestMakeUUID(t *testing.T) {
	ctx := NewApplyContext("legacy_e2e", "E2E-", "e2e00000-")
	got := MakeUUID(ctx, "0000-0000-0000-000000000001")
	want := "e2e00000-0000-0000-0000-000000000001"
	if got != want {
		t.Fatalf("MakeUUID got=%q want=%q", got, want)
	}
	if len(got) != 36 {
		t.Errorf("UUID resultante no es canónico (36 chars); got=%d (%q)", len(got), got)
	}
}

// TestMakeUUID_NonLegacyIsCanonical asegura que el hash de scenarios
// no-legacy más el sufijo estándar produzca un UUID textual de exactamente
// 36 caracteres (formato 8-4-4-4-12). Es la regresión que motivó el
// cambio a hash de 5 chars (ver namespace.go).
func TestMakeUUID_NonLegacyIsCanonical(t *testing.T) {
	_, schema := Derive("teacher_grades_only")
	ctx := NewApplyContext("teacher_grades_only", "E2E-", schema)
	got := MakeUUID(ctx, "0000-0000-0000-000000000001")
	if len(got) != 36 {
		t.Fatalf("MakeUUID(non-legacy) no es canónico: got=%d chars %q", len(got), got)
	}
}

func TestMakeUUID_NilContext(t *testing.T) {
	got := MakeUUID(nil, "abc")
	if got != "abc" {
		t.Fatalf("nil ctx debería devolver suffix tal cual; got=%q", got)
	}
}

func TestMakeCode(t *testing.T) {
	ctx := NewApplyContext("teacher_grades_only", "E2E-A1B2C3D4-", "e2ea1b2c3d4-")
	got := MakeCode(ctx, "SCHOOL", "01")
	want := "E2E-A1B2C3D4-SCHOOL-01"
	if got != want {
		t.Fatalf("MakeCode got=%q want=%q", got, want)
	}
}

func TestMakeEmail(t *testing.T) {
	ctx := NewApplyContext("teacher_grades_only", "E2E-A1B2C3D4-", "e2ea1b2c3d4-")
	got := MakeEmail(ctx, "teacher", "role_only")
	if !strings.HasSuffix(got, "@edugo.test") {
		t.Errorf("MakeEmail debe terminar en @edugo.test; got=%q", got)
	}
	if !strings.Contains(got, "teacher-role_only-") {
		t.Errorf("MakeEmail debe contener 'teacher-role_only-'; got=%q", got)
	}
	hash := schemaHashFromPrefix(ctx)
	if !strings.Contains(got, hash) {
		t.Errorf("MakeEmail debe contener el hash del scenario (%q); got=%q", hash, got)
	}
}

func TestMakeEmail_NilContext(t *testing.T) {
	got := MakeEmail(nil, "teacher", "role_only")
	if !strings.Contains(got, LegacyHash) {
		t.Errorf("nil ctx debe usar LegacyHash; got=%q", got)
	}
}
