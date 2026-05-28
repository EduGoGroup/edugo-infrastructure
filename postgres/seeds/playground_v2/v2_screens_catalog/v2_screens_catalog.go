// Package v2_screens_catalog es el primer playground de la línea v2,
// pensado para iterar el CRUD de los recursos meta del SDUI desde la
// propia app: screen_templates, screen_instances, permissions_mgmt y
// roles.
//
// A diferencia de los playgrounds v1 (focal_pantalla, focal_evaluacion,
// etc.) que sembraban recursos+pantallas ad-hoc encima de L0, este
// paquete asume que el sistema completo (L0..L4) ya corrió. L4 trae
// los 4 recursos meta con sus permisos y pantallas list/form
// configuradas, así que aquí sólo se siembra el envoltorio multi-tenant
// + 3 usuarios con grants distintos para validar la grilla de permisos.
//
// Convive con los playgrounds v1 sin colisionar: rango UUID propio
// 62000000-... y emails con sufijo @v2.edugo.local.
//
// Lo que siembra:
//  1. academic.schools          — 1 escuela "V2 Catalog".
//  2. academic.academic_units   — 1 unidad raíz.
//  3. iam.roles                 — 3 roles scope=school:
//     v2_catalog_admin   (CRUD completo sobre los 4 meta)
//     v2_catalog_viewer  (sólo .read)
//     v2_catalog_author  (.read + .create, sin update/delete)
//  4. iam.role_grants           — patrones específicos a los 4 recursos meta:
//     admin.screen_templates.*  / .read / .read+.create
//     admin.screen_instances.*  / .read / .read+.create
//     admin.permissions_mgmt.*  / .read / .read+.create
//     admin.roles.*             / .read / .read+.create
//     No se usa wildcard global *.read para mantener
//     el menú colapsado a esos 4 (más el padre admin).
//  5. auth.users                — 3 usuarios con password "12345678".
//  6. iam.user_roles            — assignments 1×1.
//  7. academic.memberships      — los 3 en la misma escuela/unidad.
//
// Idempotente: OnConflict DoNothing por id en todas las inserciones.
package v2_screens_catalog

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	AdminEmail  = "v2-admin@v2.edugo.local"
	ViewerEmail = "v2-viewer@v2.edugo.local"
	AuthorEmail = "v2-author@v2.edugo.local"
	Password    = "12345678"

	// Rango UUID 62000000-...: reservado para v2_screens_catalog.
	schoolID     = "62000000-0000-0000-0000-000000000001"
	unitID       = "62000000-0000-0000-0000-000000000002"
	adminUserID  = "62000000-0000-0000-0000-000000000010"
	viewerUserID = "62000000-0000-0000-0000-000000000011"
	authorUserID = "62000000-0000-0000-0000-000000000012"
	adminMembID  = "62000000-0000-0000-0000-000000000020"
	viewerMembID = "62000000-0000-0000-0000-000000000021"
	authorMembID = "62000000-0000-0000-0000-000000000022"

	adminRoleID  = "12000000-0000-0000-0000-000000000001"
	viewerRoleID = "12000000-0000-0000-0000-000000000002"
	authorRoleID = "12000000-0000-0000-0000-000000000003"

	adminRoleName  = "v2_catalog_admin"
	viewerRoleName = "v2_catalog_viewer"
	authorRoleName = "v2_catalog_author"

	schoolCode = "V2-CATALOG"
	schoolName = "Escuela V2 — Catálogo de Pantallas"
	unitCode   = "V2-CATALOG-MAIN"
	unitName   = "Sede Única"

	academicYear = 2026
)

// metaResources lista los 4 recursos meta del SDUI (definidos en L4)
// sobre los que este playground habilita CRUD. El path se usa para
// construir los patterns de role_grants.
var metaResources = []string{
	"admin.screen_templates",
	"admin.screen_instances",
	"admin.permissions_mgmt",
	"admin.roles",
}

