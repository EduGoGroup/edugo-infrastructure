package l4

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ApplyConceptTypes siembra el catalogo conceptual del sistema:
// 5 concept_types (tipos de institucion) y 50 concept_definitions
// (10 terminos por tipo, agrupados en 4 categorias: org, unit,
// member, content).
//
// Idempotente via clause.OnConflict en (id).
//
// Decisiones aplicadas vs legacy (data.go lineas 759-818):
//
//  1. Renombre cosmetico del tipo `workshop` ("Taller / Workshop" =>
//     "Taller"). El sufijo "/ Workshop" duplicaba el code y rompia el
//     tono UI; el code `workshop` ya identifica al tipo.
//  2. Term_keys conservados con notacion dot (`org.name_singular`,
//     `unit.level1`, etc.). NO son snake_case, pero forman un
//     namespace jerarquico i18n-style consistente en los 5 tipos.
//     Cambiar a snake_case (org_name_singular) destruiria la
//     agrupacion semantica que consume el frontend al resolver
//     terminologia. Decision: conservar dot-notation; documentar como
//     convencion de dominio.
//  3. Tipo `language_academy` mantiene sus 10 values en ingles
//     (Level / Class / Term / Student / Teacher / Parent / Course /
//     Test). Es intencional: una academia de idiomas usa terminologia
//     ESL/inglesa. NO uniformizamos a espaniol.
//  4. Tipo `workshop`: `content.subject` = "Taller" colisiona con
//     `org.name_singular` = "Taller". Es semanticamente correcto en
//     ese dominio (la "materia" de un taller es el propio taller).
//     Conservado.
//  5. UUIDs reescritos como `c4xxxxxx-...` (prefijo c4 = concept L4)
//     para distinguir del legacy `c1000000-...` y evitar colisiones
//     accidentales si alguien re-aplica el legacy por error.
//  6. is_active = true para los 5 tipos (no esta en el legacy pero
//     es default del schema; explicitado por claridad).
//
// Inventario: 5 tipos + 50 definiciones (10 x 5). Sin descartes ni
// duplicados detectados en el legacy.
func ApplyConceptTypes(tx *gorm.DB) error {
	if err := applyConceptTypeRows(tx); err != nil {
		return err
	}
	if err := applyConceptDefinitionRows(tx); err != nil {
		return err
	}
	return nil
}

// conceptTypeID encapsula el UUID + code para evitar parseos
// repetidos en la lista de definiciones.
type conceptTypeID struct {
	id   uuid.UUID
	code string
}

var (
	// IDs propios (prefijo c4 para distinguir del legacy `c1...`).
	// El sufijo coincide con el code para legibilidad.
	conceptTypePrimary    = conceptTypeID{id: uuid.MustParse("c4000000-0000-0000-0000-000000000001"), code: "primary_school"}
	conceptTypeHigh       = conceptTypeID{id: uuid.MustParse("c4000000-0000-0000-0000-000000000002"), code: "high_school"}
	conceptTypeLanguage   = conceptTypeID{id: uuid.MustParse("c4000000-0000-0000-0000-000000000003"), code: "language_academy"}
	conceptTypeTechnical  = conceptTypeID{id: uuid.MustParse("c4000000-0000-0000-0000-000000000004"), code: "technical_school"}
	conceptTypeWorkshop   = conceptTypeID{id: uuid.MustParse("c4000000-0000-0000-0000-000000000005"), code: "workshop"}
)

func strPtr(s string) *string { return &s }

// buildL4ConceptTypes retorna las 5 filas de iam.concept_types
// sembradas por L4 como entities.ConceptType. Helper compartido por
// applyConceptTypeRows y por el accessor público l4.ConceptTypes().
func buildL4ConceptTypes() []entities.ConceptType {
	return []entities.ConceptType{
		{
			ID:          conceptTypePrimary.id,
			Name:        "Escuela Primaria",
			Code:        conceptTypePrimary.code,
			Description: strPtr("Institucion de educacion basica"),
			IsActive:    true,
		},
		{
			ID:          conceptTypeHigh.id,
			Name:        "Colegio Secundario",
			Code:        conceptTypeHigh.code,
			Description: strPtr("Institucion de educacion media"),
			IsActive:    true,
		},
		{
			ID:          conceptTypeLanguage.id,
			Name:        "Academia de Idiomas",
			Code:        conceptTypeLanguage.code,
			Description: strPtr("Centro de ensenianza de idiomas"),
			IsActive:    true,
		},
		{
			ID:          conceptTypeTechnical.id,
			Name:        "Instituto Tecnico",
			Code:        conceptTypeTechnical.code,
			Description: strPtr("Formacion tecnica y profesional"),
			IsActive:    true,
		},
		{
			// Renombre vs legacy: "Taller / Workshop" => "Taller".
			ID:          conceptTypeWorkshop.id,
			Name:        "Taller",
			Code:        conceptTypeWorkshop.code,
			Description: strPtr("Cursos cortos y talleres practicos"),
			IsActive:    true,
		},
	}
}

