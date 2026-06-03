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

// resolveKontingenID menentukan kontingen_id yang akan dipakai:
//   - Admin biasa  → pakai kontingen_id dari JWT (sudah terikat satu kontingen)
//   - Superadmin   → wajib kirim ?territory_id=X, backend cari kontingen dari territory tersebut
func (h *Handler) resolveKontingenID(c *gin.Context, claims *jwt.Claims) (uint, bool) {
	// Admin biasa — kontingen_id sudah ada di JWT
	if claims.KontingenID != 0 {
		return claims.KontingenID, true
	}

	// Superadmin — cari dari query param territory_id
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

// GET /admin/tahap1?territory_id=X
// Ambil status tahap1 + daftar cabor yang sudah dipilih kontingen.
// Superadmin wajib kirim query param territory_id.
// Admin biasa tidak perlu — kontingen_id diambil dari JWT.
func (h *Handler) Get(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	data, err := h.service.GetData(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data tahap 1 berhasil diambil",
		"data":    data,
	})
}

// PUT /admin/tahap1?territory_id=X
// Upsert satu cabor (tambah atau update kuota) via form-data.
// Superadmin wajib kirim query param territory_id.
func (h *Handler) Update(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	caborIDStr := c.PostForm("cabor_id")
	if caborIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "cabor_id wajib diisi",
		})
		return
	}

	caborID, err := strconv.ParseUint(caborIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "cabor_id tidak valid",
		})
		return
	}

	putra, _ := strconv.Atoi(c.PostForm("putra"))
	putri, _ := strconv.Atoi(c.PostForm("putri"))
	pelatih, _ := strconv.Atoi(c.PostForm("pelatih"))

	if err := h.service.UpsertCabor(kontingenID, uint(caborID), putra, putri, pelatih); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data cabor berhasil disimpan",
	})
}

// DELETE /admin/tahap1/:cabor_id?territory_id=X
// Hapus satu cabor dari daftar tahap 1.
// Superadmin wajib kirim query param territory_id.
func (h *Handler) DeleteCabor(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	caborID, err := strconv.ParseUint(c.Param("cabor_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "cabor_id tidak valid",
		})
		return
	}

	if err := h.service.DeleteCabor(kontingenID, uint(caborID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cabor berhasil dihapus dari daftar",
	})
}

// POST /admin/tahap1/submit?territory_id=X
// Kunci tahap 1 — ubah status ke SUBMITTED.
// Superadmin wajib kirim query param territory_id.
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
		"message": "Tahap 1 berhasil disubmit",
	})
}
