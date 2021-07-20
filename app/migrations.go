package app

import (
	"gitlab-tg-bot/internal"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/go-pg/migrations/v7"
	"github.com/go-pg/pg/v9"
	"github.com/sirupsen/logrus"
)

func CheckAndMigrate(conf internal.Configuration) {
	db := Connect(conf)

	oldVersion, newVersion, targetVersion, err :=
		int64(0), int64(0), int64(0), error(nil)

	targetVersion, err = getVersion()

	newVersion, err = migrations.Version(db)
	if err != nil {
		errrorr := err.Error()
		if strings.HasPrefix(errrorr, "ERROR #42P01 relation \"gopg_migrations\" does not exist") {
			oldVersion, newVersion, err = migrations.Run(db, []string{"init"}...) // Если миграция ещё не проводилась, инициируем и ставим нужную версию
			if err != nil {
				panic(err)
			}
		}
	}
	if newVersion < targetVersion {
		oldVersion, newVersion, err = migrations.Run(db, []string{}...) // если установлена не та версия бд, обновляем вверх
		if err != nil {
			panic(err)
		}
		logrus.Info(oldVersion, " -> ", newVersion)
		return
	}
	if newVersion > targetVersion {
		panic("Версия миграции существующей базы больше доступной. Откатите базу")
	}

	panic(err)
}

func Connect(conf internal.Configuration) *pg.DB {
	return pg.Connect(&pg.Options{
		Addr:     conf.GetString(internal.DbConnectionString), // "localhost:1000"
		User:     conf.GetString(internal.DbUser),             //"gitlab_bot"
		Password: conf.GetString(internal.DbPassword),         //"9_9"
		Database: conf.GetString(internal.DbName),             // "gitlab_bot"
	})
	// docker run -d --name tg-gitlab-integration -e POSTGRES_PASSWORD=9_9 -e  POSTGRES_USER=gitlab_bot -e POSTGRES_DB=gitlab_bot --restart always -p "1000:5432" postgres
}

func getVersion() (int64, error) {
	files, err := ioutil.ReadDir("./migrations")
	if err != nil {
		return 0, err
	}
	var lastVer int64
	for _, f := range files {
		name := f.Name()
		dashIdx := strings.IndexRune(name, '_')
		if dashIdx != -1 {
			v := name[:dashIdx]
			ver, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			if lastVer < int64(ver) {
				lastVer = int64(ver)
			}
		}
	}
	return lastVer, nil
}
