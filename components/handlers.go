package components

import (
	"fmt"
	"log"
	tele "gopkg.in/telebot.v3"
	"github.com/zsefvlol/timezonemapper"
	// "github.com/google/uuid"
	"strings"
	"strconv"
	"time"
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
		// loc = c.Get("loc").(*time.Location)
	)

	var spends []Spend
	db.Model(&spends).Where("user_id = ?", user.ID).Order("date").Find(&spends)

	// resp := "Всего трат в этом месяце: %d"
	resp := "Spends:\n"
	for i, spend := range spends {

		resp += fmt.Sprintf("%d. %6.2f - %s\n", i + 1, spend.Value, spend.Name)
	}

	return c.Send(resp)
}

func AddSpendHandler(c tele.Context) error {
	var (
		user = c.Sender()
		text = c.Text()
		db = c.Get("db").(*gorm.DB)
	)

	vals := strings.Split(text, "-")
	if len(vals) != 2 {
		return c.Send("Неправильный формат расходов!1")
	}

	val64, err := strconv.ParseFloat(strings.TrimSpace(vals[0]), 64)
	if err != nil {
		return c.Send("Неправильный формат расходов!2")
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
