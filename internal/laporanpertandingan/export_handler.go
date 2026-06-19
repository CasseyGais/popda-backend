package laporanpertandingan

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)
// ===== POST /admin/laporan-pertandingan/:id/export/pdf =====
// Export PDF satu pertandingan.
// Body JSON: { "penandatangan": [...] }
// Tanda tangan lewat base64 PNG dari signature pad di frontend.
func (h *Handler) ExportSatuPDF(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	var req ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Body opsional — tidak apa-apa jika kosong (tidak ada tanda tangan)
		req = ExportRequest{}
	}

	laporan, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	pdfBytes, err := GenerateSatuPDF(laporan, req.Penandatangan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat PDF"})
		return
	}

	filename := fmt.Sprintf(
		"laporan_%d_%s_%s.pdf",
		laporan.ID,
		strings.ReplaceAll(laporan.NamaCabor, " ", "_"),
		laporan.TanggalPertandingan.Format("2006-01-02"),
	)

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/pdf")
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// ===== POST /admin/laporan-pertandingan/export/pdf =====
// Export PDF batch — per hari atau semua pertandingan.
//
// Body JSON:
// {
//   "tanggal":  "2026-08-10",   ← opsional, filter per hari
//   "cabor_id": 6,              ← opsional
//   "nomor_id": 70,             ← opsional
//   "penandatangan": [
//     {
//       "jabatan":       "Wasit",
//       "nama_tercetak": "Budi Santoso",
//       "nip":           "198001012010011001",
//       "signature_b64": "data:image/png;base64,iVBOR..."  ← dari signature pad
//     },
//     {
//       "jabatan":       "Ketua Panitia",
//       "nama_tercetak": "Rudi Hartono",
//       "nip":           "",
//       "signature_b64": ""
//     }
//   ]
// }
func (h *Handler) ExportBatchPDF(c *gin.Context) {
	var req ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req = ExportRequest{}
	}

	// Build filter dari request
	filter := FilterLaporan{
		CaborID: req.CaborID,
		NomorID: req.NomorID,
		Tanggal: req.Tanggal,
	}

	items, err := h.service.GetAll(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data laporan"})
		return
	}
	if len(items) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Tidak ada data laporan untuk di-export"})
		return
	}

	// Judul dokumen
	var judulParts []string
	judulParts = append(judulParts, "REKAP LAPORAN PERTANDINGAN POPDA 2026")
	if req.Tanggal != "" {
		t, err := time.Parse("2006-01-02", req.Tanggal)
		if err == nil {
			bulan := []string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
				"Juli", "Agustus", "September", "Oktober", "November", "Desember"}
			judulParts = append(judulParts, fmt.Sprintf("%d %s %d", t.Day(), bulan[int(t.Month())], t.Year()))
		}
	}
	judul := strings.Join(judulParts, " — ")

	pdfBytes, err := GenerateBatchPDF(items, req.Penandatangan, judul)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat PDF"})
		return
	}

	tanggal := time.Now().Format("2006-01-02")
	suffix := "semua"
	if req.Tanggal != "" {
		suffix = req.Tanggal
	}
	filename := fmt.Sprintf("laporan_pertandingan_%s_%s.pdf", suffix, tanggal)

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/pdf")
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// parseUintParam parse path param menjadi uint
func parseUintParam(c *gin.Context, name string) (uint, error) {
	v := c.Param(name)
	var id uint64
	_, err := fmt.Sscan(v, &id)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
