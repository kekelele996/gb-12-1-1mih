package model

import "time"

// Document is an essay/PS/RL/CV under an application.
type Document struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ApplicationID  uint      `gorm:"index;not null" json:"application_id"`
	DocType        string    `gorm:"size:16;not null" json:"doc_type"`
	Title          string    `gorm:"size:255;not null" json:"title"`
	Content        string    `gorm:"type:text" json:"content"`
	CurrentVersion int       `gorm:"default:1" json:"current_version"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
