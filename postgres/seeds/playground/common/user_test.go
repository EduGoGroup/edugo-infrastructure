package common

import (
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// TestSeedUser_Idempotent verifica que llamar SeedUser dos veces con el mismo
// UserSpec no duplique la fila.
//
// SKIP: entities.User.TableName() devuelve "auth.users". SQLite trata "auth."
// como un database adjunto inexistente, por lo que GORM no puede crear la
// tabla ni insertar.
// La cobertura de idempotencia viene por composición: SeedUser =
// OnConflictIgnore(tx, buildUser(spec)); ambos lados cubiertos en
// db_test.go (OnConflictIgnore) y user_test.go (buildUser).
func TestSeedUser_Idempotent(t *testing.T) {
	t.Skip("AutoMigrate de User no soportado en SQLite: TableName()=\"auth.users\" " +
		"requiere schema Postgres. Idempotencia cubierta por composición: " +
		"SeedUser = OnConflictIgnore + buildUser, ambos ya testeados.")
	_ = uuid.New()
	_ = bcrypt.MinCost // evitar import vacío
}

func TestBuildUser_HashesPassword(t *testing.T) {
	id := uuid.New()
	u := buildUser(UserSpec{
		ID:        id,
		Email:     "test@edugo.local",
		Password:  "12345678",
		FirstName: "Test",
		LastName:  "User",
	})
	if u.PasswordHash == "" {
		t.Fatal("expected PasswordHash no vacío")
	}
	if u.PasswordHash == "12345678" {
		t.Fatal("PasswordHash NO debe ser el plaintext")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("12345678")); err != nil {
		t.Fatalf("bcrypt.CompareHashAndPassword falló: %v", err)
	}
}

func TestBuildUser_PreservesFields(t *testing.T) {
	id := uuid.New()
	u := buildUser(UserSpec{
		ID:        id,
		Email:     "alice@edugo.local",
		Password:  "12345678",
		FirstName: "Alice",
		LastName:  "Lastname",
	})
	if u.ID != id {
		t.Fatalf("expected ID preserved, got %v", u.ID)
	}
	if u.Email != "alice@edugo.local" {
		t.Fatalf("expected Email preserved, got %q", u.Email)
	}
	if u.FirstName != "Alice" || u.LastName != "Lastname" {
		t.Fatalf("expected names preserved, got %q %q", u.FirstName, u.LastName)
	}
	if !u.IsActive {
		t.Fatal("expected IsActive=true")
	}
	// TokenVersion no se setea: la BD lo provee con default=1.
	if u.TokenVersion != 0 {
		t.Fatalf("expected TokenVersion NO seteado (=0), got %d", u.TokenVersion)
	}
}
