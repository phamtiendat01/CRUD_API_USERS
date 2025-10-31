package models

import (
	"time"

	"gorm.io/gorm"
)

// User: bảng người dùng với đầy đủ trường chuyên nghiệp
type User struct {
	ID           int        `json:"id"            gorm:"primaryKey;autoIncrement"`
	Username     string     `json:"username"      gorm:"type:varchar(50);uniqueIndex;not null"`
	Email        string     `json:"email"         gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string     `json:"-"             gorm:"type:varchar(255);not null"` // ẩn trong JSON
	FullName     string     `json:"full_name"     gorm:"type:varchar(100)"`
	Phone        string     `json:"phone"         gorm:"type:varchar(20);index"`
	Gender       string     `json:"gender"        gorm:"type:varchar(10)"` // male|female|other
	DateOfBirth  *time.Time `json:"date_of_birth" gorm:"type:date"`

	AvatarURL string `json:"avatar_url" gorm:"type:varchar(255)"`

	Street     string `json:"street"      gorm:"type:varchar(255)"`
	City       string `json:"city"        gorm:"type:varchar(100)"`
	State      string `json:"state"       gorm:"type:varchar(100)"`
	Country    string `json:"country"     gorm:"type:varchar(100)"`
	PostalCode string `json:"postal_code" gorm:"type:varchar(20)"`

	Role   string `json:"role"   gorm:"type:varchar(20);default:user"`         // user|admin
	Status string `json:"status" gorm:"type:varchar(20);default:active;index"` // active|inactive|banned

	LastLoginAt *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"` // soft delete
}
