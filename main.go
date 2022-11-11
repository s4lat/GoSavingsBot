package main

import (
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/s4lat/gosavingsbot/handlers"
	"github.com/s4lat/gosavingsbot/locale"
	"github.com/s4lat/gosavingsbot/log"
	"github.com/s4lat/gosavingsbot/middleware"
	"github.com/s4lat/gosavingsbot/models"
)

func main() {
	log.InitLoggers()
	locale.InitLocales()

	log.InfoLogger.Print("Connecting to db...")
	db, err := gorm.Open(sqlite.Open("./data/data.db"), &gorm.Config{})
	if err != nil {
		log.ErrorLogger.Fatal("Failed to connect database")
	}
	db.AutoMigrate(&models.Spend{}, &models.User{})

	pref := tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.ErrorLogger.Fatal(err)
		return
	}

	b.Use(
		middleware.PassData(map[string]interface{}{"db": db}),
		middleware.SetLang(),
	)

	b.Handle("/help", handlers.StartHandler)
	b.Handle("/set_lang", handlers.LangAskHandler)
	b.Handle("/start", handlers.StartHandler, middleware.SetLocation())
	b.Handle("/delete_my_data", handlers.AskToDeleteUserData, middleware.SetLocation())

	b.Handle("Today", handlers.DaySpendsHandler, middleware.SetLocation())
	b.Handle("Сегодня", handlers.DaySpendsHandler, middleware.SetLocation())

	b.Handle("Statistics", handlers.YearSpendsHandler, middleware.SetLocation())
	b.Handle("Статистика", handlers.YearSpendsHandler, middleware.SetLocation())

	b.Handle("Settings", handlers.SettingsHandler, middleware.SetLocation())
	b.Handle("Настройки", handlers.SettingsHandler, middleware.SetLocation())

	b.Handle(tele.OnLocation, handlers.LocationHandler)
	b.Handle(tele.OnText, handlers.OnTextHandler, middleware.SetLocation())
	b.Handle(tele.OnCallback, handlers.CallbackHandler, middleware.SetLocation())

	log.InfoLogger.Print("Starting bot...")
	b.Start()
}
