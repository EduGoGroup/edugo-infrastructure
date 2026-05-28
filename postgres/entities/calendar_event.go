package entities

import (
	"time"

	"github.com/google/uuid"
)

// CalendarEvent representa la tabla 'calendar_events' en PostgreSQL.
//
// Migracion: 094_academic_calendar_events.sql
type CalendarEvent struct {
	ID          uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	SchoolID    uuid.UUID  `db:"school_id" gorm:"type:uuid;index;not null;constraint:calendar_events_school_fkey,OnDelete:CASCADE;index:idx_calendar_school" validate:"required,uuid"`
	Title       string     `db:"title" gorm:"not null;type:varchar(200)" validate:"required,min=2,max=255"`
	Description *string    `db:"description" validate:"omitempty"`
	EventType   string     `db:"event_type" gorm:"not null;type:varchar(30);check:calendar_events_type_check,event_type IN ('holiday','exam','meeting','activity','deadline')" validate:"required,oneof=holiday exam meeting activity deadline"`
	StartDate   time.Time  `db:"start_date" gorm:"not null;type:date;index:idx_calendar_dates"`
	EndDate     *time.Time `db:"end_date" gorm:"type:date;index:idx_calendar_dates"`
	IsAllDay    bool       `db:"is_all_day" gorm:"not null;default:true"`
	CreatedBy   uuid.UUID  `db:"created_by" gorm:"type:uuid;not null" validate:"required,uuid"`
	CreatedAt   time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt   time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

func (CalendarEvent) TableName() string {
	return "academic.calendar_events"
}
