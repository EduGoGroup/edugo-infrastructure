package common

import (
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// buildUserRole construye un UserRole SIN scope (SchoolID/AcademicUnitID nil):
// el hook BeforeSave de la entidad calcula ScopePattern = "*" para ese caso. El
// id es determinístico SHA1(userID:roleID) — mismo patrón que usaban las 5
// copias upsertUserRole.
func buildUserRole(userID, roleID uuid.UUID) entities.UserRole {
	derived := uuid.NewSHA1(uuid.NameSpaceOID, []byte(userID.String()+":"+roleID.String()))
	return entities.UserRole{
		ID:        derived,
		UserID:    userID,
		RoleID:    roleID,
		IsActive:  true,
		GrantedAt: time.Now().UTC(),
	}
}

// SeedUserRole liga un usuario a un rol L4 con id determinístico. SchoolID y
// AcademicUnitID quedan nil (scope global del rol): el BeforeSave de UserRole
// calcula ScopePattern = "*" automáticamente. Idempotente por id.
//
// Si un playground necesita un id fijo (p. ej. n0n1_escuelas ataba el admin a
// super_admin con un id de constante propio), puede seguir sembrando ese caso a
// mano; el resultado es equivalente porque el id derivado también es estable.
func SeedUserRole(tx *gorm.DB, userID, roleID uuid.UUID) error {
	ur := buildUserRole(userID, roleID)
	return onConflictIgnore(tx, &ur)
}
