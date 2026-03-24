package entities

import (
	"time"

	"github.com/google/uuid"
)

// Schedule representa la tabla 'schedules' en PostgreSQL.
//
// Migracion: 092_academic_schedules.sql
type Schedule struct {
	ID                  uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey"`
	AcademicUnitID      uuid.UUID  `db:"academic_unit_id" gorm:"type:uuid;index;not null"`
	SubjectID           uuid.UUID  `db:"subject_id" gorm:"type:uuid;not null"`
	TeacherMembershipID uuid.UUID  `db:"teacher_membership_id" gorm:"type:uuid;index;not null"`
	DayOfWeek           int        `db:"day_of_week" gorm:"not null"`
	StartTime           string     `db:"start_time" gorm:"not null;type:time"`
	EndTime             string     `db:"end_time" gorm:"not null;type:time"`
	Room                *string    `db:"room" gorm:"type:varchar(50)"`
	PeriodID            *uuid.UUID `db:"period_id" gorm:"type:uuid"`
	IsActive            bool       `db:"is_active" gorm:"not null;default:true"`
	CreatedAt           time.Time  `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt           time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

func (Schedule) TableName() string {
	return "academic.schedules"
}
