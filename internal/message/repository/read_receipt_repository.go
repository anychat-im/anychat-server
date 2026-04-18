package repository

import (
	"context"

	"github.com/anychat/server/internal/message/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ReadReceiptRepository read receipt repository interface
type ReadReceiptRepository interface {
	Upsert(ctx context.Context, receipt *model.MessageReadReceipt) error
	GetByConversationAndUser(ctx context.Context, conversationID, userID string) (*model.MessageReadReceipt, error)
	GetByConversation(ctx context.Context, conversationID string) ([]*model.MessageReadReceipt, error)
	GetByUser(ctx context.Context, userID string) ([]*model.MessageReadReceipt, error)
	Delete(ctx context.Context, conversationID, userID string) error
	WithTx(tx *gorm.DB) ReadReceiptRepository
}

// readReceiptRepositoryImpl read receipt repository implementation
type readReceiptRepositoryImpl struct {
	db *gorm.DB
}

// NewReadReceiptRepository creates read receipt repository
func NewReadReceiptRepository(db *gorm.DB) ReadReceiptRepository {
	return &readReceiptRepositoryImpl{db: db}
}

// Upsert creates or updates read receipt
func (r *readReceiptRepositoryImpl) Upsert(ctx context.Context, receipt *model.MessageReadReceipt) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "conversation_id"}, {Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"conversation_type",
				"target_id",
				"last_read_seq",
				"last_read_message_id",
				"read_at",
			}),
		}).
		Create(receipt).Error
}

// GetByConversationAndUser retrieves read receipt by conversation and user
func (r *readReceiptRepositoryImpl) GetByConversationAndUser(ctx context.Context, conversationID, userID string) (*model.MessageReadReceipt, error) {
	var receipt model.MessageReadReceipt
	err := r.db.WithContext(ctx).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		First(&receipt).Error
	if err != nil {
		return nil, err
	}
	return &receipt, nil
}

// GetByConversation retrieves all read receipts for a conversation
func (r *readReceiptRepositoryImpl) GetByConversation(ctx context.Context, conversationID string) ([]*model.MessageReadReceipt, error) {
	var receipts []*model.MessageReadReceipt
	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("read_at DESC").
		Find(&receipts).Error
	return receipts, err
}

// GetByUser retrieves all read receipts for a user
func (r *readReceiptRepositoryImpl) GetByUser(ctx context.Context, userID string) ([]*model.MessageReadReceipt, error) {
	var receipts []*model.MessageReadReceipt
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("read_at DESC").
		Find(&receipts).Error
	return receipts, err
}

// Delete deletes read receipt
func (r *readReceiptRepositoryImpl) Delete(ctx context.Context, conversationID, userID string) error {
	return r.db.WithContext(ctx).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Delete(&model.MessageReadReceipt{}).Error
}

// WithTx uses transaction
func (r *readReceiptRepositoryImpl) WithTx(tx *gorm.DB) ReadReceiptRepository {
	return &readReceiptRepositoryImpl{db: tx}
}
