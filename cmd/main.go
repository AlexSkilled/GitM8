package main

import (
	"os"
	"os/signal"

	"gitlab-tg-bot/app"
	configuration "gitlab-tg-bot/conf"
)

func main() {
	conf := configuration.NewConfiguration()

	application := app.NewApp(conf)

	app.CheckMigration(conf)

	application.Start()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	<-done

	application.Stop()
}
