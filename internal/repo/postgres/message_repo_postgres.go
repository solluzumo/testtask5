package postgres

import (
	"context"
	"testtask5/internal/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MessageRepoPostgres struct {
	db       *gorm.DB
	dbLogger *zap.Logger
}

func NewMessageRepoPostgres(db *gorm.DB, appLogger *zap.Logger) *MessageRepoPostgres {
	dbLogger := appLogger.Named("message_db")
	return &MessageRepoPostgres{
		db:       db,
		dbLogger: dbLogger,
	}
}

func (mr *MessageRepoPostgres) CreateMessage(ctx context.Context, data *domain.MessageDomain) (int, error) {
	message := &domain.MessageDomain{}
	// if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
	// 	return err
	// }

	return message, nil
}

func (mr *MessageRepoPostgres) LinkMessage(ctx context.Context, data *domain.MessageDomain) error {
	// if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
	// 	return err
	// }

	return nil
}

func (mr *MessageRepoPostgres) GetMessagesByChat(ctx context.Context, chatID int, limit int) []*domain.MessageDomain {
	var messages []*domain.MessageDomain

	return messages
}
