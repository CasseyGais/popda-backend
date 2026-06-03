package tahap2

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"popda_bulutangkis/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// GET /admin/tahap2/export/pdf
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
	if len(data.NomorList) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Tidak ada data untuk di-export"})
		return
	}

	pdfBytes, err := GeneratePDF(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat file export"})
		return
	}

	namaKontingen := strings.ReplaceAll(data.Kontingen.NamaKontingen, " ", "_")
	tanggal := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("tahap2_%s_%s.pdf", namaKontingen, tanggal)

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/pdf")
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GET /admin/tahap2/export/excel
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
	if len(data.NomorList) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Tidak ada data untuk di-export"})
		return
	}

	xlsxBytes, err := GenerateExcel(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat file export"})
		return
	}

	namaKontingen := strings.ReplaceAll(data.Kontingen.NamaKontingen, " ", "_")
	tanggal := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("tahap2_%s_%s.xlsx", namaKontingen, tanggal)

	contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, contentType, xlsxBytes)
}
