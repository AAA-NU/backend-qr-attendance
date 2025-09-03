package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	BotUsername string
	QRLifetime  time.Duration
}

func Load() *Config {
	godotenv.Load(".env")
	port := getEnv("PORT", "8080")
	databaseURL := getEnv("DATABASE_URL", "postgres://username:password@localhost:5432/qr_attendance?sslmode=disable")
	botUsername := getEnv("BOT_USERNAME", "your_bot_username")

	lifetimeStr := getEnv("QR_LIFETIME_SECONDS", "12")
	lifetime, _ := strconv.Atoi(lifetimeStr)

	return &Config{
		Port:        port,
		DatabaseURL: databaseURL,
		BotUsername: botUsername,
		QRLifetime:  time.Duration(lifetime) * time.Second,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
