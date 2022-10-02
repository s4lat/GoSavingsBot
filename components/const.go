package components

import (
	"github.com/xuri/excelize/v2"
)

var INT2MONTHS = [12]string{
	"Январь",
	"Февраль",
	"Март",
	"Апрель",
	"Май",
	"Июнь",
	"Июль",
	"Август",
	"Сентябрь",
	"Октябрь",
	"Ноябрь",
	"Декабрь",
}

var HELP_MSG = 	("Для того чтобы пользоваться ботом используй кнопки клавиатуры: \n" +
		"  \"<strong>Сегодня</strong>\" - выводит список сегодняшних трат\n" + 
		"  \"<strong>Статистика</strong>\" - выводит стастику по тратам за год\n\n" + 
		"Чтобы <strong>добавить трату</strong> отправь сообщение в формате:\n  <strong>&lt;число&gt;</strong> - <strong>&lt;наименование траты&gt;</strong>\n" +
		"  Например: \"133.7 - Шаверма\"" + 
		"\n\nЧтобы удалить трату нажми на текст <strong>/delN</strong> рядом с тратой")


var (
	CSV_PREFIX = "/csv"
	EXCEL_PREFIX = "/excel"
)

var (
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