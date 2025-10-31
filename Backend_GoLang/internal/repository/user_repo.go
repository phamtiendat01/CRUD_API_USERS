package repository

import (
	"errors"

	"crud_api_us/internal/models"
)

var (
	// Dùng cho case "không tìm thấy"
	ErrNotFound = errors.New("not found")
)

// Interface dùng chung cho mọi implementation (MySQL, memory, ...).
type UserRepository interface {
	List() ([]models.User, error)
	Get(id int) (models.User, error)
	Create(u *models.User) error
	Update(id int, in *models.User) (models.User, error)
	Delete(id int) (bool, error)
}
