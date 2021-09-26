package utils

import (
	"fmt"
	"strings"
)

type NoticeMessage struct {
	Header   MessageRow
	MainInfo MessageRow
	Author   MessageRow

	TriggeredByUserName *MessageRow
	AssignedToUserName  *MessageRow
	SubInfo             *MessageRow
}

func (n *NoticeMessage) Print() string {
	sb := strings.Builder{}
	sb.WriteString(n.Header.Print() + "\n")
	sb.WriteString(n.MainInfo.Print() + "\n")
	sb.WriteString(n.Author.Print() + "\n")
	switch {
	case n.TriggeredByUserName != nil:
		sb.WriteString(n.TriggeredByUserName.Print() + "\n")
	case n.AssignedToUserName != nil:
		sb.WriteString(n.AssignedToUserName.Print() + "\n")
	case n.SubInfo != nil:
		sb.WriteString(n.SubInfo.Print() + "\n")
	}
	return sb.String()
}

func (n *NoticeMessage) PrintWithEmoji() string {
	sb := strings.Builder{}
	sb.WriteString(n.Header.PrintWithEmoji() + "\n")
	sb.WriteString(n.MainInfo.PrintWithEmoji() + "\n")
	sb.WriteString(n.Author.PrintWithEmoji() + "\n")

	switch {
	case n.TriggeredByUserName != nil:
		sb.WriteString(n.TriggeredByUserName.PrintWithEmoji() + "\n")
	case n.AssignedToUserName != nil:
		sb.WriteString(n.AssignedToUserName.PrintWithEmoji() + "\n")
	case n.SubInfo != nil:
		sb.WriteString(n.SubInfo.PrintWithEmoji() + "\n")
	}
	return sb.String()
}

type MessageRow struct {
	Message string
	Emoji   *string
}

func NewMessageRawWithEmoji(emoji string, pattern string, values ...interface{}) *MessageRow {
	return &MessageRow{
		Message: fmt.Sprintf(pattern, values...),
		Emoji:   &emoji,
	}
}

func NewMessageRaw(pattern string, values ...interface{}) *MessageRow {
	return &MessageRow{
		Message: fmt.Sprintf(pattern, values...),
	}
}

func (m *MessageRow) WriteWithEmoji(emoji string, pattern string, values ...interface{}) {
	m.Emoji = &emoji
	m.Write(pattern, values...)
}

func (m *MessageRow) Write(pattern string, values ...interface{}) {
	m.Message = fmt.Sprintf(pattern, values...)
}

func (m *MessageRow) Print() string {
	return m.Message
}

func (m *MessageRow) PrintWithEmoji() string {
	if m.Emoji != nil {
		return *m.Emoji + m.Message
	}
	return m.Message
}
