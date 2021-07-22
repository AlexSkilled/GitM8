package moq

import (
	"fmt"

	"github.com/stretchr/testify/mock"
)

type Notifier struct {
	mock.Mock
}

func (m Notifier) Notify(payload string) error {
	fmt.Println(payload)
	args := m.Called(payload)
	return args.Error(0)
}
