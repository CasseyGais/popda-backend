package tahap3

import (
	"fmt"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

// formatTanggal mengubah berbagai format tanggal DB menjadi "2006-01-02"
// Menangani: "2003-01-01", "2003-01-01T00:00:00+07:00", "2003-01-01 00:00:00"
func formatTanggal(s string) string {
	if s == "" {
		return ""
	}
	// Potong bagian waktu jika ada
	// Format ISO: "2003-01-01T00:00:00+07:00" atau "2003-01-01T00:00:00Z"
	if idx := strings.IndexAny(s, "T "); idx != -1 {
		s = s[:idx]
	}
	// Coba parse YYYY-MM-DD
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return s // kembalikan aslinya jika tidak bisa parse
	}
	return t.Format("2006-01-02")
}

// ExportData berisi semua data untuk generate file export tahap 3
type ExportData struct {
	Kontingen *Kontingen
	Atlets    []MasterAtlet
	Pelatihs  []MasterPelatih
	Officials []MasterOfficial
}

// ===== PDF GENERATOR =====

// GeneratePDF membuat satu PDF dengan 3 section: Atlet, Pelatih, Official
func GeneratePDF(data *ExportData) ([]byte, error) {
	// A4 Landscape: 297mm lebar, margin 10mm kiri+kanan → usable = 277mm
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 15, 10)
	pdf.SetAutoPageBreak(true, 15)

	tanggal := time.Now().Format("02 January 2006")

	// ===== SECTION 1: ATLET =====
	pdf.AddPage()
	writeLandscapeHeader(pdf, "REKAP ENTRY BY NAME - ATLET", data.Kontingen.NamaKontingen, tanggal, data.Kontingen.Tahap3Status)
	writeAtletTable(pdf, data.Atlets)

	// ===== SECTION 2: PELATIH =====
	pdf.AddPage()
	writeLandscapeHeader(pdf, "REKAP ENTRY BY NAME - PELATIH", data.Kontingen.NamaKontingen, tanggal, data.Kontingen.Tahap3Status)
	writePelatihTable(pdf, data.Pelatihs)

	// ===== SECTION 3: OFFICIAL =====
	pdf.AddPage()
	writeLandscapeHeader(pdf, "REKAP ENTRY BY NAME - OFFICIAL", data.Kontingen.NamaKontingen, tanggal, data.Kontingen.Tahap3Status)
	writeOfficialTable(pdf, data.Officials)

	w := &pdfWriter{}
	if err := pdf.Output(w); err != nil {
		return nil, err
	}
	return w.buf, nil
}

// GeneratePDFAtlet hanya section atlet
func GeneratePDFAtlet(kontingen *Kontingen, atlets []MasterAtlet) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 15, 10)
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()
	tanggal := time.Now().Format("02 January 2006")
	writeLandscapeHeader(pdf, "REKAP ENTRY BY NAME - ATLET", kontingen.NamaKontingen, tanggal, kontingen.Tahap3Status)
	writeAtletTable(pdf, atlets)
	w := &pdfWriter{}
	if err := pdf.Output(w); err != nil {
		return nil, err
	}
	return w.buf, nil
}

// GeneratePDFPelatih hanya section pelatih
func GeneratePDFPelatih(kontingen *Kontingen, pelatihs []MasterPelatih) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 15, 10)
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()
	tanggal := time.Now().Format("02 January 2006")
	writeLandscapeHeader(pdf, "REKAP ENTRY BY NAME - PELATIH", kontingen.NamaKontingen, tanggal, kontingen.Tahap3Status)
	writePelatihTable(pdf, pelatihs)
	w := &pdfWriter{}
	if err := pdf.Output(w); err != nil {
		return nil, err
	}
	return w.buf, nil
}

// GeneratePDFOfficial hanya section official
func GeneratePDFOfficial(kontingen *Kontingen, officials []MasterOfficial) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 15, 10)
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()
	tanggal := time.Now().Format("02 January 2006")
	writeLandscapeHeader(pdf, "REKAP ENTRY BY NAME - OFFICIAL", kontingen.NamaKontingen, tanggal, kontingen.Tahap3Status)
	writeOfficialTable(pdf, officials)
	w := &pdfWriter{}
	if err := pdf.Output(w); err != nil {
		return nil, err
	}
	return w.buf, nil
}

