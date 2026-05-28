// Package focal_evaluacion_v2 es el playground v2 para validar la
// arquitectura SDUI post-refactor (composer defaults+added/removed +
// master-detail-v1 template + scope split form-submit/resource-toolbar).
// Snapshot 002 documenta el plan de re-introducción de Fases 1-4 del
// refactor; este playground es el dataset focal para la Fase 3-4 de
// validación end-to-end.
//
// Diferencias vs focal_evaluacion (v1):
//  1. Standalone: NO compone sobre focal_pantalla. Siembra sus propios
//     users (focal-v2-admin/viewer/author), school, unit, roles y
//     memberships. v1 reutilizaba el rol y la escuela de focal_pantalla;
//     v2 los aísla para no contaminar la foto pre-refactor.
//  2. Wildcard-first en grants: cada rol recibe directamente
//     `content.assessments.*` (admin/author CRUD completo cubre
//     publish/archive/view-questions automáticamente — ver
//     feedback_wildcard_first). viewer recibe solo `*.read`.
//  3. UUIDs en rango 62000000-...: separado de 60..., 61... usados por
//     admin y focal_pantalla. Convive sin colisión.
//
// Convención (project_edugo_playgrounds_convention en memoria): los
// playgrounds son fotos inmutables y se acumulan. focal_evaluacion v1
// queda INTOCADO como snapshot de la arquitectura pre-refactor; v2 es
// una foto NUEVA paralela, no una edición.
//
// Asume que las capas L0..L4 corrieron (resources, roles, screen_
// templates, screen_instances, resource_screens ya existen). El
// playground solo siembra datos del dominio + grants específicos del
// rol de prueba.
//
// Lo que siembra:
//  1. academic.schools         — 1 escuela (FOCAL-EVAL-V2).
//  2. academic.academic_units  — 1 unidad raíz.
//  3. iam.roles                — 3 roles scope=school (admin, viewer, author).
//  4. iam.role_grants          — admin: content.assessments.* + *.read;
//                                viewer: *.read; author: *.read + content.assessments.create.
//  5. auth.users               — 3 usuarios focal-v2-{admin,viewer,author}@edugo.local.
//  6. iam.user_roles           — assignments user × rol.
//  7. academic.memberships     — admin, viewer y author en la misma escuela/unidad.
//  8. assessment.assessments   — 3 evaluaciones variadas (status, timed, thresholds distintos).
//  9. assessment.questions     — 4 preguntas por evaluación (12 total).
// 10. assessment.question_options — 4 options por cada multiple_choice (24 total).
//
// Idempotente: todas las inserciones usan OnConflict DoNothing por id.
package focal_evaluacion_v2

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// Credenciales del playground v2. Password compartido para
	// simplicidad — fixture dev, no usuario real.
	AdminEmail  = "focal-v2-admin@edugo.local"
	ViewerEmail = "focal-v2-viewer@edugo.local"
	AuthorEmail = "focal-v2-author@edugo.local"
	Password    = "12345678"

	// Rango UUID 62000000-...: reservado para focal_evaluacion_v2.
	// Separado del 61000000-... (focal_pantalla) y 60000000-... (admin).
	schoolID     = "62000000-0000-0000-0000-000000000001"
	unitID       = "62000000-0000-0000-0000-000000000002"
	adminUserID  = "62000000-0000-0000-0000-000000000003"
	viewerUserID = "62000000-0000-0000-0000-000000000004"
	authorUserID = "62000000-0000-0000-0000-000000000005"
	adminMembID  = "62000000-0000-0000-0000-000000000006"
	viewerMembID = "62000000-0000-0000-0000-000000000007"
	authorMembID = "62000000-0000-0000-0000-000000000008"

	// IDs de roles en rango 12000000-... (separado del 11000000-...
	// usado por focal_pantalla).
	adminRoleID  = "12000000-0000-0000-0000-000000000001"
	viewerRoleID = "12000000-0000-0000-0000-000000000002"
	authorRoleID = "12000000-0000-0000-0000-000000000003"

	adminRoleName  = "focal_evaluacion_v2_admin"
	viewerRoleName = "focal_evaluacion_v2_viewer"
	authorRoleName = "focal_evaluacion_v2_author"

	// Patterns wildcard-first (feedback_wildcard_first).
	// admin: CRUD completo + lectura general (publish/archive/etc.
	// cubiertos por content.assessments.*).
	adminPatternCrud = "content.assessments.*"
	adminPatternRead = "*.read"
	// viewer: solo lectura global.
	viewerPatternRead = "*.read"
	// author: lectura general + creación específica (no update/delete).
	authorPatternRead   = "*.read"
	authorPatternCreate = "content.assessments.create"

	schoolCode = "FOCAL-EVAL-V2"
	schoolName = "Escuela Focal Evaluación V2"
	unitCode   = "FOCAL-EVAL-V2-MAIN"
	unitName   = "Sede Única V2"

	academicYear = 2026

	// IDs determinísticos de assessments en rango 62000010-...
	assessment1ID = "62000010-0000-0000-0000-000000000001"
	assessment2ID = "62000010-0000-0000-0000-000000000002"
	assessment3ID = "62000010-0000-0000-0000-000000000003"
)

