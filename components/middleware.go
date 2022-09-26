package components

import (
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
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


func TimeZoneSet() func(tele.HandlerFunc) tele.HandlerFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			user := c.Sender()
			db := c.Get("db").(*gorm.DB)

			tz := TimeZone{}
			if db.Model(&tz).Where("user_id = ?", user.ID).Find(&tz).RowsAffected == 0 {
				return StartHandler(c)
			}
			return next(c)
		}
	}
}


