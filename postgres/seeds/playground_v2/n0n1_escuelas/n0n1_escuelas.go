// Package n0n1_escuelas es un playground de la línea v2 que siembra un
// ecosistema multi-escuela realista para validar N0 (onboarding / join
// requests) y N1 (estructura académica: unidades, materias, sesiones,
// membresías) sobre el sistema completo (L0..L4).
//
// Tres escuelas, una de cada concept_type:
//   - Cristo Rey (primary_school): 3 grados, 4 materias, secciones A/B → 24 sesiones.
//   - UCV (university, concept_type NUEVO sembrado aquí): 1 semestre, 3 materias,
//     secciones Mañana/Tarde → 6 sesiones.
//   - InglesLab (workshop): 1 nivel, 1 materia, secciones Mañana/Tarde → 2 sesiones.
//
// Total 32 subject_offerings, todas SIN docente asignado (teacher_membership_id
// NULL) — la asignación docente↔sesión queda para ejercitarse en la app.
//
// Usuarios (emails @edugo.local en minúscula sin acentos). Escenario on-top
// (plan 053 T5.2): NO crea filas en auth.users; los crea antes la API de
// identidad (Users()) y aquí se resuelven por email:
//   - 1 admin global (admin@edugo.local) con rol L0 super_admin (acceso "*") +
//     membresía "admin" con alcance COLEGIO en las 3 escuelas.
//   - 13 profesores con membresías "teacher" de alcance UNIDAD (17 membresías:
//     varios profesores enseñan en 2 unidades). Ningún profesor se vincula a una
//     subject_offering (eso queda para la app).
//   - 1 solicitante N0 (carlos.estudiante@edugo.local) con una school_join_request
//     PENDIENTE a InglesLab / Nivel Básico, sin firmas, para probar el doble gate.
//
// Como todo v2, asume que L0..L4 corrieron: reusa el rol L0 super_admin
// (10000000-...-001) y los concept_types L4 primary_school (c4..001) y workshop
// (c4..005); sólo crea el concept_type university (69..00c1) que faltaba.
//
// Rango UUID propio: 69000000-0000-0000-0000-XXXXXXXXXXXX. Idempotente:
// OnConflict DoNothing por PK (o clave natural) en toda inserción.
package n0n1_escuelas

import (
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/catalog"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// Emails de los usuarios que no son profesores (los profesores van en
	// teacherSeeds).
	AdminEmail       = "admin@edugo.local"
	SolicitanteEmail = "carlos.estudiante@edugo.local"

	// Rol L0 super_admin (reusado, NO se crea) para el admin global.
	superAdminRoleID = "10000000-0000-0000-0000-000000000001"

	// concept_types: primary_school y workshop ya existen en L4; university se
	// crea aquí.
	conceptTypePrimaryID    = "c4000000-0000-0000-0000-000000000001"
	conceptTypeWorkshopID   = "c4000000-0000-0000-0000-000000000005"
	conceptTypeUniversityID = "69000000-0000-0000-0000-0000000000c1"

	// Escuelas.
	schoolCristoReyID = "69000000-0000-0000-0000-000000000101"
	schoolUCVID       = "69000000-0000-0000-0000-000000000102"
	schoolInglesLabID = "69000000-0000-0000-0000-000000000103"

	// Períodos académicos (1 por escuela).
	periodCristoReyID = "69000000-0000-0000-0000-000000000111"
	periodUCVID       = "69000000-0000-0000-0000-000000000112"
	periodInglesLabID = "69000000-0000-0000-0000-000000000113"

	// Unidades académicas.
	unitCR1ID    = "69000000-0000-0000-0000-000000000121" // Cristo Rey - Primer Grado
	unitCR2ID    = "69000000-0000-0000-0000-000000000122" // Cristo Rey - Segundo Grado
	unitCR3ID    = "69000000-0000-0000-0000-000000000123" // Cristo Rey - Tercer Grado
	unitUCVS1ID  = "69000000-0000-0000-0000-000000000124" // UCV - Primer Semestre
	unitINGLNBID = "69000000-0000-0000-0000-000000000125" // InglesLab - Nivel Básico

	// Materias.
	subjMatID  = "69000000-0000-0000-0000-000000000131" // CR Matemática
	subjCasID  = "69000000-0000-0000-0000-000000000132" // CR Castellano
	subjCsoID  = "69000000-0000-0000-0000-000000000133" // CR Ciencia Sociales
	subjRelID  = "69000000-0000-0000-0000-000000000134" // CR Religión
	subjMatIID = "69000000-0000-0000-0000-000000000135" // UCV Matemática I
	subjCasIID = "69000000-0000-0000-0000-000000000136" // UCV Castellano I
	subjCsoUID = "69000000-0000-0000-0000-000000000137" // UCV Ciencia Sociales
	subjGraID  = "69000000-0000-0000-0000-000000000138" // InglesLab Gramática

	// Membresías del admin (alcance colegio en cada escuela).
	adminMembCristoReyID = "69000000-0000-0000-0000-000000001311"
	adminMembUCVID       = "69000000-0000-0000-0000-000000001312"
	adminMembInglesLabID = "69000000-0000-0000-0000-000000001313"

	// School join request pendiente (N0).
	joinRequestID = "69000000-0000-0000-0000-000000001401"

	academicYear = 2026
)

