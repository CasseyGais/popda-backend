package laporanpertandingan

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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

// GET /admin/laporan-pertandingan
// Query params opsional: cabor_id, nomor_id, babak, tanggal, pemenang
func (h *Handler) GetAll(c *gin.Context) {
	filter := FilterLaporan{}

	if v := c.Query("cabor_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			filter.CaborID = uint(id)
		}
	}
	if v := c.Query("nomor_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			filter.NomorID = uint(id)
		}
	}
	if v := c.Query("babak"); v != "" {
		filter.Babak = strings.ToUpper(v)
	}
	if v := c.Query("tanggal"); v != "" {
		filter.Tanggal = v
	}
	if v := c.Query("pemenang"); v != "" {
		filter.Pemenang = strings.ToUpper(v)
	}

	data, err := h.service.GetAll(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data laporan pertandingan berhasil diambil",
		"data":    data,
	})
}

// GET /admin/laporan-pertandingan/:id
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Laporan pertandingan berhasil diambil",
		"data":    data,
	})
}

// POST /admin/laporan-pertandingan
// Body JSON — created_by otomatis dari JWT user_id
func (h *Handler) Create(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)

	var req CreateLaporanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	data, err := h.service.Create(&req, claims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Laporan pertandingan berhasil dibuat",
		"data":    data,
	})
}

// PUT /admin/laporan-pertandingan/:id
// Partial update — hanya field yang dikirim yang berubah
// Jika atlet_a/atlet_b dikirim, isinya akan me-replace semua atlet sisi tersebut
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	var req UpdateLaporanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format request tidak valid"})
		return
	}

	data, err := h.service.Update(uint(id), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "laporan pertandingan tidak ditemukan" {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "tidak valid") {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Laporan pertandingan berhasil diupdate",
		"data":    data,
	})
}

// DELETE /admin/laporan-pertandingan/:id
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "laporan pertandingan tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Laporan pertandingan berhasil dihapus"})
}

// PUT /admin/laporan-pertandingan/:id/foto
// Upload foto bukti (papan skor, momen penting)
// Field form: foto
func (h *Handler) UploadFoto(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	file, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File foto diperlukan (field: foto)"})
		return
	}

	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + file.Filename
	dst := filepath.Join("uploads", "laporan", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan file foto"})
		return
	}

	path := "/uploads/laporan/" + filename
	if err := h.service.UpdateFile(uint(id), "foto_bukti", path); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "laporan pertandingan tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Foto bukti berhasil diupload",
		"path":    path,
	})
}

// PUT /admin/laporan-pertandingan/:id/video
// Upload video bukti (opsional)
// Field form: video
func (h *Handler) UploadVideo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File video diperlukan (field: video)"})
		return
	}

	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + file.Filename
	dst := filepath.Join("uploads", "laporan", "video", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan file video"})
		return
	}

	path := "/uploads/laporan/video/" + filename
	if err := h.service.UpdateFile(uint(id), "video_bukti", path); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "laporan pertandingan tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Video bukti berhasil diupload",
		"path":    path,
	})
}

// ===== DROPDOWN ENDPOINTS =====

// GET /admin/laporan-pertandingan/dropdown/kontingen
// Semua kontingen — untuk dropdown Tim A dan Tim B
func (h *Handler) GetKontingenDropdown(c *gin.Context) {
	data, err := h.service.GetKontingenDropdown()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data kontingen"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

// GET /admin/laporan-pertandingan/dropdown/cabor
// Cabor aktif — untuk dropdown cabang olahraga
func (h *Handler) GetCaborDropdown(c *gin.Context) {
	data, err := h.service.GetCaborDropdown()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data cabor"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

// GET /admin/laporan-pertandingan/dropdown/nomor?cabor_id=6
// Nomor/kelas aktif — bisa difilter by cabor_id
func (h *Handler) GetNomorDropdown(c *gin.Context) {
	var caborID uint
	if v := c.Query("cabor_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			caborID = uint(id)
		}
	}
	data, err := h.service.GetNomorDropdown(caborID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data nomor"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

// GET /admin/laporan-pertandingan/dropdown/atlet?kontingen_id=2&cabor_id=6&nomor_id=70
// Atlet dari trx_pendaftaran_atlet — hanya atlet yang sudah terdaftar di cabor+nomor tertentu
// Filter: kontingen_id, cabor_id, nomor_id — semua opsional
func (h *Handler) GetAtletDropdown(c *gin.Context) {
	var kontingenID, caborID, nomorID uint
	if v := c.Query("kontingen_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			kontingenID = uint(id)
		}
	}
	if v := c.Query("cabor_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			caborID = uint(id)
		}
	}
	if v := c.Query("nomor_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			nomorID = uint(id)
		}
	}

	data, err := h.service.GetAtletTerdaftarDropdown(kontingenID, caborID, nomorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data atlet"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}
