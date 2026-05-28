package layers

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// applyL1User siembra el usuario viewer@edugo.demo de L1.
//
// SEGURIDAD: el password L1_VIEWER_PASSWORD es constante de
// bootstrapping. Rotar antes del primer login productivo en cloud.
//
// NOTA: bcrypt genera un salt aleatorio por invocación, por lo que el
// hash es no-determinístico. La idempotencia se garantiza vía
// OnConflict DoNothing por columna `id`: si el usuario ya existe no
// se reescribe el hash existente. Replica el patrón de upsertL0User.
func applyL1User(tx *gorm.DB) error {
	id, err := uuid.Parse(L1_USER_VIEWER_ID)
	if err != nil {
		return fmt.Errorf("applyL1User: parse id: %w", err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(L1_VIEWER_PASSWORD), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("applyL1User: bcrypt: %w", err)
	}
	user := entities.User{
		ID:           id,
		Email:        L1_VIEWER_EMAIL,
		PasswordHash: string(hash),
		FirstName:    "Viewer",
		LastName:     "Demo",
		IsActive:     true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&user).Error
}
