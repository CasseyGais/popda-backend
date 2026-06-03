package tahap2

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

// ExportData berisi data lengkap untuk generate file export tahap 2
type ExportData struct {
	Kontingen  *Kontingen
	NomorList  []NomorExportRow
}

// NomorExportRow adalah hasil JOIN untuk export nomor tahap 2
type NomorExportRow struct {
	NamaCabor    string
	NamaNomor    string
	JenisKelamin string
	Tipe         string
}

// GeneratePDF membuat file PDF rekap tahap 2
func GeneratePDF(data *ExportData) ([]byte, error) {
	// A4 Portrait: usable = 210 - 15 - 15 = 180mm
	marginL, marginT, marginR := 15.0, 20.0, 15.0
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginL, marginT, marginR)
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()

	pageW, _ := pdf.GetPageSize()
	usableW := pageW - marginL - marginR // 180mm

	tanggal := time.Now().Format("02 January 2006")

	// ===== HEADER =====
	pdf.SetFillColor(41, 128, 185)
	pdf.Rect(marginL, marginT-5, usableW, 1.5, "F")

	pdf.SetFont("Arial", "B", 15)
	pdf.SetTextColor(41, 128, 185)
	pdf.CellFormat(usableW, 9, "REKAP ENTRY BY NUMBER - POPDA 2026", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(usableW, 6, "Kontingen: "+data.Kontingen.NamaKontingen, "", 1, "C", false, 0, "")
	pdf.CellFormat(usableW, 6, "Tanggal Cetak: "+tanggal, "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "B", 9)
	if data.Kontingen.Tahap2Status == "SUBMITTED" {
		pdf.SetTextColor(39, 174, 96)
	} else {
		pdf.SetTextColor(230, 126, 34)
	}
	pdf.CellFormat(usableW, 6, "Status: "+data.Kontingen.Tahap2Status, "", 1, "C", false, 0, "")

	pdf.SetDrawColor(220, 220, 220)
	pdf.SetLineWidth(0.3)
	pdf.Line(marginL, pdf.GetY()+2, marginL+usableW, pdf.GetY()+2)
	pdf.Ln(5)

	// ===== TABLE =====
	// No(10) + Cabor(60) + Nomor(65) + JK(25) + Tipe(20) = 180mm
	colWidths := []float64{10, 60, 65, 25, 20}
	headers := []string{"No", "Cabang Olahraga", "Nomor Pertandingan", "Jenis\nKelamin", "Tipe"}
	rowH := 8.0
	hdrH := 10.0

	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(41, 128, 185)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetDrawColor(255, 255, 255)
	pdf.SetLineWidth(0.3)
	pdf.SetX(marginL)
	for i, h := range headers {
		pdf.CellFormat(colWidths[i], hdrH, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 9)
	pdf.SetDrawColor(200, 200, 200)
	pdf.SetLineWidth(0.2)
	fill := false
	for idx, row := range data.NomorList {
		if fill {
			pdf.SetFillColor(235, 245, 251)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		pdf.SetTextColor(50, 50, 50)
		jk := row.JenisKelamin
		if jk == "L" {
			jk = "PUTRA"
		} else if jk == "P" {
			jk = "PUTRI"
		}
		pdf.SetX(marginL)
		pdf.CellFormat(colWidths[0], rowH, fmt.Sprintf("%d", idx+1), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colWidths[1], rowH, row.NamaCabor, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colWidths[2], rowH, row.NamaNomor, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colWidths[3], rowH, jk, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colWidths[4], rowH, row.Tipe, "1", 0, "C", fill, 0, "")
		pdf.Ln(-1)
		fill = !fill
	}

	// ===== RINGKASAN =====
	pdf.Ln(5)
	pdf.SetX(marginL)
	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(50, 50, 50)
	pdf.CellFormat(usableW, 7, fmt.Sprintf("Total Nomor Terdaftar: %d", len(data.NomorList)), "", 1, "L", false, 0, "")

	// Footer
	pdf.Ln(4)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(150, 150, 150)
	pdf.CellFormat(usableW, 5, "Dicetak oleh sistem POPDA 2026  •  "+tanggal, "", 1, "C", false, 0, "")

	w := &pdfWriter{}
	if err := pdf.Output(w); err != nil {
		return nil, err
	}
	return w.buf, nil
}

// GenerateExcel membuat file Excel rekap tahap 2
func GenerateExcel(data *ExportData) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Tahap 2"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	tanggal := time.Now().Format("02 January 2006")

	// Styles
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 13},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	subtitleStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	headerStyle, _ := f.NewStyle(&excelize.Style{
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
	cellStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "CCCCCC", Style: 1},
			{Type: "right", Color: "CCCCCC", Style: 1},
			{Type: "top", Color: "CCCCCC", Style: 1},
			{Type: "bottom", Color: "CCCCCC", Style: 1},
		},
	})

	// Lebar kolom
	f.SetColWidth(sheet, "A", "A", 6)
	f.SetColWidth(sheet, "B", "B", 30)
	f.SetColWidth(sheet, "C", "C", 30)
	f.SetColWidth(sheet, "D", "D", 15)
	f.SetColWidth(sheet, "E", "E", 15)

	// Header dokumen
	f.MergeCell(sheet, "A1", "E1")
	f.SetCellValue(sheet, "A1", "REKAP ENTRY BY NUMBER - POPDA 2026")
	f.SetCellStyle(sheet, "A1", "E1", titleStyle)
	f.SetRowHeight(sheet, 1, 22)

	f.MergeCell(sheet, "A2", "E2")
	f.SetCellValue(sheet, "A2", "Kontingen: "+data.Kontingen.NamaKontingen)
	f.SetCellStyle(sheet, "A2", "E2", subtitleStyle)

	f.MergeCell(sheet, "A3", "E3")
	f.SetCellValue(sheet, "A3", "Tanggal Cetak: "+tanggal+" | Status: "+data.Kontingen.Tahap2Status)
	f.SetCellStyle(sheet, "A3", "E3", subtitleStyle)

	// Header tabel (baris 5)
	headers := []string{"No", "Cabang Olahraga", "Nomor Pertandingan", "Jenis Kelamin", "Tipe"}
	cols := []string{"A", "B", "C", "D", "E"}
	for i, h := range headers {
		cell := cols[i] + "5"
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	f.SetRowHeight(sheet, 5, 18)

	// Data rows
	for idx, row := range data.NomorList {
		r := idx + 6
		rowStr := fmt.Sprintf("%d", r)
		jk := row.JenisKelamin
		if jk == "L" {
			jk = "PUTRA"
		} else if jk == "P" {
			jk = "PUTRI"
		}
		f.SetCellValue(sheet, "A"+rowStr, idx+1)
		f.SetCellValue(sheet, "B"+rowStr, row.NamaCabor)
		f.SetCellValue(sheet, "C"+rowStr, row.NamaNomor)
		f.SetCellValue(sheet, "D"+rowStr, jk)
		f.SetCellValue(sheet, "E"+rowStr, row.Tipe)
		for _, col := range cols {
			f.SetCellStyle(sheet, col+rowStr, col+rowStr, cellStyle)
		}
	}

	// Ringkasan
	summaryRow := fmt.Sprintf("%d", len(data.NomorList)+7)
	f.MergeCell(sheet, "A"+summaryRow, "E"+summaryRow)
	f.SetCellValue(sheet, "A"+summaryRow, fmt.Sprintf("Total Nomor Terdaftar: %d", len(data.NomorList)))

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// pdfWriter adalah io.Writer sederhana untuk output gofpdf
type pdfWriter struct {
	buf []byte
}

func (w *pdfWriter) Write(p []byte) (int, error) {
	w.buf = append(w.buf, p...)
	return len(p), nil
}
