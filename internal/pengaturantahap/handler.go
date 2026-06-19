package pengaturantahap

import (
	"net/http"
	"strconv"
	"strings"

	"popda_bulutangkis/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// isSuperadmin deteksi superadmin dari JWT claims.
// Case-insensitive — DB simpan "SUPERADMIN" tapi JWT bisa case apapun.
func isSuperadmin(claims *jwt.Claims) bool {
	return strings.ToLower(claims.Role) == "superadmin" || claims.KontingenID == 0
}

// GET /admin/pengaturan-tahap
// Ambil status semua tahap. Bisa diakses semua role.
// Dipakai frontend untuk tampilkan banner "Tahap X belum dibuka".
func (h *Handler) GetAll(c *gin.Context) {
	data, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data pengaturan tahap",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// PUT /admin/pengaturan-tahap/:tahap
// Update pengaturan satu tahap. Hanya superadmin.
// Body JSON (semua opsional): { "is_open": true, "tanggal_buka": "2026-06-01", "tanggal_tutup": "2026-06-30" }
func (h *Handler) Update(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	// Guard: hanya superadmin
	if !isSuperadmin(claims) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Hanya superadmin yang bisa mengubah pengaturan tahap",
		})
		return
	}

	tahapParam, err := strconv.ParseUint(c.Param("tahap"), 10, 32)
	if err != nil || tahapParam < 1 || tahapParam > 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter tahap tidak valid (harus 1, 2, atau 3)",
		})
		return
	}

	var req UpdatePengaturanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
		})
		return
	}

	result, err := h.service.Update(uint(tahapParam), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pengaturan Tahap " + strconv.Itoa(int(tahapParam)) + " berhasil diupdate",
		"data":    result,
	})
}
