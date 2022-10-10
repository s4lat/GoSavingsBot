package components

import (
	"fmt"
	"log"
	tele "gopkg.in/telebot.v3"
	"github.com/zsefvlol/timezonemapper"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/google/uuid"
	"strings"
	"strconv"
	"time"
	"bytes"
	"gorm.io/gorm"
)

func LangAskHandler(c tele.Context) error {
	selector := &tele.ReplyMarkup{}
	selector.Inline(selector.Row(
		selector.Data("🏴󠁧󠁢󠁥󠁮󠁧󠁿 English", uuid.NewString(), "setLang", "en"),
		selector.Data("🇷🇺 Русский", uuid.NewString(), "setLang", "ru"),
	))

	return c.Send("Which language do you prefer?\n\nКакой язык для тебя удобнее?", selector, "HTML")
}

func AskToDeleteUserData(c tele.Context) error {
	var (
		lang = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)

	selector := &tele.ReplyMarkup{}
	selector.Inline(selector.Row(
		selector.Data(printer.Sprintf("Yes"), uuid.NewString(), "delete_all_my_data"),
		selector.Data(printer.Sprintf("No"), uuid.NewString(), "cancel"),
	))

	return c.Send(printer.Sprintf("Are you sure you want to delete all your data? This action is <strong>permanent</strong>"), 
	"HTML", selector)
}

func TimeZoneAskHandler(c tele.Context) error { 
	var (
		lang = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)

	r := &tele.ReplyMarkup{ResizeKeyboard: true}
	r.Reply(r.Row(r.Location(printer.Sprintf("Send my location"))))

	return c.Send(printer.Sprintf("ASK_LOCATION"), r, "HTML")
}

