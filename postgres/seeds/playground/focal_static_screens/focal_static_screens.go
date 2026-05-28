// Package focal_static_screens extiende el playground focal_botonera con
// grants sobre los recursos cuyas pantallas estaticas fueron migradas en
// la fase 3 del estandar ui-adaptation (Materials + Notifications).
//
// **Motivacion**: los 3 usuarios focales (botonera-viewer / -author /
// -publisher) solo tenian grants sobre `content.assessments.*` y
// `academic.announcements.*`. Las nuevas pantallas migradas a DSListRow
// (MaterialsListScreen, MaterialDetailNativeScreen, MaterialUploadScreen,
// NotificationsListScreen) gatean por `content.materials.*` y
// `notifications.*`, asi que los focales no podian entrar a validar la
// botonera ahi. Este playground agrega los grants faltantes manteniendo
// el patron wildcard-first ya establecido.
//
// **Composicion (autosuficiente)**: depende de [focal_botonera] (que a su
// vez depende de focal_evaluacion_v2). Apply() encadena
// focal_botonera.Apply() al inicio para garantizar la cadena cuando se
// invoca con `P=focal-static-screens` standalone. Idempotente: los
// INSERTs usan OnConflict DoNothing.
//
// **Roles extendidos** (NO se crean roles nuevos, se agregan grants a
// los roles existentes 13000000-...-001/002/003 de focal_botonera):
//   - focal-viewer    -> +content.materials.read, +notifications.read
//   - focal-author    -> +content.materials.{read,create,update},
//     +notifications.{read,create,update}
//   - focal-publisher -> +content.materials.*, +notifications.*
//
// **Por que no editamos focal_botonera.go**: la convencion del proyecto
// trata cada playground como una foto inmutable; los grants extra viven
// en este paquete aparte (project_edugo_playgrounds_convention).
package focal_static_screens

import (
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground/common"
	focal_botonera "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground/focal_botonera"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// Role IDs heredados de focal_botonera (rango 13000000-...).
	viewerRoleID    = "13000000-0000-0000-0000-000000000001"
	authorRoleID    = "13000000-0000-0000-0000-000000000002"
	publisherRoleID = "13000000-0000-0000-0000-000000000003"

	// Tenant heredado de focal_botonera -> focal_evaluacion_v2.
	tenantSchoolID = "62000000-0000-0000-0000-000000000001"
	tenantUnitID   = "62000000-0000-0000-0000-000000000002"
	// Autor (publisher) de los materiales de prueba.
	publisherUserID = "63000000-0000-0000-0000-000000000003"
)

// Apply siembra los grants extra del playground focal_static_screens.
// Encadena focal_botonera.Apply al inicio para ser autosuficiente cuando
// se invoca con `P=focal-static-screens` standalone. Con `P=all` la
// dependencia ya corrio antes en el registry — la doble aplicacion es
// no-op gracias a OnConflict DoNothing. Idempotente.
func Apply(tx *gorm.DB) error {
	if err := focal_botonera.Apply(tx); err != nil {
		return fmt.Errorf("playground/focal_static_screens: dependencia focal_botonera: %w", err)
	}
	if err := upsertRoleGrants(tx); err != nil {
		return fmt.Errorf("playground/focal_static_screens: role_grants: %w", err)
	}
	if err := upsertMaterials(tx); err != nil {
		return fmt.Errorf("playground/focal_static_screens: materials: %w", err)
	}
	return nil
}

