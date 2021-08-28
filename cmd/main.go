package main

import (
	"gitlab-tg-bot/app"
	configuration "gitlab-tg-bot/conf"
)

func main() {
	conf := configuration.NewConfiguration()

	application := app.NewApp(conf)

	app.CheckMigration(conf)

	application.Start()
}
