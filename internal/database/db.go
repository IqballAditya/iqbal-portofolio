package database

import (
	"database/sql"
	_ "embed"
	"strings"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, err
	}

	pragmas := "PRAGMA foreign_keys=ON;"
	if !strings.Contains(dataSourceName, "memory") {
		pragmas = "PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;"
	}

	if _, err := db.Exec(pragmas); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	_, err := db.Exec(schemaSQL)
	return err
}
