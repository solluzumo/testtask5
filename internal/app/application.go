package app

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Application struct {
	instance *AppInstance
	logger   *zap.Logger

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewApplication(db *gorm.DB, logger *zap.Logger) *Application {
	ctx, cancel := context.WithCancel(context.Background())

	return &Application{
		instance: NewAppAppInstance(db, logger),
		logger:   logger,
		cancel:   cancel,
		ctx:      ctx,
	}
}

// Start запускает все долгоживущие компоненты
func (a *Application) Start() {
	a.startMetricsWorker()
}

// Stop корректно завершает всё
func (a *Application) Stop(timeout time.Duration) {
	a.logger.Info("начинаю graceful shutdown...")

	a.cancel()

	// Ждём завершения
	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		a.logger.Info("все воркеры завершены")
	case <-time.After(timeout):
		a.logger.Error("таймаут при завершении воркеров")
	}
}

func (a *Application) startMetricsWorker() {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		workersTimer, err := strconv.Atoi(os.Getenv("WORKERS_TIMER"))
		if err != nil {
			a.logger.Error("не удалось распарсить workers timer из окружения", zap.Error(err))
			workersTimer = 15 // дефолтное значение
		}
		ticker := time.NewTicker(time.Duration(workersTimer) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				chatCount := a.instance.Repos.ChatRepo.Count(context.Background())
				msgCount := a.instance.Repos.MessageRepo.Count(context.Background())
				a.logger.Info("метрики",
					zap.Int64("chats_total", chatCount),
					zap.Int64("messages_total", msgCount))
			case <-a.ctx.Done():
				return
			}
		}
	}()
}
