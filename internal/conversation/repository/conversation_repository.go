package repository

import (
	"context"
	"time"

	"github.com/anychat/server/internal/conversation/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ConversationRepository 会话仓库接口
type ConversationRepository interface {
	// Upsert 创建或更新会话
	Upsert(ctx context.Context, conversation *model.Conversation) error
	// GetByID 根据会话ID获取会话
	GetByID(ctx context.Context, conversationID string) (*model.Conversation, error)
	// GetByUserAndTarget 根据用户ID和目标ID获取会话
	GetByUserAndTarget(ctx context.Context, userID, conversationType, targetID string) (*model.Conversation, error)
	// ListByUser 获取用户的会话列表
	ListByUser(ctx context.Context, userID string, limit int, updatedBefore *time.Time) ([]*model.Conversation, error)
	// Delete 删除会话
	Delete(ctx context.Context, userID, conversationID string) error
	// SetPinned 设置置顶状态
	SetPinned(ctx context.Context, userID, conversationID string, pinned bool, pinTime *time.Time) error
	// SetMuted 设置免打扰状态
	SetMuted(ctx context.Context, userID, conversationID string, muted bool) error
	// SetBurnAfterReading 设置阅后即焚时长
	SetBurnAfterReading(ctx context.Context, userID, conversationID string, duration int32) error
	// SetAutoDelete 设置自动删除时长
	SetAutoDelete(ctx context.Context, userID, conversationID string, duration int32) error
	// ClearUnread 清除未读数
	ClearUnread(ctx context.Context, userID, conversationID string) error
	// IncrUnread 增加未读数
	IncrUnread(ctx context.Context, userID, conversationID string, count int32) error
	// SumUnread 统计用户总未读数
	SumUnread(ctx context.Context, userID string) (int32, error)
	// WithTx 使用事务
	WithTx(tx *gorm.DB) ConversationRepository
}

// conversationRepositoryImpl 会话仓库实现
type conversationRepositoryImpl struct {
	db *gorm.DB
}

// NewConversationRepository 创建会话仓库
func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &conversationRepositoryImpl{db: db}
}

// Upsert 创建或更新会话（冲突时更新最后消息信息）
func (r *conversationRepositoryImpl) Upsert(ctx context.Context, conversation *model.Conversation) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}, {Name: "conversation_type"}, {Name: "target_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"last_message_id",
				"last_message_content",
				"last_message_time",
				"updated_at",
			}),
		}).
		Create(conversation).Error
}

// GetByID 根据会话ID获取会话
func (r *conversationRepositoryImpl) GetByID(ctx context.Context, conversationID string) (*model.Conversation, error) {
	var conversation model.Conversation
	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		First(&conversation).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

// GetByUserAndTarget 根据用户ID和目标ID获取会话
func (r *conversationRepositoryImpl) GetByUserAndTarget(ctx context.Context, userID, conversationType, targetID string) (*model.Conversation, error) {
	var conversation model.Conversation
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND conversation_type = ? AND target_id = ?", userID, conversationType, targetID).
		First(&conversation).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

// ListByUser 获取用户会话列表（按置顶+最后消息时间排序）
func (r *conversationRepositoryImpl) ListByUser(ctx context.Context, userID string, limit int, updatedBefore *time.Time) ([]*model.Conversation, error) {
	q := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if updatedBefore != nil {
		q = q.Where("updated_at < ?", updatedBefore)
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	var conversations []*model.Conversation
	err := q.Order("is_pinned DESC, COALESCE(last_message_time, created_at) DESC").
		Limit(limit).
		Find(&conversations).Error
	return conversations, err
}

// Delete 删除会话（仅删除属于该用户的会话）
func (r *conversationRepositoryImpl) Delete(ctx context.Context, userID, conversationID string) error {
	return r.db.WithContext(ctx).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Delete(&model.Conversation{}).Error
}

// SetPinned 设置置顶状态
func (r *conversationRepositoryImpl) SetPinned(ctx context.Context, userID, conversationID string, pinned bool, pinTime *time.Time) error {
	updates := map[string]interface{}{
		"is_pinned": pinned,
		"pin_time":  pinTime,
	}
	return r.db.WithContext(ctx).Model(&model.Conversation{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Updates(updates).Error
}

// SetMuted 设置免打扰状态
func (r *conversationRepositoryImpl) SetMuted(ctx context.Context, userID, conversationID string, muted bool) error {
	return r.db.WithContext(ctx).Model(&model.Conversation{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Update("is_muted", muted).Error
}

// SetBurnAfterReading 设置阅后即焚时长
func (r *conversationRepositoryImpl) SetBurnAfterReading(ctx context.Context, userID, conversationID string, duration int32) error {
	return r.db.WithContext(ctx).Model(&model.Conversation{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Update("burn_after_reading", duration).Error
}

// SetAutoDelete 设置自动删除时长
func (r *conversationRepositoryImpl) SetAutoDelete(ctx context.Context, userID, conversationID string, duration int32) error {
	return r.db.WithContext(ctx).Model(&model.Conversation{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Update("auto_delete_duration", duration).Error
}

// ClearUnread 清除未读数
func (r *conversationRepositoryImpl) ClearUnread(ctx context.Context, userID, conversationID string) error {
	return r.db.WithContext(ctx).Model(&model.Conversation{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Update("unread_count", 0).Error
}

// IncrUnread 增加未读数
func (r *conversationRepositoryImpl) IncrUnread(ctx context.Context, userID, conversationID string, count int32) error {
	return r.db.WithContext(ctx).Model(&model.Conversation{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		UpdateColumn("unread_count", gorm.Expr("unread_count + ?", count)).Error
}

// SumUnread 统计用户所有未读数之和（免打扰会话不计入）
func (r *conversationRepositoryImpl) SumUnread(ctx context.Context, userID string) (int32, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&model.Conversation{}).
		Where("user_id = ? AND is_muted = false", userID).
		Select("COALESCE(SUM(unread_count), 0)").
		Scan(&total).Error
	return int32(total), err
}

// WithTx 返回使用事务的仓库实例
func (r *conversationRepositoryImpl) WithTx(tx *gorm.DB) ConversationRepository {
	return &conversationRepositoryImpl{db: tx}
}
