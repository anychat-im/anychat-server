package service

import (
	"context"
	"testing"
	"time"

	messagepb "github.com/anychat/server/api/proto/message"
	"github.com/anychat/server/internal/message/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupMessageServiceTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Message{}))
	return db
}

func TestAckReadTriggers_UpdatesExpireTimeWithBurnDuration(t *testing.T) {
	db := setupMessageServiceTestDB(t)
	now := time.Now()

	msg := &model.Message{
		MessageID:               "msg-burn-1",
		ConversationID:          "single_u1_u2",
		ConversationType:        model.ConversationTypeSingle,
		SenderID:                "u1",
		ContentType:             model.ContentTypeText,
		Content:                 `{"text":"hello"}`,
		Sequence:                1,
		Status:                  model.MessageStatusNormal,
		BurnAfterReadingSeconds: 30,
		CreatedAt:               now,
		UpdatedAt:               now,
	}
	require.NoError(t, db.Create(msg).Error)

	svc := &messageServiceImpl{db: db}
	resp, err := svc.AckReadTriggers(context.Background(), &messagepb.AckReadTriggersRequest{
		UserId: "u2",
		Events: []*messagepb.ReadTriggerEvent{
			{MessageId: msg.MessageID},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"msg-burn-1"}, resp.SuccessIds)
	assert.Empty(t, resp.IgnoredIds)

	var updated model.Message
	require.NoError(t, db.Where("message_id = ?", msg.MessageID).First(&updated).Error)
	require.NotNil(t, updated.BurnAfterReadingExpireTime)
	require.NotNil(t, updated.ExpireTime)
	assert.WithinDuration(t, now.Add(30*time.Second), *updated.BurnAfterReadingExpireTime, 2*time.Second)
	assert.WithinDuration(t, now.Add(30*time.Second), *updated.ExpireTime, 2*time.Second)
}

func TestAckReadTriggers_KeepEarlierExpireTimeAndIgnoreSender(t *testing.T) {
	db := setupMessageServiceTestDB(t)
	now := time.Now()
	earlierExpire := now.Add(10 * time.Second)

	msg := &model.Message{
		MessageID:               "msg-burn-2",
		ConversationID:          "single_u1_u2",
		ConversationType:        model.ConversationTypeSingle,
		SenderID:                "u1",
		ContentType:             model.ContentTypeText,
		Content:                 `{"text":"hello"}`,
		Sequence:                1,
		Status:                  model.MessageStatusNormal,
		BurnAfterReadingSeconds: 30,
		AutoDeleteExpireTime:    &earlierExpire,
		ExpireTime:              &earlierExpire,
		CreatedAt:               now,
		UpdatedAt:               now,
	}
	require.NoError(t, db.Create(msg).Error)

	svc := &messageServiceImpl{db: db}

	// 发送方上报，应该被忽略
	respSender, err := svc.AckReadTriggers(context.Background(), &messagepb.AckReadTriggersRequest{
		UserId: "u1",
		Events: []*messagepb.ReadTriggerEvent{
			{MessageId: msg.MessageID},
		},
	})
	require.NoError(t, err)
	assert.Empty(t, respSender.SuccessIds)
	assert.Equal(t, []string{"msg-burn-2"}, respSender.IgnoredIds)

	// 接收方上报，但现有过期时间更早，不应被延后
	respReceiver, err := svc.AckReadTriggers(context.Background(), &messagepb.AckReadTriggersRequest{
		UserId: "u2",
		Events: []*messagepb.ReadTriggerEvent{
			{MessageId: msg.MessageID},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"msg-burn-2"}, respReceiver.SuccessIds)

	var updated model.Message
	require.NoError(t, db.Where("message_id = ?", msg.MessageID).First(&updated).Error)
	require.NotNil(t, updated.BurnAfterReadingExpireTime)
	require.NotNil(t, updated.ExpireTime)
	assert.WithinDuration(t, now.Add(30*time.Second), *updated.BurnAfterReadingExpireTime, 2*time.Second)
	assert.WithinDuration(t, earlierExpire, *updated.ExpireTime, time.Second)
}
