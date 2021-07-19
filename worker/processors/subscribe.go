package processors

import (
	"context"
	"fmt"
	"gitlab-tg-bot/internal"
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
	SubscriptionStepType
	SubscriptionStepEnd
)

type HookType int

const (
	PushEvents HookType = iota
	IssuesEvents
	ConfidentialIssuesEvents
	MergeRequestsEvents
	TagPushEvents
	NoteEvents
	JobEvents
	PipelineEvents
	WikiPageEvents
	EndChoosingEvent
)

var eventsNames = []string{
	"Пуши",
	"Вопросы",
	"Конфиденциальные вопросы",
	"Запрос на слияние",
	"Выпуск тэга",
	"Заметки",
	"Работы (jobы)",
	"Пайп",
	"Вики",
	"Закончить выбор",
}

type SubscribeProcessor struct {
	service        interfaces.ServiceStorage
	subscribeForms map[int64]*subscribeForm
}

type subscribeForm struct {
	domain       string
	repositoryId string
	gitlab       model.GitlabUser
	currentStep  SubscriptionStep
	hookTypes    map[HookType]bool
	lastMessage  tgbotapi.Message
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
			hookTypes:   map[HookType]bool{},
		}
		form, _ = s.subscribeForms[user.Id]
	}

	if update.CallbackQuery != nil {
		switch form.currentStep {
		case SubscriptionStepDomain:
			form.domain = update.CallbackQuery.Data
			s.updateDomainMessage(user, update.CallbackQuery.Message.Chat.ID, form, bot)
		case SubscriptionStepRepository:
			form.repositoryId = update.CallbackQuery.Data
			s.updateRepositoryMessage(update.CallbackQuery.Message.Chat.ID, form, bot)
		case SubscriptionStepType:
			typeNumber, err := strconv.ParseInt(update.CallbackQuery.Data, 10, 64)
			if err != nil {
				logrus.Errorln("Ошибка парсинга ответа!", err)
			}

			if typeNumber != int64(EndChoosingEvent) {
				form.hookTypes[HookType(typeNumber)] = !form.hookTypes[HookType(typeNumber)]
				return s.updateHookTypeMessage(update.CallbackQuery.Message.Chat.ID, form, bot)
			}
		}

		form.currentStep++
	}

	switch form.currentStep {
	case SubscriptionStepDomain:
		return s.suggestDomains(user, form, update, bot, ctx)
	case SubscriptionStepRepository:
		for _, item := range user.Gitlabs {
			if item.Domain == form.domain {
				form.gitlab = item
				break
			}
		}
		return s.suggestRepositories(form, update, bot, ctx)
	case SubscriptionStepType:
		return s.suggestHookType(update, form, bot)
	case SubscriptionStepEnd:
		webhook := model.Hook{
			RepoId:                   form.repositoryId,
			PushEvents:               form.hookTypes[PushEvents],
			IssuesEvents:             form.hookTypes[IssuesEvents],
			ConfidentialIssuesEvents: form.hookTypes[ConfidentialIssuesEvents],
			MergeRequestsEvents:      form.hookTypes[MergeRequestsEvents],
			TagPushEvents:            form.hookTypes[TagPushEvents],
			NoteEvents:               form.hookTypes[NoteEvents],
			JobEvents:                form.hookTypes[JobEvents],
			PipelineEvents:           form.hookTypes[PipelineEvents],
			WikiPageEvents:           form.hookTypes[WikiPageEvents],
		}

		err = s.service.GitlabApi().AddWebHook(form.gitlab, webhook)
		if err != nil {
			logrus.Errorln(err)
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
				"Не удалось создать слушатель эвентов"+err.Error()))
			return true
		}
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Успешно создан слушатель эвентов"))
	}

	return true
}

func (s *SubscribeProcessor) IsInterceptor() bool {
	return true
}

func (s *SubscribeProcessor) DumpChatSession(chatId int64) {
	delete(s.subscribeForms, chatId)
}

