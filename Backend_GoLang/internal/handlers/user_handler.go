package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"crud_api_us/internal/repository"
	"crud_api_us/internal/services"
)

// Handler nắm Service
type UserHandler struct{ svc *services.UserService }

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{svc: services.NewUserService(repo)}
}

/************* DTO (request) *************/
type CreateUserRequest struct {
	Username   string `json:"username"     binding:"required,min=3,max=50"`
	Email      string `json:"email"        binding:"required,email"`
	Password   string `json:"password"     binding:"required,min=6,max=100"`
	FullName   string `json:"full_name"    binding:"omitempty,max=100"`
	Phone      string `json:"phone"        binding:"omitempty,max=20"`
	Gender     string `json:"gender"       binding:"omitempty,oneof=male female other"`
	DOB        string `json:"date_of_birth" binding:"omitempty"` // yyyy-mm-dd
	AvatarURL  string `json:"avatar_url"   binding:"omitempty,url"`
	Street     string `json:"street"       binding:"omitempty,max=255"`
	City       string `json:"city"         binding:"omitempty,max=100"`
	State      string `json:"state"        binding:"omitempty,max=100"`
	Country    string `json:"country"      binding:"omitempty,max=100"`
	PostalCode string `json:"postal_code"  binding:"omitempty,max=20"`
	Role       string `json:"role"         binding:"omitempty,oneof=user admin"`
	Status     string `json:"status"       binding:"omitempty,oneof=active inactive banned"`
}

type UpdateUserRequest struct {
	Username   string `json:"username"     binding:"required,min=3,max=50"`
	Email      string `json:"email"        binding:"required,email"`
	Password   string `json:"password"     binding:"omitempty,min=6,max=100"`
	FullName   string `json:"full_name"    binding:"omitempty,max=100"`
	Phone      string `json:"phone"        binding:"omitempty,max=20"`
	Gender     string `json:"gender"       binding:"omitempty,oneof=male female other"`
	DOB        string `json:"date_of_birth" binding:"omitempty"`
	AvatarURL  string `json:"avatar_url"   binding:"omitempty,url"`
	Street     string `json:"street"       binding:"omitempty,max=255"`
	City       string `json:"city"         binding:"omitempty,max=100"`
	State      string `json:"state"        binding:"omitempty,max=100"`
	Country    string `json:"country"      binding:"omitempty,max=100"`
	PostalCode string `json:"postal_code"  binding:"omitempty,max=20"`
	Role       string `json:"role"         binding:"omitempty,oneof=user admin"`
	Status     string `json:"status"       binding:"omitempty,oneof=active inactive banned"`
}

/************* Helpers *************/
func writeErr(c *gin.Context, code int, msg string) { c.JSON(code, gin.H{"error": msg}) }

/************* Handlers + Swagger *************/

// ListUsers godoc
// @Summary      Danh sách người dùng
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}  UserDoc
// @Router       /admin/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.svc.List()
	if err != nil {
		writeErr(c, http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary      Lấy người dùng theo ID
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Param        id  path      int  true  "User ID"
// @Success      200  {object} UserDoc
// @Failure      404  {object} ErrorResponse
// @Router       /admin/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	u, err := h.svc.Get(id)
	if err != nil {
		if err == repository.ErrNotFound {
			writeErr(c, http.StatusNotFound, "not found")
			return
		}
		writeErr(c, http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, u)
}

// CreateUser godoc
// @Summary      Tạo người dùng
// @Tags         Admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        user  body     CreateUserRequest  true  "User payload"
// @Success      201   {object} UserDoc
// @Failure      400   {object} ErrorResponse
// @Failure      409   {object} ErrorResponse
// @Router       /admin/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var in CreateUserRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		writeErr(c, http.StatusBadRequest, "invalid body")
		return
	}
	out, err := h.svc.Create(services.CreateParams{
		Username: in.Username, Email: in.Email, Password: in.Password,
		FullName: in.FullName, Phone: in.Phone, Gender: in.Gender, DOB: in.DOB,
		AvatarURL: in.AvatarURL, Street: in.Street, City: in.City, State: in.State,
		Country: in.Country, PostalCode: in.PostalCode, Role: in.Role, Status: in.Status,
	})
	if err != nil {
		switch err {
		case services.ErrDuplicate:
			writeErr(c, http.StatusConflict, "username/email already exists")
		case services.ErrBadInput:
			writeErr(c, http.StatusBadRequest, "invalid body")
		default:
			writeErr(c, http.StatusInternalServerError, "server error")
		}
		return
	}
	c.JSON(http.StatusCreated, out)
}

// UpdateUser godoc
// @Summary      Cập nhật người dùng
// @Tags         Admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path     int                true  "User ID"
// @Param        user  body     UpdateUserRequest  true  "User payload"
// @Success      200   {object} UserDoc
// @Failure      400   {object} ErrorResponse
// @Failure      404   {object} ErrorResponse
// @Failure      409   {object} ErrorResponse
// @Router       /admin/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var in UpdateUserRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		writeErr(c, http.StatusBadRequest, "invalid body")
		return
	}
	out, err := h.svc.Update(id, services.UpdateParams{
		Username: in.Username, Email: in.Email, Password: in.Password,
		FullName: in.FullName, Phone: in.Phone, Gender: in.Gender, DOB: in.DOB,
		AvatarURL: in.AvatarURL, Street: in.Street, City: in.City, State: in.State,
		Country: in.Country, PostalCode: in.PostalCode, Role: in.Role, Status: in.Status,
	})
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			writeErr(c, http.StatusNotFound, "not found")
		case services.ErrDuplicate:
			writeErr(c, http.StatusConflict, "username/email already exists")
		case services.ErrBadInput:
			writeErr(c, http.StatusBadRequest, "invalid body")
		default:
			writeErr(c, http.StatusInternalServerError, "server error")
		}
		return
	}
	c.JSON(http.StatusOK, out)
}

// DeleteUser godoc
// @Summary      Xoá người dùng
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  int  true  "User ID"
// @Success      204  {string} string "No Content"
// @Failure      404  {object}  ErrorResponse
// @Router       /admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ok, err := h.svc.Delete(id)
	if err != nil {
		writeErr(c, http.StatusInternalServerError, "server error")
		return
	}
	if !ok {
		writeErr(c, http.StatusNotFound, "not found")
		return
	}
	c.Status(http.StatusNoContent)
}
