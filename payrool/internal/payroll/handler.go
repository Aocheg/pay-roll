package payroll

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"payrool/internal/attendance"
	"payrool/internal/eligibility"
	"payrool/internal/student"
)

type Service struct {
	repo           *Repository
	attendanceRepo *attendance.Repository
	studentRepo    *student.Repository
}

func NewService(repo *Repository, attendanceRepo *attendance.Repository, studentRepo *student.Repository) *Service {
	return &Service{repo: repo, attendanceRepo: attendanceRepo, studentRepo: studentRepo}
}

func (s *Service) CalculatePayroll(studentID int, periodStart, periodEnd string) (Payroll, error) {
	if studentID <= 0 {
		return Payroll{}, errors.New("invalid student ID")
	}
	if err := validateDateRange(periodStart, periodEnd); err != nil {
		return Payroll{}, err
	}

	student, err := s.studentRepo.GetByID(studentID)
	if err != nil {
		return Payroll{}, fmt.Errorf("student not found: %w", err)
	}

	if _, err := s.repo.GetByStudentAndPeriod(studentID, periodStart, periodEnd); err == nil {
		return Payroll{}, errors.New("payroll already exists for this student and period")
	}

	attendanceList, err := s.attendanceRepo.ListByStudent(studentID)
	if err != nil {
		return Payroll{}, fmt.Errorf("failed to retrieve attendance: %w", err)
	}

	expectedDays := countExpectedDays(periodStart, periodEnd)
	presentDays := 0
	for _, record := range attendanceList {
		if isWithinPeriod(record.AttendanceDate, periodStart, periodEnd) {
			if strings.EqualFold(record.Status, "present") || strings.EqualFold(record.Status, "late") || strings.EqualFold(record.Status, "excused") {
				presentDays++
			}
		}
	}

	attendancePercentage := 0.0
	if expectedDays > 0 {
		attendancePercentage = (float64(presentDays) / float64(expectedDays)) * 100
	}

	eligible := eligibility.IsEligible(attendancePercentage)
	payableDays := 0
	grossAmount := 0.0
	if eligible {
		payableDays = presentDays
		grossAmount = float64(payableDays) * student.DailyRate
	}

	now := time.Now().UTC().Format(time.RFC3339)
	payroll := Payroll{
		StudentID:            studentID,
		PeriodStart:          periodStart,
		PeriodEnd:            periodEnd,
		ExpectedDays:         expectedDays,
		PresentDays:          presentDays,
		AttendancePercentage: attendancePercentage,
		Eligible:             eligible,
		DailyRate:            student.DailyRate,
		PayableDays:          payableDays,
		GrossAmount:          grossAmount,
		PaymentStatus:        "pending",
		ProcessedAt:          "",
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	id, err := s.repo.Create(payroll)
	if err != nil {
		return Payroll{}, err
	}
	payroll.ID = int(id)
	return payroll, nil
}

func (s *Service) GetPayroll(id int) (Payroll, error) {
	if id <= 0 {
		return Payroll{}, errors.New("invalid payroll ID")
	}
	return s.repo.GetByID(id)
}

func (s *Service) ListPayroll() ([]Payroll, error) {
	return s.repo.List()
}

func (s *Service) ListPayrollByStudent(studentID int) ([]Payroll, error) {
	if studentID <= 0 {
		return nil, errors.New("invalid student ID")
	}
	return s.repo.ListByStudent(studentID)
}

func (s *Service) UpdateStatus(id int, status string) error {
	if id <= 0 {
		return errors.New("invalid payroll ID")
	}
	status = strings.ToLower(strings.TrimSpace(status))
	if status == "" {
		return errors.New("payment status is required")
	}

	allowed := map[string]bool{
		"pending":  true,
		"approved": true,
		"paid":     true,
		"rejected": true,
	}
	if !allowed[status] {
		return errors.New("status must be one of: pending, approved, paid, rejected")
	}
	return s.repo.UpdateStatus(id, status)
}

func validateDateRange(periodStart, periodEnd string) error {
	if strings.TrimSpace(periodStart) == "" || strings.TrimSpace(periodEnd) == "" {
		return errors.New("period start and end dates are required")
	}

	start, err := time.Parse("2006-01-02", periodStart)
	if err != nil {
		return errors.New("period start must be in YYYY-MM-DD format")
	}
	end, err := time.Parse("2006-01-02", periodEnd)
	if err != nil {
		return errors.New("period end must be in YYYY-MM-DD format")
	}
	if end.Before(start) {
		return errors.New("period end must be on or after period start")
	}
	return nil
}

func countExpectedDays(periodStart, periodEnd string) int {
	start, err := time.Parse("2006-01-02", periodStart)
	if err != nil {
		return 0
	}
	end, err := time.Parse("2006-01-02", periodEnd)
	if err != nil {
		return 0
	}

	return int(end.Sub(start).Hours()/24) + 1
}

func isWithinPeriod(dateText, periodStart, periodEnd string) bool {
	date, err := time.Parse("2006-01-02", dateText)
	if err != nil {
		return false
	}
	start, err := time.Parse("2006-01-02", periodStart)
	if err != nil {
		return false
	}
	end, err := time.Parse("2006-01-02", periodEnd)
	if err != nil {
		return false
	}
	return !date.Before(start) && !date.After(end)
}
