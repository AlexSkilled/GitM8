package processors

import (
	"context"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/model"
	"gitlab-tg-bot/utils"

	"github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var _ interfaces.Interceptor = (*SubscribeProcessor)(nil)

type SubscribeProcessor struct {
	service interfaces.ServiceStorage
}

func NewSubscribeProcessor(services interfaces.ServiceStorage) *SubscribeProcessor {
	return &SubscribeProcessor{
		service: services,
	}
}

func (s *SubscribeProcessor) Process(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI) (isEnd bool) {
	user, err := utils.ExtractUser(ctx)
	if err != nil {
		// TODO
	}
	var gitlab model.GitlabUser

	switch len(user.Gitlabs) {
	case 0:
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Необходимо зарегестрировать аккаунт gitlab - команда /register"))
		return true
	case 1:
		gitlab = user.Gitlabs[0]
	default:
		if update.CallbackQuery != nil {
			gitlab = user.Gitlabs[0]
		} else {
			button := tgbotapi.NewInlineKeyboardButtonData("привет!", "1")

			markup := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{button})

			message := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите репозиторий")
			message.BaseChat.ReplyMarkup = markup
			_, err = bot.Send(message)
			if err != nil {
				logrus.Error(err)
			}

			return false
		}
		// TODO предложить выбрать акк на гите
	}

	_, err = s.service.GitlabApi().GetRepositories(gitlab)

	return true
}

func (s *SubscribeProcessor) IsInterceptor() bool {
	return true
}

func (s *SubscribeProcessor) DumpUserSession(userId int64) {
	panic("implement me")
}
