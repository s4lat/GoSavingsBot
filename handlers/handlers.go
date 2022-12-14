package handlers

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zsefvlol/timezonemapper"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"

	"github.com/s4lat/gosavingsbot/export"
	"github.com/s4lat/gosavingsbot/log"
	"github.com/s4lat/gosavingsbot/models"
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
)

// LangAskHandler is a handler for language asking.
// Returned on "/set_lang" and on first use of bot.
func LangAskHandler(c tele.Context) error {
	selector := &tele.ReplyMarkup{}
	selector.Inline(selector.Row(
		selector.Data("🇬🇧 English", uuid.NewString(), "setLang", "en"),
		selector.Data("🇷🇺 Русский", uuid.NewString(), "setLang", "ru"),
	))

	return c.Send("Which language do you prefer?\n\nКакой язык для тебя удобнее?", selector, "HTML")
}

// AskToDeleteUserData is a handler for asking for all user data deletion.
//
// Language required for work.
func AskToDeleteUserData(c tele.Context) error {
	var (
		lang    = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)

	selector := &tele.ReplyMarkup{}
	selector.Inline(selector.Row(
		selector.Data(printer.Sprintf("Yes"), uuid.NewString(), "delete_all_my_data"),
		selector.Data(printer.Sprintf("No"), uuid.NewString(), "cancel"),
	))

	return c.Send(printer.Sprintf(
		"Are you sure you want to delete all your data? This action is <strong>permanent</strong>",
	),
		"HTML", selector)
}

// TimeZoneAskHandler is a handler for time zone asking.
//
// Language required for work.
func TimeZoneAskHandler(c tele.Context) error {
	var (
		lang    = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)

	r := &tele.ReplyMarkup{ResizeKeyboard: true}
	r.Reply(r.Row(r.Location(printer.Sprintf("Send my location"))))

	return c.Send(printer.Sprintf("ASK_LOCATION"), r, "HTML")
}

// StartHandler is a handler that sends menu and help message.
//
// Language required for work.
func StartHandler(c tele.Context) error {
	var (
		lang    = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)

	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(
		menu.Row(menu.Text(printer.Sprintf("Today"))),
		menu.Row(menu.Text(printer.Sprintf("Statistics"))),
		menu.Row(menu.Text(printer.Sprintf("Settings"))),
	)

	return c.Send(printer.Sprintf("HELP_MSG"), menu, "HTML")
}

// SettignsHandler is a handler for user settings.
//
// Language required for work.
func SettingsHandler(c tele.Context) error {
	var (
		userID  = c.Sender().ID
		db      = c.Get("db").(*gorm.DB)
		lang    = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)

	var user models.User
	db.Find(&user, "id = ?", userID)
	return c.Send(printer.Sprintf("SETTINGS_MSG", user.TimeZone), "HTML")
}

// DaySpendsHandler is a handler for day spends stats.
// Have selector for scrolling between days.
// If key "date" in context not set, returning current date stats, else "date" stats.
//
// Location and language required for work.
func DaySpendsHandler(c tele.Context) error {
	var (
		userID  = c.Sender().ID
		db      = c.Get("db").(*gorm.DB)
		loc     = c.Get("loc").(*time.Location)
		lang    = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)

	date, ok := c.Get("date").(time.Time)
	if !ok {
		date = time.Now().In(loc)
	}

	year, month, day := date.Date()
	spends := models.GetSpendsByDayMonthYear(userID, db, day, int(month), year, loc)

	resp := ""
	resp += printer.Sprintf("Spends on <strong>%02d.%02d</strong> (%d):\n", day, int(month), len(spends))

	var totalSpend float32
	for _, spend := range spends {
		totalSpend += spend.Value
	}

	if len(spends) > 20 {
		resp += fmt.Sprintf("  [-] %.2f  -  %s (%s)\n      ...  ...  ...\n", spends[0].Value,
			spends[0].Name, "/del"+strconv.FormatInt(spends[0].ID, 10))

		spends = spends[len(spends)-10:]
	}

	for _, spend := range spends {
		resp += fmt.Sprintf("  [-] %.2f  -  %s (%s)\n", spend.Value,
			spend.Name, "/del"+strconv.FormatInt(spend.ID, 10))
	}

	resp += "\n" + printer.Sprintf("Total spend: <strong>%.2f</strong>\n", totalSpend)

	selector := &tele.ReplyMarkup{}
	selector.Inline(selector.Row(
		selector.Data("<", uuid.NewString(), "getDay", date.AddDate(0, 0, -1).Format("2/1")),
		selector.Data(printer.Sprintf("Today"), uuid.NewString(), "getDay", time.Now().In(loc).Format("2/1")),
		selector.Data(">", uuid.NewString(), "getDay", date.AddDate(0, 0, +1).Format("2/1")),
	),
		selector.Row(
			selector.Data("<<", uuid.NewString(), "getDay", date.AddDate(0, 0, -10).Format("2/1")),
			selector.Data(">>", uuid.NewString(), "getDay", date.AddDate(0, 0, +10).Format("2/1")),
		))

	log.InfoLogger.Printf("Sended DaySpends for %d", userID)
	return c.EditOrSend(resp, selector, "HTML")
}

