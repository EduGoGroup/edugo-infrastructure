package fixtures

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/catalog"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
)

// GuardianRelation crea el vínculo guardian↔student que las pantallas
// del rol "guardian" necesitan para resolver contra datos reales:
//
//   - 1 auth.users (student ficticio del scenario)
//   - 1 academic.memberships (role="student") en la escuela del guardian
//   - 1 academic.guardian_relations entre el guardian provisto por
//     RoleOnly(guardian) y el student creado.
//
// Requiere las capacidades "school" y "user" (apoderado) de la
// composición. Provee la capacidad "guardian_relation".
//
// Convenciones:
//
//   - El student NO se provee al resto de la composición como "user"
//     (esa capacidad ya está ocupada por el guardian); se exporta su
//     ID/email como constantes para los tests Kotlin.
//   - Cleanup borra sólo lo creado por esta fixture: el student.user,
//     su membership y la guardian_relation. NO toca al guardian (lo
//     creó RoleOnly y lo limpia RoleOnly).
type GuardianRelation struct {
	// Password en claro para el student creado. Si está vacío se
	// usa el default "E2EUser2026!".
	Password string
}

// Constantes exportadas en el JSON de fixtures-constants.json.
const (
	guardianRelationStudentIDConstant    = "E2EFixtureGuardianRelationStudentID"
	guardianRelationStudentEmailConstant = "E2EFixtureGuardianRelationStudentEmail"
	guardianRelationIDConstant           = "E2EFixtureGuardianRelationID"
)

// Sufijos UUID determinísticos en sub-namespaces fuera del rango
// usado por RoleOnly (segmento 5: 0a1/0a2/0a3) para evitar colisiones.
const (
	guardianRelationStudentUserSuffix       = "0000-0000-0000-0000000000a1"
	guardianRelationStudentMembershipSuffix = "0000-0000-0000-0000000000a2"
	guardianRelationSuffix                  = "0000-0000-0000-0000000000a3"
)

// Manifest implementa framework.Fixture.
func (f *GuardianRelation) Manifest() framework.FixtureManifest {
	return framework.FixtureManifest{
		Name:     "guardian_relation",
		Provides: []string{"guardian_relation"},
		Requires: []string{"school", "user"},
		Tables: []string{
			"auth.users",
			"academic.memberships",
			"academic.guardian_relations",
		},
		Constants: map[string]string{
			guardianRelationStudentIDConstant:    "{{.StudentID}}",
			guardianRelationStudentEmailConstant: "{{.StudentEmail}}",
			guardianRelationIDConstant:           "{{.RelationID}}",
		},
		Description: "Crea 1 student (user+membership) y el vínculo academic.guardian_relations con el guardian provisto por la composición.",
	}
}

