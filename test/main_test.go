package test

import (
	"gitlab-tg-bot/test/utils"
	"testing"

	"github.com/go-pg/pg/v9"
)

type TestEnv struct {
	*pg.DB
}

func TestMain(m *testing.M) {
	_, closer := utils.CreateDocker()

	defer closer()

	m.Run()

}
