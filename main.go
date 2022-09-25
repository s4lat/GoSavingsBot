package main

/*
	TODO:
		0. Time zone detection and saving
		1. Adding expenses by sending the line "expense - expense name"
		2. Export all expenses to excel/csv file
		3. Ability to change between currencies
*/

import (
	"log"
	"os"
	"time"
	"fmt"
	
	comps "my_projects/GoSavingsBot/components"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func main() {
	log.Print("Connecting to db...")
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	}
	db.AutoMigrate(&comps.Spend{})

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

		spend := comps.Spend{UserID: user.ID, Name: "Test spend", Value: 13.37}
		db.Create(&spend)
		return c.Send(fmt.Sprintf("Hi, %s!", user.Username))
	})

	log.Print("Starting bot...")
	b.Start()
}