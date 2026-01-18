package main

import (
	"net/http"
	"os"
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
	appInstance := app.NewApp(db, appLogger)

	router := chi.NewRouter()

	router.Use(middleware.LoggingMiddleWare(appLogger))
	router.Use(middleware.TimeoutMiddleware(2 * time.Second))

	app.RegisterRoutes(router, appInstance)

	srv := &http.Server{
		Addr:    ":" + os.Getenv("HTTP_PORT"),
		Handler: router,
	}
	srv.ListenAndServe()
	appLogger.Info("сервер запущен и слушает ", zap.Time("started", time.Now()))
}
