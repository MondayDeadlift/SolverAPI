package app

import (
	"SolverAPI/config"
	"SolverAPI/internal/handler"
	"SolverAPI/internal/repository/postgres"
	"SolverAPI/internal/service"
	"SolverAPI/pkg/codewars"
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // Драйвер PostgreSQL

	"github.com/labstack/echo/v4"

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres" // Алиас для избежания конфликта
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Server struct {
	Echo     *echo.Echo
	Config   *config.Config
	Codewars *codewars.Client
}

func New() (*Server, error) {
	// Инициализация конфига
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Создаем Echo-сервер
	e := cfg.HTTPServer.NewEcho()

	// Инициализируем клиенты
	cwClient := codewars.NewClient(cfg.Codewars.APIURL)

	return &Server{
		Echo:     e,
		Config:   cfg,
		Codewars: cwClient,
	}, nil
}

func (s *Server) RegisterHandlers(userService *service.UserService, kataService *service.KataService) {
	userHandler := handler.NewUserHandler(s.Codewars, userService)
	kataHandler := handler.NewKataHandler(kataService)

	//health-check
	healthHandler := handler.NewHealthHandler()

	s.Echo.GET("/health", healthHandler.Check)
	s.Echo.GET("/users/:username", userHandler.GetUser)
	s.Echo.GET("/katas/random", kataHandler.GetRandomKata) // Новый эндпоинт
}

func (s *Server) Start() error {
	return s.Echo.Start(":" + s.Config.HTTPServer.Port)
}

func (s *Server) SetupDependencies() error {
	// Инициализация БД
	db, err := sql.Open("postgres", s.Config.Database.DSN)
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(5 * time.Minute)
	// Запуск миграций
	if err := s.runMigrations(); err != nil {
		return fmt.Errorf("migrations failed: %w", err)
	}

	// Проверка подключения
	if err := db.Ping(); err != nil {
		return fmt.Errorf("DB ping failed: %w", err)
	}

	// Инициализация репозиториев
	userRepo := postgres.NewUserRepository(db)
	kataRepo := postgres.NewKataRepository(db)

	// Инициализация сервисов
	userService := service.NewUserService(userRepo, s.Codewars)
	kataService := service.NewKataService(kataRepo, s.Codewars)

	// Регистрация обработчиков
	s.RegisterHandlers(userService, kataService) // Обновляем метод
	return nil
}

func (s *Server) runMigrations() error {
	// Увеличиваем таймауты подключения
	connStr := s.Config.Database.DSN + "&connect_timeout=10"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}

	// Настройки пула соединений
	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	// Проверка подключения с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("DB ping failed: %w", err)
	}

	driver, err := migratepostgres.WithInstance(db, &migratepostgres.Config{
		DatabaseName: "codewars",
	})
	if err != nil {
		db.Close()
		return fmt.Errorf("failed to create driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		db.Close()
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		db.Close()
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	db.Close()
	return nil
}
