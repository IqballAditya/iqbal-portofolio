package database_test

import (
	"path/filepath"
	"strings"
	"testing"

	"iqbal-portfolio/internal/database"
)

func TestInitDBAndMigrate(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")

	db, err := database.InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	var journalMode string
	if err := db.QueryRow("PRAGMA journal_mode;").Scan(&journalMode); err != nil {
		t.Fatalf("failed to query journal_mode: %v", err)
	}
	if strings.ToLower(journalMode) != "wal" {
		t.Errorf("expected journal_mode wal, got %s", journalMode)
	}

	if err := database.Migrate(db); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	expectedTables := []string{
		"projects",
		"milestones",
		"skills",
		"shortlists",
		"inquiries",
	}

	for _, table := range expectedTables {
		var name string
		err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?;", table).Scan(&name)
		if err != nil {
			t.Errorf("expected table %s to exist, query error: %v", table, err)
		}
	}

	// Verify idempotency
	if err := database.Migrate(db); err != nil {
		t.Errorf("second Migrate call failed: %v", err)
	}
}
