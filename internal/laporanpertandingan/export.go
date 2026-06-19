package laporanpertandingan

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// TTDData berisi data tanda tangan satu penandatangan
type TTDData struct {
	Jabatan      string // misal "Wasit", "Ketua Panitia", "Kepala Dispora"
	NamaTercetak string // nama yang dicetak di bawah garis
	NIP          string // opsional
	SignatureB64 string // base64 PNG dari signature pad (opsional)
}

// ExportRequest adalah payload POST untuk export PDF
// Berisi filter + data tanda tangan
type ExportRequest struct {
	// Filter — kosong = semua
	Tanggal string `json:"tanggal"` // YYYY-MM-DD, filter per hari
	CaborID uint   `json:"cabor_id"`
	NomorID uint   `json:"nomor_id"`

	// Tanda tangan — array, biasanya 2-3 orang
	Penandatangan []TTDData `json:"penandatangan"`
}

// ===== HELPER LABEL =====

var labelBabak = map[string]string{
	"PENYISIHAN":         "Penyisihan",
	"8_BESAR":            "8 Besar",
	"PEREMPAT_FINAL":     "Perempat Final",
	"SEMIFINAL":          "Semifinal",
	"FINAL":              "Final",
	"PEREBUTAN_TEMPAT_3": "Perebutan Tempat 3",
	"LAINNYA":            "Lainnya",
}

var labelPemenang = map[string]string{
	"TIM_A": "Tim A",
	"TIM_B": "Tim B",
	"DRAW":  "Seri",
}

func formatTglLaporan(t TanggalDate) string {
	if t.IsZero() {
		return "-"
	}
	bulan := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	return fmt.Sprintf("%d %s %d", t.Day(), bulan[int(t.Month())], t.Year())
}

func babakLabel(b string) string {
	if v, ok := labelBabak[b]; ok {
		return v
	}
	return b
}

func pemenangLabel(p string) string {
	if v, ok := labelPemenang[p]; ok {
		return v
	}
	return p
}

func juaraLabel(j *uint8) string {
	if j == nil {
		return "-"
	}
	return fmt.Sprintf("Juara %d", *j)
}

func namaKontingenB(n *string) string {
	if n == nil {
		return "-"
	}
	return *n
}

func catatanStr(s *string) string {
	if s == nil || *s == "" {
		return "-"
	}
	return *s
}

// ===== pdfWriter =====

type pdfWriter struct{ buf []byte }

func (w *pdfWriter) Write(p []byte) (int, error) {
	w.buf = append(w.buf, p...)
	return len(p), nil
}

// ===== GENERATE SATU PERTANDINGAN =====

// GenerateSatuPDF membuat PDF portrait satu pertandingan
func GenerateSatuPDF(l *LaporanDetail, ttds []TTDData) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.SetAutoPageBreak(true, 20)
	pdf.AddPage()

	const usableW = 170.0

	renderSatuLaporan(pdf, l, usableW)
	renderTTD(pdf, ttds, usableW)
	renderFooter(pdf, usableW)

	w := &pdfWriter{}
	if err := pdf.Output(w); err != nil {
		return nil, err
	}
	return w.buf, nil
}

// ===== GENERATE BATCH (per hari / semua) =====

