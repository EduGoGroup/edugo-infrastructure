//go:build integration
// +build integration

package framework_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/fixtures"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/internal/testdb"
)

// Estos tests integration ejercitan helpers del framework contra una
// BD real (testcontainers postgres:15-alpine + production seed) para
// alcanzar cobertura sobre código que sólo es alcanzable con tx≠nil:
// UpsertString, UpsertJSON, Compose con fixtures encadenadas, y los
// branches de DeleteByPrefix.

// TestComposeAssessmentsList aplica role_only + screen_only(assessments-list)
// vía Compose ad-hoc. El branch assessments-list de screen_only crea
// 1 assessment de prueba (entity.Assessment) — distinto al branch
// grades-list que ejercita teacher_grades_only — alimentando la
// cobertura del método privado applyAssessmentsList.
func TestComposeAssessmentsList(t *testing.T) {
	if !testdb.IntegrationGate() {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}
	gdb := testdb.StartPostgres(t)
	c := framework.NewComposer(framework.NewRegistry(), framework.NewNopLogger())

	ctx, err := c.Compose(gdb, "compose_assessments_list", []framework.Fixture{
		&fixtures.RoleOnly{RoleCode: "teacher"},
		&fixtures.ScreenOnly{ScreenKey: "assessments-list"},
	})
	if err != nil {
		t.Fatalf("Compose: %v", err)
	}
	if ctx.Constants["E2EFixtureScreenOnlyAssessmentID"] == "" {
		t.Error("se esperaba E2EFixtureScreenOnlyAssessmentID poblada en ctx.Constants")
	}
}

// TestUpsertHelpersIntegration ejercita los 3 helpers UpsertBool /
// UpsertString / UpsertJSON contra una fila real de academic.schools
// creada por la fixture role_only.
func TestUpsertHelpersIntegration(t *testing.T) {
	if !testdb.IntegrationGate() {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}
	gdb := testdb.StartPostgres(t)
	c := framework.NewComposer(framework.NewRegistry(), framework.NewNopLogger())

	ctx, err := c.Compose(gdb, "upsert_helpers", []framework.Fixture{
		&fixtures.RoleOnly{RoleCode: "school_admin"},
	})
	if err != nil {
		t.Fatalf("Compose role_only: %v", err)
	}
	schoolEntity, ok := ctx.Provided["school"]
	if !ok || schoolEntity.ID == "" {
		t.Fatal("role_only no proveyó school")
	}
	schoolID, err := uuid.Parse(schoolEntity.ID)
	if err != nil {
		t.Fatalf("school ID inválido: %v", err)
	}

	if err := framework.UpsertBool(gdb, (entities.School{}).TableName(), "id", schoolID, "is_active", true); err != nil {
		t.Errorf("UpsertBool: %v", err)
	}
	if err := framework.UpsertString(gdb, (entities.School{}).TableName(), "id", schoolID, "country", "Chile"); err != nil {
		t.Errorf("UpsertString: %v", err)
	}
	if err := framework.UpsertJSON(gdb, (entities.School{}).TableName(), "id", schoolID, "metadata", []byte(`{"upsert":"json"}`)); err != nil {
		t.Errorf("UpsertJSON: %v", err)
	}

	// Camino de error: row inexistente.
	missing := uuid.New()
	if err := framework.UpsertBool(gdb, (entities.School{}).TableName(), "id", missing, "is_active", true); err == nil {
		t.Error("UpsertBool sobre fila inexistente debería fallar")
	}
	if err := framework.UpsertString(gdb, (entities.School{}).TableName(), "id", missing, "country", "X"); err == nil {
		t.Error("UpsertString sobre fila inexistente debería fallar")
	}
	if err := framework.UpsertJSON(gdb, (entities.School{}).TableName(), "id", missing, "metadata", []byte(`{}`)); err == nil {
		t.Error("UpsertJSON sobre fila inexistente debería fallar")
	}
}
