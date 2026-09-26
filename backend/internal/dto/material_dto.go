package dto

// MaterialCreateRequest adds a checklist item.
type MaterialCreateRequest struct {
	Name       string `json:"name" binding:"required,max=128"`
	Category   string `json:"category" binding:"omitempty,max=64"`
	IsRequired bool   `json:"is_required"`
}

// MaterialStatusRequest updates item status.
//   - student: status=uploaded + file_url
//   - counselor/admin: status=approved/rejected (+ optional review_remark)
type MaterialStatusRequest struct {
	Status       string `json:"status" binding:"required"`
	FileURL      string `json:"file_url" binding:"omitempty,max=255"`
	ReviewRemark string `json:"review_remark" binding:"omitempty,max=500"`
}
