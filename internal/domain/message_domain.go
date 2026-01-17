package domain

import "time"

type MessageDomain struct {
	BaseDomain
	Text      string
	ChatID    int
	CreatedAt time.Time
}