func writeLandscapeHeader(pdf *gofpdf.Fpdf, judul, namaKontingen, tanggal, status string) {
	marginL := 10.0
	pageW, _ := pdf.GetPageSize()
	marginR := 10.0
	usableW := pageW - marginL - marginR // 277mm landscape

	// Garis atas dekoratif
	pdf.SetFillColor(41, 128, 185)
	pdf.Rect(marginL, 10, usableW, 1.5, "F")

	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(41, 128, 185)
	pdf.CellFormat(usableW, 8, judul+" - POPDA 2026", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(usableW, 5, "Kontingen: "+namaKontingen+"   |   Tanggal Cetak: "+tanggal, "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "B", 9)
	if status == "SUBMITTED" {
		pdf.SetTextColor(39, 174, 96)
	} else {
		pdf.SetTextColor(230, 126, 34)
	}
	pdf.CellFormat(usableW, 5, "Status: "+status, "", 1, "C", false, 0, "")

	pdf.SetDrawColor(220, 220, 220)
	pdf.SetLineWidth(0.3)
	pdf.Line(marginL, pdf.GetY()+2, marginL+usableW, pdf.GetY()+2)
	pdf.Ln(4)
}

func writeAtletTable(pdf *gofpdf.Fpdf, atlets []MasterAtlet) {
	// A4 Landscape usable = 277mm (297 - 10 - 10)
	// No(7)+Nama(55)+JK(7)+TglLahir(22)+NISN(30)+Sekolah(60)+Kelas(25)+KabKota(45)+HP(26) = 277
	const marginL = 10.0
	colW := []float64{7, 55, 7, 22, 30, 60, 25, 45, 26}
	headers := []string{"No", "Nama Lengkap", "JK", "Tgl Lahir", "NISN", "Sekolah", "Kelas", "Kab/Kota", "No. HP"}
	rowH := 7.0
	hdrH := 9.0

	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(41, 128, 185)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetDrawColor(255, 255, 255)
	pdf.SetLineWidth(0.3)
	pdf.SetX(marginL)
	for i, h := range headers {
		pdf.CellFormat(colW[i], hdrH, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 8)
	pdf.SetDrawColor(200, 200, 200)
	pdf.SetLineWidth(0.2)
	fill := false
	for idx, a := range atlets {
		if fill {
			pdf.SetFillColor(235, 245, 251)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		pdf.SetTextColor(50, 50, 50)
		pdf.SetX(marginL)
		pdf.CellFormat(colW[0], rowH, fmt.Sprintf("%d", idx+1), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colW[1], rowH, a.NamaLengkap, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colW[2], rowH, a.JenisKelamin, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colW[3], rowH, formatTanggal(a.TanggalLahir), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colW[4], rowH, a.NISN, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colW[5], rowH, a.Sekolah, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colW[6], rowH, a.KelasJurusan, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colW[7], rowH, a.KabupatenKota, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colW[8], rowH, a.NoHP, "1", 0, "L", fill, 0, "")
		pdf.Ln(-1)
		fill = !fill
	}
	pdf.Ln(4)
	pdf.SetX(marginL)
	pdf.SetFont("Arial", "B", 8)
	pdf.SetTextColor(50, 50, 50)
	pdf.CellFormat(0, 6, fmt.Sprintf("Total Atlet: %d", len(atlets)), "", 1, "L", false, 0, "")
}

func writePelatihTable(pdf *gofpdf.Fpdf, pelatihs []MasterPelatih) {
	// No(7)+Nama(60)+JK(7)+TglLahir(22)+Jabatan(50)+HP(30)+Email(60)+KabKota(41) = 277
	const marginL = 10.0
	colW := []float64{7, 60, 7, 22, 50, 30, 60, 41}
	headers := []string{"No", "Nama Lengkap", "JK", "Tgl Lahir", "Jabatan", "No. HP", "Email", "Kab/Kota"}
	rowH := 7.0
	hdrH := 9.0

	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(41, 128, 185)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetDrawColor(255, 255, 255)
	pdf.SetLineWidth(0.3)
	pdf.SetX(marginL)
	for i, h := range headers {
		pdf.CellFormat(colW[i], hdrH, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 8)
	pdf.SetDrawColor(200, 200, 200)
	pdf.SetLineWidth(0.2)
	fill := false
	for idx, p := range pelatihs {
		if fill {
			pdf.SetFillColor(235, 245, 251)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		pdf.SetTextColor(50, 50, 50)
		pdf.SetX(marginL)
		pdf.CellFormat(colW[0], rowH, fmt.Sprintf("%d", idx+1), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colW[1], rowH, p.NamaLengkap, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colW[2], rowH, p.JenisKelamin, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colW[3], rowH, formatTanggal(p.TanggalLahir), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colW[4], rowH, p.Jabatan, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colW[5], rowH, p.NoHP, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colW[6], rowH, p.Email, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colW[7], rowH, p.KabupatenKota, "1", 0, "L", fill, 0, "")
		pdf.Ln(-1)
		fill = !fill
	}
	pdf.Ln(4)
	pdf.SetX(marginL)
	pdf.SetFont("Arial", "B", 8)
	pdf.SetTextColor(50, 50, 50)
	pdf.CellFormat(0, 6, fmt.Sprintf("Total Pelatih: %d", len(pelatihs)), "", 1, "L", false, 0, "")
}

func writeOfficialTable(pdf *gofpdf.Fpdf, officials []MasterOfficial) {
	// No(7)+Nama(60)+JK(7)+TglLahir(22)+Jabatan(50)+HP(30)+Email(60)+KabKota(41) = 277
	const marginL = 10.0
	colW := []float64{7, 60, 7, 22, 50, 30, 60, 41}
	headers := []string{"No", "Nama Lengkap", "JK", "Tgl Lahir", "Jabatan", "No. HP", "Email", "Kab/Kota"}
	rowH := 7.0
	hdrH := 9.0

	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(41, 128, 185)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetDrawColor(255, 255, 255)
	pdf.SetLineWidth(0.3)
	pdf.SetX(marginL)
	for i, h := range headers {
		pdf.CellFormat(colW[i], hdrH, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 8)
	pdf.SetDrawColor(200, 200, 200)
	pdf.SetLineWidth(0.2)
	fill := false
	for idx, o := range officials {
		if fill {
			pdf.SetFillColor(235, 245, 251)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		pdf.SetTextColor(50, 50, 50)
		pdf.SetX(marginL)
		pdf.CellFormat(colW[0], rowH, fmt.Sprintf("%d", idx+1), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colW[1], rowH, o.NamaLengkap, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colW[2], rowH, o.JenisKelamin, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colW[3], rowH, formatTanggal(o.TanggalLahir), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colW[4], rowH, o.Jabatan, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colW[5], rowH, o.NoHP, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colW[6], rowH, o.Email, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colW[7], rowH, o.KabupatenKota, "1", 0, "L", fill, 0, "")
		pdf.Ln(-1)
		fill = !fill
	}
	pdf.Ln(4)
	pdf.SetX(marginL)
	pdf.SetFont("Arial", "B", 8)
	pdf.SetTextColor(50, 50, 50)
	pdf.CellFormat(0, 6, fmt.Sprintf("Total Official: %d", len(officials)), "", 1, "L", false, 0, "")
}

// ===== EXCEL GENERATOR =====

// GenerateExcel membuat satu file Excel dengan 3 sheet: Atlet, Pelatih, Official
func GenerateExcel(data *ExportData) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	f.DeleteSheet("Sheet1")

	tanggal := time.Now().Format("02 January 2006")
	namaKontingen := data.Kontingen.NamaKontingen
	status := data.Kontingen.Tahap3Status

	if err := writeAtletSheet(f, namaKontingen, tanggal, status, data.Atlets); err != nil {
		return nil, err
	}
	if err := writePelatihSheet(f, namaKontingen, tanggal, status, data.Pelatihs); err != nil {
		return nil, err
	}
	if err := writeOfficialSheet(f, namaKontingen, tanggal, status, data.Officials); err != nil {
		return nil, err
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerateExcelAtlet membuat Excel hanya sheet atlet
func GenerateExcelAtlet(kontingen *Kontingen, atlets []MasterAtlet) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()
	f.DeleteSheet("Sheet1")
	tanggal := time.Now().Format("02 January 2006")
	if err := writeAtletSheet(f, kontingen.NamaKontingen, tanggal, kontingen.Tahap3Status, atlets); err != nil {
		return nil, err
	}
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerateExcelPelatih membuat Excel hanya sheet pelatih
func GenerateExcelPelatih(kontingen *Kontingen, pelatihs []MasterPelatih) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()
	f.DeleteSheet("Sheet1")
	tanggal := time.Now().Format("02 January 2006")
	if err := writePelatihSheet(f, kontingen.NamaKontingen, tanggal, kontingen.Tahap3Status, pelatihs); err != nil {
		return nil, err
	}
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerateExcelOfficial membuat Excel hanya sheet official
func GenerateExcelOfficial(kontingen *Kontingen, officials []MasterOfficial) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()
	f.DeleteSheet("Sheet1")
	tanggal := time.Now().Format("02 January 2006")
	if err := writeOfficialSheet(f, kontingen.NamaKontingen, tanggal, kontingen.Tahap3Status, officials); err != nil {
		return nil, err
	}
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ===== INTERNAL SHEET WRITERS =====

func makeHeaderStyle(f *excelize.File) int {
	s, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"2980B9"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "FFFFFF", Style: 1},
			{Type: "right", Color: "FFFFFF", Style: 1},
			{Type: "top", Color: "FFFFFF", Style: 1},
			{Type: "bottom", Color: "FFFFFF", Style: 1},
		},
	})
	return s
}

func makeCellStyle(f *excelize.File) int {
	s, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "CCCCCC", Style: 1},
			{Type: "right", Color: "CCCCCC", Style: 1},
			{Type: "top", Color: "CCCCCC", Style: 1},
			{Type: "bottom", Color: "CCCCCC", Style: 1},
		},
	})
	return s
}

func makeTitleStyle(f *excelize.File) int {
	s, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 13},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	return s
}

