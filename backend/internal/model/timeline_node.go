package model

import "time"

// TimelineNode is a key milestone in an application.
type TimelineNode struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ApplicationID uint      `gorm:"index;not null" json:"application_id"`
	Title         string    `gorm:"size:255;not null" json:"title"`
	NodeType      string    `gorm:"size:32;index" json:"node_type"`
	DueDate       time.Time `gorm:"type:date;index" json:"due_date"`
	IsDone        bool      `gorm:"default:false" json:"is_done"`
	ReminderSent  bool      `gorm:"default:false" json:"reminder_sent"`
	CreatedAt     time.Time `json:"created_at"`
}
