package student

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type createStudentRequest struct {
	StudentNo string  `json:"student_no"`
	FullName  string  `json:"full_name"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Program   string  `json:"program"`
	DailyRate float64 `json:"daily_rate"`
}

type createStudentResponse struct {
	ID int64 `json:"id"`
}

type updateStudentRequest struct {
	StudentNo string  `json:"student_no"`
	FullName  string  `json:"full_name"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Program   string  `json:"program"`
	DailyRate float64 `json:"daily_rate"`
}

func (h *Handler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request createStudentRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	student := Student{
		StudentNo: request.StudentNo,
		FullName:  request.FullName,
		Email:     request.Email,
		Phone:     request.Phone,
		Program:   request.Program,
		DailyRate: request.DailyRate,
	}

	id, err := h.service.CreateStudent(student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := createStudentResponse{
		ID: id,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ListStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	students, err := h.service.ListStudents()
	if err != nil {
		http.Error(w, "failed to retrieve students", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(students)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idText := strings.TrimPrefix(r.URL.Path, "/api/students/")

	id, err := strconv.Atoi(idText)
	if err != nil || id <= 0 {
		http.Error(w, "invalid student ID", http.StatusBadRequest)
		return
	}

	student, err := h.service.GetStudent(id)
	if err != nil {
		http.Error(w, "student not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(student)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idText := strings.TrimPrefix(r.URL.Path, "/api/students/")

	id, err := strconv.Atoi(idText)
	if err != nil || id <= 0 {
		http.Error(w, "invalid student ID", http.StatusBadRequest)
		return
	}

	var request updateStudentRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	student := Student{
		ID:        id,
		StudentNo: request.StudentNo,
		FullName:  request.FullName,
		Email:     request.Email,
		Phone:     request.Phone,
		Program:   request.Program,
		DailyRate: request.DailyRate,
	}

	err = h.service.UpdateStudent(student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(student)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
