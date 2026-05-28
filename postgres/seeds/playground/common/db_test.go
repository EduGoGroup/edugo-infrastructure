package common

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type fixture struct {
	ID   string `gorm:"primaryKey"`
	Name string
}

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open failed: %v", err)
	}
	if err := db.AutoMigrate(&fixture{}); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
	return db
}

func countFixtures(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	var n int64
	if err := db.Model(&fixture{}).Count(&n).Error; err != nil {
		t.Fatalf("Count failed: %v", err)
	}
	return n
}

func TestOnConflictIgnore_FirstInsert(t *testing.T) {
	db := newTestDB(t)
	if err := OnConflictIgnore(db, &fixture{ID: "a", Name: "first"}); err != nil {
		t.Fatalf("OnConflictIgnore failed: %v", err)
	}
	if got := countFixtures(t, db); got != 1 {
		t.Fatalf("expected count=1, got %d", got)
	}
}

func TestOnConflictIgnore_SecondInsertIgnored(t *testing.T) {
	db := newTestDB(t)
	if err := OnConflictIgnore(db, &fixture{ID: "a", Name: "first"}); err != nil {
		t.Fatalf("first OnConflictIgnore failed: %v", err)
	}
	if err := OnConflictIgnore(db, &fixture{ID: "a", Name: "second"}); err != nil {
		t.Fatalf("second OnConflictIgnore failed: %v", err)
	}
	if got := countFixtures(t, db); got != 1 {
		t.Fatalf("expected count=1, got %d", got)
	}
	var f fixture
	if err := db.First(&f, "id = ?", "a").Error; err != nil {
		t.Fatalf("First failed: %v", err)
	}
	if f.Name != "first" {
		t.Fatalf("expected Name=first (not updated), got %q", f.Name)
	}
}

func TestOnConflictIgnore_DifferentPKInserts(t *testing.T) {
	db := newTestDB(t)
	if err := OnConflictIgnore(db, &fixture{ID: "a", Name: "first"}); err != nil {
		t.Fatalf("insert a failed: %v", err)
	}
	if err := OnConflictIgnore(db, &fixture{ID: "b", Name: "bee"}); err != nil {
		t.Fatalf("insert b failed: %v", err)
	}
	if got := countFixtures(t, db); got != 2 {
		t.Fatalf("expected count=2, got %d", got)
	}
}
