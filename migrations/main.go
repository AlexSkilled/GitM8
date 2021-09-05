package main

import (
	"fmt"
	config "gitlab-tg-bot/conf"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/go-pg/migrations/v7"
	"github.com/go-pg/pg/v9"
	"github.com/spf13/pflag"
)

func main() {
	logrus.Error("fuck this docker")
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
		if err != nil {

		}
	}
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(oldVersion, " -> ", newVersion)

	wd, _ := os.Getwd()
	wd += "\\migrations\\message-patterns-data\\"
	patterns, err := ioutil.ReadDir(wd)
	if err != nil {
		logrus.Error(err)
		return
	}

	insertionScript := strings.Builder{}

	for _, f := range patterns {
		name := f.Name()
		if strings.Contains(name, ".sql") {
			file, err := ioutil.ReadFile(wd + name)
			if err != nil {
				logrus.Error(err)
				continue
			}

			insertionScript.Write(file)
		}
	}

	_, err = db.Exec(insertionScript.String())
	if err != nil {
		logrus.Error(err)
	}

}

func Connect() *pg.DB {
	conf, err := config.NewConfiguration()
	if err != nil {

	}
	loginOption := &pg.Options{
		Addr:     conf.GetString(config.DbConnectionString),
		User:     conf.GetString(config.DbUser),
		Password: conf.GetString(config.DbPassword),
		Database: conf.GetString(config.DbName),
	}
	logrus.Infof("Подключение к базе данных. Address: %s, User: %s, Password: %s, DbName: %s",
		loginOption.Addr, loginOption.User, loginOption.Password, loginOption.Database)
	return pg.Connect(loginOption)
	// docker run -d --name tg-gitlab-integration -e POSTGRES_PASSWORD=9_9 -e  POSTGRES_USER=gitlab_bot -e POSTGRES_DB=gitlab_bot --restart always -p "1000:5432" postgres
}
