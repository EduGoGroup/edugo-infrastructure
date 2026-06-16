package fixtures

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/catalog"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l4"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
)

// productionRoleIDs mapea roleCode → UUID del system seed roles (L0..L4).
// Post-Fase-6 (ADR-6/7), los roles se distribuyen entre las capas:
//   - L0 sembra super_admin (UUID legacy 10000000-...-0001 conservado).
//   - L1 sembra announcement_viewer (UUID propio b1000000-...).
//   - L4 sembra los 11 roles del sistema (canónicos + alias) con UUIDs
//     propios b4000000-0001-... (ADR-6 §6: legacy NO se reutiliza).
//
// Las fixtures jamás escriben sobre estos IDs: los referencian en
// user_roles para componer el RBAC del scenario.
// PRE-4: la entry `platform_admin` fue eliminada porque el rol fue
// removido del catálogo L4. Scenarios que necesiten un actor con
// permisos globales deben usar `super_admin`.
var productionRoleIDs = map[string]string{
	"super_admin":         layers.L0_ROLE_SUPER_ADMIN_ID,
	"announcement_viewer": layers.L1_ROLE_ANNOUNCEMENT_VIEWER_ID,
	"school_admin":        l4.L4_ROLE_SCHOOL_ADMIN_ID,
	"school_director":     l4.L4_ROLE_SCHOOL_DIRECTOR_ID,
	"school_coordinator":  l4.L4_ROLE_SCHOOL_COORDINATOR_ID,
	"teacher":             l4.L4_ROLE_TEACHER_ID,
	"assistant_teacher":   l4.L4_ROLE_ASSISTANT_TEACHER_ID,
	"observer":            l4.L4_ROLE_OBSERVER_ID,
	"student":             l4.L4_ROLE_STUDENT_ID,
	"guardian":            l4.L4_ROLE_GUARDIAN_ID,
	"school_assistant":    l4.L4_ROLE_SCHOOL_ASSISTANT_ID,
	"readonly_auditor":    l4.L4_ROLE_READONLY_AUDITOR_ID,
}

// RoleOnly crea el setup mínimo para una prueba focalizada en un solo
// rol del catálogo:
//
//   - 1 academic.schools  (provee "school")
//   - 1 auth.users        (provee "user")
//   - 1 iam.user_roles    (provee "user_role")
//   - 1 academic.memberships (provee "membership")
//
// El RoleCode se resuelve contra los system seed roles (L0..L4) en
// read-only. NO se crean filas en iam.roles ni iam.role_permissions:
// el rol y sus permisos viven en el catálogo sembrado por las capas
// del sistema.
type RoleOnly struct {
	// RoleCode es el rol del catálogo a referenciar (ej. "teacher",
	// "school_admin"). Se valida contra productionRoleIDs en Apply.
	RoleCode string

	// Password en claro para el usuario creado. Se hashea con bcrypt.
	// Si está vacío se usa el default "E2EUser2026!".
	Password string
}

// Manifest implementa framework.Fixture.
func (f *RoleOnly) Manifest() framework.FixtureManifest {
	return framework.FixtureManifest{
		Name:     "role_only",
		Provides: []string{"school", "user", "user_role", "membership"},
		Tables: []string{
			"academic.schools",
			"auth.users",
			"iam.user_roles",
			"academic.memberships",
		},
		Constants: map[string]string{
			"E2EFixtureRoleOnlyRoleCode":     "{{.RoleCode}}",
			"E2EFixtureRoleOnlySchoolCode":   "{{.SchoolCode}}",
			"E2EFixtureRoleOnlyUserEmail":    "{{.UserEmail}}",
			"E2EFixtureRoleOnlyUserPassword": "{{.UserPassword}}",
		},
		Description: "Crea 1 escuela + 1 user + 1 user_role + 1 membership apuntando al roleCode dado.",
	}
}

