package utils

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-pg/pg/v9"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

const (
	dsn      = "postgres://%s:%s@localhost:%s/%s?sslmode=disable"
	dialect  = "postgres"
	idleConn = 25
	maxConn  = 25
)

type Config struct {
	User   string
	Pass   string
	DbName string
	Port   string
}

func CreateDocker(conf Config) (db *pg.DB, clearFunc func(), err error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return
	}

	opts := getOptions(conf)

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		return
	}

	connectionString := fmt.Sprintf(dsn, conf.User, conf.Pass, conf.Port, conf.DbName)

	if err = pool.Retry(func() error {
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err.Error())
	}

	clearFunc = func() { pool.Purge(resource) }

	db = pg.Connect(&pg.Options{
		Addr:     "localhost:" + conf.Port,
		User:     conf.User,
		Password: conf.Pass,
		Database: conf.DbName,
	})

	return
}

func getOptions(conf Config) dockertest.RunOptions {
	return dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "alpine",
		Env: []string{
			"POSTGRES_PASSWORD=" + conf.Pass,
			"POSTGRES_USER=" + conf.User,
			"POSTGRES_DB=" + conf.DbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"993": {
				{
					HostIP: "0.0.0.0", HostPort: conf.Port,
				},
			},
		},
	}
}
