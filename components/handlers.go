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

	c.Send(fmt.Sprintf("Часовой пояс установлен в: \n<strong>%s</strong>", timezone), 
		"HTML")

	location, _ := time.LoadLocation(timezone)
	c.Set("loc", location)
	return HomeHandler(c)
}

func HomeHandler(c tele.Context) error {
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

	cutted := 0
	if len(spends) > 20 {
		resp += fmt.Sprintf("  %d. %-6.2f - %6s\n      ...  ...  ...\n", 1, spends[0].Value, spends[0].Name)
		cutted = len(spends) - 10

		spends = spends[len(spends) - 10:]
	}

	for i, spend := range spends {
		resp += fmt.Sprintf("  %3d. %-6.2f - %6s\n", i + 1 + cutted, spend.Value, spend.Name)
	}

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
			return HomeHandler(c)
		}
	
		day, err := strconv.Atoi(vals[0])
		if err != nil {
			c.Send("Something went wrong(((")
			return HomeHandler(c)
		}

		month, err := strconv.Atoi(vals[1])
		if err != nil {
			c.Send("Something went wrong(((")
			return HomeHandler(c)
		}
		
		date := time.Date(time.Now().In(loc).Year(), time.Month(month), day, 0, 0, 0, 0, loc)
		c.Set("date", date)
		return HomeHandler(c)

	default:
		c.Send("Something went wrong(((")
		return HomeHandler(c)
	}
	return nil
}

func AddSpendHandler(c tele.Context) error {
	var (
		user = c.Sender()
		text = c.Text()
		db = c.Get("db").(*gorm.DB)
	)

	vals := strings.Split(text, "-")
	if len(vals) != 2 {
		return c.Send("Неправильный формат расходов!")
	}

	val64, err := strconv.ParseFloat(strings.TrimSpace(vals[0]), 64)
	if err != nil {
		return c.Send("Неправильный формат расходов!")
	}
	name := vals[1]
	value := float32(val64)

	spend := Spend{
		Name: name, 
		Value: value, 
		UserID: user.ID, 
		Date: time.Now(),
	}
	db.Create(&spend)
	return HomeHandler(c)
}