package models

import "time"

type Message struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	ChatID    int
	Text      string `gorm:"size:5000;not null"`
	CreatedAt time.Time
}