// Apply siembra el playground v2_screens_catalog. Asume que L0..L4
// corrieron (los 4 recursos meta + sus screen_instances ya existen).
// Idempotente.
func Apply(tx *gorm.DB) error {
	if err := upsertSchool(tx); err != nil {
		return fmt.Errorf("playground_v2/v2_screens_catalog: school: %w", err)
	}
	if err := upsertAcademicUnit(tx); err != nil {
		return fmt.Errorf("playground_v2/v2_screens_catalog: academic_unit: %w", err)
	}
	if err := upsertRole(tx, adminRoleID, adminRoleName, "Admin V2 — CRUD pantallas/permisos", "CRUD completo sobre los 4 recursos meta del SDUI."); err != nil {
		return fmt.Errorf("playground_v2/v2_screens_catalog: admin_role: %w", err)
	}
	if err := upsertRoleGrants(tx, adminRoleID, grantsFor("*")); err != nil {
		return fmt.Errorf("playground_v2/v2_screens_catalog: admin_grants: %w", err)
	}
	if err := upsertRole(tx, viewerRoleID, viewerRoleName, "Viewer V2 — solo lectura", "Sólo lectura sobre los 4 recursos meta."); err != nil {
		return fmt.Errorf("playground_v2/v2_screens_catalog: viewer_role: %w", err)
	}
	if err := upsertRoleGrants(tx, viewerRoleID, grantsFor("read")); err != nil {
		return fmt.Errorf("playground_v2/v2_screens_catalog: viewer_grants: %w", err)
	}
	if err := upsertRole(tx, authorRoleID, authorRoleName, "Author V2 — lee y crea", "Lee + crea, sin editar ni eliminar."); err != nil {
		return fmt.Errorf("playground_v2/v2_screens_catalog: author_role: %w", err)
	}
	if err := upsertRoleGrants(tx, authorRoleID, append(grantsFor("read"), grantsFor("create")...)); err != nil {
		return fmt.Errorf("playground_v2/v2_screens_catalog: author_grants: %w", err)
	}
	if err := upsertUser(tx, adminUserID, AdminEmail, "Admin", "V2"); err != nil {
		return fmt.Errorf("playground_v2/v2_screens_catalog: admin_user: %w", err)
	}
	if err := upsertUser(tx, viewerUserID, ViewerEmail, "Viewer", "V2"); err != nil {
		return fmt.Errorf("playground_v2/v2_screens_catalog: viewer_user: %w", err)
	}
	if err := upsertUser(tx, authorUserID, AuthorEmail, "Author", "V2"); err != nil {
		return fmt.Errorf("playground_v2/v2_screens_catalog: author_user: %w", err)
	}
	if err := upsertUserRole(tx, adminUserID, adminRoleID); err != nil {
		return fmt.Errorf("playground_v2/v2_screens_catalog: admin_user_role: %w", err)
	}
	if err := upsertUserRole(tx, viewerUserID, viewerRoleID); err != nil {
		return fmt.Errorf("playground_v2/v2_screens_catalog: viewer_user_role: %w", err)
	}
	if err := upsertUserRole(tx, authorUserID, authorRoleID); err != nil {
		return fmt.Errorf("playground_v2/v2_screens_catalog: author_user_role: %w", err)
	}
	if err := upsertMembership(tx, adminMembID, adminUserID, "admin"); err != nil {
		return fmt.Errorf("playground_v2/v2_screens_catalog: admin_membership: %w", err)
	}
	if err := upsertMembership(tx, viewerMembID, viewerUserID, "teacher"); err != nil {
		return fmt.Errorf("playground_v2/v2_screens_catalog: viewer_membership: %w", err)
	}
	if err := upsertMembership(tx, authorMembID, authorUserID, "teacher"); err != nil {
		return fmt.Errorf("playground_v2/v2_screens_catalog: author_membership: %w", err)
	}
	return nil
}

