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
	ID             uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey"`
	SchoolID       uuid.UUID      `db:"school_id" gorm:"type:uuid;index;not null"`
	AcademicUnitID *uuid.UUID     `db:"academic_unit_id" gorm:"type:uuid"`
	AuthorID       uuid.UUID      `db:"author_id" gorm:"type:uuid;not null"`
	Title          string         `db:"title" gorm:"not null;type:varchar(200)"`
	Body           string         `db:"body" gorm:"not null"`
	Scope          string         `db:"scope" gorm:"not null;type:varchar(20)"`
	TargetRoles    pq.StringArray `db:"target_roles" gorm:"type:text[]"`
	IsPinned       bool           `db:"is_pinned" gorm:"not null;default:false"`
	PublishedAt    *time.Time     `db:"published_at"`
	ExpiresAt      *time.Time     `db:"expires_at"`
	CreatedAt      time.Time      `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt      time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

func (Announcement) TableName() string {
	return "academic.announcements"
}
