package dto

// DocumentCreateRequest creates a document.
type DocumentCreateRequest struct {
	DocType string `json:"doc_type" binding:"required"`
	Title   string `json:"title" binding:"required,max=255"`
	Content string `json:"content"`
}

// DocumentSaveRequest saves a new version.
type DocumentSaveRequest struct {
	Content       string `json:"content" binding:"required"`
	ChangeSummary string `json:"change_summary" binding:"omitempty,max=255"`
}

// RollbackRequest rolls back to a version.
type RollbackRequest struct {
	VersionNo int `json:"version_no" binding:"required"`
}

// AnnotationCreateRequest adds an annotation.
type AnnotationCreateRequest struct {
	Content     string `json:"content" binding:"required,max=1000"`
	StartOffset int    `json:"start_offset"`
	EndOffset   int    `json:"end_offset"`
}
