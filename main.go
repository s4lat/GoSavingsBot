package main

/*
	TODO:
		1. Get spends by month, day
		month - Function in models.go, getting rows for 33, then converting to user timezone and sorting by month again
		day - getting rows for 3 days, then converting to user timezone and sorting by day again
		2. Export all expenses to excel/csv file
		3. Ability to change between currencies
*/

import (
	"log"
	"os"
	"time"
	// "fmt"
	
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
	db.AutoMigrate(&comps.Spend{}, &comps.TimeZone{})

	pref := tele.Settings{
		Token:  os.Getenv("TG_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Use(comps.PassData(map[string]interface{}{"db": db}))

	b.Handle("/start", comps.HomeHandler, comps.SetLocation())
	b.Handle(tele.OnText, comps.AddSpendHandler, comps.SetLocation())
	b.Handle(tele.OnLocation, comps.LocationHandler)

	log.Print("Starting bot...")
	b.Start()
}