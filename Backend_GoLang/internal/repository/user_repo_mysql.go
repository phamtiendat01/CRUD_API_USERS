package repository

import (
	"errors"

	"crud_api_us/internal/models"

	"gorm.io/gorm"
)

type mysqlUserRepo struct{ db *gorm.DB }

func NewMySQLUserRepo(db *gorm.DB) UserRepository { return &mysqlUserRepo{db: db} }

func (r *mysqlUserRepo) List() ([]models.User, error) {
	var users []models.User
	return users, r.db.Order("id").Find(&users).Error
}

func (r *mysqlUserRepo) Get(id int) (models.User, error) {
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}
	return u, nil
}

func (r *mysqlUserRepo) Create(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *mysqlUserRepo) Update(id int, in *models.User) (models.User, error) {
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}
	u.Username = in.Username
	u.Email = in.Email
	u.FullName = in.FullName
	u.Phone = in.Phone
	u.Gender = in.Gender
	u.DateOfBirth = in.DateOfBirth
	u.AvatarURL = in.AvatarURL
	u.Street = in.Street
	u.City = in.City
	u.State = in.State
	u.Country = in.Country
	u.PostalCode = in.PostalCode
	u.Role = in.Role
	u.Status = in.Status
	if in.PasswordHash != "" {
		u.PasswordHash = in.PasswordHash
	}
	return u, r.db.Save(&u).Error
}

func (r *mysqlUserRepo) Delete(id int) (bool, error) {
	res := r.db.Delete(&models.User{}, id)
	return res.RowsAffected > 0, res.Error
}

// Auto-migrate + seed (giữ nguyên nếu bạn đã có)
