package dto

import "time"

// SendFriendRequestRequest is the request for sending a friend request
type SendFriendRequestRequest struct {
	UserID  string `json:"user_id" binding:"required" example:"user-123"`
	Message string `json:"message" binding:"max=200" example:"Hello, I'd like to add you as a friend"`
	Source  string `json:"source" binding:"required,oneof=search qrcode group contacts" example:"search"`
}

// HandleFriendRequestRequest is the request for handling a friend request
type HandleFriendRequestRequest struct {
	Action string `json:"action" binding:"required,oneof=accept reject" example:"accept"`
}

// UpdateRemarkRequest is the request for updating remark
type UpdateRemarkRequest struct {
	Remark string `json:"remark" binding:"max=50" example:"old friend"`
}

// AddToBlacklistRequest is the request for adding to blacklist
type AddToBlacklistRequest struct {
	UserId string `json:"user_id" binding:"required" example:"user-456"`
}

// FriendResponse is the friend info response
type FriendResponse struct {
	UserID    string    `json:"user_id" example:"user-123"`
	Remark    string    `json:"remark" example:"old friend"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	UserInfo  *UserInfo `json:"user_info,omitempty"`
}

// FriendListResponse is the friend list response
type FriendListResponse struct {
	Friends []*FriendResponse `json:"friends"`
	Total   int64             `json:"total" example:"10"`
}

// FriendRequestResponse is the friend request response
type FriendRequestResponse struct {
	ID           int64     `json:"id" example:"1"`
	FromUserID   string    `json:"from_user_id" example:"user-123"`
	ToUserID     string    `json:"to_user_id" example:"user-456"`
	Message      string    `json:"message" example:"Hello"`
	Source       string    `json:"source" example:"search"`
	Status       string    `json:"status" example:"pending"`
	CreatedAt    time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	FromUserInfo *UserInfo `json:"from_user_info,omitempty"`
}

// FriendRequestListResponse is the friend request list response
type FriendRequestListResponse struct {
	Requests []*FriendRequestResponse `json:"requests"`
	Total    int64                    `json:"total" example:"5"`
}

// SendFriendRequestResponse is the response for sending a friend request
type SendFriendRequestResponse struct {
	RequestID    int64 `json:"request_id" example:"1"`
	AutoAccepted bool  `json:"auto_accepted" example:"false"`
}

// BlacklistItemResponse is the blacklist item response
type BlacklistItemResponse struct {
	ID              int64     `json:"id" example:"1"`
	UserID          string    `json:"user_id" example:"user-123"`
	BlockedUserID   string    `json:"blocked_user_id" example:"user-456"`
	CreatedAt       time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	BlockedUserInfo *UserInfo `json:"blocked_user_info,omitempty"`
}

// BlacklistResponse is the blacklist list response
type BlacklistResponse struct {
	Items []*BlacklistItemResponse `json:"items"`
	Total int64                    `json:"total" example:"2"`
}

// UserInfo is the basic user info (retrieved from user-service)
type UserInfo struct {
	UserID   string  `json:"user_id" example:"user-123"`
	Nickname string  `json:"nickname" example:"John"`
	Avatar   string  `json:"avatar" example:"https://example.com/avatar.jpg"`
	Gender   *int32  `json:"gender,omitempty" example:"1"`
	Bio      *string `json:"bio,omitempty" example:"Personal signature"`
}
