package model

import "time"

// GroupPinnedMessage 群置顶消息
type GroupPinnedMessage struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	GroupID   string    `gorm:"column:group_id;not null;index:idx_group_pinned_messages_group_id" json:"groupId"`
	MessageID string    `gorm:"column:message_id;not null" json:"messageId"`
	PinnedBy  string    `gorm:"column:pinned_by;not null" json:"pinnedBy"`
	Content   string    `gorm:"column:content;type:text" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

func (GroupPinnedMessage) TableName() string {
	return "group_pinned_messages"
}
