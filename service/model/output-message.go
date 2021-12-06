package model

type MessageOptions struct {
	DisableNotification bool
}

// OutputMessage
// Lang - должен быть в формате ru_RU
type OutputMessage struct {
	ChatId int64
	MessageOptions
	Msg  string
	Lang string
}
