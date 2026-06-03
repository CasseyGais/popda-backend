package tahap1

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)
// ExportData berisi data lengkap untuk generate file export tahap 1
type ExportData struct {
	Kontingen  *Kontingen
	CaborList  []CaborExportRow
	TotalPutra int
	TotalPutri int
	TotalPelatih int
	TotalAtlet int
	TotalPersonel int
}

// CaborExportRow adalah join result untuk export — cabor_id + nama cabor
type CaborExportRow struct {
	NamaCabor     string
	Putra         int
	Putri         int
	Pelatih       int
	TotalAtlet    int
	TotalPersonel int
}

// GeneratePDF membuat file PDF rekap tahap 1 dan menulis ke writer
func GeneratePDF(data *ExportData) ([]byte, error) {
	// A4 Portrait: 210mm total
	// Margin kiri 20, kanan 20 → usable = 170mm
	// Kolom: No(8) + Cabor(62) + Putra(20) + Putri(20) + Pelatih(20) + TotAtlet(20) + TotPersonel(20) = 170
	const (
		marginL = 20.0
		marginT = 20.0
		marginR = 20.0
	)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginL, marginT, marginR)
	pdf.SetAutoPageBreak(true, 20)
	pdf.AddPage()

	const usableW = 170.0 // 210 - 20 - 20
	colWidths := []float64{8, 62, 20, 20, 20, 20, 20}

	tanggal := time.Now().Format("02 January 2006")

	// ===== HEADER =====
	pdf.SetFillColor(41, 128, 185)
	pdf.Rect(marginL, marginT-6, usableW, 2, "F")

	pdf.SetX(marginL)
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(41, 128, 185)
	pdf.CellFormat(usableW, 9, "REKAP ENTRY BY SPORT - POPDA 2026", "", 1, "C", false, 0, "")

	pdf.SetX(marginL)
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(usableW, 6, "Kontingen: "+data.Kontingen.NamaKontingen, "", 1, "C", false, 0, "")

	pdf.SetX(marginL)
	pdf.CellFormat(usableW, 6, "Tanggal Cetak: "+tanggal, "", 1, "C", false, 0, "")

	status := data.Kontingen.Tahap1Status
	pdf.SetX(marginL)
	pdf.SetFont("Arial", "B", 9)
	if status == "SUBMITTED" {
		pdf.SetTextColor(39, 174, 96)
	} else {
		pdf.SetTextColor(230, 126, 34)
	}
	pdf.CellFormat(usableW, 6, "Status: "+status, "", 1, "C", false, 0, "")

	pdf.SetDrawColor(220, 220, 220)
	pdf.SetLineWidth(0.3)
	pdf.Line(marginL, pdf.GetY()+2, marginL+usableW, pdf.GetY()+2)
	pdf.Ln(5)

	// ===== TABLE HEADER =====
	headers := []string{"No", "Cabang Olahraga", "Atlet Putra", "Atlet Putri", "Pelatih", "Total Atlet", "Total Personel"}
	rowH := 8.0
	hdrH := 9.0

	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(41, 128, 185)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetDrawColor(255, 255, 255)
	pdf.SetLineWidth(0.3)
	pdf.SetX(marginL)
	for i, h := range headers {
		pdf.CellFormat(colWidths[i], hdrH, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// ===== DATA ROWS =====
	pdf.SetFont("Arial", "", 9)
	pdf.SetDrawColor(200, 200, 200)
	pdf.SetLineWidth(0.2)
	fill := false
	for idx, row := range data.CaborList {
		if fill {
			pdf.SetFillColor(235, 245, 251)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		pdf.SetTextColor(50, 50, 50)
		pdf.SetX(marginL)
		pdf.CellFormat(colWidths[0], rowH, fmt.Sprintf("%d", idx+1), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colWidths[1], rowH, row.NamaCabor, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colWidths[2], rowH, fmt.Sprintf("%d", row.Putra), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colWidths[3], rowH, fmt.Sprintf("%d", row.Putri), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colWidths[4], rowH, fmt.Sprintf("%d", row.Pelatih), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colWidths[5], rowH, fmt.Sprintf("%d", row.TotalAtlet), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colWidths[6], rowH, fmt.Sprintf("%d", row.TotalPersonel), "1", 0, "C", fill, 0, "")
		pdf.Ln(-1)
		fill = !fill
	}

	// ===== BARIS TOTAL =====
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(41, 128, 185)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetDrawColor(255, 255, 255)
	pdf.SetX(marginL)
	pdf.CellFormat(colWidths[0]+colWidths[1], rowH, "TOTAL", "1", 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[2], rowH, fmt.Sprintf("%d", data.TotalPutra), "1", 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[3], rowH, fmt.Sprintf("%d", data.TotalPutri), "1", 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[4], rowH, fmt.Sprintf("%d", data.TotalPelatih), "1", 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[5], rowH, fmt.Sprintf("%d", data.TotalAtlet), "1", 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[6], rowH, fmt.Sprintf("%d", data.TotalPersonel), "1", 0, "C", true, 0, "")
	pdf.Ln(-1)

	// ===== FOOTER =====
	pdf.Ln(8)
	pdf.SetX(marginL)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(150, 150, 150)
	pdf.CellFormat(usableW, 5, "Dicetak oleh sistem POPDA 2026  •  "+tanggal, "", 1, "C", false, 0, "")

	pdfBuf := &pdfWriter{}
	if err := pdf.Output(pdfBuf); err != nil {
		return nil, err
	}
	return pdfBuf.buf, nil
}

// GenerateExcel membuat file Excel rekap tahap 1
func GenerateExcel(data *ExportData) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Tahap 1"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	tanggal := time.Now().Format("02 January 2006")

	// Styles
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 13},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
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
	totalStyle, _ := f.NewStyle(&excelize.Style{
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

	// Set kolom width
	f.SetColWidth(sheet, "A", "A", 6)
	f.SetColWidth(sheet, "B", "B", 35)
	f.SetColWidth(sheet, "C", "G", 14)

	// Header dokumen
	f.MergeCell(sheet, "A1", "G1")
	f.SetCellValue(sheet, "A1", "REKAP ENTRY BY SPORT - POPDA 2026")
	f.SetCellStyle(sheet, "A1", "G1", titleStyle)
	f.SetRowHeight(sheet, 1, 22)

	f.MergeCell(sheet, "A2", "G2")
	f.SetCellValue(sheet, "A2", "Kontingen: "+data.Kontingen.NamaKontingen)
	f.SetCellStyle(sheet, "A2", "G2", subtitleStyle)

	f.MergeCell(sheet, "A3", "G3")
	f.SetCellValue(sheet, "A3", "Tanggal Cetak: "+tanggal+" | Status: "+data.Kontingen.Tahap1Status)
	f.SetCellStyle(sheet, "A3", "G3", subtitleStyle)

	// Header tabel (baris 5)
	headers := []string{"No", "Cabang Olahraga", "Atlet Putra", "Atlet Putri", "Pelatih", "Total Atlet", "Total Personel"}
	cols := []string{"A", "B", "C", "D", "E", "F", "G"}
	for i, h := range headers {
		cell := cols[i] + "5"
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	f.SetRowHeight(sheet, 5, 18)

	// Data rows
	for idx, row := range data.CaborList {
		r := idx + 6
		rowStr := fmt.Sprintf("%d", r)
		f.SetCellValue(sheet, "A"+rowStr, idx+1)
		f.SetCellValue(sheet, "B"+rowStr, row.NamaCabor)
		f.SetCellValue(sheet, "C"+rowStr, row.Putra)
		f.SetCellValue(sheet, "D"+rowStr, row.Putri)
		f.SetCellValue(sheet, "E"+rowStr, row.Pelatih)
		f.SetCellValue(sheet, "F"+rowStr, row.TotalAtlet)
		f.SetCellValue(sheet, "G"+rowStr, row.TotalPersonel)
		for _, col := range cols {
			f.SetCellStyle(sheet, col+rowStr, col+rowStr, cellStyle)
		}
	}

	// Total row
	totalRow := fmt.Sprintf("%d", len(data.CaborList)+6)
	f.MergeCell(sheet, "A"+totalRow, "B"+totalRow)
	f.SetCellValue(sheet, "A"+totalRow, "TOTAL")
	f.SetCellValue(sheet, "C"+totalRow, data.TotalPutra)
	f.SetCellValue(sheet, "D"+totalRow, data.TotalPutri)
	f.SetCellValue(sheet, "E"+totalRow, data.TotalPelatih)
	f.SetCellValue(sheet, "F"+totalRow, data.TotalAtlet)
	f.SetCellValue(sheet, "G"+totalRow, data.TotalPersonel)
	for _, col := range cols {
		f.SetCellStyle(sheet, col+totalRow, col+totalRow, totalStyle)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// pdfWriter adalah io.Writer sederhana untuk menampung output gofpdf
type pdfWriter struct {
	buf []byte
}

func (w *pdfWriter) Write(p []byte) (int, error) {
	w.buf = append(w.buf, p...)
	return len(p), nil
}
