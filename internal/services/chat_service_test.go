package services

import (
	"context"
	"errors"
	"testing"
	"testtask5/internal/domain"
	"testtask5/internal/repo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockChatRepository struct {
	mock.Mock
}

type MockMessageRepository struct {
	mock.Mock
}

func (m *MockChatRepository) CreateChat(ctx context.Context, data *domain.ChatDomain) (*domain.ChatDomain, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(*domain.ChatDomain), args.Error(1)
}

func (m *MockChatRepository) Count(ctx context.Context) int64 {
	args := m.Called(ctx)
	return int64(args.Int(0))
}

func (m *MockMessageRepository) Count(ctx context.Context) int64 {
	args := m.Called(ctx)
	return int64(args.Int(0))
}

func (m *MockChatRepository) ChatExists(ctx context.Context, param repo.FilterParam) (bool, error) {
	args := m.Called(ctx, param)
	return args.Bool(0), args.Error(1)
}

func (m *MockChatRepository) FindChatById(ctx context.Context, data *domain.ChatDomain) (*domain.ChatDomain, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChatDomain), args.Error(1)
}

func (m *MockChatRepository) DeleteChat(ctx context.Context, chatID int) error {
	args := m.Called(ctx, chatID)
	return args.Error(0)
}

func (m *MockMessageRepository) CreateMessage(ctx context.Context, data *domain.MessageDomain) (*domain.MessageDomain, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(*domain.MessageDomain), args.Error(1)
}

func (m *MockMessageRepository) GetMessagesByChaWithLimit(ctx context.Context, chatID int, limit int) []*domain.MessageDomain {
	args := m.Called(ctx, chatID, limit)
	return args.Get(0).([]*domain.MessageDomain)
}

func (m *MockMessageRepository) DeleteMessages(ctx context.Context, chatID int) error {
	args := m.Called(ctx, chatID)
	return args.Error(0)
}
func setupTest(t *testing.T) (*ChatService, *MockChatRepository, *MockMessageRepository) {
	logger := zap.NewNop()
	mockChatRepo := new(MockChatRepository)
	mockMessageRepo := new(MockMessageRepository)
	svc := NewChatService(mockMessageRepo, mockChatRepo, logger)
	return svc, mockChatRepo, mockMessageRepo
}

func TestChatService_CreateChat(t *testing.T) {
	ctx := context.Background()

	t.Run("чата не существует - успешное создание", func(t *testing.T) {
		svc, mockChatRepo, _ := setupTest(t) // Инициализация ВНУТРИ теста
		chat := &domain.ChatDomain{Title: "test"}

		mockChatRepo.On("ChatExists", ctx, repo.FilterParam{Field: "title", Value: "test"}).
			Return(false, nil)

		mockChatRepo.On("CreateChat", ctx, chat).
			Return(chat, nil)

		result, err := svc.CreateChat(ctx, chat)

		assert.NoError(t, err)
		assert.Equal(t, chat, result)
		mockChatRepo.AssertExpectations(t)
	})

	t.Run("чат уже существует — ошибка", func(t *testing.T) {
		svc, mockChatRepo, _ := setupTest(t) // Инициализация ВНУТРИ теста
		chat := &domain.ChatDomain{Title: "existing"}

		mockChatRepo.On("ChatExists", ctx, repo.FilterParam{Field: "title", Value: "existing"}).
			Return(true, nil)

		result, err := svc.CreateChat(ctx, chat)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrChatAlreadyExists)
		mockChatRepo.AssertExpectations(t)
	})

	t.Run("ошибка при проверке существования — возвращается ошибка", func(t *testing.T) {
		svc, mockChatRepo, _ := setupTest(t) // Инициализация ВНУТРИ теста
		chat := &domain.ChatDomain{Title: "invalid"}

		mockChatRepo.On("ChatExists", ctx, repo.FilterParam{Field: "title", Value: "invalid"}).
			Return(false, errors.New("field not allowed"))

		result, err := svc.CreateChat(ctx, chat)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrFieldIsNotAllowed)
		mockChatRepo.AssertExpectations(t)
	})
}

func TestChatService_ChatExists(t *testing.T) {
	ctx := context.Background()
	svc, mockChatRepo, _ := setupTest(t)

	param := repo.FilterParam{Field: "title", Value: "test"}

	mockChatRepo.On("ChatExists", ctx, param).
		Return(true, nil)

	exists, err := svc.ChatExists(ctx, param)

	assert.NoError(t, err)
	assert.True(t, exists)
	mockChatRepo.AssertExpectations(t)
}

