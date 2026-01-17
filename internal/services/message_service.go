package services

import (
	"testtask5/internal/repo/postgres"

	"go.uber.org/zap"
)

type MessageService struct {
	messageRepo   *postgres.MessageRepoPostgres
	serviceLogger *zap.Logger
}

func NewMessageService(mRepo *postgres.MessageRepoPostgres, appLogger *zap.Logger) *MessageService {
	serviceLogger := appLogger.Named("message_service")
	return &MessageService{
		messageRepo:   mRepo,
		serviceLogger: serviceLogger,
	}
}
