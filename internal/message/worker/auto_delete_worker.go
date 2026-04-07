package worker

import (
	"context"
	"time"

	"github.com/anychat/server/internal/message/model"
	"github.com/anychat/server/internal/message/repository"
	"github.com/anychat/server/pkg/logger"
	"github.com/anychat/server/pkg/notification"
	"go.uber.org/zap"
)

type AutoDeleteWorker struct {
	messageRepo     repository.MessageRepository
	notificationPub notification.Publisher
	batchSize       int
	interval        time.Duration
	stopCh          chan struct{}
}

func NewAutoDeleteWorker(
	messageRepo repository.MessageRepository,
	notificationPub notification.Publisher,
	batchSize int,
	interval time.Duration,
) *AutoDeleteWorker {
	return &AutoDeleteWorker{
		messageRepo:     messageRepo,
		notificationPub: notificationPub,
		batchSize:       batchSize,
		interval:        interval,
		stopCh:          make(chan struct{}),
	}
}

func (w *AutoDeleteWorker) Start() {
	logger.Info("AutoDeleteWorker starting", zap.Int("batchSize", w.batchSize), zap.Duration("interval", w.interval))

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-w.stopCh:
			logger.Info("AutoDeleteWorker stopped")
			return
		case <-ticker.C:
			w.cleanup()
		}
	}
}

func (w *AutoDeleteWorker) Stop() {
	close(w.stopCh)
}

func (w *AutoDeleteWorker) cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	now := time.Now()

	for {
		select {
		case <-w.stopCh:
			return
		default:
		}

		expiredMessages, err := w.messageRepo.GetExpiredMessages(ctx, now, w.batchSize)
		if err != nil {
			logger.Error("Failed to get expired messages", zap.Error(err))
			break
		}

		if len(expiredMessages) == 0 {
			break
		}

		messageIDs := make([]string, 0, len(expiredMessages))
		reasons := map[string][]string{
			"auto_delete":        {},
			"burn_after_reading": {},
			"both":               {},
		}
		for _, msg := range expiredMessages {
			messageIDs = append(messageIDs, msg.MessageID)
			reason := inferDeleteReason(msg)
			reasons[reason] = append(reasons[reason], msg.MessageID)
		}

		if err := w.messageRepo.BatchUpdateStatus(ctx, messageIDs, 2); err != nil {
			logger.Error("Failed to batch update message status", zap.Error(err))
			break
		}

		logger.Info("Deleted expired messages", zap.Int("count", len(messageIDs)))

		for reason, ids := range reasons {
			if len(ids) == 0 {
				continue
			}
			w.publishNotification(ctx, ids, reason)
		}
	}
}

func inferDeleteReason(msg *model.Message) string {
	hasAuto := msg.AutoDeleteExpireTime != nil
	hasBurn := msg.BurnAfterReadingExpireTime != nil

	if hasAuto && hasBurn {
		if msg.AutoDeleteExpireTime.Before(*msg.BurnAfterReadingExpireTime) {
			return "auto_delete"
		}
		if msg.BurnAfterReadingExpireTime.Before(*msg.AutoDeleteExpireTime) {
			return "burn_after_reading"
		}
		return "both"
	}
	if hasBurn {
		return "burn_after_reading"
	}
	return "auto_delete"
}

func (w *AutoDeleteWorker) publishNotification(ctx context.Context, messageIDs []string, reason string) {
	notif := notification.NewNotification(notification.TypeMessageAutoDeleted, "", notification.PriorityNormal).
		AddPayloadField("message_ids", messageIDs).
		AddPayloadField("reason", reason)

	if err := w.notificationPub.Publish(notif); err != nil {
		logger.Warn("Failed to publish auto delete notification", zap.Error(err))
	}
}

func (w *AutoDeleteWorker) StartAsync() {
	go w.Start()
}
