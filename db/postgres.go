package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	connStr := os.Getenv("DATABASE_URL") // e.g. postgres://user:pass@host:port/dbname?sslmode=disable
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL")
}
