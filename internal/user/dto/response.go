package dto

import "time"

// UserProfileResponse user profile response
type UserProfileResponse struct {
	UserID    string     `json:"user_id"`
	Nickname  string     `json:"nickname"`
	Avatar    string     `json:"avatar"`
	Signature string     `json:"signature"`
	Gender    int        `json:"gender"`
	Birthday  *time.Time `json:"birthday,omitempty"`
	Region    string     `json:"region"`
	Phone     *string    `json:"phone,omitempty"`
	Email     *string    `json:"email,omitempty"`
	QRCodeURL string     `json:"qrcode_url"`
	CreatedAt time.Time  `json:"created_at"`
}

// UserInfoResponse user info response (when querying other users)
type UserInfoResponse struct {
	UserID    string `json:"user_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Signature string `json:"signature"`
	Gender    int    `json:"gender"`
	Region    string `json:"region"`
	IsFriend  bool   `json:"is_friend"`
	IsBlocked bool   `json:"is_blocked"`
}

// UserSettingsResponse user settings response
type UserSettingsResponse struct {
	UserID                string `json:"user_id"`
	NotificationEnabled   bool   `json:"notification_enabled"`
	SoundEnabled          bool   `json:"sound_enabled"`
	VibrationEnabled      bool   `json:"vibration_enabled"`
	MessagePreviewEnabled bool   `json:"message_preview_enabled"`
	FriendVerifyRequired  bool   `json:"friend_verify_required"`
	SearchByPhone         bool   `json:"search_by_phone"`
	SearchByID            bool   `json:"search_by_id"`
	Language              string `json:"language"`
}

// QRCodeResponse QR code response
type QRCodeResponse struct {
	QRCodeURL string    `json:"qrcode_url"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SearchUsersResponse search users response
type SearchUsersResponse struct {
	Total int64            `json:"total"`
	Users []*UserBriefInfo `json:"users"`
}

// UserBriefInfo user brief info
type UserBriefInfo struct {
	UserID    string `json:"user_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Signature string `json:"signature"`
}

// BindPhoneResponse bind phone response
type BindPhoneResponse struct {
	PhoneNumber string `json:"phone_number"`
	IsPrimary   bool   `json:"is_primary"`
}

// ChangePhoneResponse change phone response
type ChangePhoneResponse struct {
	OldPhoneNumber string `json:"old_phone_number"`
	NewPhoneNumber string `json:"new_phone_number"`
}

// BindEmailResponse bind email response
type BindEmailResponse struct {
	Email     string `json:"email"`
	IsPrimary bool   `json:"is_primary"`
}

// ChangeEmailResponse change email response
type ChangeEmailResponse struct {
	OldEmail string `json:"old_email"`
	NewEmail string `json:"new_email"`
}
