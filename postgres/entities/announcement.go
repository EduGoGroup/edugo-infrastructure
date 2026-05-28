package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Announcement representa la tabla 'announcements' en PostgreSQL.
//
// Migracion: 093_academic_announcements.sql
type Announcement struct {
	ID             uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	SchoolID       uuid.UUID      `db:"school_id" gorm:"type:uuid;index;not null;constraint:announcements_school_fkey,OnDelete:CASCADE;index:idx_announcements_school" validate:"required,uuid"`
	AcademicUnitID *uuid.UUID     `db:"academic_unit_id" gorm:"type:uuid;constraint:announcements_unit_fkey,OnDelete:CASCADE" validate:"omitempty,uuid"`
	AuthorID       uuid.UUID      `db:"author_id" gorm:"type:uuid;not null" validate:"required,uuid"`
	Title          string         `db:"title" gorm:"not null;type:varchar(200)" validate:"required,min=2,max=255"`
	Body           string         `db:"body" gorm:"not null" validate:"required"`
	Scope          string         `db:"scope" gorm:"not null;type:varchar(20);check:announcements_scope_check,scope IN ('school','unit','role')" validate:"required,oneof=school unit role"`
	TargetRoles    pq.StringArray `db:"target_roles" gorm:"type:text[]" validate:"-"`
	IsPinned       bool           `db:"is_pinned" gorm:"not null;default:false"`
	PublishedAt    *time.Time     `db:"published_at"`
	ExpiresAt      *time.Time     `db:"expires_at"`
	CreatedAt      time.Time      `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt      time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

func (Announcement) TableName() string {
	return "academic.announcements"
}
