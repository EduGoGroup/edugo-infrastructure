// Package focal_evaluacion compone sobre focal_pantalla agregando datos del
// dominio de evaluaciones (assessments + questions + question_options) y los
// grants necesarios para que los roles de focal_pantalla puedan operar el
// flujo completo de evaluación.
//
// Convención (project_edugo_playgrounds_convention en memoria): los playgrounds
// son fotos inmutables y se acumulan. NO duplicamos aquí la creación de
// users/roles/school/units/memberships — eso lo hace focal_pantalla. Este
// paquete asume que focal_pantalla ya corrió (o correrá antes vía registry
// `all`) y se limita a añadir lo nuevo.
//
// Asume de focal_pantalla:
//   - schoolID  = "61000000-0000-0000-0000-000000000003"
//   - unitID    = "61000000-0000-0000-0000-000000000004"
//   - adminUserID = "61000000-0000-0000-0000-000000000001"
//   - roles: admin (11000000-...-0001), viewer (11000000-...-0002), author (11000000-...-0003)
//
// Decisión de grants (wildcard-first, ver feedback_wildcard_first en memoria):
// focal_pantalla ya otorgó `*.read` global a los tres roles. Eso cubre
// `content.assessments.read` automáticamente, así que NO lo repetimos por
// rol. Solo agregamos lo que NO está cubierto por `*.read`:
//   - admin  → `content.assessments.*` (CRUD completo)
//   - viewer → nada nuevo (ya tiene *.read; assessments.read cubierto)
//   - author → `content.assessments.create` (crear sin editar/eliminar)
//
// Rangos de UUIDs propios (separados de focal_pantalla para no pisar):
//   - role_grants:      61000001-0000-0000-0000-0000000000XX
//   - assessments:      61000001-0000-0000-0000-00000000000X (offset 1..3)
//   - questions:        61000001-0000-0000-0000-0000000XYYYY (X=assessment, YYYY=sort)
//   - question_options: 61000001-0000-0000-0000-000000XYYYZ (Z=option index)
//
// Lo que siembra:
//  1. iam.role_grants                — 2 grants nuevos (admin + author).
//  2. assessment.assessment          — 3 evaluaciones en la única escuela.
//  3. assessment.questions           — 4 preguntas por evaluación (12 total).
//  4. assessment.question_options    — 4 options para cada multiple_choice
//                                       (2 multiple_choice × 3 assessments × 4 = 24).
//
// NO siembra assessments_assignments ni assessment_attempts (queda para una
// iteración posterior si el flujo lo pide).
//
// Idempotente: todas las inserciones usan OnConflict DoNothing por id.
package focal_evaluacion

import (
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// Constantes heredadas de focal_pantalla. Las redeclaramos como
	// privadas para no introducir dependencia de paquete; el contrato es:
	// si focal_pantalla cambia un UUID, hay que reflejarlo acá.
	schoolID    = "61000000-0000-0000-0000-000000000003"
	adminUserID = "61000000-0000-0000-0000-000000000001"

	adminRoleID  = "11000000-0000-0000-0000-000000000001"
	authorRoleID = "11000000-0000-0000-0000-000000000003"

	// Patterns nuevos a aplicar. viewer no necesita nada — su `*.read`
	// global ya cubre `content.assessments.read`.
	adminPatternCrud    = "content.assessments.*"
	authorPatternCreate = "content.assessments.create"

	// IDs determinísticos de assessments (rango 61000001-...-00000001..03).
	assessment1ID = "61000001-0000-0000-0000-000000000001"
	assessment2ID = "61000001-0000-0000-0000-000000000002"
	assessment3ID = "61000001-0000-0000-0000-000000000003"
)

// Apply siembra el playground focal_evaluacion. Asume que focal_pantalla ya
// corrió (provee users/roles/school/units). Idempotente.
func Apply(tx *gorm.DB) error {
	if err := upsertAdminAssessmentsGrant(tx); err != nil {
		return fmt.Errorf("playground/focal_evaluacion: admin_grant: %w", err)
	}
	if err := upsertAuthorAssessmentsGrant(tx); err != nil {
		return fmt.Errorf("playground/focal_evaluacion: author_grant: %w", err)
	}
	if err := upsertAssessments(tx); err != nil {
		return fmt.Errorf("playground/focal_evaluacion: assessments: %w", err)
	}
	if err := upsertQuestionsAndOptions(tx); err != nil {
		return fmt.Errorf("playground/focal_evaluacion: questions: %w", err)
	}
	return nil
}

// upsertAdminAssessmentsGrant otorga `content.assessments.*` al rol admin
// de focal_pantalla. Su `*.read` global ya cubre lectura; este grant
// completa create/update/delete/publish/grade/etc.
func upsertAdminAssessmentsGrant(tx *gorm.DB) error {
	return upsertRoleGrants(tx, adminRoleID, []string{adminPatternCrud})
}

