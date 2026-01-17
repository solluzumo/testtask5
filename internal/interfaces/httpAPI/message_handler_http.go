package httpHandlers

import (
	"testtask5/internal/services"

	"go.uber.org/zap"
)

type MessageAPIHTTP struct {
	messageService *services.MessageService
	apiLogger      *zap.Logger
}

func NewMessageAPIHTTP(mService *services.MessageService, appLogger *zap.Logger) *MessageAPIHTTP {
	apiLogger := appLogger.Named("message_api_http")
	return &MessageAPIHTTP{
		messageService: mService,
		apiLogger:      apiLogger,
	}
}