// concept_type a sembrar (los otros 2 ya existen en L4).
var universityConceptType = struct {
	id, code, name, desc string
}{
	id:   conceptTypeUniversityID,
	code: "university",
	name: "Universidad",
	desc: "Institución de educación superior",
}

// schoolSeed describe una escuela a sembrar.
type schoolSeed struct {
	id, name, code, conceptTypeID string
}

var schoolSeeds = []schoolSeed{
	{id: schoolCristoReyID, name: "Cristo Rey", code: "CRISTO-REY", conceptTypeID: conceptTypePrimaryID},
	{id: schoolUCVID, name: "UCV", code: "UCV", conceptTypeID: conceptTypeUniversityID},
	{id: schoolInglesLabID, name: "InglesLab", code: "INGLESLAB", conceptTypeID: conceptTypeWorkshopID},
}

// periodSeed describe un período académico (1 por escuela).
type periodSeed struct {
	id, schoolID, name, code string
	startY, startM, startD   int
	endY, endM, endD         int
}

var periodSeeds = []periodSeed{
	{id: periodCristoReyID, schoolID: schoolCristoReyID, name: "Marzo 2026 - Diciembre 2026", code: "CR-2026", startY: 2026, startM: 3, startD: 1, endY: 2026, endM: 12, endD: 15},
	{id: periodUCVID, schoolID: schoolUCVID, name: "Verano 2026", code: "UCV-V2026", startY: 2026, startM: 7, startD: 1, endY: 2026, endM: 9, endD: 15},
	{id: periodInglesLabID, schoolID: schoolInglesLabID, name: "Verano 2026", code: "IL-V2026", startY: 2026, startM: 7, startD: 1, endY: 2026, endM: 9, endD: 15},
}

// unitSeed describe una unidad académica.
type unitSeed struct {
	id, schoolID, name, code string
}

var unitSeeds = []unitSeed{
	{id: unitCR1ID, schoolID: schoolCristoReyID, name: "Primer Grado", code: "CR-1"},
	{id: unitCR2ID, schoolID: schoolCristoReyID, name: "Segundo Grado", code: "CR-2"},
	{id: unitCR3ID, schoolID: schoolCristoReyID, name: "Tercer Grado", code: "CR-3"},
	{id: unitUCVS1ID, schoolID: schoolUCVID, name: "Primer Semestre", code: "UCV-S1"},
	{id: unitINGLNBID, schoolID: schoolInglesLabID, name: "Nivel Básico", code: "IL-NB"},
}

// subjectSeed describe una materia (academic_unit_id = NULL).
type subjectSeed struct {
	id, schoolID, name, code string
}

var subjectSeeds = []subjectSeed{
	{id: subjMatID, schoolID: schoolCristoReyID, name: "Matemática", code: "MAT"},
	{id: subjCasID, schoolID: schoolCristoReyID, name: "Castellano", code: "CAS"},
	{id: subjCsoID, schoolID: schoolCristoReyID, name: "Ciencia Sociales", code: "CSO"},
	{id: subjRelID, schoolID: schoolCristoReyID, name: "Religión", code: "REL"},
	{id: subjMatIID, schoolID: schoolUCVID, name: "Matemática I", code: "MATI"},
	{id: subjCasIID, schoolID: schoolUCVID, name: "Castellano I", code: "CASI"},
	{id: subjCsoUID, schoolID: schoolUCVID, name: "Ciencia Sociales", code: "CSO-U"},
	{id: subjGraID, schoolID: schoolInglesLabID, name: "Gramática", code: "GRA"},
}

