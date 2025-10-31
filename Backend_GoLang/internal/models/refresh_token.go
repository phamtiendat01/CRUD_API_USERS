package models

import "time"

type RefreshToken struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	TokenID   string    `gorm:"type:varchar(64);uniqueIndex;not null"` // jti
	UserID    int       `gorm:"index;not null"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ExpiresAt time.Time `gorm:"index;not null"`
	Revoked   bool      `gorm:"index;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
