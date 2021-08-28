package main

import (
	"gitlab-tg-bot/app"
	configuration "gitlab-tg-bot/conf"

	"github.com/go-pg/pg/v9"
)

func main() {
	conf := configuration.NewConfiguration()

	application := app.NewApp(conf)

	app.CheckMigration(pg.Connect(&pg.Options{
		Addr:     conf.GetString(configuration.DbHost) + conf.GetString(configuration.DbPort),
		User:     conf.GetString(configuration.DbUser),
		Password: conf.GetString(configuration.DbPassword),
		Database: conf.GetString(configuration.DbName),
	}))

	application.Start()
}
