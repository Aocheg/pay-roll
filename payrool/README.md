🎓 Payrool — Student Payroll Management System
A RESTful API for managing student payroll based on attendance tracking. Built with Go, using a SQLite database and a clean layered architecture (handler → service → repository).

📋 Table of Contents
Overview
Features
Project Structure
Getting Started
Environment Variables
API Reference
Roles & Permissions
Running Tests
Overview
Payrool automates the process of calculating student stipends or allowances based on their attendance records. It supports role-based access control so that admins, managers, attendance officers, and staff each have appropriate access levels.

How it works:

Register students with a daily rate
Record daily attendance for each student
Calculate payroll for a given period — the system checks attendance percentage and eligibility automatically
Approve and mark payroll as paid
A student must have at least 70% attendance in a period to be eligible for payment.

Features
✅ JWT-based authentication (no third-party libraries)
✅ Role-based access control (Admin, Manager, Attendance Officer, Staff)
✅ Student management (create, read, update)
✅ Attendance tracking with status support (present, absent, late, excused)
✅ Automatic payroll calculation based on attendance percentage
✅ Payroll approval workflow (pending → approved → paid)
✅ SQLite database (zero setup, file-based)
✅ Clean layered architecture
✅ Unit and integration tests
Project Structure
payrool/
├── cmd/
│   └── server/
│       └── main.go               # Entry point, routes
├── internal/
│   ├── attendance/
│   │   ├── model.go              # Attendance struct
│   │   ├── handler.go            # HTTP handlers
│   │   ├── service.go            # Business logic
│   │   ├── repository.go         # Database queries
│   │   └── attendance_test.go    # Tests
│   ├── database/
│   │   └── database.go           # DB connection & table creation
│   ├── eligibility/
│   │   └── service.go            # Eligibility logic (≥70% attendance)
│   ├── payroll/
│   │   ├── model.go
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── payroll_test.go
│   ├── student/
│   │   ├── model.go
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── *_test.go
│   └── user/
│       ├── model.go
│       ├── auth.go               # JWT generation & validation
│       ├── handler.go
│       ├── service.go
│       └── repository.go
├── go.mod
├── go.sum
└── payrool.db                    # SQLite database (auto-created)
Getting Started
Prerequisites
Go 1.21 or higher
Installation
# Clone the repository
git clone https://github.com/yourusername/pay-roll.git
cd pay-roll/payrool

# Install dependencies
go mod download

# Run the server
go run ./cmd/server
The server will start at http://localhost:8080.

You should see:

Database connection successful
PAY ROOL server running on http://localhost:8080
Environment Variables
Create a .env file or export these variables before running:

Variable	Description	Default
PAYROOL_JWT_SECRET	Secret key for signing JWT tokens	payrool-dev-secret
export PAYROOL_JWT_SECRET=your-strong-secret-key-here
⚠️ Always set a strong PAYROOL_JWT_SECRET in production. The default value is for development only.

API Reference
Health Checks
GET /health     → "PAY ROOL is running"
GET /ready      → "PAY ROOL is ready" (checks DB connection)
Authentication
Register a User
POST /api/auth/register
Request body:

{
  "username": "john_admin",
  "password": "strongpassword",
  "full_name": "John Doe",
  "email": "john@example.com",
  "role": "admin"
}
Response:

{ "id": 1 }
Available roles: admin, manager, attendance_officer, staff

Login
POST /api/auth/login
Request body:

{
  "username": "john_admin",
  "password": "strongpassword"
}
Response:

{
  "token": "eyJ...",
  "username": "john_admin",
  "role": "admin",
  "full_name": "John Doe"
}
Use the token in all subsequent requests as a Bearer token: Authorization: Bearer <token>

Students
Create a Student
POST /api/students
Authorization: Bearer <token>
Request body:

{
  "student_no": "STU001",
  "full_name": "Ada Lovelace",
  "email": "ada@example.com",
  "phone": "08012345678",
  "program": "Software Engineering",
  "daily_rate": 7500
}
List All Students
GET /api/students
Authorization: Bearer <token>
Get a Student
GET /api/students/{id}
Authorization: Bearer <token>
Update a Student
PUT /api/students/{id}
Authorization: Bearer <token>
Attendance
Record Attendance
POST /api/attendance
Authorization: Bearer <token>
Request body:

{
  "student_id": 1,
  "attendance_date": "2026-10-01",
  "status": "present"
}
Valid statuses: present, absent, late, excused

List All Attendance Records
GET /api/attendance
Authorization: Bearer <token>
Get Attendance by Student
GET /api/attendance/student/{student_id}
Authorization: Bearer <token>
Get a Single Attendance Record
GET /api/attendance/{id}
Authorization: Bearer <token>
Update Attendance
PUT /api/attendance/{id}
Authorization: Bearer <token>
Delete Attendance
DELETE /api/attendance/{id}
Authorization: Bearer <token>
Payroll
Calculate Payroll
POST /api/payroll/calculate
Authorization: Bearer <token>
Request body:

{
  "student_id": 1,
  "period_start": "2026-10-01",
  "period_end": "2026-10-31"
}
Response:

{ "id": 1 }
The system automatically calculates attendance percentage, checks eligibility (≥70%), and computes gross amount based on daily rate × payable days.

List All Payroll Records
GET /api/payroll
Authorization: Bearer <token>
Get Payroll by Student
GET /api/payroll/student/{student_id}
Authorization: Bearer <token>
Get a Single Payroll Record
GET /api/payroll/{id}
Authorization: Bearer <token>
Approve Payroll
POST /api/payroll/{id}/approve
Authorization: Bearer <token>
Mark Payroll as Paid
POST /api/payroll/{id}/pay
Authorization: Bearer <token>
Users (Admin only)
List All Users
GET /api/users
Authorization: Bearer <token>   (admin only)
Get a User
GET /api/users/{id}
Authorization: Bearer <token>   (admin only)
Roles & Permissions
Endpoint	Admin	Manager	Attendance Officer	Staff
Register/Login	✅	✅	✅	✅
List/Get Users	✅	❌	❌	❌
Students (all)	✅	✅	✅	✅
Attendance (write)	✅	✅	✅	❌
Attendance (read)	✅	✅	✅	✅
Calculate Payroll	✅	✅	❌	❌
View Payroll	✅	✅	❌	✅
Approve/Pay Payroll	✅	✅	❌	❌
Running Tests
cd payrool

# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests for a specific package
go test ./internal/payroll/... -v
go test ./internal/student/... -v
go test ./internal/attendance/... -v
Example Workflow
# 1. Register an admin user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123","full_name":"Admin User","role":"admin"}'

# 2. Login and get token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 3. Create a student (use token from step 2)
curl -X POST http://localhost:8080/api/students \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"student_no":"STU001","full_name":"Ada Lovelace","daily_rate":7500}'

# 4. Record attendance
curl -X POST http://localhost:8080/api/attendance \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"student_id":1,"attendance_date":"2026-10-01","status":"present"}'

# 5. Calculate payroll
curl -X POST http://localhost:8080/api/payroll/calculate \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"student_id":1,"period_start":"2026-10-01","period_end":"2026-10-31"}'

# 6. Approve payroll
curl -X POST http://localhost:8080/api/payroll/1/approve \
  -H "Authorization: Bearer <token>"

# 7. Mark as paid
curl -X POST http://localhost:8080/api/payroll/1/pay \
  -H "Authorization: Bearer <token>"
Built With
Go — Core language
modernc.org/sqlite — Pure Go SQLite driver (no CGO)
golang.org/x/crypto — bcrypt for password hashing
