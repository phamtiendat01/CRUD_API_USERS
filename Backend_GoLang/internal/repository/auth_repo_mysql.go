package repository

import (
	"errors"

	"crud_api_us/internal/models"

	"gorm.io/gorm"
)

type mysqlAuthRepo struct{ db *gorm.DB }

func NewMySQLAuthRepo(db *gorm.DB) AuthRepository { return &mysqlAuthRepo{db: db} }

func (r *mysqlAuthRepo) FindByUsernameOrEmail(identifier string) (models.User, error) {
	var u models.User
	if err := r.db.Where("username = ? OR email = ?", identifier, identifier).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}
	return u, nil
}

func (r *mysqlAuthRepo) SaveRefreshToken(t *models.RefreshToken) error {
	return r.db.Create(t).Error
}

func (r *mysqlAuthRepo) RevokeRefreshTokenByJTI(jti string) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("token_id = ? AND revoked = ?", jti, false).
		Update("revoked", true).Error
}
