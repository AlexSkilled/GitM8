package processors

import (
	"context"
	"fmt"
	"strconv"

	"gitlab-tg-bot/internal/emoji"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"
	"gitlab-tg-bot/utils"
	"gitlab-tg-bot/worker/commands"

	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

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

func (s *SubscribeProcessor) Handle(ctx context.Context, in *tgmodel.MessageIn) (out *tgmodel.MessageOut) {
	user, err := utils.ExtractUser(ctx)
	if err != nil {
		logrus.Errorln(err)
		return nil
	}

	chatId := in.Chat.ID

	form, ok := s.subscribeForms[chatId]
	if !ok {
		s.subscribeForms[chatId] = &subscribeForm{
			currentStep: SubscriptionStep_OfferDomain,
			hookTypes:   map[HookType]bool{},
			tgUserId:    in.From.ID,
		}
		form, _ = s.subscribeForms[chatId]
	}

	switch form.currentStep {
	case SubscriptionStep_OfferDomain:
		form.currentStep++
		return s.suggestDomains(user, form, in, ctx)
	case SubscriptionStep_ObtainDomain:
		form.currentStep++
		form.domain = in.Text
		return s.Handle(ctx, in)
	case SubscriptionStep_OfferRepository:
		form.currentStep++
		for _, item := range user.Gitlabs {
			if item.Domain == form.domain {
				form.gitlab = item
				break
			}
		}
		return s.suggestRepositories(form, in, ctx)
	case SubscriptionStep_ObtainRepository:
		form.currentStep++
		form.repositoryId = in.Text
		s.updateRepositoryMessage(form)
		return s.Handle(ctx, in)
	case SubscriptionStep_OfferType:
		form.currentStep++
		return s.suggestHookType()
	case SubscriptionStep_ObtainType:
		typeNumber, err := strconv.ParseInt(in.Text, 10, 64)
		if err != nil {
			logrus.Errorf("Ошибка парсинга ответа! %v, в ответе %v", err, in)
		}

		if typeNumber != int64(EndChoosingEvent) {
			form.hookTypes[HookType(typeNumber)] = !form.hookTypes[HookType(typeNumber)]
			return s.updateHookTypeMessage()
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
	s.finalizeChoose(form)
	_, err = s.service.Subscription().Subscribe(form.gitlab, in.Chat.ID, webhook)

	delete(s.subscribeForms, chatId)

	if err != nil {
		return &tgmodel.MessageOut{
			Text: "Не удалось создать слушатель эвентов" + err.Error(),
		}
	}
	return &tgmodel.MessageOut{
		Text: "Успешно создан слушатель эвентов",
	}

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

func (s *SubscribeProcessor) suggestDomains(user model.User, form *subscribeForm, message *tgmodel.MessageIn, ctx context.Context) *tgmodel.MessageOut {
	switch len(user.Gitlabs) {
	case 0:
		return &tgmodel.MessageOut{Text: "Необходимо зарегестрировать аккаунт gitlab - команда " + commands.Register}
	case 1:
		form.domain = user.Gitlabs[0].Domain
		form.currentStep++

		out := s.Handle(ctx, message)
		out.Text = fmt.Sprintf("Зарегистрирован только один домен - %s. Этап пропускается.\n%s", user.Gitlabs[0].Domain, out.Text)
		return out
	default:
		btns := &tgmodel.InlineKeyboard{Columns: 2}
		for _, item := range user.Gitlabs {
			btns.AddButton("Аккаунт: "+item.Username+". Домен: "+item.Domain, item.Domain)
		}

		return &tgmodel.MessageOut{
			Text:          "Выберите домен",
			InlineButtons: btns,
		}
	}
}

//func (s *SubscribeProcessor) updateDomainMessage(user model.User, chatId int64, form *subscribeForm) {
//
//	buttons := make([]tgbotapi.InlineKeyboardButton, len(user.Gitlabs))
//	for i, item := range user.Gitlabs {
//		buttons[i] = tgbotapi.InlineKeyboardButton{Text: item.Domain}
//		buttons[i].CallbackData = nil
//
//		if item.Domain == form.domain {
//			buttons[i].Text += s.whiteCheck
//		}
//	}
//
//	markup := utils.NewTgMessageButtonsMarkup(buttons, 2)
//
//	message := tgbotapi.NewEditMessageReplyMarkup(chatId, form.lastMessage.MessageID, markup)
//
//	//_, err := bot.Send(message)
//	//if err != nil {
//	//	logrus.Error(err)
//	//}
//}

func (s *SubscribeProcessor) suggestRepositories(form *subscribeForm, message *tgmodel.MessageIn, ctx context.Context) (out *tgmodel.MessageOut) {
	repos, err := s.service.Subscription().GetRepositories(form.gitlab)
	if err != nil {
		logrus.Errorln(err)
	}
	switch len(repos) {
	case 0:
		return &tgmodel.MessageOut{
			Text: fmt.Sprintf("Для домена %s отсутствуют репозитории", form.gitlab.Domain),
		}
	case 1:
		form.repositoryId = repos[0].Id
		form.currentStep++

		out = s.Handle(ctx, message)
		out.Text = fmt.Sprintf("Найден единственный репозиторий - %s.\n%s", repos[0].Name, out.Text)
		return out
	default:
		btns := &tgmodel.InlineKeyboard{Columns: 2}
		for _, item := range repos {
			btns.AddButton(item.Name, item.Id)
		}

		return &tgmodel.MessageOut{
			Text:          "Выберите репозиторий",
			InlineButtons: btns,
		}
	}
}

func (s *SubscribeProcessor) updateRepositoryMessage(form *subscribeForm) (out *tgmodel.MessageOut) {
	repos, err := s.service.Subscription().GetRepositories(form.gitlab)
	if err != nil {
		logrus.Error(err)
		return
	}

	buttons := &tgmodel.InlineKeyboard{}
	for _, item := range repos {
		text := item.Name
		if item.Id == form.repositoryId {
			text += s.whiteCheck
		}
		buttons.AddButton(text, item.Id)
	}

	return &tgmodel.MessageOut{
		InlineButtons: buttons,
	}
}

func (s *SubscribeProcessor) suggestHookType() (out *tgmodel.MessageOut) {
	btns := &tgmodel.InlineKeyboard{Columns: 2}
	for i, item := range eventsNames {
		btns.AddButton(item, strconv.Itoa(i))
	}

	return &tgmodel.MessageOut{
		Text:          "Выберите тип подписки",
		InlineButtons: btns,
	}
}

func (s *SubscribeProcessor) updateHookTypeMessage() *tgmodel.MessageOut {
	btns := &tgmodel.InlineKeyboard{Columns: 2}
	for i, item := range eventsNames {
		btns.AddButton(item, strconv.Itoa(i))
	}

	return &tgmodel.MessageOut{
		Text:          "Выберите тип подписки",
		InlineButtons: btns,
	}
}

func (s *SubscribeProcessor) finalizeChoose(form *subscribeForm) *tgmodel.MessageOut {
	btns := &tgmodel.InlineKeyboard{Columns: 2}
	for i, item := range eventsNames {

		whiteCheck := s.whiteCheck

		if selected, _ := form.hookTypes[HookType(i)]; selected {
			btns.AddButton(item+whiteCheck, strconv.Itoa(i))
		}
	}

	return &tgmodel.MessageOut{
		Text:          "",
		InlineButtons: btns,
		Keyboard:      nil,
	}
}
