package main

import (
	"fmt"
	"gitlab-tg-bot/internal"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/go-pg/migrations/v7"
	"github.com/go-pg/pg/v9"
	"github.com/spf13/pflag"
)

func main() {
	db := Connect()
	var err error

	oldVersion, newVersion, err := migrations.Run(db, pflag.Args()...)
	if err != nil {
		if len(err.Error()) > 38 && err.Error()[0:38] == "table \"gopg_migrations\" does not exist" {
			args := []string{"init"}
			_, _, err = migrations.Run(db, args...)
			if err != nil {
			}
			fmt.Print("Initialized! Version is 0.\n")
		} else {
		}
		oldVersion, newVersion, err = migrations.Run(db, pflag.Args()...)
	}
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(oldVersion, " -> ", newVersion)
}

func Connect() *pg.DB {
	conf, err := internal.NewConfiguration()
	if err != nil {

	}
	loginOption := &pg.Options{
		Addr:     conf.GetString(internal.DbConnectionString),
		User:     conf.GetString(internal.DbUser),
		Password: conf.GetString(internal.DbPassword),
		Database: conf.GetString(internal.DbName),
	}
	logrus.Infof("Подключение к базе данных. Address: %s, User: %s, Password: %s, DbName: %s",
		loginOption.Addr, loginOption.User, loginOption.Password, loginOption.Database)
	return pg.Connect(loginOption)
}
