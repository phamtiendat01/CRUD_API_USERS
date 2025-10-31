package services

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"crud_api_us/internal/models"
	"crud_api_us/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	CookieName string
}

func LoadJWTConfigFromEnv() JWTConfig {
	_ = godotenv.Load()
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-me"
	}
	access, _ := time.ParseDuration(getEnv("ACCESS_TOKEN_TTL", "15m"))
	refresh, _ := time.ParseDuration(getEnv("REFRESH_TOKEN_TTL", "168h"))
	cname := getEnv("REFRESH_COOKIE_NAME", "refresh_token")
	return JWTConfig{Secret: secret, AccessTTL: access, RefreshTTL: refresh, CookieName: cname}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

type AuthService struct {
	users repository.UserRepository
	auth  repository.AuthRepository
	jwt   JWTConfig
}

func NewAuthService(users repository.UserRepository, auth repository.AuthRepository, cfg JWTConfig) *AuthService {
	return &AuthService{users: users, auth: auth, jwt: cfg}
}

// ---------- helpers ----------
func (s *AuthService) hashPassword(pw string) (string, error) {
	if strings.TrimSpace(pw) == "" {
		return "", errors.New("bad_password")
	}
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}
func (s *AuthService) checkPassword(hash, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}

type Claims struct {
	UserID   int    `json:"uid"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) makeToken(user models.User, ttl time.Duration, jti string) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(ttl)
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject:   strconv.Itoa(user.ID),
			ID:        jti,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwt.Secret))
	return signed, exp, err
}

// ---------- API ----------
type LoginResult struct {
	AccessToken string
	AccessExp   time.Time
	Refresh     string
	RefreshExp  time.Time
	User        models.User
}

func (s *AuthService) Login(identifier, password string) (LoginResult, error) {
	user, err := s.auth.FindByUsernameOrEmail(identifier)
	if err != nil {
		return LoginResult{}, err
	}
	if err := s.checkPassword(user.PasswordHash, password); err != nil {
		return LoginResult{}, errors.New("invalid_credentials")
	}
	// access
	accessJTI := uuid.NewString()
	access, accessExp, err := s.makeToken(user, s.jwt.AccessTTL, accessJTI)
	if err != nil {
		return LoginResult{}, err
	}
	// refresh
	refreshJTI := uuid.NewString()
	refresh, refreshExp, err := s.makeToken(user, s.jwt.RefreshTTL, refreshJTI)
	if err != nil {
		return LoginResult{}, err
	}
	// persist refresh JTI
	if err := s.auth.SaveRefreshToken(&models.RefreshToken{
		TokenID:   refreshJTI,
		UserID:    user.ID,
		ExpiresAt: refreshExp,
	}); err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		AccessToken: access, AccessExp: accessExp,
		Refresh: refresh, RefreshExp: refreshExp, User: user,
	}, nil
}

func (s *AuthService) Logout(refreshJTI string) error {
	if strings.TrimSpace(refreshJTI) == "" {
		return nil
	}
	return s.auth.RevokeRefreshTokenByJTI(refreshJTI)
}
