package main

/*
	TODO:
		2. Export all expenses to excel/csv file
		3. commands - download data for any month from list, for year, for previous years
		interface example:
		(September: 40000.00 (/csvM09))
		Total: 80000.00 (/csvY2022)
		/excelY

		Excel only for year stats
		4. Ability to change between currencies
		5. Add loging
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