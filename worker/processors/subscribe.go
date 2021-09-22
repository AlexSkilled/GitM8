package processors

import (
	"context"
	"fmt"
	"gitlab-tg-bot/internal/emoji"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"
	"gitlab-tg-bot/utils"
	"strconv"

	"github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var _ interfaces.Interceptor = (*SubscribeProcessor)(nil)

type SubscriptionStep int32

const (
	SubscriptionStep_OfferDomain SubscriptionStep = iota
	SubscriptionStep_ObtainDomain
	SubscriptionStep_OfferRepository
	SubscriptionStep_ObtainRepository
	SubscriptionStep_OfferType
	SubscriptionStep_ObtainType

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
	whiteCheck     string
}

type subscribeForm struct {
	domain       string
	repositoryId string
	tgUserId     int64
	gitlab       model.GitUser
	currentStep  SubscriptionStep
	hookTypes    map[HookType]bool
	lastMessage  tgbotapi.Message
}

func NewSubscribeProcessor(services interfaces.ServiceStorage) *SubscribeProcessor {
	return &SubscribeProcessor{
		service:        services,
		subscribeForms: map[int64]*subscribeForm{},
		whiteCheck:     emoji.GetEmoji(emoji.WhiteCheckMark),
	}
}

func (s *SubscribeProcessor) Process(ctx context.Context, message *tgbotapi.Message, bot *tgbotapi.BotAPI) (isEnd bool) {
	user, err := utils.ExtractUser(ctx)
	if err != nil {
		logrus.Errorln(err)
		return true
	}

	chatId := message.Chat.ID

	form, ok := s.subscribeForms[chatId]
	if !ok {
		s.subscribeForms[chatId] = &subscribeForm{
			currentStep: SubscriptionStep_OfferDomain,
			hookTypes:   map[HookType]bool{},
			tgUserId:    message.From.ID,
		}
		form, _ = s.subscribeForms[chatId]
	}

	switch form.currentStep {
	case SubscriptionStep_OfferDomain:
		form.currentStep++
		return s.suggestDomains(user, form, message, bot, ctx)
	case SubscriptionStep_ObtainDomain:
		form.currentStep++
		form.domain = message.Text
		return s.Process(ctx, message, bot)
	case SubscriptionStep_OfferRepository:
		form.currentStep++
		for _, item := range user.Gitlabs {
			if item.Domain == form.domain {
				form.gitlab = item
				break
			}
		}
		return s.suggestRepositories(form, message, bot, ctx)
	case SubscriptionStep_ObtainRepository:
		form.currentStep++
		form.repositoryId = message.Text
		s.updateRepositoryMessage(message.Chat.ID, form, bot)
		return s.Process(ctx, message, bot)
	case SubscriptionStep_OfferType:
		form.currentStep++
		return s.suggestHookType(message, form, bot)
	case SubscriptionStep_ObtainType:
		typeNumber, err := strconv.ParseInt(message.Text, 10, 64)
		if err != nil {
			logrus.Errorf("Ошибка парсинга ответа! %v, в ответе %v", err, message)
		}

		if typeNumber != int64(EndChoosingEvent) {
			form.hookTypes[HookType(typeNumber)] = !form.hookTypes[HookType(typeNumber)]
			return s.updateHookTypeMessage(message.Chat.ID, form, bot)
		}
	}
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
	s.finalizeChoose(form, bot)
	_, err = s.service.Subscription().Subscribe(form.gitlab, message.Chat.ID, webhook)

	delete(s.subscribeForms, chatId)

	if err != nil {
		logrus.Errorln(err)
		_, _ = bot.Send(tgbotapi.NewMessage(message.Chat.ID,
			"Не удалось создать слушатель эвентов"+err.Error()))
		return true
	}
	bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Успешно создан слушатель эвентов"))

	return true

	//if update.CallbackQuery != nil {
	//	switch form.currentStep {
	//	case SubscriptionStepDomain:
	//		form.domain = update.CallbackQuery.Data
	//		s.updateDomainMessage(user, update.CallbackQuery.Message.Chat.ID, form, bot)
	//	case SubscriptionStepRepository:
	//		form.repositoryId = update.CallbackQuery.Data
	//		s.updateRepositoryMessage(update.CallbackQuery.Message.Chat.ID, form, bot)
	//	case SubscriptionStepType:
	//		typeNumber, err := strconv.ParseInt(update.CallbackQuery.Data, 10, 64)
	//		if err != nil {
	//			logrus.Errorln("Ошибка парсинга ответа!", err)
	//		}
	//
	//		if typeNumber != int64(EndChoosingEvent) {
	//			form.hookTypes[HookType(typeNumber)] = !form.hookTypes[HookType(typeNumber)]
	//			return s.updateHookTypeMessage(update.CallbackQuery.Message.Chat.ID, form, bot)
	//		}
	//	}
	//
	//	form.currentStep++
	//}
	//
	//switch form.currentStep {
	//case SubscriptionStepDomain:
	//	return s.suggestDomains(user, form, update, bot, ctx)
	//case SubscriptionStepRepository:
	//	for _, item := range user.Gitlabs {
	//		if item.Domain == form.domain {
	//			form.gitlab = item
	//			break
	//		}
	//	}
	//	return s.suggestRepositories(form, update, bot, ctx)
	//case SubscriptionStepType:
	//	return s.suggestHookType(update, form, bot)
	//case SubscriptionStepEnd:
	//	webhook := model.Hook{
	//		RepoId:                   form.repositoryId,
	//		PushEvents:               form.hookTypes[PushEvents],
	//		IssuesEvents:             form.hookTypes[IssuesEvents],
	//		ConfidentialIssuesEvents: form.hookTypes[ConfidentialIssuesEvents],
	//		MergeRequestsEvents:      form.hookTypes[MergeRequestsEvents],
	//		TagPushEvents:            form.hookTypes[TagPushEvents],
	//		NoteEvents:               form.hookTypes[NoteEvents],
	//		JobEvents:                form.hookTypes[JobEvents],
	//		PipelineEvents:           form.hookTypes[PipelineEvents],
	//		WikiPageEvents:           form.hookTypes[WikiPageEvents],
	//	}
	//	s.finalizeChoose(form, bot)
	//	_, err = s.service.Subscription().Subscribe(form.gitlab, update.CallbackQuery.Message.Chat.ID, webhook)
	//
	//	delete(s.subscribeForms, chatId)
	//
	//	if err != nil {
	//		logrus.Errorln(err)
	//		_, _ = bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
	//			"Не удалось создать слушатель эвентов"+err.Error()))
	//		return true
	//	}
	//	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Успешно создан слушатель эвентов"))
	//}

	//return true
}

func (s *SubscribeProcessor) IsInterceptor() bool {
	return true
}

func (s *SubscribeProcessor) DumpChatSession(chatId int64) {
	delete(s.subscribeForms, chatId)
}

func (s *SubscribeProcessor) suggestDomains(user model.User, form *subscribeForm, message *tgbotapi.Message,
	bot *tgbotapi.BotAPI, ctx context.Context) bool {
	switch len(user.Gitlabs) {
	case 0:
		_, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Необходимо зарегестрировать аккаунт gitlab - команда /register"))
		if err != nil {
			logrus.Errorln(err)
		}
		return true
	case 1:
		form.domain = user.Gitlabs[0].Domain
		form.currentStep++
		_, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID,
			fmt.Sprintf("Зарегистрирован только один домен - %s. Этап пропускается", user.Gitlabs[0].Domain)))
		if err != nil {
			logrus.Errorln(err)
			return true
		}
		return s.Process(ctx, message, bot)
	default:
		buttons := make([]tgbotapi.InlineKeyboardButton, len(user.Gitlabs))
		for _, item := range user.Gitlabs {
			buttons = append(buttons,
				tgbotapi.NewInlineKeyboardButtonData("Аккаунт: "+item.Username+". Домен: "+item.Domain, item.Domain))
		}

		markup := utils.NewTgMessageButtonsMarkup(buttons, 2)
		messageOut := tgbotapi.NewMessage(message.Chat.ID, "Выберите домен")
		messageOut.BaseChat.ReplyMarkup = markup

		msg, err := bot.Send(messageOut)
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
			buttons[i].Text += s.whiteCheck
		}
	}

	markup := utils.NewTgMessageButtonsMarkup(buttons, 2)

	message := tgbotapi.NewEditMessageReplyMarkup(chatId, form.lastMessage.MessageID, markup)

	_, err := bot.Send(message)
	if err != nil {
		logrus.Error(err)
	}
}

