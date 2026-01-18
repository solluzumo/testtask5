package dto

import (
	"testtask5/internal/domain"
	"time"
)

type CreateChatResponse struct {
	ID        int                     `json:"id"`
	Title     string                  `json:"title"`
	CreatedAt time.Time               `json:"created_at"`
	Messages  []*domain.MessageDomain `json:"messages"`
}