func makeSubtitleStyle(f *excelize.File) int {
	s, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	return s
}

func writeSheetHeader(f *excelize.File, sheet, judul, namaKontingen, tanggal, status string, colCount int) {
	lastCol, _ := excelize.ColumnNumberToName(colCount)
	titleStyle := makeTitleStyle(f)
	subStyle := makeSubtitleStyle(f)

	f.MergeCell(sheet, "A1", lastCol+"1")
	f.SetCellValue(sheet, "A1", judul+" - POPDA 2026")
	f.SetCellStyle(sheet, "A1", lastCol+"1", titleStyle)
	f.SetRowHeight(sheet, 1, 22)

	f.MergeCell(sheet, "A2", lastCol+"2")
	f.SetCellValue(sheet, "A2", "Kontingen: "+namaKontingen)
	f.SetCellStyle(sheet, "A2", lastCol+"2", subStyle)

	f.MergeCell(sheet, "A3", lastCol+"3")
	f.SetCellValue(sheet, "A3", "Tanggal Cetak: "+tanggal+" | Status: "+status)
	f.SetCellStyle(sheet, "A3", lastCol+"3", subStyle)
}

func writeAtletSheet(f *excelize.File, namaKontingen, tanggal, status string, atlets []MasterAtlet) error {
	sheet := "Atlet"
	f.NewSheet(sheet)

	// Kolom: No, Nama Lengkap, JK, Tanggal Lahir, NISN, Sekolah, Kelas, Kab/Kota, No. HP
	headers := []string{"No", "Nama Lengkap", "JK", "Tanggal Lahir", "NISN", "Sekolah", "Kelas/Jurusan", "Kab/Kota", "No. HP"}
	widths := []float64{6, 30, 5, 15, 15, 30, 15, 22, 15}
	colNames := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}

	writeSheetHeader(f, sheet, "REKAP ENTRY BY NAME - ATLET", namaKontingen, tanggal, status, len(headers))
	headerStyle := makeHeaderStyle(f)
	cellStyle := makeCellStyle(f)

	for i, h := range headers {
		cell := colNames[i] + "5"
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
		f.SetColWidth(sheet, colNames[i], colNames[i], widths[i])
	}
	f.SetRowHeight(sheet, 5, 18)

	for idx, a := range atlets {
		r := fmt.Sprintf("%d", idx+6)
		f.SetCellValue(sheet, "A"+r, idx+1)
		f.SetCellValue(sheet, "B"+r, a.NamaLengkap)
		f.SetCellValue(sheet, "C"+r, a.JenisKelamin)
		f.SetCellValue(sheet, "D"+r, formatTanggal(a.TanggalLahir))
		f.SetCellValue(sheet, "E"+r, a.NISN)
		f.SetCellValue(sheet, "F"+r, a.Sekolah)
		f.SetCellValue(sheet, "G"+r, a.KelasJurusan)
		f.SetCellValue(sheet, "H"+r, a.KabupatenKota)
		f.SetCellValue(sheet, "I"+r, a.NoHP)
		for _, col := range colNames {
			f.SetCellStyle(sheet, col+r, col+r, cellStyle)
		}
	}
	return nil
}

