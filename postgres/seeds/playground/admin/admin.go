// Package admin define el playground más simple: un único usuario
// `admin@edugo.local` con password `12345678`, ligado al rol L0 super_admin
// (sembrado por system/layers/l0_roles.go) y con grant wildcard `*`
// (acceso total). Incluye también una escuela + unidad académica +
// membership para que el usuario pueda seleccionar contexto en la UI
// multi-escuela / multi-unidad.
//
// L0 sólo crea el rol super_admin; este playground agrega:
//  1. auth.users           — usuario admin.
//  2. iam.user_roles       — assignment user × super_admin.
//  3. iam.role_grants      — pattern `*` para super_admin (lo sembraría L4,
//     pero como sólo aplicamos hasta L0 lo sembramos aquí).
//  4. academic.schools     — 1 escuela Playground.
//  5. academic.academic_units — 1 unidad raíz.
//  6. academic.memberships — admin en la escuela con role=admin.
//
// Idempotente: usa OnConflict DoNothing por id en todas las inserciones.
package admin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	UserEmail    = "admin@edugo.local"
	UserPassword = "12345678"

	// Segundo usuario: no-global, una sola escuela + una sola unidad. Sirve
	// para validar el caso opuesto al admin: el dropdown debe ocultar
	// "Cambiar escuela" y "Cambiar unidad" porque solo hay una opción.
	ViewerEmail = "viewer@edugo.local"

	userID       = "60000000-0000-0000-0000-000000000001"
	schoolID     = "60000000-0000-0000-0000-000000000002"
	unitID       = "60000000-0000-0000-0000-000000000003"
	membershipID = "60000000-0000-0000-0000-000000000004"

	// Escuela A extra: segunda unidad para probar "Cambiar unidad".
	unitAAnnexID = "60000000-0000-0000-0000-000000000005"

	// Escuela B: segunda escuela con una sola unidad — para probar "Cambiar
	// escuela" + chain a "Cambiar unidad" cuando aplica.
	schoolBID     = "60000000-0000-0000-0000-000000000006"
	unitBID       = "60000000-0000-0000-0000-000000000007"
	membershipBID = "60000000-0000-0000-0000-000000000008"

	// Viewer playground: usuario no-global con rol playground_viewer y
	// membership a la escuela Norte (1 unidad).
	viewerUserID       = "60000000-0000-0000-0000-000000000009"
	viewerRoleID       = "10000000-0000-0000-0000-000000000002"
	viewerRoleName     = "playground_viewer"
	viewerMembershipID = "60000000-0000-0000-0000-00000000000a"
	// Dos grants:
	// - "announcements.read" literal: lo entiende el matcher del menú
	//   (patternTouchesResource solo soporta literales y `prefix.*`,
	//   no `*.suffix`). Hace que announcements aparezca en el menú en
	//   modo "view".
	// - "*.read" amplio: cubre el resto de endpoints de lectura que la
	//   UI llama al arrancar (notifications, etc.) sin enumerar.
	viewerReadPatternMenu = "announcements.read"
	viewerReadPatternAll  = "*.read"

	schoolCode = "PG-ADMIN"
	schoolName = "Escuela Playground Admin"
	unitCode   = "PG-ADMIN-MAIN"
	unitName   = "Sede Principal"

	unitAAnnexCode = "PG-ADMIN-ANEXO"
	unitAAnnexName = "Sede Anexo"

	schoolBCode = "PG-NORTE"
	schoolBName = "Escuela Playground Norte"
	unitBCode   = "PG-NORTE-MAIN"
	unitBName   = "Sede Norte"

	academicYear = 2026
)

