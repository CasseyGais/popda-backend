package sertifikat

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// TTDSertifikat berisi data satu penandatangan untuk sertifikat PDF
type TTDSertifikat struct {
	Jabatan      string `json:"jabatan"`
	NamaTercetak string `json:"nama_tercetak"`
	NIP          string `json:"nip"`
	SignatureB64 string `json:"signature_b64"` // base64 PNG dari signature pad
}

// SertifikatPDFData berisi semua data yang diperlukan untuk generate satu lembar sertifikat.
type SertifikatPDFData struct {
	NomorSertifikat string // misal "POPDA/2026/ATL/001", bisa kosong
	NamaPenerima    string // dari nama_lengkap tabel master
	Judul           string // misal "Sertifikat Peserta POPDA 2026"
	TipePenerima    string // ATLET / PELATIH / OFFICIAL
	TanggalTerbit   string // YYYY-MM-DD
	Catatan         string // opsional
	// Penandatangan — opsional, jika kosong tampilkan garis kosong
	Penandatangan []TTDSertifikat
}

// GenerateSertifikatPDF membuat PDF landscape satu halaman per sertifikat.
// Jika slice berisi lebih dari satu item, setiap item jadi satu halaman.
// Ukuran: A4 Landscape (297mm x 210mm)
func GenerateSertifikatPDF(items []SertifikatPDFData) ([]byte, error) {
	// A4 Landscape: 297 x 210 mm
	// Margin kiri/kanan 20, atas/bawah 18 → usable 257mm x 174mm
	const (
		pageW   = 297.0
		pageH   = 210.0
		marginL = 20.0
		marginT = 18.0
		marginR = 20.0
		marginB = 18.0
		usableW = pageW - marginL - marginR // 257mm
		usableH = pageH - marginT - marginB // 174mm
		centerX = pageW / 2                 // 148.5mm
	)

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(marginL, marginT, marginR)
	pdf.SetAutoPageBreak(false, 0)

	for _, item := range items {
		pdf.AddPage()

		// ===== BORDER DEKORATIF =====
		// Garis bingkai luar tebal
		pdf.SetDrawColor(20, 80, 40) // hijau tua
		pdf.SetLineWidth(2.5)
		pdf.Rect(marginL-5, marginT-5, usableW+10, usableH+10, "D")

		// Garis bingkai dalam tipis
		pdf.SetDrawColor(180, 150, 20) // emas
		pdf.SetLineWidth(0.8)
		pdf.Rect(marginL-3, marginT-3, usableW+6, usableH+6, "D")

		// ===== BACKGROUND WATERMARK TEXT (tipis di tengah) =====
		pdf.SetFont("Arial", "B", 60)
		pdf.SetTextColor(230, 240, 230) // abu-abu sangat terang
		pdf.SetXY(marginL, pageH/2-18)
		pdf.CellFormat(usableW, 36, "POPDA 2026", "", 0, "C", false, 0, "")

		// ===== ORNAMEN SUDUT (kotak kecil) =====
		cornerSize := 12.0
		pdf.SetFillColor(20, 80, 40)
		// Kiri atas
		pdf.Rect(marginL-5, marginT-5, cornerSize, cornerSize, "F")
		// Kanan atas
		pdf.Rect(pageW-marginR-cornerSize+5, marginT-5, cornerSize, cornerSize, "F")
		// Kiri bawah
		pdf.Rect(marginL-5, pageH-marginB-cornerSize+5, cornerSize, cornerSize, "F")
		// Kanan bawah
		pdf.Rect(pageW-marginR-cornerSize+5, pageH-marginB-cornerSize+5, cornerSize, cornerSize, "F")

		// ===== JUDUL UTAMA "PIAGAM PENGHARGAAN" =====
		y := marginT + 8.0
		pdf.SetFont("Arial", "B", 26)
		pdf.SetTextColor(20, 80, 40) // hijau tua
		pdf.SetXY(marginL, y)
		pdf.CellFormat(usableW, 12, "PIAGAM PENGHARGAAN", "", 1, "C", false, 0, "")
		y += 12

		// Nomor sertifikat (jika ada)
		if item.NomorSertifikat != "" {
			pdf.SetFont("Arial", "", 10)
			pdf.SetTextColor(80, 80, 80)
			pdf.SetXY(marginL, y)
			pdf.CellFormat(usableW, 7, "No. "+item.NomorSertifikat, "", 1, "C", false, 0, "")
			y += 7
		}
		y += 2

		// Garis pemisah emas
		pdf.SetDrawColor(180, 150, 20)
		pdf.SetLineWidth(0.6)
		pdf.Line(centerX-60, y, centerX+60, y)
		y += 6

		// ===== DIBERIKAN KEPADA =====
		pdf.SetFont("Arial", "I", 12)
		pdf.SetTextColor(80, 80, 80)
		pdf.SetXY(marginL, y)
		pdf.CellFormat(usableW, 7, "Diberikan kepada", "", 1, "C", false, 0, "")
		y += 10

		// ===== NAMA PENERIMA (besar & bold) =====
		pdf.SetFont("Arial", "B", 22)
		pdf.SetTextColor(20, 20, 80) // biru tua
		pdf.SetXY(marginL, y)
		pdf.CellFormat(usableW, 11, item.NamaPenerima, "", 1, "C", false, 0, "")
		y += 11

		// Garis bawah nama
		pdf.SetDrawColor(20, 20, 80)
		pdf.SetLineWidth(0.4)
		nameWidth := float64(len(item.NamaPenerima)) * 3.2
		if nameWidth > 180 {
			nameWidth = 180
		}
		if nameWidth < 60 {
			nameWidth = 60
		}
		pdf.Line(centerX-nameWidth/2, y, centerX+nameWidth/2, y)
		y += 6

		// ===== SEBAGAI =====
		pdf.SetFont("Arial", "I", 12)
		pdf.SetTextColor(80, 80, 80)
		pdf.SetXY(marginL, y)
		pdf.CellFormat(usableW, 7, "Sebagai", "", 1, "C", false, 0, "")
		y += 9

		// ===== JUDUL / PREDIKAT =====
		pdf.SetFont("Arial", "B", 14)
		pdf.SetTextColor(180, 100, 0) // coklat/emas
		pdf.SetXY(marginL, y)
		pdf.CellFormat(usableW, 8, item.Judul, "", 1, "C", false, 0, "")
		y += 8

		// Tipe penerima (atlet/pelatih/official)
		tipeTeks := map[string]string{
			"ATLET":    "Atlet",
			"PELATIH":  "Pelatih",
			"OFFICIAL": "Official",
		}
		tipeLabel, ok := tipeTeks[item.TipePenerima]
		if !ok {
			tipeLabel = item.TipePenerima
		}
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(100, 100, 100)
		pdf.SetXY(marginL, y)
		pdf.CellFormat(usableW, 6, tipeLabel, "", 1, "C", false, 0, "")
		y += 8

		// ===== KALIMAT DALAM RANGKA =====
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(60, 60, 60)
		pdf.SetXY(marginL, y)
		namaKegiatan := "Dalam Rangka Pekan Olahraga Pelajar Daerah (POPDA) Provinsi Banten Tahun 2026"
		pdf.CellFormat(usableW, 6, namaKegiatan, "", 1, "C", false, 0, "")
		y += 7

		// Catatan tambahan (jika ada)
		if item.Catatan != "" {
			pdf.SetFont("Arial", "I", 9)
			pdf.SetTextColor(100, 100, 100)
			pdf.SetXY(marginL, y)
			pdf.CellFormat(usableW, 5, item.Catatan, "", 1, "C", false, 0, "")
			y += 6
		}

		// ===== BLOK TANDA TANGAN =====
		tanggalFormatted := formatTanggal(item.TanggalTerbit)
		y += 2
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(60, 60, 60)
		pdf.SetXY(marginL, y)
		pdf.CellFormat(usableW, 6, "Serang, "+tanggalFormatted, "", 1, "C", false, 0, "")
		y += 3

		renderTTDSertifikat(pdf, item.Penandatangan, marginL, y, usableW)

		// ===== FOOTER kecil =====
		pdf.SetFont("Arial", "I", 7)
		pdf.SetTextColor(180, 180, 180)
		pdf.SetXY(marginL, pageH-marginB-2)
		pdf.CellFormat(usableW, 4,
			fmt.Sprintf("Dicetak oleh Sistem POPDA 2026  •  %s", time.Now().Format("02 January 2006")),
			"", 0, "C", false, 0, "")
	}

	w := &pdfWriter{}
	if err := pdf.Output(w); err != nil {
		return nil, err
	}
	return w.buf, nil
}

