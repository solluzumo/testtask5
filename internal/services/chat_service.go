package services

import (
	"context"
	"testtask5/internal/domain"
	"testtask5/internal/repo/postgres"
	"time"

	"go.uber.org/zap"
)

type ChatService struct {
	chatRepo      *postgres.ChatRepoPostgres
	messageRepo   *postgres.MessageRepoPostgres
	serviceLogger *zap.Logger
}

func NewChatService(mRepo *postgres.MessageRepoPostgres, cRepo *postgres.ChatRepoPostgres, appLogger *zap.Logger) *ChatService {
	serviceLogger := appLogger.Named("chat_service")
	return &ChatService{
		chatRepo:      cRepo,
		messageRepo:   mRepo,
		serviceLogger: serviceLogger,
	}
}

func (cs *ChatService) CreateChat(ctx context.Context, chat *domain.ChatDomain) (*domain.ChatDomain, error) {

	chatID, err := cs.chatRepo.CreateChat(ctx, chat.Title)
	if err != nil {
		return nil, err
	}
	//Задаём дату и id
	chat.CreatedAt = time.Now()
	chat.ID = chatID

	return chat, nil
}

func (cs *ChatService) GetChatByTitle(ctx context.Context, title string, limit int) (*domain.ChatDomain, error) {
	chat, err := cs.chatRepo.FindChatByTitle(ctx, title)
	if err != nil {
		cs.serviceLogger.Warn("чат с таким title не существует", zap.String("title", title))
		return nil, domain.ErrChatNotFound
	}

	messages := cs.messageRepo.GetMessagesByChat(ctx, chat.ID, limit)

	chat.Messages = messages

	return chat, nil
}

func (cs *ChatService) GetChatById(ctx context.Context, chatID int) (*domain.ChatDomain, error) {
	chat, err := cs.chatRepo.FindChatByID(ctx, chatID)
	if err != nil {
		cs.serviceLogger.Warn("чат с таким id не существует", zap.Int("chatID", chatID))
		return nil, domain.ErrChatNotFound
	}
	return chat, nil
}

func (cs *ChatService) DeleteChatByID(ctx context.Context, chatID int) error {

	chat, _ := cs.GetChatById(ctx, chatID)

	if chat == nil {
		cs.serviceLogger.Warn("чат с таким id не существует", zap.Int("chatID", chatID))
		return domain.ErrChatNotFound
	}

	result := cs.chatRepo.DeleteChat(ctx, chatID)

	if !result {
		return nil
	}
	return nil
}

func (cs *ChatService) SendMessage(ctx context.Context, message *domain.MessageDomain) (*domain.MessageDomain, error) {
	chat, _ := cs.GetChatById(ctx, message.ChatID)
	if chat == nil {
		cs.serviceLogger.Warn("чат с таким id не существует", zap.Int("chatID", message.ChatID))
		return nil, domain.ErrChatNotFound
	}

	//создаём сообщение
	messageID, err := cs.messageRepo.CreateMessage(ctx, message)
	if err != nil {
		cs.serviceLogger.Error("не удалось создать сообщение", zap.Error(err))
		return nil, err
	}

	//Задаём дату и id
	message.CreatedAt = time.Now()
	message.ID = messageID

	return message, nil
}
