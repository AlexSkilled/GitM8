package processors

import (
	"context"
	"gitlab-tg-bot/internal/interfaces"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Start struct {
	services interfaces.ServiceStorage
}

var _ interfaces.TgProcessor = (*Start)(nil)

func NewStartProcessor(services interfaces.ServiceStorage) *Start {
	return &Start{
		services: services,
	}
}

func (s *Start) IsInterceptor() bool {
	return false
}

func (s *Start) Process(ctx context.Context, message *tgbotapi.Message, bot *tgbotapi.BotAPI) (isEnd bool) {
	messageText := "Новый профиль телеграм добавлен."
	_, _ = bot.Send(tgbotapi.NewMessage(message.Chat.ID, messageText))

	return true
}
