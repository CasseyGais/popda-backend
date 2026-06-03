package tahap2

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

// resolveKontingenID sama dengan tahap1 — superadmin wajib kirim ?territory_id=X
func (h *Handler) resolveKontingenID(c *gin.Context, claims *jwt.Claims) (uint, bool) {
	if claims.KontingenID != 0 {
		return claims.KontingenID, true
	}

	territoryIDStr := c.Query("territory_id")
	if territoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Superadmin wajib kirim query parameter territory_id",
		})
		return 0, false
	}

	territoryID, err := strconv.ParseUint(territoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "territory_id tidak valid",
		})
		return 0, false
	}

	kontingenID, err := h.service.GetKontingenIDByTerritory(uint(territoryID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Kontingen untuk territory ini tidak ditemukan",
		})
		return 0, false
	}

	return kontingenID, true
}

// GET /admin/tahap2?territory_id=X
// Ambil status tahap2 + daftar nomor dari cabor yang dipilih di tahap 1.
// Superadmin wajib kirim query param territory_id.
func (h *Handler) Get(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	data, err := h.service.GetData(kontingenID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data tahap 2 berhasil diambil",
		"data":    data,
	})
}

// POST /admin/tahap2/nomor/:nomor_id?territory_id=X
// Daftarkan satu nomor ke kontingen.
func (h *Handler) DaftarNomor(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	nomorID, err := strconv.ParseUint(c.Param("nomor_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "nomor_id tidak valid",
		})
		return
	}

	if err := h.service.DaftarNomor(kontingenID, uint(nomorID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Nomor berhasil didaftarkan",
	})
}

// DELETE /admin/tahap2/nomor/:nomor_id?territory_id=X
// Batalkan pendaftaran satu nomor dari kontingen.
func (h *Handler) BatalNomor(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	nomorID, err := strconv.ParseUint(c.Param("nomor_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "nomor_id tidak valid",
		})
		return
	}

	if err := h.service.BatalNomor(kontingenID, uint(nomorID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pendaftaran nomor berhasil dibatalkan",
	})
}

// POST /admin/tahap2/submit?territory_id=X
// Kunci tahap 2 — ubah status ke SUBMITTED.
func (h *Handler) Submit(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	if err := h.service.Submit(kontingenID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tahap 2 berhasil disubmit",
	})
}