// Apply siembra el playground focal_evaluacion_v2. Asume que L0..L4
// corrieron. Standalone — no depende de focal_pantalla. Idempotente.
func Apply(tx *gorm.DB) error {
	// 1. School.
	if err := common.SeedSchool(tx, common.SchoolSpec{
		ID:   common.MustParseUUID(schoolID),
		Name: schoolName,
		Code: schoolCode,
	}); err != nil {
		return fmt.Errorf("playground/focal_evaluacion_v2: school: %w", err)
	}

	// 2. Academic unit (sin helper en common — entity-específica).
	if err := upsertAcademicUnit(tx); err != nil {
		return fmt.Errorf("playground/focal_evaluacion_v2: academic_unit: %w", err)
	}

	// 3. Roles (admin/viewer/author scope=school). Wildcard-first:
	// el catálogo de permisos no se enumera por rol; los grants más
	// abajo apuntan a patterns amplios.
	roles := []common.RoleSpec{
		{
			ID:          common.MustParseUUID(adminRoleID),
			Name:        adminRoleName,
			DisplayName: "Admin Evaluación V2 — CRUD",
			Description: "CRUD completo sobre evaluaciones. Playground focal-evaluacion v2.",
		},
		{
			ID:          common.MustParseUUID(viewerRoleID),
			Name:        viewerRoleName,
			DisplayName: "Viewer Evaluación V2 — solo lectura",
			Description: "Solo lectura. Playground focal-evaluacion v2.",
		},
		{
			ID:          common.MustParseUUID(authorRoleID),
			Name:        authorRoleName,
			DisplayName: "Author Evaluación V2 — crea sin editar",
			Description: "Crea evaluaciones sin poder modificarlas ni eliminarlas. Playground focal-evaluacion v2.",
		},
	}
	for _, r := range roles {
		if err := common.SeedRole(tx, r); err != nil {
			return fmt.Errorf("playground/focal_evaluacion_v2: roles: %w", err)
		}
	}

	// 4. Role grants. Wildcard-first: admin usa `content.assessments.*`
	// que cubre publish/archive/view-questions/view-results
	// automáticamente. Sin enumeración atómica.
	grantSpecs := []struct {
		roleID   string
		patterns []string
	}{
		{adminRoleID, []string{adminPatternCrud, adminPatternRead}},
		{viewerRoleID, []string{viewerPatternRead}},
		{authorRoleID, []string{authorPatternRead, authorPatternCreate}},
	}
	for _, gs := range grantSpecs {
		rid := common.MustParseUUID(gs.roleID)
		for _, pattern := range gs.patterns {
			if err := common.SeedRoleGrant(tx, rid, pattern); err != nil {
				return fmt.Errorf("playground/focal_evaluacion_v2: role_grants: %w", err)
			}
		}
	}

	// 5. Users.
	userSpecs := []common.UserSpec{
		{ID: common.MustParseUUID(adminUserID), Email: AdminEmail, Password: Password, FirstName: "Admin", LastName: "EvalV2"},
		{ID: common.MustParseUUID(viewerUserID), Email: ViewerEmail, Password: Password, FirstName: "Viewer", LastName: "EvalV2"},
		{ID: common.MustParseUUID(authorUserID), Email: AuthorEmail, Password: Password, FirstName: "Author", LastName: "EvalV2"},
	}
	for _, us := range userSpecs {
		if err := common.SeedUser(tx, us); err != nil {
			return fmt.Errorf("playground/focal_evaluacion_v2: users: %w", err)
		}
	}

	// 6. User-role assignments.
	userRolePairs := [][2]string{
		{adminUserID, adminRoleID},
		{viewerUserID, viewerRoleID},
		{authorUserID, authorRoleID},
	}
	for _, p := range userRolePairs {
		if err := common.SeedUserRole(tx, common.MustParseUUID(p[0]), common.MustParseUUID(p[1])); err != nil {
			return fmt.Errorf("playground/focal_evaluacion_v2: user_roles: %w", err)
		}
	}

	// 7. Memberships. AcademicUnitID por puntero para que el login
	// resuelva contexto completo sin switch-context.
	auid := common.MustParseUUID(unitID)
	sid := common.MustParseUUID(schoolID)
	membSpecs := []common.MembershipSpec{
		{ID: common.MustParseUUID(adminMembID), UserID: common.MustParseUUID(adminUserID), SchoolID: sid, AcademicUnitID: &auid, Role: "admin"},
		{ID: common.MustParseUUID(viewerMembID), UserID: common.MustParseUUID(viewerUserID), SchoolID: sid, AcademicUnitID: &auid, Role: "teacher"},
		{ID: common.MustParseUUID(authorMembID), UserID: common.MustParseUUID(authorUserID), SchoolID: sid, AcademicUnitID: &auid, Role: "teacher"},
	}
	for _, ms := range membSpecs {
		if err := common.SeedMembership(tx, ms); err != nil {
			return fmt.Errorf("playground/focal_evaluacion_v2: memberships: %w", err)
		}
	}

	// 8-10. Assessments + questions + question_options.
	if err := upsertAssessments(tx); err != nil {
		return fmt.Errorf("playground/focal_evaluacion_v2: assessments: %w", err)
	}
	if err := upsertQuestionsAndOptions(tx); err != nil {
		return fmt.Errorf("playground/focal_evaluacion_v2: questions: %w", err)
	}
	return nil
}