// upsertAuthorAssessmentsGrant otorga solo `content.assessments.create`
// al rol author. La lectura ya está cubierta por `*.read` global. Sin
// update/delete — paralelo al modelo del author de focal_pantalla con
// anuncios (crea sin poder modificar).
func upsertAuthorAssessmentsGrant(tx *gorm.DB) error {
	return upsertRoleGrants(tx, authorRoleID, []string{authorPatternCreate})
}

// upsertRoleGrants inserta una lista de patterns como allow-grants para un
// rol. ID determinístico SHA1(role_id:pattern:effect) — misma fórmula que
// focal_pantalla — para que el reseed sea idempotente y un mismo pattern
// no se duplique aunque cambiemos el namespace de UUIDs.
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

// upsertAssessments crea 3 evaluaciones variadas sobre la única escuela
// del playground. La variedad cubre: status (draft vs published), timed
// vs no timed, distintos pass_threshold, distintos max_attempts.
//
// CHECK constraints relevantes en assessment.assessment (ver entity):
//   - source_type ∈ ('manual','ai_generated')
//   - pass_threshold ∈ [0,100]
//   - status ∈ ('draft','generated','published','archived','closed')
//   - available_until > available_from (cuando ambos != null)
func upsertAssessments(tx *gorm.DB) error {
	sid, err := uuid.Parse(schoolID)
	if err != nil {
		return err
	}
	cid, err := uuid.Parse(adminUserID)
	if err != nil {
		return err
	}

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

	titleMath := "Examen Diagnóstico de Matemáticas"
	titleSolar := "Quiz Sobre Sistema Solar"
	titleRead := "Evaluación Final de Lectura"

	items := []entities.Assessment{
		{
			ID:                 mustParseUUID(assessment1ID),
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
			ID:                 mustParseUUID(assessment2ID),
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
			ID:                 mustParseUUID(assessment3ID),
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

// questionSpec describe una pregunta en código antes de materializar
// el struct entities.Question. Mantiene el código de upsertQuestionsAndOptions
// declarativo y compacto.
type questionSpec struct {
	idSuffix     string // sufijo del UUID (rango 61000001-...-0000XYYYY)
	sortOrder    int
	text         string
	qType        string  // multiple_choice | true_false | short_answer | open_ended
	correct      string  // valor textual; para multiple_choice es el option_text correcto
	explanation  string
	points       float64
	difficulty   string // easy | medium | hard
	options      []string
}

// upsertQuestionsAndOptions siembra 4 preguntas por cada uno de los 3
// assessments (12 total), con mezcla de tipos: 2 multiple_choice + 1
// true_false + 1 short_answer.
//
// CHECK constraints relevantes en assessment.questions (ver entity):
//   - question_type ∈ ('multiple_choice','true_false','short_answer','open_ended')
//   - difficulty ∈ ('easy','medium','hard') si != null
//
// Para multiple_choice insertamos 4 question_options con sort_order 1..4.
// El `correct_answer` se guarda como texto literal — la entidad no exige
// un FK a la option correcta; la app deserializa y compara.
func upsertQuestionsAndOptions(tx *gorm.DB) error {
	type bundle struct {
		assessmentID  string
		assessmentIdx int // 1..3, para namespace de UUIDs
		questions     []questionSpec
	}

	bundles := []bundle{
		{
			assessmentID:  assessment1ID,
			assessmentIdx: 1,
			questions: []questionSpec{
				{
					idSuffix:    "00010001",
					sortOrder:   1,
					text:        "¿Cuánto es 7 × 8?",
					qType:       "multiple_choice",
					correct:     "56",
					explanation: "7 × 8 = 56 (tabla del 7).",
					points:      10,
					difficulty:  "easy",
					options:     []string{"54", "56", "58", "64"},
				},
				{
					idSuffix:    "00010002",
					sortOrder:   2,
					text:        "¿Cuál de las siguientes es una fracción equivalente a 1/2?",
					qType:       "multiple_choice",
					correct:     "2/4",
					explanation: "Multiplicando numerador y denominador por 2 se obtiene 2/4.",
					points:      15,
					difficulty:  "medium",
					options:     []string{"1/3", "2/4", "3/5", "1/4"},
				},
				{
					idSuffix:    "00010003",
					sortOrder:   3,
					text:        "El cero es un número par.",
					qType:       "true_false",
					correct:     "true",
					explanation: "Cero es divisible por 2 sin residuo, por lo tanto es par.",
					points:      10,
					difficulty:  "easy",
				},
				{
					idSuffix:    "00010004",
					sortOrder:   4,
					text:        "Escribe el resultado de 12² (doce al cuadrado).",
					qType:       "short_answer",
					correct:     "144",
					explanation: "12 × 12 = 144.",
					points:      20,
					difficulty:  "medium",
				},
			},
		},
		{
			assessmentID:  assessment2ID,
			assessmentIdx: 2,
			questions: []questionSpec{
				{
					idSuffix:    "00020001",
					sortOrder:   1,
					text:        "¿Cuál es el planeta más cercano al Sol?",
					qType:       "multiple_choice",
					correct:     "Mercurio",
					explanation: "Mercurio es el primer planeta del sistema solar.",
					points:      10,
					difficulty:  "easy",
					options:     []string{"Venus", "Mercurio", "Marte", "La Tierra"},
				},
				{
					idSuffix:    "00020002",
					sortOrder:   2,
					text:        "¿Cuántos planetas tiene actualmente el sistema solar?",
					qType:       "multiple_choice",
					correct:     "8",
					explanation: "Plutón fue reclasificado como planeta enano en 2006.",
					points:      15,
					difficulty:  "medium",
					options:     []string{"7", "8", "9", "10"},
				},
				{
					idSuffix:    "00020003",
					sortOrder:   3,
					text:        "Júpiter es el planeta más grande del sistema solar.",
					qType:       "true_false",
					correct:     "true",
					explanation: "Júpiter supera en masa a todos los demás planetas juntos.",
					points:      10,
					difficulty:  "easy",
				},
				{
					idSuffix:    "00020004",
					sortOrder:   4,
					text:        "Nombra el satélite natural de la Tierra.",
					qType:       "short_answer",
					correct:     "Luna",
					explanation: "La Tierra tiene un único satélite natural: la Luna.",
					points:      20,
					difficulty:  "easy",
				},
			},
		},
		{
			assessmentID:  assessment3ID,
			assessmentIdx: 3,
			questions: []questionSpec{
				{
					idSuffix:    "00030001",
					sortOrder:   1,
					text:        "¿Qué tipo de texto presenta hechos verificables sin opinión del autor?",
					qType:       "multiple_choice",
					correct:     "Informativo",
					explanation: "El texto informativo busca presentar datos objetivos.",
					points:      10,
					difficulty:  "easy",
					options:     []string{"Narrativo", "Informativo", "Lírico", "Dramático"},
				},
				{
					idSuffix:    "00030002",
					sortOrder:   2,
					text:        "¿Cuál es la función principal del texto argumentativo?",
					qType:       "multiple_choice",
					correct:     "Persuadir al lector",
					explanation: "El texto argumentativo busca convencer mediante razones.",
					points:      15,
					difficulty:  "medium",
					options:     []string{"Entretener", "Informar", "Persuadir al lector", "Describir lugares"},
				},
				{
					idSuffix:    "00030003",
					sortOrder:   3,
					text:        "Una metáfora compara dos elementos usando la palabra \"como\".",
					qType:       "true_false",
					correct:     "false",
					explanation: "Esa es la definición de símil; la metáfora compara sin nexo.",
					points:      10,
					difficulty:  "medium",
				},
				{
					idSuffix:    "00030004",
					sortOrder:   4,
					text:        "Define con tus palabras qué es el tema principal de un texto.",
					qType:       "short_answer",
					correct:     "Idea central que el autor desarrolla a lo largo del texto.",
					explanation: "Se acepta cualquier respuesta que mencione \"idea central\" o \"asunto principal\".",
					points:      20,
					difficulty:  "hard",
				},
			},
		},
	}

	for _, b := range bundles {
		aid, err := uuid.Parse(b.assessmentID)
		if err != nil {
			return err
		}
		for _, qs := range b.questions {
			qid := mustParseUUID("61000001-0000-0000-0000-0000" + qs.idSuffix)

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
// IDs determinísticos: se agrega un sufijo del índice de option (1..4)
// al idSuffix de la question. Como question idSuffix tiene 8 chars
// ("000X0YYY"), las options usan 9 (concat con el índice 1..4).
func upsertQuestionOptions(tx *gorm.DB, questionID uuid.UUID, qIDSuffix string, options []string) error {
	for i, optText := range options {
		// Opt ID: 61000001-0000-0000-0000-000<qIDSuffix><i+1>
		// qIDSuffix es de 8 chars; i+1 ocupa 1 char (1..4). Total = 12 chars
		// → encajan en los 12 chars finales del UUID (000000000000).
		// Compongo "000" + qIDSuffix (8) + str(i+1) (1) = 12 chars.
		optIDStr := fmt.Sprintf("61000001-0000-0000-0000-000%s%d", qIDSuffix, i+1)
		oid := mustParseUUID(optIDStr)
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

// mustParseUUID es helper interno: convierte un literal UUID en uuid.UUID
// o entra en pánico si el literal es inválido. Solo se usa para constantes
// hardcodeadas dentro de este paquete, donde el formato está bajo control.
func mustParseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		panic(fmt.Sprintf("focal_evaluacion: UUID inválido %q: %v", s, err))
	}
	return id
}
