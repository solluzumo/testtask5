package models

import "time"

type ChatModel struct {
	BaseModel
	Title     string `gorm:"size:200;not null"`
	CreatedAt time.Time
}
