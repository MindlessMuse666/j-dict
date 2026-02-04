package config

import (
	"os"
	"strconv"
	"time"
)

// Config хранит всю конфигурацию приложения
type Config struct {
	DB       DBConfig
	Server   ServerConfig
	Kafka    KafkaConfig
	JWT      JWTConfig
	LogLevel string
}

// DBConfig конфигурация базы данных
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// ServerConfig конфигурация HTTP сервера
type ServerConfig struct {
	Port string
}

// KafkaConfig конфигурация Kafka брокера
type KafkaConfig struct {
	Broker string
}

// JWTConfig конфигурация JWT токенов
type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

// Load загружает конфигурацию из переменных окружения
func Load() Config {
	return Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "app_user"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "jp_ru_dict"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Kafka: KafkaConfig{
			Broker: getEnv("KAFKA_BROKER", "localhost:9092"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "secret"),
			Expiry: getEnvAsDuration("JWT_EXPIRY", 24*time.Hour),
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv получает строковое значение из env или возвращает дефолтное
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt получает int значение из env или возвращает дефолтное
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// getEnvAsDuration получает time.Duration из env или возвращает дефолтное
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}
