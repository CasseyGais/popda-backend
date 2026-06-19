package sertifikat

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GET /admin/sertifikat/:id/export/pdf
// Generate PDF landscape satu sertifikat berdasarkan ID.
func (h *Handler) ExportPDF(c *gin.Context) {
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

	// Build data untuk PDF
	item := SertifikatPDFData{
		NamaPenerima:  data.NamaPenerima,
		Judul:         data.Judul,
		TipePenerima:  data.TipePenerima,
		TanggalTerbit: data.TanggalTerbit,
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

// GET /admin/sertifikat/export/batch/pdf?tipe=ATLET&atlet_id=4
// Generate PDF batch — satu file PDF berisi banyak halaman sertifikat.
// Query params sama dengan GET /admin/sertifikat (filter opsional).
func (h *Handler) ExportBatchPDF(c *gin.Context) {
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
		}
		if s.NomorSertifikat != nil {
			item.NomorSertifikat = *s.NomorSertifikat
		}
		if s.Catatan != nil {
			item.Catatan = *s.Catatan
		}
		items = append(items, item)
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
