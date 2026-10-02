package student

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func testDatabase(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	query := `
	CREATE TABLE students (
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

	_, err = db.Exec(query)
	if err != nil {
		t.Fatalf("failed to create students table: %v", err)
	}

	return db
}

func TestRepositoryCreate(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)

	student := Student{
		StudentNo: "STU001",
		FullName:  "John Doe",
		Email:     "john@example.com",
		Phone:     "08012345678",
		Program:   "Software Engineering",
		DailyRate: 5000,
		Active:    true,
	}

	id, err := repo.Create(student)
	if err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	if id <= 0 {
		t.Fatalf("expected a valid student ID, got %d", id)
	}

	var name string

	err = db.QueryRow(
		"SELECT full_name FROM students WHERE id = ?",
		id,
	).Scan(&name)

	if err != nil {
		t.Fatalf("failed to retrieve created student: %v", err)
	}

	if name != student.FullName {
		t.Fatalf(
			"expected student name %q, got %q",
			student.FullName,
			name,
		)
	}
}

func TestRepositoryGetByID(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)

	student := Student{
		StudentNo: "STU002",
		FullName:  "Jane Doe",
		Email:     "jane@example.com",
		Phone:     "08098765432",
		Program:   "Data Science",
		DailyRate: 6000,
		Active:    true,
	}

	id, err := repo.Create(student)
	if err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	found, err := repo.GetByID(int(id))
	if err != nil {
		t.Fatalf("failed to get student: %v", err)
	}

	if found.ID != int(id) {
		t.Fatalf("expected ID %d, got %d", id, found.ID)
	}

	if found.StudentNo != student.StudentNo {
		t.Fatalf(
			"expected student number %q, got %q",
			student.StudentNo,
			found.StudentNo,
		)
	}

	if found.FullName != student.FullName {
		t.Fatalf(
			"expected student name %q, got %q",
			student.FullName,
			found.FullName,
		)
	}

	if found.DailyRate != student.DailyRate {
		t.Fatalf(
			"expected daily rate %.2f, got %.2f",
			student.DailyRate,
			found.DailyRate,
		)
	}
}

func TestRepositoryList(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)

	students := []Student{
		{
			StudentNo: "STU001",
			FullName:  "John Doe",
			Email:     "john@example.com",
			Phone:     "08012345678",
			Program:   "Software Engineering",
			DailyRate: 5000,
			Active:    true,
		},
		{
			StudentNo: "STU002",
			FullName:  "Jane Doe",
			Email:     "jane@example.com",
			Phone:     "08098765432",
			Program:   "Data Science",
			DailyRate: 6000,
			Active:    true,
		},
		{
			StudentNo: "STU003",
			FullName:  "David Smith",
			Email:     "david@example.com",
			Phone:     "08123456789",
			Program:   "Cybersecurity",
			DailyRate: 5500,
			Active:    true,
		},
	}

	for _, student := range students {
		_, err := repo.Create(student)

		if err != nil {
			t.Fatalf("failed to create test student: %v", err)
		}
	}

	found, err := repo.List()

	if err != nil {
		t.Fatalf("failed to list students: %v", err)
	}

	if len(found) != len(students) {
		t.Fatalf(
			"expected %d students, got %d",
			len(students),
			len(found),
		)
	}

	if found[0].StudentNo != "STU001" {
		t.Fatalf(
			"expected first student STU001, got %s",
			found[0].StudentNo,
		)
	}

	if found[1].StudentNo != "STU002" {
		t.Fatalf(
			"expected second student STU002, got %s",
			found[1].StudentNo,
		)
	}

	if found[2].StudentNo != "STU003" {
		t.Fatalf(
			"expected third student STU003, got %s",
			found[2].StudentNo,
		)
	}
}
