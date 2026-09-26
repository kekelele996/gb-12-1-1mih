package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/dto"
	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
	"github.com/gbstudyapply/gbstudyapply/internal/service"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

// MessageHandler exposes message endpoints.
type MessageHandler struct {
	svc    *service.MessageService
	logger *slog.Logger
}

// NewMessageHandler creates a MessageHandler.
func NewMessageHandler(svc *service.MessageService, logger *slog.Logger) *MessageHandler {
	return &MessageHandler{svc: svc, logger: logger}
}

// List handles GET /messages?unread=1.
func (h *MessageHandler) List(c *gin.Context) {
	unread := c.Query("unread") == "1"
	items, err := h.svc.ListByReceiver(middleware.GetUserID(c), unread)
	if err != nil {
		c.Error(err)
		return
	}
	unreadCount, _ := h.svc.CountUnread(middleware.GetUserID(c))
	c.JSON(http.StatusOK, dto.OK(gin.H{"items": items, "unread_count": unreadCount}))
}

// Send handles POST /messages.
func (h *MessageHandler) Send(c *gin.Context) {
	var req dto.MessageSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam+": "+err.Error()))
		return
	}
	m, err := h.svc.Send(middleware.GetUserID(c), req.ReceiverID, req.Content)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(m))
}

// MarkRead handles PUT /messages/:id/read.
func (h *MessageHandler) MarkRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid message id"))
		return
	}
	m, err := h.svc.MarkRead(middleware.GetUserID(c), uint(id))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(m))
}
