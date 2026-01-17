package postgres

import (
	"context"
	"testtask5/internal/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ChatRepoPostgres struct {
	db       *gorm.DB
	dbLogger *zap.Logger
}

func NewChatRepoPostgres(db *gorm.DB, appLogger *zap.Logger) *ChatRepoPostgres {
	dbLogger := appLogger.Named("chat_db")
	return &ChatRepoPostgres{
		db:       db,
		dbLogger: dbLogger,
	}
}

func (cr *ChatRepoPostgres) CreateChat(ctx context.Context, title string) (int, error) {
	chat := &domain.ChatDomain{}

	return chat, nil
}

func (cr *ChatRepoPostgres) FindChatByTitle(ctx context.Context, title string) (*domain.ChatDomain, error) {
	chat := &domain.ChatDomain{}

	return chat, nil
}

func (cr *ChatRepoPostgres) FindChatByID(ctx context.Context, chatID int) (*domain.ChatDomain, error) {
	chat := &domain.ChatDomain{}

	return chat, nil
}

func (cr *ChatRepoPostgres) DeleteChat(ctx context.Context, chatID int) bool {

	return true
}
