package service

import (
	"context"
	"fmt"
	"time"

	conversationpb "github.com/anychat/server/api/proto/conversation"
	"github.com/anychat/server/internal/conversation/model"
	"github.com/anychat/server/internal/conversation/repository"
	"github.com/anychat/server/pkg/logger"
	"github.com/anychat/server/pkg/notification"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// ConversationService 会话服务接口
type ConversationService interface {
	GetConversations(ctx context.Context, req *conversationpb.GetConversationsRequest) (*conversationpb.GetConversationsResponse, error)
	GetConversation(ctx context.Context, userID, conversationID string) (*conversationpb.Conversation, error)
	CreateOrUpdateConversation(ctx context.Context, req *conversationpb.CreateOrUpdateConversationRequest) (*conversationpb.Conversation, error)
	DeleteConversation(ctx context.Context, userID, conversationID string) error
	SetPinned(ctx context.Context, userID, conversationID string, pinned bool) error
	SetMuted(ctx context.Context, userID, conversationID string, muted bool) error
	SetBurnAfterReading(ctx context.Context, userID, conversationID string, duration int32) error
	SetAutoDelete(ctx context.Context, userID, conversationID string, duration int32) error
	ClearUnread(ctx context.Context, userID, conversationID string) error
	GetTotalUnread(ctx context.Context, userID string) (int32, error)
	IncrUnread(ctx context.Context, userID, conversationID string, count int32) error
}

// conversationServiceImpl 会话服务实现
type conversationServiceImpl struct {
	conversationRepo repository.ConversationRepository
	notificationPub  notification.Publisher
}

// NewConversationService 创建会话服务
func NewConversationService(
	conversationRepo repository.ConversationRepository,
	notificationPub notification.Publisher,
) ConversationService {
	return &conversationServiceImpl{
		conversationRepo: conversationRepo,
		notificationPub:  notificationPub,
	}
}

// GetConversations 获取用户会话列表
func (s *conversationServiceImpl) GetConversations(ctx context.Context, req *conversationpb.GetConversationsRequest) (*conversationpb.GetConversationsResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 20
	}

	var updatedBefore *time.Time
	if req.UpdatedBefore != nil {
		t := time.Unix(*req.UpdatedBefore, 0)
		updatedBefore = &t
	}

	conversations, err := s.conversationRepo.ListByUser(ctx, req.UserId, limit, updatedBefore)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}

	pbConversations := make([]*conversationpb.Conversation, 0, len(conversations))
	for _, c := range conversations {
		pbConversations = append(pbConversations, toProtoConversation(c))
	}

	return &conversationpb.GetConversationsResponse{
		Conversations: pbConversations,
		HasMore:       len(conversations) == limit,
	}, nil
}

// GetConversation 获取单个会话
func (s *conversationServiceImpl) GetConversation(ctx context.Context, userID, conversationID string) (*conversationpb.Conversation, error) {
	conversation, err := s.conversationRepo.GetByID(ctx, conversationID)
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("conversation not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	if conversation.UserID != userID {
		return nil, fmt.Errorf("conversation not found")
	}
	return toProtoConversation(conversation), nil
}

// CreateOrUpdateConversation 创建或更新会话（消息到达时调用）
func (s *conversationServiceImpl) CreateOrUpdateConversation(ctx context.Context, req *conversationpb.CreateOrUpdateConversationRequest) (*conversationpb.Conversation, error) {
	// 先尝试查找已有会话
	existing, err := s.conversationRepo.GetByUserAndTarget(ctx, req.UserId, req.ConversationType, req.TargetId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing conversation: %w", err)
	}

	var msgTime *time.Time
	if req.LastMessageTimestamp > 0 {
		t := time.Unix(req.LastMessageTimestamp, 0)
		msgTime = &t
	}

	if existing != nil {
		// 更新已有会话的最后消息信息
		existing.LastMessageID = req.LastMessageId
		existing.LastMessageContent = req.LastMessageContent
		existing.LastMessageTime = msgTime
		if err := s.conversationRepo.Upsert(ctx, existing); err != nil {
			return nil, fmt.Errorf("failed to update conversation: %w", err)
		}
		return toProtoConversation(existing), nil
	}

	// 创建新会话
	conversation := &model.Conversation{
		ConversationID:     uuid.New().String(),
		ConversationType:   req.ConversationType,
		UserID:             req.UserId,
		TargetID:           req.TargetId,
		LastMessageID:      req.LastMessageId,
		LastMessageContent: req.LastMessageContent,
		LastMessageTime:    msgTime,
	}

	if err := s.conversationRepo.Upsert(ctx, conversation); err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return toProtoConversation(conversation), nil
}

