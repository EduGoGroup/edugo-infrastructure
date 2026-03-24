package entities

import (
	"time"

	"github.com/google/uuid"
)

// AcademicPeriod representa la tabla 'academic_periods' en PostgreSQL.
//
// Migracion: 039_academic_periods.sql
type AcademicPeriod struct {
	ID           uuid.UUID `db:"id" gorm:"type:uuid;primaryKey"`
	SchoolID     uuid.UUID `db:"school_id" gorm:"type:uuid;index;not null"`
	Name         string    `db:"name" gorm:"not null;type:varchar(100)"`
	Code         *string   `db:"code" gorm:"type:varchar(20)"`
	Type         string    `db:"type" gorm:"not null;type:varchar(20)"`
	StartDate    time.Time `db:"start_date" gorm:"not null;type:date"`
	EndDate      time.Time `db:"end_date" gorm:"not null;type:date"`
	IsActive     bool      `db:"is_active" gorm:"not null;default:false"`
	AcademicYear int       `db:"academic_year" gorm:"not null"`
	SortOrder    int       `db:"sort_order" gorm:"default:0"`
	CreatedAt    time.Time `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

func (AcademicPeriod) TableName() string {
	return "academic.academic_periods"
}
