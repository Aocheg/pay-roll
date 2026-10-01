package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "payrool.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	fmt.Println("Database connection successful")

	return db, nil
}

func CreateTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		student_no TEXT NOT NULL UNIQUE,
		full_name TEXT NOT NULL,
		email TEXT,
		phone TEXT,
		program TEXT,
		daily_rate REAL NOT NULL,
		active INTEGER NOT NULL DEFAULT 1
	);
	`

	_, err := db.Exec(query)

	return err
}