// offeringGroup describe un conjunto de subject_offerings generadas por loop:
// para una escuela/unidad/período, el producto cartesiano subjects × sections.
type offeringGroup struct {
	schoolID, unitID, periodID string
	subjectIDs                 []string
	sections                   []string
}

var offeringGroups = []offeringGroup{
	// Cristo Rey: Primer/Segundo/Tercer Grado × 4 materias × secciones A/B = 24.
	{schoolID: schoolCristoReyID, unitID: unitCR1ID, periodID: periodCristoReyID, subjectIDs: []string{subjMatID, subjCasID, subjCsoID, subjRelID}, sections: []string{"A", "B"}},
	{schoolID: schoolCristoReyID, unitID: unitCR2ID, periodID: periodCristoReyID, subjectIDs: []string{subjMatID, subjCasID, subjCsoID, subjRelID}, sections: []string{"A", "B"}},
	{schoolID: schoolCristoReyID, unitID: unitCR3ID, periodID: periodCristoReyID, subjectIDs: []string{subjMatID, subjCasID, subjCsoID, subjRelID}, sections: []string{"A", "B"}},
	// UCV: Primer Semestre × 3 materias × secciones Mañana/Tarde = 6.
	{schoolID: schoolUCVID, unitID: unitUCVS1ID, periodID: periodUCVID, subjectIDs: []string{subjMatIID, subjCasIID, subjCsoUID}, sections: []string{"Mañana", "Tarde"}},
	// InglesLab: Nivel Básico × 1 materia × secciones Mañana/Tarde = 2.
	{schoolID: schoolInglesLabID, unitID: unitINGLNBID, periodID: periodInglesLabID, subjectIDs: []string{subjGraID}, sections: []string{"Mañana", "Tarde"}},
}

// teacherSeed describe un profesor: nombre, apellido y email. El id lo pone la
// API de identidad al crearlo; aquí se resuelve por email.
type teacherSeed struct {
	first, last, email string
}

var teacherSeeds = []teacherSeed{
	{first: "Pedro", last: "Perez", email: "pedro.perez@edugo.local"},
	{first: "Juan", last: "Perez", email: "juan.perez@edugo.local"},
	{first: "Pilar", last: "Irragory", email: "pilar.irragory@edugo.local"},
	{first: "Juan", last: "Soto", email: "juan.soto@edugo.local"},
	{first: "Many", last: "Ramirez", email: "many.ramirez@edugo.local"},
	{first: "Francisco", last: "I", email: "francisco.i@edugo.local"},
	{first: "Wladimir", last: "Guerrero", email: "wladimir.guerrero@edugo.local"},
	{first: "Luis", last: "Sojo", email: "luis.sojo@edugo.local"},
	{first: "Juan", last: "23", email: "juan.23@edugo.local"},
	{first: "Felipe", last: "Lira", email: "felipe.lira@edugo.local"},
	{first: "Gustavo", last: "Dudamel", email: "gustavo.dudamel@edugo.local"},
	{first: "Luis", last: "Silva", email: "luis.silva@edugo.local"},
	{first: "Sheldon", last: "Cooper", email: "sheldon.cooper@edugo.local"},
}

// teacherMembershipSeed describe una membresía "teacher" (alcance UNIDAD).
// Cada profesor puede tener 1 o 2 membresías (enseña en 1 o 2 unidades).
type teacherMembershipSeed struct {
	id, teacherEmail, schoolID, unitID string
}

