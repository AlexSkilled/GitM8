package test

import (
	"gitlab-tg-bot/app"
	configuration "gitlab-tg-bot/conf"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/test/utils"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/go-pg/pg/v9"
)

type TestEnv struct {
	*pg.DB
	application app.App
	conf        interfaces.Configuration
}

var testEnvironment TestEnv

func TestMain(m *testing.M) {
	dir, _ := os.Getwd()
	dir += "/../conf/bot-conf-test.yml"
	_ = pflag.String("conf-path", dir, "Path to configuration file")
	pflag.Parse()

	testEnvironment.conf = configuration.NewConfiguration()

	conf := utils.DockerConfig{
		User:   testEnvironment.conf.GetString(configuration.DbUser),
		Pass:   testEnvironment.conf.GetString(configuration.DbPassword),
		DbName: testEnvironment.conf.GetString(configuration.DbName),
		Port:   testEnvironment.conf.GetString(configuration.DbPort),
	}

	db, closer, err := utils.CreateDocker(conf)
	defer closer()
	if err != nil {
		logrus.Fatal(err)
	}

	testEnvironment.DB = db
	testEnvironment.migrate()

	testEnvironment.application = app.NewApp(testEnvironment.conf)

	m.Run()
}

func (t *TestEnv) migrate() {
	dir, _ := os.Getwd()
	dir += "/../migrations/"
	logrus.Println("Накатывание миграций из", dir)

	filesMigrations, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	type migrationFile struct {
		Number int64
		Script string
	}

	migrationFiles := make([]migrationFile, 0, len(filesMigrations))
	for _, f := range filesMigrations {
		name := f.Name()
		if strings.Contains(name, ".up.sql") {
			number, err := strconv.ParseInt(name[0:strings.IndexRune(name, '_')], 10, 32)
			if err != nil {
				logrus.Errorf("неверное название файла миграции '%s' \nОшибка:%v", name, err)
				continue
			}

			script, err := ioutil.ReadFile(dir + name)
			if err != nil {
				logrus.Errorf("ошибка при попытке прочитать файл миграции %s \n%v", dir+name, err)
				continue
			}

			migrationFiles = append(migrationFiles,
				migrationFile{
					Number: number,
					Script: string(script),
				})
		}
	}

	sort.Slice(migrationFiles, func(i, j int) bool {
		return migrationFiles[i].Number < migrationFiles[j].Number
	})
	finalScript := ``
	for _, item := range migrationFiles {
		finalScript += item.Script
	}

	_, err = t.Exec(finalScript)
	if err != nil {
		logrus.Fatalf("ошибка при накатке миграций:%v", err)
	}

}
