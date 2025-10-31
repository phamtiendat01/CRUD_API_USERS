package repository

import "crud_api_us/internal/models"

type AuthRepository interface {
	FindByUsernameOrEmail(identifier string) (models.User, error)
	SaveRefreshToken(token *models.RefreshToken) error
	RevokeRefreshTokenByJTI(jti string) error
}
