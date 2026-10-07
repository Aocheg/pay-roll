package attendance

import (
	"database/sql"
	"testing"

	"payrool/internal/student"

	_ "modernc.org/sqlite"
)

func setupAttendanceTestDB(t *testing.T) (*sql.DB, *student.Repository, *Repository) {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if _, err = db.Exec(`
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
		CREATE TABLE attendance (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			student_id INTEGER NOT NULL,
			attendance_date TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			UNIQUE(student_id, attendance_date)
		);
	`); err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	t.Cleanup(func() { db.Close() })
	return db, student.NewRepository(db), NewRepository(db)
}

func TestCreateAttendance(t *testing.T) {
	_, studentRepo, attendanceRepo := setupAttendanceTestDB(t)
	service := NewService(attendanceRepo, studentRepo)

	studentID, err := studentRepo.Create(student.Student{
		StudentNo: "ATT001",
		FullName:  "Alice Doe",
		DailyRate: 5000,
		Active:    true,
	})
	if err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	id, err := service.CreateAttendance(Attendance{
		StudentID:      int(studentID),
		AttendanceDate: "2026-10-06",
		Status:         "present",
	})
	if err != nil {
		t.Fatalf("failed to create attendance: %v", err)
	}
	if id <= 0 {
		t.Fatalf("expected valid attendance ID, got %d", id)
	}

	item, err := attendanceRepo.GetByID(int(id))
	if err != nil {
		t.Fatalf("failed to fetch attendance: %v", err)
	}
	if item.Status != "present" {
		t.Fatalf("expected present status, got %q", item.Status)
	}
}

func TestCreateAttendanceDuplicateRejects(t *testing.T) {
	_, studentRepo, attendanceRepo := setupAttendanceTestDB(t)
	service := NewService(attendanceRepo, studentRepo)

	studentID, err := studentRepo.Create(student.Student{
		StudentNo: "ATT002",
		FullName:  "Bob Doe",
		DailyRate: 5000,
		Active:    true,
	})
	if err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	_, err = service.CreateAttendance(Attendance{
		StudentID:      int(studentID),
		AttendanceDate: "2026-10-06",
		Status:         "present",
	})
	if err != nil {
		t.Fatalf("failed to insert first attendance: %v", err)
	}

	_, err = service.CreateAttendance(Attendance{
		StudentID:      int(studentID),
		AttendanceDate: "2026-10-06",
		Status:         "absent",
	})
	if err == nil {
		t.Fatal("expected duplicate attendance to be rejected")
	}
}
