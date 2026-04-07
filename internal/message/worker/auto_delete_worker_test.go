package worker

import (
	"context"
	"testing"
	"time"

	"github.com/anychat/server/internal/message/model"
	"github.com/anychat/server/internal/message/repository"
	"github.com/anychat/server/pkg/logger"
	"github.com/anychat/server/pkg/notification"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type mockPublisher struct {
	notifications []*notification.Notification
}

func (m *mockPublisher) Publish(n *notification.Notification) error {
	m.notifications = append(m.notifications, n)
	return nil
}

func (m *mockPublisher) PublishToUser(string, *notification.Notification) error    { return nil }
func (m *mockPublisher) PublishToUsers([]string, *notification.Notification) error { return nil }
func (m *mockPublisher) PublishToGroup(string, *notification.Notification) error   { return nil }
func (m *mockPublisher) PublishBroadcast(*notification.Notification) error         { return nil }

func setupAutoDeleteWorkerTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Message{}))
	return db
}

func TestAutoDeleteWorker_PublishReasonByMessageType(t *testing.T) {
	logger.Log = zap.NewNop()

	db := setupAutoDeleteWorkerTestDB(t)
	repo := repository.NewMessageRepository(db)
	pub := &mockPublisher{}
	w := NewAutoDeleteWorker(repo, pub, 100, time.Minute)

	now := time.Now()
	past := now.Add(-time.Minute)
	future := now.Add(time.Minute)

	autoEarlier := now.Add(-50 * time.Second)
	burnLater := now.Add(-10 * time.Second)
	burnEarlier := now.Add(-40 * time.Second)
	autoLater := now.Add(-5 * time.Second)
	sameTime := now.Add(-20 * time.Second)

	require.NoError(t, repo.Create(context.Background(), &model.Message{
		MessageID:               "msg-auto",
		ConversationID:          "conv-1",
		ConversationType:        model.ConversationTypeSingle,
		SenderID:                "u1",
		ContentType:             model.ContentTypeText,
		Content:                 `{"text":"a"}`,
		Sequence:                1,
		Status:                  model.MessageStatusNormal,
		BurnAfterReadingSeconds: 0,
		AutoDeleteExpireTime:    &past,
		ExpireTime:              &past,
	}))
	require.NoError(t, repo.Create(context.Background(), &model.Message{
		MessageID:                  "msg-burn",
		ConversationID:             "conv-1",
		ConversationType:           model.ConversationTypeSingle,
		SenderID:                   "u1",
		ContentType:                model.ContentTypeText,
		Content:                    `{"text":"b"}`,
		Sequence:                   2,
		Status:                     model.MessageStatusNormal,
		BurnAfterReadingSeconds:    60,
		BurnAfterReadingExpireTime: &past,
		ExpireTime:                 &past,
	}))
	require.NoError(t, repo.Create(context.Background(), &model.Message{
		MessageID:                  "msg-both-auto-win",
		ConversationID:             "conv-1",
		ConversationType:           model.ConversationTypeSingle,
		SenderID:                   "u1",
		ContentType:                model.ContentTypeText,
		Content:                    `{"text":"d"}`,
		Sequence:                   4,
		Status:                     model.MessageStatusNormal,
		BurnAfterReadingSeconds:    60,
		AutoDeleteExpireTime:       &autoEarlier,
		BurnAfterReadingExpireTime: &burnLater,
		ExpireTime:                 &autoEarlier,
	}))
	require.NoError(t, repo.Create(context.Background(), &model.Message{
		MessageID:                  "msg-both-burn-win",
		ConversationID:             "conv-1",
		ConversationType:           model.ConversationTypeSingle,
		SenderID:                   "u1",
		ContentType:                model.ContentTypeText,
		Content:                    `{"text":"e"}`,
		Sequence:                   5,
		Status:                     model.MessageStatusNormal,
		BurnAfterReadingSeconds:    60,
		AutoDeleteExpireTime:       &autoLater,
		BurnAfterReadingExpireTime: &burnEarlier,
		ExpireTime:                 &burnEarlier,
	}))
	require.NoError(t, repo.Create(context.Background(), &model.Message{
		MessageID:                  "msg-both-same",
		ConversationID:             "conv-1",
		ConversationType:           model.ConversationTypeSingle,
		SenderID:                   "u1",
		ContentType:                model.ContentTypeText,
		Content:                    `{"text":"f"}`,
		Sequence:                   6,
		Status:                     model.MessageStatusNormal,
		BurnAfterReadingSeconds:    60,
		AutoDeleteExpireTime:       &sameTime,
		BurnAfterReadingExpireTime: &sameTime,
		ExpireTime:                 &sameTime,
	}))
	require.NoError(t, repo.Create(context.Background(), &model.Message{
		MessageID:               "msg-future",
		ConversationID:          "conv-1",
		ConversationType:        model.ConversationTypeSingle,
		SenderID:                "u1",
		ContentType:             model.ContentTypeText,
		Content:                 `{"text":"c"}`,
		Sequence:                3,
		Status:                  model.MessageStatusNormal,
		BurnAfterReadingSeconds: 0,
		ExpireTime:              &future,
	}))

	w.cleanup()

	var autoMsg model.Message
	require.NoError(t, db.Where("message_id = ?", "msg-auto").First(&autoMsg).Error)
	assert.Equal(t, int16(model.MessageStatusDeleted), autoMsg.Status)

	var burnMsg model.Message
	require.NoError(t, db.Where("message_id = ?", "msg-burn").First(&burnMsg).Error)
	assert.Equal(t, int16(model.MessageStatusDeleted), burnMsg.Status)

	var futureMsg model.Message
	require.NoError(t, db.Where("message_id = ?", "msg-future").First(&futureMsg).Error)
	assert.Equal(t, int16(model.MessageStatusNormal), futureMsg.Status)

	reasonToIDs := map[string][]string{}
	for _, n := range pub.notifications {
		reason, _ := n.Payload["reason"].(string)
		ids, _ := n.Payload["message_ids"].([]string)
		reasonToIDs[reason] = append(reasonToIDs[reason], ids...)
	}

	assert.ElementsMatch(t, []string{"msg-auto", "msg-both-auto-win"}, reasonToIDs["auto_delete"])
	assert.ElementsMatch(t, []string{"msg-burn", "msg-both-burn-win"}, reasonToIDs["burn_after_reading"])
	assert.Contains(t, reasonToIDs["both"], "msg-both-same")
}
