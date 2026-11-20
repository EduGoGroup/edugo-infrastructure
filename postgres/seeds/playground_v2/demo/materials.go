package demo

import (
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"gorm.io/gorm"
)

// MaterialFileName es el nombre del PDF de T1 que la AR-DEMO.2 sube a
// MaterialFileKey(). El fixture no toca el bucket.
const MaterialFileName = "fracciones-equivalentes.pdf"

// materialFileSizeBytes es el tamaño declarado del PDF. Solo se muestra; si el
// PDF que se sube pesa otra cosa, la descarga no cambia.
const materialFileSizeBytes = 245760

// buildMaterials arma T1: un tema de Matemáticas en estado `ready` (la tutora
// solo ve los materiales en ese estado), público, con un archivo F1 cuya clave es
// la del bucket, y asignado a la sesión de Matemáticas sin ventana. Sin trabajos
// de pipeline: con uno pendiente, «Generar evaluación» queda muerto.
func buildMaterials(w *world) {
	created := w.day(-4).Add(10 * time.Hour)
	unit, subject := UnitID, SubjectMathID
	w.materials = []entities.Material{{
		ID: MaterialID, SchoolID: SchoolID, UploadedByMembershipID: TeacherMembershipID,
		SubjectID: &subject, AcademicUnitID: &unit, Title: "Fracciones equivalentes",
		Description: ptr("Guía con ejemplos y ejercicios de fracciones equivalentes."),
		Summary: ptr("## Fracciones equivalentes\n\n" +
			"- Dos fracciones son equivalentes si representan la misma parte del entero.\n" +
			"- Se obtienen multiplicando o dividiendo el numerador y el denominador por el mismo número.\n" +
			"- Ejemplo: 1/2 = 2/4 = 3/6."),
		Status: "ready", IsPublic: true, CreatedAt: created, UpdatedAt: created,
	}}
	w.materialFiles = []entities.MaterialFile{{
		ID: MaterialFileID, MaterialID: MaterialID, FileURL: MaterialFileKey(), FileName: MaterialFileName,
		FileType: "application/pdf", FileSizeBytes: materialFileSizeBytes, CreatedAt: created,
	}}
	assigned := created.Add(time.Hour)
	w.materialAsgs = []entities.MaterialAssignment{{
		ID: newID(domMaterialAssignment, 1), MaterialID: MaterialID, SubjectOfferingID: OfferingMathID,
		AssignedByMembershipID: TeacherMembershipID, AssignedAt: assigned, CreatedAt: assigned, UpdatedAt: assigned,
	}}
}

func writeMaterials(tx *gorm.DB, w *world) error {
	if err := insert(tx, &w.materials); err != nil {
		return fmt.Errorf("materials: %w", err)
	}
	if err := insert(tx, &w.materialFiles); err != nil {
		return fmt.Errorf("material_file: %w", err)
	}
	if err := insert(tx, &w.materialAsgs); err != nil {
		return fmt.Errorf("material_assignment: %w", err)
	}
	return nil
}
