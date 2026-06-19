package masterofficial

import (
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

// resolveKontingenID menentukan kontingen_id yang akan dipakai.
// Tiga kondisi untuk deteksi superadmin (defense in depth):
//  1. claims.Role == "superadmin"
//  2. claims.KontingenID == 0 (token lama)
//  3. ?territory_id ada di query → selalu override JWT
func (h *Handler) resolveKontingenID(c *gin.Context, claims *jwt.Claims) (uint, bool) {
	isSuperadmin := claims.Role == "superadmin" || claims.KontingenID == 0

	// If territory_id is in query, always use it (overrides JWT)
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

	// No territory_id in query — superadmin must provide one
	if isSuperadmin {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Superadmin wajib kirim query parameter territory_id"})
		return 0, false
	}

	// Admin biasa — use kontingen_id from JWT
	return claims.KontingenID, true
}

// GET /admin/master/official?territory_id=X
func (h *Handler) GetAll(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	data, err := h.service.GetByKontingenID(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data official"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data official berhasil diambil",
		"data":    data,
	})
}

// GET /admin/master/official/:id
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	data, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Data official berhasil diambil", "data": data})
}

// POST /admin/master/official?territory_id=X
func (h *Handler) Create(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	var req CreateMasterOfficialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format request tidak valid", "error": err.Error()})
		return
	}

	data, err := h.service.Create(kontingenID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Official berhasil dibuat", "data": data})
}

// PUT /admin/master/official/:id
func (h *Handler) Update(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	var req UpdateMasterOfficialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format request tidak valid", "error": err.Error()})
		return
	}

	data, err := h.service.Update(uint(id), kontingenID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "official tidak ditemukan" {
			status = http.StatusNotFound
		} else if err.Error() == "tidak diizinkan mengubah data official kontingen lain" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Official berhasil diupdate", "data": data})
}

// DELETE /admin/master/official/:id
func (h *Handler) Delete(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	if err := h.service.Delete(uint(id), kontingenID); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "official tidak ditemukan" {
			status = http.StatusNotFound
		} else if err.Error() == "tidak diizinkan menghapus official kontingen lain" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Official berhasil dihapus"})
}

// PUT /admin/master/official/:id/foto
func (h *Handler) UpdateFoto(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	file, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File foto diperlukan"})
		return
	}

	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + file.Filename
	dst := filepath.Join("uploads", "official", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan file foto"})
		return
	}

	path := "/uploads/official/" + filename
	if err := h.service.UpdateFile(uint(id), kontingenID, "foto", path); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Foto official berhasil diupdate", "path": path})
}

// PUT /admin/master/official/:id/file/:kolom
// kolom: file_ktp | file_surat_tugas
func (h *Handler) UploadFile(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	kolom := c.Param("kolom")
	allowed := map[string]bool{"file_ktp": true, "file_surat_tugas": true}
	if !allowed[kolom] {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Nama kolom tidak valid"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File diperlukan"})
		return
	}

	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + file.Filename
	dst := filepath.Join("uploads", "official", "docs", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan file"})
		return
	}

	path := "/uploads/official/docs/" + filename
	if err := h.service.UpdateFile(uint(id), kontingenID, kolom, path); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "File berhasil diupload", "path": path})
}

// ===== TRX PENDAFTARAN =====

// GET /admin/master/official/trx?territory_id=X
func (h *Handler) GetTrx(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	data, err := h.service.GetTrxByKontingen(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data transaksi official"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Data transaksi official berhasil diambil", "data": data})
}

// POST /admin/master/official/trx — daftarkan official
// Body: { "official_id": 1 }
func (h *Handler) CreateTrx(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	var req struct {
		OfficialID uint `json:"official_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format request tidak valid", "error": err.Error()})
		return
	}

	data, err := h.service.CreateTrx(kontingenID, req.OfficialID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "official tidak ditemukan" {
			status = http.StatusNotFound
		} else if err.Error() == "tidak diizinkan mendaftarkan official kontingen lain" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Official berhasil didaftarkan", "data": data})
}

// DELETE /admin/master/official/trx/:id
func (h *Handler) DeleteTrx(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	if err := h.service.DeleteTrx(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Pendaftaran official berhasil dihapus"})
}
