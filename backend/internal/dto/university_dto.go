package dto

// UniversityCreateRequest creates/updates a university (admin).
type UniversityCreateRequest struct {
	Name                string `json:"name" binding:"required,max=128"`
	Country             string `json:"country" binding:"required,max=64"`
	City                string `json:"city" binding:"omitempty,max=64"`
	Ranking             int    `json:"ranking"`
	TopMajors           string `json:"top_majors"`
	ApplicationDeadline string `json:"application_deadline" binding:"omitempty,max=32"`
	TuitionRange        string `json:"tuition_range" binding:"omitempty,max=64"`
	Requirements        string `json:"requirements"`
}