// Apply implementa framework.Fixture. Las validaciones independientes
// de la BD se ejecutan ANTES de tocar tx, de modo que un test pueda
// invocar Apply con tx=nil y obtener errores limpios sin panics.
func (f *GuardianRelation) Apply(tx *gorm.DB, ctx *framework.ApplyContext) error {
	if ctx == nil {
		return fmt.Errorf("guardian_relation: nil ApplyContext")
	}
	school, ok := ctx.Provided["school"]
	if !ok || school.ID == "" {
		return fmt.Errorf("guardian_relation: capability %q no provista por la composición", "school")
	}
	guardian, ok := ctx.Provided["user"]
	if !ok || guardian.ID == "" {
		return fmt.Errorf("guardian_relation: capability %q no provista por la composición", "user")
	}
	schoolUUID, err := uuid.Parse(school.ID)
	if err != nil {
		return fmt.Errorf("guardian_relation: school.ID inválido (%q): %w", school.ID, err)
	}
	guardianUUID, err := uuid.Parse(guardian.ID)
	if err != nil {
		return fmt.Errorf("guardian_relation: user.ID inválido (%q): %w", guardian.ID, err)
	}
	if tx == nil {
		return fmt.Errorf("guardian_relation: nil transaction")
	}

	password := f.Password
	if password == "" {
		password = "E2EUser2026!"
	}

	studentID := framework.MakeUUID(ctx, guardianRelationStudentUserSuffix)
	membershipID := framework.MakeUUID(ctx, guardianRelationStudentMembershipSuffix)
	relationID := framework.MakeUUID(ctx, guardianRelationSuffix)

	for _, id := range []string{studentID, membershipID, relationID} {
		if err := framework.AssertNotProductionNamespace(id); err != nil {
			return err
		}
	}

	studentEmail := framework.MakeEmail(ctx, "student", "guardian_relation")

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("guardian_relation: bcrypt: %w", err)
	}

	studentUUID, err := uuid.Parse(studentID)
	if err != nil {
		return fmt.Errorf("guardian_relation: studentID UUID inválido (%q): %w", studentID, err)
	}
	membershipUUID, err := uuid.Parse(membershipID)
	if err != nil {
		return fmt.Errorf("guardian_relation: membershipID UUID inválido (%q): %w", membershipID, err)
	}
	relationUUID, err := uuid.Parse(relationID)
	if err != nil {
		return fmt.Errorf("guardian_relation: relationID UUID inválido (%q): %w", relationID, err)
	}

	student := entities.User{
		ID:           studentUUID,
		Email:        studentEmail,
		PasswordHash: string(hashed),
		FirstName:    "E2E",
		LastName:     "Student",
		IsActive:     true,
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&student).Error; err != nil {
		return fmt.Errorf("guardian_relation: insert student user: %w", err)
	}
	// Booleano crítico (F2·H5): asegurar IsActive=true incluso si la
	// fila ya existía con otro estado.
	if err := framework.UpsertBool(tx, student.TableName(), "id", student.ID, "is_active", true); err != nil {
		return err
	}

	invitationTypeID, err := catalog.ResolveInvitationTypeID(tx, "student")
	if err != nil {
		return fmt.Errorf("guardian_relation: resolve invitation_type: %w", err)
	}
	membership := entities.Membership{
		ID:               membershipUUID,
		UserID:           student.ID,
		SchoolID:         schoolUUID,
		InvitationTypeID: invitationTypeID,
		Metadata:         json.RawMessage(`{"e2e":true,"fixture":"guardian_relation"}`),
		Status:           "active",
		EnrolledAt:       time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC),
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&membership).Error; err != nil {
		return fmt.Errorf("guardian_relation: insert student membership: %w", err)
	}
	if err := framework.UpsertString(tx, membership.TableName(), "id", membership.ID, "status", "active"); err != nil {
		return err
	}

	relation := entities.GuardianRelation{
		ID:               relationUUID,
		GuardianID:       guardianUUID,
		StudentID:        student.ID,
		RelationshipType: "parent",
		IsPrimary:        true,
		IsActive:         true,
		Status:           "active",
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&relation).Error; err != nil {
		return fmt.Errorf("guardian_relation: insert guardian_relation: %w", err)
	}
	if err := framework.UpsertBool(tx, relation.TableName(), "id", relation.ID, "is_active", true); err != nil {
		return err
	}

	ctx.Provide("guardian_relation", framework.ProvidedEntity{
		Kind: "guardian_relation",
		ID:   relationID,
		Extra: map[string]string{
			"guardian_id":   guardianUUID.String(),
			"student_id":    studentID,
			"student_email": studentEmail,
		},
	})

	ctx.SetConstant(guardianRelationStudentIDConstant, studentID)
	ctx.SetConstant(guardianRelationStudentEmailConstant, studentEmail)
	ctx.SetConstant(guardianRelationIDConstant, relationID)
	return nil
}

// Cleanup implementa framework.Fixture. Borra en orden inverso al de
// creación: guardian_relations → memberships → users. Sólo el student
// que creó esta fixture (membership y user filtrados por SchemaPrefix);
// NO el guardian (lo creó y lo limpia RoleOnly).
func (f *GuardianRelation) Cleanup(tx *gorm.DB, ctx *framework.ApplyContext) error {
	if tx == nil {
		return fmt.Errorf("guardian_relation cleanup: nil transaction")
	}
	if ctx == nil || ctx.SchemaPrefix == "" {
		return fmt.Errorf("guardian_relation cleanup: SchemaPrefix vacío")
	}
	prefix := ctx.SchemaPrefix
	tables := []struct {
		name string
		col  string
	}{
		{"academic.guardian_relations", "id"},
		{"academic.memberships", "id"},
		{"auth.users", "id"},
	}
	for _, t := range tables {
		if _, err := framework.DeleteByPrefix(tx, t.name, t.col, prefix); err != nil {
			return fmt.Errorf("guardian_relation cleanup %s: %w", t.name, err)
		}
	}
	return nil
}

// E2EFixtureGuardianRelationStudentID es la clave bajo la cual se
// exporta el UUID del student creado por la fixture en el JSON de
// constantes consumido por los tests Kotlin.
const E2EFixtureGuardianRelationStudentID = guardianRelationStudentIDConstant

// E2EFixtureGuardianRelationStudentEmail es la clave del email del
// student.
const E2EFixtureGuardianRelationStudentEmail = guardianRelationStudentEmailConstant

// E2EFixtureGuardianRelationID es la clave del UUID de la fila en
// academic.guardian_relations.
const E2EFixtureGuardianRelationID = guardianRelationIDConstant
