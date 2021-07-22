package main

import (
	"gitlab-tg-bot/app"
	configuration "gitlab-tg-bot/conf"

	"github.com/sirupsen/logrus"
)

func main() {
	conf, err := configuration.NewConfiguration()
	if err != nil {
		logrus.Errorln("Ошибка при конфигурации приложения.")
		panic(err)
	}
	application := app.NewApp(conf)
	application.Start()
}
