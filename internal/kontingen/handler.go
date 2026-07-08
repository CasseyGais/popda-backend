package kontingen

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"popda_bulutangkis/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ================= GET =================
func (h *Handler) GetIdentitas(c *gin.Context) {
	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims := userData.(*jwt.Claims)
	kontingenID := claims.KontingenID // ✅ Pakai KontingenID, bukan UserID

	data, err := h.service.GetIdentitas(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data identitas",
			"error":   err.Error(),
		})
		return
	}

	if data == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Data identitas kosong",
			"data":    gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data identitas berhasil diambil",
		"data":    data,
	})
}

// ================= UPDATE =================
func (h *Handler) Update(c *gin.Context) {
	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims := userData.(*jwt.Claims)
	kontingenID := claims.KontingenID // ✅ Pakai KontingenID, bukan UserID

	input := IdentitasKontingen{
		// ===== INSTANSI =====
		Alamat:        c.PostForm("alamat"),
		EmailInstansi: c.PostForm("email_instansi"),
		PhoneInstansi: c.PostForm("phone_instansi"),

		// ===== KEPALA  =====
		KepalaNama:    c.PostForm("kepala_nama"),
		KepalaJabatan: c.PostForm("kepala_jabatan"),
		KepalaNIP:     c.PostForm("kepala_nip"),
		KepalaTelepon: c.PostForm("kepala_telepon"),

		// ===== PIC =====
		PICNama:    c.PostForm("pic_nama"),
		PICJabatan: c.PostForm("pic_jabatan"),
		PICTelepon: c.PostForm("pic_telepon"),
	}

	// ===== Upload Foto Kepala =====
	if file, err := c.FormFile("kepala_foto"); err == nil {
		filename := time.Now().Format("20060102150405") + "_" + file.Filename
		dst := filepath.Join("uploads", "kepala", filename)

		if err := c.SaveUploadedFile(file, dst); err == nil {
			input.KepalaFoto = "/uploads/kepala/" + filename
		}
	}

	// ===== Upload Foto PIC =====
	if file, err := c.FormFile("pic_foto"); err == nil {
		filename := time.Now().Format("20060102150405") + "_" + file.Filename
		dst := filepath.Join("uploads", "pic", filename)

		if err := c.SaveUploadedFile(file, dst); err == nil {
			input.PICFoto = "/uploads/pic/" + filename
		}
	}

	if err := h.service.UpdateIdentitas(kontingenID, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengupdate data identitas",
			"error":   err.Error(),
		})
		return
	}

	updated, _ := h.service.GetIdentitas(kontingenID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Identitas kontingen berhasil diperbarui",
		"data":    updated,
	})
}

// ================= KONTINGEN HANDLERS =================
func (h *Handler) GetKontingen(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kontingen diperlukan"})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	kontingen, err := h.service.GetKontingenByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if kontingen == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kontingen tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kontingen})
}

func (h *Handler) GetKontingenByTerritory(c *gin.Context) {
	territoryIDStr := c.Param("territory_id")
	if territoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Territory ID diperlukan"})
		return
	}

	var territoryID uint
	if _, err := fmt.Sscanf(territoryIDStr, "%d", &territoryID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Territory ID tidak valid"})
		return
	}

	kontingen, err := h.service.GetKontingenByTerritoryID(territoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if kontingen == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kontingen untuk territory ini tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kontingen})
}

func (h *Handler) CreateKontingen(c *gin.Context) {
	var input Kontingen

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set CreatedAt ke current time
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	if err := h.service.CreateKontingen(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kontingen berhasil dibuat",
		"data":    input,
	})
}

func (h *Handler) UpdateKontingen(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kontingen diperlukan"})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var input Kontingen
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set UpdatedAt ke current time
	input.UpdatedAt = time.Now()

	if err := h.service.UpdateKontingen(id, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, _ := h.service.GetKontingenByID(id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Kontingen berhasil diperbarui",
		"data":    updated,
	})
}

// ================= RESET TAHAP =================

// resolveKontingenIDForReset meng-resolve kontingen_id dari query ?territory_id=X
// (superadmin) atau dari JWT claims (user biasa).
// Endpoint reset hanya bisa diakses superadmin, tapi helper ini tetap generik
// agar mudah dipakai ulang.
func (h *Handler) resolveKontingenIDForReset(c *gin.Context, claims *jwt.Claims) (uint, bool) {
	if territoryIDStr := c.Query("territory_id"); territoryIDStr != "" {
		tid, err := strconv.ParseUint(territoryIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "territory_id tidak valid",
			})
			return 0, false
		}

		k, err := h.service.GetKontingenByTerritoryID(uint(tid))
		if err != nil || k == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Kontingen untuk territory ini tidak ditemukan",
			})
			return 0, false
		}
		return k.ID, true
	}

	// Fallback ke KontingenID dari JWT
	if claims.KontingenID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Superadmin wajib kirim query parameter territory_id",
		})
		return 0, false
	}
	return claims.KontingenID, true
}

// POST /admin/tahap1/reset — hanya SUPERADMIN
func (h *Handler) ResetTahap1(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	kontingenID, ok := h.resolveKontingenIDForReset(c, claims)
	if !ok {
		return
	}

	if err := h.service.ResetTahap(kontingenID, 1); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mereset tahap 1",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tahap 1 berhasil direset ke DRAFT",
	})
}

// POST /admin/tahap2/reset — hanya SUPERADMIN
func (h *Handler) ResetTahap2(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	kontingenID, ok := h.resolveKontingenIDForReset(c, claims)
	if !ok {
		return
	}

	if err := h.service.ResetTahap(kontingenID, 2); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mereset tahap 2",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tahap 2 berhasil direset ke DRAFT",
	})
}

// POST /admin/tahap3/reset — hanya SUPERADMIN
func (h *Handler) ResetTahap3(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	kontingenID, ok := h.resolveKontingenIDForReset(c, claims)
	if !ok {
		return
	}

	if err := h.service.ResetTahap(kontingenID, 3); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mereset tahap 3",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tahap 3 berhasil direset ke DRAFT",
	})
}
