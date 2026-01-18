package domain

import "time"

type MessageDomain struct {
	ID        int
	Text      string
	ChatID    int
	CreatedAt time.Time
}