// teacherMembershipSeeds: 17 membresías teacher, ids 69..1321..69..1331.
var teacherMembershipSeeds = []teacherMembershipSeed{
	// Pedro Perez (1011) → CR Primer Grado + UCV Primer Semestre.
	{id: "69000000-0000-0000-0000-000000001321", teacherEmail: "pedro.perez@edugo.local", schoolID: schoolCristoReyID, unitID: unitCR1ID},
	{id: "69000000-0000-0000-0000-000000001322", teacherEmail: "pedro.perez@edugo.local", schoolID: schoolUCVID, unitID: unitUCVS1ID},
	// Juan Perez (1012) → CR Primer Grado + UCV Primer Semestre.
	{id: "69000000-0000-0000-0000-000000001323", teacherEmail: "juan.perez@edugo.local", schoolID: schoolCristoReyID, unitID: unitCR1ID},
	{id: "69000000-0000-0000-0000-000000001324", teacherEmail: "juan.perez@edugo.local", schoolID: schoolUCVID, unitID: unitUCVS1ID},
	// Pilar Irragory (1013) → CR Primer Grado.
	{id: "69000000-0000-0000-0000-000000001325", teacherEmail: "pilar.irragory@edugo.local", schoolID: schoolCristoReyID, unitID: unitCR1ID},
	// Juan Soto (1014) → CR Segundo Grado.
	{id: "69000000-0000-0000-0000-000000001326", teacherEmail: "juan.soto@edugo.local", schoolID: schoolCristoReyID, unitID: unitCR2ID},
	// Many Ramirez (1015) → CR Segundo Grado.
	{id: "69000000-0000-0000-0000-000000001327", teacherEmail: "many.ramirez@edugo.local", schoolID: schoolCristoReyID, unitID: unitCR2ID},
	// Francisco I (1016) → CR Segundo Grado.
	{id: "69000000-0000-0000-0000-000000001328", teacherEmail: "francisco.i@edugo.local", schoolID: schoolCristoReyID, unitID: unitCR2ID},
	// Wladimir Guerrero (1017) → CR Tercer Grado.
	{id: "69000000-0000-0000-0000-000000001329", teacherEmail: "wladimir.guerrero@edugo.local", schoolID: schoolCristoReyID, unitID: unitCR3ID},
	// Luis Sojo (1018) → CR Tercer Grado.
	{id: "69000000-0000-0000-0000-00000000132a", teacherEmail: "luis.sojo@edugo.local", schoolID: schoolCristoReyID, unitID: unitCR3ID},
	// Juan 23 (1019) → CR Tercer Grado.
	{id: "69000000-0000-0000-0000-00000000132b", teacherEmail: "juan.23@edugo.local", schoolID: schoolCristoReyID, unitID: unitCR3ID},
	// Felipe Lira (101a) → UCV Primer Semestre + InglesLab Nivel Básico.
	{id: "69000000-0000-0000-0000-00000000132c", teacherEmail: "felipe.lira@edugo.local", schoolID: schoolUCVID, unitID: unitUCVS1ID},
	{id: "69000000-0000-0000-0000-00000000132d", teacherEmail: "felipe.lira@edugo.local", schoolID: schoolInglesLabID, unitID: unitINGLNBID},
	// Gustavo Dudamel (101b) → UCV Primer Semestre + InglesLab Nivel Básico.
	{id: "69000000-0000-0000-0000-00000000132e", teacherEmail: "gustavo.dudamel@edugo.local", schoolID: schoolUCVID, unitID: unitUCVS1ID},
	{id: "69000000-0000-0000-0000-00000000132f", teacherEmail: "gustavo.dudamel@edugo.local", schoolID: schoolInglesLabID, unitID: unitINGLNBID},
	// Luis Silva (101c) → UCV Primer Semestre.
	{id: "69000000-0000-0000-0000-000000001330", teacherEmail: "luis.silva@edugo.local", schoolID: schoolUCVID, unitID: unitUCVS1ID},
	// Sheldon Cooper (101d) → UCV Primer Semestre.
	{id: "69000000-0000-0000-0000-000000001331", teacherEmail: "sheldon.cooper@edugo.local", schoolID: schoolUCVID, unitID: unitUCVS1ID},
}

// adminMembershipSeed describe una membresía "admin" con alcance COLEGIO.
type adminMembershipSeed struct {
	id, schoolID string
}

var adminMembershipSeeds = []adminMembershipSeed{
	{id: adminMembCristoReyID, schoolID: schoolCristoReyID},
	{id: adminMembUCVID, schoolID: schoolUCVID},
	{id: adminMembInglesLabID, schoolID: schoolInglesLabID},
}

