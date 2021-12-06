package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	config "gitlab-tg-bot/conf"

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
		if strings.Contains(err.Error(), "\"gopg_migrations\" does not exist") {
			args := []string{"init"}
			_, _, err = migrations.Run(db, args...)
			if err != nil {
				panic(err)
			}
		}
		oldVersion, newVersion, err = migrations.Run(db, pflag.Args()...)
		if err != nil {
			fmt.Printf("Ошибка при накатке миграции! %s", err.Error())
			os.Exit(1)
		}
	}

	fmt.Println(oldVersion, " -> ", newVersion)

	wd, _ := os.Getwd()
	wd = path.Join(wd, "migrations", "message-patterns-data")
	patterns, err := ioutil.ReadDir(wd)
	if err != nil {
		logrus.Error(err)
		return
	}

	insertionScript := strings.Builder{}

	for _, f := range patterns {
		name := f.Name()
		if strings.HasSuffix(name, ".sql") {
			file, err := ioutil.ReadFile(path.Join(wd, name))
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
	conf := config.NewConfiguration()

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
