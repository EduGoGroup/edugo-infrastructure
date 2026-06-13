// Package focal_pantalla define un playground minimalista pensado para
// iterar sobre el diseño y el pulido visual de UNA sola pantalla
// (announcements-list). El menú dinámico del KMP filtra los items por los
// grants efectivos del rol; al limitar los grants a `announcements.*` o
// `announcements.read`, el menú colapsa a un solo item ("Anuncios"). Eso
// permite enfocar el trabajo de UI sin la distracción de las 70+ pantallas
// que siembra el sistema completo (L4).
//
// Convive con el playground admin sin pisarlo: usa IDs y emails propios.
// La política de la convención de playgrounds (project_edugo_playgrounds_convention
// en memoria) prohíbe editar admin para satisfacer este caso — siempre se
// crea un paquete nuevo y los anteriores quedan como foto histórica.
//
// Lo que siembra:
//  1. auth.users           — 3 usuarios (focal-admin, focal-viewer y focal-author).
//  2. iam.roles            — 3 roles scope=school (focal_pantalla_admin, focal_pantalla_viewer, focal_pantalla_author).
//  3. iam.role_grants      — admin: announcements.* + *.read; viewer: announcements.read + *.read;
//     author: announcements.read + announcements.create + *.read (caso intermedio: crea pero no modifica/elimina).
//  4. iam.user_roles       — assignments user × rol.
//  5. academic.schools     — 1 escuela.
//  6. academic.academic_units — 1 unidad raíz.
//  7. academic.memberships — admin, viewer y author en la misma escuela/unidad.
//  8. academic.announcements — 5 anuncios de prueba variados (pinned, scopes
//     distintos, con y sin expiración) para validar lista/edit/create.
//
// Idempotente: usa OnConflict DoNothing por id en todas las inserciones.
package focal_pantalla

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/catalog"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// Credenciales (password compartido para simplicidad de prueba; la app es
	// dev-only y los playgrounds son fixtures, no usuarios reales).
	AdminEmail  = "focal-admin@edugo.local"
	ViewerEmail = "focal-viewer@edugo.local"
	AuthorEmail = "focal-author@edugo.local"
	Password    = "12345678"

	// Rango UUID 61000000-...: reservado para focal_pantalla. Convive con el
	// rango 60000000-... del playground admin sin colisión.
	adminUserID  = "61000000-0000-0000-0000-000000000001"
	viewerUserID = "61000000-0000-0000-0000-000000000002"
	schoolID     = "61000000-0000-0000-0000-000000000003"
	unitID       = "61000000-0000-0000-0000-000000000004"
	adminMembID  = "61000000-0000-0000-0000-000000000005"
	viewerMembID = "61000000-0000-0000-0000-000000000006"
	authorUserID = "61000000-0000-0000-0000-000000000007"
	authorMembID = "61000000-0000-0000-0000-000000000008"

	// IDs de roles en rango 11000000-... (separado del 10000000-... del admin).
	adminRoleID    = "11000000-0000-0000-0000-000000000001"
	viewerRoleID   = "11000000-0000-0000-0000-000000000002"
	authorRoleID   = "11000000-0000-0000-0000-000000000003"
	adminRoleName  = "focal_pantalla_admin"
	viewerRoleName = "focal_pantalla_viewer"
	authorRoleName = "focal_pantalla_author"

	// Grants del rol admin: CRUD completo sobre anuncios + lectura general.
	// IMPORTANTE: el recurso "announcements" está colgado bajo el padre
	// "academic" en L4 (parent_id != null). El matcher del menú construye
	// el resourcePath jerárquicamente, por lo que el path real del recurso
	// es "academic.announcements" — los grants deben usar ese prefijo
	// completo, no solo "announcements.*".
	adminPatternCrud = "academic.announcements.*"
	adminPatternRead = "*.read"

	// Grants del rol viewer: lectura sobre anuncios + lectura general. Misma
	// regla: el path completo es "academic.announcements", por eso el
	// literal incluye el prefijo academic.
	viewerPatternMenu = "academic.announcements.read"
	viewerPatternAll  = "*.read"

	// Grants del rol author: lectura + creación de anuncios. SIN update ni
	// delete — caso "termino medio" para validar que la UI oculta editar /
	// eliminar pero deja crear. La lectura general (*.read) cubre endpoints
	// auxiliares como notifications/unread-count.
	authorPatternRead   = "academic.announcements.read"
	authorPatternCreate = "academic.announcements.create"
	authorPatternAll    = "*.read"

	schoolCode = "FOCAL-PANTALLA"
	schoolName = "Escuela Focal — Pantalla"
	unitCode   = "FOCAL-PANTALLA-MAIN"
	unitName   = "Sede Única"

	academicYear = 2026
)

