package components

import (
	"encoding/csv"
	"bytes"
	"fmt"
	"github.com/xuri/excelize/v2"
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

func SpendsToExcel(spends []Spend) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	
	header_style, _ := f.NewStyle(&excelize.Style{Border: XLSX_FULL_BORDER, Fill: XLSX_GRAY1_FILL, Font: &XLSX_WHITE_FONT})
	col_style1, _ := f.NewStyle(&excelize.Style{Border: XLSX_FULL_BORDER, Fill: XLSX_GRAY2_FILL,})
	col_style2, _ := f.NewStyle(&excelize.Style{Border: XLSX_FULL_BORDER, Fill: XLSX_GRAY3_FILL,})

	f.NewSheet("Общее")
	f.SetCellValue("Общее", "A1", "Месяц")
	f.SetCellValue("Общее", "B1", "Потрачено")
	f.SetCellStyle("Общее", "A1", "B1", header_style)

	for i, month := range INT2MONTHS {
		f.NewSheet(month)
		f.SetCellValue(month, "A1", "Название траты")
		f.SetColWidth(month, "A", "A", 40)
		f.SetCellValue(month, "B1", "Потрачено")
		f.SetColWidth(month, "B", "B", 10)
		f.SetCellValue(month, "C1", "Время")
		f.SetColWidth(month, "C", "C", 8)
		f.SetCellValue(month, "D1", "Дата")
		f.SetColWidth(month, "D", "D", 10)

		f.SetCellValue(month, "E2", "Всего")
		f.SetCellFormula(month, "F2", "=SUM(B:B)")
		f.SetColWidth(month, "E", "E", 6)
		f.SetCellStyle(month, "E2", "F2", header_style)

		f.SetCellValue("Общее", fmt.Sprintf("A%d", i+2), month)
		f.SetCellFormula("Общее", fmt.Sprintf("B%d", i+2), fmt.Sprintf("=%s!F2", month))
		if i % 2 == 0 {
			f.SetCellStyle("Общее", fmt.Sprintf("A%d",i+2), fmt.Sprintf("B%d", i+2), col_style1)
		} else {
			f.SetCellStyle("Общее", fmt.Sprintf("A%d", i+2), fmt.Sprintf("B%d", i+2), col_style2)
		}
	}

	f.SetCellValue("Общее", "A14", "Всего")
	f.SetCellFormula("Общее", "B14", "=SUM(B2:B13)")
	f.SetCellStyle("Общее", "A14", "B14", header_style)
	row_by_month := []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	for _, spend := range spends {
		month := int(spend.Date.Month())

		hours, mins, _ := spend.Date.Clock()
		year, _, day := spend.Date.Date()

		sheet := INT2MONTHS[month - 1]
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row_by_month[month - 1]), spend.Name)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row_by_month[month - 1]), spend.Value)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row_by_month[month - 1]), 
			fmt.Sprintf("%02d:%02d", hours, mins))
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row_by_month[month - 1]), 
			fmt.Sprintf("%02d.%02d.%04d", day, month, year))
		row_by_month[month - 1] += 1
	}

	for i, month := range INT2MONTHS {
		f.SetCellStyle(month, "A1", "D1", header_style)
		for j := 2; j < row_by_month[i]; j++ {
			if j % 2 == 0 {
				f.SetCellStyle(month, fmt.Sprintf("A%d", j), fmt.Sprintf("D%d", j), col_style1)
			} else {
				f.SetCellStyle(month, fmt.Sprintf("A%d", j), fmt.Sprintf("D%d", j), col_style2)
			}
		}
	}

	f.DeleteSheet("Sheet1")
	return f.WriteToBuffer()
}