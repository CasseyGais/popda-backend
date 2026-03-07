package tahap3

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

	// TODO: Get kontingenID from user_territories mapping
	kontingenID := claims.UserID // Temporary - will be fixed with proper auth

	data, err := h.service.GetData(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handler) Submit(c *gin.Context) {

	claims := c.MustGet("user").(*jwt.Claims)

	// TODO: Get kontingenID from user_territories mapping
	kontingenID := claims.UserID // Temporary - will be fixed with proper auth

	if err := h.service.Submit(kontingenID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tahap 3 berhasil disubmit"})
}
