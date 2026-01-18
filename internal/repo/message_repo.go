package repo

import (
	"context"
	"testtask5/internal/domain"
)

type MessageRepostiory interface {
	CreateMessage(ctx context.Context, data *domain.MessageDomain) (*domain.MessageDomain, error)
	GetMessagesByChaWithLimit(ctx context.Context, chatID int, limit int) []*domain.MessageDomain
	DeleteMessages(ctx context.Context, chatID int) error
	Count(ctx context.Context) int64
}
