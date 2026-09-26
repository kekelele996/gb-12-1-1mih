package dto

// ApplicationCreateRequest creates a project.
type ApplicationCreateRequest struct {
	UniversityID uint   `json:"university_id" binding:"required"`
	Major        string `json:"major" binding:"required,max=128"`
	Round        string `json:"round" binding:"omitempty,max=32"`
}

// ApplicationStatusRequest updates project status.
type ApplicationStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
