package app

import (
	"github.com/go-pg/pg/v9"
	_ "github.com/go-pg/pg/v9/orm"
	"gitlab-tg-bot/data"
)

type App struct {
	ProviderStorage data.ProviderStorage
	ServiceStorage  service.Holder
	//Telegram       tgbotapi.BotAPI
}

func NewApp() App {
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:1000",
		User:     "gitlab_bot",
		Password: "9_9",
		Database: "gitlab_bot",
	})

	providerStorage := data.NewProviderStorage(db)
	app := App{
		ProviderStorage: providerStorage,
		ServiceStorage:  service.NewStorage(providerStorage),
	}

	return app
}
