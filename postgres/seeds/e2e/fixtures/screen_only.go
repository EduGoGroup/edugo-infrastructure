package fixtures

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/catalog"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
)

// ScreenOnly crea el contenido mínimo necesario para que una pantalla
// del production seed cargue dentro de un scenario E2E.
//
// IMPORTANTE: NO toca filas en ui_config.screen_templates,
// ui_config.screen_instances ni ui_config.resource_screens — todas
// pertenecen al production seed (UUIDs `10000000-...`) y están
// protegidas por framework.AssertNotProductionNamespace. Esta fixture
// asume que esas filas ya existen y se limita a poblar contenido de
// dominio (ej. una assessment de prueba para `assessments-list`) en el
// namespace del scenario.
//
// Para pantallas que requieren múltiples FKs poco prácticos de armar
// (ej. `grades-list`: membership + subject + period), se cae al caso
// "default" — la fixture no falla, sólo se registra que no hay
// contenido específico. La pantalla seguirá siendo navegable porque el
// production seed ya provee el chrome.
type ScreenOnly struct {
	// ScreenKey identifica la pantalla del production seed
	// (ej. "assessments-list", "grades-list").
	ScreenKey string
}

// Manifest implementa framework.Fixture.
//
// Nota: el conjunto efectivo de tablas tocadas depende del ScreenKey
// (cada branch de Apply crea un subconjunto distinto). Declaramos aquí
// la unión de todas las tablas posibles para que el cleanup las
// recorra; las tablas que no recibieron filas en un scenario concreto
// simplemente no encontrarán matches en el DELETE LIKE.
func (f *ScreenOnly) Manifest() framework.FixtureManifest {
	return framework.FixtureManifest{
		Name:     "screen_only",
		Provides: []string{"screen_data"},
		Requires: []string{"school"},
		Tables: []string{
			"assessment.assessment",
			"academic.grades",
			"academic.memberships",
			"auth.users",
			"academic.academic_periods",
			"academic.academic_units",
			"academic.subjects",
		},
		Constants: map[string]string{
			"E2EFixtureScreenOnlyScreenKey": "{{.ScreenKey}}",
		},
		Description: "Crea el contenido mínimo de dominio (assessment, grade, etc.) para que una pantalla del production seed cargue.",
	}
}

// Apply implementa framework.Fixture. Las validaciones que no
// dependen de la BD se ejecutan ANTES de cualquier acceso a tx, de
// modo que un test puede invocar Apply con tx=nil y obtener errores
// limpios sin panics.
func (f *ScreenOnly) Apply(tx *gorm.DB, ctx *framework.ApplyContext) error {
	if f.ScreenKey == "" {
		return fmt.Errorf("screen_only: ScreenKey requerido")
	}
	if ctx == nil {
		return fmt.Errorf("screen_only: nil ApplyContext")
	}
	school, ok := ctx.Provided["school"]
	if !ok || school.ID == "" {
		return fmt.Errorf("screen_only: capability %q no provista por la composición", "school")
	}
	schoolUUID, err := uuid.Parse(school.ID)
	if err != nil {
		return fmt.Errorf("screen_only: school.ID inválido (%q): %w", school.ID, err)
	}
	if tx == nil {
		return fmt.Errorf("screen_only: nil transaction")
	}

	ctx.SetConstant("E2EFixtureScreenOnlyScreenKey", f.ScreenKey)

	switch f.ScreenKey {
	case "assessments-list":
		return f.applyAssessmentsList(tx, ctx, schoolUUID)
	case "grades-list":
		return f.applyGradesList(tx, ctx, schoolUUID)
	default:
		// La pantalla no requiere contenido específico desde esta
		// fixture. Se registra una constante para diagnóstico y se
		// devuelve sin error (C-REQ-9.2: behavior pragmático).
		ctx.SetConstant("E2EFixtureScreenOnlyContent", "none")
		ctx.Provide("screen_data", framework.ProvidedEntity{
			Kind: "screen_data",
			ID:   "",
			Extra: map[string]string{
				"screen_key": f.ScreenKey,
				"content":    "none",
				"reason":     "screenKey sin contenido específico — pantalla sólo necesita el chrome del production seed",
			},
		})
		return nil
	}
}

