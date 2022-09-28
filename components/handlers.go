package components

import (
	"fmt"
	"log"
	tele "gopkg.in/telebot.v3"
	"github.com/zsefvlol/timezonemapper"
	"github.com/google/uuid"
	"strings"
	"strconv"
	"time"
	// "reflect"
	"gorm.io/gorm"
)

func TimeZoneHandler(c tele.Context) error { 
	var (
		user = c.Sender()
	)

	r := &tele.ReplyMarkup{ResizeKeyboard: true}

	r.Reply(r.Row(r.Location("Отправить моё местоположение")))
	return c.Send(fmt.Sprintf("Привет, %s!\n", user.Username) +
		"Отправь мне свое местоположение, чтобы я смог установить правильный часовой пояс" +
		"\n\n<i>Если боишься деанонимизации, можешь прикрепить любую геопозицию в том же часовом поясе</i>", 
	r, "HTML")
}

func LocationHandler(c tele.Context) error {
	var (
		user = c.Sender()
		loc = c.Message().Location
		db = c.Get("db").(*gorm.DB)
	)

	timezone := timezonemapper.LatLngToTimezoneString(float64(loc.Lat), float64(loc.Lng))
	tz := TimeZone{UserID: user.ID, TZ: timezone}
	if db.Model(&tz).Where("user_id = ?", user.ID).Updates(&tz).RowsAffected == 0 {
		log.Printf(fmt.Sprintf("Adding info about timezone for %s(%d)", user.Username, user.ID))
	    db.Create(&tz)
	} else {
		log.Printf(fmt.Sprintf("Updating info about timezone for %s(%d)", user.Username, user.ID))
	}
	location, _ := time.LoadLocation(timezone)

	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	btnDaySpends := menu.Text("Сегодня")
	btnMonthSpends := menu.Text("Статистика")
	menu.Reply(
		menu.Row(btnDaySpends,),
		menu.Row(btnMonthSpends,),
	)

	c.Send(fmt.Sprintf("Часовой пояс установлен в: \n<strong>%s</strong>", timezone), 
		"HTML", menu)

	c.Set("loc", location)
	c.Set("tz_name", timezone)
	return DaySpendsHandler(c)
}

func StartHandler(c tele.Context) error {
	var (
		// user = c.Sender()
		// db = c.Get("db").(*gorm.DB)
		tz = c.Get("tz_name").(string)
	)

	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	btnDaySpends := menu.Text("Сегодня")
	btnMonthSpends := menu.Text("Статистика")
	menu.Reply(
		menu.Row(btnDaySpends,),
		menu.Row(btnMonthSpends,),
	)

	c.Send(fmt.Sprintf("Твой часовой пояс: <strong> %s </strong>", tz), menu, "HTML")
	return DaySpendsHandler(c)
}

func DaySpendsHandler(c tele.Context) error {
	var (
		user = c.Sender()
		db = c.Get("db").(*gorm.DB)
		loc = c.Get("loc").(*time.Location)
		date_interface = c.Get("date")
	)

	var date time.Time
	if date_interface != nil {
		date = date_interface.(time.Time)
	} else {
		date = time.Now().In(loc)
	}

	year, month, day := date.Date()

	spends := GetSpendsByDayMonthYear(user.ID, db, day, int(month), year, loc)
	resp := ""

	resp += fmt.Sprintf("Траты за <strong>%02d.%02d</strong> (%d):\n", day, int(month), len(spends))

	var total_spend float32
	for _, spend := range spends {
		total_spend += spend.Value
	}

	if len(spends) > 20 {
		hours, mins, _ := spends[0].Date.In(loc).Clock()
		resp += fmt.Sprintf("  [%02d:%02d] ", hours, mins)
		resp += fmt.Sprintf("%.2f  -  %s (%s)\n      ...  ...  ...\n", spends[0].Value, 
			spends[0].Name, "/del" + strconv.FormatInt(spends[0].ID, 10))

		spends = spends[len(spends) - 10:]
	}

	for _, spend := range spends {
		hours, mins, _ := spend.Date.In(loc).Clock()
		resp += fmt.Sprintf("  [%02d:%02d] ", hours, mins)
		resp += fmt.Sprintf("%.2f  -  %s (%s)\n", spend.Value, 
			spend.Name, "/del" + strconv.FormatInt(spend.ID, 10))
	}

	resp += fmt.Sprintf("\nВсего потрачено: <strong>%.2f</strong>\n", total_spend)

	selector := &tele.ReplyMarkup{}
	selector.Inline(selector.Row(
		selector.Data("<", uuid.NewString(), "get_day", date.AddDate(0, 0, -1).Format("2/1")),
		selector.Data("Сегодня", uuid.NewString(), "get_day", time.Now().In(loc).Format("2/1")),
		selector.Data(">", uuid.NewString(), "get_day", date.AddDate(0, 0, +1).Format("2/1")),
	), 
	selector.Row(
		selector.Data("<<", uuid.NewString(), "get_day", date.AddDate(0, 0, -10).Format("2/1")),
		selector.Data(">>", uuid.NewString(), "get_day", date.AddDate(0, 0, +10).Format("2/1")),
	))

	return c.EditOrSend(resp, selector, "HTML")
}

