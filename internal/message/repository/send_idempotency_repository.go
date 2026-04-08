package repository

import (
	"context"

	"github.com/anychat/server/internal/message/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SendIdempotencyRepository 发送幂等仓库接口
type SendIdempotencyRepository interface {
	CreateIfNotExists(ctx context.Context, rec *model.MessageSendIdempotency) error
	GetForUpdateByKey(ctx context.Context, senderID, conversationID, localID string) (*model.MessageSendIdempotency, error)
	BindMessageID(ctx context.Context, senderID, conversationID, localID, messageID string) error
	WithTx(tx *gorm.DB) SendIdempotencyRepository
}

type sendIdempotencyRepositoryImpl struct {
	db *gorm.DB
}

// NewSendIdempotencyRepository 创建发送幂等仓库
func NewSendIdempotencyRepository(db *gorm.DB) SendIdempotencyRepository {
	return &sendIdempotencyRepositoryImpl{db: db}
}

// CreateIfNotExists 创建幂等记录（存在则忽略）
func (r *sendIdempotencyRepositoryImpl) CreateIfNotExists(ctx context.Context, rec *model.MessageSendIdempotency) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(rec).Error
}

// GetForUpdateByKey 按键查询并加行锁
func (r *sendIdempotencyRepositoryImpl) GetForUpdateByKey(ctx context.Context, senderID, conversationID, localID string) (*model.MessageSendIdempotency, error) {
	var rec model.MessageSendIdempotency
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("sender_id = ? AND conversation_id = ? AND local_id = ?", senderID, conversationID, localID).
		First(&rec).Error
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

// BindMessageID 绑定实际消息ID
func (r *sendIdempotencyRepositoryImpl) BindMessageID(ctx context.Context, senderID, conversationID, localID, messageID string) error {
	return r.db.WithContext(ctx).
		Model(&model.MessageSendIdempotency{}).
		Where("sender_id = ? AND conversation_id = ? AND local_id = ?", senderID, conversationID, localID).
		Update("message_id", messageID).Error
}

// WithTx 使用事务
func (r *sendIdempotencyRepositoryImpl) WithTx(tx *gorm.DB) SendIdempotencyRepository {
	return &sendIdempotencyRepositoryImpl{db: tx}
}
