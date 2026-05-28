package layers

import (
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// applyL0Users siembra el usuario super_admin@edugo.system de L0 con
// su ligadura al rol super_admin. Idempotente por id.
//
// SEGURIDAD: el password L0_SUPER_ADMIN_PASSWORD es constante de
// bootstrapping. Rotar antes del primer login productivo en cloud.
//
// NOTA: bcrypt genera un salt aleatorio por invocación, por lo que el
// hash es no-determinístico. La idempotencia se garantiza vía
// OnConflict DoNothing por columna `id`: si el usuario ya existe no
// se reescribe el hash existente.
func applyL0Users(tx *gorm.DB) error {
	if err := upsertL0User(tx); err != nil {
		return err
	}
	if err := upsertL0UserRole(tx); err != nil {
		return err
	}
	return nil
}

func upsertL0User(tx *gorm.DB) error {
	id, err := uuid.Parse(L0_USER_SUPER_ADMIN_ID)
	if err != nil {
		return fmt.Errorf("upsertL0User: parse id: %w", err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(L0_SUPER_ADMIN_PASSWORD), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("upsertL0User: bcrypt: %w", err)
	}
	user := entities.User{
		ID:           id,
		Email:        L0_SUPER_ADMIN_EMAIL,
		PasswordHash: string(hash),
		FirstName:    "Super",
		LastName:     "Admin",
		IsActive:     true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&user).Error
}

func upsertL0UserRole(tx *gorm.DB) error {
	userID, err := uuid.Parse(L0_USER_SUPER_ADMIN_ID)
	if err != nil {
		return fmt.Errorf("upsertL0UserRole: parse user_id: %w", err)
	}
	roleID, err := uuid.Parse(L0_ROLE_SUPER_ADMIN_ID)
	if err != nil {
		return fmt.Errorf("upsertL0UserRole: parse role_id: %w", err)
	}
	// UUID determinístico derivado de (user_id, role_id). Permite
	// idempotencia por `id` aunque la UNIQUE compuesta
	// (user_id, role_id, school_id, academic_unit_id) no dispare
	// conflicto cuando school_id y academic_unit_id son NULL
	// (PostgreSQL trata NULL != NULL en UNIQUE constraints).
	derivedID := uuid.NewSHA1(uuid.NameSpaceOID, []byte(userID.String()+":"+roleID.String()))
	ur := entities.UserRole{
		ID:        derivedID,
		UserID:    userID,
		RoleID:    roleID,
		IsActive:  true,
		GrantedAt: time.Now().UTC(),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&ur).Error
}
