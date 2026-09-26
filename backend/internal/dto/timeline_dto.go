package dto

import "time"

// TimelineCreateRequest adds a node.
type TimelineCreateRequest struct {
	Title    string    `json:"title" binding:"required,max=255"`
	NodeType string    `json:"node_type" binding:"omitempty,max=32"`
	DueDate  time.Time `json:"due_date" binding:"required"`
}
