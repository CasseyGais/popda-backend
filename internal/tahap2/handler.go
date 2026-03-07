package tahap2

import (
	"net/http"

	"popda_bulutangkis/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Get(c *gin.Context) {

	claims := c.MustGet("user").(*jwt.Claims)

	// TODO: Get kontingenID dari user_territories mapping
	kontingenID := claims.UserID // Sementara - akan diperbaiki dengan auth yang proper

	data, err := h.service.GetData(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handler) Update(c *gin.Context) {

	claims := c.MustGet("user").(*jwt.Claims)

	// TODO: Get kontingenID dari user_territories mapping
	kontingenID := claims.UserID // Sementara - akan diperbaiki dengan auth yang proper

	var req struct {
		NomorIDs []uint `json:"nomor_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid"})
		return
	}

	if err := h.service.Update(kontingenID, req.NomorIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil disimpan"})
}

func (h *Handler) Submit(c *gin.Context) {

	claims := c.MustGet("user").(*jwt.Claims)

	// TODO: Get kontingenID dari user_territories mapping
	kontingenID := claims.UserID // Sementara - akan diperbaiki dengan auth yang proper

	if err := h.service.Submit(kontingenID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tahap 2 berhasil disubmit"})
}