// Apply siembra el playground admin. Asume que L0 ya corrió (super_admin
// existe en iam.roles).
func Apply(tx *gorm.DB) error {
	if err := upsertAdminUser(tx); err != nil {
		return fmt.Errorf("playground/admin: user: %w", err)
	}
	if err := upsertAdminUserRole(tx); err != nil {
		return fmt.Errorf("playground/admin: user_role: %w", err)
	}
	if err := upsertSuperAdminWildcardGrant(tx); err != nil {
		return fmt.Errorf("playground/admin: role_grant: %w", err)
	}
	if err := upsertSchool(tx); err != nil {
		return fmt.Errorf("playground/admin: school: %w", err)
	}
	if err := upsertAcademicUnit(tx); err != nil {
		return fmt.Errorf("playground/admin: academic_unit: %w", err)
	}
	if err := upsertAcademicUnitAnnex(tx); err != nil {
		return fmt.Errorf("playground/admin: academic_unit (annex): %w", err)
	}
	if err := upsertMembership(tx); err != nil {
		return fmt.Errorf("playground/admin: membership: %w", err)
	}
	if err := upsertSchoolB(tx); err != nil {
		return fmt.Errorf("playground/admin: school_b: %w", err)
	}
	if err := upsertAcademicUnitB(tx); err != nil {
		return fmt.Errorf("playground/admin: academic_unit_b: %w", err)
	}
	if err := upsertMembershipB(tx); err != nil {
		return fmt.Errorf("playground/admin: membership_b: %w", err)
	}
	if err := upsertViewerRole(tx); err != nil {
		return fmt.Errorf("playground/admin: viewer_role: %w", err)
	}
	if err := upsertViewerRoleGrant(tx); err != nil {
		return fmt.Errorf("playground/admin: viewer_role_grant: %w", err)
	}
	if err := upsertViewerUser(tx); err != nil {
		return fmt.Errorf("playground/admin: viewer_user: %w", err)
	}
	if err := upsertViewerUserRole(tx); err != nil {
		return fmt.Errorf("playground/admin: viewer_user_role: %w", err)
	}
	if err := upsertViewerMembership(tx); err != nil {
		return fmt.Errorf("playground/admin: viewer_membership: %w", err)
	}
	return nil
}

func upsertAdminUser(tx *gorm.DB) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt: %w", err)
	}
	u := entities.User{
		ID:           id,
		Email:        UserEmail,
		PasswordHash: string(hash),
		FirstName:    "Admin",
		LastName:     "Playground",
		IsActive:     true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&u).Error
}

func upsertAdminUserRole(tx *gorm.DB) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	rid, err := uuid.Parse(layers.L0_ROLE_SUPER_ADMIN_ID)
	if err != nil {
		return err
	}
	// UUID determinístico para idempotencia (UNIQUE compuesto trata NULL!=NULL).
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

func upsertSuperAdminWildcardGrant(tx *gorm.DB) error {
	rid, err := uuid.Parse(layers.L0_ROLE_SUPER_ADMIN_ID)
	if err != nil {
		return err
	}
	pattern := "*"
	effect := "allow"
	// ID determinístico — alineado con la convención de L4 (SHA1 sobre
	// role_id + pattern + effect).
	gid := uuid.NewSHA1(uuid.NameSpaceOID, []byte(rid.String()+":"+pattern+":"+effect))
	g := entities.RoleGrant{
		ID:      gid,
		RoleID:  rid,
		Pattern: pattern,
		Effect:  effect,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "role_id"}, {Name: "pattern"}, {Name: "effect"}},
		DoNothing: true,
	}).Create(&g).Error
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

