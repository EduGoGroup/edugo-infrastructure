package common

import (
	"testing"

	"github.com/google/uuid"
)

// TestSeedUserRole_Idempotent verifica que llamar SeedUserRole dos veces con
// el mismo (userID, roleID) no duplique la fila.
//
// SKIP: entities.UserRole.TableName() devuelve "iam.user_roles". SQLite trata
// "iam." como un database adjunto inexistente. Adicionalmente, UserRole tiene
// un BeforeSave hook y un uniqueIndex compuesto que exigen soporte Postgres.
// La cobertura de idempotencia viene por composición: SeedUserRole =
// OnConflictIgnore(tx, buildUserRole(userID, roleID)); ambos lados cubiertos
// en db_test.go (OnConflictIgnore) y user_role_test.go (buildUserRole).
func TestSeedUserRole_Idempotent(t *testing.T) {
	t.Skip("AutoMigrate de UserRole no soportado en SQLite: TableName()=\"iam.user_roles\" " +
		"requiere schema Postgres. Idempotencia cubierta por composición: " +
		"SeedUserRole = OnConflictIgnore + buildUserRole, ambos ya testeados.")
	_ = uuid.New()
}

func TestBuildUserRole_DeterministicID(t *testing.T) {
	userID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	roleID := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	ur1 := buildUserRole(userID, roleID)
	ur2 := buildUserRole(userID, roleID)
	if ur1.ID != ur2.ID {
		t.Fatalf("expected ID determinístico, got %v != %v", ur1.ID, ur2.ID)
	}
}

func TestBuildUserRole_DefaultScopeFields(t *testing.T) {
	userID := uuid.New()
	roleID := uuid.New()
	ur := buildUserRole(userID, roleID)
	if ur.SchoolID != nil {
		t.Fatalf("expected SchoolID nil, got %v", *ur.SchoolID)
	}
	if ur.AcademicUnitID != nil {
		t.Fatalf("expected AcademicUnitID nil, got %v", *ur.AcademicUnitID)
	}
	if !ur.IsActive {
		t.Fatal("expected IsActive=true")
	}
	if ur.UserID != userID {
		t.Fatalf("expected UserID preserved, got %v", ur.UserID)
	}
	if ur.RoleID != roleID {
		t.Fatalf("expected RoleID preserved, got %v", ur.RoleID)
	}
	if ur.GrantedAt.IsZero() {
		t.Fatal("expected GrantedAt no zero")
	}
}

func TestBuildUserRole_DifferentInputsProduceDifferentIDs(t *testing.T) {
	userA := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	userB := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
	roleID := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	a := buildUserRole(userA, roleID)
	b := buildUserRole(userB, roleID)
	if a.ID == b.ID {
		t.Fatalf("expected IDs distintos, ambos %v", a.ID)
	}
}
