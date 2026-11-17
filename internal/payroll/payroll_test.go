package payroll

import (
	"database/sql"
	"fmt"
	"testing"

	"payrool/internal/attendance"
	"payrool/internal/student"

	_ "modernc.org/sqlite"
)

func setupPayrollTestDB(t *testing.T) (*sql.DB, *student.Repository, *attendance.Repository) {
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
		CREATE TABLE payroll (
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
			UNIQUE(student_id, period_start, period_end)
		);
	`); err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	t.Cleanup(func() { db.Close() })
	return db, student.NewRepository(db), attendance.NewRepository(db)
}

func TestPayrollEligibleStudent(t *testing.T) {
	db, studentRepo, attendanceRepo := setupPayrollTestDB(t)
	repo := NewRepository(db)
	service := NewService(repo, attendanceRepo, studentRepo)

	studentID, err := studentRepo.Create(student.Student{
		StudentNo: "PAY001",
		FullName:  "Ada Doe",
		DailyRate: 7500,
		Active:    true,
	})
	if err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	for day := 1; day <= 20; day++ {
		date := fmt.Sprintf("2026-10-%02d", day)
		_, err := attendanceRepo.Create(attendance.Attendance{
			StudentID:      int(studentID),
			AttendanceDate: date,
			Status:         "present",
			CreatedAt:      "2026-10-01T00:00:00Z",
			UpdatedAt:      "2026-10-01T00:00:00Z",
		})
		if err != nil {
			t.Fatalf("failed to seed attendance: %v", err)
		}
	}

	payroll, err := service.CalculatePayroll(int(studentID), "2026-10-01", "2026-10-20")
	if err != nil {
		t.Fatalf("failed to calculate payroll: %v", err)
	}

	if !payroll.Eligible {
		t.Fatal("expected student to be eligible")
	}
	if payroll.ExpectedDays != 20 {
		t.Fatalf("expected 20 expected days, got %d", payroll.ExpectedDays)
	}
	if payroll.GrossAmount != 150000 {
		t.Fatalf("expected amount 150000, got %.2f", payroll.GrossAmount)
	}
}

func TestPayrollBelowThresholdStudent(t *testing.T) {
	db, studentRepo, attendanceRepo := setupPayrollTestDB(t)
	repo := NewRepository(db)
	service := NewService(repo, attendanceRepo, studentRepo)

	studentID, err := studentRepo.Create(student.Student{
		StudentNo: "PAY002",
		FullName:  "Grace Doe",
		DailyRate: 7500,
		Active:    true,
	})
	if err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	for day := 1; day <= 20; day++ {
		status := "present"
		if day == 9 {
			status = "absent"
		}
		date := fmt.Sprintf("2026-10-%02d", day)
		_, err := attendanceRepo.Create(attendance.Attendance{
			StudentID:      int(studentID),
			AttendanceDate: date,
			Status:         status,
			CreatedAt:      "2026-10-01T00:00:00Z",
			UpdatedAt:      "2026-10-01T00:00:00Z",
		})
		if err != nil {
			t.Fatalf("failed to seed attendance: %v", err)
		}
	}

	payroll, err := service.CalculatePayroll(int(studentID), "2026-10-01", "2026-10-20")
	if err != nil {
		t.Fatalf("failed to calculate payroll: %v", err)
	}

	if payroll.Eligible {
		t.Fatal("expected student to be ineligible")
	}
	if payroll.GrossAmount != 0 {
		t.Fatalf("expected zero payment, got %.2f", payroll.GrossAmount)
	}
}
