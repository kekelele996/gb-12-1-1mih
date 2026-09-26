package repository

import (
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

// RecommendationRepository handles recommendation persistence.
type RecommendationRepository struct{ db *gorm.DB }

// NewRecommendationRepository creates the repository.
func NewRecommendationRepository(db *gorm.DB) *RecommendationRepository {
	return &RecommendationRepository{db: db}
}

// Create inserts a recommendation.
func (r *RecommendationRepository) Create(rec *model.Recommendation) error {
	return translate(r.db.Create(rec).Error)
}

// ListByStudent returns recommendations for a student.
func (r *RecommendationRepository) ListByStudent(studentID uint) ([]model.Recommendation, error) {
	var items []model.Recommendation
	if err := r.db.Where("student_id = ?", studentID).Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
