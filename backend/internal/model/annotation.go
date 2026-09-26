package model

import "time"

// Annotation is a counselor comment on a document with offset range.
type Annotation struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	DocumentID  uint      `gorm:"index;not null" json:"document_id"`
	CounselorID uint      `gorm:"index;not null" json:"counselor_id"`
	Content     string    `gorm:"size:1000;not null" json:"content"`
	StartOffset int       `json:"start_offset"`
	EndOffset   int       `json:"end_offset"`
	CreatedAt   time.Time `json:"created_at"`
}
