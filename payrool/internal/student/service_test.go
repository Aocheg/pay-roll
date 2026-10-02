package student

import "testing"

func TestServiceCreateStudent(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)
	service := NewService(repo)

	student := Student{
		StudentNo: "STU001",
		FullName:  "John Doe",
		Email:     "john@example.com",
		Phone:     "08012345678",
		Program:   "Software Engineering",
		DailyRate: 5000,
		Active:    true,
	}

	id, err := service.CreateStudent(student)

	if err != nil {
		t.Fatalf("expected student creation to succeed: %v", err)
	}

	if id <= 0 {
		t.Fatalf("expected valid student ID, got %d", id)
	}
}

func TestServiceRejectsMissingStudentNumber(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)
	service := NewService(repo)

	student := Student{
		FullName:  "John Doe",
		DailyRate: 5000,
		Active:    true,
	}

	_, err := service.CreateStudent(student)

	if err == nil {
		t.Fatal("expected error for missing student number")
	}
}

func TestServiceRejectsMissingName(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)
	service := NewService(repo)

	student := Student{
		StudentNo: "STU002",
		DailyRate: 5000,
		Active:    true,
	}

	_, err := service.CreateStudent(student)

	if err == nil {
		t.Fatal("expected error for missing student name")
	}
}

func TestServiceRejectsNegativeDailyRate(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)
	service := NewService(repo)

	student := Student{
		StudentNo: "STU003",
		FullName:  "Jane Doe",
		DailyRate: -5000,
		Active:    true,
	}

	_, err := service.CreateStudent(student)

	if err == nil {
		t.Fatal("expected error for negative daily rate")
	}
}
