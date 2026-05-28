package fixtures

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
)

// GlobalUserNoMembership crea un usuario "global" con rol L0 super_admin
// PERO sin ninguna fila en `academic.memberships`. Esta combinación
// reproduce el bug class detectado en sesión 2026-05-12: un super_admin
// que iniciaba sesión recibía `schools[]` vacío en el login response y
// quedaba bloqueado en SchoolSelector (no podía elegir escuela ni unidad
// porque switch-context exigía membership).
//
// Inserta:
//
//   - 1 fila en `auth.users` (email global-super@e2e.test, password hash).
//   - 1 fila en `iam.user_roles` con role_id = L0 super_admin UUID y
//     school_id = NULL, academic_unit_id = NULL (rol GLOBAL).
//   - NINGUNA fila en `academic.memberships` — clave del test.
//
// Provides: `global_user`. Sin Requires (es self-contained; el rol
// super_admin del catálogo lo siembra L0 vía system.ApplySystem).
//
// Idempotente vía OnConflict DoNothing por id en cada inserción.
type GlobalUserNoMembership struct {
	// Password en claro para el usuario creado. Si está vacío se usa el
	// default "GlobalSuper2026!". Se hashea con bcrypt antes de
	// persistir.
	Password string
}

// Manifest implementa framework.Fixture.
func (f *GlobalUserNoMembership) Manifest() framework.FixtureManifest {
	return framework.FixtureManifest{
		Name:     "global_user_no_membership",
		Provides: []string{"global_user"},
		Tables: []string{
			"auth.users",
			"iam.user_roles",
		},
		Constants: map[string]string{
			"E2EFixtureGlobalUserEmail":    "{{.UserEmail}}",
			"E2EFixtureGlobalUserPassword": "{{.UserPassword}}",
			"E2EFixtureGlobalUserID":       "{{.UserID}}",
		},
		Description: "Crea un user global con rol super_admin (school_id=NULL) y sin membership. Reproduce el flujo de SchoolSelector → switch-context para validar el fix L4 + identity.",
	}
}

// Apply implementa framework.Fixture. Idempotente.
func (f *GlobalUserNoMembership) Apply(tx *gorm.DB, ctx *framework.ApplyContext) error {
	if ctx == nil {
		return fmt.Errorf("global_user_no_membership: nil ApplyContext")
	}

	password := f.Password
	if password == "" {
		password = "GlobalSuper2026!"
	}

	userID := framework.MakeUUID(ctx, "0000-0000-0000-000000000100")
	userRoleID := framework.MakeUUID(ctx, "0000-0000-0000-000000000101")

	for _, id := range []string{userID, userRoleID} {
		if err := framework.AssertNotProductionNamespace(id); err != nil {
			return err
		}
	}

	userEmail := framework.MakeEmail(ctx, "global-super", "no_membership")

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("global_user_no_membership: bcrypt: %w", err)
	}

	userUUID := uuid.MustParse(userID)
	user := entities.User{
		ID:           userUUID,
		Email:        userEmail,
		PasswordHash: string(hashed),
		FirstName:    "Global",
		LastName:     "SuperAdmin",
		IsActive:     true,
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoNothing: true,
	}).Create(&user).Error; err != nil {
		return fmt.Errorf("global_user_no_membership: insert user: %w", err)
	}
	if err := framework.UpsertBool(tx, user.TableName(), "id", user.ID, "is_active", true); err != nil {
		return err
	}

	roleUUID, err := uuid.Parse(layers.L0_ROLE_SUPER_ADMIN_ID)
	if err != nil {
		return fmt.Errorf("global_user_no_membership: parse role uuid: %w", err)
	}

	// SchoolID y AcademicUnitID quedan NULL deliberadamente: este es un
	// rol global. La columna user_roles.school_id IS NULL es la firma
	// que distingue al super_admin sin tenant del super_admin con
	// membership de demo.
	userRole := entities.UserRole{
		ID:        uuid.MustParse(userRoleID),
		UserID:    userUUID,
		RoleID:    roleUUID,
		SchoolID:  nil,
		IsActive:  true,
		GrantedAt: time.Now().UTC(),
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&userRole).Error; err != nil {
		return fmt.Errorf("global_user_no_membership: insert user_role: %w", err)
	}
	if err := framework.UpsertBool(tx, userRole.TableName(), "id", userRole.ID, "is_active", true); err != nil {
		return err
	}

	ctx.Provide("global_user", framework.ProvidedEntity{
		Kind:  "user",
		ID:    userID,
		Extra: map[string]string{"email": userEmail, "password": password},
	})

	ctx.SetConstant("E2EFixtureGlobalUserEmail", userEmail)
	ctx.SetConstant("E2EFixtureGlobalUserPassword", password)
	ctx.SetConstant("E2EFixtureGlobalUserID", userID)
	return nil
}

// Cleanup implementa framework.Fixture. Borra exclusivamente por prefijo
// del scenario (SchemaPrefix) — nunca toca filas del production seed.
func (f *GlobalUserNoMembership) Cleanup(tx *gorm.DB, ctx *framework.ApplyContext) error {
	prefix := ctx.SchemaPrefix
	tables := []struct {
		name string
		col  string
	}{
		{"iam.user_roles", "id"},
		{"auth.users", "id"},
	}
	for _, t := range tables {
		if _, err := framework.DeleteByPrefix(tx, t.name, t.col, prefix); err != nil {
			return fmt.Errorf("global_user_no_membership cleanup %s: %w", t.name, err)
		}
	}
	return nil
}