func StartHandler(c tele.Context) error {
	var (
		lang = c.Get("lang").(*language.Tag)
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

func SettingsHandler(c tele.Context) error {
	var (
		userID = c.Sender().ID
		db = c.Get("db").(*gorm.DB)
		lang = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)

	var user User
	db.Find(&user, "id = ?", userID)
	return c.Send(printer.Sprintf("SETTINGS_MSG", user.TimeZone), "HTML")
}

func DaySpendsHandler(c tele.Context) error {
	var (
		userID = c.Sender().ID
		db = c.Get("db").(*gorm.DB)
		loc = c.Get("loc").(*time.Location)
		dateInterface = c.Get("date")
		lang = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)

	var date time.Time
	if dateInterface != nil {
		date = dateInterface.(time.Time)
	} else {
		date = time.Now().In(loc)
	}

	year, month, day := date.Date()
	spends := GetSpendsByDayMonthYear(userID, db, day, int(month), year, loc)
	
	resp := ""
	resp += printer.Sprintf("Spends on <strong>%02d.%02d</strong> (%d):\n", day, int(month), len(spends))

	var totalSpend float32
	for _, spend := range spends {
		totalSpend += spend.Value
	}

	if len(spends) > 20 {
		resp += fmt.Sprintf("  [-] %.2f  -  %s (%s)\n      ...  ...  ...\n", spends[0].Value, 
			spends[0].Name, "/del" + strconv.FormatInt(spends[0].ID, 10))

		spends = spends[len(spends) - 10:]
	}

	for _, spend := range spends {
		resp += fmt.Sprintf("  [-] %.2f  -  %s (%s)\n", spend.Value, 
			spend.Name, "/del" + strconv.FormatInt(spend.ID, 10))
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

	return c.EditOrSend(resp, selector, "HTML")
}

func YearSpendsHandler(c tele.Context) error {
	var (
		userID = c.Sender().ID
		db = c.Get("db").(*gorm.DB)
		loc = c.Get("loc").(*time.Location)
		year_interface = c.Get("year")
		lang = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)


	var year int
	if year_interface != nil {
		year = year_interface.(int)
	} else {
		year = 2022
	}

	spends := GetSpendsByYear(userID, db, year, loc)

	var year_total float32
	var months_totals [12]float32
	for _, spend := range spends {
		month := int(spend.Date.Month())
		months_totals[month - 1] += spend.Value
		year_total += spend.Value
	}
	resp := printer.Sprintf("Year: <strong>%#d</strong>\n", year)

	for i, month_total := range months_totals {
		resp += printer.Sprintf("%s: <strong>%.2f</strong>\n", printer.Sprintf(INT2MONTHS[i]), month_total)
	}
	resp += "\n" + printer.Sprintf("Total spend: <strong>%.2f</strong>\n", year_total)
	resp += fmt.Sprintf("%s%d %s%d", CSV_PREFIX, year, EXCEL_PREFIX, year)

	selector := &tele.ReplyMarkup{}
	selector.Inline(selector.Row(
		selector.Data("<", uuid.NewString(), "getYear", strconv.Itoa(year - 1)),
		selector.Data(strconv.Itoa(time.Now().In(loc).Year()), uuid.NewString(), 
			"getYear", strconv.Itoa(time.Now().In(loc).Year())),
		selector.Data(">", uuid.NewString(), "getYear", strconv.Itoa(year + 1)),
	))

	return c.EditOrSend(resp, selector, "HTML")
}

func CallbackHandler(c tele.Context) error {
	var (
		args = c.Args()
		userID = c.Sender().ID 
		db = c.Get("db").(*gorm.DB)
		lang = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)

	switch args[1] {
	case "setLang":
		lang, err := language.Parse(args[2])
		if err != nil {
			return LangAskHandler(c)
		}

		user := User{ID: userID, Lang: args[2]}
		if db.Model(&user).Where("id = ?", userID).Updates(&user).RowsAffected == 0 {
			db.Create(&user)
		}

		c.EditOrSend(printer.Sprintf("🏴󠁧󠁢󠁥󠁮󠁧󠁿 <strong>English</strong> is selected"), "HTML")

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
			c.Send(printer.Sprintf("Something went wrong\n<i>Try sending /start and repeat your actions</i>"))
			return DaySpendsHandler(c)
		}
	
		day, err := strconv.Atoi(vals[0])
		if err != nil {
			c.Send(printer.Sprintf("Something went wrong\n<i>Try sending /start and repeat your actions</i>"))
			return DaySpendsHandler(c)
		}

		month, err := strconv.Atoi(vals[1])
		if err != nil {
			c.Send(printer.Sprintf("Something went wrong\n<i>Try sending /start and repeat your actions</i>"))
			return DaySpendsHandler(c)
		}
		
		date := time.Date(time.Now().In(loc).Year(), time.Month(month), day, 0, 0, 0, 0, loc)
		c.Set("date", date)
		return DaySpendsHandler(c)

	case "getYear":
		year, err := strconv.Atoi(args[2])
		if err != nil {
			c.Send(printer.Sprintf("Something went wrong\n<i>Try sending /start and repeat your actions</i>"))
			return YearSpendsHandler(c)
		}
		
		c.Set("year", year)
		return YearSpendsHandler(c)
	case "delete_all_my_data":
		db.Delete(&User{}, "id = ?", userID)
		db.Delete(&Spend{}, "user_id = ?", userID)

		return c.Send(printer.Sprintf("All of your data has been erased"))

	case "cancel":
		return c.Delete()
	default:
		c.Send(printer.Sprintf("Something went wrong\n<i>Try sending /start and repeat your actions</i>"))
		return DaySpendsHandler(c)
	}
	return nil
}

func OnTextHandler(c tele.Context) error {
	text := c.Text()

	if strings.HasPrefix(text, "/del") {
		return DelSpendHandler(c)
	} 

	if strings.HasPrefix(text, CSV_PREFIX) {
		return ExportHandler(c)
	}

	if strings.HasPrefix(text, EXCEL_PREFIX) {
		return ExportHandler(c)
	}
	return AddSpendHandler(c)
}

func LocationHandler(c tele.Context) error {
	var (
		userID = c.Sender().ID
		loc = c.Message().Location
		db = c.Get("db").(*gorm.DB)
	)

	timezone := timezonemapper.LatLngToTimezoneString(float64(loc.Lat), float64(loc.Lng))
	user := User{ID: userID, TimeZone: timezone}
	if db.Model(&user).Where("id = ?", userID).Updates(&user).RowsAffected == 0 {
		log.Printf(fmt.Sprintf("Adding info about timezone for %d", userID))
	    db.Create(&user)
	} else {
		log.Printf(fmt.Sprintf("Updating info about timezone for %d", userID))
	}
	location, _ := time.LoadLocation(timezone)

	c.Set("loc", location)
	return StartHandler(c)
}

func AddSpendHandler(c tele.Context) error {
	var (
		userID = c.Sender().ID
		text = c.Text()
		db = c.Get("db").(*gorm.DB)
		lang = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)


	vals := strings.Split(text, " ")
	if len(vals) < 2 {
		return c.Send(printer.Sprintf("Wrong spend format!\n/help - for more info"), "HTML")
	}

	// trim spaces and replacing "," -> "." before parsing
	val64, err := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(vals[0], ",", ".")), 64)
	if err != nil {
		return c.Send(printer.Sprintf("Wrong spend format!\n/help - for more info"), "HTML")
	}

	name := strings.TrimSpace(strings.Join(vals[1:], " "))
	value := float32(val64)

	spend := Spend{
		Name: name, 
		Value: value, 
		UserID: userID, 
		Date: time.Now(),
	}
	db.Create(&spend)
	return DaySpendsHandler(c)
}

