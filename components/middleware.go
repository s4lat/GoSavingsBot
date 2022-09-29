package components

import (
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
	"time"
	"log"
)

func PassData(data map[string]interface{}) func(tele.HandlerFunc) tele.HandlerFunc {
	return func (next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			for k, v := range data{
				c.Set(k, v)
			}
			return next(c)
		}
	}
}


func SetLocation() func(tele.HandlerFunc) tele.HandlerFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			user := c.Sender()
			db := c.Get("db").(*gorm.DB)

			tz := TimeZone{}
			if db.Find(&tz, "user_id = ?", user.ID).RowsAffected == 0 {
				return TimeZoneAskHandler(c)
			}

			loc, err := time.LoadLocation(tz.TZ)
			if err != nil {
				log.Print(err)
				return TimeZoneAskHandler(c)
			}

			c.Set("loc", loc)
			c.Set("tz_name", tz.TZ)
			return next(c)
		}
	}
}


