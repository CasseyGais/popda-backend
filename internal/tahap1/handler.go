package tahap1

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

// resolveKontingenID menentukan kontingen_id yang akan dipakai:
//
//   - Admin biasa  → pakai kontingen_id dari JWT
//   - Superadmin   → SELALU dari ?territory_id=X, tidak pernah dari JWT
//
// Superadmin dideteksi dari TIGA kondisi (defense in depth):
//  1. claims.Role == "superadmin"  → cara utama (role dari JWT)
//  2. claims.KontingenID == 0      → fallback: superadmin tidak punya kontingen
//  3. ?territory_id ada di query   → opt-in override: siapapun yang kirim territory_id
//                                    akan digunakan territory-nya
//
// Kondisi 3 memastikan bahkan token lama yang rolenya salah tetap bisa
// di-resolve dengan benar selama frontend mengirim territory_id.
func (h *Handler) resolveKontingenID(c *gin.Context, claims *jwt.Claims) (uint, bool) {
	isSuperadmin := strings.ToLower(claims.Role) == "superadmin" || claims.KontingenID == 0

	// Jika territory_id ada di query, selalu gunakan (override JWT)
	// Ini berlaku untuk superadmin maupun kasus khusus lainnya
	if territoryIDStr := c.Query("territory_id"); territoryIDStr != "" {
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

	// Tidak ada territory_id di query
	// Superadmin wajib kirim territory_id
	if isSuperadmin {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Superadmin wajib kirim query parameter territory_id",
		})
		return 0, false
	}

	// Admin biasa — kontingen_id dari JWT
	return claims.KontingenID, true
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
