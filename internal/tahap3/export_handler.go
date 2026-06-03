package tahap3

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"popda_bulutangkis/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const xlsxContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

// GET /admin/tahap3/export/pdf  — semua: atlet + pelatih + official
func (h *Handler) ExportPDF(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	data, err := h.service.GetExportData(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if len(data.Atlets) == 0 && len(data.Pelatihs) == 0 && len(data.Officials) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Tidak ada data untuk di-export"})
		return
	}

	pdfBytes, err := GeneratePDF(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat file export"})
		return
	}

	filename := buildFilename("tahap3", data.Kontingen.NamaKontingen, "pdf")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/pdf")
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GET /admin/tahap3/export/excel  — 3 sheet: Atlet, Pelatih, Official
func (h *Handler) ExportExcel(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	data, err := h.service.GetExportData(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if len(data.Atlets) == 0 && len(data.Pelatihs) == 0 && len(data.Officials) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Tidak ada data untuk di-export"})
		return
	}

	xlsxBytes, err := GenerateExcel(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat file export"})
		return
	}

	filename := buildFilename("tahap3", data.Kontingen.NamaKontingen, "xlsx")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", xlsxContentType)
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, xlsxContentType, xlsxBytes)
}

// GET /admin/tahap3/export/atlet/pdf
func (h *Handler) ExportAtletPDF(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	kontingen, atlets, err := h.service.GetExportAtlets(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if len(atlets) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Tidak ada data untuk di-export"})
		return
	}

	pdfBytes, err := GeneratePDFAtlet(kontingen, atlets)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat file export"})
		return
	}

	filename := buildFilename("tahap3_atlet", kontingen.NamaKontingen, "pdf")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/pdf")
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GET /admin/tahap3/export/atlet/excel
func (h *Handler) ExportAtletExcel(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	kontingen, atlets, err := h.service.GetExportAtlets(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if len(atlets) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Tidak ada data untuk di-export"})
		return
	}

	xlsxBytes, err := GenerateExcelAtlet(kontingen, atlets)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat file export"})
		return
	}

	filename := buildFilename("tahap3_atlet", kontingen.NamaKontingen, "xlsx")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", xlsxContentType)
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, xlsxContentType, xlsxBytes)
}

// GET /admin/tahap3/export/pelatih/pdf
func (h *Handler) ExportPelatihPDF(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	kontingen, pelatihs, err := h.service.GetExportPelatihs(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if len(pelatihs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Tidak ada data untuk di-export"})
		return
	}

	pdfBytes, err := GeneratePDFPelatih(kontingen, pelatihs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat file export"})
		return
	}

	filename := buildFilename("tahap3_pelatih", kontingen.NamaKontingen, "pdf")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/pdf")
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GET /admin/tahap3/export/pelatih/excel
func (h *Handler) ExportPelatihExcel(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	kontingen, pelatihs, err := h.service.GetExportPelatihs(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if len(pelatihs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Tidak ada data untuk di-export"})
		return
	}

	xlsxBytes, err := GenerateExcelPelatih(kontingen, pelatihs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat file export"})
		return
	}

	filename := buildFilename("tahap3_pelatih", kontingen.NamaKontingen, "xlsx")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", xlsxContentType)
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, xlsxContentType, xlsxBytes)
}

// GET /admin/tahap3/export/official/pdf
func (h *Handler) ExportOfficialPDF(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	kontingen, officials, err := h.service.GetExportOfficials(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if len(officials) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Tidak ada data untuk di-export"})
		return
	}

	pdfBytes, err := GeneratePDFOfficial(kontingen, officials)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat file export"})
		return
	}

	filename := buildFilename("tahap3_official", kontingen.NamaKontingen, "pdf")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/pdf")
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GET /admin/tahap3/export/official/excel
func (h *Handler) ExportOfficialExcel(c *gin.Context) {
	claims := c.MustGet("user").(*jwt.Claims)
	kontingenID, ok := h.resolveKontingenID(c, claims)
	if !ok {
		return
	}

	kontingen, officials, err := h.service.GetExportOfficials(kontingenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if len(officials) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Tidak ada data untuk di-export"})
		return
	}

	xlsxBytes, err := GenerateExcelOfficial(kontingen, officials)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat file export"})
		return
	}

	filename := buildFilename("tahap3_official", kontingen.NamaKontingen, "xlsx")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", xlsxContentType)
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, xlsxContentType, xlsxBytes)
}

// buildFilename membuat nama file sesuai konvensi dokumentasi
// format: <prefix>_<nama_kontingen_underscore>_<YYYY-MM-DD>.<ext>
func buildFilename(prefix, namaKontingen, ext string) string {
	slug := strings.ReplaceAll(namaKontingen, " ", "_")
	tanggal := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s_%s_%s.%s", prefix, slug, tanggal, ext)
}
