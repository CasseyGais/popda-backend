package tahap3

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

// resolveKontingenID menentukan kontingen_id yang akan dipakai.
// Tiga kondisi untuk deteksi superadmin (defense in depth):
//  1. claims.Role == "superadmin"
//  2. claims.KontingenID == 0 (token lama)
//  3. ?territory_id ada di query → selalu override JWT
func (h *Handler) resolveKontingenID(c *gin.Context, claims *jwt.Claims) (uint, bool) {
	isSuperadmin := strings.ToLower(claims.Role) == "superadmin" || claims.KontingenID == 0

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

// ===== OVERVIEW =====

// GET /admin/tahap3?territory_id=X
func (h *Handler) Get(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok { return }

	data, err := h.service.GetData(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Data tahap 3 berhasil diambil", "data": data})
}

// ===== ATLET =====

// GET /admin/tahap3/atlet?territory_id=X
func (h *Handler) GetAtlets(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok { return }

	data, err := h.service.GetAtlets(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Data atlet berhasil diambil", "data": data})
}

// GET /admin/tahap3/atlet/:id
func (h *Handler) GetAtletByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}
	data, err := h.service.GetAtletByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

// POST /admin/tahap3/atlet?territory_id=X
func (h *Handler) CreateAtlet(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok { return }

	var req CreateAtletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format request tidak valid", "error": err.Error()})
		return
	}
	data, err := h.service.CreateAtlet(kontingenID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Atlet berhasil dibuat", "data": data})
}

// PUT /admin/tahap3/atlet/:id
func (h *Handler) UpdateAtlet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}
	var req UpdateAtletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format request tidak valid", "error": err.Error()})
		return
	}
	data, err := h.service.UpdateAtlet(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Atlet berhasil diupdate", "data": data})
}

// DELETE /admin/tahap3/atlet/:id
func (h *Handler) DeleteAtlet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}
	if err := h.service.DeleteAtlet(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Atlet berhasil dihapus"})
}

// PUT /admin/tahap3/atlet/:id/foto — upload foto atlet
func (h *Handler) UploadAtletFoto(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
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
	dst := filepath.Join("uploads", "atlet", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan file"})
		return
	}
	path := "/uploads/atlet/" + filename
	if err := h.service.repo.UpdateAtletFoto(uint(id), path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Foto atlet berhasil diupload", "path": path})
}

// PUT /admin/tahap3/atlet/:id/file/:kolom — upload dokumen atlet
// kolom: file_kartu_pelajar | file_akte_kelahiran | file_kk | file_surat_keterangan_sekolah | file_surat_izin_ortu
func (h *Handler) UploadAtletFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}
	kolom := c.Param("kolom")
	allowed := map[string]bool{
		"file_kartu_pelajar": true, "file_akte_kelahiran": true,
		"file_kk": true, "file_surat_keterangan_sekolah": true, "file_surat_izin_ortu": true,
	}
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
	dst := filepath.Join("uploads", "atlet", "docs", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan file"})
		return
	}
	path := "/uploads/atlet/docs/" + filename
	if err := h.service.repo.UpdateAtletFile(uint(id), kolom, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "File berhasil diupload", "path": path})
}

// ===== PELATIH =====

// GET /admin/tahap3/pelatih?territory_id=X
func (h *Handler) GetPelatihs(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok { return }

	data, err := h.service.GetPelatihs(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Data pelatih berhasil diambil", "data": data})
}

// GET /admin/tahap3/pelatih/:id
func (h *Handler) GetPelatihByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}
	data, err := h.service.GetPelatihByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

// POST /admin/tahap3/pelatih?territory_id=X
func (h *Handler) CreatePelatih(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok { return }

	var req CreatePelatihRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format request tidak valid", "error": err.Error()})
		return
	}
	data, err := h.service.CreatePelatih(kontingenID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Pelatih berhasil dibuat", "data": data})
}

// PUT /admin/tahap3/pelatih/:id
func (h *Handler) UpdatePelatih(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}
	var req UpdatePelatihRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format request tidak valid", "error": err.Error()})
		return
	}
	data, err := h.service.UpdatePelatih(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Pelatih berhasil diupdate", "data": data})
}

// DELETE /admin/tahap3/pelatih/:id
func (h *Handler) DeletePelatih(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}
	if err := h.service.DeletePelatih(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Pelatih berhasil dihapus"})
}

// PUT /admin/tahap3/pelatih/:id/file/:kolom — upload dokumen pelatih
// kolom: foto | file_ktp | file_surat_tugas | file_sertifikat_pelatih
func (h *Handler) UploadPelatihFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}
	kolom := c.Param("kolom")
	allowed := map[string]bool{"foto": true, "file_ktp": true, "file_surat_tugas": true, "file_sertifikat_pelatih": true}
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
	dst := filepath.Join("uploads", "pelatih", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan file"})
		return
	}
	path := "/uploads/pelatih/" + filename
	if err := h.service.repo.UpdatePelatihFile(uint(id), kolom, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "File berhasil diupload", "path": path})
}

// ===== OFFICIAL =====

