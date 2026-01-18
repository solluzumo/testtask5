## Тестовое задание на позицию Go Developer.

Сервис для управления чатами и сообщениями. 

Реализован REST API, работа с PostgreSQL через GORM.

Архитектура: Clean Architecture (Handler -> Service -> Repository)

База данных: PostgreSQL

ORM: GORM

Миграции: Goose

Роутер: Chi

Логгер: Zap

Тесты: Testify, Mockery

Для запуска требуется только `Docker` и `docker-compose`.

### Запуск
```bash
git clone https://github.com/solluzumo/testtask5.git
cd testtask5
docker compose -f testing.docker-compose.yml build --no-cache
docker compose -f testing.docker-compose.yml up -d
```

После запуска API будет доступно по адресу: http://localhost:8080

### Остановка
```bash
docker compose -f testing.docker-compose.yml down
```
### Тесты
В проекте написаны Unit-тесты для слоя бизнес-логики с использованием моков. Тесты расположены рядом с сервисами в internal/services/chat_service_test.go

Запуск тестов:

```Bash
go test ./internal/services -v
```

### Немного про сервисы
Те методы, которые необходимо было реализовать предполагали, что нужен агрегат, поэтому chat service зависит от репозитория chat И репозитория message в том числе.
Однако, если мы хотим реализовать уже более узкие методы, например удаление или редактирование сообщения, то есть такие методы, для которых о chat знать не нужно, то мы можем реализовать их в message service(именно поэтому я всё такие оставил его, хотя в рамках задачи он не нужен)
