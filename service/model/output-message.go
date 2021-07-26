package model

type MessageOptions struct {
	DisableNotification bool
}

type OutputMessage struct {
	ChatId int64
	MessageOptions
	Msg string
}