// Apply implementa framework.Fixture.
func (f *RoleOnly) Apply(tx *gorm.DB, ctx *framework.ApplyContext) error {
	if f.RoleCode == "" {
		return fmt.Errorf("role_only: RoleCode requerido")
	}
	roleUUID, ok := productionRoleIDs[f.RoleCode]
	if !ok {
		return fmt.Errorf("unknown role code: %q (available: %s)", f.RoleCode, availableRoleCodes())
	}

	password := f.Password
	if password == "" {
		password = "E2EUser2026!"
	}

	schoolID := framework.MakeUUID(ctx, "0000-0000-0000-000000000001")
	userID := framework.MakeUUID(ctx, "0000-0000-0000-000000000010")
	userRoleID := framework.MakeUUID(ctx, "0000-0000-0000-000000000020")
	membershipID := framework.MakeUUID(ctx, "0000-0000-0000-000000000030")

	for _, id := range []string{schoolID, userID, userRoleID, membershipID} {
		if err := framework.AssertNotProductionNamespace(id); err != nil {
			return err
		}
	}

	schoolCode := framework.MakeCode(ctx, "SCHOOL", "01")
	userEmail := framework.MakeEmail(ctx, f.RoleCode, "role_only")

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("role_only: bcrypt: %w", err)
	}

	school := entities.School{
		ID:               uuid.MustParse(schoolID),
		Name:             "RoleOnly Test School",
		Code:             schoolCode,
		Country:          "Chile",
		SubscriptionTier: "free",
		MaxTeachers:      10,
		MaxStudents:      50,
		Metadata:         json.RawMessage(`{"e2e":true,"fixture":"role_only"}`),
		IsActive:         true,
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&school).Error; err != nil {
		return fmt.Errorf("role_only: insert school: %w", err)
	}
	// Booleano crítico (F2·H5): asegurar IsActive=true incluso si la
	// fila ya existía con otro estado.
	if err := framework.UpsertBool(tx, school.TableName(), "id", school.ID, "is_active", true); err != nil {
		return err
	}

	user := entities.User{
		ID:           uuid.MustParse(userID),
		Email:        userEmail,
		PasswordHash: string(hashed),
		FirstName:    "E2E",
		LastName:     "RoleOnly",
		IsActive:     true,
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoNothing: true,
	}).Create(&user).Error; err != nil {
		return fmt.Errorf("role_only: insert user: %w", err)
	}
	if err := framework.UpsertBool(tx, user.TableName(), "id", user.ID, "is_active", true); err != nil {
		return err
	}

	schoolUUID := uuid.MustParse(schoolID)
	userRole := entities.UserRole{
		ID:        uuid.MustParse(userRoleID),
		UserID:    user.ID,
		RoleID:    uuid.MustParse(roleUUID),
		SchoolID:  &schoolUUID,
		IsActive:  true,
		GrantedAt: time.Now().UTC(),
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&userRole).Error; err != nil {
		return fmt.Errorf("role_only: insert user_role: %w", err)
	}
	if err := framework.UpsertBool(tx, userRole.TableName(), "id", userRole.ID, "is_active", true); err != nil {
		return err
	}

	invitationTypeID, err := catalog.ResolveInvitationTypeID(tx, membershipRoleFor(f.RoleCode))
	if err != nil {
		return fmt.Errorf("role_only: resolve invitation_type: %w", err)
	}
	membership := entities.Membership{
		ID:               uuid.MustParse(membershipID),
		UserID:           user.ID,
		SchoolID:         schoolUUID,
		InvitationTypeID: invitationTypeID,
		Metadata:         json.RawMessage(`{"e2e":true,"fixture":"role_only"}`),
		Status:           "active",
		EnrolledAt:       time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC),
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&membership).Error; err != nil {
		return fmt.Errorf("role_only: insert membership: %w", err)
	}
	if err := framework.UpsertString(tx, membership.TableName(), "id", membership.ID, "status", "active"); err != nil {
		return err
	}

	ctx.Provide("school", framework.ProvidedEntity{Kind: "school", ID: schoolID, Code: schoolCode})
	ctx.Provide("user", framework.ProvidedEntity{
		Kind:  "user",
		ID:    userID,
		Extra: map[string]string{"email": userEmail, "password": password},
	})
	ctx.Provide("user_role", framework.ProvidedEntity{Kind: "user_role", ID: userRoleID})
	ctx.Provide("membership", framework.ProvidedEntity{Kind: "membership", ID: membershipID})

	ctx.SetConstant("E2EFixtureRoleOnlyRoleCode", f.RoleCode)
	ctx.SetConstant("E2EFixtureRoleOnlySchoolCode", schoolCode)
	ctx.SetConstant("E2EFixtureRoleOnlySchoolID", schoolID)
	ctx.SetConstant("E2EFixtureRoleOnlyUserEmail", userEmail)
	ctx.SetConstant("E2EFixtureRoleOnlyUserPassword", password)
	ctx.SetConstant("E2EFixtureRoleOnlyUserID", userID)
	return nil
}

// Cleanup implementa framework.Fixture. Borra en orden inverso al de
// creación: membership, user_role, user, school. Cada DELETE filtra
// por SchemaPrefix (sólo toca filas de este scenario).
func (f *RoleOnly) Cleanup(tx *gorm.DB, ctx *framework.ApplyContext) error {
	prefix := ctx.SchemaPrefix
	tables := []struct {
		name string
		col  string
	}{
		{"academic.memberships", "id"},
		{"iam.user_roles", "id"},
		{"auth.users", "id"},
		{"academic.schools", "id"},
	}
	for _, t := range tables {
		if _, err := framework.DeleteByPrefix(tx, t.name, t.col, prefix); err != nil {
			return fmt.Errorf("role_only cleanup %s: %w", t.name, err)
		}
	}
	return nil
}

// membershipRoleFor mapea el rol RBAC del catálogo a la key del tipo de
// invitación con que se siembra la membership (MP-08): se resuelve a
// academic.memberships.invitation_type_id vía catalog.ResolveInvitationTypeID.
// Las keys válidas son exactamente:
// ('teacher','student','guardian','coordinator','admin','assistant').
// Los roles "observer-like" (observer, readonly_auditor,
// school_assistant) caen bajo "assistant" — el tipo más cercano a
// "ayudante sin permisos plenos". Para un código desconocido devuelve
// "student" como fallback seguro.
func membershipRoleFor(roleCode string) string {
	switch roleCode {
	// PRE-4: "platform_admin" removido del catálogo. La función
	// debe seguir aceptándolo (legado en fixtures_test.go) y
	// devolver "admin" — el mapeo a memberships.role no requiere
	// que el rol exista en iam.roles.
	case "school_admin", "school_director", "platform_admin", "super_admin":
		return "admin"
	case "school_coordinator":
		return "coordinator"
	case "teacher", "assistant_teacher":
		return "teacher"
	case "guardian":
		return "guardian"
	case "observer", "readonly_auditor", "school_assistant", "announcement_viewer":
		return "assistant"
	default:
		return "student"
	}
}

// availableRoleCodes devuelve los roleCodes conocidos en mensaje de
// error legible.
func availableRoleCodes() string {
	keys := make([]string, 0, len(productionRoleIDs))
	for k := range productionRoleIDs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}
