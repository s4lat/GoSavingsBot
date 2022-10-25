package components

import (
	"log"
	"time"

	"golang.org/x/text/language"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

// PassData - passing all (key, val) from data map to context by c.Set(k, v).
func PassData(data map[string]interface{}) func(tele.HandlerFunc) tele.HandlerFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			for k, v := range data {
				c.Set(k, v)
			}
			return next(c)
		}
	}
}

// SetLang - Checks if user language is set:
//
//	If context have 'setLang' in args, passing context to LangAskHandler.
//	If user.Lang not set sends context to LangAskHandler.
//	If user.Lang is set, execute c.Set("lang", &lang), where &lang is ptr to language.Tag.
func SetLang() func(tele.HandlerFunc) tele.HandlerFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			var (
				userID = c.Sender().ID
				args   = c.Args()
				db     = c.Get("db").(*gorm.DB)
			)

			if len(args) > 1 && args[1] == "setLang" {
				lang, err := language.Parse(args[2])

				if err != nil {
					return LangAskHandler(c)
				}

				c.Set("lang", &lang)
				return next(c)
			}

			user := User{}
			if db.Find(&user, "id = ?", userID).RowsAffected == 0 {
				return LangAskHandler(c)
			}

			if len(user.Lang) == 0 {
				return LangAskHandler(c)
			}

			lang, err := language.Parse(user.Lang)
			if err != nil {
				return LangAskHandler(c)
			}

			c.Set("lang", &lang)
			return next(c)
		}
	}
}

// SetLocation - checks if user TimeZone is set:
//
//	If context have 'setLang' in args, passing context to 'next'.
//	If user.TimeZone not set, sends context to TimeZoneAskHandler.
//	If user.TimeZone is set, execute c.Set("loc", loc), where loc is ptr to time.Location.
func SetLocation() func(tele.HandlerFunc) tele.HandlerFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			var (
				userID = c.Sender().ID
				args   = c.Args()
				db     = c.Get("db").(*gorm.DB)
			)

			if len(args) > 1 && args[1] == "setLang" {
				return next(c)
			}

			user := User{}
			if db.Find(&user, "id = ?", userID).RowsAffected == 0 {
				return TimeZoneAskHandler(c)
			}

			if len(user.TimeZone) == 0 {
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
