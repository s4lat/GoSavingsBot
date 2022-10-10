package components

import (
	tele "gopkg.in/telebot.v3"
	"golang.org/x/text/language"
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

func SetLang() func(tele.HandlerFunc) tele.HandlerFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			var (
				user_id = c.Sender().ID
				args = c.Args()
				db = c.Get("db").(*gorm.DB)
			)

			if (len(args) > 1 && args[1] == "set_lang") {
				lang, err := language.Parse(args[2])

				if err != nil{
					return LangAskHandler(c)
				}

				c.Set("lang", &lang)
				return next(c)
			}

			user := User{}
			if db.Find(&user, "id = ?", user_id).RowsAffected == 0 {
				return LangAskHandler(c)
			}

			if len(user.Lang) == 0 {
				return LangAskHandler(c)
			}

			lang, err := language.Parse(user.Lang)
			if err != nil{
				return LangAskHandler(c)
			}


			c.Set("lang", &lang)
			return next(c)
		}
	}
}

func SetLocation() func(tele.HandlerFunc) tele.HandlerFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			var (
				user_id = c.Sender().ID
				args = c.Args()
				db = c.Get("db").(*gorm.DB)
			)

			if (len(args) > 1 && args[1] == "set_lang") {
				return next(c)
			}

			user := User{}
			if db.Find(&user, "id = ?", user_id).RowsAffected == 0 {
				return TimeZoneAskHandler(c)
			}

			if (len(user.TimeZone) == 0) {
				return TimeZoneAskHandler(c)
			}

			loc, err := time.LoadLocation(user.TimeZone)
			if err != nil {
				log.Print(err)
				return TimeZoneAskHandler(c)
			}

			c.Set("loc", loc)
			return next(c)
		}
	}
}


