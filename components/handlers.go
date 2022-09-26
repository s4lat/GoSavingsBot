package components

import (
	"fmt"
	"log"
	tele "gopkg.in/telebot.v3"
	"github.com/zsefvlol/timezonemapper"
	// "github.com/google/uuid"
	"gorm.io/gorm"
)

func StartHandler(c tele.Context) error {
	var (
		user = c.Sender()
		db = c.Get("db").(*gorm.DB)
	)

	tz := TimeZone{}
	if db.Model(&tz).Where("user_id = ?", user.ID).Find(&tz).RowsAffected != 0 {
		return c.Send(fmt.Sprintf("Твой часовой пояс: %s", tz.TZ))
	}

	r := &tele.ReplyMarkup{ResizeKeyboard: true}

	r.Reply(r.Row(r.Location("Отправить местоположение")))
	return c.Send(fmt.Sprintf("Шалом, %s!\n", user.Username) +
		"Отправь мне свое местоположение, чтобы я смог установить правильный часовой пояс", 
	r)
}

func TimeZoneHandler(c tele.Context) error {
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
		log.Printf(fmt.Sprintf("Updated info about timezone for %s(%d)", user.Username, user.ID))
	}

	return c.Send(fmt.Sprintf("Часовой пояс установлен в: \n<strong>%s</strong>", timezone), 
		"HTML")	
}
