package dto

import (
	"time"

	"github.com/anychat/server/internal/auth/model"
)

type SendCodeRequest struct {
	Target     string                       `json:"target" binding:"required"`
	TargetType model.VerificationTargetType `json:"target_type" binding:"required"`
	Purpose    model.VerificationPurpose    `json:"purpose" binding:"required"`
	DeviceID   string                       `json:"device_id"`
}

type SendCodeResponse struct {
	CodeID    string `json:"code_id"`
	ExpiresIn int64  `json:"expires_in"`
	Sent      bool   `json:"sent"`
	Message   string `json:"message"`
}

type VerifyCodeRequest struct {
	Target     string                       `json:"target" binding:"required"`
	TargetType model.VerificationTargetType `json:"target_type" binding:"required"`
	Purpose    model.VerificationPurpose    `json:"purpose" binding:"required"`
	Code       string                       `json:"code" binding:"required,len=6"`
}

type VerifyCodeResponse struct {
	Valid   bool   `json:"valid"`
	CodeID  string `json:"code_id"`
	Message string `json:"message"`
}

type CheckCodeStatusResponse struct {
	Status    model.VerificationCodeStatus `json:"status"`
	ExpiresAt time.Time                    `json:"expires_at"`
}
