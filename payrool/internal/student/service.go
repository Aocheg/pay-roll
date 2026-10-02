package student

import (
	"errors"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateStudent(student Student) (int64, error) {
	student.StudentNo = strings.TrimSpace(student.StudentNo)
	student.FullName = strings.TrimSpace(student.FullName)
	student.Email = strings.TrimSpace(student.Email)
	student.Phone = strings.TrimSpace(student.Phone)
	student.Program = strings.TrimSpace(student.Program)

	if student.StudentNo == "" {
		return 0, errors.New("student number is required")
	}

	if student.FullName == "" {
		return 0, errors.New("student name is required")
	}

	if student.DailyRate < 0 {
		return 0, errors.New("daily rate cannot be negative")
	}

	if !student.Active {
		student.Active = true
	}

	return s.repo.Create(student)
}
