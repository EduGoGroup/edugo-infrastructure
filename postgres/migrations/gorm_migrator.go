package migrations

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// openGORM wraps an existing *sql.DB connection into a *gorm.DB,
// reusing the caller's connection pool without opening a new one.
func openGORM(db *sql.DB) (*gorm.DB, error) {
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("error abriendo GORM desde *sql.DB: %w", err)
	}
	return gdb, nil
}

// autoMigrateAll runs gorm.AutoMigrate for all system entities.
// Order matters: entities with no external FKs first, dependents after.
// AutoMigrate is idempotent: existing tables are only altered to add missing columns/indexes.
func autoMigrateAll(gdb *gorm.DB) error {
	return gdb.AutoMigrate(
		// Auth (no cross-schema deps)
		&entities.User{},
		&entities.RefreshToken{},
		&entities.LoginAttempt{},

		// IAM (depends on iam ENUM types created in pre_gorm.sql)
		&entities.Resource{},
		&entities.Role{},
		&entities.Permission{},
		&entities.UserRole{},
		&entities.RoleGrant{},
		&entities.UserGrant{},

		// Academic base (concept_types before schools due to FK)
		&entities.ConceptType{},
		&entities.School{},
		&entities.AcademicUnit{},
		&entities.Membership{},
		&entities.Subject{},
		&entities.GuardianRelation{},
		&entities.SchoolInvitation{},
		&entities.SchoolJoinRequest{},
		&entities.ConceptDefinition{},
		&entities.SchoolConcept{},
		&entities.AcademicPeriod{},
		// Sesiones de materia (ADR 0009 / plan 010 N1.7). subject_offerings
		// referencia subjects, academic_units, academic_periods y memberships
		// (todos migrados arriba); enrollments referencia subject_offerings y
		// memberships, por eso va inmediatamente despues.
		&entities.SubjectOffering{},
		&entities.SubjectOfferingEnrollment{},
		&entities.Grade{},
		&entities.Attendance{},
		&entities.Schedule{},
		&entities.Announcement{},
		&entities.CalendarEvent{},
		&entities.Color{},

		// Content
		&entities.Course{},
		&entities.Material{},
		&entities.MaterialVersion{},
		&entities.Progress{},

		// Assessment (assessment first, then dependents)
		&entities.Assessment{},
		&entities.AssessmentAttempt{},
		&entities.AssessmentAttemptAnswer{},
		&entities.AssessmentMaterial{},
		&entities.Question{},
		&entities.QuestionOption{},
		&entities.AssessmentAssignment{},
		&entities.AttemptReview{},

		// UI Config (templates before instances due to FK)
		&entities.ScreenTemplate{},
		&entities.ScreenInstance{},
		&entities.ResourceScreen{},
		&entities.ScreenUserPreference{},

		// Audit
		&entities.AuditEvent{},

		// Notifications
		&entities.Notification{},
	)
}
