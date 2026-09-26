package repository

import (
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

// AnnotationRepository handles annotation persistence.
type AnnotationRepository struct{ db *gorm.DB }

// NewAnnotationRepository creates the repository.
func NewAnnotationRepository(db *gorm.DB) *AnnotationRepository { return &AnnotationRepository{db: db} }

// Create inserts an annotation.
func (r *AnnotationRepository) Create(a *model.Annotation) error { return translate(r.db.Create(a).Error) }

// ListByDocument returns annotations of a document.
func (r *AnnotationRepository) ListByDocument(documentID uint) ([]model.Annotation, error) {
	var items []model.Annotation
	if err := r.db.Where("document_id = ?", documentID).Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
