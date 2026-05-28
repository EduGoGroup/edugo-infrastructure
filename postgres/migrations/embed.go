package migrations

import (
	"database/sql"
	"embed"
	"fmt"
)

// sqlFiles embeds the two SQL scripts that bracket gorm.AutoMigrate.
//
//go:embed sql/pre_gorm.sql sql/post_gorm.sql
var sqlFiles embed.FS

// ApplyPreGORM executes schemas, extensions, ENUM types, and shared functions.
// Must run BEFORE gorm.AutoMigrate().
func ApplyPreGORM(db *sql.DB) error {
	content, err := sqlFiles.ReadFile("sql/pre_gorm.sql")
	if err != nil {
		return fmt.Errorf("error reading pre_gorm.sql: %w", err)
	}
	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("error executing pre_gorm.sql: %w", err)
	}
	return nil
}

// ApplyAll executes the full migration in order: ApplyPreGORM → AutoMigrate → ApplyPostGORM.
// Kept for backward compatibility with existing tooling (cmd/runner).
// Prefer calling Migrate() directly for full control over options.
func ApplyAll(db *sql.DB) error {
	if err := ApplyPreGORM(db); err != nil {
		return err
	}
	gdb, err := openGORM(db)
	if err != nil {
		return err
	}
	if err := autoMigrateAll(gdb); err != nil {
		return err
	}
	return ApplyPostGORM(db)
}

// ApplyPostGORM executes triggers, views, IAM functions, partial indexes,
// analytics tables, and any FKs that GORM cannot express via entity tags.
// Must run AFTER gorm.AutoMigrate().
func ApplyPostGORM(db *sql.DB) error {
	content, err := sqlFiles.ReadFile("sql/post_gorm.sql")
	if err != nil {
		return fmt.Errorf("error reading post_gorm.sql: %w", err)
	}
	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("error executing post_gorm.sql: %w", err)
	}
	return nil
}