// renderTTDSertifikat render blok tanda tangan di halaman sertifikat landscape.
// Jika penandatangan kosong, tampilkan 2 kolom placeholder.
func renderTTDSertifikat(pdf *gofpdf.Fpdf, ttds []TTDSertifikat, marginL, startY, usableW float64) {
	// Default: 2 kolom placeholder jika tidak ada data penandatangan
	if len(ttds) == 0 {
		ttds = []TTDSertifikat{
			{Jabatan: "Ketua Pelaksana,", NamaTercetak: "..............................", NIP: ""},
			{Jabatan: "Kepala Dinas Pemuda dan Olahraga,", NamaTercetak: "..............................", NIP: ""},
		}
	}

	n := float64(len(ttds))
	colW := usableW / n
	sigH := 20.0
	y := startY

	// Jabatan per kolom
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(60, 60, 60)
	pdf.SetX(marginL)
	for _, ttd := range ttds {
		pdf.CellFormat(colW, 5, ttd.Jabatan+",", "", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	y += 5

	// Area tanda tangan — embed gambar jika ada, kosong jika tidak
	for i, ttd := range ttds {
		xPos := marginL + float64(i)*colW
		imgX := xPos + colW/2 - 22
		imgY := y + 1

		if ttd.SignatureB64 != "" {
			imgData, err := decodeSigB64(ttd.SignatureB64)
			if err == nil && len(imgData) > 0 {
				imgName := fmt.Sprintf("sert_sig_%d_%d", i, time.Now().UnixNano())
				pdf.RegisterImageReader(imgName, "PNG", bytes.NewReader(imgData))
				pdf.Image(imgName, imgX, imgY, 44, sigH-2, false, "PNG", 0, "")
			}
		}
	}
	pdf.Ln(sigH)

	// Garis tanda tangan
	pdf.SetDrawColor(80, 80, 80)
	pdf.SetLineWidth(0.3)
	for i := range ttds {
		x1 := marginL + float64(i)*colW + colW*0.08
		x2 := marginL + float64(i)*colW + colW*0.92
		pdf.Line(x1, pdf.GetY(), x2, pdf.GetY())
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

	// NIP (jika ada salah satu yang punya NIP)
	hasNIP := false
	for _, ttd := range ttds {
		if ttd.NIP != "" {
			hasNIP = true
			break
		}
	}
	if hasNIP {
		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(80, 80, 80)
		pdf.SetX(marginL)
		for _, ttd := range ttds {
			nip := ""
			if ttd.NIP != "" {
				nip = "NIP. " + ttd.NIP
			}
			pdf.CellFormat(colW, 5, nip, "", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}
}

// decodeSigB64 decode base64 PNG dari signature pad.
// Toleran terhadap "data:image/png;base64,..." maupun plain base64.
func decodeSigB64(b64 string) ([]byte, error) {
	if idx := strings.Index(b64, ","); idx != -1 {
		b64 = b64[idx+1:]
	}
	return base64.StdEncoding.DecodeString(b64)
}
func formatTanggal(s string) string {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return s
	}
	bulan := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	return fmt.Sprintf("%d %s %d", t.Day(), bulan[int(t.Month())], t.Year())
}

// pdfWriter adalah io.Writer sederhana untuk menampung output gofpdf.
// Didefinisikan lokal di package sertifikat — tidak bisa reuse dari package lain.
type pdfWriter struct {
	buf []byte
}

func (w *pdfWriter) Write(p []byte) (int, error) {
	w.buf = append(w.buf, p...)
	return len(p), nil
}