func writePelatihSheet(f *excelize.File, namaKontingen, tanggal, status string, pelatihs []MasterPelatih) error {
	sheet := "Pelatih"
	f.NewSheet(sheet)

	// Kolom: No, Nama Lengkap, JK, Tanggal Lahir, NIK, Jabatan, No. HP, Email, Kab/Kota
	headers := []string{"No", "Nama Lengkap", "JK", "Tanggal Lahir", "NIK", "Jabatan", "No. HP", "Email", "Kab/Kota"}
	widths := []float64{6, 30, 5, 15, 20, 22, 16, 30, 22}
	colNames := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}

	writeSheetHeader(f, sheet, "REKAP ENTRY BY NAME - PELATIH", namaKontingen, tanggal, status, len(headers))
	headerStyle := makeHeaderStyle(f)
	cellStyle := makeCellStyle(f)

	for i, h := range headers {
		cell := colNames[i] + "5"
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
		f.SetColWidth(sheet, colNames[i], colNames[i], widths[i])
	}
	f.SetRowHeight(sheet, 5, 18)

	for idx, p := range pelatihs {
		r := fmt.Sprintf("%d", idx+6)
		f.SetCellValue(sheet, "A"+r, idx+1)
		f.SetCellValue(sheet, "B"+r, p.NamaLengkap)
		f.SetCellValue(sheet, "C"+r, p.JenisKelamin)
		f.SetCellValue(sheet, "D"+r, formatTanggal(p.TanggalLahir))
		f.SetCellValue(sheet, "E"+r, p.NIK)
		f.SetCellValue(sheet, "F"+r, p.Jabatan)
		f.SetCellValue(sheet, "G"+r, p.NoHP)
		f.SetCellValue(sheet, "H"+r, p.Email)
		f.SetCellValue(sheet, "I"+r, p.KabupatenKota)
		for _, col := range colNames {
			f.SetCellStyle(sheet, col+r, col+r, cellStyle)
		}
	}
	return nil
}