// Apply siembra el playground focal_pantalla. Asume que L0 corrió (templates
// + screen instance announcements-list ya existen). Idempotente.
func Apply(tx *gorm.DB) error {
	if err := upsertSchool(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: school: %w", err)
	}
	if err := upsertAcademicUnit(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: academic_unit: %w", err)
	}
	if err := upsertAdminRole(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: admin_role: %w", err)
	}
	if err := upsertAdminRoleGrants(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: admin_role_grants: %w", err)
	}
	if err := upsertAdminUser(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: admin_user: %w", err)
	}
	if err := upsertAdminUserRole(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: admin_user_role: %w", err)
	}
	if err := upsertAdminMembership(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: admin_membership: %w", err)
	}
	if err := upsertViewerRole(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: viewer_role: %w", err)
	}
	if err := upsertViewerRoleGrants(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: viewer_role_grants: %w", err)
	}
	if err := upsertViewerUser(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: viewer_user: %w", err)
	}
	if err := upsertViewerUserRole(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: viewer_user_role: %w", err)
	}
	if err := upsertViewerMembership(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: viewer_membership: %w", err)
	}
	if err := upsertAuthorRole(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: author_role: %w", err)
	}
	if err := upsertAuthorRoleGrants(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: author_role_grants: %w", err)
	}
	if err := upsertAuthorUser(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: author_user: %w", err)
	}
	if err := upsertAuthorUserRole(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: author_user_role: %w", err)
	}
	if err := upsertAuthorMembership(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: author_membership: %w", err)
	}
	if err := upsertAnnouncements(tx); err != nil {
		return fmt.Errorf("playground/focal_pantalla: announcements: %w", err)
	}
	return nil
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

// upsertAdminRole crea el rol focal_pantalla_admin (scope=school, sin
// wildcard global). Sus permisos se definen vía grants en
// upsertAdminRoleGrants — todo el ciclo CRUD sobre anuncios.
func upsertAdminRole(tx *gorm.DB) error {
	id, err := uuid.Parse(adminRoleID)
	if err != nil {
		return err
	}
	desc := "CRUD completo sobre anuncios. Playground focal-pantalla."
	r := entities.Role{
		ID:          id,
		Name:        adminRoleName,
		DisplayName: "Admin Focal — CRUD anuncios",
		Description: &desc,
		Scope:       "school",
		IsActive:    true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&r).Error
}

func upsertAdminRoleGrants(tx *gorm.DB) error {
	return upsertRoleGrants(tx, adminRoleID, []string{adminPatternCrud, adminPatternRead})
}

// upsertViewerRole crea el rol focal_pantalla_viewer — solo lectura sobre
// anuncios, para validar la UI con un usuario sin permisos de mutación.
func upsertViewerRole(tx *gorm.DB) error {
	id, err := uuid.Parse(viewerRoleID)
	if err != nil {
		return err
	}
	desc := "Solo lectura de anuncios. Playground focal-pantalla."
	r := entities.Role{
		ID:          id,
		Name:        viewerRoleName,
		DisplayName: "Viewer Focal — solo lectura",
		Description: &desc,
		Scope:       "school",
		IsActive:    true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&r).Error
}

func upsertViewerRoleGrants(tx *gorm.DB) error {
	return upsertRoleGrants(tx, viewerRoleID, []string{viewerPatternMenu, viewerPatternAll})
}

// upsertAuthorRole crea el rol focal_pantalla_author — caso intermedio:
// puede leer y crear anuncios, pero NO modificar ni eliminar. Sirve para
// validar que la UI oculta los botones de update/delete pero deja crear.
func upsertAuthorRole(tx *gorm.DB) error {
	id, err := uuid.Parse(authorRoleID)
	if err != nil {
		return err
	}
	desc := "Crea anuncios sin poder modificarlos ni eliminarlos. Playground focal-pantalla."
	r := entities.Role{
		ID:          id,
		Name:        authorRoleName,
		DisplayName: "Author Focal — crea sin editar",
		Description: &desc,
		Scope:       "school",
		IsActive:    true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&r).Error
}

func upsertAuthorRoleGrants(tx *gorm.DB) error {
	return upsertRoleGrants(tx, authorRoleID, []string{authorPatternRead, authorPatternCreate, authorPatternAll})
}

// upsertRoleGrants inserta una lista de patterns como allow-grants para un
// rol. ID determinístico SHA1(role_id:pattern:effect) alineado a la
// convención de L4 para que el reseed sea idempotente.
func upsertRoleGrants(tx *gorm.DB, roleID string, patterns []string) error {
	rid, err := uuid.Parse(roleID)
	if err != nil {
		return err
	}
	effect := "allow"
	for _, pattern := range patterns {
		gid := uuid.NewSHA1(uuid.NameSpaceOID, []byte(rid.String()+":"+pattern+":"+effect))
		g := entities.RoleGrant{
			ID:      gid,
			RoleID:  rid,
			Pattern: pattern,
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

func upsertAdminUser(tx *gorm.DB) error {
	return upsertUser(tx, adminUserID, AdminEmail, "Admin", "Focal")
}

func upsertViewerUser(tx *gorm.DB) error {
	return upsertUser(tx, viewerUserID, ViewerEmail, "Viewer", "Focal")
}

func upsertAuthorUser(tx *gorm.DB) error {
	return upsertUser(tx, authorUserID, AuthorEmail, "Author", "Focal")
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

func upsertAdminUserRole(tx *gorm.DB) error {
	return upsertUserRole(tx, adminUserID, adminRoleID)
}

func upsertViewerUserRole(tx *gorm.DB) error {
	return upsertUserRole(tx, viewerUserID, viewerRoleID)
}

func upsertAuthorUserRole(tx *gorm.DB) error {
	return upsertUserRole(tx, authorUserID, authorRoleID)
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

func upsertAdminMembership(tx *gorm.DB) error {
	return upsertMembership(tx, adminMembID, adminUserID, "admin")
}

func upsertViewerMembership(tx *gorm.DB) error {
	return upsertMembership(tx, viewerMembID, viewerUserID, "teacher")
}

func upsertAuthorMembership(tx *gorm.DB) error {
	return upsertMembership(tx, authorMembID, authorUserID, "teacher")
}

// upsertMembership liga al usuario a la única escuela/unidad del playground.
// AcademicUnitID se pasa por puntero para forzar que el contexto al loguear
// resuelva con unit completa sin necesidad de switch-context.
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
	invitationTypeID, err := catalog.ResolveInvitationTypeID(tx, roleKind)
	if err != nil {
		return err
	}
	m := entities.Membership{
		ID:               id,
		UserID:           uid,
		SchoolID:         sid,
		AcademicUnitID:   &auid,
		InvitationTypeID: invitationTypeID,
		Metadata:         json.RawMessage(`{}`),
		IsActive:         true,
		EnrolledAt:       time.Now().UTC(),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&m).Error
}

// upsertAnnouncements siembra 5 anuncios variados sobre la única escuela
// del playground. La variedad cubre los ejes que la UI debería pintar:
// pinned vs no pinned (chip de estado), scope school/unit/role (subtitle
// transformado en el contract), con/sin expires_at (badge "vencido"), y
// distintos created_at (orden de la lista). Autor: focal-admin.
//
// IDs determinísticos en rango 61000010-... para que la siembra sea
// idempotente y los anuncios sean referenciables desde tests si hace
// falta. OnConflict DoNothing por id.
func upsertAnnouncements(tx *gorm.DB) error {
	sid, err := uuid.Parse(schoolID)
	if err != nil {
		return err
	}
	uid, err := uuid.Parse(unitID)
	if err != nil {
		return err
	}
	aid, err := uuid.Parse(adminUserID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	publishedNow := now
	publishedLastWeek := now.AddDate(0, 0, -7)
	publishedYesterday := now.AddDate(0, 0, -1)
	publishedLastMonth := now.AddDate(0, -1, 0)
	expiresNextWeek := now.AddDate(0, 0, 7)
	expiredYesterday := now.AddDate(0, 0, -1)

	items := []entities.Announcement{
		{
			ID:             mustParseUUID("61000010-0000-0000-0000-000000000001"),
			SchoolID:       sid,
			AcademicUnitID: &uid,
			AuthorID:       aid,
			Title:          "Bienvenida al ciclo lectivo 2026",
			Body:           "Damos la bienvenida a estudiantes y familias al nuevo año. Las clases inician el lunes a las 08:00.",
			Scope:          "school",
			IsPinned:       true,
			PublishedAt:    &publishedLastWeek,
		},
		{
			ID:             mustParseUUID("61000010-0000-0000-0000-000000000002"),
			SchoolID:       sid,
			AcademicUnitID: &uid,
			AuthorID:       aid,
			Title:          "Reunión de apoderados — próximo viernes",
			Body:           "Se cita a la reunión general de apoderados el viernes a las 19:00 en el salón de actos.",
			Scope:          "school",
			IsPinned:       false,
			PublishedAt:    &publishedYesterday,
			ExpiresAt:      &expiresNextWeek,
		},
		{
			ID:             mustParseUUID("61000010-0000-0000-0000-000000000003"),
			SchoolID:       sid,
			AcademicUnitID: &uid,
			AuthorID:       aid,
			Title:          "Suspensión de actividades el sábado",
			Body:           "Por mantenimiento eléctrico programado, las actividades del sábado quedan suspendidas.",
			Scope:          "unit",
			IsPinned:       true,
			PublishedAt:    &publishedNow,
		},
		{
			ID:             mustParseUUID("61000010-0000-0000-0000-000000000004"),
			SchoolID:       sid,
			AcademicUnitID: &uid,
			AuthorID:       aid,
			Title:          "Convocatoria taller de robótica",
			Body:           "Docentes interesados en facilitar el taller de robótica deben inscribirse antes del 31 de mayo.",
			Scope:          "role",
			TargetRoles:    pq.StringArray{"teacher"},
			IsPinned:       false,
			PublishedAt:    &publishedNow,
		},
		{
			ID:             mustParseUUID("61000010-0000-0000-0000-000000000005"),
			SchoolID:       sid,
			AcademicUnitID: &uid,
			AuthorID:       aid,
			Title:          "Recordatorio: vencimiento de aranceles",
			Body:           "El plazo para regularizar los aranceles del primer trimestre venció ayer. Pasar por administración.",
			Scope:          "school",
			IsPinned:       false,
			PublishedAt:    &publishedLastMonth,
			ExpiresAt:      &expiredYesterday,
		},
	}

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&items).Error
}

// mustParseUUID es helper interno: convierte un literal UUID en uuid.UUID
// o entra en pánico si el literal es inválido. Solo se usa para constantes
// hardcodeadas dentro de este paquete, donde el formato está bajo control.
func mustParseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		panic(fmt.Sprintf("focal_pantalla: UUID inválido %q: %v", s, err))
	}
	return id
}
