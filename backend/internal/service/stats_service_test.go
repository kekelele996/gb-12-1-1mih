package service

import (
	"testing"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

func TestComputeAppStats(t *testing.T) {
	projects := []model.ApplicationProject{
		{Status: "planning"},
		{Status: "submitted"},
		{Status: "admitted"},
		{Status: "admitted"},
		{Status: "rejected"},
	}
	s := ComputeAppStats(projects)
	if s.Total != 5 {
		t.Errorf("Total = %d, want 5", s.Total)
	}
	if s.Admitted != 2 {
		t.Errorf("Admitted = %d, want 2", s.Admitted)
	}
	if s.Applied != 3 {
		t.Errorf("Applied = %d, want 3", s.Applied)
	}
	if s.ByStatus["planning"] != 1 {
		t.Errorf("planning = %d", s.ByStatus["planning"])
	}
}