// GenerateBatchPDF membuat PDF portrait berisi banyak laporan
// Setiap laporan dimulai dari posisi baru (bukan halaman baru — pakai spacer)
// Tanda tangan HANYA di halaman terakhir
func GenerateBatchPDF(items []LaporanDetail, ttds []TTDData, judul string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.SetAutoPageBreak(true, 20)
	pdf.AddPage()

	const usableW = 170.0
	tanggalCetak := time.Now().Format("02 January 2006")

	// ===== HEADER BATCH =====
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(20, 80, 40)
	pdf.CellFormat(usableW, 9, judul, "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(usableW, 5, fmt.Sprintf("Dicetak: %s  |  Total Pertandingan: %d", tanggalCetak, len(items)), "", 1, "C", false, 0, "")

	pdf.SetDrawColor(20, 80, 40)
	pdf.SetLineWidth(0.5)
	pdf.Line(20, pdf.GetY()+2, 20+usableW, pdf.GetY()+2)
	pdf.Ln(6)

	// ===== ISI SETIAP LAPORAN =====
	for i, l := range items {
		item := l // copy
		renderSatuLaporan(pdf, &item, usableW)

		// Garis pemisah antar laporan (kecuali yang terakhir)
		if i < len(items)-1 {
			pdf.SetDrawColor(180, 180, 180)
			pdf.SetLineWidth(0.3)
			pdf.Line(20, pdf.GetY()+3, 20+usableW, pdf.GetY()+3)
			pdf.Ln(8)
		}
	}

	// ===== TANDA TANGAN DI AKHIR =====
	pdf.Ln(6)
	renderTTD(pdf, ttds, usableW)
	renderFooter(pdf, usableW)

	w := &pdfWriter{}
	if err := pdf.Output(w); err != nil {
		return nil, err
	}
	return w.buf, nil
}

// ===== RENDER SATU LAPORAN =====

func renderSatuLaporan(pdf *gofpdf.Fpdf, l *LaporanDetail, usableW float64) {
	const marginL = 20.0
	rowH := 6.5

	// ─── Sub-header: nomor laporan + tanggal ───
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(20, 20, 80)
	pdf.SetX(marginL)
	pdf.CellFormat(usableW, 7,
		fmt.Sprintf("LAPORAN PERTANDINGAN  #%d", l.ID),
		"", 1, "L", false, 0, "")

	// Garis bawah sub-header
	pdf.SetDrawColor(20, 20, 80)
	pdf.SetLineWidth(0.3)
	pdf.Line(marginL, pdf.GetY(), marginL+usableW, pdf.GetY())
	pdf.Ln(2)

	// ─── Baris info ───
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(40, 40, 40)

	labelW := 55.0
	valueW := usableW - labelW

	type row struct{ label, value string }
	rows := []row{
		{"Tanggal Pertandingan", formatTglLaporan(l.TanggalPertandingan)},
		{"Waktu", strings.TrimSuffix(l.WaktuPertandingan, ":00") + " WIB"},		{"Venue / Lapangan", l.Venue},
		{"Cabang Olahraga", l.NamaCabor},
		{"Nomor / Kelas", l.NamaNomor},
		{"Babak", babakLabel(l.Babak)},
		{"Tim A", l.NamaKontingenA},
		{"Tim B", namaKontingenB(l.NamaKontingenB)},
	}

	// Atlet A
	if len(l.AtletA) > 0 {
		names := make([]string, len(l.AtletA))
		for i, a := range l.AtletA {
			names[i] = fmt.Sprintf("%d. %s", a.Urutan, a.NamaLengkap)
		}
		rows = append(rows, row{"Atlet Sisi A", strings.Join(names, " | ")})
	}
	// Atlet B
	if len(l.AtletB) > 0 {
		names := make([]string, len(l.AtletB))
		for i, a := range l.AtletB {
			names[i] = fmt.Sprintf("%d. %s", a.Urutan, a.NamaLengkap)
		}
		rows = append(rows, row{"Atlet Sisi B", strings.Join(names, " | ")})
	}

	rows = append(rows,
		row{"Hasil Pertandingan", l.HasilPertandingan},
		row{"Pemenang", pemenangLabel(l.Pemenang)},
		row{"Juara Ke", juaraLabel(l.JuaraKe)},
		row{"Wasit / Juri", l.Wasit},
		row{"Catatan Khusus", catatanStr(l.CatatanKhusus)},
	)

	fill := false
	for _, r := range rows {
		if fill {
			pdf.SetFillColor(245, 248, 245)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		pdf.SetX(marginL)
		pdf.SetFont("Arial", "B", 9)
		pdf.SetTextColor(60, 60, 60)
		pdf.CellFormat(labelW, rowH, r.label, "LTB", 0, "L", fill, 0, "")
		pdf.SetFont("Arial", "", 9)
		pdf.SetTextColor(20, 20, 20)
		pdf.CellFormat(valueW, rowH, r.value, "RTB", 1, "L", fill, 0, "")
		fill = !fill
	}
	pdf.Ln(2)
}

// ===== RENDER TANDA TANGAN =====

func renderTTD(pdf *gofpdf.Fpdf, ttds []TTDData, usableW float64) {
	if len(ttds) == 0 {
		return
	}
	const marginL = 20.0
	tanggalSerang := "Serang, " + time.Now().Format("02 January 2006")

	pdf.Ln(4)
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(60, 60, 60)
	pdf.SetX(marginL)
	pdf.CellFormat(usableW, 6, tanggalSerang, "", 1, "R", false, 0, "")
	pdf.Ln(2)

	colW := usableW / float64(len(ttds))

	// Jabatan
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(60, 60, 60)
	pdf.SetX(marginL)
	for _, ttd := range ttds {
		pdf.CellFormat(colW, 5, ttd.Jabatan+",", "", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Area tanda tangan (embed gambar atau ruang kosong)
	sigH := 22.0
	for i, ttd := range ttds {
		xPos := marginL + float64(i)*colW
		imgX := xPos + colW/2 - 25

		if ttd.SignatureB64 != "" {
			imgData, err := decodeSignatureB64(ttd.SignatureB64)
			if err == nil {
				imgName := fmt.Sprintf("sig_%d_%d", i, time.Now().UnixNano())
				pdf.RegisterImageReader(imgName, "PNG", bytes.NewReader(imgData))
				pdf.Image(imgName, imgX, pdf.GetY()+1, 50, sigH-4, false, "PNG", 0, "")
			}
		}
	}
	pdf.Ln(sigH)

	// Garis tanda tangan
	pdf.SetDrawColor(80, 80, 80)
	pdf.SetLineWidth(0.3)
	for i := range ttds {
		lineX1 := marginL + float64(i)*colW + colW*0.1
		lineX2 := marginL + float64(i)*colW + colW*0.9
		pdf.Line(lineX1, pdf.GetY(), lineX2, pdf.GetY())
	}
	pdf.Ln(2)

	// Nama tercetak
	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(20, 20, 20)
	pdf.SetX(marginL)
	for _, ttd := range ttds {
		pdf.CellFormat(colW, 5, ttd.NamaTercetak, "", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// NIP (jika ada)
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetX(marginL)
	hasNIP := false
	for _, ttd := range ttds {
		if ttd.NIP != "" {
			hasNIP = true
			break
		}
	}
	if hasNIP {
		for _, ttd := range ttds {
			nip := ""
			if ttd.NIP != "" {
				nip = "NIP. " + ttd.NIP
			}
			pdf.CellFormat(colW, 5, nip, "", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}
	pdf.Ln(2)
}

// ===== FOOTER =====

func renderFooter(pdf *gofpdf.Fpdf, usableW float64) {
	pdf.Ln(4)
	pdf.SetFont("Arial", "I", 7)
	pdf.SetTextColor(180, 180, 180)
	pdf.SetX(20)
	pdf.CellFormat(usableW, 4,
		"Dicetak oleh Sistem POPDA 2026  •  "+time.Now().Format("02 January 2006 15:04"),
		"", 0, "C", false, 0, "")
}

// ===== HELPER DECODE SIGNATURE =====

// decodeSignatureB64 decode base64 PNG dari signature pad
// Input boleh: "data:image/png;base64,iVBOR..." atau murni base64
func decodeSignatureB64(b64 string) ([]byte, error) {
	// Hapus data URI prefix jika ada
	if idx := strings.Index(b64, ","); idx != -1 {
		b64 = b64[idx+1:]
	}
	return base64.StdEncoding.DecodeString(b64)
}