func TestChatService_GetChatById(t *testing.T) {
	ctx := context.Background()

	t.Run("успешное получение чата с сообщениями", func(t *testing.T) {
		svc, mockChatRepo, mockMessageRepo := setupTest(t) // Инициализация ВНУТРИ теста
		chatID := 123
		limit := 10

		chatFromRepo := &domain.ChatDomain{
			ID:    chatID,
			Title: "test chat",
		}
		messages := []*domain.MessageDomain{
			{ID: 1, Text: "msg1"},
			{ID: 2, Text: "msg2"},
		}

		mockChatRepo.On("FindChatById", ctx, &domain.ChatDomain{ID: chatID}).
			Return(chatFromRepo, nil)
		mockMessageRepo.On("GetMessagesByChaWithLimit", ctx, chatID, limit).
			Return(messages)

		result, err := svc.GetChatById(ctx, chatID, limit)

		assert.NoError(t, err)
		assert.Equal(t, chatFromRepo.ID, result.ID)
		assert.Equal(t, chatFromRepo.Title, result.Title)
		assert.Equal(t, messages, result.Messages)
		mockChatRepo.AssertExpectations(t)
		mockMessageRepo.AssertExpectations(t)
	})

	t.Run("чат не найден", func(t *testing.T) {
		svc, mockChatRepo, _ := setupTest(t) // Инициализация ВНУТРИ теста
		chatID := 999
		limit := 5

		mockChatRepo.On("FindChatById", ctx, &domain.ChatDomain{ID: chatID}).
			Return(nil, errors.New("record not found"))

		result, err := svc.GetChatById(ctx, chatID, limit)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrChatNotFound)
		mockChatRepo.AssertExpectations(t)
	})
}

func TestChatService_DeleteChatByID(t *testing.T) {
	ctx := context.Background()

	t.Run("успешное удаление", func(t *testing.T) {
		svc, mockChatRepo, _ := setupTest(t) // НОВЫЙ МОК
		chatID := 123

		mockChatRepo.On("ChatExists", ctx, repo.FilterParam{Field: "id", Value: chatID}).
			Return(true, nil)
		mockChatRepo.On("DeleteChat", ctx, chatID).
			Return(nil)

		err := svc.DeleteChatByID(ctx, chatID)

		assert.NoError(t, err)
		mockChatRepo.AssertExpectations(t)
	})

	t.Run("чат не найден", func(t *testing.T) {
		svc, mockChatRepo, _ := setupTest(t) // НОВЫЙ МОК
		chatID := 999

		mockChatRepo.On("ChatExists", ctx, repo.FilterParam{Field: "id", Value: chatID}).
			Return(false, nil)

		err := svc.DeleteChatByID(ctx, chatID)

		assert.ErrorIs(t, err, domain.ErrChatNotFound)
		mockChatRepo.AssertExpectations(t)
	})

	t.Run("ошибка при проверке существования", func(t *testing.T) {
		svc, mockChatRepo, _ := setupTest(t) // НОВЫЙ МОК
		chatID := 123

		// Теперь этот On не будет конфликтовать с первым тестом, так как мок чистый
		mockChatRepo.On("ChatExists", ctx, repo.FilterParam{Field: "id", Value: chatID}).
			Return(false, errors.New("field not allowed"))

		err := svc.DeleteChatByID(ctx, chatID)

		assert.ErrorIs(t, err, domain.ErrFieldIsNotAllowed)
		mockChatRepo.AssertExpectations(t)
	})
}

func TestChatService_SendMessage(t *testing.T) {
	ctx := context.Background()

	t.Run("успешная отправка сообщения", func(t *testing.T) {
		svc, mockChatRepo, mockMessageRepo := setupTest(t)
		message := &domain.MessageDomain{
			ChatID: 123,
			Text:   "hello",
		}
		savedMessage := &domain.MessageDomain{
			ID:     1,
			ChatID: 123,
			Text:   "hello",
		}

		mockChatRepo.On("ChatExists", ctx, repo.FilterParam{Field: "id", Value: 123}).
			Return(true, nil)
		mockMessageRepo.On("CreateMessage", ctx, message).
			Return(savedMessage, nil)

		result, err := svc.SendMessage(ctx, message)

		assert.NoError(t, err)
		assert.Equal(t, savedMessage, result)
		mockChatRepo.AssertExpectations(t)
		mockMessageRepo.AssertExpectations(t)
	})

	t.Run("чат не найден", func(t *testing.T) {
		svc, mockChatRepo, _ := setupTest(t)
		message := &domain.MessageDomain{ChatID: 999, Text: "hi"}

		mockChatRepo.On("ChatExists", ctx, repo.FilterParam{Field: "id", Value: 999}).
			Return(false, nil)

		result, err := svc.SendMessage(ctx, message)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrChatNotFound)
		mockChatRepo.AssertExpectations(t)
	})

	t.Run("ошибка при проверке существования", func(t *testing.T) {
		svc, mockChatRepo, _ := setupTest(t) // Чистый мок!
		message := &domain.MessageDomain{ChatID: 123, Text: "hi"}

		// Теперь ChatExists вернет ошибку, и код не пойдет дальше создавать сообщение
		mockChatRepo.On("ChatExists", ctx, repo.FilterParam{Field: "id", Value: 123}).
			Return(false, errors.New("field not allowed"))

		result, err := svc.SendMessage(ctx, message)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrFieldIsNotAllowed)
		mockChatRepo.AssertExpectations(t)
	})
}