// upsertMaterials siembra 6 materiales en estado `ready` para validar la
// grilla migrada de MaterialsListScreen sin pasar por el CRUD (que aun
// no esta probado). Variedad de FileType / Subject / Grade para que el
// row del DSListRow muestre supporting text rico. IDs en rango
// 63000030-... para no colisionar con announcements (63000020-...) ni
// users/roles del playground (63000000-...).
func upsertMaterials(tx *gorm.DB) error {
	sid := common.MustParseUUID(tenantSchoolID)
	uid := common.MustParseUUID(tenantUnitID)
	teacherID := common.MustParseUUID(publisherUserID)
	now := time.Now().UTC()

	str := func(s string) *string { return &s }

	items := []entities.Material{
		{
			ID:                    common.MustParseUUID("63000030-0000-0000-0000-000000000001"),
			SchoolID:              sid,
			UploadedByTeacherID:   teacherID,
			AcademicUnitID:        &uid,
			Title:                 "Guia de algebra basica",
			Description:           str("PDF con ejercicios resueltos de algebra elemental para 7mo grado."),
			Subject:               str("Matematicas"),
			Grade:                 str("7"),
			FileURL:               "https://files.edugo.local/playground/algebra-basica.pdf",
			FileType:              "application/pdf",
			FileSizeBytes:         1_240_000,
			Status:                "ready",
			ProcessingCompletedAt: &now,
			IsPublic:              false,
		},
		{
			ID:                    common.MustParseUUID("63000030-0000-0000-0000-000000000002"),
			SchoolID:              sid,
			UploadedByTeacherID:   teacherID,
			AcademicUnitID:        &uid,
			Title:                 "Video clase: sistema solar",
			Description:           str("Explicacion en video de 8 minutos sobre los planetas del sistema solar."),
			Subject:               str("Ciencias"),
			Grade:                 str("5"),
			FileURL:               "https://files.edugo.local/playground/sistema-solar.mp4",
			FileType:              "video/mp4",
			FileSizeBytes:         48_500_000,
			Status:                "ready",
			ProcessingCompletedAt: &now,
			IsPublic:              true,
		},
		{
			ID:                    common.MustParseUUID("63000030-0000-0000-0000-000000000003"),
			SchoolID:              sid,
			UploadedByTeacherID:   teacherID,
			AcademicUnitID:        &uid,
			Title:                 "Presentacion: revolucion industrial",
			Description:           str("Slides en formato pptx con linea de tiempo y figuras clave."),
			Subject:               str("Historia"),
			Grade:                 str("9"),
			FileURL:               "https://files.edugo.local/playground/revolucion-industrial.pptx",
			FileType:              "application/vnd.openxmlformats-officedocument.presentationml.presentation",
			FileSizeBytes:         3_700_000,
			Status:                "ready",
			ProcessingCompletedAt: &now,
			IsPublic:              false,
		},
		{
			ID:                    common.MustParseUUID("63000030-0000-0000-0000-000000000004"),
			SchoolID:              sid,
			UploadedByTeacherID:   teacherID,
			AcademicUnitID:        &uid,
			Title:                 "Hoja de practica: comprension lectora",
			Description:           str("Texto + 10 preguntas de seleccion multiple."),
			Subject:               str("Lengua"),
			Grade:                 str("6"),
			FileURL:               "https://files.edugo.local/playground/comprension-lectora.pdf",
			FileType:              "application/pdf",
			FileSizeBytes:         620_000,
			Status:                "ready",
			ProcessingCompletedAt: &now,
			IsPublic:              false,
		},
		{
			ID:                    common.MustParseUUID("63000030-0000-0000-0000-000000000005"),
			SchoolID:              sid,
			UploadedByTeacherID:   teacherID,
			AcademicUnitID:        &uid,
			Title:                 "Mapa: regiones de Colombia",
			Description:           str("Imagen JPG en alta resolucion para usar en clase de geografia."),
			Subject:               str("Geografia"),
			Grade:                 str("8"),
			FileURL:               "https://files.edugo.local/playground/mapa-regiones.jpg",
			FileType:              "image/jpeg",
			FileSizeBytes:         2_100_000,
			Status:                "ready",
			ProcessingCompletedAt: &now,
			IsPublic:              true,
		},
		{
			ID:                    common.MustParseUUID("63000030-0000-0000-0000-000000000006"),
			SchoolID:              sid,
			UploadedByTeacherID:   teacherID,
			AcademicUnitID:        &uid,
			Title:                 "Audio: dictado en ingles",
			Description:           str("Pista MP3 de 3 minutos para ejercicio de listening."),
			Subject:               str("Ingles"),
			Grade:                 str("10"),
			FileURL:               "https://files.edugo.local/playground/dictado-ingles.mp3",
			FileType:              "audio/mpeg",
			FileSizeBytes:         2_900_000,
			Status:                "ready",
			ProcessingCompletedAt: &now,
			IsPublic:              false,
		},
	}

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&items).Error
}

// upsertRoleGrants agrega los patterns de materials + notifications a
// los 3 roles existentes. Mismo shape wildcard-first que ya usan los
// recursos focales (content.assessments / academic.announcements).
func upsertRoleGrants(tx *gorm.DB) error {
	type grantSpec struct {
		roleID   string
		patterns []string
	}
	specs := []grantSpec{
		{viewerRoleID, []string{
			"content.materials.read",
			"notifications.read",
		}},
		{authorRoleID, []string{
			"content.materials.read",
			"content.materials.create",
			"content.materials.update",
			"notifications.read",
			"notifications.create",
			"notifications.update",
		}},
		{publisherRoleID, []string{
			"content.materials.*",
			"notifications.*",
		}},
	}
	for _, s := range specs {
		rid := common.MustParseUUID(s.roleID)
		for _, pattern := range s.patterns {
			if err := common.SeedRoleGrant(tx, rid, pattern); err != nil {
				return err
			}
		}
	}
	return nil
}