// applyAssessmentsList crea 1 assessment de prueba para que la pantalla
// `assessments-list` muestre datos en el scenario.
//
// N4 (ADR 0019): el esquema nuevo de assessment.assessment está anclado al
// modelo de sesión. La evaluación exige tres FKs NOT NULL: school_id, subject_id
// (→academic.subjects, catálogo de escuela) y created_by_membership_id
// (→academic.memberships del docente autor). La fixture compone subject + docente
// (user + membership) con UUIDs determinísticos sufijo 0a55e5500002..0a55e5500004
// dentro del namespace del scenario, antes de insertar la evaluación.
func (f *ScreenOnly) applyAssessmentsList(tx *gorm.DB, ctx *framework.ApplyContext, schoolUUID uuid.UUID) error {
	// 1) Subject (catálogo de escuela; FK RESTRICT de assessment.subject_id).
	subjectIDStr := framework.MakeUUID(ctx, "0000-0000-0000-0a55e5500002")
	if err := framework.AssertNotProductionNamespace(subjectIDStr); err != nil {
		return err
	}
	subjectID, err := uuid.Parse(subjectIDStr)
	if err != nil {
		return fmt.Errorf("screen_only: subject UUID inválido (%q): %w", subjectIDStr, err)
	}
	subject := entities.Subject{
		ID:       subjectID,
		SchoolID: schoolUUID,
		Name:     "ScreenOnly Sample Assessment Subject",
		IsActive: true,
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&subject).Error; err != nil {
		return fmt.Errorf("screen_only: insert subject: %w", err)
	}
	if err := framework.UpsertBool(tx, subject.TableName(), "id", subject.ID, "is_active", true); err != nil {
		return err
	}

	// 2) Teacher user + membership (autor; FK RESTRICT de created_by_membership_id).
	teacherUserIDStr := framework.MakeUUID(ctx, "0000-0000-0000-0a55e5500003")
	if err := framework.AssertNotProductionNamespace(teacherUserIDStr); err != nil {
		return err
	}
	teacherUserID, err := uuid.Parse(teacherUserIDStr)
	if err != nil {
		return fmt.Errorf("screen_only: teacher user UUID inválido (%q): %w", teacherUserIDStr, err)
	}
	teacherEmail := framework.MakeEmail(ctx, "teacher", "screen_only_assessments")
	hashed, err := bcrypt.GenerateFromPassword([]byte("E2EUser2026!"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("screen_only: bcrypt: %w", err)
	}
	teacherUser := entities.User{
		ID:           teacherUserID,
		Email:        teacherEmail,
		PasswordHash: string(hashed),
		FirstName:    "ScreenOnly",
		LastName:     "Teacher",
		IsActive:     true,
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoNothing: true,
	}).Create(&teacherUser).Error; err != nil {
		return fmt.Errorf("screen_only: insert teacher user: %w", err)
	}
	if err := framework.UpsertBool(tx, teacherUser.TableName(), "id", teacherUser.ID, "is_active", true); err != nil {
		return err
	}

	teacherMembershipIDStr := framework.MakeUUID(ctx, "0000-0000-0000-0a55e5500004")
	if err := framework.AssertNotProductionNamespace(teacherMembershipIDStr); err != nil {
		return err
	}
	teacherMembershipID, err := uuid.Parse(teacherMembershipIDStr)
	if err != nil {
		return fmt.Errorf("screen_only: teacher membership UUID inválido (%q): %w", teacherMembershipIDStr, err)
	}
	teacherInvitationTypeID, err := catalog.ResolveInvitationTypeID(tx, "teacher")
	if err != nil {
		return fmt.Errorf("screen_only: resolve teacher invitation_type: %w", err)
	}
	teacherMembership := entities.Membership{
		ID:               teacherMembershipID,
		UserID:           teacherUser.ID,
		SchoolID:         schoolUUID,
		InvitationTypeID: teacherInvitationTypeID,
		Metadata:         json.RawMessage(`{"e2e":true,"fixture":"screen_only","screen_key":"assessments-list"}`),
		IsActive:         true,
		EnrolledAt:       time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&teacherMembership).Error; err != nil {
		return fmt.Errorf("screen_only: insert teacher membership: %w", err)
	}
	if err := framework.UpsertBool(tx, teacherMembership.TableName(), "id", teacherMembership.ID, "is_active", true); err != nil {
		return err
	}

	// 3) Assessment (esquema nuevo: school_id + subject_id + created_by_membership_id).
	assessmentID := framework.MakeUUID(ctx, "0000-0000-0000-0a55e5500001")
	if err := framework.AssertNotProductionNamespace(assessmentID); err != nil {
		return err
	}
	parsed, err := uuid.Parse(assessmentID)
	if err != nil {
		return fmt.Errorf("screen_only: assessment UUID inválido (%q): %w", assessmentID, err)
	}

	title := "ScreenOnly Sample Assessment"
	desc := "Assessment de prueba creada por la fixture screen_only para poblar la pantalla assessments-list."

	assessment := entities.Assessment{
		ID:                    parsed,
		SchoolID:              schoolUUID,
		CreatedByMembershipID: teacherMembership.ID,
		SubjectID:             subject.ID,
		Title:                 title,
		Description:           &desc,
		SourceType:            "manual",
		Status:                "draft",
		QuestionsCount:        0,
		PassThreshold:         70,
		IsTimed:               false,
		ShuffleQuestions:      false,
		ShowCorrectAnswers:    true,
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&assessment).Error; err != nil {
		return fmt.Errorf("screen_only: insert assessment: %w", err)
	}

	ctx.Provide("screen_data", framework.ProvidedEntity{
		Kind: "screen_data",
		ID:   assessmentID,
		Extra: map[string]string{
			"screen_key": f.ScreenKey,
			"content":    "assessment",
		},
	})
	ctx.SetConstant("E2EFixtureScreenOnlyAssessmentID", assessmentID)
	ctx.SetConstant("E2EFixtureScreenOnlyAssessmentTitle", title)
	ctx.SetConstant("E2EFixtureScreenOnlyAssessmentSubjectID", subjectIDStr)
	ctx.SetConstant("E2EFixtureScreenOnlyAssessmentAuthorMembershipID", teacherMembershipIDStr)
	return nil
}