func DelSpendHandler(c tele.Context) error {
	var (
		userID = c.Sender().ID
		text = c.Text()[4:]
		db = c.Get("db").(*gorm.DB)
		loc = c.Get("loc").(*time.Location)
		lang = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)

	spendID, err := strconv.Atoi(text)
	if err != nil {
		return c.Send(printer.Sprintf("Wrong command format!"))
	}
	var spend Spend
	if db.Find(&spend, "id = ? AND userID", spendID, userID).RowsAffected == 0 {
		return c.Send(printer.Sprintf("There is no such spend"))
	}
	db.Delete(&spend)
	c.Send(printer.Sprintf("Spend <strong>\"%.2f  -  %s\"</strong> has been deleted!", spend.Value, spend.Name), 
		"HTML")


	date := time.Date(spend.Date.Year(), spend.Date.Month(), spend.Date.Day(), 0, 0, 0, 0, loc)
	c.Set("date", date)
	return DaySpendsHandler(c)
}

func ExportHandler(c tele.Context) error {
	var (
		userID = c.Sender().ID
		text = c.Text()
		db = c.Get("db").(*gorm.DB)
		loc = c.Get("loc").(*time.Location)
		lang = c.Get("lang").(*language.Tag)
		printer = message.NewPrinter(*lang)
	)

	var start int
	if strings.HasPrefix(text, CSV_PREFIX) {
		start = len(CSV_PREFIX)
	} else if strings.HasPrefix(text, EXCEL_PREFIX) {
		start = len(EXCEL_PREFIX)
	}

	year, err := strconv.Atoi(text[start:])
	if err != nil {
		return c.Send(printer.Sprintf("Wrong command format!"))
	}
	spends := GetSpendsByYear(userID, db, year, loc)
	if len(spends) == 0 {
		return c.Send(printer.Sprintf("No spends during this period"))
	}

	var reader *bytes.Buffer
	var filename string
	if strings.HasPrefix(text, CSV_PREFIX) {
		filename = fmt.Sprintf("%04d.csv", year)
		reader = SpendsToCSV(spends)

	} else if strings.HasPrefix(text, EXCEL_PREFIX) {
		filename = fmt.Sprintf("%04d.xlsx", year)

		reader, err = SpendsToExcel(spends, printer)
		if err != nil {
			return c.Send(printer.Sprintf("Something went wrong\nTry sending /start and repeat your actions"), 
			"HTML")
		}
	} else {
		return c.Send(printer.Sprintf("Wrong command format!"))
	}

	file := &tele.Document{File: tele.FromReader(reader), FileName: filename}

	return c.Send(file)
}