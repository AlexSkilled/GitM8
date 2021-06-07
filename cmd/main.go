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
	notifier, err := internal.NewNotifier(conf)
	if err != nil {
		logrus.Errorln("Ошибка при подключении к Telegram API.")
		panic(err)
	}
	handler := internal.NewHandler(conf, notifier)
	panic(http.ListenAndServe(conf.GetString(internal.ServerUrl), handler))
}
