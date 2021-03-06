package app

import (
	config "gitlab-tg-bot/conf"
	"gitlab-tg-bot/data"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/message-handling/langs"
	"gitlab-tg-bot/service"
	"gitlab-tg-bot/transport"
	"gitlab-tg-bot/worker"

	"github.com/sirupsen/logrus"

	"github.com/go-pg/pg/v9"
	_ "github.com/go-pg/pg/v9/orm"
)

type App struct {
	ProviderStorage *data.ProviderStorage
	ServiceStorage  interfaces.ServiceStorage
	Conf            interfaces.Configuration

	tg interfaces.TelegramWorker
}

func NewApp(conf interfaces.Configuration) App {
	db := pg.Connect(&pg.Options{
		Addr:     conf.GetString(config.DbConnectionString),
		User:     conf.GetString(config.DbUser),
		Password: conf.GetString(config.DbPassword),
		Database: conf.GetString(config.DbName),
	})

	providerStorage := data.NewProviderStorage(db)

	langs.Init(conf.GetString(config.DefaultLanguage), nil)

	app := App{
		ProviderStorage: providerStorage,
		ServiceStorage:  service.NewStorage(providerStorage, conf),
		Conf:            conf,
	}

	go app.printConfig()

	return app
}

func (a *App) Start() {
	a.tg = worker.NewTelegramWorker(a.Conf, a.ServiceStorage)
	go a.tg.Start()

	go transport.ServeHTTP(a.Conf, a.ServiceStorage, a.tg)
}

func (a *App) Stop() {
	a.tg.Stop()
}

func (a *App) printConfig() {
	logrus.Infof("Сервер запущен на порте %s", a.Conf.GetString(config.ServerUrl))
	logrus.Infof("Секретный ключ для хуков %s", a.Conf.GetString(config.SecretKey))
	logrus.Infof("Ссылка на перехватчик хуков %s", a.Conf.GetString(config.WebHookUrl))
	logrus.Infof("Подключение к базе. Connection string: %s, Username: %s, Password: %s, DbName: %s",
		a.Conf.GetString(config.DbConnectionString),
		a.Conf.GetString(config.DbUser),
		a.Conf.GetString(config.DbPassword),
		a.Conf.GetString(config.DbName))
}
