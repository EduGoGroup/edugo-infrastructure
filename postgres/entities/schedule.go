package entities

import (
	"time"

	"github.com/google/uuid"
)

// Schedule representa la tabla 'schedules' en PostgreSQL.
//
// Migracion: 092_academic_schedules.sql
type Schedule struct {
	ID                  uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	AcademicUnitID      uuid.UUID  `db:"academic_unit_id" gorm:"type:uuid;index;not null;constraint:schedules_unit_fkey,OnDelete:CASCADE;index:idx_schedules_unit" validate:"required,uuid"`
	SubjectID           uuid.UUID  `db:"subject_id" gorm:"type:uuid;not null;constraint:schedules_subject_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	TeacherMembershipID uuid.UUID  `db:"teacher_membership_id" gorm:"type:uuid;index;not null;constraint:schedules_teacher_fkey,OnDelete:CASCADE;index:idx_schedules_teacher" validate:"required,uuid"`
	DayOfWeek           int        `db:"day_of_week" gorm:"not null;check:schedules_dow_check,day_of_week BETWEEN 0 AND 6" validate:"required"`
	StartTime           string     `db:"start_time" gorm:"not null;type:time" validate:"required"`
	EndTime             string     `db:"end_time" gorm:"not null;type:time" validate:"required"`
	Room                *string    `db:"room" gorm:"type:varchar(50)" validate:"omitempty"`
	PeriodID            *uuid.UUID `db:"period_id" gorm:"type:uuid;constraint:schedules_period_fkey" validate:"omitempty,uuid"`
	IsActive            bool       `db:"is_active" gorm:"not null;default:true"`
	CreatedAt           time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt           time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

func (Schedule) TableName() string {
	return "academic.schedules"
}
