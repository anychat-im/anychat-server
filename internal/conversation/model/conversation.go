package model

import (
	"time"
)

// Conversation 会话模型
type Conversation struct {
	ConversationID     string     `gorm:"column:conversation_id;primaryKey"`
	ConversationType   string     `gorm:"column:conversation_type;not null"` // single/group/system
	UserID             string     `gorm:"column:user_id;not null"`
	TargetID           string     `gorm:"column:target_id;not null"`
	LastMessageID      string     `gorm:"column:last_message_id"`
	LastMessageContent string     `gorm:"column:last_message_content"`
	LastMessageTime    *time.Time `gorm:"column:last_message_time"`
	UnreadCount        int32      `gorm:"column:unread_count;default:0"`
	IsPinned           bool       `gorm:"column:is_pinned;default:false"`
	IsMuted            bool       `gorm:"column:is_muted;default:false"`
	PinTime            *time.Time `gorm:"column:pin_time"`
	BurnAfterReading   int32      `gorm:"column:burn_after_reading;default:0"`   // 阅后即焚时长(秒),0表示未启用
	AutoDeleteDuration int32      `gorm:"column:auto_delete_duration;default:0"` // 自动删除时长(秒),0表示未启用
	CreatedAt          time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 指定表名
func (Conversation) TableName() string {
	return "conversations"
}

// ConversationType 会话类型常量
const (
	ConversationTypeSingle = "single"
	ConversationTypeGroup  = "group"
	ConversationTypeSystem = "system"
)
