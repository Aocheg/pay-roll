package attendance

type Attendance struct {
	ID             int    `json:"id"`
	StudentID      int    `json:"student_id"`
	AttendanceDate string `json:"attendance_date"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
