package dto

import "time"

// ========== Request DTOs ==========

// CreateGroupRequest request for creating a group
type CreateGroupRequest struct {
	Name      string   `json:"name" binding:"required,min=1,max=100" example:"Tech Discussion Group"`
	Avatar    string   `json:"avatar" binding:"omitempty,url" example:"https://example.com/avatar.jpg"`
	MemberIDs []string `json:"member_ids" binding:"required,min=1" example:"user-123,user-456"`
}

// UpdateGroupRequest request for updating group info
type UpdateGroupRequest struct {
	Name         *string `json:"name,omitempty" binding:"omitempty,min=1,max=100" example:"New group name"`
	Avatar       *string `json:"avatar,omitempty" binding:"omitempty,url" example:"https://example.com/new-avatar.jpg"`
	Announcement *string `json:"announcement,omitempty" binding:"omitempty,max=1000" example:"Group announcement content"`
	Description  *string `json:"description,omitempty" binding:"omitempty,max=1000" example:"Group description content"`
}

// InviteMembersRequest request for inviting members
type InviteMembersRequest struct {
	UserIDs []string `json:"user_ids" binding:"required,min=1" example:"user-123,user-456"`
}

// UpdateMemberRoleRequest request for updating member role
type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin member" example:"admin"`
}

// UpdateMemberNicknameRequest request for updating group nickname
type UpdateMemberNicknameRequest struct {
	Nickname string `json:"nickname" binding:"required,max=50" example:"Group nickname"`
}

// UpdateMemberRemarkRequest request for updating group remark (only visible to self)
type UpdateMemberRemarkRequest struct {
	Remark string `json:"remark" binding:"max=20" example:"Product Discussion Group"` // Empty string clears remark
}

// TransferOwnershipRequest request for transferring group owner
type TransferOwnershipRequest struct {
	NewOwnerID string `json:"new_owner_id" binding:"required" example:"user-456"`
}

// JoinGroupRequest request for joining a group
type JoinGroupRequest struct {
	Message string `json:"message" binding:"max=200" example:"Hello, I would like to join the group"`
}

// HandleJoinRequestRequest request for handling join request
type HandleJoinRequestRequest struct {
	Accept bool `json:"accept" example:"true"`
}

// UpdateGroupSettingsRequest request for updating group settings
type UpdateGroupSettingsRequest struct {
	JoinVerify        *bool `json:"join_verify,omitempty" example:"true"`
	AllowMemberInvite *bool `json:"allow_member_invite,omitempty" example:"true"`
	AllowViewHistory  *bool `json:"allow_view_history,omitempty" example:"true"`
	AllowAddFriend    *bool `json:"allow_add_friend,omitempty" example:"true"`
	AllowMemberModify *bool `json:"allow_member_modify,omitempty" example:"false"`
}

// PinGroupMessageRequest request for pinning group message
type PinGroupMessageRequest struct {
	MessageID string `json:"message_id" binding:"required" example:"msg-123"`
	Content   string `json:"content,omitempty" example:"Message summary"`
}

// SetGroupMuteRequest request for setting group mute
type SetGroupMuteRequest struct {
	Enabled *bool `json:"enabled" binding:"required" example:"true"`
}

// MuteMemberRequest request for muting a member
type MuteMemberRequest struct {
	Type            string `json:"type" binding:"required,oneof=permanent temporary" example:"temporary"`
	DurationMinutes int32  `json:"duration_minutes,omitempty" example:"60"`
}

// ========== Response DTOs ==========

