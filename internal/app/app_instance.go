package app

import (
	httpHandlers "testtask5/internal/interfaces/httpAPI"
	"testtask5/internal/repo/postgres"
	"testtask5/internal/services"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppInstance struct {
	Repos    *AppRepos
	Services *AppServices
	API      *AppAPIs
}

type AppRepos struct {
	ChatRepo    *postgres.ChatRepoPostgres
	MessageRepo *postgres.MessageRepoPostgres
}

type AppServices struct {
	ChatService    *services.ChatService
	MessageService *services.MessageService
}

type AppAPIs struct {
	ChatAPI    *httpHandlers.ChatAPIHTTP
	MessageAPI *httpHandlers.MessageAPIHTTP
}

func NewAppAppInstance(db *gorm.DB, appLogger *zap.Logger) *AppInstance {
	appRepos := &AppRepos{
		ChatRepo:    postgres.NewChatRepoPostgres(db, appLogger),
		MessageRepo: postgres.NewMessageRepoPostgres(db, appLogger),
	}
	appServices := &AppServices{
		ChatService:    services.NewChatService(appRepos.MessageRepo, appRepos.ChatRepo, appLogger),
		MessageService: services.NewMessageService(appRepos.MessageRepo, appLogger),
	}
	appAPIs := &AppAPIs{
		ChatAPI:    httpHandlers.NewChatAPIHTTP(appServices.ChatService, appLogger),
		MessageAPI: httpHandlers.NewMessageAPIHTTP(appServices.MessageService, appLogger),
	}
	return &AppInstance{
		Repos:    appRepos,
		Services: appServices,
		API:      appAPIs,
	}
}
