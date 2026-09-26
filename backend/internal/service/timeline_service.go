package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

// TimelineService handles timeline nodes.
type TimelineService struct {
	repo   *repository.TimelineNodeRepository
	logger *slog.Logger
}

// NewTimelineService creates a TimelineService.
func NewTimelineService(repo *repository.TimelineNodeRepository, logger *slog.Logger) *TimelineService {
	return &TimelineService{repo: repo, logger: logger}
}

// Create adds a node to an application.
func (s *TimelineService) Create(applicationID uint, n *model.TimelineNode) (*model.TimelineNode, error) {
	n.ApplicationID = applicationID
	if err := s.repo.Create(n); err != nil {
		return nil, fmt.Errorf("timeline create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogTimelineCreateSuccess, applicationID, n.Title), "id", n.ID)
	return n, nil
}

// ListByApplication returns nodes of an application.
func (s *TimelineService) ListByApplication(applicationID uint) ([]model.TimelineNode, error) {
	return s.repo.ListByApplication(applicationID)
}

// MarkDone marks a node as done.
func (s *TimelineService) MarkDone(id uint) (*model.TimelineNode, error) {
	n, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("TimelineNode[id=%d] not found", id))
		}
		return nil, fmt.Errorf("timeline done find: %w", err)
	}
	if err := s.repo.MarkDone(id); err != nil {
		return nil, fmt.Errorf("timeline done update: %w", err)
	}
	n.IsDone = true
	s.logger.Info(fmt.Sprintf(constants.LogTimelineDone, id), "id", id)
	return n, nil
}

// Upcoming returns nodes due within days for a user's applications.
func (s *TimelineService) Upcoming(nodes []model.TimelineNode, days int) []model.TimelineNode {
	cutoff := time.Now().AddDate(0, 0, days)
	var out []model.TimelineNode
	for _, n := range nodes {
		if !n.IsDone && !n.ReminderSent && !n.DueDate.IsZero() && n.DueDate.Before(cutoff) {
			out = append(out, n)
		}
	}
	return out
}