// YearSpendsHandler is a handler for year spends stats.
// Have selector for scrolling between years.
// Have shortcuts for export csv and excel (/csvYEAR and /excelYEAR).
//
// Location and language required for work.
func YearSpendsHandler(c tele.Context) error {
	var (
		userID  = c.Sender().ID
		db      = c.Get("db").(*gorm.DB)
		loc     = c.Get("loc").(*time.Location)
		lang    = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)

	year, ok := c.Get("year").(int)
	if !ok {
		year = time.Now().In(loc).Year()
	}

	spends := models.GetSpendsByYear(userID, db, year, loc)

	var yearTotal float32
	var monthsTotals [12]float32
	for _, spend := range spends {
		month := int(spend.Date.Month())
		monthsTotals[month-1] += spend.Value
		yearTotal += spend.Value
	}
	resp := printer.Sprintf("Year: <strong>%#d</strong>\n", year)

	for i, monthTotal := range monthsTotals {
		resp += printer.Sprintf("%s: <strong>%.2f</strong>\n", printer.Sprintf(INT2MONTHS[i]), monthTotal)
	}
	resp += "\n" + printer.Sprintf("Total spend: <strong>%.2f</strong>\n", yearTotal)
	resp += fmt.Sprintf("%s%d %s%d", CSVPrefix, year, ExcelPrefix, year)

	selector := &tele.ReplyMarkup{}
	selector.Inline(selector.Row(
		selector.Data("<", uuid.NewString(), "getYear", strconv.Itoa(year-1)),
		selector.Data(strconv.Itoa(time.Now().In(loc).Year()), uuid.NewString(),
			"getYear", strconv.Itoa(time.Now().In(loc).Year())),
		selector.Data(">", uuid.NewString(), "getYear", strconv.Itoa(year+1)),
	))

	log.InfoLogger.Printf("Sended YearStats for %d", userID)
	return c.EditOrSend(resp, selector, "HTML")
}