// DeleteConversation 删除会话并发送通知
func (s *conversationServiceImpl) DeleteConversation(ctx context.Context, userID, conversationID string) error {
	if err := s.conversationRepo.Delete(ctx, userID, conversationID); err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	// 发布会话删除通知（多端同步）
	notif := notification.NewNotification(notification.TypeConversationDeleted, userID, notification.PriorityNormal).
		AddPayloadField("conversation_id", conversationID)
	if err := s.notificationPub.PublishToUser(userID, notif); err != nil {
		logger.Warn("Failed to publish conversation deleted notification",
			zap.String("userID", userID),
			zap.Error(err))
	}

	return nil
}

// SetPinned 设置置顶状态并发送通知
func (s *conversationServiceImpl) SetPinned(ctx context.Context, userID, conversationID string, pinned bool) error {
	var pinTime *time.Time
	if pinned {
		t := time.Now()
		pinTime = &t
	}

	if err := s.conversationRepo.SetPinned(ctx, userID, conversationID, pinned, pinTime); err != nil {
		return fmt.Errorf("failed to set pinned: %w", err)
	}

	// 发布置顶状态同步通知（多端同步）
	notif := notification.NewNotification(notification.TypeConversationPinUpdated, userID, notification.PriorityNormal).
		AddPayloadField("conversation_id", conversationID).
		AddPayloadField("is_pinned", pinned)
	if err := s.notificationPub.PublishToUser(userID, notif); err != nil {
		logger.Warn("Failed to publish conversation pin notification",
			zap.String("userID", userID),
			zap.Error(err))
	}

	return nil
}

// SetMuted 设置免打扰状态并发送通知
func (s *conversationServiceImpl) SetMuted(ctx context.Context, userID, conversationID string, muted bool) error {
	if err := s.conversationRepo.SetMuted(ctx, userID, conversationID, muted); err != nil {
		return fmt.Errorf("failed to set muted: %w", err)
	}

	// 发布免打扰设置同步通知（多端同步）
	notif := notification.NewNotification(notification.TypeConversationMuteUpdated, userID, notification.PriorityNormal).
		AddPayloadField("conversation_id", conversationID).
		AddPayloadField("is_muted", muted)
	if err := s.notificationPub.PublishToUser(userID, notif); err != nil {
		logger.Warn("Failed to publish conversation mute notification",
			zap.String("userID", userID),
			zap.Error(err))
	}

	return nil
}

// SetBurnAfterReading 设置阅后即焚时长并发送通知
func (s *conversationServiceImpl) SetBurnAfterReading(ctx context.Context, userID, conversationID string, duration int32) error {
	if err := s.conversationRepo.SetBurnAfterReading(ctx, userID, conversationID, duration); err != nil {
		return fmt.Errorf("failed to set burn after reading: %w", err)
	}

	// 发布阅后即焚配置变更通知（多端同步）
	notif := notification.NewNotification(notification.TypeConversationBurnUpdated, userID, notification.PriorityNormal).
		AddPayloadField("conversation_id", conversationID).
		AddPayloadField("burn_after_reading", duration)
	if err := s.notificationPub.PublishToUser(userID, notif); err != nil {
		logger.Warn("Failed to publish conversation burn notification",
			zap.String("userID", userID),
			zap.Error(err))
	}

	return nil
}

