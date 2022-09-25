package main

/*
	TODO:
		0. Expenses struct, saving to db:
			name, cost, date
		1. Adding expenses by sending the line "expense - expense name"
		2. Export all expenses to excel table
		3. Ability to change between currencies
*/

import (
	"log"
	"os"
	"time"
	"fmt"

	tele "gopkg.in/telebot.v3"
	// "gorm.io/gorm"
	// "gorm.io/driver/sqlite"
)

func main() {

	pref := tele.Settings{
		Token:  os.Getenv("TG_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(c tele.Context) error {
		user := c.Sender()
		return c.Send(fmt.Sprintf("Hi, %s!", user.Username))
	})

	log.Print("Starting bot...")
	b.Start()
}