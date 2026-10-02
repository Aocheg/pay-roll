package main

import (
	"fmt"
	"net/http"
	"payrool/internal/database"
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

	http.HandleFunc("/health", health)

	fmt.Println("PAY ROOL server running on http://localhost:8080")

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "PAY ROOL is running")
}
