package sertifikat

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GET /admin/sertifikat
// Query params opsional: tipe, atlet_id, pelatih_id, official_id
// Semua role yang diizinkan (SUPERADMIN / STAFF_LAPANGAN) bisa akses semua data.
func (h *Handler) GetAll(c *gin.Context) {
	filter := make(map[string]interface{})

	if tipe := c.Query("tipe"); tipe != "" {
		filter["tipe_penerima"] = strings.ToUpper(tipe)
	}
	if v := c.Query("atlet_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			filter["atlet_id"] = uint(id)
		}
	}
	if v := c.Query("pelatih_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			filter["pelatih_id"] = uint(id)
		}
	}
	if v := c.Query("official_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			filter["official_id"] = uint(id)
		}
	}

	data, err := h.service.GetAll(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data sertifikat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data sertifikat berhasil diambil",
		"data":    data,
	})
}

// GET /admin/sertifikat/:id
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

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Sertifikat berhasil diambil", "data": data})
}

// POST /admin/sertifikat
// Body JSON: { tipe_penerima, atlet_id/pelatih_id/official_id, judul, tanggal_terbit, ... }
// JANGAN kirim nama_penerima — diisi otomatis dari nama_lengkap tabel master.
func (h *Handler) Create(c *gin.Context) {
	var req CreateSertifikatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format request tidak valid", "error": err.Error()})
		return
	}

	data, err := h.service.Create(&req)
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "tidak ditemukan") {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Sertifikat berhasil dibuat",
		"data":    data,
	})
}

// PUT /admin/sertifikat/:id
// Partial update. NamaPenerima tidak bisa diubah.
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	var req UpdateSertifikatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format request tidak valid"})
		return
	}

	data, err := h.service.Update(uint(id), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "sertifikat tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Sertifikat berhasil diupdate", "data": data})
}

// DELETE /admin/sertifikat/:id
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "sertifikat tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Sertifikat berhasil dihapus"})
}

// PUT /admin/sertifikat/:id/file
// Upload file PDF sertifikat via multipart/form-data.
// Field: file
func (h *Handler) UploadFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File diperlukan (field: file)"})
		return
	}

	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + file.Filename
	dst := filepath.Join("uploads", "sertifikat", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan file"})
		return
	}

	path := "/uploads/sertifikat/" + filename
	if err := h.service.UpdateFile(uint(id), path); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "sertifikat tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File sertifikat berhasil diupload",
		"path":    path,
	})
}

// ===== DROPDOWN PENERIMA =====

// GET /admin/sertifikat/penerima/atlet
// Semua atlet dari semua kontingen untuk dropdown form buat sertifikat.
// Response: [{ id, nama_lengkap, kontingen_id, nama_kontingen }]
func (h *Handler) GetAtletDropdown(c *gin.Context) {
	data, err := h.service.GetAtletDropdown()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data atlet"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data atlet berhasil diambil",
		"data":    data,
	})
}

// GET /admin/sertifikat/penerima/pelatih
func (h *Handler) GetPelatihDropdown(c *gin.Context) {
	data, err := h.service.GetPelatihDropdown()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data pelatih"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pelatih berhasil diambil",
		"data":    data,
	})
}

// GET /admin/sertifikat/penerima/official
func (h *Handler) GetOfficialDropdown(c *gin.Context) {
	data, err := h.service.GetOfficialDropdown()
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
