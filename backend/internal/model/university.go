package model

import "time"

// University is a target school in the catalog.
type University struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	Name                string    `gorm:"size:128;uniqueIndex;not null" json:"name"`
	Country             string    `gorm:"size:64;index" json:"country"`
	City                string    `gorm:"size:64" json:"city"`
	Ranking             int       `gorm:"index" json:"ranking"`
	TopMajors           string    `gorm:"type:json" json:"top_majors"`
	ApplicationDeadline string    `gorm:"size:32" json:"application_deadline"`
	TuitionRange        string    `gorm:"size:64" json:"tuition_range"`
	Requirements        string    `gorm:"type:json" json:"requirements"`
	CreatedAt           time.Time `json:"created_at"`
}
