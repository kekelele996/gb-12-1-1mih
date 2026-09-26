package repository

import (
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

// MessageRepository handles message persistence.
type MessageRepository struct{ db *gorm.DB }

// NewMessageRepository creates the repository.
func NewMessageRepository(db *gorm.DB) *MessageRepository { return &MessageRepository{db: db} }

// Create inserts a message.
func (r *MessageRepository) Create(m *model.Message) error { return translate(r.db.Create(m).Error) }

// FindByID locates a message by id.
func (r *MessageRepository) FindByID(id uint) (*model.Message, error) {
	var m model.Message
	if err := translate(r.db.First(&m, id).Error); err != nil {
		return nil, err
	}
	return &m, nil
}

// Update persists a message.
func (r *MessageRepository) Update(m *model.Message) error { return translate(r.db.Save(m).Error) }

// ListByReceiver returns messages for a user, optionally unread only.
func (r *MessageRepository) ListByReceiver(receiverID uint, unreadOnly bool) ([]model.Message, error) {
	var items []model.Message
	q := r.db.Where("receiver_id = ?", receiverID)
	if unreadOnly {
		q = q.Where("is_read = ?", false)
	}
	if err := q.Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// CountUnread returns the unread message count for a user.
func (r *MessageRepository) CountUnread(receiverID uint) (int64, error) {
	var total int64
	if err := r.db.Model(&model.Message{}).Where("receiver_id = ? AND is_read = ?", receiverID, false).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}
