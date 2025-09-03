package services

import (
	"fmt"
	"log"

	"sync"
	"time"

	"github.com/aaanu/backend-qr-attendance/internal/models"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type QRService struct {
	db          *gorm.DB
	botUsername string
	lifetime    time.Duration
	currentQR   *models.QRCode
	qrImage     []byte
	mutex       sync.RWMutex
}

func NewQRService(db *gorm.DB, botUsername string, lifetime time.Duration) *QRService {
	return &QRService{
		db:          db,
		botUsername: botUsername,
		lifetime:    lifetime,
	}
}

func (s *QRService) StartQRGenerator() {
	// Генерируем первый QR-код сразу
	s.generateNewQR()

	// Запускаем таймер для регулярной генерации
	ticker := time.NewTicker(s.lifetime)
	defer ticker.Stop()

	for range ticker.C {
		s.generateNewQR()
	}
}

func (s *QRService) generateNewQR() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Деактивируем старые QR-коды
	s.db.Model(&models.QRCode{}).Where("is_active = ?", true).Update("is_active", false)

	// Очищаем таблицу от старых QR-кодов
	var count int64
	if err := s.db.Model(&models.QRCode{}).Count(&count).Error; err != nil {
		log.Printf("Error counting QR codes: %v", err)
		return
	}

	if count > 50 {
		if err := s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.QRCode{}).Error; err != nil {
			log.Printf("Error deleting QR codes: %v", err)
			return
		}
	}

	// Создаем новый QR-код
	qr := &models.QRCode{
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.lifetime),
		IsActive:  true,
	}

	if err := s.db.Create(qr).Error; err != nil {
		log.Printf("Error creating QR code: %v", err)
		return
	}

	// Генерируем ссылку для Telegram бота
	botURL := fmt.Sprintf("https://t.me/%s?start=%s", s.botUsername, qr.UUID)

	// Генерируем QR-код изображение
	qrImg, err := qrcode.Encode(botURL, qrcode.Medium, 256)
	if err != nil {
		log.Printf("Error generating QR image: %v", err)
		return
	}

	s.currentQR = qr
	s.qrImage = qrImg

	log.Printf("Generated new QR code: %s", qr.UUID)
}

func (s *QRService) GetCurrentQR() (*models.QRCode, []byte) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.currentQR, s.qrImage
}

func (s *QRService) VerifyQR(uuid string) (bool, error) {
	var qr models.QRCode

	err := s.db.Where("uuid = ? AND is_active = ?", uuid, true).First(&qr).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	// Проверяем, не истек ли QR-код
	if qr.IsExpired() {
		// Деактивируем истекший QR-код
		s.db.Model(&qr).Update("is_active", false)
		return false, nil
	}

	return true, nil
}

// GetQRStats возвращает статистику по QR-кодам (для дополнительного функционала)
func (s *QRService) GetQRStats() (int64, error) {
	var count int64
	err := s.db.Model(&models.QRCode{}).Count(&count).Error
	return count, err
}
