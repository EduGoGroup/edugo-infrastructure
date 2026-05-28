package framework

import (
	"strings"
	"testing"
)

// Tests "puros" para validar las precondiciones de UpsertBool/String/JSON.
// Los tests integration que necesitan tx se mantienen en
// upsert_integration_test.go bajo build tag `integration`.

func TestUpsertBool_NilTx(t *testing.T) {
	if err := UpsertBool(nil, "t", "id", 1, "col", true); err == nil {
		t.Fatal("se esperaba error nil tx")
	}
}

func TestUpsertBool_EmptyArgs(t *testing.T) {
	if err := UpsertBool(nil, "", "id", 1, "col", true); err == nil {
		t.Fatal("nil tx debería ganar antes de validación de args")
	}
}

func TestUpsertString_NilTx(t *testing.T) {
	if err := UpsertString(nil, "t", "id", 1, "col", "v"); err == nil {
		t.Fatal("se esperaba error nil tx")
	}
}

func TestUpsertJSON_NilTx(t *testing.T) {
	if err := UpsertJSON(nil, "t", "id", 1, "col", []byte(`{}`)); err == nil {
		t.Fatal("se esperaba error nil tx")
	}
}

func TestApplyContext_Provide(t *testing.T) {
	ctx := NewApplyContext("x", "tp", "sp")
	ctx.Provide("school", ProvidedEntity{Kind: "school", ID: "id1", Code: "C1"})
	got, ok := ctx.Provided["school"]
	if !ok {
		t.Fatal("school no quedó registrado")
	}
	if got.ID != "id1" || got.Code != "C1" {
		t.Errorf("ProvidedEntity inesperado: %+v", got)
	}
}

func TestApplyContext_SetConstant(t *testing.T) {
	ctx := NewApplyContext("x", "tp", "sp")
	ctx.SetConstant("Foo", "Bar")
	if ctx.Constants["Foo"] != "Bar" {
		t.Errorf("constants=%v", ctx.Constants)
	}
}

func TestIsSafeIdentifier(t *testing.T) {
	cases := []struct {
		s    string
		safe bool
	}{
		{"academic.schools", true},
		{"iam_roles", true},
		{`"oddName"`, true},
		{"drop table users", false},
		{"a; b", false},
		{"", false},
	}
	for _, tc := range cases {
		if got := isSafeIdentifier(tc.s); got != tc.safe {
			t.Errorf("isSafeIdentifier(%q) = %v, want %v", tc.s, got, tc.safe)
		}
	}
}

func TestFormatPrefixedClause(t *testing.T) {
	got := FormatPrefixedClause("academic.schools", "id", "e2ea1b2c3d4-")
	if !strings.Contains(got, "academic.schools.id") {
		t.Errorf("missing identifier: %q", got)
	}
	if !strings.Contains(got, "e2ea1b2c3d4-") {
		t.Errorf("missing prefix: %q", got)
	}
}

func TestDeleteByPrefix_NilTx(t *testing.T) {
	if _, err := DeleteByPrefix(nil, "t", "id", "p"); err == nil {
		t.Fatal("se esperaba error nil tx")
	}
}
