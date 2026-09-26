package model

import "time"

// Recommendation is a counselor's school-selection plan.
type Recommendation struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	StudentID    uint      `gorm:"index;not null" json:"student_id"`
	CounselorID  uint      `gorm:"index;not null" json:"counselor_id"`
	UniversityIDs string   `gorm:"type:json" json:"university_ids"`
	Reason       string    `gorm:"type:text" json:"reason"`
	CreatedAt    time.Time `json:"created_at"`
}