func upsertAcademicUnit(tx *gorm.DB) error {
	u := entities.AcademicUnit{
		ID:           common.MustParseUUID(unitID),
		SchoolID:     common.MustParseUUID(schoolID),
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

// upsertAssessments crea 3 evaluaciones variadas. Mismo dataset
// conceptual que v1 — la idea es que la pantalla luzca exactamente
// igual que en v1 (mismo número de filas, columnas pobladas) pero
// sobre la arquitectura nueva.
//
// CHECK constraints en assessment.assessment:
//   - source_type ∈ ('manual','ai_generated')
//   - pass_threshold ∈ [0,100]
//   - status ∈ ('draft','generated','published','archived','closed')
//   - available_until > available_from (cuando ambos != null)
func upsertAssessments(tx *gorm.DB) error {
	sid := common.MustParseUUID(schoolID)
	cid := common.MustParseUUID(adminUserID)

	now := time.Now().UTC()
	availableFrom := now
	availableUntil := now.Add(30 * 24 * time.Hour)

	descMath := "Evaluación diagnóstica para medir nivel inicial en matemáticas."
	descSolar := "Quiz breve para repasar conceptos del sistema solar antes de la unidad."
	descRead := "Evaluación final de comprensión lectora del segundo trimestre."

	passMath := 60.0
	passSolar := 70.0
	passRead := 75.0

	timeLimitMath := 60.0
	timeLimitRead := 45.0

	maxMath := 3
	maxSolar := 1
	maxRead := 2

	titleMath := "Examen Diagnóstico de Matemáticas (V2)"
	titleSolar := "Quiz Sobre Sistema Solar (V2)"
	titleRead := "Evaluación Final de Lectura (V2)"

	items := []entities.Assessment{
		{
			ID:                 common.MustParseUUID(assessment1ID),
			SourceType:         "manual",
			SchoolID:           &sid,
			CreatedByUserID:    &cid,
			QuestionsCount:     4,
			Title:              &titleMath,
			Description:        &descMath,
			PassThreshold:      &passMath,
			MaxAttempts:        &maxMath,
			TimeLimitMinutes:   &timeLimitMath,
			IsTimed:            true,
			ShuffleQuestions:   true,
			ShowCorrectAnswers: true,
			AvailableFrom:      &availableFrom,
			AvailableUntil:     &availableUntil,
			Status:             "published",
		},
		{
			ID:                 common.MustParseUUID(assessment2ID),
			SourceType:         "manual",
			SchoolID:           &sid,
			CreatedByUserID:    &cid,
			QuestionsCount:     4,
			Title:              &titleSolar,
			Description:        &descSolar,
			PassThreshold:      &passSolar,
			MaxAttempts:        &maxSolar,
			TimeLimitMinutes:   nil,
			IsTimed:            false,
			ShuffleQuestions:   false,
			ShowCorrectAnswers: true,
			AvailableFrom:      &availableFrom,
			AvailableUntil:     &availableUntil,
			Status:             "draft",
		},
		{
			ID:                 common.MustParseUUID(assessment3ID),
			SourceType:         "manual",
			SchoolID:           &sid,
			CreatedByUserID:    &cid,
			QuestionsCount:     4,
			Title:              &titleRead,
			Description:        &descRead,
			PassThreshold:      &passRead,
			MaxAttempts:        &maxRead,
			TimeLimitMinutes:   &timeLimitRead,
			IsTimed:            true,
			ShuffleQuestions:   false,
			ShowCorrectAnswers: false,
			AvailableFrom:      &availableFrom,
			AvailableUntil:     &availableUntil,
			Status:             "published",
		},
	}

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&items).Error
}

// questionSpec describe una pregunta antes de materializar la entity.
// Mantiene el código de upsertQuestionsAndOptions declarativo.
type questionSpec struct {
	idSuffix    string // sufijo del UUID (rango 62000020-...-0000XYYYY)
	sortOrder   int
	text        string
	qType       string // multiple_choice | true_false | short_answer | open_ended
	correct     string
	explanation string
	points      float64
	difficulty  string
	options     []string
}

// upsertQuestionsAndOptions siembra 4 preguntas por assessment (12
// total) con mezcla de tipos. Para multiple_choice agrega 4
// question_options (24 total).
//
// CHECK constraints en assessment.questions:
//   - question_type ∈ ('multiple_choice','true_false','short_answer','open_ended')
//   - difficulty ∈ ('easy','medium','hard') si != null
func upsertQuestionsAndOptions(tx *gorm.DB) error {
	type bundle struct {
		assessmentID string
		questions    []questionSpec
	}

	bundles := []bundle{
		{
			assessmentID: assessment1ID,
			questions: []questionSpec{
				{
					idSuffix: "00010001", sortOrder: 1,
					text:        "¿Cuánto es 7 × 8?",
					qType:       "multiple_choice",
					correct:     "56",
					explanation: "7 × 8 = 56 (tabla del 7).",
					points:      10, difficulty: "easy",
					options: []string{"54", "56", "58", "64"},
				},
				{
					idSuffix: "00010002", sortOrder: 2,
					text:        "¿Cuál de las siguientes es una fracción equivalente a 1/2?",
					qType:       "multiple_choice",
					correct:     "2/4",
					explanation: "Multiplicando numerador y denominador por 2 se obtiene 2/4.",
					points:      15, difficulty: "medium",
					options: []string{"1/3", "2/4", "3/5", "1/4"},
				},
				{
					idSuffix: "00010003", sortOrder: 3,
					text:        "El cero es un número par.",
					qType:       "true_false",
					correct:     "true",
					explanation: "Cero es divisible por 2 sin residuo, por lo tanto es par.",
					points:      10, difficulty: "easy",
				},
				{
					idSuffix: "00010004", sortOrder: 4,
					text:        "Escribe el resultado de 12² (doce al cuadrado).",
					qType:       "short_answer",
					correct:     "144",
					explanation: "12 × 12 = 144.",
					points:      20, difficulty: "medium",
				},
			},
		},
		{
			assessmentID: assessment2ID,
			questions: []questionSpec{
				{
					idSuffix: "00020001", sortOrder: 1,
					text:        "¿Cuál es el planeta más cercano al Sol?",
					qType:       "multiple_choice",
					correct:     "Mercurio",
					explanation: "Mercurio es el primer planeta del sistema solar.",
					points:      10, difficulty: "easy",
					options: []string{"Venus", "Mercurio", "Marte", "La Tierra"},
				},
				{
					idSuffix: "00020002", sortOrder: 2,
					text:        "¿Cuántos planetas tiene actualmente el sistema solar?",
					qType:       "multiple_choice",
					correct:     "8",
					explanation: "Plutón fue reclasificado como planeta enano en 2006.",
					points:      15, difficulty: "medium",
					options: []string{"7", "8", "9", "10"},
				},
				{
					idSuffix: "00020003", sortOrder: 3,
					text:        "Júpiter es el planeta más grande del sistema solar.",
					qType:       "true_false",
					correct:     "true",
					explanation: "Júpiter supera en masa a todos los demás planetas juntos.",
					points:      10, difficulty: "easy",
				},
				{
					idSuffix: "00020004", sortOrder: 4,
					text:        "Nombra el satélite natural de la Tierra.",
					qType:       "short_answer",
					correct:     "Luna",
					explanation: "La Tierra tiene un único satélite natural: la Luna.",
					points:      20, difficulty: "easy",
				},
			},
		},
		{
			assessmentID: assessment3ID,
			questions: []questionSpec{
				{
					idSuffix: "00030001", sortOrder: 1,
					text:        "¿Qué tipo de texto presenta hechos verificables sin opinión del autor?",
					qType:       "multiple_choice",
					correct:     "Informativo",
					explanation: "El texto informativo busca presentar datos objetivos.",
					points:      10, difficulty: "easy",
					options: []string{"Narrativo", "Informativo", "Lírico", "Dramático"},
				},
				{
					idSuffix: "00030002", sortOrder: 2,
					text:        "¿Cuál es la función principal del texto argumentativo?",
					qType:       "multiple_choice",
					correct:     "Persuadir al lector",
					explanation: "El texto argumentativo busca convencer mediante razones.",
					points:      15, difficulty: "medium",
					options: []string{"Entretener", "Informar", "Persuadir al lector", "Describir lugares"},
				},
				{
					idSuffix: "00030003", sortOrder: 3,
					text:        "Una metáfora compara dos elementos usando la palabra \"como\".",
					qType:       "true_false",
					correct:     "false",
					explanation: "Esa es la definición de símil; la metáfora compara sin nexo.",
					points:      10, difficulty: "medium",
				},
				{
					idSuffix: "00030004", sortOrder: 4,
					text:        "Define con tus palabras qué es el tema principal de un texto.",
					qType:       "short_answer",
					correct:     "Idea central que el autor desarrolla a lo largo del texto.",
					explanation: "Se acepta cualquier respuesta que mencione \"idea central\" o \"asunto principal\".",
					points:      20, difficulty: "hard",
				},
			},
		},
	}

	for _, b := range bundles {
		aid := common.MustParseUUID(b.assessmentID)
		for _, qs := range b.questions {
			qid := common.MustParseUUID("62000020-0000-0000-0000-0000" + qs.idSuffix)
			correct := qs.correct
			expl := qs.explanation
			diff := qs.difficulty
			q := entities.Question{
				ID:            qid,
				AssessmentID:  aid,
				SortOrder:     qs.sortOrder,
				QuestionText:  qs.text,
				QuestionType:  qs.qType,
				CorrectAnswer: &correct,
				Explanation:   &expl,
				Points:        qs.points,
				Difficulty:    &diff,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoNothing: true,
			}).Create(&q).Error; err != nil {
				return fmt.Errorf("question %s: %w", qid, err)
			}
			if qs.qType == "multiple_choice" && len(qs.options) > 0 {
				if err := upsertQuestionOptions(tx, qid, qs.idSuffix, qs.options); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// upsertQuestionOptions siembra las 4 opciones de una multiple_choice.
// IDs determinísticos: 62000020-0000-0000-0000-000<qIDSuffix><index>.
func upsertQuestionOptions(tx *gorm.DB, questionID uuid.UUID, qIDSuffix string, options []string) error {
	for i, optText := range options {
		optIDStr := fmt.Sprintf("62000020-0000-0000-0000-000%s%d", qIDSuffix, i+1)
		oid := common.MustParseUUID(optIDStr)
		opt := entities.QuestionOption{
			ID:         oid,
			QuestionID: questionID,
			OptionText: optText,
			SortOrder:  i + 1,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).Create(&opt).Error; err != nil {
			return fmt.Errorf("option %s: %w", oid, err)
		}
	}
	return nil
}