// GET /admin/tahap3/official?territory_id=X
func (h *Handler) GetOfficials(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok { return }

	data, err := h.service.GetOfficials(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Data official berhasil diambil", "data": data})
}

// GET /admin/tahap3/official/:id
func (h *Handler) GetOfficialByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}
	data, err := h.service.GetOfficialByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

// POST /admin/tahap3/official?territory_id=X
func (h *Handler) CreateOfficial(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok { return }

	var req CreateOfficialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format request tidak valid", "error": err.Error()})
		return
	}
	data, err := h.service.CreateOfficial(kontingenID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Official berhasil dibuat", "data": data})
}

// PUT /admin/tahap3/official/:id
func (h *Handler) UpdateOfficial(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}
	var req UpdateOfficialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format request tidak valid", "error": err.Error()})
		return
	}
	data, err := h.service.UpdateOfficial(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Official berhasil diupdate", "data": data})
}

// DELETE /admin/tahap3/official/:id
func (h *Handler) DeleteOfficial(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}
	if err := h.service.DeleteOfficial(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Official berhasil dihapus"})
}

// PUT /admin/tahap3/official/:id/file/:kolom — upload dokumen official
// kolom: foto | file_ktp | file_surat_tugas
func (h *Handler) UploadOfficialFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}
	kolom := c.Param("kolom")
	allowed := map[string]bool{"foto": true, "file_ktp": true, "file_surat_tugas": true}
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
	dst := filepath.Join("uploads", "official", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan file"})
		return
	}
	path := "/uploads/official/" + filename
	if err := h.service.repo.UpdateOfficialFile(uint(id), kolom, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "File berhasil diupload", "path": path})
}

// ===== STATISTIK =====

// GET /admin/tahap3/statistik/atlet
// Hitung jumlah atlet dari seluruh kontingen.
// Bisa diakses semua role — superadmin tidak perlu territory_id.
func (h *Handler) GetStatistikAtlet(c *gin.Context) {
	data, err := h.service.GetStatistikAtlet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghitung statistik atlet",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Statistik atlet berhasil diambil",
		"data":    data,
	})
}

// ===== TRX PENDAFTARAN =====

// POST /admin/tahap3/trx/atlet
func (h *Handler) CreateTrxAtlet(c *gin.Context) {
	var req CreateTrxAtletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format request tidak valid", "error": err.Error()})
		return
	}
	data, err := h.service.CreateTrxAtlet(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Atlet berhasil didaftarkan", "data": data})
}

// DELETE /admin/tahap3/trx/atlet/:id
func (h *Handler) DeleteTrxAtlet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}
	if err := h.service.DeleteTrxAtlet(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Pendaftaran atlet berhasil dihapus"})
}

// POST /admin/tahap3/trx/pelatih
func (h *Handler) CreateTrxPelatih(c *gin.Context) {
	var req CreateTrxPelatihRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format request tidak valid", "error": err.Error()})
		return
	}
	data, err := h.service.CreateTrxPelatih(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Pelatih berhasil didaftarkan", "data": data})
}

// DELETE /admin/tahap3/trx/pelatih/:id
func (h *Handler) DeleteTrxPelatih(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}
	if err := h.service.DeleteTrxPelatih(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Pendaftaran pelatih berhasil dihapus"})
}

// POST /admin/tahap3/trx/official
func (h *Handler) CreateTrxOfficial(c *gin.Context) {
	var req CreateTrxOfficialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format request tidak valid", "error": err.Error()})
		return
	}
	data, err := h.service.CreateTrxOfficial(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Official berhasil didaftarkan", "data": data})
}

// DELETE /admin/tahap3/trx/official/:id
func (h *Handler) DeleteTrxOfficial(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}
	if err := h.service.DeleteTrxOfficial(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Pendaftaran official berhasil dihapus"})
}

// POST /admin/tahap3/submit?territory_id=X
// Kunci tahap 3:
// - Otomatis bulk insert semua atlet/pelatih/official ke trx_pendaftaran_*
// - Set tahap3_status = SUBMITTED di tabel kontingen
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
		"message": "Tahap 3 berhasil disubmit. Semua atlet, pelatih, dan official telah didaftarkan.",
	})
}

// ===== REFERENSI TAHAP SEBELUMNYA =====

// GET /admin/tahap3/cabor?territory_id=X
// Ambil daftar cabor yang dipilih kontingen di tahap 1.
// Dipakai frontend sebagai filter cabor saat input atlet/pelatih.
func (h *Handler) GetCaborTerpilih(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	data, err := h.service.GetCaborTerpilih(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar cabor terpilih berhasil diambil",
		"data":    data,
	})
}

// GET /admin/tahap3/nomor?territory_id=X
// Ambil daftar nomor pertandingan yang dicentang kontingen di tahap 2.
// Dipakai frontend sebagai dropdown saat assign atlet ke nomor.
func (h *Handler) GetNomorTerdaftar(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	data, err := h.service.GetNomorTerdaftar(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar nomor terdaftar berhasil diambil",
		"data":    data,
	})
}
