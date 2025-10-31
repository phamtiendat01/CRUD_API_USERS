package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"crud_api_us/internal/repository"
	"crud_api_us/internal/services"
)

/************ Handler ************/
type AuthHandler struct {
	users *UserHandler // dùng lại service tạo user
	auth  *services.AuthService
	cfg   services.JWTConfig
}

func NewAuthHandler(userRepo repository.UserRepository, authRepo repository.AuthRepository, cfg services.JWTConfig) *AuthHandler {
	return &AuthHandler{
		users: NewUserHandler(userRepo),
		auth:  services.NewAuthService(userRepo, authRepo, cfg),
		cfg:   cfg,
	}
}

/************ DTO (request) ************/
type RegisterRequest struct {
	Email           string `json:"email"            binding:"required,email"`
	Password        string `json:"password"         binding:"required,min=6,max=100"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`              // username hoặc email
	Password   string `json:"password"  binding:"required,min=6,max=100"` // mật khẩu
}

/************ DTO (docs/response) ************/
// Lưu ý: các struct dưới chỉ dùng cho Swagger docs

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserDoc struct {
	ID          int     `json:"id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	Role        string  `json:"role"`
	FullName    string  `json:"full_name,omitempty"`
	Phone       string  `json:"phone,omitempty"`
	Gender      string  `json:"gender,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	Street      string  `json:"street,omitempty"`
	City        string  `json:"city,omitempty"`
	State       string  `json:"state,omitempty"`
	Country     string  `json:"country,omitempty"`
	PostalCode  string  `json:"postal_code,omitempty"`
	Status      string  `json:"status,omitempty"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
}

type RegisterResponse struct {
	Message string  `json:"message"`
	User    UserDoc `json:"user"`
}

type LoginResponse struct {
	TokenType   string  `json:"token_type"`   // "Bearer"
	AccessToken string  `json:"access_token"` // JWT
	ExpiresIn   int     `json:"expires_in"`   // giây
	User        UserDoc `json:"user"`
}

/************ Helpers ************/
func (h *AuthHandler) setRefreshCookie(c *gin.Context, token string, exp time.Time) {
	c.SetCookie(h.cfg.CookieName, token, int(time.Until(exp).Seconds()),
		"/", "", false, true) // Path=/, HttpOnly; Secure=false cho localhost
}
func (h *AuthHandler) clearRefreshCookie(c *gin.Context) {
	c.SetCookie(h.cfg.CookieName, "", -1, "/", "", false, true)
}

/************ Endpoints ************/

// Register godoc
// @Summary      Đăng ký tài khoản mới
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        req  body     RegisterRequest  true  "Register payload"
// @Success      201  {object} RegisterResponse
// @Failure      400  {object} ErrorResponse
// @Failure      409  {object} ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var in RegisterRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		writeErr(c, http.StatusBadRequest, "invalid body")
		return
	}
	// username = phần trước @ của email
	at := strings.Index(in.Email, "@")
	username := in.Email
	if at > 0 {
		username = in.Email[:at]
	}

	out, err := h.users.svc.Create(services.CreateParams{
		Username: username,
		Email:    strings.ToLower(strings.TrimSpace(in.Email)),
		Password: in.Password,
		Role:     "user", // mặc định user
	})
	if err != nil {
		switch err {
		case services.ErrDuplicate:
			writeErr(c, http.StatusConflict, "email/username already exists")
		case services.ErrBadInput:
			writeErr(c, http.StatusBadRequest, "invalid body")
		default:
			writeErr(c, http.StatusInternalServerError, "server error")
		}
		return
	}

	// Không auto-login, FE sẽ chuyển sang form đăng nhập
	c.JSON(http.StatusCreated, gin.H{
		"message": "registered successfully, please login",
		"user":    out,
	})
}

// Login godoc
// @Summary      Đăng nhập (lấy access/refresh token)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        req  body     LoginRequest  true  "Login payload"
// @Success      200  {object} LoginResponse
// @Failure      400  {object} ErrorResponse
// @Failure      401  {object} ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var in LoginRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		writeErr(c, http.StatusBadRequest, "invalid body")
		return
	}
	res, err := h.auth.Login(in.Identifier, in.Password)
	if err != nil {
		if err == repository.ErrNotFound || strings.Contains(err.Error(), "invalid_credentials") {
			writeErr(c, http.StatusUnauthorized, "invalid credentials")
			return
		}
		writeErr(c, http.StatusInternalServerError, "server error")
		return
	}

	// set refresh token cookie
	h.setRefreshCookie(c, res.Refresh, res.RefreshExp)

	c.JSON(http.StatusOK, gin.H{
		"token_type":   "Bearer",
		"access_token": res.AccessToken,
		"expires_in":   int(time.Until(res.AccessExp).Seconds()),
		"user": gin.H{
			"id": res.User.ID, "username": res.User.Username, "email": res.User.Email, "role": res.User.Role,
		},
	})
}

// Logout godoc
// @Summary      Đăng xuất (revoke refresh token & clear cookie)
// @Tags         Auth
// @Produce      json
// @Success      204  {string} string "No Content"
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	cookie, err := c.Cookie(h.cfg.CookieName)
	if err == nil && cookie != "" {
		// lấy jti từ refresh token
		tok, _ := jwt.ParseWithClaims(cookie, &services.Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(h.cfg.Secret), nil
		})
		if claims, ok := tok.Claims.(*services.Claims); ok && tok.Valid {
			_ = h.auth.Logout(claims.ID)
		}
	}
	h.clearRefreshCookie(c)
	c.Status(http.StatusNoContent)
}

// Me godoc
// @Summary      Thông tin người dùng hiện tại
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object} UserDoc
// @Failure      401  {object} ErrorResponse
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	uid := c.GetInt("uid")
	u, err := h.users.svc.Get(uid)
	if err != nil {
		writeErr(c, http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, u)
}
