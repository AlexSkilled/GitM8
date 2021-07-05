package app

import (
	"github.com/go-pg/pg/v9"
	_ "github.com/go-pg/pg/v9/orm"
	"github.com/sirupsen/logrus"
	"gitlab-tg-bot/data"
	"gitlab-tg-bot/internal"
	"gitlab-tg-bot/service"
	"net/http"
)

type App struct {
	ProviderStorage data.ProviderStorage
	ServiceStorage  service.Storage
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

func (a *App) Start() {
	conf, err := internal.NewConfiguration()
	if err != nil {
		logrus.Errorln("Ошибка при конфигурации приложения.")
		panic(err)
	}

	notifier, err := internal.NewBot(conf)
	if err != nil {
		logrus.Errorln("Ошибка при подключении к Telegram API.")
		panic(err)
	}

	handler := internal.NewHandler(conf, notifier)
	panic(http.ListenAndServeTLS(conf.GetString(internal.ServerUrl), conf.GetString(internal.ServerCertPath),
		conf.GetString(internal.ServerKeyPath), handler))
}
