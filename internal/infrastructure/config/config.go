package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration.
type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	ServerPort    string
	RabbitMQDSN   string
	RabbitMQQueue string
	RedisHost     string
	RedisPassword string
	RedisUsername string
	RedisPort     string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "postgres"),
		DBName:        getEnv("DB_NAME", "recommendation_db"),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		RabbitMQDSN:   getEnv("RABBITMQ_DSN", "amqp://guest:guest@localhost:5672/"),
		RabbitMQQueue: getEnv("RABBITMQ_QUEUE", "sender_recommendation_queue"),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisUsername: getEnv("REDIS_USERNAME", ""),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
	}
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

func (c *Config) RedisURL() string {
	return fmt.Sprintf("redis://%s:%s@%s:%s?protocol=3", c.RedisUsername, c.RedisPassword, c.RedisHost, c.RedisPort)
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