// applyGradesList crea el contenido mínimo necesario para que la
// pantalla `grades-list` muestre datos reales en el scenario.
//
// Para que una fila en academic.grades sea válida hacen falta tres FKs
// NOT NULL (membership, subject, period) más la TeacherID opcional. La
// fixture compone todo con UUIDs determinísticos sufijo
// 0a55e5500111..115 dentro del namespace del scenario:
//
//   - 1 academic.subjects
//   - 1 academic.academic_units
//   - 1 academic.academic_periods (atado a la unidad anterior)
//   - 1 auth.users (alumno)
//   - 1 academic.memberships (alumno)
//   - 1 academic.grades (apunta al teacher provisto por role_only si
//     `ctx.Provided["user"]` está disponible)
//
// El teacher se toma de la capacidad "user" provista por role_only
// cuando el scenario incluye ese rol; es opcional porque la columna
// `teacher_id` admite NULL.
func (f *ScreenOnly) applyGradesList(tx *gorm.DB, ctx *framework.ApplyContext, schoolUUID uuid.UUID) error {
	// 1) Subject
	subjectIDStr := framework.MakeUUID(ctx, "0000-0000-0000-0a55e5500111")
	if err := framework.AssertNotProductionNamespace(subjectIDStr); err != nil {
		return err
	}
	subjectID, err := uuid.Parse(subjectIDStr)
	if err != nil {
		return fmt.Errorf("screen_only: subject UUID inválido (%q): %w", subjectIDStr, err)
	}
	subject := entities.Subject{
		ID:       subjectID,
		SchoolID: schoolUUID,
		Name:     "ScreenOnly Sample Subject",
		IsActive: true,
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&subject).Error; err != nil {
		return fmt.Errorf("screen_only: insert subject: %w", err)
	}
	if err := framework.UpsertBool(tx, subject.TableName(), "id", subject.ID, "is_active", true); err != nil {
		return err
	}

	// 2) AcademicUnit (el período se ata a la unidad además de la escuela).
	unitIDStr := framework.MakeUUID(ctx, "0000-0000-0000-0a55e5500110")
	if err := framework.AssertNotProductionNamespace(unitIDStr); err != nil {
		return err
	}
	unitID, err := uuid.Parse(unitIDStr)
	if err != nil {
		return fmt.Errorf("screen_only: unit UUID inválido (%q): %w", unitIDStr, err)
	}
	unit := entities.AcademicUnit{
		ID:           unitID,
		SchoolID:     schoolUUID,
		Name:         "ScreenOnly Sample Unit",
		Code:         "SCREENONLY-UNIT",
		Type:         "class",
		AcademicYear: 2026,
		IsActive:     true,
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&unit).Error; err != nil {
		return fmt.Errorf("screen_only: insert academic_unit: %w", err)
	}

	// 3) AcademicPeriod
	periodIDStr := framework.MakeUUID(ctx, "0000-0000-0000-0a55e5500112")
	if err := framework.AssertNotProductionNamespace(periodIDStr); err != nil {
		return err
	}
	periodID, err := uuid.Parse(periodIDStr)
	if err != nil {
		return fmt.Errorf("screen_only: period UUID inválido (%q): %w", periodIDStr, err)
	}
	period := entities.AcademicPeriod{
		ID:             periodID,
		SchoolID:       schoolUUID,
		AcademicUnitID: unitID,
		Name:           "ScreenOnly Sample Period",
		Type:           "semester",
		StartDate:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC),
		AcademicYear:   2026,
		IsActive:       true,
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&period).Error; err != nil {
		return fmt.Errorf("screen_only: insert academic_period: %w", err)
	}
	if err := framework.UpsertBool(tx, period.TableName(), "id", period.ID, "is_active", true); err != nil {
		return err
	}

	// 3) Student user
	studentIDStr := framework.MakeUUID(ctx, "0000-0000-0000-0a55e5500113")
	if err := framework.AssertNotProductionNamespace(studentIDStr); err != nil {
		return err
	}
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		return fmt.Errorf("screen_only: student UUID inválido (%q): %w", studentIDStr, err)
	}
	studentEmail := framework.MakeEmail(ctx, "student", "screen_only_grades")
	hashed, err := bcrypt.GenerateFromPassword([]byte("E2EUser2026!"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("screen_only: bcrypt: %w", err)
	}
	student := entities.User{
		ID:           studentID,
		Email:        studentEmail,
		PasswordHash: string(hashed),
		FirstName:    "ScreenOnly",
		LastName:     "Student",
		IsActive:     true,
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoNothing: true,
	}).Create(&student).Error; err != nil {
		return fmt.Errorf("screen_only: insert student user: %w", err)
	}
	if err := framework.UpsertBool(tx, student.TableName(), "id", student.ID, "is_active", true); err != nil {
		return err
	}

	// 4) Student membership
	studentMembershipIDStr := framework.MakeUUID(ctx, "0000-0000-0000-0a55e5500114")
	if err := framework.AssertNotProductionNamespace(studentMembershipIDStr); err != nil {
		return err
	}
	studentMembershipID, err := uuid.Parse(studentMembershipIDStr)
	if err != nil {
		return fmt.Errorf("screen_only: student membership UUID inválido (%q): %w", studentMembershipIDStr, err)
	}
	studentInvitationTypeID, err := catalog.ResolveInvitationTypeID(tx, "student")
	if err != nil {
		return fmt.Errorf("screen_only: resolve student invitation_type: %w", err)
	}
	studentMembership := entities.Membership{
		ID:               studentMembershipID,
		UserID:           student.ID,
		SchoolID:         schoolUUID,
		InvitationTypeID: studentInvitationTypeID,
		Metadata:         json.RawMessage(`{"e2e":true,"fixture":"screen_only","screen_key":"grades-list"}`),
		IsActive:         true,
		EnrolledAt:       time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&studentMembership).Error; err != nil {
		return fmt.Errorf("screen_only: insert student membership: %w", err)
	}
	if err := framework.UpsertBool(tx, studentMembership.TableName(), "id", studentMembership.ID, "is_active", true); err != nil {
		return err
	}

	// 5) Grade. TeacherID es opcional: si role_only ya provee "user"
	// (caso teacher_grades_only) lo referenciamos; si no, se queda nil.
	gradeIDStr := framework.MakeUUID(ctx, "0000-0000-0000-0a55e5500115")
	if err := framework.AssertNotProductionNamespace(gradeIDStr); err != nil {
		return err
	}
	gradeID, err := uuid.Parse(gradeIDStr)
	if err != nil {
		return fmt.Errorf("screen_only: grade UUID inválido (%q): %w", gradeIDStr, err)
	}
	gradeValue := 85.0
	gradeLetter := "A"
	var teacherIDPtr *uuid.UUID
	if teacherEntity, ok := ctx.Provided["user"]; ok && teacherEntity.ID != "" {
		parsedTeacher, perr := uuid.Parse(teacherEntity.ID)
		if perr != nil {
			return fmt.Errorf("screen_only: teacher (provided user) UUID inválido (%q): %w", teacherEntity.ID, perr)
		}
		teacherIDPtr = &parsedTeacher
	}
	grade := entities.Grade{
		ID:               gradeID,
		MembershipID:     studentMembership.ID,
		SubjectID:        subject.ID,
		PeriodID:         period.ID,
		GradeValue:       &gradeValue,
		GradeLetter:      &gradeLetter,
		AssessmentScores: json.RawMessage("[]"),
		TeacherID:        teacherIDPtr,
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&grade).Error; err != nil {
		return fmt.Errorf("screen_only: insert grade: %w", err)
	}

	ctx.Provide("screen_data", framework.ProvidedEntity{
		Kind: "screen_data",
		ID:   gradeIDStr,
		Extra: map[string]string{
			"screen_key": f.ScreenKey,
			"content":    "grade",
		},
	})
	ctx.SetConstant("E2EFixtureScreenOnlyGradeID", gradeIDStr)
	ctx.SetConstant("E2EFixtureScreenOnlySubjectID", subjectIDStr)
	ctx.SetConstant("E2EFixtureScreenOnlyAcademicPeriodID", periodIDStr)
	ctx.SetConstant("E2EFixtureScreenOnlyStudentMembershipID", studentMembershipIDStr)
	ctx.SetConstant("E2EFixtureScreenOnlyStudentUserEmail", studentEmail)
	return nil
}

