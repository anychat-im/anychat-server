package dto

import "github.com/anychat/server/internal/auth/model"

// SendVerificationCodeRequest send verification code request
type SendVerificationCodeRequest struct {
	Target     string                       `json:"target" binding:"required"`
	TargetType model.VerificationTargetType `json:"target_type" binding:"required"`
	Purpose    model.VerificationPurpose    `json:"purpose" binding:"required"`
	DeviceID   string                       `json:"device_id"`
	IPAddress  string                       `json:"ip_address"`
}

// SendVerificationCodeResponse send verification code response
type SendVerificationCodeResponse struct {
	CodeID    string `json:"code_id"`
	ExpiresIn int64  `json:"expires_in"`
}

// RegisterRequest register request
type RegisterRequest struct {
	PhoneNumber   string           `json:"phone_number" binding:"required_without=Email"`
	Email         string           `json:"email" binding:"required_without=PhoneNumber,omitempty,email"`
	Password      string           `json:"password" binding:"required,min=8,max=32"`
	VerifyCode    string           `json:"verify_code" binding:"required"`
	Nickname      string           `json:"nickname"`
	DeviceType    model.DeviceType `json:"device_type" binding:"required"`
	DeviceID      string           `json:"device_id" binding:"required"`
	ClientVersion string           `json:"client_version" binding:"required"`
}

// RegisterResponse register response
type RegisterResponse struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
}

// LoginRequest login request
type LoginRequest struct {
	Account       string           `json:"account" binding:"required"`
	Password      string           `json:"password" binding:"required"`
	DeviceType    model.DeviceType `json:"device_type" binding:"required"`
	DeviceID      string           `json:"device_id" binding:"required"`
	ClientVersion string           `json:"client_version" binding:"required"`
	IpAddress     string           `json:"ip_address"`
}

// LoginResponse login response
type LoginResponse struct {
	UserID       string    `json:"user_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"` // seconds
	User         *UserInfo `json:"user"`
}

// UserInfo user info
type UserInfo struct {
	UserID   string  `json:"user_id"`
	Nickname string  `json:"nickname"`
	Avatar   string  `json:"avatar"`
	Phone    *string `json:"phone,omitempty"`
	Email    *string `json:"email,omitempty"`
}

// RefreshTokenRequest refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse refresh token response
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
}

// ChangePasswordRequest change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=32"`
	DeviceID    string `json:"device_id" binding:"required"`
}

// ResetPasswordRequest reset password request
type ResetPasswordRequest struct {
	Account     string `json:"account" binding:"required"`
	VerifyCode  string `json:"verify_code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=32"`
}

// LogoutRequest logout request
type LogoutRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
}

// TokenInfo token info
type TokenInfo struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
