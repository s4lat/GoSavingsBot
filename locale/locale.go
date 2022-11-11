//nolint:lll
package locale

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// InitLocales - initializing locales for supported languages.
func InitLocales() {
	message.SetString(language.Russian,
		"HELP_MSG",
		("Для того чтобы пользоваться ботом используй кнопки клавиатуры: \n" +
			"  \"<strong>Сегодня</strong>\" - выводит список сегодняшних трат\n" +
			"  \"<strong>Статистика</strong>\" - выводит стастику по тратам за год\n\n" +
			"Чтобы <strong>добавить трату</strong> отправь сообщение в формате:\n" +
			"  <strong>&lt;число&gt;</strong> - <strong>&lt;наименование траты&gt;</strong>\n" +
			"  Например: \"133.7 новые кроссовки\"\n" +
			"\nЧтобы удалить трату нажми на текст <strong>/delN</strong> рядом с тратой"),
	)

	message.SetString(language.English,
		"HELP_MSG",
		("To interact with the bot, use the keyboard buttons: \n" +
			"  \"<strong>Today</strong>\" - displays a list of today's spends\n" +
			"  \"<strong>Statistics</strong>\" - displays spending statistics for the year\n\n" +
			"To <strong>add spend</strong> send a message in the format:\n" +
			"  <strong>&lt;number&gt;</strong> <strong>&lt;spend name&gt;</strong>\n" +
			"  For example: \"133.7 new shoes\"\n" +
			"\nTo delete a spend, click on the text <strong>/delN</strong> next to the spend"),
	)

	message.SetString(language.Russian,
		"Send my location",
		"Отправить моё местоположение",
	)

	message.SetString(language.Russian,
		"ASK_LOCATION",
		("Отправь мне свое местоположение, чтобы я смог установить правильный часовой пояс\n" +
			"\n<i>Если боишься деанонимизации, можешь прикрепить любую геопозицию в том же часовом поясе</i>"),
	)

	message.SetString(language.English,
		"ASK_LOCATION",
		("Send me your location so I can set the correct time zone\n" +
			"\n<i>If you are afraid of being deanonymized, you can attach any location in the same time zone</i>"),
	)

	message.SetString(language.Russian,
		"Your time zone: <strong> %s </strong>",
		"Твой часовой пояс: <strong> %s </strong>",
	)

	message.SetString(language.Russian,
		"Spends on <strong>%02d.%02d</strong> (%d):\n",
		"Траты за <strong>%02d.%02d</strong> (%d):\n",
	)

	message.SetString(language.Russian,
		"Total spend: <strong>%.2f</strong>\n",
		"Всего потрачено: <strong>%.2f</strong>\n",
	)

	message.SetString(language.Russian,
		"Year: <strong>%#d</strong>\n",
		"Год: <strong>%#d</strong>\n",
	)

	message.SetString(language.Russian,
		"Что-то пошло не так, попробуйте еще раз",
		"Something went wrong, try again",
	)

	message.SetString(language.Russian,
		"Wrong command format!",
		"Неверный формат команды!",
	)

	message.SetString(language.Russian,
		"No spends during this period",
		"Нет трат за этот период",
	)

	message.SetString(language.Russian,
		"There is no such spend",
		"Нет такой траты!",
	)

	message.SetString(language.Russian,
		"Spend <strong>\"%.2f  -  %s\"</strong> has been deleted!",
		"Трата <strong>\"%.2f  -  %s\"</strong> была удалена!",
	)

	message.SetString(language.Russian,
		"Wrong spend format!\n/help - for more info",
		"Неправильный формат расходов!\n/help - как пользоваться ботом",
	)

	message.SetString(language.Russian,
		"Something went wrong\nTry sending /start and repeat your actions",
		"Что-то пошло не так!\nПопробуйте отправить /start и повторить свои действия",
	)

	message.SetString(language.Russian,
		"🇬🇧 <strong>English</strong> is selected",
		"Выбран 🇷🇺 <strong>Русский</strong> язык",
	)

	message.SetString(language.Russian,
		"SETTINGS_MSG",
		("Язык: 🇷🇺 <strong>Русский</strong> (/set_lang)\n" +
			"Твой часовой пояс: <strong>%s</strong>\n<i>(чтобы изменить его отправь мне новую геопозицию)</i>\n\n" +
			"/delete_my_data <i>&lt;= жми сюда, если хочешь удалить всю свою информацию из базы данных бота</i>"),
	)

	message.SetString(language.English,
		"SETTINGS_MSG",
		("Language: 🇬🇧 <strong>English</strong> (/set_lang)\n" +
			"Time zone: <strong>%s</strong>\n<i>(to change it send me new location)</i>\n\n" +
			"/delete_my_data <i>&lt;= click here if you want to delete all your information from the bot's database</i>"),
	)

	message.SetString(language.Russian,
		"Are you sure you want to delete all your data? This action is <strong>permanent</strong>",
		"Вы уверены, что хотите удалить все свои данные? Это действие <strong>необратимо</strong>",
	)
	message.SetString(language.Russian,
		"All of your data has been erased",
		"Все твои данные были стерты",
	)

	message.SetString(language.Russian, "Yes", "Да")
	message.SetString(language.Russian, "No", "Нет")

	message.SetString(language.Russian, "Today", "Сегодня")
	message.SetString(language.Russian, "Statistics", "Статистика")
	message.SetString(language.Russian, "Settings", "Настройки")

	message.SetString(language.Russian, "Total", "Всего")
	message.SetString(language.Russian, "Month", "Месяц")
	message.SetString(language.Russian, "Spended", "Потрачено")
	message.SetString(language.Russian, "Spend name", "Название траты")
	message.SetString(language.Russian, "Clock", "Время")
	message.SetString(language.Russian, "Date", "Дата")

	message.SetString(language.Russian, "January", "Январь")
	message.SetString(language.Russian, "February", "Февраль")
	message.SetString(language.Russian, "March", "Март")
	message.SetString(language.Russian, "April", "Апрель")
	message.SetString(language.Russian, "May", "Май")
	message.SetString(language.Russian, "June", "Июнь")
	message.SetString(language.Russian, "July", "Июль")
	message.SetString(language.Russian, "August", "Август")
	message.SetString(language.Russian, "September", "Сентябрь")
	message.SetString(language.Russian, "October", "Октябрь")
	message.SetString(language.Russian, "November", "Ноябрь")
	message.SetString(language.Russian, "December", "Декабрь")
}
