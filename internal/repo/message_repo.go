package repo

import (
	"context"
	"testtask5/internal/domain"
)

type MessageRepostiory interface {
	CreateMessage(ctx context.Context, data *domain.MessageDomain) domain.MessageDomain
	FindMessagesByChat(ctx context.Context, chatID int, limit int) []domain.MessageDomain
	LinkMessage(ctx context.Context, data *domain.MessageDomain) error
}