// CallbackHandler is a handler for callbacks.
// Have shortcuts for export csv/excel (/csvYEAR and /excelYEAR).
//
// Switches depending on args[1]:
//
//	args[1] == "setLang" - trying to set user lang to args[2]
//	args[1] == "getDay" - passing context to DaySpendsHandler with "date" setted to parsed from args[2] day.
//	args[1] == "getYear" - passing context to YearSpendsHandler with "year" setted to parsed from args[2] year.
//	args[1] == "delete_all_my_data" - deleting all user data from database.
//	args[1] == "cancel" - deletes the message from which the callback came.
//
// Language required for work.
func CallbackHandler(c tele.Context) error {
	var (
		args    = c.Args()
		userID  = c.Sender().ID
		db      = c.Get("db").(*gorm.DB)
		lang    = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)

	log.InfoLogger.Printf("Received callback from %d with args %+q", userID, args)
	switch args[1] {
	case "setLang":
		lang, err := language.Parse(args[2])
		if err != nil {
			return LangAskHandler(c)
		}

		user := models.User{ID: userID, Lang: args[2]}
		if db.Model(&user).Where("id = ?", userID).Updates(&user).RowsAffected == 0 {
			db.Create(&user)
			log.InfoLogger.Printf("New user(%d) registered", userID)
		}

		log.InfoLogger.Printf("Language '%s' is set for '%d'", args[2], userID)
		c.EditOrSend(printer.Sprintf("🇬🇧 <strong>English</strong> is selected"), "HTML")

		db.Find(&user, "id = ?", userID)
		c.Set("lang", &lang)
		if len(user.TimeZone) == 0 {
			return TimeZoneAskHandler(c)
		}
		return StartHandler(c)

	case "getDay":
		loc := c.Get("loc").(*time.Location)

		vals := strings.Split(args[2], "/")
		if len(vals) != 2 {
			return c.Send(printer.Sprintf("Something went wrong\n<i>Try sending /start and repeat your actions</i>"))
		}

		day, err := strconv.Atoi(vals[0])
		if err != nil {
			return c.Send(printer.Sprintf("Something went wrong\n<i>Try sending /start and repeat your actions</i>"))
		}

		month, err := strconv.Atoi(vals[1])
		if err != nil {
			return c.Send(printer.Sprintf("Something went wrong\n<i>Try sending /start and repeat your actions</i>"))
		}

		date := time.Date(time.Now().In(loc).Year(), time.Month(month), day, 0, 0, 0, 0, loc)
		c.Set("date", date)
		return DaySpendsHandler(c)

	case "getYear":
		year, err := strconv.Atoi(args[2])
		if err != nil {
			log.ErrorLogger.Print(err)
			return c.Send(printer.Sprintf("Something went wrong\n<i>Try sending /start and repeat your actions</i>"))
		}

		c.Set("year", year)
		return YearSpendsHandler(c)
	case "delete_all_my_data":
		db.Delete(&models.User{}, "id = ?", userID)
		db.Delete(&models.Spend{}, "user_id = ?", userID)

		log.InfoLogger.Printf("Deleted all user data for (%d)", userID)
		return c.Send(printer.Sprintf("All of your data has been erased"))

	case "cancel":
		return c.Delete()
	default:
		c.Send(printer.Sprintf("Something went wrong\n<i>Try sending /start and repeat your actions</i>"))
		return DaySpendsHandler(c)
	}
}

// OnTextHandler is a handler for:
//  1. Spends add msg: <cost> <spend_name>
//  2. /delN command
//  3. /excelYEAR command
//  4. /csvYEAR command
func OnTextHandler(c tele.Context) error {
	text := c.Text()

	switch {
	case strings.HasPrefix(text, "/del"):
		return DelSpendHandler(c)
	case strings.HasPrefix(text, CSVPrefix):
		return ExportHandler(c)
	case strings.HasPrefix(text, ExcelPrefix):
		return ExportHandler(c)
	default:
		return AddSpendHandler(c)
	}
}

// LocationHandler is a handler for location messages.
// Used for calculating user time zone.
func LocationHandler(c tele.Context) error {
	var (
		userID = c.Sender().ID
		loc    = c.Message().Location
		db     = c.Get("db").(*gorm.DB)
	)

	// Converting location latitude and longitude to timezone string
	timezone := timezonemapper.LatLngToTimezoneString(float64(loc.Lat), float64(loc.Lng))
	user := models.User{ID: userID, TimeZone: timezone}
	if db.Model(&user).Where("id = ?", userID).Updates(&user).RowsAffected == 0 {
		db.Create(&user)
		log.InfoLogger.Printf("New user(%d) registered", userID)
	}

	log.InfoLogger.Printf("Setted timezone to '%s' for '%d'", timezone, userID)
	location, _ := time.LoadLocation(timezone)

	c.Set("loc", location)
	return StartHandler(c)
}

