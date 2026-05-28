package entities

import (
	"time"

	"github.com/google/uuid"
)

// Attendance representa la tabla 'attendance' en PostgreSQL.
//
// Migracion: 091_academic_attendance.sql
type Attendance struct {
	ID           uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	MembershipID uuid.UUID  `db:"membership_id" gorm:"type:uuid;index;not null;constraint:attendance_membership_fkey,OnDelete:CASCADE;uniqueIndex:attendance_unique" validate:"required,uuid"`
	SubjectID    *uuid.UUID `db:"subject_id" gorm:"type:uuid;constraint:attendance_subject_fkey,OnDelete:CASCADE;uniqueIndex:attendance_unique" validate:"omitempty,uuid"`
	Date         time.Time  `db:"date" gorm:"not null;type:date;uniqueIndex:attendance_unique;index:idx_attendance_date"`
	Status       string     `db:"status" gorm:"not null;type:varchar(20);check:attendance_status_check,status IN ('present','absent','late','excused','remote')" validate:"required,oneof=present absent late excused remote"`
	RecordedBy   uuid.UUID  `db:"recorded_by" gorm:"type:uuid;not null" validate:"required,uuid"`
	Notes        *string    `db:"notes" validate:"omitempty"`
	CreatedAt    time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
}

func (Attendance) TableName() string {
	return "academic.attendance"
}
