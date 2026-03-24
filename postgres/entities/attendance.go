package entities

import (
	"time"

	"github.com/google/uuid"
)

// Attendance representa la tabla 'attendance' en PostgreSQL.
//
// Migracion: 091_academic_attendance.sql
type Attendance struct {
	ID           uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey"`
	MembershipID uuid.UUID  `db:"membership_id" gorm:"type:uuid;index;not null"`
	SubjectID    *uuid.UUID `db:"subject_id" gorm:"type:uuid"`
	Date         time.Time  `db:"date" gorm:"not null;type:date"`
	Status       string     `db:"status" gorm:"not null;type:varchar(20)"`
	RecordedBy   uuid.UUID  `db:"recorded_by" gorm:"type:uuid;not null"`
	Notes        *string    `db:"notes"`
	CreatedAt    time.Time  `db:"created_at" gorm:"not null;autoCreateTime"`
}

func (Attendance) TableName() string {
	return "academic.attendance"
}
