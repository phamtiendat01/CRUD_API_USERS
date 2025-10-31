package services

import (
	"errors"
	"strings"
	"time"

	"crud_api_us/internal/models"
	"crud_api_us/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrDuplicate = errors.New("duplicate") // email/username trùng
	ErrBadInput  = errors.New("bad_input") // dữ liệu không hợp lệ
)

type UserService struct{ repo repository.UserRepository }

func NewUserService(r repository.UserRepository) *UserService { return &UserService{repo: r} }

// ====== Helpers ======
func hashPassword(pw string) (string, error) {
	if strings.TrimSpace(pw) == "" {
		return "", ErrBadInput
	}
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}
func parseDOB(s string) (*time.Time, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, ErrBadInput
	}
	return &t, nil
}

// ====== Service API ======
type CreateParams struct {
	Username, Email, Password, FullName, Phone, Gender, DOB,
	AvatarURL, Street, City, State, Country, PostalCode,
	Role, Status string
}

type UpdateParams struct {
	Username, Email, Password, FullName, Phone, Gender, DOB,
	AvatarURL, Street, City, State, Country, PostalCode,
	Role, Status string
}

func (s *UserService) List() ([]models.User, error)    { return s.repo.List() }
func (s *UserService) Get(id int) (models.User, error) { return s.repo.Get(id) }
func (s *UserService) Delete(id int) (bool, error)     { return s.repo.Delete(id) }

func (s *UserService) Create(p CreateParams) (models.User, error) {
	dob, err := parseDOB(p.DOB)
	if err != nil {
		return models.User{}, err
	}

	hash, err := hashPassword(p.Password)
	if err != nil {
		return models.User{}, err
	}

	u := models.User{
		Username:     strings.TrimSpace(p.Username),
		Email:        strings.TrimSpace(p.Email),
		PasswordHash: hash,
		FullName:     p.FullName,
		Phone:        p.Phone,
		Gender:       p.Gender,
		DateOfBirth:  dob,
		AvatarURL:    p.AvatarURL,
		Street:       p.Street, City: p.City, State: p.State, Country: p.Country, PostalCode: p.PostalCode,
		Role:   defaultIfEmpty(p.Role, "user"),
		Status: defaultIfEmpty(p.Status, "active"),
	}
	if err := s.repo.Create(&u); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "Duplicate entry") {
			return models.User{}, ErrDuplicate
		}
		return models.User{}, err
	}
	return u, nil
}

func (s *UserService) Update(id int, p UpdateParams) (models.User, error) {
	dob, err := parseDOB(p.DOB)
	if err != nil {
		return models.User{}, err
	}

	u := models.User{
		Username: strings.TrimSpace(p.Username),
		Email:    strings.TrimSpace(p.Email),
		FullName: p.FullName, Phone: p.Phone, Gender: p.Gender, DateOfBirth: dob,
		AvatarURL: p.AvatarURL, Street: p.Street, City: p.City, State: p.State, Country: p.Country, PostalCode: p.PostalCode,
		Role: p.Role, Status: p.Status,
	}
	if strings.TrimSpace(p.Password) != "" {
		if u.PasswordHash, err = hashPassword(p.Password); err != nil {
			return models.User{}, err
		}
	}
	out, err := s.repo.Update(id, &u)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "Duplicate entry") {
			return models.User{}, ErrDuplicate
		}
		return models.User{}, err
	}
	return out, nil
}

func defaultIfEmpty(s, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}
