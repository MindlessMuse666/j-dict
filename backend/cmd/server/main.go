package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"jp-ru-dict/backend/internal/config"
	"jp-ru-dict/backend/internal/db"
	"jp-ru-dict/backend/internal/handler"
	"jp-ru-dict/backend/internal/middleware"
	"jp-ru-dict/backend/internal/repository"
	"jp-ru-dict/backend/internal/service"
	"jp-ru-dict/backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Server struct {
	router *gin.Engine
	server *http.Server
}

func NewServer(database *sql.DB, cfg *config.Config) *Server {
	// Инициализация репозиториев
	authRepo := repository.NewAuthRepository(database)
	wordsRepo := repository.NewWordsRepository(database)

	// Инициализация сервисов
	authService := service.NewAuthService(authRepo, cfg.JWT.Secret, cfg.JWT.Expiry)
	wordsService := service.NewWordsService(wordsRepo)

	// Инициализация обработчиков
	authHandler := handler.NewAuthHandler(authService)
	wordsHandler := handler.NewWordsHandler(wordsService)
	healthHandler := handler.NewHealthHandler()

	// Настройка Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Настройка CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Публичные маршруты
	router.POST("/api/auth/register", authHandler.Register)
	router.POST("/api/auth/login", authHandler.Login)

	// Защищенные маршруты (требуют аутентификации)
	authorized := router.Group("/api")
	authorized.Use(middleware.AuthMiddleware(authService))
	{
		authorized.GET("/auth/me", authHandler.Me)
		authorized.POST("/auth/logout", authHandler.Logout)

		authorized.GET("/words", wordsHandler.GetWords)
		authorized.GET("/words/search", wordsHandler.SearchWords)
		authorized.GET("/words/:id", wordsHandler.GetWord)
		authorized.POST("/words", wordsHandler.CreateWord)
		authorized.PATCH("/words/:id", wordsHandler.UpdateWord)
		authorized.DELETE("/words/:id", wordsHandler.DeleteWord)
	}

	// Системные маршруты
	router.GET("/api/health", healthHandler.HealthCheck)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		router: router,
		server: server,
	}
}

func (s *Server) Run() {
	log.Printf("Сервер запущен на порту %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func runMigrations(cfg config.DBConfig) error {
	// Формируем строку подключения для миграций
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode,
	)

	m, err := migrate.New(
		"file://db/migrations",
		connStr,
	)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Миграции успешно применены")
	return nil
}

func main() {
	// Обработка флагов командной строки
	var migrateFlag bool
	flag.BoolVar(&migrateFlag, "migrate", false, "Выполнить миграции")
	flag.Parse()

	// Загрузка конфигурации
	cfg := config.Load()

	// Инициализация логгера
	logger.Init(cfg.LogLevel)

	// Если указан флаг миграций - выполняем их
	if migrateFlag {
		if err := runMigrations(cfg.DB); err != nil {
			log.Fatalf("Ошибка выполнения миграций: %v", err)
		}
		return
	}

	// Подключение к базе данных
	database, err := db.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer database.Close()

	// Инициализация и запуск сервера
	server := NewServer(database, &cfg)
	go server.Run()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка завершения работы сервера: %v", err)
	}

	log.Println("Сервер остановлен")
}
