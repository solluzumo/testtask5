package httpHandlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"testtask5/internal/domain"
	"testtask5/internal/dto"
	"testtask5/internal/services"

	"go.uber.org/zap"
)

type ChatAPIHTTP struct {
	chatService *services.ChatService
	apiLogger   *zap.Logger
}

func NewChatAPIHTTP(mService *services.ChatService, appLogger *zap.Logger) *ChatAPIHTTP {
	apiLogger := appLogger.Named("chat_api_http")
	return &ChatAPIHTTP{
		chatService: mService,
		apiLogger:   apiLogger,
	}
}

func (ch *ChatAPIHTTP) CreateChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var createChatRequest dto.CreateChatRequest
	var chatDomain domain.ChatDomain

	if err := json.NewDecoder(r.Body).Decode(&createChatRequest); err != nil {
		ch.apiLogger.Warn("не получилось расшифровать тело запроса: ",
			zap.String("path", r.URL.Path),
			zap.Error(err))
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//Валидация на пустую строку
	if createChatRequest.Title == "" {
		ch.apiLogger.Warn("title должен быть не пустой строкой")
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	//Убираем лишние пробелы по краям строки
	trimmedTitle := strings.TrimSpace(createChatRequest.Title)

	//Валидация по длине строки
	if len(trimmedTitle) > 200 {
		ch.apiLogger.Warn("превышена максимальная длина title (200)")
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	chatDomain.Title = trimmedTitle
	//Создаём чат
	result, err := ch.chatService.CreateChat(ctx, &chatDomain)
	if err != nil {
		ch.apiLogger.Error("не удалось создать чат: ", zap.Error(err))
		http.Error(w, "Не удалось создать чат", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(
		&dto.CreateChatResponse{
			ID:        result.ID,
			Title:     result.Title,
			CreatedAt: result.CreatedAt,
		})
}

func (ch *ChatAPIHTTP) GetChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var getChatRequest dto.GetChatRequest

	queryParams := r.URL.Query()

	//читаем и валидируем limit (20 по умолчанию)
	limitStr := queryParams.Get("limit")
	limitInt := 20

	if limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "неверный формат limit", http.StatusBadRequest)
			return
		}
		limitInt = parsed
	}

	if err := json.NewDecoder(r.Body).Decode(&getChatRequest); err != nil {
		ch.apiLogger.Warn("не получилось расшифровать тело запроса: ",
			zap.String("path", r.URL.Path),
			zap.Error(err))
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//Валидация на пустую строку
	if getChatRequest.Title == "" {
		ch.apiLogger.Warn("title должен быть не пустой строкой")
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	//Убираем лишние пробелы по краям строки
	trimmedTitle := strings.TrimSpace(getChatRequest.Title)

	//Валидация по длине строки
	if len(trimmedTitle) > 200 {
		ch.apiLogger.Warn("превышена максимальная длина title (200)")
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	//Получаем чат и сообщения
	result, err := ch.chatService.GetChatByTitle(ctx, trimmedTitle, limitInt)
	if err != nil {
		ch.apiLogger.Error("не удалось создать чат: ", zap.Error(err))
		http.Error(w, "Не удалось создать чат", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(
		&dto.CreateChatResponse{
			ID:        result.ID,
			Title:     result.Title,
			CreatedAt: result.CreatedAt,
		})
}

func (ch *ChatAPIHTTP) DeleteChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	chatID := r.PathValue("id")

	//валидация айдишника
	if chatID == "" {
		ch.apiLogger.Warn("chatid не должен быть пустым")
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(chatID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if id == 0 {
		http.Error(w, "ID must be greater than 0", http.StatusBadRequest)
		return
	}
	//конец валидации айдишника

	//Получаем чат и сообщения
	if err := ch.chatService.DeleteChatByID(ctx, id); err != nil {
		ch.apiLogger.Error("не удалось создать чат: ", zap.Error(err))
		http.Error(w, "Не удалось создать чат", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&dto.DeleteChatResponse{
		Content:    "чат успешно удалён",
		StatusCode: http.StatusNoContent,
	})
}

func (ch *ChatAPIHTTP) SendMessage(w http.ResponseWriter, r *http.Request) {
	var createMessageRequest dto.CreateMessageRequest
	var messageDomain domain.MessageDomain

	ctx := r.Context()

	chatID := r.PathValue("id")

	if err := json.NewDecoder(r.Body).Decode(&createMessageRequest); err != nil {
		ch.apiLogger.Warn("не получилось расшифровать тело запроса: ",
			zap.String("path", r.URL.Path),
			zap.Error(err))
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//Валидация на пустую строку
	if createMessageRequest.Text == "" {
		ch.apiLogger.Warn("title должен быть не пустой строкой")
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	//Убираем лишние пробелы по краям строки
	trimmedText := strings.TrimSpace(createMessageRequest.Text)

	//Валидация на длину строки
	if len(trimmedText) > 5000 {
		ch.apiLogger.Warn("превышена максимальная длина title (200)")
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	//валидация айдишника
	if chatID == "" {
		ch.apiLogger.Warn("title должен быть не пустой строкой")
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(chatID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if id == 0 {
		http.Error(w, "ID must be greater than 0", http.StatusBadRequest)
		return
	}
	//конец валидации айдишника

	messageDomain.ChatID = id
	messageDomain.Text = trimmedText

	//Отправляем сообщение в чат
	result, err := ch.chatService.SendMessage(ctx, &messageDomain)
	if err != nil {
		ch.apiLogger.Error("не удалось создать чат: ", zap.Error(err))

		if errors.Is(err, domain.ErrChatNotFound) {
			http.Error(w, "Чат с таким id не найден", http.StatusNotFound)
			return
		}

		http.Error(w, "Не удалось создать чат", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}
