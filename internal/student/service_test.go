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
		Active:    false,
	}

	id, err := service.CreateStudent(student)

	if err != nil {
		t.Fatalf("expected student creation to succeed: %v", err)
	}

	if id <= 0 {
		t.Fatalf("expected valid student ID, got %d", id)
	}

	created, err := repo.GetByID(int(id))
	if err != nil {
		t.Fatalf("failed to retrieve created student: %v", err)
	}

	if !created.Active {
		t.Fatal("expected newly created student to be active")
	}
}

func TestServiceRejectsMissingStudentNumber(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)
	service := NewService(repo)

	student := Student{
		FullName:  "John Doe",
		DailyRate: 5000,
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
	}

	_, err := service.CreateStudent(student)

	if err == nil {
		t.Fatal("expected error for negative daily rate")
	}
}

func TestServiceGetStudent(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)
	service := NewService(repo)

	id, err := repo.Create(Student{
		StudentNo: "STU004",
		FullName:  "David Smith",
		DailyRate: 5500,
		Active:    true,
	})

	if err != nil {
		t.Fatalf("failed to create test student: %v", err)
	}

	found, err := service.GetStudent(int(id))
	if err != nil {
		t.Fatalf("failed to get student: %v", err)
	}

	if found.FullName != "David Smith" {
		t.Fatalf(
			"expected student name %q, got %q",
			"David Smith",
			found.FullName,
		)
	}
}

func TestServiceRejectsInvalidStudentID(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)
	service := NewService(repo)

	_, err := service.GetStudent(0)

	if err == nil {
		t.Fatal("expected error for invalid student ID")
	}
}

func TestServiceListStudents(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)
	service := NewService(repo)

	_, err := repo.Create(Student{
		StudentNo: "STU005",
		FullName:  "Student One",
		DailyRate: 5000,
		Active:    true,
	})

	if err != nil {
		t.Fatalf("failed to create first student: %v", err)
	}

	_, err = repo.Create(Student{
		StudentNo: "STU006",
		FullName:  "Student Two",
		DailyRate: 6000,
		Active:    true,
	})

	if err != nil {
		t.Fatalf("failed to create second student: %v", err)
	}

	students, err := service.ListStudents()
	if err != nil {
		t.Fatalf("failed to list students: %v", err)
	}

	if len(students) != 2 {
		t.Fatalf("expected 2 students, got %d", len(students))
	}
}

func TestServiceUpdateStudent(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)
	service := NewService(repo)

	id, err := repo.Create(Student{
		StudentNo: "STU007",
		FullName:  "Original Name",
		Email:     "original@example.com",
		Phone:     "08000000000",
		Program:   "Software Engineering",
		DailyRate: 5000,
		Active:    true,
	})

	if err != nil {
		t.Fatalf("failed to create test student: %v", err)
	}

	updatedStudent := Student{
		ID:        int(id),
		StudentNo: "STU007",
		FullName:  "Updated Name",
		Email:     "updated@example.com",
		Phone:     "08111111111",
		Program:   "Backend Engineering",
		DailyRate: 7000,
		Active:    true,
	}

	err = service.UpdateStudent(updatedStudent)
	if err != nil {
		t.Fatalf("failed to update student: %v", err)
	}

	found, err := service.GetStudent(int(id))
	if err != nil {
		t.Fatalf("failed to retrieve updated student: %v", err)
	}

	if found.FullName != "Updated Name" {
		t.Fatalf(
			"expected updated name %q, got %q",
			"Updated Name",
			found.FullName,
		)
	}

	if found.DailyRate != 7000 {
		t.Fatalf(
			"expected updated daily rate %.2f, got %.2f",
			7000.0,
			found.DailyRate,
		)
	}
}

func TestServiceDeactivateStudent(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)
	service := NewService(repo)

	id, err := repo.Create(Student{
		StudentNo: "STU008",
		FullName:  "Inactive Student",
		DailyRate: 5000,
		Active:    true,
	})

	if err != nil {
		t.Fatalf("failed to create test student: %v", err)
	}

	err = service.DeactivateStudent(int(id))
	if err != nil {
		t.Fatalf("failed to deactivate student: %v", err)
	}

	found, err := service.GetStudent(int(id))
	if err != nil {
		t.Fatalf("failed to retrieve student: %v", err)
	}

	if found.Active {
		t.Fatal("expected student to be inactive")
	}
}