// SetAutoDelete 设置自动删除时长并发送通知
func (s *conversationServiceImpl) SetAutoDelete(ctx context.Context, userID, conversationID string, duration int32) error {
	if err := s.conversationRepo.SetAutoDelete(ctx, userID, conversationID, duration); err != nil {
		return fmt.Errorf("failed to set auto delete: %w", err)
	}

	// 发布自动删除配置变更通知（多端同步）
	notif := notification.NewNotification(notification.TypeConversationAutoDeleteUpdated, userID, notification.PriorityNormal).
		AddPayloadField("conversation_id", conversationID).
		AddPayloadField("auto_delete_duration", duration)
	if err := s.notificationPub.PublishToUser(userID, notif); err != nil {
		logger.Warn("Failed to publish conversation auto delete notification",
			zap.String("userID", userID),
			zap.Error(err))
	}

	return nil
}

// ClearUnread 清除未读数并发送通知
func (s *conversationServiceImpl) ClearUnread(ctx context.Context, userID, conversationID string) error {
	if err := s.conversationRepo.ClearUnread(ctx, userID, conversationID); err != nil {
		return fmt.Errorf("failed to clear unread: %w", err)
	}

	// 获取最新总未读数
	total, err := s.conversationRepo.SumUnread(ctx, userID)
	if err != nil {
		logger.Warn("Failed to get total unread after clear", zap.String("userID", userID), zap.Error(err))
	}

	// 发布未读数更新通知（多端同步）
	notif := notification.NewNotification(notification.TypeConversationUnreadUpdated, userID, notification.PriorityNormal).
		AddPayloadField("conversation_id", conversationID).
		AddPayloadField("unread_count", 0).
		AddPayloadField("total_unread", total)
	if err := s.notificationPub.PublishToUser(userID, notif); err != nil {
		logger.Warn("Failed to publish unread notification",
			zap.String("userID", userID),
			zap.Error(err))
	}

	return nil
}

// GetTotalUnread 获取用户总未读数
func (s *conversationServiceImpl) GetTotalUnread(ctx context.Context, userID string) (int32, error) {
	total, err := s.conversationRepo.SumUnread(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get total unread: %w", err)
	}
	return total, nil
}

// IncrUnread 增加未读数并发送通知
func (s *conversationServiceImpl) IncrUnread(ctx context.Context, userID, conversationID string, count int32) error {
	if err := s.conversationRepo.IncrUnread(ctx, userID, conversationID, count); err != nil {
		return fmt.Errorf("failed to incr conversation unread: %w", err)
	}

	// 发布未读数更新通知
	notif := notification.NewNotification(notification.TypeConversationUnreadUpdated, userID, notification.PriorityNormal).
		AddPayloadField("conversation_id", conversationID)
	if err := s.notificationPub.PublishToUser(userID, notif); err != nil {
		logger.Warn("Failed to publish unread incr notification",
			zap.String("userID", userID),
			zap.Error(err))
	}

	return nil
}

// toProtoConversation 将model.Conversation转换为protobuf Conversation
func toProtoConversation(s *model.Conversation) *conversationpb.Conversation {
	pb := &conversationpb.Conversation{
		ConversationId:     s.ConversationID,
		ConversationType:   s.ConversationType,
		UserId:             s.UserID,
		TargetId:           s.TargetID,
		LastMessageId:      s.LastMessageID,
		LastMessageContent: s.LastMessageContent,
		UnreadCount:        s.UnreadCount,
		IsPinned:           s.IsPinned,
		IsMuted:            s.IsMuted,
		BurnAfterReading:   s.BurnAfterReading,
		AutoDeleteDuration: s.AutoDeleteDuration,
		CreatedAt:          timestamppb.New(s.CreatedAt),
		UpdatedAt:          timestamppb.New(s.UpdatedAt),
	}
	if s.LastMessageTime != nil {
		pb.LastMessageTime = timestamppb.New(*s.LastMessageTime)
	}
	if s.PinTime != nil {
		pb.PinTime = timestamppb.New(*s.PinTime)
	}
	return pb
}
