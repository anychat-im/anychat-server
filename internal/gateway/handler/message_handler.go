package handler

import (
	"net/http"

	messagepb "github.com/anychat/server/api/proto/message"
	"github.com/anychat/server/internal/gateway/client"
	gwmiddleware "github.com/anychat/server/internal/gateway/middleware"
	"github.com/anychat/server/pkg/response"
	"github.com/gin-gonic/gin"
)

// MessageHandler 消息HTTP处理器
type MessageHandler struct {
	clientManager *client.Manager
}

// NewMessageHandler 创建消息处理器
func NewMessageHandler(clientManager *client.Manager) *MessageHandler {
	return &MessageHandler{
		clientManager: clientManager,
	}
}

type ackReadTriggersRequest struct {
	Events []readTriggerEvent `json:"events" binding:"required,min=1"`
}

type readTriggerEvent struct {
	MessageID      string `json:"message_id" binding:"required"`
	ClientAt       *int64 `json:"client_at,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
}

// AckReadTriggers 阅后即焚阅读触发回执
// @Summary      阅后即焚阅读触发回执
// @Description  客户端批量上报消息阅读触发事件，服务端据此启动阅后即焚计时
// @Tags         消息
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      ackReadTriggersRequest  true  "阅读触发事件"
// @Success      200      {object}  response.Response{data=object}  "成功"
// @Failure      400      {object}  response.Response  "参数错误"
// @Failure      401      {object}  response.Response  "未授权"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /messages/read-triggers [post]
func (h *MessageHandler) AckReadTriggers(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)

	var req ackReadTriggersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	events := make([]*messagepb.ReadTriggerEvent, 0, len(req.Events))
	for _, event := range req.Events {
		pbEvent := &messagepb.ReadTriggerEvent{
			MessageId: event.MessageID,
		}
		if event.ClientAt != nil {
			pbEvent.ClientAt = event.ClientAt
		}
		if event.IdempotencyKey != "" {
			pbEvent.IdempotencyKey = &event.IdempotencyKey
		}
		events = append(events, pbEvent)
	}

	resp, err := h.clientManager.Message().AckReadTriggers(c.Request.Context(), &messagepb.AckReadTriggersRequest{
		UserId: userID,
		Events: events,
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}

	response.Success(c, gin.H{
		"success_ids": resp.SuccessIds,
		"ignored_ids": resp.IgnoredIds,
	})
}
