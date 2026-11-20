package demo

import (
	"fmt"

	"github.com/google/uuid"
)

// Las cuatro cuentas con login. Su id NO lo elige el fixture: es el `sub` que
// identity-core asignó al crearlas (2026-11-20), porque el canje busca la fila
// de auth.users por ese `sub` (ADR 0041). Un UUID no es un secreto; las claves
// viven fuera de git.
const (
	TeacherUserID  = "be0fc9bb-f76a-49c3-95e9-5d69e765a7db" // Valentina Rojas, docente
	Student1UserID = "d4544906-e877-4f8d-8d0d-0e41d252fdbf" // Matías Fuentes, alumno 1
	Student2UserID = "585da1ba-af33-4588-b086-6467f85a1642" // Antonia Herrera, alumna 2
	GuardianUserID = "bf56600b-fed0-4ff7-bf2c-45f513a37ccd" // Carolina Muñoz, tutora

	TeacherEmail  = "valentina.rojas@demo.edugo.test"
	Student1Email = "matias.fuentes@demo.edugo.test"
	Student2Email = "antonia.herrera@demo.edugo.test"
	GuardianEmail = "carolina.munoz@demo.edugo.test"
)

// domain es el tercer y cuarto dígito hexadecimal del rango propio del fixture:
// todo id que el fixture elige tiene la forma de0000DD-0000-0000-0000-NNNNNNNNNNNN,
// con DD el dominio y N un contador. El rango de…… estaba libre (2026-11-19).
type domain uint8

const (
	domSchool domain = iota + 1
	domUnit
	domPeriod
	domSubject
	domOffering
	domUser
	domMembership
	domUserRole
	domGuardianRelation
	domAssessment
	domQuestion
	domOption
	domAssignment
	domAttempt
	domAnswer
	domMaterial
	domMaterialFile
	domMaterialAssignment
	domGrade
	domAttendance
	domAnnouncement
	domUserGrant
	domSchoolConcept
)

// IDPrefix es el prefijo de todo id que elige el fixture.
const IDPrefix = "de"

// newID arma el id n del dominio d dentro del rango propio.
func newID(d domain, n int) uuid.UUID {
	return uuid.MustParse(fmt.Sprintf("%s0000%02x-0000-0000-0000-%012x", IDPrefix, uint8(d), n))
}

// Ids fijos de la estructura, exportados para las pruebas y para la AR-DEMO.2.
var (
	SchoolID   = newID(domSchool, 1)
	RootUnitID = newID(domUnit, 1)
	UnitID     = newID(domUnit, 2) // «8.º Básico A», la unidad de trabajo
	PeriodID   = newID(domPeriod, 1)

	SubjectMathID     = newID(domSubject, 1)
	SubjectLanguageID = newID(domSubject, 2)
	SubjectScienceID  = newID(domSubject, 3)

	OfferingMathID     = newID(domOffering, 1)
	OfferingLanguageID = newID(domOffering, 2)
	OfferingScienceID  = newID(domOffering, 3)

	TeacherMembershipID  = newID(domMembership, 1)
	Student1MembershipID = newID(domMembership, 2)
	Student2MembershipID = newID(domMembership, 3)

	GuardianRelationID = newID(domGuardianRelation, 1)

	AssessmentDraftID    = newID(domAssessment, 1) // E-B
	AssessmentExamID     = newID(domAssessment, 2) // E-X
	AssessmentPracticeID = newID(domAssessment, 3) // E-P
	AssessmentReviewID   = newID(domAssessment, 4) // E-R

	MaterialID     = newID(domMaterial, 1)     // T1
	MaterialFileID = newID(domMaterialFile, 1) // F1
)

// fillerMembershipID es la membresía del alumno de relleno i (0-based); las tres
// primeras membresías son las de las cuentas con login.
func fillerMembershipID(i int) uuid.UUID { return newID(domMembership, 10+i) }

// fillerUserID es el usuario del alumno de relleno i (0-based).
func fillerUserID(i int) uuid.UUID { return newID(domUser, 1+i) }

// MaterialFileKey es la clave del objeto PDF de T1 en el bucket. El fixture la
// escribe en content.material_file.file_url (learning usa ese valor como clave,
// file_download_url.go) y la AR-DEMO.2 sube el PDF ahí.
func MaterialFileKey() string {
	return fmt.Sprintf("materials/%s/files/%s", MaterialID, MaterialFileID)
}

// NotificationID es la fórmula de la API (learning, notification/repository.go:77;
// platform, notification/repository.go:108): SHA-1 sobre el OID de
// user+type+resource_type+resource_id, sin separadores. Sembrar el mismo id que
// escribiría la API es lo que hace idempotente un aviso repetido.
func NotificationID(userID uuid.UUID, typ, resourceType string, resourceID uuid.UUID) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(userID.String()+typ+resourceType+resourceID.String()))
}
