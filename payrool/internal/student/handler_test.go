package student

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestHandlerCreateStudent(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	requestBody := map[string]interface{}{
		"student_no": "STU001",
		"full_name":  "John Doe",
		"email":      "john@example.com",
		"phone":      "08012345678",
		"program":    "Software Engineering",
		"daily_rate": 5000,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("failed to encode request body: %v", err)
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/students",
		bytes.NewReader(body),
	)

	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	handler.CreateStudent(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			recorder.Code,
		)
	}

	var response createStudentResponse

	err = json.NewDecoder(recorder.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID <= 0 {
		t.Fatalf("expected valid student ID, got %d", response.ID)
	}

	created, err := repo.GetByID(int(response.ID))
	if err != nil {
		t.Fatalf("failed to retrieve created student: %v", err)
	}

	if created.FullName != "John Doe" {
		t.Fatalf(
			"expected student name %q, got %q",
			"John Doe",
			created.FullName,
		)
	}
}

func TestHandlerCreateStudentRejectsInvalidJSON(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/students",
		bytes.NewBufferString(`invalid json`),
	)

	recorder := httptest.NewRecorder()

	handler.CreateStudent(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}

func TestHandlerCreateStudentRejectsInvalidStudent(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	requestBody := map[string]interface{}{
		"student_no": "",
		"full_name":  "John Doe",
		"daily_rate": 5000,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("failed to encode request body: %v", err)
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/students",
		bytes.NewReader(body),
	)

	recorder := httptest.NewRecorder()

	handler.CreateStudent(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}

func TestHandlerCreateStudentRejectsWrongMethod(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/students",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.CreateStudent(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusMethodNotAllowed,
			recorder.Code,
		)
	}
}

func TestHandlerListStudents(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	_, err := service.CreateStudent(Student{
		StudentNo: "STU001",
		FullName:  "John Doe",
		Email:     "john@example.com",
		Phone:     "08012345678",
		Program:   "Software Engineering",
		DailyRate: 5000,
	})

	if err != nil {
		t.Fatalf("failed to create test student: %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/students", nil)
	response := httptest.NewRecorder()

	handler.ListStudents(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var students []Student

	err = json.NewDecoder(response.Body).Decode(&students)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(students) != 1 {
		t.Fatalf("expected 1 student, got %d", len(students))
	}

	if students[0].StudentNo != "STU001" {
		t.Fatalf("expected student number STU001, got %s", students[0].StudentNo)
	}

	if students[0].FullName != "John Doe" {
		t.Fatalf("expected full name John Doe, got %s", students[0].FullName)
	}
}

func TestHandlerGetStudent(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	id, err := service.CreateStudent(Student{
		StudentNo: "STU001",
		FullName:  "John Doe",
		Email:     "john@example.com",
		Phone:     "08012345678",
		Program:   "Software Engineering",
		DailyRate: 5000,
	})

	if err != nil {
		t.Fatalf("failed to create test student: %v", err)
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/students/"+strconv.FormatInt(id, 10),
		nil,
	)

	response := httptest.NewRecorder()

	handler.GetStudent(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var student Student

	err = json.NewDecoder(response.Body).Decode(&student)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if student.ID != int(id) {
		t.Fatalf("expected student ID %d, got %d", id, student.ID)
	}

	if student.StudentNo != "STU001" {
		t.Fatalf("expected student number STU001, got %s", student.StudentNo)
	}

	if student.FullName != "John Doe" {
		t.Fatalf("expected full name John Doe, got %s", student.FullName)
	}
}

func TestHandlerUpdateStudent(t *testing.T) {
	db := testDatabase(t)

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	id, err := service.CreateStudent(Student{
		StudentNo: "STU001",
		FullName:  "John Doe",
		Email:     "john@example.com",
		Phone:     "08012345678",
		Program:   "Software Engineering",
		DailyRate: 5000,
	})

	if err != nil {
		t.Fatalf("failed to create test student: %v", err)
	}

	requestBody := `{
		"student_no": "STU001",
		"full_name": "John Gabriel Doe",
		"email": "john.gabriel@example.com",
		"phone": "08098765432",
		"program": "Backend Engineering",
		"daily_rate": 7500
	}`

	request := httptest.NewRequest(
		http.MethodPut,
		"/api/students/"+strconv.FormatInt(id, 10),
		strings.NewReader(requestBody),
	)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	handler.UpdateStudent(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	updatedStudent, err := service.GetStudent(int(id))
	if err != nil {
		t.Fatalf("failed to retrieve updated student: %v", err)
	}

	if updatedStudent.FullName != "John Gabriel Doe" {
		t.Fatalf(
			"expected full name John Gabriel Doe, got %s",
			updatedStudent.FullName,
		)
	}

	if updatedStudent.Email != "john.gabriel@example.com" {
		t.Fatalf(
			"expected email john.gabriel@example.com, got %s",
			updatedStudent.Email,
		)
	}

	if updatedStudent.Program != "Backend Engineering" {
		t.Fatalf(
			"expected program Backend Engineering, got %s",
			updatedStudent.Program,
		)
	}

	if updatedStudent.DailyRate != 7500 {
		t.Fatalf(
			"expected daily rate 7500, got %.2f",
			updatedStudent.DailyRate,
		)
	}
}
