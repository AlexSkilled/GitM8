package main

import (
	"fmt"
	"os"

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
	return pg.Connect(&pg.Options{
		Addr:     "localhost:1000",
		User:     "gitlab_bot",
		Password: "9_9",
		Database: "gitlab_bot",
	})
	// docker run -d --name tg-gitlab-integration -e POSTGRES_PASSWORD=9_9 -e  POSTGRES_USER=gitlab_bot -e POSTGRES_DB=gitlab_bot -p "1000:5432" postgres
}
