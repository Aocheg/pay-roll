package user

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Active       bool   `json:"active"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

const (
	RoleAdmin             = "admin"
	RoleManager           = "manager"
	RoleAttendanceOfficer = "attendance_officer"
	RoleStaff             = "staff"
)
