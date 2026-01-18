package httpHandlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"testtask5/internal/domain"
	"testtask5/internal/dto"
	"testtask5/internal/services"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type ChatAPIHTTP struct {
	chatService *services.ChatService
	apiLogger   *zap.Logger
}

func NewChatAPIHTTP(mService *services.ChatService, appLogger *zap.Logger) *ChatAPIHTTP {
	return &ChatAPIHTTP{
		chatService: mService,
		apiLogger:   appLogger.Named("chat_api_http"),
	}
}

// Создать чат
func (ch *ChatAPIHTTP) CreateChat(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateChatRequest
	if err := ch.decodeJSON(r, &req); err != nil {
		ch.respondError(w, "некорректный JSON", http.StatusBadRequest, err)
		return
	}

	// валидация title
	if err := req.Validate(); err != nil {
		ch.respondError(w, "ошибка валидации title", http.StatusBadRequest, err)
		return
	}

	chat := &domain.ChatDomain{Title: req.Title}
	result, err := ch.chatService.CreateChat(r.Context(), chat)
	if err != nil {
		ch.handleDomainError(w, err)
		return
	}

	ch.respondJSON(w, http.StatusCreated, &dto.CreateChatResponse{
		ID:        result.ID,
		Title:     result.Title,
		CreatedAt: result.CreatedAt,
	})
}

// Получить чат и limit сообщений
func (ch *ChatAPIHTTP) GetChat(w http.ResponseWriter, r *http.Request) {
	id, err := ch.parseID(r)
	if err != nil {
		ch.respondError(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	limit := ch.parseLimit(r)

	result, err := ch.chatService.GetChatById(r.Context(), id, limit)
	if err != nil {
		ch.handleDomainError(w, err)
		return
	}

	ch.respondJSON(w, http.StatusOK, &dto.CreateChatResponse{
		ID:        result.ID,
		Title:     result.Title,
		CreatedAt: result.CreatedAt,
		Messages:  result.Messages,
	})
}

// Удаление чата
func (ch *ChatAPIHTTP) DeleteChat(w http.ResponseWriter, r *http.Request) {
	id, err := ch.parseID(r)
	if err != nil {
		ch.respondError(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	if err := ch.chatService.DeleteChatByID(r.Context(), id); err != nil {
		ch.handleDomainError(w, err)
		return
	}

	ch.respondJSON(w, http.StatusOK, &dto.DeleteChatResponse{
		Content:    "чат успешно удалён",
		StatusCode: http.StatusOK, // Обычно No Content (204) не возвращает тело, но оставил как у вас
	})
}

// Отправка сообщения
func (ch *ChatAPIHTTP) SendMessage(w http.ResponseWriter, r *http.Request) {
	id, err := ch.parseID(r)
	if err != nil {
		ch.respondError(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	var req dto.CreateMessageRequest
	if err := ch.decodeJSON(r, &req); err != nil {
		ch.respondError(w, "некорректный JSON", http.StatusBadRequest, err)
		return
	}

	// валидация text
	if err := req.Validate(); err != nil {
		ch.respondError(w, "ошибка валидации text", http.StatusBadRequest, err)
		return
	}

	msgDomain := &domain.MessageDomain{
		ChatID: id,
		Text:   req.Text,
	}

	result, err := ch.chatService.SendMessage(r.Context(), msgDomain)
	if err != nil {
		ch.handleDomainError(w, err)
		return
	}

	ch.respondJSON(w, http.StatusOK, result)
}

// parseID извлекает и валидирует ID из URL
func (ch *ChatAPIHTTP) parseID(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		return 0, errors.New("chatID отсутствует")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		return 0, errors.New("chatID должен быть положительным числом")
	}
	return id, nil
}

// parseLimit парсит параметр limit или возвращает дефолтное значение
func (ch *ChatAPIHTTP) parseLimit(r *http.Request) int {
	defaultLimit := 20
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		return defaultLimit
	}
	if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
		return parsed
	}
	return defaultLimit
}

// decodeJSON декодирует тело запроса
func (ch *ChatAPIHTTP) decodeJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// respondJSON отправляет стандартизированный JSON ответ
func (ch *ChatAPIHTTP) respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			ch.apiLogger.Error("ошибка при кодировании ответа", zap.Error(err))
		}
	}
}

// respondError отправляет ошибку в формате JSON
func (ch *ChatAPIHTTP) respondError(w http.ResponseWriter, message string, code int, err error) {
	if err != nil {
		ch.apiLogger.Warn(message, zap.Error(err))
	} else {
		ch.apiLogger.Warn(message)
	}
	ch.respondJSON(w, code, map[string]string{"error": message})
}

// handleDomainError маппит ошибки домена на HTTP коды
func (ch *ChatAPIHTTP) handleDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrChatNotFound):
		ch.respondError(w, "чат не найден", http.StatusNotFound, err)
	case errors.Is(err, domain.ErrChatAlreadyExists):
		ch.respondError(w, "чат с таким названием уже существует", http.StatusBadRequest, err)
	case errors.Is(err, context.Canceled):
		// Клиент ушел, отвечать некому, просто логируем
		ch.apiLogger.Info("запрос отменён клиентом")
	default:
		ch.apiLogger.Error("внутренняя ошибка сервера", zap.Error(err))
		ch.respondError(w, "internal server error", http.StatusInternalServerError, nil)
	}
}
