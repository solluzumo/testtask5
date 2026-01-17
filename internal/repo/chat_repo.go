package repo

import (
	"context"
	"testtask5/internal/domain"
)

type ChatRepostiory interface {
	CreateChat(ctx context.Context, title string) (*domain.ChatDomain, error)
	FindChatByTitle(ctx context.Context, title string) (*domain.ChatDomain, error)
	FindChatByID(ctx context.Context, chatID int) (*domain.ChatDomain, error)
	DeleteChat(ctx context.Context, chatID int) bool
}
