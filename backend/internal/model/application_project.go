package model

import "time"

// ApplicationProject tracks one school application for a student.
type ApplicationProject struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	StudentID    uint      `gorm:"index;not null" json:"student_id"`
	CounselorID  uint      `gorm:"index" json:"counselor_id"`
	UniversityID uint      `gorm:"index;not null" json:"university_id"`
	Major        string    `gorm:"size:128;not null" json:"major"`
	Round        string    `gorm:"size:32" json:"round"`
	Status       string    `gorm:"size:16;default:planning;index" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
