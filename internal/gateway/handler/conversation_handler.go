package handler

import (
	"net/http"
	"strconv"

	conversationpb "github.com/anychat/server/api/proto/conversation"
	"github.com/anychat/server/internal/gateway/client"
	gwmiddleware "github.com/anychat/server/internal/gateway/middleware"
	"github.com/anychat/server/pkg/response"
	"github.com/gin-gonic/gin"
)

// ConversationHandler conversation HTTP处理器
type ConversationHandler struct {
	clientManager *client.Manager
}

// NewConversationHandler 创建conversation处理器
func NewConversationHandler(clientManager *client.Manager) *ConversationHandler {
	return &ConversationHandler{clientManager: clientManager}
}

// GetConversations 获取会话列表
// @Summary      获取会话列表
// @Description  获取当前用户的会话列表，支持增量同步（通过updatedBefore参数）
// @Tags         会话
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit          query  int    false  "返回数量（默认20，最大100）"
// @Param        updatedBefore  query  int64  false  "Unix时间戳，仅返回此时间之前更新的会话（增量同步）"
// @Success      200  {object}  response.Response{data=object}  "成功"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /conversations [get]
func (h *ConversationHandler) GetConversations(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)

	req := &conversationpb.GetConversationsRequest{UserId: userID}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = int32(limit)
		}
	}
	if beforeStr := c.Query("updatedBefore"); beforeStr != "" {
		if t, err := strconv.ParseInt(beforeStr, 10, 64); err == nil {
			req.UpdatedBefore = &t
		}
	}

	resp, err := h.clientManager.Conversation().GetConversations(c.Request.Context(), req)
	if err != nil {
		handleGRPCError(c, err)
		return
	}
	response.Success(c, resp)
}

// GetConversation 获取单个会话
// @Summary      获取单个会话
// @Description  获取指定会话的详情
// @Tags         会话
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        conversationId  path  string  true  "会话ID"
// @Success      200  {object}  response.Response{data=object}  "成功"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      404  {object}  response.Response  "会话不存在"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /conversations/{conversationId} [get]
func (h *ConversationHandler) GetConversation(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)
	conversationID := c.Param("conversationId")

	resp, err := h.clientManager.Conversation().GetConversation(c.Request.Context(), &conversationpb.GetConversationRequest{
		UserId:         userID,
		ConversationId: conversationID,
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}
	response.Success(c, resp)
}

// DeleteConversation 删除会话
// @Summary      删除会话
// @Description  删除指定会话（不影响消息，仅从会话列表中移除）
// @Tags         会话
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        conversationId  path  string  true  "会话ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /conversations/{conversationId} [delete]
func (h *ConversationHandler) DeleteConversation(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)
	conversationID := c.Param("conversationId")

	_, err := h.clientManager.Conversation().DeleteConversation(c.Request.Context(), &conversationpb.DeleteConversationRequest{
		UserId:         userID,
		ConversationId: conversationID,
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}
	response.Success(c, nil)
}

// setPinnedRequest 置顶请求体
type setPinnedRequest struct {
	Pinned bool `json:"pinned"`
}

// SetPinned 设置会话置顶
// @Summary      设置会话置顶
// @Description  置顶或取消置顶指定会话
// @Tags         会话
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        conversationId  path  string           true  "会话ID"
// @Param        request    body  setPinnedRequest  true  "置顶状态"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /conversations/{conversationId}/pin [put]
func (h *ConversationHandler) SetPinned(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)
	conversationID := c.Param("conversationId")

	var req setPinnedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.clientManager.Conversation().SetPinned(c.Request.Context(), &conversationpb.SetPinnedRequest{
		UserId:         userID,
		ConversationId: conversationID,
		Pinned:         req.Pinned,
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}
	response.Success(c, nil)
}

// setMutedRequest 免打扰请求体
type setMutedRequest struct {
	Muted bool `json:"muted"`
}

// SetMuted 设置会话免打扰
// @Summary      设置会话免打扰
// @Description  开启或关闭指定会话的免打扰
// @Tags         会话
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        conversationId  path  string          true  "会话ID"
// @Param        request    body  setMutedRequest  true  "免打扰状态"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /conversations/{conversationId}/mute [put]
func (h *ConversationHandler) SetMuted(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)
	conversationID := c.Param("conversationId")

	var req setMutedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.clientManager.Conversation().SetMuted(c.Request.Context(), &conversationpb.SetMutedRequest{
		UserId:         userID,
		ConversationId: conversationID,
		Muted:          req.Muted,
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}
	response.Success(c, nil)
}

// MarkRead 标记会话已读（清除未读数）
// @Summary      标记会话已读
// @Description  清除指定会话的未读消息数
// @Tags         会话
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        conversationId  path  string  true  "会话ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /conversations/{conversationId}/read [post]
func (h *ConversationHandler) MarkRead(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)
	conversationID := c.Param("conversationId")

	_, err := h.clientManager.Conversation().ClearUnread(c.Request.Context(), &conversationpb.ClearUnreadRequest{
		UserId:         userID,
		ConversationId: conversationID,
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}
	response.Success(c, nil)
}

// GetTotalUnread 获取总未读数
// @Summary      获取总未读数
// @Description  获取当前用户所有会话的总未读消息数（免打扰会话不计入）
// @Tags         会话
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=object}  "成功"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /conversations/unread/total [get]
func (h *ConversationHandler) GetTotalUnread(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)

	resp, err := h.clientManager.Conversation().GetTotalUnread(c.Request.Context(), &conversationpb.GetTotalUnreadRequest{
		UserId: userID,
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}
	response.Success(c, resp)
}

// setBurnAfterReadingRequest 阅后即焚请求体
type setBurnAfterReadingRequest struct {
	Duration int32 `json:"duration"` // 秒,0表示取消
}

// SetBurnAfterReading 设置阅后即焚
// @Summary      设置阅后即焚
// @Description  设置会话阅后即焚时长，0表示取消
// @Tags         会话
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        conversationId  path  string                     true  "会话ID"
// @Param        request    body  setBurnAfterReadingRequest  true  "阅后即焚时长(秒)"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /conversations/{conversationId}/burn [put]
func (h *ConversationHandler) SetBurnAfterReading(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)
	conversationID := c.Param("conversationId")

	var req setBurnAfterReadingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.clientManager.Conversation().SetBurnAfterReading(c.Request.Context(), &conversationpb.SetBurnAfterReadingRequest{
		UserId:         userID,
		ConversationId: conversationID,
		Duration:       req.Duration,
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}
	response.Success(c, nil)
}

type setAutoDeleteRequest struct {
	Duration int32 `json:"duration"` // 秒,0表示取消
}

// SetAutoDelete 设置自动删除
// @Summary      设置自动删除
// @Description  设置会话自动删除时长，0表示取消
// @Tags         会话
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        conversationId  path  string                  true  "会话ID"
// @Param        request    body  setAutoDeleteRequest   true  "自动删除时长(秒)"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /conversations/{conversationId}/auto_delete [put]
func (h *ConversationHandler) SetAutoDelete(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)
	conversationID := c.Param("conversationId")

	var req setAutoDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.clientManager.Conversation().SetAutoDelete(c.Request.Context(), &conversationpb.SetAutoDeleteRequest{
		UserId:         userID,
		ConversationId: conversationID,
		Duration:       req.Duration,
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}
	response.Success(c, nil)
}