func (s *SubscribeProcessor) suggestRepositories(form *subscribeForm, message *tgbotapi.Message,
	bot *tgbotapi.BotAPI, ctx context.Context) (isEnd bool) {
	repos, err := s.service.Subscription().GetRepositories(form.gitlab)
	if err != nil {
		logrus.Errorln(err)
	}
	switch len(repos) {
	case 0:
		_, err = bot.Send(tgbotapi.NewMessage(message.Chat.ID,
			fmt.Sprintf("Для домена %s отсутствуют репозитории", form.gitlab.Domain)))
		return true
	case 1:
		form.repositoryId = repos[0].Id
		form.currentStep++
		bot.Send(tgbotapi.NewMessage(message.Chat.ID,
			fmt.Sprintf("Найден единственный репозиторий - %s.", repos[0].Name)))
		return s.Process(ctx, message, bot)
	default:
		buttons := make([]tgbotapi.InlineKeyboardButton, len(repos))
		for i, item := range repos {
			buttons[i] = tgbotapi.NewInlineKeyboardButtonData(item.Name, item.Id)
		}
		markup := utils.NewTgMessageButtonsMarkup(buttons, 2)
		messageOut := tgbotapi.NewMessage(message.Chat.ID, "Выберите репозиторий")
		messageOut.BaseChat.ReplyMarkup = markup

		msg, err := bot.Send(messageOut)
		if err != nil {
			logrus.Error(err)
		}
		form.lastMessage = msg
	}

	return false
}
func (s *SubscribeProcessor) updateRepositoryMessage(chatId int64, form *subscribeForm,
	bot *tgbotapi.BotAPI) {

	repos, err := s.service.Subscription().GetRepositories(form.gitlab)
	if err != nil {
		logrus.Error(err)
		return
	}

	buttons := make([]tgbotapi.InlineKeyboardButton, len(repos))
	for i, item := range repos {
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(item.Name, item.Id)
		if item.Id == form.repositoryId {
			buttons[i].Text += s.whiteCheck
		}
	}

	markup := utils.NewTgMessageButtonsMarkup(buttons, 2)

	message := tgbotapi.NewEditMessageReplyMarkup(chatId, form.lastMessage.MessageID, markup)

	_, err = bot.Send(message)
	if err != nil {
		logrus.Error(err)
	}
}

