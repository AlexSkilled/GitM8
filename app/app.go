package app

import (
	"gitlab-tg-bot/data"
	"gitlab-tg-bot/internal"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service"
	"gitlab-tg-bot/transport"
	"gitlab-tg-bot/worker"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/go-pg/pg/v9"
	_ "github.com/go-pg/pg/v9/orm"
)

type App struct {
	ProviderStorage *data.ProviderStorage
	ServiceStorage  interfaces.ServiceStorage
	Conf            internal.Configuration
	//Telegram       tgbotapi.BotAPI
}

func NewApp(conf internal.Configuration) App {
	CheckMigration(conf)

	db := pg.Connect(&pg.Options{
		Addr:     conf.GetString(internal.DbConnectionString),
		User:     conf.GetString(internal.DbUser),
		Password: conf.GetString(internal.DbPassword),
		Database: conf.GetString(internal.DbName),
	})

	providerStorage := data.NewProviderStorage(db)
	app := App{
		ProviderStorage: providerStorage,
		ServiceStorage:  service.NewStorage(providerStorage, conf),
		Conf:            conf,
	}

	go app.printConfig()

	return app
}

func (a *App) Start() {
	tgIntegration := worker.NewTelegramWorker(a.Conf, a.ServiceStorage)
	go tgIntegration.Start()
	handler := transport.NewPublicHandler(a.Conf, tgIntegration)
	panic(http.ListenAndServe(a.Conf.GetString(internal.ServerUrl), handler))
}

func (a *App) printConfig() {
	logrus.Infof("Сервер запущен на порте %s", a.Conf.GetString(internal.ServerUrl))
	logrus.Infof("Секретный ключ для хуков %s", a.Conf.GetString(internal.SecretKey))
	logrus.Infof("Ссылка на перехватчик хуков %s", a.Conf.GetString(internal.WebHookUrl))
	logrus.Infof("Подключение к базе. Connection string: %s, Username: %s, Password: %s, DbName: %s",
		a.Conf.GetString(internal.DbConnectionString),
		a.Conf.GetString(internal.DbUser),
		a.Conf.GetString(internal.DbPassword),
		a.Conf.GetString(internal.DbName))
}
