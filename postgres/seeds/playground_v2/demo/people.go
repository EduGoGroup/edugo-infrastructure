package demo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/catalog"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/common"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// fillerNames son los dieciocho alumnos de relleno, sin login: dan cuerpo a las
// listas y son los dueños de los intentos por corregir de E-R.
var fillerNames = [][2]string{
	{"Benjamín", "Soto"}, {"Florencia", "Díaz"}, {"Vicente", "Morales"}, {"Isidora", "Reyes"},
	{"Tomás", "Castillo"}, {"Josefa", "Espinoza"}, {"Agustín", "Torres"}, {"Emilia", "Contreras"},
	{"Joaquín", "Silva"}, {"Martina", "Sepúlveda"}, {"Lucas", "Pizarro"}, {"Catalina", "Araya"},
	{"Maximiliano", "Rojas"}, {"Sofía", "Navarro"}, {"Cristóbal", "Vergara"}, {"Amanda", "Carrasco"},
	{"Felipe", "Molina"}, {"Trinidad", "Fuentealba"},
}

// FillerCount es el número de alumnos de relleno.
var FillerCount = len(fillerNames)

// membershipRow es una membresía por sembrar; Role es la key del tipo de
// invitación, que se resuelve a su id al escribir.
type membershipRow struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Role   string
}

// student es un alumno de demo, con o sin login.
type student struct {
	UserID       uuid.UUID
	MembershipID uuid.UUID
}

// students devuelve los veinte alumnos: primero Matías y Antonia, luego el relleno.
func students() []student {
	out := []student{
		{UserID: uuid.MustParse(Student1UserID), MembershipID: Student1MembershipID},
		{UserID: uuid.MustParse(Student2UserID), MembershipID: Student2MembershipID},
	}
	for i := range fillerNames {
		out = append(out, student{UserID: fillerUserID(i), MembershipID: fillerMembershipID(i)})
	}
	return out
}

// fillerEmail arma el correo de un alumno de relleno, en un subdominio propio
// para que no se confunda con una cuenta que entra.
func fillerEmail(first, last string) string {
	r := strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ñ", "n")
	return r.Replace(strings.ToLower(first+"."+last)) + "@relleno.demo.edugo.test"
}

// buildPeople arma los usuarios, sus roles (con school_id, como `base`), las
// membresías en la unidad de trabajo (la tutora no lleva), las inscripciones de
// cada alumno en las tres sesiones y el vínculo de la tutora con Antonia.
func buildPeople(w *world) {
	teacherID := uuid.MustParse(TeacherUserID)
	guardianID := uuid.MustParse(GuardianUserID)
	student2ID := uuid.MustParse(Student2UserID)

	w.users = []common.UserSpec{
		{ID: teacherID, Email: TeacherEmail, FirstName: "Valentina", LastName: "Rojas"},
		{ID: uuid.MustParse(Student1UserID), Email: Student1Email, FirstName: "Matías", LastName: "Fuentes"},
		{ID: student2ID, Email: Student2Email, FirstName: "Antonia", LastName: "Herrera"},
		{ID: guardianID, Email: GuardianEmail, FirstName: "Carolina", LastName: "Muñoz"},
	}
	for i, n := range fillerNames {
		w.users = append(w.users, common.UserSpec{
			ID: fillerUserID(i), Email: fillerEmail(n[0], n[1]), FirstName: n[0], LastName: n[1],
		})
	}

	grantedAt := w.day(periodStartDays)
	school := SchoolID
	role := func(n int, userID uuid.UUID, roleID string) entities.UserRole {
		return entities.UserRole{
			ID: newID(domUserRole, n), UserID: userID, RoleID: uuid.MustParse(roleID),
			SchoolID: &school, IsActive: true, GrantedAt: grantedAt,
		}
	}
	w.userRoles = []entities.UserRole{
		role(1, teacherID, l4.L4_ROLE_TEACHER_ID),
		role(2, guardianID, l4.L4_ROLE_GUARDIAN_ID),
	}
	w.memberships = []membershipRow{{ID: TeacherMembershipID, UserID: teacherID, Role: "teacher"}}

	for i, s := range students() {
		w.userRoles = append(w.userRoles, role(10+i, s.UserID, l4.L4_ROLE_STUDENT_ID))
		w.memberships = append(w.memberships, membershipRow{ID: s.MembershipID, UserID: s.UserID, Role: "student"})
		for _, o := range w.offerings {
			w.enrollments = append(w.enrollments, entities.SubjectOfferingEnrollment{
				OfferingID: o.ID, SubjectID: o.SubjectID, PeriodID: o.PeriodID,
				StudentMembershipID: s.MembershipID, EnrolledAt: grantedAt,
			})
		}
	}

	w.guardian = common.GuardianRelationSpec{
		ID: GuardianRelationID, GuardianID: guardianID, StudentID: student2ID,
		SchoolID: SchoolID, RelationshipType: "parent", IsPrimary: true, IsActive: true, Status: "active",
	}
}

// writePeople escribe usuarios, membresías, roles, sesiones, inscripciones y el
// vínculo de la tutora, en el orden de sus claves foráneas.
func writePeople(tx *gorm.DB, w *world) error {
	for _, u := range w.users {
		if err := common.SeedUser(tx, u); err != nil {
			return fmt.Errorf("usuario %s: %w", u.Email, err)
		}
	}
	if err := writeMemberships(tx, w); err != nil {
		return err
	}
	if err := insert(tx, &w.userRoles); err != nil {
		return fmt.Errorf("roles: %w", err)
	}
	if err := writeOfferings(tx, w); err != nil {
		return err
	}
	if err := insert(tx, &w.enrollments); err != nil {
		return fmt.Errorf("inscripciones: %w", err)
	}
	if err := common.SeedGuardianRelation(tx, w.guardian); err != nil {
		return fmt.Errorf("vínculo de la tutora: %w", err)
	}
	return nil
}

// writeMemberships resuelve cada key de tipo de invitación (las siembra L4) y
// escribe las membresías activas en la unidad de trabajo, con enrolled_at al
// inicio del período.
func writeMemberships(tx *gorm.DB, w *world) error {
	typeIDs := map[string]uuid.UUID{}
	unit := UnitID
	rows := make([]entities.Membership, 0, len(w.memberships))
	for _, m := range w.memberships {
		typeID, ok := typeIDs[m.Role]
		if !ok {
			id, err := catalog.ResolveInvitationTypeID(tx, m.Role)
			if err != nil {
				return fmt.Errorf("membresías: %w", err)
			}
			typeIDs[m.Role], typeID = id, id
		}
		rows = append(rows, entities.Membership{
			ID: m.ID, UserID: m.UserID, SchoolID: SchoolID, AcademicUnitID: &unit,
			InvitationTypeID: typeID, Metadata: json.RawMessage(`{}`), Status: "active",
			EnrolledAt: w.day(periodStartDays),
		})
	}
	if err := insert(tx, &rows); err != nil {
		return fmt.Errorf("membresías: %w", err)
	}
	return nil
}
