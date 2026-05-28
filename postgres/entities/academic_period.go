package entities

import (
	"time"

	"github.com/google/uuid"
)

// AcademicPeriod representa la tabla 'academic_periods' en PostgreSQL.
//
// Migracion: 039_academic_periods.sql
type AcademicPeriod struct {
	ID           uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	SchoolID     uuid.UUID `db:"school_id" gorm:"type:uuid;index;not null;constraint:academic_periods_school_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	Name         string    `db:"name" gorm:"not null;type:varchar(100)" validate:"required,min=2,max=255"`
	Code         *string   `db:"code" gorm:"type:varchar(20)" validate:"omitempty"`
	Type         string    `db:"type" gorm:"not null;type:varchar(20);check:academic_periods_type_check,type IN ('semester','trimester','bimester','quarter')" validate:"required,oneof=semester trimester bimester quarter"`
	StartDate    time.Time `db:"start_date" gorm:"not null;type:date"`
	EndDate      time.Time `db:"end_date" gorm:"not null;type:date"`
	// NOTE: partial unique index idx_academic_periods_active ON (school_id) WHERE is_active = true must be created in post_gorm.sql
	IsActive     bool      `db:"is_active" gorm:"not null;default:false"`
	AcademicYear int       `db:"academic_year" gorm:"not null" validate:"required"`
	SortOrder    int       `db:"sort_order" gorm:"default:0" validate:"omitempty"`
	CreatedAt    time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt    time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

func (AcademicPeriod) TableName() string {
	return "academic.academic_periods"
}
