package domain

import (
	"time"
)

type ChatDomain struct {
	BaseDomain
	Title     string
	CreatedAt time.Time
	Messages  []*MessageDomain
}
