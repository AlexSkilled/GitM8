package processors

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/model"
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

func (s *Start) Process(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI) (isEnd bool) {
	err := s.services.User().Register(model.User{
		Id:         update.Message.From.ID,
		Name:       update.Message.From.FirstName,
		TgUsername: update.Message.From.UserName,
	})
	var messageText string
	if err != nil {
		messageText = "Ошибка при регистрации."
	} else {
		messageText = "Регистрация прошла успешно."
	}

	_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, messageText))

	return true
}