// Users lista los 15 usuarios que el escenario necesita en auth.users: el
// admin global, los 13 profesores y el solicitante N0.
func Users() []common.ScenarioUser {
	out := make([]common.ScenarioUser, 0, len(teacherSeeds)+2)
	out = append(out, common.ScenarioUser{Email: AdminEmail, FirstName: "Admin", LastName: "Total"})
	for _, t := range teacherSeeds {
		out = append(out, common.ScenarioUser{Email: t.email, FirstName: t.first, LastName: t.last})
	}
	out = append(out, common.ScenarioUser{Email: SolicitanteEmail, FirstName: "Carlos", LastName: "Estudiante"})
	return out
}

// Apply siembra el playground n0n1_escuelas. Asume que L0..L4 corrieron (rol
// super_admin y concept_types primary_school/workshop ya existen) y que users
// trae el id de cada usuario de Users(). Orden: concept_type university →
// schools → periods → units → subjects → subject_offerings → user_role admin →
// memberships → join_request. Idempotente.
func Apply(tx *gorm.DB, users common.UserIDs) error {
	if err := users.Require(Users()); err != nil {
		return fmt.Errorf("playground_v2/n0n1_escuelas: %w", err)
	}
	if err := upsertConceptType(tx); err != nil {
		return fmt.Errorf("playground_v2/n0n1_escuelas: concept_type: %w", err)
	}
	for _, s := range schoolSeeds {
		if err := upsertSchool(tx, s); err != nil {
			return fmt.Errorf("playground_v2/n0n1_escuelas: school %s: %w", s.code, err)
		}
	}
	for _, p := range periodSeeds {
		if err := upsertPeriod(tx, p); err != nil {
			return fmt.Errorf("playground_v2/n0n1_escuelas: period %s: %w", p.code, err)
		}
	}
	for _, u := range unitSeeds {
		if err := upsertUnit(tx, u); err != nil {
			return fmt.Errorf("playground_v2/n0n1_escuelas: unit %s: %w", u.code, err)
		}
	}
	for _, s := range subjectSeeds {
		if err := upsertSubject(tx, s); err != nil {
			return fmt.Errorf("playground_v2/n0n1_escuelas: subject %s: %w", s.code, err)
		}
	}
	if err := upsertOfferings(tx); err != nil {
		return fmt.Errorf("playground_v2/n0n1_escuelas: offerings: %w", err)
	}

	// Usuarios: los crea la API de identidad; aquí solo se usan sus ids.
	adminUser := users[AdminEmail]

	// UserRole admin → super_admin (acceso global "*" vía BeforeSave del UserRole).
	// El id se deriva SHA1(userID:roleID) en el común; la constante explícita
	// adminUserRoleID que usaba el seed original no se referencia en ningún otro
	// lado, así que el resultado es funcionalmente idéntico (mismo vínculo, scope "*").
	if err := common.SeedUserRole(tx, adminUser, common.MustParseUUID(superAdminRoleID)); err != nil {
		return fmt.Errorf("playground_v2/n0n1_escuelas: admin_user_role: %w", err)
	}

	// Membresías admin (alcance colegio) en las 3 escuelas.
	for _, m := range adminMembershipSeeds {
		if err := common.SeedMembership(tx, common.MembershipSpec{
			ID: common.MustParseUUID(m.id), UserID: adminUser, SchoolID: common.MustParseUUID(m.schoolID), AcademicUnitID: nil, Role: "admin",
		}); err != nil {
			return fmt.Errorf("playground_v2/n0n1_escuelas: admin_membership %s: %w", m.id, err)
		}
	}
	// Membresías teacher (alcance unidad).
	for _, m := range teacherMembershipSeeds {
		unitID := common.MustParseUUID(m.unitID)
		if err := common.SeedMembership(tx, common.MembershipSpec{
			ID: common.MustParseUUID(m.id), UserID: users[m.teacherEmail], SchoolID: common.MustParseUUID(m.schoolID), AcademicUnitID: &unitID, Role: "teacher",
		}); err != nil {
			return fmt.Errorf("playground_v2/n0n1_escuelas: teacher_membership %s: %w", m.id, err)
		}
	}

	// School join request pendiente (N0): Carlos → InglesLab / Nivel Básico.
	if err := upsertJoinRequest(tx, users[SolicitanteEmail]); err != nil {
		return fmt.Errorf("playground_v2/n0n1_escuelas: join_request: %w", err)
	}

	return nil
}

