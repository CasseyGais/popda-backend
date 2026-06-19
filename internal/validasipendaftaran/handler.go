package validasipendaftaran

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

// isSuperadmin deteksi superadmin — case-insensitive karena DB simpan "SUPERADMIN".
func isSuperadmin(claims *jwt.Claims) bool {
	return strings.ToLower(claims.Role) == "superadmin" || claims.KontingenID == 0
}

// resolveKontingenID resolve kontingen_id untuk endpoint yang bisa diakses semua role.
// Superadmin wajib kirim ?territory_id, admin biasa pakai JWT.
func (h *Handler) resolveKontingenID(c *gin.Context, claims *jwt.Claims) (uint, bool) {
	isSA := isSuperadmin(claims)

	if territoryIDStr := c.Query("territory_id"); territoryIDStr != "" {
		territoryID, err := strconv.ParseUint(territoryIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "territory_id tidak valid"})
			return 0, false
		}
		kontingenID, err := h.service.GetKontingenIDByTerritory(uint(territoryID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Kontingen untuk territory ini tidak ditemukan"})
			return 0, false
		}
		return kontingenID, true
	}

	if isSA {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Superadmin wajib kirim query parameter territory_id"})
		return 0, false
	}

	return claims.KontingenID, true
}

// GET /admin/validasi-pendaftaran/status 🔒 (Semua Role)
// Widget dashboard — hanya status validasi kontingen yang sedang login.
// Superadmin bisa baca kontingen tertentu dengan ?territory_id=X.
func (h *Handler) GetStatus(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	data, err := h.service.GetStatus(kontingenID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

// GET /admin/validasi-pendaftaran 🔒 (Superadmin Only)
// List semua kontingen + status validasi ketiga tahap.
// Query params opsional: ?status=PENDING&tahap=1&territory_id=2
func (h *Handler) GetAll(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	if !isSuperadmin(claims) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Hanya superadmin yang bisa mengakses halaman validasi pendaftaran",
		})
		return
	}

	filterStatus := c.Query("status")
	filterTahapStr := c.Query("tahap")
	filterTerritoryIDStr := c.Query("territory_id")

	var filterTahap int
	if filterTahapStr != "" {
		t, err := strconv.Atoi(filterTahapStr)
		if err == nil && t >= 1 && t <= 3 {
			filterTahap = t
		}
	}

	var filterTerritoryID uint
	if filterTerritoryIDStr != "" {
		t, err := strconv.ParseUint(filterTerritoryIDStr, 10, 32)
		if err == nil {
			filterTerritoryID = uint(t)
		}
	}

	data, err := h.service.GetAll(filterStatus, filterTahap, filterTerritoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data validasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

// GET /admin/rekap-pendaftaran 🔒 (Semua Role)
// Ambil semua data pendaftaran kontingen dalam satu response.
// Admin biasa: data kontingen sendiri. Superadmin: wajib ?territory_id=X.
func (h *Handler) GetRekap(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	data, err := h.service.GetRekap(kontingenID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Data rekap tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}
// Set status VALID atau REVISI untuk satu tahap satu kontingen.
func (h *Handler) SetValidasi(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	if !isSuperadmin(claims) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Hanya superadmin yang bisa mengubah status validasi",
		})
		return
	}

	kontingenID, err := strconv.ParseUint(c.Param("kontingen_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "kontingen_id tidak valid"})
		return
	}

	tahap, err := strconv.Atoi(c.Param("tahap"))
	if err != nil || tahap < 1 || tahap > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "tahap harus 1, 2, atau 3"})
		return
	}

	var req SetValidasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format request tidak valid"})
		return
	}

	result, err := h.service.SetValidasi(uint(kontingenID), tahap, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Validasi Tahap " + strconv.Itoa(tahap) + " kontingen " + result.NamaKontingen + " berhasil disimpan",
		"data":    result,
	})
}
