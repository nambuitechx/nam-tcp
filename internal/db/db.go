package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func GetDBConnection() *sql.DB {
	db, err := sql.Open("sqlite", "app.db")
	if err != nil {
		log.Fatalf("db error: %v", err)
	}

	db.SetMaxOpenConns(1)
	return db
}

func Up(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
            id    TEXT PRIMARY KEY,
            email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
        );
		CREATE TABLE IF NOT EXISTS user_pats (
            id    TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
            hash_token TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			expires_at INTEGER NOT NULL
        );
	`)
	if err != nil {
		log.Printf("migration error: %v", err)
	}
}
