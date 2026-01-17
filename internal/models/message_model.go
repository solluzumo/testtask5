package models

import "time"

type MessageModel struct {
	BaseModel
	ChatID    int
	Text      string `gorm:"size:5000;not null"`
	CreatedAt time.Time
}
