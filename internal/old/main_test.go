package old

import (
	"gitlab-tg-bot/internal/old/model"
	"reflect"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
)

func TestMain(m *testing.M) {
	faker.AddProvider("datetime", func(v reflect.Value) (interface{}, error) {
		return model.Datetime(time.Now()), nil
	})
	m.Run()
}
