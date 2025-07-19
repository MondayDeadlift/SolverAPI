package config

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type HTTPServer struct {
	Host         string        `env:"SERVER_HOST" envDefault:"localhost"`
	Port         string        `env:"SERVER_PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"30s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"30s"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" envDefault:"60s"`
}

type Config struct {
	Environment string `env:"ENV" envDefault:"development"`
	HTTPServer  HTTPServer
	Codewars    CodewarsConfig
	Database    DatabaseConfig
}

type CodewarsConfig struct {
	APIURL string `env:"CODEWARS_API_URL" envDefault:"https://www.codewars.com/api/v1"`
}

type DatabaseConfig struct {
	DSN string `env:"DB_DSN" envDefault:"postgres://user:password@localhost:5432/codewars?sslmode=disable"`
}

func Load() (*Config, error) {
	//Загрузка .env файла
	if err := godotenv.Load("config/local.env"); err != nil {
		log.Println("Файл .env не найден, используются переменные окружения")
	}

	//Парсинг переменных в структуру
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	//Проверка валидации
	if cfg.HTTPServer.Port == "" {
		return nil, errors.New("SERVER_PORT не может быть пустым")
	}

	return cfg, nil
}

// NewEcho создает экземпляр Echo с настройками из конфига
func (s *HTTPServer) NewEcho() *echo.Echo {
	e := echo.New()

	// Настройка HTTP-сервера
	e.Server = &http.Server{
		Addr:         s.Host + ":" + s.Port,
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
		IdleTimeout:  s.IdleTimeout,
	}

	return e
}

// NewServer создает и возвращает настроенный Echo-сервер
func NewServer() (*echo.Echo, *Config, error) {
	// Загружаем конфиг
	cfg, err := Load()
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка загрузки конфига: %w", err)
	}

	// Создаем Echo-сервер
	s := cfg.HTTPServer.NewEcho()

	return s, cfg, nil
}
