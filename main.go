package main

/*
	TODO:
		1. Switching languages
		2. Pic for stats: pillars or circle
		3. Add loging
		4. Comments
		5. Pretty README about:
			1. Functions
			2. Deploy
			3. Hosted @GoSavingsBot
*/

import (
	"log"
	"os"
	"time"
	
	comps "my_projects/GoSavingsBot/components"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func main() {
	log.Print("Connecting to db...")
	db, err := gorm.Open(sqlite.Open("./data/data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	}
	db.AutoMigrate(&comps.Spend{}, &comps.TimeZone{})

	pref := tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Use(comps.PassData(map[string]interface{}{"db": db}))

	b.Handle("/start", comps.StartHandler, comps.SetLocation())
	b.Handle("Сегодня", comps.DaySpendsHandler, comps.SetLocation())
	b.Handle("Статистика", comps.YearSpendsHandler, comps.SetLocation())
	b.Handle(tele.OnText, comps.OnTextHandler, comps.SetLocation())
	b.Handle(tele.OnLocation, comps.LocationHandler)
	b.Handle(tele.OnCallback, comps.CallbackHandler, comps.SetLocation())
	b.Handle("/help", func (c tele.Context) error {
		return c.Send(comps.HELP_MSG, "HTML")
	})

	log.Print("Starting bot...")
	b.Start()
}