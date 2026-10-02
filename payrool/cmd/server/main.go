package main

import (
	"fmt"
	"net/http"

	"payrool/internal/database"
	"payrool/internal/student"
)

func main() {
	db, err := database.Open()

	if err != nil {
		fmt.Println("Database error:", err)
		return
	}

	defer db.Close()

	err = database.CreateTables(db)

	if err != nil {
		fmt.Println("Table creation error:", err)
		return
	}

	studentRepo := student.NewRepository(db)
	studentService := student.NewService(studentRepo)
	studentHandler := student.NewHandler(studentService)

	http.HandleFunc("/health", health)
	http.HandleFunc("/api/students", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			studentHandler.CreateStudent(w, r)
			return
		}

		if r.Method == http.MethodGet {
			studentHandler.ListStudents(w, r)
			return
		}

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/api/students/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			studentHandler.GetStudent(w, r)
			return
		}

		if r.Method == http.MethodPut {
			studentHandler.UpdateStudent(w, r)
			return
		}

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	fmt.Println("PAY ROOL server running on http://localhost:8080")

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "PAY ROOL is running")
}
