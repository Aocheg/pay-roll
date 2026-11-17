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

	CREATE TABLE IF NOT EXISTS attendance (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		student_id INTEGER NOT NULL,
		attendance_date TEXT NOT NULL,
		status TEXT NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		UNIQUE(student_id, attendance_date),
		FOREIGN KEY(student_id) REFERENCES students(id)
	);

	CREATE TABLE IF NOT EXISTS payroll (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		student_id INTEGER NOT NULL,
		period_start TEXT NOT NULL,
		period_end TEXT NOT NULL,
		expected_days INTEGER NOT NULL,
		present_days INTEGER NOT NULL,
		attendance_percentage REAL NOT NULL,
		eligible INTEGER NOT NULL,
		daily_rate REAL NOT NULL,
		payable_days INTEGER NOT NULL,
		gross_amount REAL NOT NULL,
		payment_status TEXT NOT NULL DEFAULT 'pending',
		processed_at TEXT,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		UNIQUE(student_id, period_start, period_end),
		FOREIGN KEY(student_id) REFERENCES students(id)
	);

	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		full_name TEXT NOT NULL,
		email TEXT,
		role TEXT NOT NULL DEFAULT 'staff',
		active INTEGER NOT NULL DEFAULT 1,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);
	`

	_, err := db.Exec(query)
	return err
}
