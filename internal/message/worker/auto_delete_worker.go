package worker

import (
	"context"
	"time"

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

		messageIDs, err := w.messageRepo.GetExpiredMessageIDs(ctx, now, w.batchSize)
		if err != nil {
			logger.Error("Failed to get expired messages", zap.Error(err))
			break
		}

		if len(messageIDs) == 0 {
			break
		}

		if err := w.messageRepo.BatchUpdateStatus(ctx, messageIDs, 2); err != nil {
			logger.Error("Failed to batch update message status", zap.Error(err))
			break
		}

		logger.Info("Deleted expired messages", zap.Int("count", len(messageIDs)))

		w.publishNotification(ctx, messageIDs)
	}
}

func (w *AutoDeleteWorker) publishNotification(ctx context.Context, messageIDs []string) {
	notif := notification.NewNotification(notification.TypeMessageAutoDeleted, "", notification.PriorityNormal).
		AddPayloadField("message_ids", messageIDs)

	if err := w.notificationPub.Publish(notif); err != nil {
		logger.Warn("Failed to publish auto delete notification", zap.Error(err))
	}
}

func (w *AutoDeleteWorker) StartAsync() {
	go w.Start()
}
