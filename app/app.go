package app

import (
	"gitlab-tg-bot/data"
	"gitlab-tg-bot/internal"
	"gitlab-tg-bot/service"
	"gitlab-tg-bot/transport"
	"gitlab-tg-bot/worker"
	"net/http"

	"github.com/go-pg/pg/v9"
	_ "github.com/go-pg/pg/v9/orm"
)

type App struct {
	ProviderStorage *data.ProviderStorage
	ServiceStorage  service.Storage
	Conf            internal.Configuration
	//Telegram       tgbotapi.BotAPI
}

func NewApp(conf internal.Configuration) App {
	db := pg.Connect(&pg.Options{
		Addr:     conf.GetString(internal.DbConnectionString),
		User:     conf.GetString(internal.DbUser),
		Password: conf.GetString(internal.DbPassword),
		Database: conf.GetString(internal.DbName),
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