func YearSpendsHandler(c tele.Context) error {
	var (
		user = c.Sender()
		db = c.Get("db").(*gorm.DB)
		loc = c.Get("loc").(*time.Location)
		year_interface = c.Get("year")
	)

	var year int
	if year_interface != nil {
		year = year_interface.(int)
	} else {
		year = 2022
	}

	year_total, months_totals := GetYearStats(user.ID, db, year, loc)
	resp := fmt.Sprintf("<i>Год: <strong>%d</strong></i>\n", year)

	for i, month_total := range months_totals {
		resp += fmt.Sprintf("%s: <strong>%.2f</strong>\n", int2months[i], month_total)
	}
	resp += fmt.Sprintf("\nВсего потрачено: <strong> %.2f </strong>", year_total)

	selector := &tele.ReplyMarkup{}
	selector.Inline(selector.Row(
		selector.Data("<", uuid.NewString(), "get_year", strconv.Itoa(year - 1)),
		selector.Data(strconv.Itoa(time.Now().In(loc).Year()), uuid.NewString(), 
			"get_year", strconv.Itoa(time.Now().In(loc).Year())),
		selector.Data(">", uuid.NewString(), "get_year", strconv.Itoa(year + 1)),
	))

	return c.EditOrSend(resp, selector, "HTML")
}

func CallbackHandler(c tele.Context) error {
	var (
		args = c.Args()
		loc = c.Get("loc").(*time.Location)
	)

	switch args[1] {
	case "get_day":
		vals := strings.Split(args[2], "/")
		if len(vals) != 2 {
			c.Send("Something went wrong(((")
			return DaySpendsHandler(c)
		}
	
		day, err := strconv.Atoi(vals[0])
		if err != nil {
			c.Send("Something went wrong(((")
			return DaySpendsHandler(c)
		}

		month, err := strconv.Atoi(vals[1])
		if err != nil {
			c.Send("Something went wrong(((")
			return DaySpendsHandler(c)
		}
		
		date := time.Date(time.Now().In(loc).Year(), time.Month(month), day, 0, 0, 0, 0, loc)
		c.Set("date", date)
		return DaySpendsHandler(c)

	case "get_year":
		year, err := strconv.Atoi(args[2])
		if err != nil {
			c.Send("Something went wrong(((")
			return YearSpendsHandler(c)
		}
		
		c.Set("year", year)
		return YearSpendsHandler(c)

	default:
		c.Send("Something went wrong(((")
		return DaySpendsHandler(c)
	}
	return nil
}

func UpdateSpendsHandler(c tele.Context) error {
	var (
		user = c.Sender()
		text = c.Text()
		db = c.Get("db").(*gorm.DB)
	)

	if strings.HasPrefix(text, "/del") {
		return DelSpendHandler(c)
	}

	vals := strings.Split(text, "-")
	if len(vals) != 2 {
		return c.Send("Неправильный формат расходов!")
	}

	val64, err := strconv.ParseFloat(strings.TrimSpace(vals[0]), 64)
	if err != nil {
		return c.Send("Неправильный формат расходов!")
	}

	name := strings.TrimSpace(vals[1])
	value := float32(val64)

	spend := Spend{
		Name: name, 
		Value: value, 
		UserID: user.ID, 
		Date: time.Now(),
	}
	db.Create(&spend)
	return DaySpendsHandler(c)
}

func DelSpendHandler(c tele.Context) error {
	var (
		user = c.Sender()
		text = c.Text()[4:]
		db = c.Get("db").(*gorm.DB)
		loc = c.Get("loc").(*time.Location)
	)

	spend_id, err := strconv.Atoi(text)
	if err != nil {
		return c.Send("Неверный формат команды!")
	}
	var spend Spend
	if db.Find(&spend, "id = ? AND user_id", spend_id, user.ID).RowsAffected == 0 {
		return c.Send("Нет такой траты!")
	}
	db.Delete(&spend)
	c.Send(fmt.Sprintf("Трата <strong>\"%.2f  -  %s\"</strong> - удалена!", spend.Value, spend.Name), "HTML")


	date := time.Date(spend.Date.Year(), spend.Date.Month(), spend.Date.Day(), 0, 0, 0, 0, loc)
	c.Set("date", date)
	return DaySpendsHandler(c)
}