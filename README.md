# Тестовое задание на позицию Go Developer.

Сервис для управления чатами и сообщениями. 

Реализован REST API, работа с PostgreSQL через GORM.

## Стек
- Архитектура: Clean Architecture (Handler -> Service -> Repository)

- База данных: PostgreSQL

- ORM: GORM

- Миграции: Goose

- Роутер: Chi

- Логгер: Zap

- Тесты: Testify, Mockery

## Запуск
```bash
git clone https://github.com/solluzumo/testtask5.git
cd testtask5
docker compose -f testing.docker-compose.yml build --no-cache
docker compose -f testing.docker-compose.yml up -d
```

После запуска API будет доступно по адресу: http://localhost:8080

## Остановка
```bash
docker compose -f testing.docker-compose.yml down
```
## Тесты
В проекте написаны Unit-тесты для слоя бизнес-логики с использованием моков. Тесты расположены рядом с сервисами в internal/services/chat_service_test.go

Запуск тестов:

```Bash
go test ./internal/services -v
```

## Немного про сервисы
Для того, чтобы выполнить техническое задание, я решил описать агрегат - ChatService, который агрегирует в себе Chat и Message, в нём и реализованый все методы.

Несмотря на то, что это не предусмотрено, я решил оставить message service и message handler, ведь они могут понадобиться для методов, работающих только с сущностью message - удаление, редактирование сообщений и тд.

## Что не надо было, но сделал
Gracefull shutdown, фоновые воркеры через горутины, которые собирают некоторую метрику, интервал для сбора метрики задаётся в docker compose, по мелочи: обработка случая, когда клиент ушёл, не дождавшись ответа
