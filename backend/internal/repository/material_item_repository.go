package repository

import (
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

// MaterialItemRepository handles material checklist persistence.
type MaterialItemRepository struct{ db *gorm.DB }

// NewMaterialItemRepository creates the repository.
func NewMaterialItemRepository(db *gorm.DB) *MaterialItemRepository {
	return &MaterialItemRepository{db: db}
}

// Create inserts a material item.
func (r *MaterialItemRepository) Create(m *model.MaterialItem) error { return translate(r.db.Create(m).Error) }

// FindByID locates an item by id.
func (r *MaterialItemRepository) FindByID(id uint) (*model.MaterialItem, error) {
	var m model.MaterialItem
	if err := translate(r.db.First(&m, id).Error); err != nil {
		return nil, err
	}
	return &m, nil
}

// Update persists an item.
func (r *MaterialItemRepository) Update(m *model.MaterialItem) error { return translate(r.db.Save(m).Error) }

// ListByApplication returns items of an application.
func (r *MaterialItemRepository) ListByApplication(applicationID uint) ([]model.MaterialItem, error) {
	var items []model.MaterialItem
	if err := r.db.Where("application_id = ?", applicationID).Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
