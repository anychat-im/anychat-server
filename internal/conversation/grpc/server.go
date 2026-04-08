package grpc

import (
	"context"

	commonpb "github.com/anychat/server/api/proto/common"
	conversationpb "github.com/anychat/server/api/proto/conversation"
	"github.com/anychat/server/internal/conversation/service"
	"github.com/anychat/server/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server Conversation gRPC服务器
type Server struct {
	conversationpb.UnimplementedConversationServiceServer
	conversationService service.ConversationService
}

// NewServer 创建gRPC服务器
func NewServer(conversationService service.ConversationService) *Server {
	return &Server{conversationService: conversationService}
}

// GetConversations 获取用户会话列表
func (s *Server) GetConversations(ctx context.Context, req *conversationpb.GetConversationsRequest) (*conversationpb.GetConversationsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	resp, err := s.conversationService.GetConversations(ctx, req)
	if err != nil {
		logger.Error("GetConversations failed", zap.String("userID", req.UserId), zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

// GetConversation 获取单个会话
func (s *Server) GetConversation(ctx context.Context, req *conversationpb.GetConversationRequest) (*conversationpb.Conversation, error) {
	if req.UserId == "" || req.ConversationId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and conversation_id are required")
	}
	conversation, err := s.conversationService.GetConversation(ctx, req.UserId, req.ConversationId)
	if err != nil {
		logger.Error("GetConversation failed", zap.String("conversationID", req.ConversationId), zap.Error(err))
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return conversation, nil
}

// CreateOrUpdateConversation 创建或更新会话
func (s *Server) CreateOrUpdateConversation(ctx context.Context, req *conversationpb.CreateOrUpdateConversationRequest) (*conversationpb.Conversation, error) {
	if req.UserId == "" || req.TargetId == "" || req.ConversationType == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id, target_id and conversation_type are required")
	}
	conversation, err := s.conversationService.CreateOrUpdateConversation(ctx, req)
	if err != nil {
		logger.Error("CreateOrUpdateConversation failed", zap.String("userID", req.UserId), zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return conversation, nil
}

// DeleteConversation 删除会话
func (s *Server) DeleteConversation(ctx context.Context, req *conversationpb.DeleteConversationRequest) (*commonpb.Empty, error) {
	if req.UserId == "" || req.ConversationId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and conversation_id are required")
	}
	if err := s.conversationService.DeleteConversation(ctx, req.UserId, req.ConversationId); err != nil {
		logger.Error("DeleteConversation failed", zap.String("conversationID", req.ConversationId), zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &commonpb.Empty{}, nil
}

// SetPinned 设置置顶
func (s *Server) SetPinned(ctx context.Context, req *conversationpb.SetPinnedRequest) (*commonpb.Empty, error) {
	if req.UserId == "" || req.ConversationId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and conversation_id are required")
	}
	if err := s.conversationService.SetPinned(ctx, req.UserId, req.ConversationId, req.Pinned); err != nil {
		logger.Error("SetPinned failed", zap.String("conversationID", req.ConversationId), zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &commonpb.Empty{}, nil
}

// SetMuted 设置免打扰
func (s *Server) SetMuted(ctx context.Context, req *conversationpb.SetMutedRequest) (*commonpb.Empty, error) {
	if req.UserId == "" || req.ConversationId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and conversation_id are required")
	}
	if err := s.conversationService.SetMuted(ctx, req.UserId, req.ConversationId, req.Muted); err != nil {
		logger.Error("SetMuted failed", zap.String("conversationID", req.ConversationId), zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &commonpb.Empty{}, nil
}

// ClearUnread 清除未读数
func (s *Server) ClearUnread(ctx context.Context, req *conversationpb.ClearUnreadRequest) (*commonpb.Empty, error) {
	if req.UserId == "" || req.ConversationId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and conversation_id are required")
	}
	if err := s.conversationService.ClearUnread(ctx, req.UserId, req.ConversationId); err != nil {
		logger.Error("ClearUnread failed", zap.String("conversationID", req.ConversationId), zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &commonpb.Empty{}, nil
}

// GetTotalUnread 获取总未读数
func (s *Server) GetTotalUnread(ctx context.Context, req *conversationpb.GetTotalUnreadRequest) (*conversationpb.GetTotalUnreadResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	total, err := s.conversationService.GetTotalUnread(ctx, req.UserId)
	if err != nil {
		logger.Error("GetTotalUnread failed", zap.String("userID", req.UserId), zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &conversationpb.GetTotalUnreadResponse{TotalUnread: total}, nil
}

// IncrUnread 增加未读数
func (s *Server) IncrUnread(ctx context.Context, req *conversationpb.IncrUnreadRequest) (*commonpb.Empty, error) {
	if req.UserId == "" || req.ConversationId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and conversation_id are required")
	}
	if err := s.conversationService.IncrUnread(ctx, req.UserId, req.ConversationId, req.Count); err != nil {
		logger.Error("IncrUnread failed", zap.String("conversationID", req.ConversationId), zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &commonpb.Empty{}, nil
}

// SetBurnAfterReading 设置阅后即焚
func (s *Server) SetBurnAfterReading(ctx context.Context, req *conversationpb.SetBurnAfterReadingRequest) (*commonpb.Empty, error) {
	if req.UserId == "" || req.ConversationId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and conversation_id are required")
	}
	if err := s.conversationService.SetBurnAfterReading(ctx, req.UserId, req.ConversationId, req.Duration); err != nil {
		logger.Error("SetBurnAfterReading failed", zap.String("conversationID", req.ConversationId), zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &commonpb.Empty{}, nil
}

// SetAutoDelete 设置自动删除
func (s *Server) SetAutoDelete(ctx context.Context, req *conversationpb.SetAutoDeleteRequest) (*commonpb.Empty, error) {
	if req.UserId == "" || req.ConversationId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and conversation_id are required")
	}
	if err := s.conversationService.SetAutoDelete(ctx, req.UserId, req.ConversationId, req.Duration); err != nil {
		logger.Error("SetAutoDelete failed", zap.String("conversationID", req.ConversationId), zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &commonpb.Empty{}, nil
}
