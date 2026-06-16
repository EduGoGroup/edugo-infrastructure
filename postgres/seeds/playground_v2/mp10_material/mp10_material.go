// Package mp10_material es un fixture de la línea v2 que siembra MATERIAL
// PUBLICADO (content.materials) para las materias de los hijos del representante
// del mundo `base`, para validar E2E la pantalla "Material del hijo" (plan 024
// micro-plan M3).
//
// Contexto: en `base` la tabla content.materials queda VACÍA, así que la lista
// de material del hijo no tiene nada que mostrar. Este fixture compone ENCIMA de
// `base`: invoca primero base.Apply (idempotente: trunca y resiembra el mundo de
// dev) y reutiliza sus escuelas, unidades, materias, membresías e inscripciones
// tal cual; luego agrega filas a content.materials. Por eso
// `make docker-playground-v2 P=mp10_material` deja un mundo completo y
// consistente — el migrator aplica `system` + ESTE fixture, que a su vez siembra
// `base` adentro (el runner NO encadena `base` automáticamente).
//
// Contrato del lector (academic, postgresWardMaterialRepository.ListByStudent):
// un material aparece en "Material del hijo" si y solo si
//   - status = 'ready'  (NO 'published': ese valor NO existe en el CHECK de
//     content.materials — status ∈ {draft,uploaded,processing,ready,failed}; el
//     lector filtra exactamente por 'ready'),
//   - deleted_at IS NULL,
//   - subject_id ∈ las materias en las que el hijo está inscrito
//     (subject_offering_enrollments) en su unidad activa con membership active.
//
// NO filtra por is_public (la columna se devuelve pero no gatea); aun así se
// siembra is_public=true para reflejar material publicado/visible.
//
// Tejido representante↔hijo↔inscripción en `base` que este fixture explota
// (ver seedGuardianTejido + seedSubjectOfferings de base):
//   - Miguel Castro (tutor.castro@edugo.test) ↔ hijos Carlos y Diego en S1
//     (Colegio San Ignacio, b1…01).
//   - Carlos (membership bb…01, unidad ac…03 = 5to A): inscrito en Matemáticas
//     (dd…01) en 5to A → verá "Tarea: Geometría básica".
//   - Diego  (membership bb…04, unidad ac…04 = 5to B): inscrito en Matemáticas
//     (dd…01) y Ciencias (dd…02) en 5to B → verá "Guía de Matemáticas:
//     Fracciones" y "Lectura: El Sistema Solar".
//
// NO se siembra content.material_file: file_count queda en 0, que es suficiente
// para validar la LISTA (no se ejercita descarga). uploaded_by_membership_id es
// la membresía del docente de la materia en `base` (FK RESTRICT a
// academic.memberships).
//
// Idempotente: OnConflict DoNothing por id. Rango UUID propio
// e1000000-0000-0000-0000-0000000000NN (no colisiona con los demás fixtures).
package mp10_material

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/base"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// Escuela y materias de `base` (Colegio San Ignacio).
	schoolID      = "b1000000-0000-0000-0000-000000000001"
	subjectMathID = "dd000000-0000-0000-0000-000000000001" // Matematicas
	subjectSciID  = "dd000000-0000-0000-0000-000000000002" // Ciencias Naturales

	// Unidades de `base`: 5to A (Carlos) y 5to B (Diego).
	unit5AID = "ac000000-0000-0000-0000-000000000003"
	unit5BID = "ac000000-0000-0000-0000-000000000004"

	// Docentes (membresías de `base`) que suben el material.
	teacher5AID = "bb000000-0000-0000-0000-000000000008" // docente 5to A
	teacher5BID = "bb000000-0000-0000-0000-000000000010" // docente 5to B

	// Material publicado = status 'ready' (único valor que el lector lista; ver
	// doc del paquete). NO existe 'published' en content.materials.
	statusReady = "ready"
)

// strPtr es un helper para los campos *string del entity (description).
func strPtr(s string) *string { return &s }

// Apply siembra el material publicado de los hijos del representante sobre el
// mundo `base`. Primero compone `base` (idempotente) para garantizar que existan
// las materias/unidades/membresías/inscripciones a las que apuntan las FKs del
// material; luego inserta las filas en content.materials. Idempotente por id.
func Apply(tx *gorm.DB) error {
	// Compone el mundo de datos `base` (trunca y resiembra; idempotente). El
	// runner aplica `system` + este fixture, pero NO encadena `base` solo.
	if err := base.Apply(tx); err != nil {
		return fmt.Errorf("playground_v2/mp10_material: base: %w", err)
	}

	schoolUUID := common.MustParseUUID(schoolID)
	subjectMathUUID := common.MustParseUUID(subjectMathID)
	subjectSciUUID := common.MustParseUUID(subjectSciID)
	unit5AUUID := common.MustParseUUID(unit5AID)
	unit5BUUID := common.MustParseUUID(unit5BID)
	teacher5AUUID := common.MustParseUUID(teacher5AID)
	teacher5BUUID := common.MustParseUUID(teacher5BID)

	materials := []entities.Material{
		// Matemáticas (5to B) — lo verá Diego (inscrito en Matemáticas en 5to B).
		{
			ID:                     common.MustParseUUID("e1000000-0000-0000-0000-000000000001"),
			SchoolID:               schoolUUID,
			UploadedByMembershipID: teacher5BUUID,
			SubjectID:              &subjectMathUUID,
			AcademicUnitID:         &unit5BUUID,
			Title:                  "Guía de Matemáticas: Fracciones",
			Description:            strPtr("Guía de práctica sobre suma y resta de fracciones."),
			Status:                 statusReady,
			IsPublic:               true,
		},
		// Ciencias Naturales (5to B) — lo verá Diego (inscrito en Ciencias en 5to B).
		{
			ID:                     common.MustParseUUID("e1000000-0000-0000-0000-000000000002"),
			SchoolID:               schoolUUID,
			UploadedByMembershipID: teacher5BUUID,
			SubjectID:              &subjectSciUUID,
			AcademicUnitID:         &unit5BUUID,
			Title:                  "Lectura: El Sistema Solar",
			Description:            strPtr("Lectura introductoria sobre los planetas del sistema solar."),
			Status:                 statusReady,
			IsPublic:               true,
		},
		// Matemáticas (5to A) — lo verá Carlos (inscrito en Matemáticas en 5to A).
		{
			ID:                     common.MustParseUUID("e1000000-0000-0000-0000-000000000003"),
			SchoolID:               schoolUUID,
			UploadedByMembershipID: teacher5AUUID,
			SubjectID:              &subjectMathUUID,
			AcademicUnitID:         &unit5AUUID,
			Title:                  "Tarea: Geometría básica",
			Description:            strPtr("Ejercicios de geometría básica: perímetro y área."),
			Status:                 statusReady,
			IsPublic:               true,
		},
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&materials).Error; err != nil {
		return fmt.Errorf("playground_v2/mp10_material: materials: %w", err)
	}

	return nil
}
