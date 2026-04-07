package repository

import (
	"context"
	"time"

	"github.com/anychat/server/internal/group/model"
	"gorm.io/gorm"
)

type GroupPinnedMessageRepository interface {
	Upsert(ctx context.Context, pinned *model.GroupPinnedMessage) error
	Delete(ctx context.Context, groupID, messageID string) error
	ListByGroup(ctx context.Context, groupID string) ([]*model.GroupPinnedMessage, error)
	WithTx(tx *gorm.DB) GroupPinnedMessageRepository
}

type groupPinnedMessageRepositoryImpl struct {
	db *gorm.DB
}

func NewGroupPinnedMessageRepository(db *gorm.DB) GroupPinnedMessageRepository {
	return &groupPinnedMessageRepositoryImpl{db: db}
}

func (r *groupPinnedMessageRepositoryImpl) Upsert(ctx context.Context, pinned *model.GroupPinnedMessage) error {
	updates := map[string]any{
		"pinned_by":  pinned.PinnedBy,
		"content":    pinned.Content,
		"created_at": time.Now(),
	}
	return r.db.WithContext(ctx).
		Model(&model.GroupPinnedMessage{}).
		Where("group_id = ? AND message_id = ?", pinned.GroupID, pinned.MessageID).
		Assign(updates).
		FirstOrCreate(pinned).Error
}

func (r *groupPinnedMessageRepositoryImpl) Delete(ctx context.Context, groupID, messageID string) error {
	return r.db.WithContext(ctx).
		Where("group_id = ? AND message_id = ?", groupID, messageID).
		Delete(&model.GroupPinnedMessage{}).Error
}

func (r *groupPinnedMessageRepositoryImpl) ListByGroup(ctx context.Context, groupID string) ([]*model.GroupPinnedMessage, error) {
	var records []*model.GroupPinnedMessage
	err := r.db.WithContext(ctx).
		Where("group_id = ?", groupID).
		Order("created_at DESC").
		Find(&records).Error
	return records, err
}

func (r *groupPinnedMessageRepositoryImpl) WithTx(tx *gorm.DB) GroupPinnedMessageRepository {
	return &groupPinnedMessageRepositoryImpl{db: tx}
}
