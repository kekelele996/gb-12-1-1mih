package model

import "time"

// User represents a platform account (student/counselor/admin).
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Email        string    `gorm:"size:128;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	RealName     string    `gorm:"size:64" json:"real_name"`
	Phone        string    `gorm:"size:32" json:"phone"`
	Role         string    `gorm:"size:16;default:student;index" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