func (s *SubscribeProcessor) suggestHookType(message *tgbotapi.Message, form *subscribeForm,
	bot *tgbotapi.BotAPI) (isEnd bool) {
	buttons := make([]tgbotapi.InlineKeyboardButton, len(eventsNames))
	for i, item := range eventsNames {
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(item, strconv.Itoa(i))
	}

	markup := utils.NewTgMessageButtonsMarkup(buttons, 2)
	messageOut := tgbotapi.NewMessage(message.Chat.ID, "Выберите тип подписки")
	messageOut.BaseChat.ReplyMarkup = markup

	msg, err := bot.Send(messageOut)
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
			buttons[i].Text += s.whiteCheck
		}
	}

	markup := utils.NewTgMessageButtonsMarkup(buttons, 2)

	message := tgbotapi.NewEditMessageReplyMarkup(chatId, form.lastMessage.MessageID, markup)

	msg, err := bot.Send(message)
	if err != nil {
		logrus.Error(err)
	}

	form.lastMessage = msg
	return false
}

func (s *SubscribeProcessor) finalizeChoose(form *subscribeForm,
	bot *tgbotapi.BotAPI) {
	buttons := make([]tgbotapi.InlineKeyboardButton, 0, 1)
	for i, item := range eventsNames {

		whiteCheck := s.whiteCheck

		if selected, _ := form.hookTypes[HookType(i)]; selected {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(item+whiteCheck, strconv.Itoa(i)))
		}
	}

	markup := utils.NewTgMessageButtonsMarkup(buttons, 2)

	message := tgbotapi.NewEditMessageReplyMarkup(form.lastMessage.Chat.ID, form.lastMessage.MessageID, markup)

	_, err := bot.Send(message)
	if err != nil {
		logrus.Error(err)
	}
}