// Cleanup implementa framework.Fixture. Borra exclusivamente las filas
// con SchemaPrefix del scenario.
func (f *ScreenOnly) Cleanup(tx *gorm.DB, ctx *framework.ApplyContext) error {
	if tx == nil {
		return fmt.Errorf("screen_only cleanup: nil transaction")
	}
	if ctx == nil || ctx.SchemaPrefix == "" {
		return fmt.Errorf("screen_only cleanup: SchemaPrefix vacío")
	}
	prefix := ctx.SchemaPrefix
	// Orden de borrado: respetamos las FKs. Para grades-list:
	// grades → memberships → users → academic_periods → academic_units → subjects.
	// Para assessments-list (N4): assessment.assessment PRIMERO (FK RESTRICT a
	// subjects y memberships del docente autor) → memberships → users → subjects.
	// Las tablas que no tienen filas en el scenario devuelven 0 rows.
	tables := []struct {
		name string
		col  string
	}{
		{"assessment.assessment", "id"},
		{"academic.grades", "id"},
		{"academic.memberships", "id"},
		{"auth.users", "id"},
		{"academic.academic_periods", "id"},
		{"academic.academic_units", "id"},
		{"academic.subjects", "id"},
	}
	for _, t := range tables {
		if _, err := framework.DeleteByPrefix(tx, t.name, t.col, prefix); err != nil {
			return fmt.Errorf("screen_only cleanup %s: %w", t.name, err)
		}
	}
	return nil
}

// SupportedScreenKeys devuelve los screenKeys que la fixture sabe
// poblar con contenido específico. El resto se acepta pero sin
// inserts adicionales. Útil para tests y mensajes de diagnóstico.
func SupportedScreenKeys() []string {
	return []string{"assessments-list", "grades-list"}
}

// FormatSupportedScreenKeys devuelve los screenKeys soportados como
// string CSV — útil para mensajes de error legibles.
func FormatSupportedScreenKeys() string {
	return strings.Join(SupportedScreenKeys(), ", ")
}
