package main

import (
	"gitlab-tg-bot/internal"
	"net/http"

	"github.com/sirupsen/logrus"
)

func main() {
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
