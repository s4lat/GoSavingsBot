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

	CSVPrefix   = "/csv"
	ExcelPrefix = "/excel"

	XLSXFullBorder = []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}
	XLSXGrayFill1 = excelize.Fill{Type: "pattern", Color: []string{"#404040"}, Pattern: 1}
	XLSXGrayFill2 = excelize.Fill{Type: "pattern", Color: []string{"#D9D9D9"}, Pattern: 1}
	XLSXGrayFill3 = excelize.Fill{Type: "pattern", Color: []string{"#F2F2F2"}, Pattern: 1}
	XLSXWhiteFont = excelize.Font{Color: "#FFFFFF"}
)
