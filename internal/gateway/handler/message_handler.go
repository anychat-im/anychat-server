package handler

import (
	"net/http"

	messagepb "github.com/anychat/server/api/proto/message"
	"github.com/anychat/server/internal/gateway/client"
	gwmiddleware "github.com/anychat/server/internal/gateway/middleware"
	"github.com/anychat/server/pkg/response"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
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

type recallMessageRequest struct {
	MessageID string `json:"message_id" binding:"required"`
}

type ackReadTriggersRequest struct {
	Events []readTriggerEvent `json:"events" binding:"required,min=1"`
}

type readTriggerEvent struct {
	MessageID      string `json:"message_id" binding:"required"`
	ClientAt       *int64 `json:"client_at,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
}

// RecallMessage 撤回消息
// @Summary      撤回消息
// @Description  撤回指定消息，只能撤回自己发送的消息，且需在2分钟内
// @Tags         消息
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      recallMessageRequest  true  "消息ID"
// @Success      200      {object}  response.Response  "成功"
// @Failure      400      {object}  response.Response  "参数错误"
// @Failure      401      {object}  response.Response  "未授权"
// @Failure      403      {object}  response.Response  "无权限或已超时"
// @Failure      404      {object}  response.Response  "消息不存在"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /messages/recall [post]
func (h *MessageHandler) RecallMessage(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)

	var req recallMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := metadata.AppendToOutgoingContext(c.Request.Context(), "x-user-id", userID)
	_, err := h.clientManager.Message().RecallMessage(ctx, &messagepb.RecallMessageRequest{
		MessageId: req.MessageID,
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}

	response.Success(c, nil)
}

// DeleteMessage 删除消息
// @Summary      删除消息
// @Description  删除指定消息，只能删除自己发送的消息
// @Tags         消息
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        messageId  path      string  true  "消息ID"
// @Success      200      {object}  response.Response  "成功"
// @Failure      401      {object}  response.Response  "未授权"
// @Failure      403      {object}  response.Response  "无权限"
// @Failure      404      {object}  response.Response  "消息不存在"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /messages/{messageId} [delete]
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)
	messageID := c.Param("messageId")

	if messageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message_id is required"})
		return
	}

	ctx := metadata.AppendToOutgoingContext(c.Request.Context(), "x-user-id", userID)
	_, err := h.clientManager.Message().DeleteMessage(ctx, &messagepb.DeleteMessageRequest{
		MessageId: messageID,
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}

	response.Success(c, nil)
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

	ctx := metadata.AppendToOutgoingContext(c.Request.Context(), "x-user-id", userID)
	resp, err := h.clientManager.Message().AckReadTriggers(ctx, &messagepb.AckReadTriggersRequest{
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
