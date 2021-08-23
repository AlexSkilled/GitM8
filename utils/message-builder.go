package utils

import (
	"errors"
	"fmt"
	"strings"
)

type MessageBuilder struct {
	sb     strings.Builder
	values []interface{}
}

func (m *MessageBuilder) WriteStringN(text string) {
	m.sb.WriteString("\n" + text)
}

func (m *MessageBuilder) WriteStringNF(text string, value ...interface{}) {
	m.sb.WriteString("\n" + text)
	m.values = append(m.values, value...)
}

func (m *MessageBuilder) String() (string, error) {
	message := m.sb.String()

	if strings.Count(message, "%s") != len(m.values) {
		return "", errors.New("количество вхождений %s не равно количеству переданных аргументов")
	}

	message = fmt.Sprintf(message, m.values...)

	return message, nil
}
