package model

import "time"

// MaterialItem is a checklist item for an application.
type MaterialItem struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ApplicationID uint      `gorm:"index;not null" json:"application_id"`
	Name          string    `gorm:"size:128;not null" json:"name"`
	Category      string    `gorm:"size:64" json:"category"`
	IsRequired    bool      `gorm:"default:true" json:"is_required"`
	Status        string    `gorm:"size:16;default:pending" json:"status"`
	FileURL       string    `gorm:"size:255" json:"file_url"`
	UploadedAt    *time.Time `json:"uploaded_at"`
	CreatedAt     time.Time `json:"created_at"`
}
