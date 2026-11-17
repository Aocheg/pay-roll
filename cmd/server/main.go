package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"payrool/internal/attendance"
	"payrool/internal/database"
	"payrool/internal/payroll"
	"payrool/internal/student"
	"payrool/internal/user"
)

func main() {
	db, err := database.Open()
	if err != nil {
		fmt.Println("Database error:", err)
		return
	}
	defer db.Close()
	if err := database.CreateTables(db); err != nil {
		fmt.Println("Table creation error:", err)
		return
	}

	studentHandler := student.NewHandler(student.NewService(student.NewRepository(db)))
	attendanceHandler := attendance.NewHandler(attendance.NewService(attendance.NewRepository(db), student.NewRepository(db)))
	payrollHandler := payroll.NewHandler(payroll.NewService(payroll.NewRepository(db), attendance.NewRepository(db), student.NewRepository(db)))
	userHandler := user.NewHandler(user.NewService(user.NewRepository(db)))

	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/ready", ready(db))

	mux.HandleFunc("/api/auth/register", userHandler.RegisterUser)
	mux.HandleFunc("/api/auth/login", userHandler.Login)

	mux.HandleFunc("/api/users", user.RequireRole(user.RoleAdmin)(userHandler.ListUsers))
	mux.HandleFunc("/api/users/", user.RequireRole(user.RoleAdmin)(userHandler.GetUser))

	mux.HandleFunc("/api/students", user.RequireRole(user.RoleAdmin, user.RoleManager, user.RoleAttendanceOfficer, user.RoleStaff)(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			studentHandler.CreateStudent(w, r)
		case http.MethodGet:
			studentHandler.ListStudents(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	mux.HandleFunc("/api/students/", user.RequireRole(user.RoleAdmin, user.RoleManager, user.RoleAttendanceOfficer, user.RoleStaff)(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			studentHandler.GetStudent(w, r)
		case http.MethodPut:
			studentHandler.UpdateStudent(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/attendance", user.RequireRole(user.RoleAdmin, user.RoleManager, user.RoleAttendanceOfficer)(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			attendanceHandler.CreateAttendance(w, r)
		case http.MethodGet:
			attendanceHandler.ListAttendance(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	mux.HandleFunc("/api/attendance/student/", user.RequireRole(user.RoleAdmin, user.RoleManager, user.RoleAttendanceOfficer, user.RoleStaff)(attendanceHandler.ListAttendanceByStudent))
	mux.HandleFunc("/api/attendance/", user.RequireRole(user.RoleAdmin, user.RoleManager, user.RoleAttendanceOfficer)(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			attendanceHandler.GetAttendance(w, r)
		case http.MethodPut:
			attendanceHandler.UpdateAttendance(w, r)
		case http.MethodDelete:
			attendanceHandler.DeleteAttendance(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/payroll/calculate", user.RequireRole(user.RoleAdmin, user.RoleManager)(payrollHandler.CalculatePayroll))
	mux.HandleFunc("/api/payroll/student/", user.RequireRole(user.RoleAdmin, user.RoleManager, user.RoleStaff)(payrollHandler.ListPayrollByStudent))
	mux.HandleFunc("/api/payroll", user.RequireRole(user.RoleAdmin, user.RoleManager, user.RoleStaff)(payrollHandler.ListPayroll))
	mux.HandleFunc("/api/payroll/", user.RequireRole(user.RoleAdmin, user.RoleManager)(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/approve") {
			payrollHandler.ApprovePayroll(w, r)
			return
		}
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/pay") {
			payrollHandler.MarkPayrollPaid(w, r)
			return
		}
		if r.Method == http.MethodGet {
			payrollHandler.GetPayroll(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}))

	fmt.Println("PAY ROOL server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Server error:", err)
	}
}

func health(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "PAY ROOL is running") }
func ready(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, "database not ready", http.StatusServiceUnavailable)
			return
		}
		fmt.Fprintln(w, "PAY ROOL is ready")
	}
}
