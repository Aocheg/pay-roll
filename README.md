# 🎓 Payrool — Student Payroll Management System

> A lightweight, robust RESTful API for managing student stipends and allowances based on automated attendance tracking. Built with Go, powered by an embedded SQLite database, and structured with a clean layered architecture (`handler` → `service` → `repository`).

---

## 📋 Table of Contents

- [Overview](#-overview)
  - [How It Works](#how-it-works)
  - [Eligibility Rule](#eligibility-rule)
- [Key Features](#-key-features)
- [Project Structure](#-project-structure)
- [Getting Started](#-getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation & Setup](#installation--setup)
- [Environment Variables](#-environment-variables)
- [Roles & Permissions](#-roles--permissions)
- [API Reference](#-api-reference)
  - [Health Checks](#health-checks)
  - [Authentication](#authentication)
  - [Student Management](#student-management)
  - [Attendance Tracking](#attendance-tracking)
  - [Payroll Management](#payroll-management)
  - [User Administration](#user-administration)
- [End-to-End Workflow Example](#-end-to-end-workflow-example)
- [Running Tests](#-running-tests)
- [Built With](#-built-with)

---

## 🌟 Overview

**Payrool** automates the computation and disbursement lifecycle of student stipends based on verifiable attendance records. It enforces strict Role-Based Access Control (RBAC), ensuring administrators, managers, attendance officers, and staff have dedicated, auditable operational boundaries.

### How It Works

```
┌────────────────────┐      ┌────────────────────┐      ┌────────────────────┐      ┌────────────────────┐
│ 1. Register Student│ ───► │ 2. Log Attendance  │ ───► │3. Calculate Payroll│ ───► │ 4. Approve & Pay   │
│   (Set Daily Rate) │      │(Present/Late/etc.) │      │ (Threshold: ≥ 70%) │      │(Disburse Stipend)  │
└────────────────────┘      └────────────────────┘      └────────────────────┘      └────────────────────┘
```

1. **Student Registration** — Onboard students into the system with their program details and an assigned daily stipend rate.
2. **Attendance Tracking** — Log daily attendance records with granular status flags (`present`, `absent`, `late`, `excused`).
3. **Automated Payroll Calculation** — Trigger payroll calculations for a specific date range. The system computes the attendance percentage, verifies stipend eligibility against business rules, and computes gross payable amounts.
4. **Approval & Payout** — Authorized managers and admins review, approve, and mark payroll records as disbursed.

### Eligibility Rule

> [!IMPORTANT]
> **Attendance Threshold:** A student must maintain an attendance rate of **at least 70%** within the billing period to be eligible for payment.
> 
> $$\text{Attendance Percentage} = \left(\frac{\text{Present Days}}{\text{Expected Days}}\right) \times 100$$
> 
> If `attendance_percentage < 70%`, `eligible` is set to `false` and gross payable amount evaluates to `0.00`.

---

## ✨ Key Features

- **🔐 Zero-Dependency JWT Authentication:** Secure token generation, signing, and verification implemented natively using Go's standard library (no external JWT bloat).
- **👥 Role-Based Access Control (RBAC):** Distinct permission tiers for `admin`, `manager`, `attendance_officer`, and `staff`.
- **🎓 Student Management:** Full lifecycle management (create, view, update) including configurable daily allowance rates.
- **📅 Daily Attendance Tracking:** Comprehensive attendance logging supporting multiple statuses (`present`, `absent`, `late`, `excused`).
- **⚙️ Automated Payroll Engine:** Dynamic calculation of attendance ratios, automatic eligibility enforcement, and gross amount computation (`daily_rate` × `payable_days`).
- **🔄 Multi-Stage Approval Workflow:** Governed state machine transitioning payroll records from `pending` → `approved` → `paid`.
- **🗄️ Embedded SQLite Database:** Zero configuration required; leverages `modernc.org/sqlite` (pure Go, zero CGO requirements).
- **🏗️ Clean Layered Architecture:** Strict separation of concerns across handlers, services, and repositories for testability and maintainability.
- **🧪 Unit & Integration Test Suites:** Comprehensive automated test coverage across domain logic, eligibility rules, and repository operations.

---

## 📁 Project Structure

```text
pay-roll/
├── cmd/
│   └── server/
│       └── main.go               # HTTP routing, server initialization, dependency wiring
├── internal/
│   ├── attendance/
│   │   ├── handler.go            # HTTP request/response handlers for attendance
│   │   ├── model.go              # Attendance domain models & validation
│   │   ├── repository.go         # SQLite queries for attendance persistence
│   │   ├── service.go            # Attendance business logic
│   │   └── attendance_test.go    # Attendance unit & integration tests
│   ├── database/
│   │   └── database.go           # Database connection pool & table migration schemas
│   ├── eligibility/
│   │   └── service.go            # Business rule engine for attendance eligibility (≥70%)
│   ├── payroll/
│   │   ├── handler.go            # HTTP handlers for calculation, approvals, and payouts
│   │   ├── model.go              # Payroll domain models & lifecycle statuses
│   │   ├── repository.go         # SQLite payroll records persistence
│   │   ├── service.go            # Payroll processing & computation engine
│   │   └── payroll_test.go       # Payroll calculation tests
│   ├── student/
│   │   ├── handler.go            # HTTP handlers for student records
│   │   ├── model.go              # Student domain models & rate definitions
│   │   ├── repository.go         # SQLite student queries
│   │   ├── service.go            # Student management business logic
│   │   └── student_test.go       # Student unit tests
│   └── user/
│       ├── auth.go               # Custom JWT generation, claim parsing, and validation
│       ├── handler.go            # Authentication & user administration handlers
│       ├── model.go              # User models, roles, and credential hashes
│       ├── repository.go         # SQLite user queries
│       └── service.go            # Password hashing (bcrypt) and auth services
├── go.mod                        # Go module dependencies
├── go.sum                        # Go module checksums
└── payrool.db                    # Embedded SQLite database (auto-created on first run)
```

---

## 🚀 Getting Started

### Prerequisites

- **Go:** Version `1.21` or higher installed ([Download Go](https://go.dev/dl/))
- **Git:** Version control client

### Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Aocheg/pay-roll.git
   cd pay-roll
   ```

2. **Download dependencies:**
   ```bash
   go mod download
   ```

3. **Start the API server:**
   ```bash
   go run ./cmd/server
   ```

4. **Verify startup:**
   The server initializes the SQLite database schema automatically and binds to port `8080`:
   ```text
   Database connection successful
   PAY ROOL server running on http://localhost:8080
   ```

---

## ⚙️ Environment Variables

Configure application settings via environment variables or a `.env` file:

| Variable | Description | Default | Required in Production |
| :--- | :--- | :--- | :---: |
| `PAYROOL_JWT_SECRET` | Secret key used to sign and verify JWT authentication tokens | `payrool-dev-secret` | **Yes** |

```bash
# Example: export custom secret before starting the server
export PAYROOL_JWT_SECRET="your-ultra-secure-random-secret-key"
```

> [!WARNING]
> The default secret `payrool-dev-secret` is intended **strictly for local development**. Always supply a cryptographically secure random string in staging and production environments.

---

## 🔒 Roles & Permissions

The application implements granular Role-Based Access Control (RBAC):

| Endpoint / Action | Admin | Manager | Attendance Officer | Staff |
| :--- | :---: | :---: | :---: | :---: |
| **Authentication** (`/api/auth/*`) | ✅ | ✅ | ✅ | ✅ |
| **User Administration** (`/api/users/*`) | ✅ | ❌ | ❌ | ❌ |
| **Student Management** (`/api/students/*`) | ✅ | ✅ | ✅ | ✅ |
| **Record / Modify Attendance** (`POST/PUT/DELETE /api/attendance/*`) | ✅ | ✅ | ✅ | ❌ |
| **Read Attendance Records** (`GET /api/attendance/*`) | ✅ | ✅ | ✅ | ✅ |
| **Calculate Payroll** (`POST /api/payroll/calculate`) | ✅ | ✅ | ❌ | ❌ |
| **View Payroll Records** (`GET /api/payroll/*`) | ✅ | ✅ | ❌ | ✅ |
| **Approve & Mark Paid** (`POST /api/payroll/{id}/approve`, `/pay`) | ✅ | ✅ | ❌ | ❌ |

---

## 📡 API Reference

All protected endpoints require an `Authorization` header containing a valid Bearer token:
```http
Authorization: Bearer <your_jwt_token>
```

### Health Checks

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | Basic application health check |
| `GET` | `/ready` | Database readiness check |

---

### Authentication

#### Register User
- **Method:** `POST`
- **Endpoint:** `/api/auth/register`
- **Access:** Public
- **Request Body:**
  ```json
  {
    "username": "john_admin",
    "password": "strongpassword",
    "full_name": "John Doe",
    "email": "john@example.com",
    "role": "admin"
  }
  ```
  *Allowed roles:* `admin`, `manager`, `attendance_officer`, `staff`
- **Response:** `201 Created`
  ```json
  {
    "id": 1
  }
  ```

#### Login
- **Method:** `POST`
- **Endpoint:** `/api/auth/login`
- **Access:** Public
- **Request Body:**
  ```json
  {
    "username": "john_admin",
    "password": "strongpassword"
  }
  ```
- **Response:** `200 OK`
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "username": "john_admin",
    "role": "admin",
    "full_name": "John Doe"
  }
  ```

---

### Student Management

#### Create Student
- **Method:** `POST`
- **Endpoint:** `/api/students`
- **Access:** `admin`, `manager`, `attendance_officer`, `staff`
- **Request Body:**
  ```json
  {
    "student_no": "STU001",
    "full_name": "Ada Lovelace",
    "email": "ada@example.com",
    "phone": "08012345678",
    "program": "Software Engineering",
    "daily_rate": 7500.00
  }
  ```

#### List All Students
- **Method:** `GET`
- **Endpoint:** `/api/students`
- **Access:** `admin`, `manager`, `attendance_officer`, `staff`

#### Get Student by ID
- **Method:** `GET`
- **Endpoint:** `/api/students/{id}`
- **Access:** `admin`, `manager`, `attendance_officer`, `staff`

#### Update Student
- **Method:** `PUT`
- **Endpoint:** `/api/students/{id}`
- **Access:** `admin`, `manager`, `attendance_officer`, `staff`

---

### Attendance Tracking

#### Record Daily Attendance
- **Method:** `POST`
- **Endpoint:** `/api/attendance`
- **Access:** `admin`, `manager`, `attendance_officer`
- **Request Body:**
  ```json
  {
    "student_id": 1,
    "attendance_date": "2026-10-01",
    "status": "present"
  }
  ```
  *Valid statuses:* `present`, `absent`, `late`, `excused`

#### List All Attendance Records
- **Method:** `GET`
- **Endpoint:** `/api/attendance`
- **Access:** `admin`, `manager`, `attendance_officer`

#### List Attendance by Student
- **Method:** `GET`
- **Endpoint:** `/api/attendance/student/{student_id}`
- **Access:** `admin`, `manager`, `attendance_officer`, `staff`

#### Get Attendance Record by ID
- **Method:** `GET`
- **Endpoint:** `/api/attendance/{id}`
- **Access:** `admin`, `manager`, `attendance_officer`

#### Update Attendance Record
- **Method:** `PUT`
- **Endpoint:** `/api/attendance/{id}`
- **Access:** `admin`, `manager`, `attendance_officer`

#### Delete Attendance Record
- **Method:** `DELETE`
- **Endpoint:** `/api/attendance/{id}`
- **Access:** `admin`, `manager`, `attendance_officer`

---

### Payroll Management

#### Calculate Payroll
- **Method:** `POST`
- **Endpoint:** `/api/payroll/calculate`
- **Access:** `admin`, `manager`
- **Request Body:**
  ```json
  {
    "student_id": 1,
    "period_start": "2026-10-01",
    "period_end": "2026-10-31"
  }
  ```
- **Response:** `200 OK`
  ```json
  {
    "id": 1
  }
  ```
  *Note:* The system aggregates all attendance logs within `period_start` and `period_end`, evaluates attendance against the 70% threshold, and computes:
  $$\text{gross\_amount} = \text{daily\_rate} \times \text{payable\_days}$$

#### List All Payroll Records
- **Method:** `GET`
- **Endpoint:** `/api/payroll`
- **Access:** `admin`, `manager`, `staff`

#### List Payroll by Student
- **Method:** `GET`
- **Endpoint:** `/api/payroll/student/{student_id}`
- **Access:** `admin`, `manager`, `staff`

#### Get Single Payroll Record
- **Method:** `GET`
- **Endpoint:** `/api/payroll/{id}`
- **Access:** `admin`, `manager`

#### Approve Payroll
- **Method:** `POST`
- **Endpoint:** `/api/payroll/{id}/approve`
- **Access:** `admin`, `manager`
- **Description:** Updates payroll status from `pending` to `approved`.

#### Mark Payroll as Paid
- **Method:** `POST`
- **Endpoint:** `/api/payroll/{id}/pay`
- **Access:** `admin`, `manager`
- **Description:** Updates payroll status from `approved` to `paid`.

---

### User Administration

#### List All Users
- **Method:** `GET`
- **Endpoint:** `/api/users`
- **Access:** `admin` only

#### Get User by ID
- **Method:** `GET`
- **Endpoint:** `/api/users/{id}`
- **Access:** `admin` only

---

## 💡 End-to-End Workflow Example

Here is a full lifecycle execution flow using `curl`:

```bash
# 1. Register an administrative account
curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123",
    "full_name": "System Administrator",
    "email": "admin@example.com",
    "role": "admin"
  }'

# 2. Login to obtain JWT Bearer token
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# 3. Create a student with daily stipend rate
curl -s -X POST http://localhost:8080/api/students \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "student_no": "STU001",
    "full_name": "Ada Lovelace",
    "email": "ada@example.com",
    "phone": "08012345678",
    "program": "Software Engineering",
    "daily_rate": 7500
  }'

# 4. Record attendance entries for student ID 1
curl -s -X POST http://localhost:8080/api/attendance \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 1,
    "attendance_date": "2026-10-01",
    "status": "present"
  }'

# 5. Calculate payroll for October 2026
curl -s -X POST http://localhost:8080/api/payroll/calculate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 1,
    "period_start": "2026-10-01",
    "period_end": "2026-10-31"
  }'

# 6. Approve payroll (ID: 1)
curl -s -X POST http://localhost:8080/api/payroll/1/approve \
  -H "Authorization: Bearer $TOKEN"

# 7. Mark payroll as paid
curl -s -X POST http://localhost:8080/api/payroll/1/pay \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🧪 Running Tests

Execute automated unit and integration tests across the codebase:

```bash
# Run all tests across all internal packages
go test ./...

# Run all tests with detailed verbose logging
go test ./... -v

# Run tests for specific domain modules
go test ./internal/payroll/... -v
go test ./internal/student/... -v
go test ./internal/attendance/... -v
```

---

## 🛠️ Built With

- **[Go](https://go.dev/)** — High-performance, statically typed compiled backend language.
- **[modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)** — Pure Go SQLite driver (CGO-free for zero-overhead cross-platform deployment).
- **[golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto)** — Robust bcrypt password hashing and verification.
