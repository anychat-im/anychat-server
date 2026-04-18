package handler

import (
	"strconv"

	friendpb "github.com/anychat/server/api/proto/friend"
	friendmodel "github.com/anychat/server/internal/friend/model"
	"github.com/anychat/server/internal/gateway/client"
	gwmiddleware "github.com/anychat/server/internal/gateway/middleware"
	"github.com/anychat/server/pkg/response"
	"github.com/gin-gonic/gin"
)

// FriendHandler friend HTTP handler
type FriendHandler struct {
	clientManager *client.Manager
}

// NewFriendHandler creates friend handler
func NewFriendHandler(clientManager *client.Manager) *FriendHandler {
	return &FriendHandler{
		clientManager: clientManager,
	}
}

// GetFriends get friend list
// @Summary      get friend list
// @Description  Get all friends of current user, supports incremental sync
// @Tags         friend
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        lastUpdateTime  query  int64  false  "last update timestamp (incremental sync)"
// @Success      200  {object}  response.Response{data=object}  "success"
// @Failure      401  {object}  response.Response  "unauthorized"
// @Failure      500  {object}  response.Response  "server error"
// @Router       /friends [get]
func (h *FriendHandler) GetFriends(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)

	// Parse query parameters
	var lastUpdateTime *int64
	if timeStr := c.Query("last_update_time"); timeStr != "" {
		if t, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
			lastUpdateTime = &t
		}
	}

	resp, err := h.clientManager.Friend().GetFriendList(c.Request.Context(), &friendpb.GetFriendListRequest{
		UserId:         userID,
		LastUpdateTime: lastUpdateTime,
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}

	response.Success(c, resp)
}