// GroupResponse response for group info
type GroupResponse struct {
	GroupID      string    `json:"group_id" example:"group-123"`
	Name         string    `json:"name" example:"Tech Discussion Group"`
	DisplayName  string    `json:"display_name" example:"Product Discussion Group"` // Remark takes priority, falls back to Name
	Avatar       string    `json:"avatar" example:"https://example.com/avatar.jpg"`
	Announcement string    `json:"announcement" example:"Welcome to join"`
	Description  string    `json:"description" example:"Group description"`
	OwnerID      string    `json:"owner_id" example:"user-123"`
	MemberCount  int32     `json:"member_count" example:"10"`
	MaxMembers   int32     `json:"max_members" example:"500"`
	JoinVerify   bool      `json:"join_verify" example:"true"`
	IsMuted      bool      `json:"is_muted" example:"false"`
	MyRole       string    `json:"my_role,omitempty" example:"member"`
	CreatedAt    time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// GroupListResponse response for group list
type GroupListResponse struct {
	Groups     []*GroupResponse `json:"groups"`
	Total      int64            `json:"total" example:"5"`
	UpdateTime int64            `json:"update_time,omitempty" example:"1609459200"`
}

// GroupMemberResponse response for group member
type GroupMemberResponse struct {
	UserID        string     `json:"user_id" example:"user-123"`
	GroupNickname *string    `json:"group_nickname,omitempty" example:"Group nickname"`
	Role          string     `json:"role" example:"member"`
	IsMuted       bool       `json:"is_muted" example:"false"`
	MutedUntil    *time.Time `json:"muted_until,omitempty" example:"2024-01-01T01:00:00Z"`
	JoinedAt      time.Time  `json:"joined_at" example:"2024-01-01T00:00:00Z"`
	UserInfo      *UserInfo  `json:"user_info,omitempty"`
}

// GroupMemberListResponse response for group member list
type GroupMemberListResponse struct {
	Members []*GroupMemberResponse `json:"members"`
	Total   int64                  `json:"total" example:"10"`
	Page    int                    `json:"page" example:"1"`
}

// JoinGroupResponse response for joining a group
type JoinGroupResponse struct {
	NeedVerify bool   `json:"need_verify" example:"true"`
	RequestID  *int64 `json:"request_id,omitempty" example:"123"`
	Message    string `json:"message" example:"Request submitted, waiting for approval"`
}

// JoinRequestResponse response for join request
type JoinRequestResponse struct {
	ID        int64     `json:"id" example:"123"`
	GroupID   string    `json:"group_id" example:"group-123"`
	UserID    string    `json:"user_id" example:"user-123"`
	InviterID *string   `json:"inviter_id,omitempty" example:"user-456"`
	Message   string    `json:"message" example:"Hello, I would like to join"`
	Status    string    `json:"status" example:"pending"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UserInfo  *UserInfo `json:"user_info,omitempty"`
}

// JoinRequestListResponse response for join request list
type JoinRequestListResponse struct {
	Requests []*JoinRequestResponse `json:"requests"`
	Total    int64                  `json:"total" example:"3"`
}

// GroupSettingsResponse response for group settings
type GroupSettingsResponse struct {
	GroupID           string `json:"group_id" example:"group-123"`
	JoinVerify        bool   `json:"join_verify" example:"true"`
	AllowMemberInvite bool   `json:"allow_member_invite" example:"true"`
	AllowViewHistory  bool   `json:"allow_view_history" example:"true"`
	AllowAddFriend    bool   `json:"allow_add_friend" example:"true"`
	AllowMemberModify bool   `json:"allow_member_modify" example:"false"`
}

// PinnedMessageResponse response for pinned message
type PinnedMessageResponse struct {
	MessageID   string `json:"message_id" example:"msg-123"`
	Content     string `json:"content" example:"Message summary"`
	PinnedBy    string `json:"pinned_by" example:"user-123"`
	PinnedAt    int64  `json:"pinned_at" example:"1710000000"`
	ContentType string `json:"content_type,omitempty" example:"text"`
	MessageSeq  *int64 `json:"message_seq,omitempty" example:"8123"`
}

// PinnedMessageListResponse response for pinned message list
type PinnedMessageListResponse struct {
	Total      int32                    `json:"total" example:"2"`
	Version    int64                    `json:"version" example:"1775700000"`
	TopMessage *PinnedMessageResponse   `json:"top_message,omitempty"`
	Messages   []*PinnedMessageResponse `json:"messages"`
}

// GroupQRCodeResponse response for group QR code
type GroupQRCodeResponse struct {
	Token    string `json:"token" example:"550e8400-e29b-41d4-a716-446655440000"`
	DeepLink string `json:"deep_link" example:"anychat://join/group?token=550e8400-e29b-41d4-a716-446655440000"`
	ExpireAt int64  `json:"expire_at" example:"1754006400"` // Unix timestamp in seconds
}

// GroupQRCodePreviewResponse response for QR code preview
type GroupQRCodePreviewResponse struct {
	GroupID     string `json:"group_id" example:"group-123"`
	Name        string `json:"name" example:"Product Team"`
	Avatar      string `json:"avatar" example:"https://example.com/group.png"`
	MemberCount int32  `json:"member_count" example:"42"`
	NeedVerify  bool   `json:"need_verify" example:"true"`
}

// JoinGroupByQRCodeResponse response for joining via QR code
type JoinGroupByQRCodeResponse struct {
	Joined     bool   `json:"joined" example:"true"` // true=direct join, false=request submitted
	GroupID    string `json:"group_id" example:"group-123"`
	NeedVerify bool   `json:"need_verify" example:"false"` // whether approval is needed
	RequestID  *int64 `json:"request_id,omitempty" example:"123"`
}

// ========== Shared Types ==========

// UserInfo basic user info (from user-service)
type UserInfo struct {
	UserID   string  `json:"user_id" example:"user-123"`
	Nickname string  `json:"nickname" example:"Zhang San"`
	Avatar   string  `json:"avatar" example:"https://example.com/avatar.jpg"`
	Gender   *int32  `json:"gender,omitempty" example:"1"`
	Bio      *string `json:"bio,omitempty" example:"Personal signature"`
}
