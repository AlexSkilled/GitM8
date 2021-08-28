package utils

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-pg/pg/v9"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	dsn      = "postgres://%s:%s@localhost%s/%s?sslmode=disable"
	dialect  = "postgres"
	idleConn = 25
	maxConn  = 25

	containerName = "gitlab_tests"
)

type DockerConfig struct {
	User   string
	Pass   string
	DbName string
	Port   string
}

func CreateDocker(conf DockerConfig) (db *pg.DB, clearFunc func(), err error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return
	}

	opts := getOptions(conf)

	resource, ok := pool.ContainerByName(containerName)
	if ok {
		// Если такой контейнер уже существует (тесты были прерваны пользователем) удаляем и создаём на его месте новый
		pool.Purge(resource)
	}

	resource, err = pool.RunWithOptions(&opts)

	clearFunc = func() { pool.Purge(resource) }

	defer func() {
		if err != nil {
			clearFunc()
		}
	}()

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

	db = pg.Connect(&pg.Options{
		Addr:     "localhost" + conf.Port,
		User:     conf.User,
		Password: conf.Pass,
		Database: conf.DbName,
	})

	return
}

func getOptions(conf DockerConfig) dockertest.RunOptions {
	return dockertest.RunOptions{
		Name:       containerName,
		Repository: "postgres",
		Tag:        "alpine",
		Env: []string{
			"POSTGRES_PASSWORD=" + conf.Pass,
			"POSTGRES_USER=" + conf.User,
			"POSTGRES_DB=" + conf.DbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{
					HostIP: "0.0.0.0", HostPort: conf.Port[1:],
				},
			},
		},
	}
}
