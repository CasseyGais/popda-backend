package sertifikat

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ExportPayload adalah body request untuk endpoint export PDF sertifikat.
// Penandatangan opsional — jika kosong, PDF tampilkan garis tanda tangan kosong.
type ExportPayload struct {
	Penandatangan []TTDSertifikat `json:"penandatangan"`
}

// POST /admin/sertifikat/:id/export/pdf
// Generate PDF landscape satu sertifikat.
// Body JSON (opsional): { "penandatangan": [...] }
func (h *Handler) ExportPDF(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	// Body opsional — tidak apa-apa jika kosong
	var payload ExportPayload
	_ = c.ShouldBindJSON(&payload)

	data, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	item := SertifikatPDFData{
		NamaPenerima:  data.NamaPenerima,
		Judul:         data.Judul,
		TipePenerima:  data.TipePenerima,
		TanggalTerbit: data.TanggalTerbit,
		Penandatangan: payload.Penandatangan,
	}
	if data.NomorSertifikat != nil {
		item.NomorSertifikat = *data.NomorSertifikat
	}
	if data.Catatan != nil {
		item.Catatan = *data.Catatan
	}

	pdfBytes, err := GenerateSertifikatPDF([]SertifikatPDFData{item})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat file PDF"})
		return
	}

	namaFile := strings.ReplaceAll(data.NamaPenerima, " ", "_")
	tanggal := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("sertifikat_%s_%s.pdf", namaFile, tanggal)

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/pdf")
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// POST /admin/sertifikat/export/batch/pdf
// Generate PDF batch — satu file berisi semua/filtered sertifikat.
// Body JSON (opsional): { "penandatangan": [...] }
// Query params filter: tipe, atlet_id, pelatih_id, official_id
func (h *Handler) ExportBatchPDF(c *gin.Context) {
	// Body opsional
	var payload ExportPayload
	_ = c.ShouldBindJSON(&payload)

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

	list, err := h.service.GetAll(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data sertifikat"})
		return
	}
	if len(list) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Tidak ada data sertifikat untuk di-export"})
		return
	}

	items := make([]SertifikatPDFData, 0, len(list))
	for _, s := range list {
		item := SertifikatPDFData{
			NamaPenerima:  s.NamaPenerima,
			Judul:         s.Judul,
			TipePenerima:  s.TipePenerima,
			TanggalTerbit: s.TanggalTerbit,
			// Tanda tangan hanya di halaman terakhir — tidak di setiap halaman saat batch
			// Pass penandatangan hanya ke item terakhir
		}
		if s.NomorSertifikat != nil {
			item.NomorSertifikat = *s.NomorSertifikat
		}
		if s.Catatan != nil {
			item.Catatan = *s.Catatan
		}
		items = append(items, item)
	}

	// Taruh penandatangan hanya di item terakhir
	if len(payload.Penandatangan) > 0 && len(items) > 0 {
		items[len(items)-1].Penandatangan = payload.Penandatangan
	}

	pdfBytes, err := GenerateSertifikatPDF(items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat file PDF"})
		return
	}

	tanggal := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("sertifikat_batch_%s.pdf", tanggal)

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/pdf")
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
