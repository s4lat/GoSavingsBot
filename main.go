package main

import (
	comps "github.com/s4lat/GoSavingsBot/components"
	tele "gopkg.in/telebot.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"time"
)

func main() {
	comps.InitLoggers()
	comps.InitLocales()

	comps.InfoLogger.Print("Connecting to db...")
	db, err := gorm.Open(sqlite.Open("./data/data.db"), &gorm.Config{})
	if err != nil {
		comps.ErrorLogger.Fatal("Failed to connect database")
	}
	db.AutoMigrate(&comps.Spend{}, &comps.User{})

	pref := tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		comps.ErrorLogger.Fatal(err)
		return
	}

	b.Use(
		comps.PassData(map[string]interface{}{"db": db}),
		comps.SetLang(),
	)

	b.Handle("/start", comps.StartHandler, comps.SetLocation())
	b.Handle("/help", comps.StartHandler)
	b.Handle("/set_lang", comps.LangAskHandler)
	b.Handle("/delete_my_data", comps.AskToDeleteUserData, comps.SetLocation())

	b.Handle("Today", comps.DaySpendsHandler, comps.SetLocation())
	b.Handle("Сегодня", comps.DaySpendsHandler, comps.SetLocation())

	b.Handle("Statistics", comps.YearSpendsHandler, comps.SetLocation())
	b.Handle("Статистика", comps.YearSpendsHandler, comps.SetLocation())

	b.Handle("Settings", comps.SettingsHandler, comps.SetLocation())
	b.Handle("Настройки", comps.SettingsHandler, comps.SetLocation())

	b.Handle(tele.OnText, comps.OnTextHandler, comps.SetLocation())
	b.Handle(tele.OnLocation, comps.LocationHandler)
	b.Handle(tele.OnCallback, comps.CallbackHandler, comps.SetLocation())

	comps.InfoLogger.Print("Starting bot...")
	b.Start()
}
