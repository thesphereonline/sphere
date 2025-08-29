package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
)

func ApplyMigrations(db *sql.DB, migrationsDir string) error {
	// naive runner: execute each .sql in lexicographic order
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".sql" {
			continue
		}
		path := filepath.Join(migrationsDir, e.Name())
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		log.Printf("Applying migration %s", e.Name())
		if _, err := db.Exec(string(b)); err != nil {
			return err
		}
	}
	return nil
}
