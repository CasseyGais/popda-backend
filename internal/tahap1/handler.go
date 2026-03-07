package tahap1

import (
	"net/http"
	"strconv"

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

	caborStr := c.PostFormArray("caborList[]")
	var cabor []uint
	for _, str := range caborStr {
		if id, err := strconv.ParseUint(str, 10, 32); err == nil {
			cabor = append(cabor, uint(id))
		}
	}

	a, _ := strconv.Atoi(c.PostForm("jumlahAtlet"))
	p, _ := strconv.Atoi(c.PostForm("jumlahPelatih"))
	o, _ := strconv.Atoi(c.PostForm("jumlahOfficial"))

	if err := h.service.Update(kontingenID, cabor, a, p, o); err != nil {
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

	c.JSON(http.StatusOK, gin.H{"message": "Tahap 1 berhasil disubmit"})
}