func writeOfficialSheet(f *excelize.File, namaKontingen, tanggal, status string, officials []MasterOfficial) error {
	sheet := "Official"
	f.NewSheet(sheet)

	// Kolom: No, Nama Lengkap, JK, Tanggal Lahir, NIK, Jabatan, No. HP, Email, Kab/Kota
	headers := []string{"No", "Nama Lengkap", "JK", "Tanggal Lahir", "NIK", "Jabatan", "No. HP", "Email", "Kab/Kota"}
	widths := []float64{6, 30, 5, 15, 20, 22, 16, 30, 22}
	colNames := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}

	writeSheetHeader(f, sheet, "REKAP ENTRY BY NAME - OFFICIAL", namaKontingen, tanggal, status, len(headers))
	headerStyle := makeHeaderStyle(f)
	cellStyle := makeCellStyle(f)

	for i, h := range headers {
		cell := colNames[i] + "5"
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
		f.SetColWidth(sheet, colNames[i], colNames[i], widths[i])
	}
	f.SetRowHeight(sheet, 5, 18)

	for idx, o := range officials {
		r := fmt.Sprintf("%d", idx+6)
		f.SetCellValue(sheet, "A"+r, idx+1)
		f.SetCellValue(sheet, "B"+r, o.NamaLengkap)
		f.SetCellValue(sheet, "C"+r, o.JenisKelamin)
		f.SetCellValue(sheet, "D"+r, formatTanggal(o.TanggalLahir))
		f.SetCellValue(sheet, "E"+r, o.NIK)
		f.SetCellValue(sheet, "F"+r, o.Jabatan)
		f.SetCellValue(sheet, "G"+r, o.NoHP)
		f.SetCellValue(sheet, "H"+r, o.Email)
		f.SetCellValue(sheet, "I"+r, o.KabupatenKota)
		for _, col := range colNames {
			f.SetCellStyle(sheet, col+r, col+r, cellStyle)
		}
	}
	return nil
}

// pdfWriter adalah io.Writer sederhana untuk output gofpdf
type pdfWriter struct {
	buf []byte
}

func (w *pdfWriter) Write(p []byte) (int, error) {
	w.buf = append(w.buf, p...)
	return len(p), nil
}
