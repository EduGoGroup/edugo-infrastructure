package entities

import (
	"time"

	"github.com/google/uuid"
)

// CalendarEvent representa la tabla 'calendar_events' en PostgreSQL.
//
// Migracion: 094_academic_calendar_events.sql
type CalendarEvent struct {
	ID          uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey"`
	SchoolID    uuid.UUID  `db:"school_id" gorm:"type:uuid;index;not null"`
	Title       string     `db:"title" gorm:"not null;type:varchar(200)"`
	Description *string    `db:"description"`
	EventType   string     `db:"event_type" gorm:"not null;type:varchar(30)"`
	StartDate   time.Time  `db:"start_date" gorm:"not null;type:date"`
	EndDate     *time.Time `db:"end_date" gorm:"type:date"`
	IsAllDay    bool       `db:"is_all_day" gorm:"not null;default:true"`
	CreatedBy   uuid.UUID  `db:"created_by" gorm:"type:uuid;not null"`
	CreatedAt   time.Time  `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

func (CalendarEvent) TableName() string {
	return "academic.calendar_events"
}
