package demo

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// attendanceWindowDays es cuántos días atrás llega la asistencia: las tres
// últimas semanas, solo días hábiles, sin hoy.
const attendanceWindowDays = 21

// accountGrades son las notas de las dos cuentas con login, en el orden de las
// sesiones (Matemáticas, Lenguaje, Ciencias). Fijas para que el ensayo sepa qué
// ver; el relleno sale de una fórmula.
var accountGrades = map[uuid.UUID][3]float64{
	Student1MembershipID: {6.2, 5.5, 6.8},
	Student2MembershipID: {6.7, 6.1, 5.9},
}

// buildRecords arma las notas, la asistencia y los anuncios.
func buildRecords(w *world) {
	buildGrades(w)
	buildAttendance(w)
	buildAnnouncements(w)
}

// buildGrades: una nota por alumno y materia en el período activo, escala 1 a 7.
// Matemáticas va finalizada; Lenguaje y Ciencias, no.
func buildGrades(w *world) {
	teacher := TeacherMembershipID
	created := w.day(-10).Add(11 * time.Hour)
	finalized := w.day(-7).Add(18 * time.Hour)
	n := 0
	for i, s := range students() {
		fixed, isAccount := accountGrades[s.MembershipID]
		for j, o := range w.offerings {
			n++
			value := fixed[j]
			if !isAccount {
				value = math.Round((3.5+float64((i*7+j*11)%36)/10)*10) / 10
			}
			g := entities.Grade{
				ID: newID(domGrade, n), MembershipID: s.MembershipID, SubjectID: o.SubjectID, PeriodID: PeriodID,
				GradeValue: ptr(value), AssessmentScores: json.RawMessage(`[]`), Source: "manual",
				TeacherID: &teacher, CreatedAt: created, UpdatedAt: created,
			}
			if o.ID == OfferingMathID {
				g.FinalizedAt = ptr(finalized)
				g.UpdatedAt = finalized
			}
			w.grades = append(w.grades, g)
		}
	}
}

// weekdays devuelve los días hábiles de la ventana de asistencia, del más viejo
// al más reciente.
func (w *world) weekdays() []time.Time {
	var out []time.Time
	for d := -attendanceWindowDays; d <= -1; d++ {
		day := w.day(d)
		if wd := day.Weekday(); wd != time.Saturday && wd != time.Sunday {
			out = append(out, day)
		}
	}
	return out
}

// attendanceStatus decide el estado del alumno i en la materia j el día k de la
// ventana. Las ausencias y atrasos de las cuentas van por índice sobre los días
// generados, así que existen cualquiera sea el día en que se aplica.
func attendanceStatus(i, j, k, days int) string {
	switch {
	case i == 0 && j == 0 && k == 2, i == 1 && j == 1 && k == 4:
		return "absent"
	case i == 0 && j == 1 && k == days-2, i == 1 && j == 0 && k == days-3:
		return "late"
	case i < 2:
		return "present"
	case (i*5+k*3+j)%23 == 0:
		return "absent"
	case (i*5+k*3+j)%17 == 1:
		return "late"
	default:
		return "present"
	}
}

// buildAttendance: cada alumno, cada día hábil de las tres últimas semanas, en
// cada materia (subject_id relleno: si queda nulo el índice único no evita
// duplicados, deuda 098). La pasa la docente.
func buildAttendance(w *world) {
	recorder := uuid.MustParse(TeacherUserID)
	days := w.weekdays()
	n := 0
	for i, s := range students() {
		for j, o := range w.offerings {
			subject := o.SubjectID
			for k, day := range days {
				n++
				w.attendance = append(w.attendance, entities.Attendance{
					ID: newID(domAttendance, n), MembershipID: s.MembershipID, SubjectID: &subject,
					Date: day, Status: attendanceStatus(i, j, k, len(days)), RecordedBy: recorder,
					CreatedAt: day.Add(time.Duration(9+j) * time.Hour),
				})
			}
		}
	}
}

// buildAnnouncements: tres anuncios de la docente en las dos últimas semanas
// (propuesta 9: la demo no tiene admin). Uno de escuela y fijado, uno de su
// unidad (el que puede editar en vivo) y uno de escuela reciente.
func buildAnnouncements(w *world) {
	author := uuid.MustParse(TeacherUserID)
	unit := UnitID
	targets := pq.StringArray{"student", "guardian"}
	mk := func(n int, unitID *uuid.UUID, scope, title, body string, pinned bool, at time.Time) entities.Announcement {
		return entities.Announcement{
			ID: newID(domAnnouncement, n), SchoolID: SchoolID, AcademicUnitID: unitID, AuthorID: author,
			Title: title, Body: body, Scope: scope, TargetRoles: targets, IsPinned: pinned,
			PublishedAt: ptr(at), CreatedAt: at, UpdatedAt: at,
		}
	}
	w.announcements = []entities.Announcement{
		mk(1, nil, "school", "Bienvenidos al semestre",
			"Les damos la bienvenida al semestre. Revisen el calendario de evaluaciones en la aplicación.",
			true, w.day(-13).Add(8*time.Hour)),
		mk(2, &unit, "unit", "Prueba de fracciones en 8.º Básico A",
			"La prueba de fracciones ya está disponible en la aplicación. Tienen dos intentos.",
			false, w.day(-6).Add(12*time.Hour)),
		mk(3, nil, "school", "Reunión de apoderados",
			"La reunión de apoderados será el próximo jueves a las 18:30 en el salón principal.",
			false, w.day(-2).Add(16*time.Hour)),
	}
}

func writeRecords(tx *gorm.DB, w *world) error {
	if err := insert(tx, &w.grades); err != nil {
		return fmt.Errorf("grades: %w", err)
	}
	if err := insert(tx, &w.attendance); err != nil {
		return fmt.Errorf("attendance: %w", err)
	}
	if err := insert(tx, &w.announcements); err != nil {
		return fmt.Errorf("announcements: %w", err)
	}
	return nil
}
