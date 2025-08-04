package db

import (
	"context"
	"embed"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations(conn *pgx.Conn) error {
	files, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read migration dir: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		content, err := migrationFiles.ReadFile("migrations/" + file.Name())
		if err != nil {
			return fmt.Errorf("read migration file: %w", err)
		}

		sqlStatements := splitSQLStatements(string(content))

		log.Printf("Applying migration %s with %d statements...", file.Name(), len(sqlStatements))
		for i, stmt := range sqlStatements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			_, err := conn.Exec(context.Background(), stmt)
			if err != nil {
				return fmt.Errorf("exec statement #%d in migration %s: %w", i+1, file.Name(), err)
			}
		}
		log.Printf("Migration %s applied successfully", file.Name())
	}

	return nil
}

// splitSQLStatements naively splits by semicolon;
// for production, consider a proper SQL parser to handle semicolons inside strings/comments.
func splitSQLStatements(sql string) []string {
	// Simple split by semicolon
	return strings.Split(sql, ";")
}
