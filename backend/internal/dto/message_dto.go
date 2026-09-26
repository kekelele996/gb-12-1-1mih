package dto

// MessageSendRequest sends a message.
type MessageSendRequest struct {
	ReceiverID uint   `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required,max=2000"`
}
