package framework

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestConstantsExporter_WriteThenRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fixtures-constants.json")
	exp := NewConstantsExporter(path)
	exp.Now = func() time.Time { return time.Unix(1700000000, 0).UTC() }

	ctx := NewApplyContext("teacher_grades_only", "E2E-A1B2C-", "e2ea1b2c-")
	ctx.SetConstant("E2EFixtureRoleOnlyUserEmail", "teacher-role_only-a1b2c@edugo.test")
	if err := exp.WriteFromContext(ctx); err != nil {
		t.Fatalf("WriteFromContext: %v", err)
	}

	got, err := exp.Read()
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got.SchemaVersion != ConstantsExportSchemaVersion {
		t.Errorf("schemaVersion=%q, want %q", got.SchemaVersion, ConstantsExportSchemaVersion)
	}
	sc, ok := got.Scenarios["teacher_grades_only"]
	if !ok {
		t.Fatalf("scenario teacher_grades_only no presente: %+v", got.Scenarios)
	}
	if sc.TenantPrefix != "E2E-A1B2C-" {
		t.Errorf("TenantPrefix=%q", sc.TenantPrefix)
	}
	if sc.Constants["E2EFixtureRoleOnlyUserEmail"] != "teacher-role_only-a1b2c@edugo.test" {
		t.Errorf("Constants leído mal: %v", sc.Constants)
	}
}

func TestConstantsExporter_MergeMultipleScenarios(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fixtures-constants.json")
	exp := NewConstantsExporter(path)

	ctx1 := NewApplyContext("a", "E2E-AAAAA-", "e2eaaaaa-")
	ctx1.SetConstant("k1", "v1")
	if err := exp.WriteFromContext(ctx1); err != nil {
		t.Fatalf("write a: %v", err)
	}
	ctx2 := NewApplyContext("b", "E2E-BBBBB-", "e2ebbbbb-")
	ctx2.SetConstant("k2", "v2")
	if err := exp.WriteFromContext(ctx2); err != nil {
		t.Fatalf("write b: %v", err)
	}

	got, err := exp.Read()
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(got.Scenarios) != 2 {
		t.Fatalf("se esperaban 2 scenarios, got=%d (%v)", len(got.Scenarios), got.Scenarios)
	}
	if got.Scenarios["a"].Constants["k1"] != "v1" {
		t.Errorf("scenario a.k1 mal: %v", got.Scenarios["a"])
	}
	if got.Scenarios["b"].Constants["k2"] != "v2" {
		t.Errorf("scenario b.k2 mal: %v", got.Scenarios["b"])
	}
}

func TestConstantsExporter_OverwritesScenario(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fixtures-constants.json")
	exp := NewConstantsExporter(path)

	ctx1 := NewApplyContext("a", "E2E-AAAAA-", "e2eaaaaa-")
	ctx1.SetConstant("k", "v1")
	if err := exp.WriteFromContext(ctx1); err != nil {
		t.Fatalf("write 1: %v", err)
	}
	ctx2 := NewApplyContext("a", "E2E-AAAAA-", "e2eaaaaa-")
	ctx2.SetConstant("k", "v2")
	if err := exp.WriteFromContext(ctx2); err != nil {
		t.Fatalf("write 2: %v", err)
	}

	got, _ := exp.Read()
	if got.Scenarios["a"].Constants["k"] != "v2" {
		t.Errorf("se esperaba sobrescritura; got=%v", got.Scenarios["a"].Constants)
	}
}

func TestConstantsExporter_MissingFileIsEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "no-existe.json")
	exp := NewConstantsExporter(path)
	got, err := exp.Read()
	if err != nil {
		t.Fatalf("Read sobre archivo inexistente debería ser no-op: %v", err)
	}
	if len(got.Scenarios) != 0 {
		t.Errorf("scenarios debería ser 0; got=%d", len(got.Scenarios))
	}
}

func TestConstantsExporter_NilContext(t *testing.T) {
	exp := NewConstantsExporter("/tmp/nope.json")
	if err := exp.WriteFromContext(nil); err == nil {
		t.Fatal("se esperaba error con ctx nil")
	}
}

func TestConstantsExporter_EmptyScenarioName(t *testing.T) {
	exp := NewConstantsExporter("/tmp/nope.json")
	ctx := NewApplyContext("", "", "")
	if err := exp.WriteFromContext(ctx); err == nil {
		t.Fatal("se esperaba error con ScenarioName vacío")
	}
}

func TestConstantsExporter_FormatSerializable(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "f.json")
	exp := NewConstantsExporter(path)
	ctx := NewApplyContext("a", "E2E-AAAAA-", "e2eaaaaa-")
	ctx.SetConstant("z", "last")
	ctx.SetConstant("a", "first")
	if err := exp.WriteFromContext(ctx); err != nil {
		t.Fatalf("write: %v", err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read raw: %v", err)
	}
	var roundtrip ConstantsExport
	if err := json.Unmarshal(raw, &roundtrip); err != nil {
		t.Fatalf("roundtrip: %v", err)
	}
	if !strings.Contains(string(raw), `"schemaVersion": "1"`) {
		t.Errorf("schemaVersion no aparece en JSON: %s", raw)
	}
}

func TestConstantsExport_DefaultPath(t *testing.T) {
	exp := NewConstantsExporter("")
	if exp.Path != DefaultExportPath {
		t.Errorf("Path default esperaba %q, got %q", DefaultExportPath, exp.Path)
	}
}

func TestSortedKeys(t *testing.T) {
	keys := SortedKeys(map[string]string{"b": "1", "a": "2", "c": "3"})
	want := []string{"a", "b", "c"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("SortedKeys[%d] = %q, want %q", i, k, want[i])
		}
	}
}
