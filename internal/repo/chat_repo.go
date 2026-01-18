package repo

import (
	"context"
	"testtask5/internal/domain"
)

type ChatRepostiory interface {
	CreateChat(ctx context.Context, data *domain.ChatDomain) (*domain.ChatDomain, error)
	FindChatById(ctx context.Context, data *domain.ChatDomain) (*domain.ChatDomain, error)
	ChatExists(ctx context.Context, param FilterParam) (bool, error)
	DeleteChat(ctx context.Context, chatID int) error
	Count(ctx context.Context) int64
}
