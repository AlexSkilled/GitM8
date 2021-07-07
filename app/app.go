package app

import (
	"github.com/go-pg/pg/v9"
	_ "github.com/go-pg/pg/v9/orm"
	"gitlab-tg-bot/data"
	"gitlab-tg-bot/internal"
	"gitlab-tg-bot/service"
	"gitlab-tg-bot/transport"
	"gitlab-tg-bot/worker"
	"net/http"
)

type App struct {
	ProviderStorage *data.ProviderStorage
	ServiceStorage  service.Storage
	Conf            internal.Configuration
	//Telegram       tgbotapi.BotAPI
}

func NewApp(conf internal.Configuration) App {
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:1000",
		User:     "gitlab_bot",
		Password: "9_9",
		Database: "gitlab_bot",
	})

	providerStorage := data.NewProviderStorage(db)
	app := App{
		ProviderStorage: &providerStorage,
		ServiceStorage:  service.NewStorage(providerStorage),
		Conf:            conf,
	}

	return app
}

func (a *App) Start() {
	tgIntegration := worker.NewTelegramWorker(a.Conf, &a.ServiceStorage)
	go tgIntegration.Start()
	handler := transport.NewPublicHandler(a.Conf, tgIntegration)
	panic(http.ListenAndServe(a.Conf.GetString(internal.ServerUrl), handler))
}