func upsertConceptType(tx *gorm.DB) error {
	id, err := uuid.Parse(universityConceptType.id)
	if err != nil {
		return err
	}
	desc := universityConceptType.desc
	ct := entities.ConceptType{
		ID:          id,
		Name:        universityConceptType.name,
		Code:        universityConceptType.code,
		Description: &desc,
		IsActive:    true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&ct).Error
}

func upsertSchool(tx *gorm.DB, s schoolSeed) error {
	ctID := common.MustParseUUID(s.conceptTypeID)
	return common.SeedSchool(tx, common.SchoolSpec{
		ID:            common.MustParseUUID(s.id),
		Name:          s.name,
		Code:          s.code,
		Country:       "Venezuela",
		ConceptTypeID: &ctID,
	})
}

func upsertPeriod(tx *gorm.DB, p periodSeed) error {
	// AcademicUnitID se deja en uuid.Nil (período por escuela, sin anclar a
	// unidad), igual que el seed original.
	return common.SeedActivePeriod(tx, common.PeriodSpec{
		ID:           common.MustParseUUID(p.id),
		SchoolID:     common.MustParseUUID(p.schoolID),
		Name:         p.name,
		Code:         p.code,
		Type:         "semester",
		StartDate:    time.Date(p.startY, time.Month(p.startM), p.startD, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(p.endY, time.Month(p.endM), p.endD, 0, 0, 0, 0, time.UTC),
		AcademicYear: academicYear,
		SortOrder:    0,
	})
}

func upsertUnit(tx *gorm.DB, u unitSeed) error {
	return common.SeedAcademicUnit(tx, common.UnitSpec{
		ID:           common.MustParseUUID(u.id),
		SchoolID:     common.MustParseUUID(u.schoolID),
		Name:         u.name,
		Code:         u.code,
		Type:         "grade",
		AcademicYear: academicYear,
	})
}

func upsertSubject(tx *gorm.DB, s subjectSeed) error {
	return common.SeedSubject(tx, common.SubjectSpec{
		ID:       common.MustParseUUID(s.id),
		SchoolID: common.MustParseUUID(s.schoolID),
		Name:     s.name,
		Code:     s.code,
	})
}

// upsertOfferings genera las 32 subject_offerings por loop sobre offeringGroups
// (subjects × sections), asignando ids secuenciales 69..0201..69..0220. Todas
// SIN docente (teacher_membership_id NULL) y sin capacity. Idempotente por id.
func upsertOfferings(tx *gorm.DB) error {
	seq := 0x201 // primer id secuencial: 69..0201.
	for _, g := range offeringGroups {
		sid := common.MustParseUUID(g.schoolID)
		auid := common.MustParseUUID(g.unitID)
		pid := common.MustParseUUID(g.periodID)
		for _, subjIDStr := range g.subjectIDs {
			subjID := common.MustParseUUID(subjIDStr)
			for _, section := range g.sections {
				id := common.MustParseUUID(fmt.Sprintf("69000000-0000-0000-0000-%012x", seq))
				seq++
				label := section
				if err := common.SeedOffering(tx, common.OfferingSpec{
					ID:                  id,
					SchoolID:            sid,
					SubjectID:           subjID,
					AcademicUnitID:      auid,
					PeriodID:            pid,
					SectionLabel:        &label,
					TeacherMembershipID: nil,
					Capacity:            nil,
				}); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// upsertJoinRequest crea la solicitud N0 PENDIENTE (sin firmas) del solicitante
// Carlos a InglesLab / Nivel Básico, role "student". Idempotente por id.
func upsertJoinRequest(tx *gorm.DB, uid uuid.UUID) error {
	id, err := uuid.Parse(joinRequestID)
	if err != nil {
		return err
	}
	sid, err := uuid.Parse(schoolInglesLabID)
	if err != nil {
		return err
	}
	auid, err := uuid.Parse(unitINGLNBID)
	if err != nil {
		return err
	}
	invitationTypeID, err := catalog.ResolveInvitationTypeID(tx, "student")
	if err != nil {
		return err
	}
	r := entities.SchoolJoinRequest{
		ID:               id,
		UserID:           uid,
		SchoolID:         sid,
		AcademicUnitID:   auid,
		InvitationTypeID: invitationTypeID,
		InvitationID:     nil,
		Status:           "pending",
		RequestedAt:      time.Now().UTC(),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&r).Error
}
