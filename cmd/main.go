package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testtask5/internal/app"
	"testtask5/internal/middleware"
	"time"

	_ "github.com/lib/pq" // драйвер PostgreSQL
	"go.uber.org/zap"

	"github.com/go-chi/chi"
)

func main() {

	//БАЗОВЫЙ ЛОГГЕР
	logger := app.NewZapLogger()
	defer logger.Sync()

	//ГЛОБАЛЬНЫЙ ЛОГГЕР ПРИЛОЖЕНИЯ
	appLogger := logger.Named("app")

	//БАЗА ДАННЫХ
	db, err := app.InitDb()
	if err != nil {
		appLogger.Error("не удалось подключиться к базе данных: ", zap.Error(err))
	}
	appLogger.Info("база данных создана ", zap.Time("started", time.Now()))

	//Создаём app для DI
	appInstance := app.NewAppAppInstance(db, appLogger)

	//Создаём application для фоновых задач
	application := app.NewApplication(db, logger)
	application.Start()

	router := chi.NewRouter()

	router.Use(middleware.LoggingMiddleWare(appLogger))
	router.Use(middleware.TimeoutMiddleware(2 * time.Second))

	app.RegisterRoutes(router, appInstance)

	srv := &http.Server{
		Addr:    ":" + os.Getenv("HTTP_PORT"),
		Handler: router,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		appLogger.Info("сервер запущен и слушает ", zap.Time("started", time.Now()))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Не удалось запустить сервер", zap.Error(err))
		}
	}()

	// Ждём Ctrl+C
	<-ctx.Done()
	appLogger.Info("Завершаем работу сервера")

	//Завершаем работу воркеров
	application.Stop(10 * time.Second)

	// Завершаем HTTP-сервер
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		appLogger.Fatal("Не удалось корректно завершить работу", zap.Error(err))
	}

}
