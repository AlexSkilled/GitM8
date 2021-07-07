package interfaces

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type PublicNotifier interface {
	Notify(chatIds []int32, payload string) error
}

type TgProcessor interface {
	IsInterceptor() bool
	Process(update tgbotapi.Update, bot *tgbotapi.BotAPI) (isEnd bool)
}

type PublicProcessor interface {
	Process(payload []byte) (string, error)
}

type StepName int32
