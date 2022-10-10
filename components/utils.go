package components

import (
	"encoding/csv"
	"bytes"
	"fmt"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/message"
)

func SpendsToCSV(spends []Spend) *bytes.Buffer {
	records := make([][]string, len(spends) + 1)

	records[0] = []string{"name", "value", "clock", "date"}
	for i := 0; i < len(spends); i++ {
		hours, mins, _ := spends[i].Date.Clock()
		year, month, day := spends[i].Date.Date()
		records[i + 1] = []string{
			spends[i].Name, 
			fmt.Sprintf("%0.2f", spends[i].Value), 
			fmt.Sprintf("%02d:%02d", hours, mins),
			fmt.Sprintf("%02d.%02d.%d", day, month, year),
		}
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	w.WriteAll(records)

	return &buf
}

func SpendsToExcel(spends []Spend, printer *message.Printer) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	
	header_style, _ := f.NewStyle(&excelize.Style{Border: XLSX_FULL_BORDER, Fill: XLSX_GRAY1_FILL, Font: &XLSX_WHITE_FONT})
	col_style1, _ := f.NewStyle(&excelize.Style{Border: XLSX_FULL_BORDER, Fill: XLSX_GRAY2_FILL,})
	col_style2, _ := f.NewStyle(&excelize.Style{Border: XLSX_FULL_BORDER, Fill: XLSX_GRAY3_FILL,})

	f.NewSheet(printer.Sprintf("Total"))
	f.SetCellValue(printer.Sprintf("Total"), "A1", printer.Sprintf("Month"))
	f.SetCellValue(printer.Sprintf("Total"), "B1", printer.Sprintf("Spended"))
	f.SetCellStyle(printer.Sprintf("Total"), "A1", "B1", header_style)

	for i, month := range INT2MONTHS {
		f.NewSheet(printer.Sprintf(month))
		f.SetCellValue(printer.Sprintf(month), "A1", printer.Sprintf("Spend name"))
		f.SetColWidth(printer.Sprintf(month), "A", "A", 40)
		f.SetCellValue(printer.Sprintf(month), "B1", printer.Sprintf("Spended"))
		f.SetColWidth(printer.Sprintf(month), "B", "B", 10)
		f.SetCellValue(printer.Sprintf(month), "C1", printer.Sprintf("Clock"))
		f.SetColWidth(printer.Sprintf(month), "C", "C", 8)
		f.SetCellValue(printer.Sprintf(month), "D1", printer.Sprintf("Date"))
		f.SetColWidth(printer.Sprintf(month), "D", "D", 10)

		f.SetCellValue(printer.Sprintf(month), "E2", printer.Sprintf("Total"))
		f.SetCellFormula(printer.Sprintf(month), "F2", "=SUM(B:B)")
		f.SetColWidth(printer.Sprintf(month), "E", "E", 6)
		f.SetCellStyle(printer.Sprintf(month), "E2", "F2", header_style)

		f.SetCellValue(printer.Sprintf("Total"), fmt.Sprintf("A%d", i+2), printer.Sprintf(month))
		f.SetCellFormula(printer.Sprintf("Total"), fmt.Sprintf("B%d", i+2), fmt.Sprintf("=%s!F2", printer.Sprintf(month)))
		if i % 2 == 0 {
			f.SetCellStyle(printer.Sprintf("Total"), fmt.Sprintf("A%d",i+2), fmt.Sprintf("B%d", i+2), col_style1)
		} else {
			f.SetCellStyle(printer.Sprintf("Total"), fmt.Sprintf("A%d", i+2), fmt.Sprintf("B%d", i+2), col_style2)
		}
	}

	f.SetCellValue(printer.Sprintf("Total"), "A14", printer.Sprintf("Total"))
	f.SetCellFormula(printer.Sprintf("Total"), "B14", "=SUM(B2:B13)")
	f.SetCellStyle(printer.Sprintf("Total"), "A14", "B14", header_style)
	row_by_month := []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	for _, spend := range spends {
		month := int(spend.Date.Month())

		hours, mins, _ := spend.Date.Clock()
		year, _, day := spend.Date.Date()

		sheet := printer.Sprintf(INT2MONTHS[month - 1])
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row_by_month[month - 1]), spend.Name)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row_by_month[month - 1]), spend.Value)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row_by_month[month - 1]), 
			fmt.Sprintf("%02d:%02d", hours, mins))
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row_by_month[month - 1]), 
			fmt.Sprintf("%02d.%02d.%04d", day, month, year))
		row_by_month[month - 1] += 1
	}

	for i, month := range INT2MONTHS {
		f.SetCellStyle(printer.Sprintf(month), "A1", "D1", header_style)
		for j := 2; j < row_by_month[i]; j++ {
			if j % 2 == 0 {
				f.SetCellStyle(printer.Sprintf(month), fmt.Sprintf("A%d", j), fmt.Sprintf("D%d", j), col_style1)
			} else {
				f.SetCellStyle(printer.Sprintf(month), fmt.Sprintf("A%d", j), fmt.Sprintf("D%d", j), col_style2)
			}
		}
	}

	f.DeleteSheet("Sheet1")
	return f.WriteToBuffer()
}