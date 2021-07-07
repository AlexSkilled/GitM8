package internal

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Processor interface {
	Process(payload []byte) (msg string, skip bool, err error)
}

type Notifier interface {
	Notify(payload string) error
}

type UpdateWatcher interface {
	Polling()
}

type Configuration interface {
	GetBool(string) bool
	GetInt(string) int
	GetInt32(string) int32
	GetInt64(string) int64
	GetString(string) string
}

type PublicNotifier interface {
	Notify(chatIds []int32, payload string) error
}

type TgProcessor interface {
	IsInterceptor() bool
	Process(update tgbotapi.Update) bool
}

type PublicProcessor interface {
	Process(payload []byte) (string, error)
}
