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

func CheckMigration(conf internal.Configuration) {
	db := Connect(conf)

	newVersion, targetVersion, err :=
		int64(0), int64(0), error(nil)

	targetVersion, err = getVersion()

	newVersion, err = migrations.Version(db)
	if err != nil {
		logrus.Errorln(err)
		panic(err)
	}
	if newVersion != targetVersion {
		panic("Версия базы не совпадает с миграциями в /migrations")
	}
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