// SendFriendRequest send friend request
// @Summary      send friend request
// @Description  Send friend request to specified user
// @Tags         friend
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  object  true  "request info"
// @Success      200  {object}  response.Response{data=object}  "success"
// @Failure      400  {object}  response.Response  "parameter error"
// @Failure      401  {object}  response.Response  "unauthorized"
// @Failure      500  {object}  response.Response  "server error"
// @Router       /friends/requests [post]
func (h *FriendHandler) SendFriendRequest(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)

	var req struct {
		UserID  string `json:"user_id" binding:"required" example:"user-456"`
		Message string `json:"message" example:"你好,我想加你为好友"`
		Source  int16  `json:"source" binding:"required,oneof=1 2 3 4" example:"1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}
	source := friendmodel.FriendRequestSource(req.Source)
	if !source.IsValid() {
		response.ParamError(c, "invalid source")
		return
	}

	resp, err := h.clientManager.Friend().SendFriendRequest(c.Request.Context(), &friendpb.SendFriendRequestRequest{
		FromUserId: userID,
		ToUserId:   req.UserID,
		Message:    req.Message,
		Source:     friendpb.FriendRequestSource(source),
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}

	response.Success(c, resp)
}

// HandleFriendRequest handle friend request
// @Summary      handle friend request
// @Description  Handle friend request with numeric action
// @Tags         friend
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path  int  true  "request ID"
// @Param        request  body  object  true  "handle action: 1-accept 2-reject"
// @Success      200  {object}  response.Response  "success"
// @Failure      400  {object}  response.Response  "parameter error"
// @Failure      401  {object}  response.Response  "unauthorized"
// @Failure      500  {object}  response.Response  "server error"
// @Router       /friends/requests/{id} [put]
func (h *FriendHandler) HandleFriendRequest(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "invalid request id")
		return
	}

	var req struct {
		Action int16 `json:"action" binding:"required,oneof=1 2" example:"1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	action := friendmodel.FriendRequestAction(req.Action)
	if !action.IsValid() {
		response.ParamError(c, "invalid action")
		return
	}

	_, err = h.clientManager.Friend().HandleFriendRequest(c.Request.Context(), &friendpb.HandleFriendRequestRequest{
		UserId:    userID,
		RequestId: requestID,
		Action:    friendpb.FriendRequestAction(action),
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}

	response.Success(c, nil)
}

// GetFriendRequests get friend request list
// @Summary      get friend request list
// @Description  Get friend request list by numeric request_type
// @Tags         friend
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request_type  query  int  false  "request type: 1-received 2-sent"  default(1)
// @Success      200  {object}  response.Response{data=object}  "success"
// @Failure      401  {object}  response.Response  "unauthorized"
// @Failure      500  {object}  response.Response  "server error"
// @Router       /friends/requests [get]
func (h *FriendHandler) GetFriendRequests(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)

	requestType := friendmodel.FriendRequestQueryTypeReceived
	if v := c.Query("request_type"); v != "" {
		parsed, err := strconv.ParseInt(v, 10, 16)
		if err != nil {
			response.ParamError(c, "invalid request_type")
			return
		}
		requestType = friendmodel.FriendRequestQueryType(parsed)
	}
	if !requestType.IsValid() {
		response.ParamError(c, "invalid request_type")
		return
	}

	resp, err := h.clientManager.Friend().GetFriendRequests(c.Request.Context(), &friendpb.GetFriendRequestsRequest{
		UserId:      userID,
		RequestType: friendpb.FriendRequestQueryType(requestType),
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}

	response.Success(c, resp)
}

// DeleteFriend delete friend
// @Summary      delete friend
// @Description  Delete specified friend
// @Tags         friend
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "friend user ID"
// @Success      200  {object}  response.Response  "success"
// @Failure      401  {object}  response.Response  "unauthorized"
// @Failure      500  {object}  response.Response  "server error"
// @Router       /friends/{id} [delete]
func (h *FriendHandler) DeleteFriend(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)
	friendID := c.Param("id")

	_, err := h.clientManager.Friend().DeleteFriend(c.Request.Context(), &friendpb.DeleteFriendRequest{
		UserId:   userID,
		FriendId: friendID,
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}

	response.Success(c, nil)
}

// UpdateRemark update friend remark
// @Summary      update friend remark
// @Description  Update remark name for specified friend
// @Tags         friend
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path  string  true  "friend user ID"
// @Param        request  body  object  true  "remark info"
// @Success      200  {object}  response.Response  "success"
// @Failure      400  {object}  response.Response  "parameter error"
// @Failure      401  {object}  response.Response  "unauthorized"
// @Failure      500  {object}  response.Response  "server error"
// @Router       /friends/{id}/remark [put]
func (h *FriendHandler) UpdateRemark(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)
	friendID := c.Param("id")

	var req struct {
		Remark string `json:"remark" binding:"max=50" example:"老朋友"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	_, err := h.clientManager.Friend().UpdateRemark(c.Request.Context(), &friendpb.UpdateRemarkRequest{
		UserId:   userID,
		FriendId: friendID,
		Remark:   req.Remark,
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}

	response.Success(c, nil)
}

// AddToBlacklist add to blacklist
// @Summary      add to blacklist
// @Description  Add specified user to blacklist
// @Tags         friend
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  object  true  "user ID"
// @Success      200  {object}  response.Response  "success"
// @Failure      400  {object}  response.Response  "parameter error"
// @Failure      401  {object}  response.Response  "unauthorized"
// @Failure      500  {object}  response.Response  "server error"
// @Router       /friends/blacklist [post]
func (h *FriendHandler) AddToBlacklist(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)

	var req struct {
		UserID string `json:"user_id" binding:"required" example:"user-456"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	_, err := h.clientManager.Friend().AddToBlacklist(c.Request.Context(), &friendpb.AddToBlacklistRequest{
		UserId:        userID,
		BlockedUserId: req.UserID,
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}

	response.Success(c, nil)
}

// RemoveFromBlacklist remove from blacklist
// @Summary      remove from blacklist
// @Description  Remove specified user from blacklist
// @Tags         friend
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "blocked user ID"
// @Success      200  {object}  response.Response  "success"
// @Failure      401  {object}  response.Response  "unauthorized"
// @Failure      500  {object}  response.Response  "server error"
// @Router       /friends/blacklist/{id} [delete]
func (h *FriendHandler) RemoveFromBlacklist(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)
	blockedUserID := c.Param("id")

	_, err := h.clientManager.Friend().RemoveFromBlacklist(c.Request.Context(), &friendpb.RemoveFromBlacklistRequest{
		UserId:        userID,
		BlockedUserId: blockedUserID,
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}

	response.Success(c, nil)
}

// GetBlacklist get blacklist
// @Summary      get blacklist
// @Description  Get current user's blacklist
// @Tags         friend
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=object}  "success"
// @Failure      401  {object}  response.Response  "unauthorized"
// @Failure      500  {object}  response.Response  "server error"
// @Router       /friends/blacklist [get]
func (h *FriendHandler) GetBlacklist(c *gin.Context) {
	userID := gwmiddleware.GetUserID(c)

	resp, err := h.clientManager.Friend().GetBlacklist(c.Request.Context(), &friendpb.GetBlacklistRequest{
		UserId: userID,
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}

	response.Success(c, resp)
}
