package dto

// RecommendationCreateRequest makes a school-selection plan.
type RecommendationCreateRequest struct {
	StudentID     uint   `json:"student_id" binding:"required"`
	UniversityIDs []uint `json:"university_ids" binding:"required,min=1"`
	Reason        string `json:"reason" binding:"omitempty,max=2000"`
}
