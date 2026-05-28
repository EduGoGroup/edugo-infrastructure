package common

import (
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// buildUserRole construye un UserRole sin scope (SchoolID y
// AcademicUnitID en nil): el BeforeSave de la entidad calcula
// ScopePattern = "*" para ese caso. ID determinístico
// SHA1(userID:roleID) — mismo patrón que upsertUserRoles en
// focal_evaluacion_v2 y focal_botonera.
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

// SeedUserRole liga un usuario a un rol con ID determinístico. SchoolID
// y AcademicUnitID quedan nil (scope global del rol); el BeforeSave de
// UserRole calcula ScopePattern = "*" automáticamente.
//
// Si en el futuro se necesita scope school/unit, agregar variante
// SeedUserRoleScoped — no abrir el spec ahora (YAGNI).
func SeedUserRole(tx *gorm.DB, userID, roleID uuid.UUID) error {
	ur := buildUserRole(userID, roleID)
	return OnConflictIgnore(tx, &ur)
}
