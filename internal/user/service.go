package user

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterUser(u User, password string) (int64, error) {
	u.Username = strings.TrimSpace(u.Username)
	u.FullName = strings.TrimSpace(u.FullName)
	u.Email = strings.TrimSpace(u.Email)
	u.Role = strings.TrimSpace(u.Role)
	password = strings.TrimSpace(password)

	if u.Username == "" {
		return 0, errors.New("username is required")
	}
	if u.FullName == "" {
		return 0, errors.New("full name is required")
	}
	if password == "" {
		return 0, errors.New("password is required")
	}
	if u.Role == "" {
		u.Role = "staff"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	u.PasswordHash = string(hash)
	u.Active = true
	u.CreatedAt = now
	u.UpdatedAt = now

	return s.repo.Create(u)
}

func (s *Service) Login(username, password string) (User, string, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || password == "" {
		return User{}, "", errors.New("username and password are required")
	}

	u, err := s.repo.GetByUsername(username)
	if err != nil {
		return User{}, "", errors.New("invalid username or password")
	}
	if !u.Active {
		return User{}, "", errors.New("account is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return User{}, "", errors.New("invalid username or password")
	}

	token, err := GenerateToken(u)
	if err != nil {
		return User{}, "", err
	}

	return u, token, nil
}

func (s *Service) GetUser(id int) (User, error) {
	if id <= 0 {
		return User{}, errors.New("invalid user ID")
	}
	return s.repo.GetByID(id)
}

func (s *Service) ListUsers() ([]User, error) {
	return s.repo.List()
}

func (s *Service) UpdateUser(u User) error {
	if u.ID <= 0 {
		return errors.New("invalid user ID")
	}
	u.Username = strings.TrimSpace(u.Username)
	u.FullName = strings.TrimSpace(u.FullName)
	u.Email = strings.TrimSpace(u.Email)
	u.Role = strings.TrimSpace(u.Role)
	if u.Username == "" {
		return errors.New("username is required")
	}
	if u.FullName == "" {
		return errors.New("full name is required")
	}
	if u.Role == "" {
		u.Role = "staff"
	}
	u.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return s.repo.Update(u)
}
