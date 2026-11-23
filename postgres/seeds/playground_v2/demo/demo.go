// Package demo es el mundo de datos de la demo en KMP (MVP-DEMO, HU-DEMO.5): una
// escuela, «Colegio Los Aromos», con su unidad «8.º Básico A», una docente, dos
// alumnos y una tutora con login, y dieciocho alumnos de relleno sin login. Cada
// cuenta encuentra lo que su escena del libreto necesita: tres materias, cuatro
// evaluaciones (E-B en borrador, E-X por rendir, E-P de práctica y E-R con
// intentos por corregir), un tema con su PDF (T1), notas, asistencia, anuncios y
// sus avisos.
//
// Es un fixture autónomo pero NO entra en "all": va en un registro aparte de
// playground_v2.go, porque solo tiene sentido en la base de demo, que el migrador
// recrea con --playground-v2=demo (fuerza siempre la recreación).
//
// Reglas que respeta (libreto, hoja «El mundo de demo»):
//   - Los ids de las cuatro cuentas son el `sub` de identity-core (ids.go); los
//     demás salen de un rango propio, de0000DD-…, partido por dominio.
//   - Todo es relativo al momento de aplicar, truncado al día en UTC: el período,
//     la asistencia, los anuncios, los intentos entregados y los avisos.
//   - Solo corre sobre una base sin otra escuela: no trunca nada (ver Apply).
//   - No encola nada en auth.identity_sync_outbox: la sincronización con
//     identity-core está apagada en el ambiente de demo (ADR 0046).
package demo

import (
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// world es el mundo entero como datos en Go, armado antes de tocar la base a
// partir de un único `now`. Las pruebas sin base lo inspeccionan tal cual.
type world struct {
	now   time.Time // momento de aplicar (UTC)
	today time.Time // now truncado al día

	school        common.SchoolSpec
	units         []entities.AcademicUnit
	period        common.PeriodSpec
	subjects      []common.SubjectSpec
	users         []common.UserSpec
	memberships   []membershipRow
	userRoles     []entities.UserRole
	offerings     []common.OfferingSpec
	enrollments   []entities.SubjectOfferingEnrollment
	guardian      common.GuardianRelationSpec
	grants        []entities.UserGrant
	assessments   []entities.Assessment
	questions     []entities.Question
	options       []entities.QuestionOption
	assignments   []entities.AssessmentAssignment
	attempts      []entities.AssessmentAttempt
	answers       []entities.AssessmentAttemptAnswer
	materials     []entities.Material
	materialFiles []entities.MaterialFile
	materialAsgs  []entities.MaterialAssignment
	notifications []entities.Notification
	grades        []entities.Grade
	attendance    []entities.Attendance
	announcements []entities.Announcement
}

// day devuelve el día de hoy desplazado n días (n negativo = atrás).
func (w *world) day(n int) time.Time { return w.today.AddDate(0, 0, n) }

// buildWorld arma el mundo para el momento now. Es puro: no toca la base.
func buildWorld(now time.Time) *world {
	now = now.UTC()
	w := &world{now: now, today: now.Truncate(24 * time.Hour)}
	buildSchool(w)
	buildPeople(w)
	buildGrants(w)
	buildAssessments(w)
	buildAttempts(w)
	buildMaterials(w)
	buildNotifications(w)
	buildRecords(w)
	return w
}

// Apply siembra el mundo de demo en una transacción. Se niega si la base ya
// tiene una escuela que no es la suya (un `base` debajo, o cualquier otro mundo):
// el fixture no trunca nada y se aplica siempre sobre una base recién recreada.
// Sobre su propio mundo es re-ejecutable (OnConflict DoNothing por id), pero no
// refresca las fechas: para eso se recrea la base.
func Apply(gdb *gorm.DB) error {
	w := buildWorld(time.Now())
	return gdb.Transaction(func(tx *gorm.DB) error {
		if err := ensureEmptyBase(tx); err != nil {
			return err
		}
		if err := write(tx, w); err != nil {
			return fmt.Errorf("playground_v2/demo: %w", err)
		}
		return nil
	})
}

// ensureEmptyBase falla si hay alguna escuela distinta de la de demo.
func ensureEmptyBase(tx *gorm.DB) error {
	var others int64
	if err := tx.Table("academic.schools").Where("id <> ?", SchoolID).Count(&others).Error; err != nil {
		return fmt.Errorf("playground_v2/demo: contando escuelas: %w", err)
	}
	if others > 0 {
		return fmt.Errorf("playground_v2/demo: la base ya tiene %d escuela(s) que no son la de demo; "+
			"demo solo se aplica sobre una base recién recreada (--playground-v2=demo)", others)
	}
	return nil
}

// write escribe el mundo en el orden de sus claves foráneas.
func write(tx *gorm.DB, w *world) error {
	steps := []struct {
		name string
		fn   func(*gorm.DB, *world) error
	}{
		{"escuela", writeSchool},
		{"personas", writePeople},
		{"recorte de permisos", func(tx *gorm.DB, w *world) error { return insert(tx, &w.grants) }},
		{"evaluaciones", writeAssessments},
		{"intentos", writeAttempts},
		{"materiales", writeMaterials},
		{"avisos", func(tx *gorm.DB, w *world) error { return insert(tx, &w.notifications) }},
		{"notas, asistencia y anuncios", writeRecords},
	}
	for _, s := range steps {
		if err := s.fn(tx, w); err != nil {
			return fmt.Errorf("%s: %w", s.name, err)
		}
	}
	return nil
}

// insert crea las filas ignorando las que ya existen por su clave primaria.
func insert[T any](tx *gorm.DB, rows *[]T) error {
	if len(*rows) == 0 {
		return nil
	}
	return tx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(rows, 200).Error
}

func ptr[T any](v T) *T { return &v }
