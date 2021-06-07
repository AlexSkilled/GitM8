package moq

import (
	"github.com/stretchr/testify/mock"
)

type Config struct {
	mock.Mock
}

func (c Config) GetBool(s string) bool {
	args := c.Called(s)
	return args.Bool(0)
}

func (c Config) GetInt(s string) int {
	args := c.Called(s)
	return args.Int(0)
}

func (c Config) GetInt32(s string) int32 {
	args := c.Called(s)
	return int32(args.Int(0))
}

func (c Config) GetInt64(s string) int64 {
	args := c.Called(s)
	return int64(args.Int(0))
}

func (c Config) GetString(s string) string {
	args := c.Called(s)
	return args.String(0)
}
