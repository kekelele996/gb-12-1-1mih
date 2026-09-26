package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

// MessageService handles internal messages.
type MessageService struct {
	repo   *repository.MessageRepository
	logger *slog.Logger
}

// NewMessageService creates a MessageService.
func NewMessageService(repo *repository.MessageRepository, logger *slog.Logger) *MessageService {
	return &MessageService{repo: repo, logger: logger}
}

// Send creates a message from a user (or system, sender 0).
func (s *MessageService) Send(senderID, receiverID uint, content string) (*model.Message, error) {
	m := &model.Message{SenderID: senderID, ReceiverID: receiverID, Content: content}
	if err := s.repo.Create(m); err != nil {
		return nil, fmt.Errorf("message send: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogMessageSendSuccess, senderID, receiverID), "id", m.ID)
	return m, nil
}

// ListByReceiver returns messages for a user.
func (s *MessageService) ListByReceiver(receiverID uint, unreadOnly bool) ([]model.Message, error) {
	return s.repo.ListByReceiver(receiverID, unreadOnly)
}

// CountUnread returns unread count.
func (s *MessageService) CountUnread(receiverID uint) (int64, error) {
	return s.repo.CountUnread(receiverID)
}

// MarkRead marks a message as read (receiver only).
func (s *MessageService) MarkRead(userID, id uint) (*model.Message, error) {
	m, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("Message[id=%d] not found", id))
		}
		return nil, fmt.Errorf("message read find: %w", err)
	}
	if m.ReceiverID != userID {
		return nil, util.NewAppError(403, constants.CodeForbidden,
			fmt.Sprintf("Message[id=%d] read failed: not receiver", id))
	}
	m.IsRead = true
	if err := s.repo.Update(m); err != nil {
		return nil, fmt.Errorf("message read update: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogMessageRead, id), "id", id)
	return m, nil
}
