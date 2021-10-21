package interfaces

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PublicNotifier interface {
	Notify(chatIds []int32, payload string) error
}

type TgProcessor interface {
	IsInterceptor() bool
	Handle(ctx context.Context, message *tgbotapi.Message, bot *tgbotapi.BotAPI) (isEnd bool)
}

type PublicProcessor interface {
	Process(payload []byte) (string, error)
}

type StepName int32

type Interceptor interface {
	TgProcessor
	DumpChatSession(userId int64)
}