// AddSpendHandler is a handler that adds users spend.
//
// Language required for work.
func AddSpendHandler(c tele.Context) error {
	var (
		userID  = c.Sender().ID
		text    = c.Text()
		db      = c.Get("db").(*gorm.DB)
		lang    = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)

	// splitting text message by spaces
	vals := strings.Split(text, " ")
	if len(vals) < 2 {
		return c.Send(printer.Sprintf("Wrong spend format!\n/help - for more info"), "HTML")
	}

	// trim spaces and replacing "," -> "." in number before parsing
	val64, err := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(vals[0], ",", ".")), 64)
	if err != nil {
		return c.Send(printer.Sprintf("Wrong spend format!\n/help - for more info"), "HTML")
	}

	// using all text excluding first number for name
	name := strings.TrimSpace(strings.Join(vals[1:], " "))
	value := float32(val64)

	spend := models.Spend{
		Name:   name,
		Value:  value,
		UserID: userID,
		Date:   time.Now(),
	}
	db.Create(&spend)

	log.InfoLogger.Printf("Added spend for '%d'", userID)
	return DaySpendsHandler(c)
}

// DelSpendHandler is a handler for spend deleting.
//
// Language and location required for work.
func DelSpendHandler(c tele.Context) error {
	var (
		userID  = c.Sender().ID
		text    = c.Text()[4:] // Cutting "/del" substr from message txt
		db      = c.Get("db").(*gorm.DB)
		loc     = c.Get("loc").(*time.Location)
		lang    = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)

	// converting N from '/delN' to int and using as spendID
	spendID, err := strconv.Atoi(text)
	if err != nil {
		return c.Send(printer.Sprintf("Wrong command format!"))
	}

	var spend models.Spend
	if db.Find(&spend, "id = ? AND user_id = ?", spendID, userID).RowsAffected == 0 {
		return c.Send(printer.Sprintf("There is no such spend"))
	}
	db.Delete(&spend)
	log.InfoLogger.Printf("Deleted spend %d for '%d'", spendID, userID)

	c.Send(printer.Sprintf("Spend <strong>\"%.2f  -  %s\"</strong> has been deleted!", spend.Value, spend.Name),
		"HTML")

	date := time.Date(spend.Date.Year(), spend.Date.Month(), spend.Date.Day(), 0, 0, 0, 0, loc)
	c.Set("date", date)
	return DaySpendsHandler(c)
}

// ExportHandler is a handler for export year spends to excel or csv.
//
// Language and location required for work.
func ExportHandler(c tele.Context) error {
	var (
		userID  = c.Sender().ID
		text    = c.Text()
		db      = c.Get("db").(*gorm.DB)
		loc     = c.Get("loc").(*time.Location)
		lang    = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)

	// checking if message has csv or excel prefix, setting 'start' to len(prefix)
	// start used for cutting substr "/csv" or "/excel" from msg
	var start int
	if strings.HasPrefix(text, CSVPrefix) {
		start = len(CSVPrefix)
	} else if strings.HasPrefix(text, ExcelPrefix) {
		start = len(ExcelPrefix)
	}

	year, err := strconv.Atoi(text[start:])
	if err != nil {
		return c.Send(printer.Sprintf("Wrong command format!"))
	}
	spends := models.GetSpendsByYear(userID, db, year, loc)
	if len(spends) == 0 {
		return c.Send(printer.Sprintf("No spends during this period"))
	}

	// reader used as buffer for generated csv/excel file content
	var reader *bytes.Buffer
	var filename string

	switch {
	case strings.HasPrefix(text, CSVPrefix):
		{
			filename = fmt.Sprintf("%04d.csv", year)
			reader, err = export.SpendsToCSV(spends)
			if err != nil {
				log.ErrorLogger.Print(err)
				return c.Send(printer.Sprintf("Something went wrong\nTry sending /start and repeat your actions"),
					"HTML")
			}
		}
	case strings.HasPrefix(text, ExcelPrefix):
		{
			filename = fmt.Sprintf("%04d.xlsx", year)

			reader, err = export.SpendsToExcel(spends, printer)
			if err != nil {
				log.ErrorLogger.Print(err)
				return c.Send(printer.Sprintf("Something went wrong\nTry sending /start and repeat your actions"),
					"HTML")
			}
		}
	default:
		{
			return c.Send(printer.Sprintf("Wrong command format!"))
		}
	}

	log.InfoLogger.Printf("Sendend '%s' to '%d'", filename, userID)
	file := &tele.Document{File: tele.FromReader(reader), FileName: filename}

	return c.Send(file)
}
