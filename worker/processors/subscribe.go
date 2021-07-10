package processors

import (
	"context"
	"fmt"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/model"
	"gitlab-tg-bot/utils"
	"strconv"

	"github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var _ interfaces.Interceptor = (*SubscribeProcessor)(nil)

type SubscriptionStep int32

const (
	SubscriptionStepDomain SubscriptionStep = iota
	SubscriptionStepRepository
	SubscriptionStepEnd
)

type SubscribeProcessor struct {
	service        interfaces.ServiceStorage
	subscribeForms map[int64]*subscribeForm
}

type subscribeForm struct {
	domain       string
	repositoryId int64
	currentStep  SubscriptionStep
}

func NewSubscribeProcessor(services interfaces.ServiceStorage) *SubscribeProcessor {
	return &SubscribeProcessor{
		service:        services,
		subscribeForms: map[int64]*subscribeForm{},
	}
}

func (s *SubscribeProcessor) Process(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI) (isEnd bool) {
	user, err := utils.ExtractUser(ctx)
	if err != nil {
		// TODO
		logrus.Errorln(err)
		return true
	}

	form, ok := s.subscribeForms[user.Id]
	if !ok {
		s.subscribeForms[user.Id] = &subscribeForm{
			currentStep: SubscriptionStepDomain,
		}
		form, _ = s.subscribeForms[user.Id]
	}

	if update.CallbackQuery != nil {
		switch form.currentStep {
		case SubscriptionStepDomain:
			form.domain = update.CallbackQuery.Data

		case SubscriptionStepRepository:
			form.repositoryId, err = strconv.ParseInt(update.CallbackQuery.Data, 10, 32)
		}

		form.currentStep++
	}

	switch form.currentStep {
	case SubscriptionStepDomain:
		return s.getOrSuggestDomains(user, form, update, bot, ctx)
	case SubscriptionStepRepository:
		var gitlab model.GitlabUser
		for _, item := range user.Gitlabs {
			if item.Domain == form.domain {
				gitlab = item
				break
			}
		}
		return s.getOrSuggestRepositories(gitlab, form, update, bot, ctx)
	case SubscriptionStepEnd:
	}

	return true
}

func (s *SubscribeProcessor) IsInterceptor() bool {
	return true
}

func (s *SubscribeProcessor) DumpChatSession(chatId int64) {
	delete(s.subscribeForms, chatId)
}

func (s *SubscribeProcessor) getOrSuggestDomains(user model.User, form *subscribeForm, update tgbotapi.Update, bot *tgbotapi.BotAPI, ctx context.Context) bool {
	switch len(user.Gitlabs) {
	case 0:
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Необходимо зарегестрировать аккаунт gitlab - команда /register"))
		return true
	case 1:
		form.domain = user.Gitlabs[0].Domain
		form.currentStep++
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Зарегистрирован только один домен - %s. Этап пропускается", user.Gitlabs[0].Domain)))
		return s.Process(ctx, update, bot)
	default:
		buttons := make([]tgbotapi.InlineKeyboardButton, len(user.Gitlabs))
		for _, item := range user.Gitlabs {
			buttons = append(buttons,
				tgbotapi.NewInlineKeyboardButtonData("Аккаунт: "+item.Username+". Домен: "+item.Domain, item.Domain))
		}

		message := utils.NewTgMessageWithButtons(update.Message.Chat.ID, "Выберите домен", buttons)

		_, err := bot.Send(message)
		if err != nil {
			logrus.Errorln(err)
		}
	}
	return false
}

func (s *SubscribeProcessor) getOrSuggestRepositories(user model.GitlabUser, form *subscribeForm, update tgbotapi.Update, bot *tgbotapi.BotAPI, ctx context.Context) (isEnd bool) {
	repos, err := s.service.GitlabApi().GetRepositories(user)
	if err != nil {
		logrus.Errorln(err)
	}
	switch len(repos) {
	case 0:
		_, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Для домена %s отсутствуют репозитории", user.Domain)))
		return true
	case 1:
		form.repositoryId = int64(repos[0].Id)
		form.currentStep++
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Найден единственный репозиторий - %s.", repos[0].Name)))
		return s.Process(ctx, update, bot)
	default:
		buttons := make([]tgbotapi.InlineKeyboardButton, len(repos))
		for _, item := range repos {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(item.RepoName, strconv.Itoa(int(item.Id))))
		}
		message := utils.NewTgMessageWithButtons(update.Message.Chat.ID, "Выберите репозиторий", buttons)
		_, err = bot.Send(message)
		if err != nil {
			logrus.Error(err)
		}
	}

	return false
}
