package attendance

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service *Service
}

type createAttendanceResponse struct {
	ID int64 `json:"id"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateAttendance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request Attendance
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateAttendance(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createAttendanceResponse{ID: id})
}

func (h *Handler) ListAttendance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	attendanceList, err := h.service.ListAttendance()
	if err != nil {
		http.Error(w, "failed to retrieve attendance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendanceList)
}

func (h *Handler) GetAttendance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idText := strings.TrimPrefix(r.URL.Path, "/api/attendance/")
	id, err := strconv.Atoi(idText)
	if err != nil || id <= 0 {
		http.Error(w, "invalid attendance ID", http.StatusBadRequest)
		return
	}

	attendanceItem, err := h.service.GetAttendance(id)
	if err != nil {
		http.Error(w, "attendance not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendanceItem)
}

func (h *Handler) ListAttendanceByStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	studentIDText := strings.TrimPrefix(r.URL.Path, "/api/attendance/student/")
	studentID, err := strconv.Atoi(studentIDText)
	if err != nil || studentID <= 0 {
		http.Error(w, "invalid student ID", http.StatusBadRequest)
		return
	}

	attendanceList, err := h.service.ListAttendanceByStudent(studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendanceList)
}

func (h *Handler) UpdateAttendance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idText := strings.TrimPrefix(r.URL.Path, "/api/attendance/")
	id, err := strconv.Atoi(idText)
	if err != nil || id <= 0 {
		http.Error(w, "invalid attendance ID", http.StatusBadRequest)
		return
	}

	var request Attendance
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	request.ID = id

	if err := h.service.UpdateAttendance(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(request)
}

func (h *Handler) DeleteAttendance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idText := strings.TrimPrefix(r.URL.Path, "/api/attendance/")
	id, err := strconv.Atoi(idText)
	if err != nil || id <= 0 {
		http.Error(w, "invalid attendance ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAttendance(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
