package dto

import (
	"testtask5/internal/domain"
	"time"
)

type GetChatResonse struct {
	ChatID    int                    `json:"chat_id"`
	Title     string                 `json:"chat_title"`
	CreatedAt time.Time              `json:"created_at"`
	Messages  []domain.MessageDomain `json:"message"`
}
