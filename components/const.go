package components

import (
	"github.com/xuri/excelize/v2"
)

var (
	INT2MONTHS = [12]string{
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	}

	CSV_PREFIX = "/csv"
	EXCEL_PREFIX = "/excel"

	XLSX_FULL_BORDER = []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
	}
	XLSX_GRAY1_FILL = excelize.Fill{Type: "pattern", Color: []string{"#404040"}, Pattern: 1}
	XLSX_GRAY2_FILL = excelize.Fill{Type: "pattern", Color: []string{"#D9D9D9"}, Pattern: 1}
	XLSX_GRAY3_FILL = excelize.Fill{Type: "pattern", Color: []string{"#F2F2F2"}, Pattern: 1}
	XLSX_WHITE_FONT = excelize.Font{Color: "#FFFFFF"}
)