func (s *SubscribeProcessor) suggestDomains(user model.User, form *subscribeForm, update tgbotapi.Update,
	bot *tgbotapi.BotAPI, ctx context.Context) bool {
	switch len(user.Gitlabs) {
	case 0:
		_, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Необходимо зарегестрировать аккаунт gitlab - команда /register"))
		if err != nil {
			logrus.Errorln(err)
		}
		return true
	case 1:
		form.domain = user.Gitlabs[0].Domain
		form.currentStep++
		_, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Зарегистрирован только один домен - %s. Этап пропускается", user.Gitlabs[0].Domain)))
		if err != nil {
			logrus.Errorln(err)
			return true
		}
		return s.Process(ctx, update, bot)
	default:
		buttons := make([]tgbotapi.InlineKeyboardButton, len(user.Gitlabs))
		for _, item := range user.Gitlabs {
			buttons = append(buttons,
				tgbotapi.NewInlineKeyboardButtonData("Аккаунт: "+item.Username+". Домен: "+item.Domain, item.Domain))
		}

		markup := utils.NewTgMessageButtonsMarkup(buttons, 2)
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите домен")
		message.BaseChat.ReplyMarkup = markup

		msg, err := bot.Send(message)
		if err != nil {
			logrus.Errorln(err)
			return true
		}
		form.lastMessage = msg
	}
	return false
}
func (s *SubscribeProcessor) updateDomainMessage(user model.User, chatId int64, form *subscribeForm,
	bot *tgbotapi.BotAPI) {

	buttons := make([]tgbotapi.InlineKeyboardButton, len(user.Gitlabs))
	for i, item := range user.Gitlabs {
		buttons[i] = tgbotapi.InlineKeyboardButton{Text: item.Domain}
		buttons[i].CallbackData = nil

		if item.Domain == form.domain {
			buttons[i].Text += internal.GetEmoji(internal.WhiteCheckMark)
		}
	}

	markup := utils.NewTgMessageButtonsMarkup(buttons, 2)

	message := tgbotapi.NewEditMessageReplyMarkup(chatId, form.lastMessage.MessageID, markup)

	_, err := bot.Send(message)
	if err != nil {
		logrus.Error(err)
	}
}

func (s *SubscribeProcessor) suggestRepositories(form *subscribeForm, update tgbotapi.Update,
	bot *tgbotapi.BotAPI, ctx context.Context) (isEnd bool) {
	repos, err := s.service.GitlabApi().GetRepositories(form.gitlab)
	if err != nil {
		logrus.Errorln(err)
	}
	switch len(repos) {
	case 0:
		_, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Для домена %s отсутствуют репозитории", form.gitlab.Domain)))
		return true
	case 1:
		form.repositoryId = repos[0].Id
		form.currentStep++
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Найден единственный репозиторий - %s.", repos[0].Name)))
		return s.Process(ctx, update, bot)
	default:
		buttons := make([]tgbotapi.InlineKeyboardButton, len(repos))
		for i, item := range repos {
			buttons[i] = tgbotapi.NewInlineKeyboardButtonData(item.Name, item.Id)
		}
		markup := utils.NewTgMessageButtonsMarkup(buttons, 2)
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите репозиторий")
		message.BaseChat.ReplyMarkup = markup

		msg, err := bot.Send(message)
		if err != nil {
			logrus.Error(err)
		}
		form.lastMessage = msg
	}

	return false
}
func (s *SubscribeProcessor) updateRepositoryMessage(chatId int64, form *subscribeForm,
	bot *tgbotapi.BotAPI) {

	repos, err := s.service.GitlabApi().GetRepositories(form.gitlab)
	if err != nil {
		logrus.Error(err)
		return
	}

	buttons := make([]tgbotapi.InlineKeyboardButton, len(repos))
	for i, item := range repos {
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(item.Name, item.Id)
		if item.Id == form.repositoryId {
			buttons[i].Text += internal.GetEmoji(internal.WhiteCheckMark)
		}
	}

	markup := utils.NewTgMessageButtonsMarkup(buttons, 2)

	message := tgbotapi.NewEditMessageReplyMarkup(chatId, form.lastMessage.MessageID, markup)

	_, err = bot.Send(message)
	if err != nil {
		logrus.Error(err)
	}
}

func (s *SubscribeProcessor) suggestHookType(update tgbotapi.Update, form *subscribeForm,
	bot *tgbotapi.BotAPI) (isEnd bool) {
	buttons := make([]tgbotapi.InlineKeyboardButton, len(eventsNames))
	for i, item := range eventsNames {
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(item, strconv.Itoa(i))
	}
	markup := utils.NewTgMessageButtonsMarkup(buttons, 2)
	message := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите тип подписки")
	message.BaseChat.ReplyMarkup = markup

	msg, err := bot.Send(message)
	if err != nil {
		logrus.Error(err)
	}
	form.lastMessage = msg
	return false
}

func (s *SubscribeProcessor) updateHookTypeMessage(chatId int64, form *subscribeForm,
	bot *tgbotapi.BotAPI) bool {

	buttons := make([]tgbotapi.InlineKeyboardButton, len(eventsNames))
	for i, item := range eventsNames {
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(item, strconv.Itoa(i))
		if selected, _ := form.hookTypes[HookType(i)]; selected {
			buttons[i].Text += internal.GetEmoji(internal.WhiteCheckMark)
		}
	}

	markup := utils.NewTgMessageButtonsMarkup(buttons, 2)

	message := tgbotapi.NewEditMessageReplyMarkup(chatId, form.lastMessage.MessageID, markup)

	_, err := bot.Send(message)
	if err != nil {
		logrus.Error(err)
	}
	return false
}