func applyConceptTypeRows(tx *gorm.DB) error {
	rows := buildL4ConceptTypes()
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "code", "description", "is_active", "updated_at"}),
	}).Create(&rows).Error; err != nil {
		return fmt.Errorf("ApplyConceptTypes: upsert concept_types: %w", err)
	}
	return nil
}

// commonTerms define la plantilla de 10 terminos que cada tipo
// instancia con sus propios values. El orden + categories son
// estables entre tipos para que el frontend pueda iterar
// independientemente del tipo.
var commonTermsTemplate = []struct {
	TermKey  string
	Category string
}{
	{"org.name_singular", "org"},
	{"org.name_plural", "org"},
	{"unit.level1", "unit"},
	{"unit.level2", "unit"},
	{"unit.period", "unit"},
	{"member.student", "member"},
	{"member.teacher", "member"},
	{"member.guardian", "member"},
	{"content.subject", "content"},
	{"content.assessment", "content"},
}

// termValuesByType[code] devuelve los 10 values para ese tipo, en
// el mismo orden que commonTermsTemplate. Mantener la longitud == 10.
var termValuesByType = map[string][10]string{
	"primary_school": {
		"Escuela", "Escuelas",
		"Grado", "Clase", "Periodo",
		"Estudiante", "Profesor", "Acudiente",
		"Materia", "Evaluacion",
	},
	"high_school": {
		"Colegio", "Colegios",
		"Anio", "Division", "Trimestre",
		"Alumno", "Docente", "Tutor",
		"Asignatura", "Examen",
	},
	// Decision: language_academy usa terminologia ESL/inglesa
	// intencionalmente (es academia de idiomas).
	"language_academy": {
		"Academia", "Academias",
		"Level", "Class", "Term",
		"Student", "Teacher", "Parent",
		"Course", "Test",
	},
	"technical_school": {
		"Instituto", "Institutos",
		"Semestre", "Seccion", "Cuatrimestre",
		"Aprendiz", "Instructor", "Representante",
		"Modulo", "Prueba",
	},
	// Decision: workshop.content.subject = "Taller" colisiona con
	// org.name_singular = "Taller"; es semanticamente correcto
	// (la "materia" de un taller es el taller mismo).
	"workshop": {
		"Taller", "Talleres",
		"Modulo", "Grupo", "Ciclo",
		"Participante", "Facilitador", "Responsable",
		"Taller", "Ejercicio",
	},
}

func buildDefinitionsFor(typeID uuid.UUID, code string) ([]entities.ConceptDefinition, error) {
	values, ok := termValuesByType[code]
	if !ok {
		return nil, fmt.Errorf("ApplyConceptTypes: missing term values for code %q", code)
	}
	out := make([]entities.ConceptDefinition, 0, len(commonTermsTemplate))
	for i, tmpl := range commonTermsTemplate {
		out = append(out, entities.ConceptDefinition{
			// ID lo genera la BD via default gen_random_uuid().
			ConceptTypeID: typeID,
			TermKey:       tmpl.TermKey,
			TermValue:     values[i],
			Category:      tmpl.Category,
			SortOrder:     i + 1,
		})
	}
	return out, nil
}

// buildL4ConceptDefinitions retorna las 50 filas (5 tipos × 10
// términos) de iam.concept_definitions sembradas por L4. Helper
// compartido por applyConceptDefinitionRows y por el accessor público
// l4.ConceptDefinitions().
//
// Nota: el `ID` se deja en uuid.Nil porque el schema lo genera con
// `gen_random_uuid()`. El loader del cross-checker no necesita el ID
// exacto — usa la natural key (concept_type_id, term_key) para
// deduplicar contra legacy, así que el accessor sintetiza un UUID
// determinístico (NewSHA1 sobre concept_type_id+":"+term_key) para que
// dos llamadas consecutivas devuelvan slices estables.
func buildL4ConceptDefinitions() ([]entities.ConceptDefinition, error) {
	allTypes := []conceptTypeID{
		conceptTypePrimary,
		conceptTypeHigh,
		conceptTypeLanguage,
		conceptTypeTechnical,
		conceptTypeWorkshop,
	}

	rows := make([]entities.ConceptDefinition, 0, len(allTypes)*len(commonTermsTemplate))
	for _, t := range allTypes {
		defs, err := buildDefinitionsFor(t.id, t.code)
		if err != nil {
			return nil, err
		}
		rows = append(rows, defs...)
	}
	return rows, nil
}

func applyConceptDefinitionRows(tx *gorm.DB) error {
	rows, err := buildL4ConceptDefinitions()
	if err != nil {
		return err
	}

	// Upsert por la unica natural key disponible
	// (concept_type_id, term_key) — `id` se genera en cada Create()
	// y por tanto NO sirve como llave de conflicto idempotente.
	// El schema garantiza UNIQUE(concept_type_id, term_key) via
	// concept_definitions_type_key_unique.
	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "concept_type_id"},
			{Name: "term_key"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"term_value", "category", "sort_order", "updated_at"}),
	}).Create(&rows).Error; err != nil {
		return fmt.Errorf("ApplyConceptTypes: upsert concept_definitions: %w", err)
	}
	return nil
}
