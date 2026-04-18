package dto

import (
	"time"

	"github.com/anychat/server/internal/user/model"
)

// UpdateProfileRequest update profile request
type UpdateProfileRequest struct {
	Nickname  *string    `json:"nickname"`
	Avatar    *string    `json:"avatar"`
	Signature *string    `json:"signature"`
	Gender    *int       `json:"gender"`
	Birthday  *time.Time `json:"birthday"`
	Region    *string    `json:"region"`
}

// UpdateSettingsRequest update settings request
type UpdateSettingsRequest struct {
	NotificationEnabled   *bool   `json:"notification_enabled"`
	SoundEnabled          *bool   `json:"sound_enabled"`
	VibrationEnabled      *bool   `json:"vibration_enabled"`
	MessagePreviewEnabled *bool   `json:"message_preview_enabled"`
	FriendVerifyRequired  *bool   `json:"friend_verify_required"`
	SearchByPhone         *bool   `json:"search_by_phone"`
	SearchByID            *bool   `json:"search_by_id"`
	Language              *string `json:"language"`
}

// UpdatePushTokenRequest update push token request
type UpdatePushTokenRequest struct {
	DeviceID  string `json:"device_id" binding:"required"`
	PushToken string `json:"push_token" binding:"required"`
	Platform  model.PushPlatform `json:"platform" binding:"required,oneof=1 2"` // 1-iOS/2-Android
}

// BindPhoneRequest bind phone request
type BindPhoneRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	VerifyCode  string `json:"verify_code" binding:"required"`
}

// ChangePhoneRequest change phone request
type ChangePhoneRequest struct {
	OldPhoneNumber string  `json:"old_phone_number" binding:"required"`
	NewPhoneNumber string  `json:"new_phone_number" binding:"required"`
	NewVerifyCode  string  `json:"new_verify_code" binding:"required"`
	OldVerifyCode  *string `json:"old_verify_code"`
	DeviceID       string  `json:"-"`
}

// BindEmailRequest bind email request
type BindEmailRequest struct {
	Email      string `json:"email" binding:"required"`
	VerifyCode string `json:"verify_code" binding:"required"`
}

// ChangeEmailRequest change email request
type ChangeEmailRequest struct {
	OldEmail      string  `json:"old_email" binding:"required"`
	NewEmail      string  `json:"new_email" binding:"required"`
	NewVerifyCode string  `json:"new_verify_code" binding:"required"`
	OldVerifyCode *string `json:"old_verify_code"`
	DeviceID      string  `json:"-"`
}

// SearchUsersRequest search users request
type SearchUsersRequest struct {
	Keyword  string `form:"keyword" binding:"required"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}
