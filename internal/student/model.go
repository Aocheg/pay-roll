package student

type Student struct {
	ID        int     `json:"id"`
	StudentNo string  `json:"student_no"`
	FullName  string  `json:"full_name"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Program   string  `json:"program"`
	DailyRate float64 `json:"daily_rate"`
	Active    bool    `json:"active"`
}
