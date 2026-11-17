package attendance

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"payrool/internal/student"
)

var validStatuses = map[string]bool{
	"present": true,
	"absent":  true,
	"late":    true,
	"excused": true,
}

type Service struct {
	repo        *Repository
	studentRepo *student.Repository
}

func NewService(repo *Repository, studentRepo *student.Repository) *Service {
	return &Service{
		repo:        repo,
		studentRepo: studentRepo,
	}
}

func (s *Service) CreateAttendance(a Attendance) (int64, error) {
	if a.StudentID <= 0 {
		return 0, errors.New("student ID is required")
	}

	if _, err := s.studentRepo.GetByID(a.StudentID); err != nil {
		return 0, fmt.Errorf("student not found: %w", err)
	}

	if err := validateDate(a.AttendanceDate); err != nil {
		return 0, err
	}

	a.Status = strings.ToLower(strings.TrimSpace(a.Status))
	if a.Status == "" {
		return 0, errors.New("attendance status is required")
	}
	if !validStatuses[a.Status] {
		return 0, errors.New("status must be one of: present, absent, late, excused")
	}

	now := time.Now().UTC().Format(time.RFC3339)
	a.CreatedAt = now
	a.UpdatedAt = now

	if _, err := s.repo.GetByStudentAndDate(a.StudentID, a.AttendanceDate); err == nil {
		return 0, errors.New("attendance for this student and date already exists")
	}

	return s.repo.Create(a)
}

func (s *Service) GetAttendance(id int) (Attendance, error) {
	if id <= 0 {
		return Attendance{}, errors.New("invalid attendance ID")
	}
	return s.repo.GetByID(id)
}

func (s *Service) ListAttendance() ([]Attendance, error) {
	return s.repo.List()
}

func (s *Service) ListAttendanceByStudent(studentID int) ([]Attendance, error) {
	if studentID <= 0 {
		return nil, errors.New("invalid student ID")
	}
	return s.repo.ListByStudent(studentID)
}

func (s *Service) UpdateAttendance(a Attendance) error {
	if a.ID <= 0 {
		return errors.New("invalid attendance ID")
	}

	if _, err := s.studentRepo.GetByID(a.StudentID); err != nil {
		return fmt.Errorf("student not found: %w", err)
	}

	if err := validateDate(a.AttendanceDate); err != nil {
		return err
	}

	a.Status = strings.ToLower(strings.TrimSpace(a.Status))
	if a.Status == "" {
		return errors.New("attendance status is required")
	}
	if !validStatuses[a.Status] {
		return errors.New("status must be one of: present, absent, late, excused")
	}

	a.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return s.repo.Update(a)
}

func (s *Service) DeleteAttendance(id int) error {
	if id <= 0 {
		return errors.New("invalid attendance ID")
	}
	return s.repo.Delete(id)
}

func validateDate(dateText string) error {
	dateText = strings.TrimSpace(dateText)
	if dateText == "" {
		return errors.New("attendance date is required")
	}

	if _, err := time.Parse("2006-01-02", dateText); err != nil {
		return errors.New("attendance date must be in YYYY-MM-DD format")
	}

	return nil
}
