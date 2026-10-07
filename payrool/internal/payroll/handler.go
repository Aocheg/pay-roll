package payroll

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service *Service
}

type calculatePayrollRequest struct {
	StudentID   int    `json:"student_id"`
	PeriodStart string `json:"period_start"`
	PeriodEnd   string `json:"period_end"`
}

type calculatePayrollResponse struct {
	ID int64 `json:"id"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CalculatePayroll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request calculatePayrollRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	payroll, err := h.service.CalculatePayroll(request.StudentID, request.PeriodStart, request.PeriodEnd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(calculatePayrollResponse{ID: int64(payroll.ID)})
}

func (h *Handler) GetPayroll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idText := strings.TrimPrefix(r.URL.Path, "/api/payroll/")
	id, err := strconv.Atoi(idText)
	if err != nil || id <= 0 {
		http.Error(w, "invalid payroll ID", http.StatusBadRequest)
		return
	}

	payroll, err := h.service.GetPayroll(id)
	if err != nil {
		http.Error(w, "payroll not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payroll)
}

func (h *Handler) ListPayroll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	payrolls, err := h.service.ListPayroll()
	if err != nil {
		http.Error(w, "failed to retrieve payroll", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payrolls)
}

func (h *Handler) ListPayrollByStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	studentIDText := strings.TrimPrefix(r.URL.Path, "/api/payroll/student/")
	studentID, err := strconv.Atoi(studentIDText)
	if err != nil || studentID <= 0 {
		http.Error(w, "invalid student ID", http.StatusBadRequest)
		return
	}

	payrolls, err := h.service.ListPayrollByStudent(studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payrolls)
}

func (h *Handler) ApprovePayroll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idText := strings.TrimPrefix(r.URL.Path, "/api/payroll/")
	idText = strings.TrimSuffix(idText, "/approve")
	id, err := strconv.Atoi(idText)
	if err != nil || id <= 0 {
		http.Error(w, "invalid payroll ID", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateStatus(id, "approved"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "approved"})
}

func (h *Handler) MarkPayrollPaid(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idText := strings.TrimPrefix(r.URL.Path, "/api/payroll/")
	idText = strings.TrimSuffix(idText, "/pay")
	id, err := strconv.Atoi(idText)
	if err != nil || id <= 0 {
		http.Error(w, "invalid payroll ID", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateStatus(id, "paid"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "paid"})
}
