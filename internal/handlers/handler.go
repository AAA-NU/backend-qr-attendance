package handlers

import (
	"net/http"

	"github.com/aaanu/backend-qr-attendance/internal/services"
	"github.com/gin-gonic/gin"
)

type QRHandler struct {
	qrService *services.QRService
}

func NewQRHandler(qrService *services.QRService) *QRHandler {
	return &QRHandler{qrService: qrService}
}

func (h *QRHandler) ShowQRPage(c *gin.Context) {
	qr, _ := h.qrService.GetCurrentQR()

	data := gin.H{
		"title": "QR-код для отметки посещаемости",
	}

	if qr != nil {
		data["qr_uuid"] = qr.UUID
		data["expires_at"] = qr.ExpiresAt.Format("15:04:05")
	}

	c.HTML(http.StatusOK, "index.html", data)
}

func (h *QRHandler) GetCurrentQR(c *gin.Context) {
	qr, qrImage := h.qrService.GetCurrentQR()

	if qr == nil || qrImage == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QR code not found"})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Data(http.StatusOK, "image/png", qrImage)
}

func (h *QRHandler) VerifyQR(c *gin.Context) {
	uuid := c.Param("uuid")

	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID is required"})
		return
	}

	isValid, err := h.qrService.VerifyQR(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": isValid,
		"uuid":  uuid,
	})
}