// grantsFor genera la lista de patterns para una acción dada sobre los
// 4 recursos meta. action="*" produce el wildcard CRUD completo.
func grantsFor(action string) []string {
	patterns := make([]string, 0, len(metaResources))
	for _, r := range metaResources {
		patterns = append(patterns, r+"."+action)
	}
	return patterns
}

func upsertSchool(tx *gorm.DB) error {
	id, err := uuid.Parse(schoolID)
	if err != nil {
		return err
	}
	s := entities.School{
		ID:               id,
		Name:             schoolName,
		Code:             schoolCode,
		Country:          "Chile",
		SubscriptionTier: "basic",
		MaxTeachers:      0,
		MaxStudents:      0,
		IsActive:         true,
		Metadata:         json.RawMessage(`{}`),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&s).Error
}

func upsertAcademicUnit(tx *gorm.DB) error {
	id, err := uuid.Parse(unitID)
	if err != nil {
		return err
	}
	sid, err := uuid.Parse(schoolID)
	if err != nil {
		return err
	}
	u := entities.AcademicUnit{
		ID:           id,
		SchoolID:     sid,
		Name:         unitName,
		Code:         unitCode,
		Type:         "school",
		AcademicYear: academicYear,
		Metadata:     json.RawMessage(`{}`),
		IsActive:     true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&u).Error
}

func upsertRole(tx *gorm.DB, idStr, name, displayName, description string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	desc := description
	r := entities.Role{
		ID:          id,
		Name:        name,
		DisplayName: displayName,
		Description: &desc,
		Scope:       "school",
		IsActive:    true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&r).Error
}

// upsertRoleGrants inserta patterns allow para un rol. ID derivado por
// SHA1(role:pattern:effect) alineado a la convención de L4.
func upsertRoleGrants(tx *gorm.DB, roleIDStr string, patterns []string) error {
	rid, err := uuid.Parse(roleIDStr)
	if err != nil {
		return err
	}
	effect := "allow"
	for _, p := range patterns {
		gid := uuid.NewSHA1(uuid.NameSpaceOID, []byte(rid.String()+":"+p+":"+effect))
		g := entities.RoleGrant{
			ID:      gid,
			RoleID:  rid,
			Pattern: p,
			Effect:  effect,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "role_id"}, {Name: "pattern"}, {Name: "effect"}},
			DoNothing: true,
		}).Create(&g).Error; err != nil {
			return err
		}
	}
	return nil
}

func upsertUser(tx *gorm.DB, idStr, email, first, last string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt: %w", err)
	}
	u := entities.User{
		ID:           id,
		Email:        email,
		PasswordHash: string(hash),
		FirstName:    first,
		LastName:     last,
		IsActive:     true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&u).Error
}

func upsertUserRole(tx *gorm.DB, userIDStr, roleIDStr string) error {
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		return err
	}
	rid, err := uuid.Parse(roleIDStr)
	if err != nil {
		return err
	}
	derived := uuid.NewSHA1(uuid.NameSpaceOID, []byte(uid.String()+":"+rid.String()))
	ur := entities.UserRole{
		ID:        derived,
		UserID:    uid,
		RoleID:    rid,
		IsActive:  true,
		GrantedAt: time.Now().UTC(),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&ur).Error
}

func upsertMembership(tx *gorm.DB, idStr, userIDStr, roleKind string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		return err
	}
	sid, err := uuid.Parse(schoolID)
	if err != nil {
		return err
	}
	auid, err := uuid.Parse(unitID)
	if err != nil {
		return err
	}
	m := entities.Membership{
		ID:             id,
		UserID:         uid,
		SchoolID:       sid,
		AcademicUnitID: &auid,
		Role:           roleKind,
		Metadata:       json.RawMessage(`{}`),
		IsActive:       true,
		EnrolledAt:     time.Now().UTC(),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&m).Error
}
