package repository

import (
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

// TimelineNodeRepository handles timeline persistence.
type TimelineNodeRepository struct{ db *gorm.DB }

// NewTimelineNodeRepository creates the repository.
func NewTimelineNodeRepository(db *gorm.DB) *TimelineNodeRepository { return &TimelineNodeRepository{db: db} }

// Create inserts a node.
func (r *TimelineNodeRepository) Create(n *model.TimelineNode) error { return translate(r.db.Create(n).Error) }

// FindByID locates a node by id.
func (r *TimelineNodeRepository) FindByID(id uint) (*model.TimelineNode, error) {
	var n model.TimelineNode
	if err := translate(r.db.First(&n, id).Error); err != nil {
		return nil, err
	}
	return &n, nil
}

// Update persists a node.
func (r *TimelineNodeRepository) Update(n *model.TimelineNode) error { return translate(r.db.Save(n).Error) }

// ListByApplication returns nodes of an application.
func (r *TimelineNodeRepository) ListByApplication(applicationID uint) ([]model.TimelineNode, error) {
	var items []model.TimelineNode
	if err := r.db.Where("application_id = ?", applicationID).Order("due_date ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// ListAll returns all nodes (for reminder scan).
func (r *TimelineNodeRepository) ListAll() ([]model.TimelineNode, error) {
	var items []model.TimelineNode
	if err := r.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// MarkDone sets a node as done.
func (r *TimelineNodeRepository) MarkDone(id uint) error {
	return r.db.Model(&model.TimelineNode{}).Where("id = ?", id).Update("is_done", true).Error
}

// MarkReminderSent marks a node's reminder as sent.
func (r *TimelineNodeRepository) MarkReminderSent(id uint) error {
	return r.db.Model(&model.TimelineNode{}).Where("id = ?", id).Update("reminder_sent", true).Error
}
