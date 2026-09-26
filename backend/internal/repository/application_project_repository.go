package repository

import (
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

// ApplicationProjectRepository handles project persistence.
type ApplicationProjectRepository struct{ db *gorm.DB }

// NewApplicationProjectRepository creates the repository.
func NewApplicationProjectRepository(db *gorm.DB) *ApplicationProjectRepository {
	return &ApplicationProjectRepository{db: db}
}

// Create inserts a project.
func (r *ApplicationProjectRepository) Create(a *model.ApplicationProject) error {
	return translate(r.db.Create(a).Error)
}

// FindByID locates a project by id.
func (r *ApplicationProjectRepository) FindByID(id uint) (*model.ApplicationProject, error) {
	var a model.ApplicationProject
	if err := translate(r.db.First(&a, id).Error); err != nil {
		return nil, err
	}
	return &a, nil
}

// Update persists a project.
func (r *ApplicationProjectRepository) Update(a *model.ApplicationProject) error {
	return translate(r.db.Save(a).Error)
}

// ListByStudent returns projects of a student.
func (r *ApplicationProjectRepository) ListByStudent(studentID uint) ([]model.ApplicationProject, error) {
	var items []model.ApplicationProject
	if err := r.db.Where("student_id = ?", studentID).Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// ListByCounselor returns projects managed by a counselor.
func (r *ApplicationProjectRepository) ListByCounselor(counselorID uint) ([]model.ApplicationProject, error) {
	var items []model.ApplicationProject
	if err := r.db.Where("counselor_id = ?", counselorID).Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// ListAll returns all projects (admin/dashboard).
func (r *ApplicationProjectRepository) ListAll() ([]model.ApplicationProject, error) {
	var items []model.ApplicationProject
	if err := r.db.Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
