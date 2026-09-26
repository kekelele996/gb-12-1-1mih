package model

import "time"

// DocumentVersion is a saved historical version of a document.
type DocumentVersion struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	DocumentID    uint      `gorm:"index;not null" json:"document_id"`
	Content       string    `gorm:"type:text" json:"content"`
	VersionNo     int       `json:"version_no"`
	ChangeSummary string    `gorm:"size:255" json:"change_summary"`
	CreatedBy     uint      `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}
