package components

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
		"  Например: \"133.7 - Шаверма\"")