package entities

import "github.com/google/uuid"

// MembershipSubject representa la tabla 'membership_subjects' en PostgreSQL.
// Tabla de union entre memberships y subjects con FK reales.
//
// Migracion: 038_academic_membership_subjects.sql
type MembershipSubject struct {
	MembershipID uuid.UUID `db:"membership_id" gorm:"type:uuid;primaryKey"`
	SubjectID    uuid.UUID `db:"subject_id" gorm:"type:uuid;primaryKey"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (MembershipSubject) TableName() string {
	return "academic.membership_subjects"
}