func upsertMembership(tx *gorm.DB) error {
	id, err := uuid.Parse(membershipID)
	if err != nil {
		return err
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	sid, err := uuid.Parse(schoolID)
	if err != nil {
		return err
	}
	m := entities.Membership{
		ID:         id,
		UserID:     uid,
		SchoolID:   sid,
		Role:       "admin",
		Metadata:   json.RawMessage(`{}`),
		IsActive:   true,
		EnrolledAt: time.Now().UTC(),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&m).Error
}

func upsertAcademicUnitAnnex(tx *gorm.DB) error {
	id, err := uuid.Parse(unitAAnnexID)
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
		Name:         unitAAnnexName,
		Code:         unitAAnnexCode,
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

func upsertSchoolB(tx *gorm.DB) error {
	id, err := uuid.Parse(schoolBID)
	if err != nil {
		return err
	}
	s := entities.School{
		ID:               id,
		Name:             schoolBName,
		Code:             schoolBCode,
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

func upsertAcademicUnitB(tx *gorm.DB) error {
	id, err := uuid.Parse(unitBID)
	if err != nil {
		return err
	}
	sid, err := uuid.Parse(schoolBID)
	if err != nil {
		return err
	}
	u := entities.AcademicUnit{
		ID:           id,
		SchoolID:     sid,
		Name:         unitBName,
		Code:         unitBCode,
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

func upsertMembershipB(tx *gorm.DB) error {
	id, err := uuid.Parse(membershipBID)
	if err != nil {
		return err
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	sid, err := uuid.Parse(schoolBID)
	if err != nil {
		return err
	}
	m := entities.Membership{
		ID:         id,
		UserID:     uid,
		SchoolID:   sid,
		Role:       "admin",
		Metadata:   json.RawMessage(`{}`),
		IsActive:   true,
		EnrolledAt: time.Now().UTC(),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&m).Error
}

// upsertViewerRole crea un rol IAM scope=school sin wildcard. Sus permisos
// se limitan al grant academic.announcements.read sembrado por
// upsertViewerRoleGrant. Pensado para validar UI con un usuario no-global.
func upsertViewerRole(tx *gorm.DB) error {
	id, err := uuid.Parse(viewerRoleID)
	if err != nil {
		return err
	}
	desc := "Solo lectura de anuncios. Usuario no-global del playground admin."
	r := entities.Role{
		ID:          id,
		Name:        viewerRoleName,
		DisplayName: "Visor de Anuncios",
		Description: &desc,
		Scope:       "school",
		IsActive:    true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&r).Error
}

// upsertViewerRoleGrant siembra dos grants para el viewer: el literal de
// announcements (para que el menú lo vea) y `*.read` (para no romper
// endpoints generales de read como notifications). Ninguno cubre verbos
// mutativos, así que el viewer sigue siendo read-only.
func upsertViewerRoleGrant(tx *gorm.DB) error {
	rid, err := uuid.Parse(viewerRoleID)
	if err != nil {
		return err
	}
	patterns := []string{viewerReadPatternMenu, viewerReadPatternAll}
	effect := "allow"
	for _, pattern := range patterns {
		gid := uuid.NewSHA1(uuid.NameSpaceOID, []byte(rid.String()+":"+pattern+":"+effect))
		g := entities.RoleGrant{
			ID:      gid,
			RoleID:  rid,
			Pattern: pattern,
			Effect:  effect,
		}
		err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "role_id"}, {Name: "pattern"}, {Name: "effect"}},
			DoNothing: true,
		}).Create(&g).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func upsertViewerUser(tx *gorm.DB) error {
	id, err := uuid.Parse(viewerUserID)
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt: %w", err)
	}
	u := entities.User{
		ID:           id,
		Email:        ViewerEmail,
		PasswordHash: string(hash),
		FirstName:    "Viewer",
		LastName:     "Playground",
		IsActive:     true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&u).Error
}

func upsertViewerUserRole(tx *gorm.DB) error {
	uid, err := uuid.Parse(viewerUserID)
	if err != nil {
		return err
	}
	rid, err := uuid.Parse(viewerRoleID)
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

// upsertViewerMembership liga al viewer a Escuela Norte + Sede Norte. Esta
// escuela tiene una sola unit académica, lo que produce el caso 1×1 ideal
// para validar que el dropdown oculte "Cambiar escuela" y "Cambiar unidad".
func upsertViewerMembership(tx *gorm.DB) error {
	id, err := uuid.Parse(viewerMembershipID)
	if err != nil {
		return err
	}
	uid, err := uuid.Parse(viewerUserID)
	if err != nil {
		return err
	}
	sid, err := uuid.Parse(schoolBID)
	if err != nil {
		return err
	}
	auid, err := uuid.Parse(unitBID)
	if err != nil {
		return err
	}
	m := entities.Membership{
		ID:             id,
		UserID:         uid,
		SchoolID:       sid,
		AcademicUnitID: &auid,
		Role:           "teacher",
		Metadata:       json.RawMessage(`{}`),
		IsActive:       true,
		EnrolledAt:     time.Now().UTC(),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&m).Error
}
