package main

import (
	"log"

	"github.com/aaanu/backend-qr-attendance/internal/config"
	"github.com/aaanu/backend-qr-attendance/internal/database"
	"github.com/aaanu/backend-qr-attendance/internal/handlers"
	"github.com/aaanu/backend-qr-attendance/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// Инициализация базы данных
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Инициализация сервисов
	qrService := services.NewQRService(db, cfg.BotUsername, cfg.QRLifetime)

	// Запуск генератора QR-кодов
	go qrService.StartQRGenerator()

	// Настройка роутера
	r := gin.Default()

	// Обработчики
	qrHandler := handlers.NewQRHandler(qrService)

	// Маршруты
	r.GET("/", qrHandler.ShowQRPage)
	r.GET("/qr/current", qrHandler.GetCurrentQR)
	r.POST("/api/verify/:uuid", qrHandler.VerifyQR)

	// Статические файлы
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("app/web/templates/*")

	log.Printf("Server starting on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
