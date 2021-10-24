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

	tg "github.com/AlexSkilled/go_tg/pkg"
	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
	"github.com/sirupsen/logrus"
)

type SubscriptionStep int32

const (
	SubscriptionStep_OfferDomain SubscriptionStep = iota
	SubscriptionStep_ObtainDomain
	SubscriptionStep_ObtainRepository
	SubscriptionStep_ObtainType
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
}

func NewSubscribeProcessor(services interfaces.ServiceStorage) *SubscribeProcessor {
	return &SubscribeProcessor{
		service:        services,
		subscribeForms: map[int64]*subscribeForm{},
		whiteCheck:     emoji.GetEmoji(emoji.WhiteCheckMark),
	}
}

func (s *SubscribeProcessor) Handle(ctx context.Context, in *tgmodel.MessageIn) (out tg.TgMessage) {
	user, err := utils.ExtractUser(ctx)
	if err != nil {
		logrus.Errorln(err)
		return nil
	}

	form, ok := s.subscribeForms[in.Chat.ID]
	if !ok {
		s.subscribeForms[in.Chat.ID] = &subscribeForm{
			currentStep: SubscriptionStep_OfferDomain,
			hookTypes:   map[HookType]bool{},
			tgUserId:    in.From.ID,
		}
		form, _ = s.subscribeForms[in.Chat.ID]
	}

	switch form.currentStep {
	case SubscriptionStep_OfferDomain:
		form.currentStep++
		return s.suggestDomains(user, form, in.Chat.ID)
	case SubscriptionStep_ObtainDomain:
		form.currentStep++
		for _, item := range user.Gitlabs {
			if item.Domain == in.Text {
				form.gitlab = item
				break
			}
		}
		multipleMessage := &tg.MultipleMessage{
			in.Chat.ID: []tg.TgMessage{
				s.updateDomains(form, user, in.MessageID),
				s.suggestRepositories(form, in.Chat.ID),
			},
		}
		return multipleMessage
	case SubscriptionStep_ObtainRepository:
		form.currentStep++
		form.repositoryId = in.Text
		multipleMessage := &tg.MultipleMessage{
			in.Chat.ID: []tg.TgMessage{
				s.updateRepositoryMessage(form, in.MessageID),
				s.suggestHookType(),
			},
		}
		return multipleMessage
	case SubscriptionStep_ObtainType:
		typeNumber, err := strconv.ParseInt(in.Text, 10, 64)
		if err != nil {
			logrus.Errorf("Ошибка парсинга ответа! %v, в ответе %v", err, in)
		}

		if typeNumber != int64(EndChoosingEvent) {
			form.hookTypes[HookType(typeNumber)] = !form.hookTypes[HookType(typeNumber)]
			return s.updateHookTypeMessage(form, in.MessageID)
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

	delete(s.subscribeForms, in.Chat.ID)

	if err != nil {
		return &tgmodel.MessageOut{
			Text: "Не удалось создать слушатель эвентов" + err.Error(),
		}
	}
	return &tgmodel.MessageOut{
		Text: "Успешно создан слушатель эвентов",
	}
}

func (s *SubscribeProcessor) IsInterceptor() bool {
	return true
}

func (s *SubscribeProcessor) Dump(chatId int64) {
	delete(s.subscribeForms, chatId)
}

func (s *SubscribeProcessor) suggestDomains(user model.User, form *subscribeForm, chatId int64) tg.TgMessage {
	switch len(user.Gitlabs) {
	case 0:
		s.Dump(chatId)
		return &tgmodel.MessageOut{Text: "Необходимо зарегестрировать аккаунт gitlab - команда " + commands.Register}
	case 1:
		form.currentStep++
		form.gitlab = user.Gitlabs[0]
		out := &tg.MultipleMessage{
			chatId: []tg.TgMessage{
				&tgmodel.MessageOut{
					Text: fmt.Sprintf("Зарегистрирован только один домен - %s.\nЭтап пропускается.", user.Gitlabs[0].Domain),
				},
				s.suggestRepositories(form, chatId),
			},
		}
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

func (s *SubscribeProcessor) updateDomains(form *subscribeForm, user model.User, messageId int) (out tg.TgMessage) {
	buttons := &tgmodel.InlineKeyboard{}
	for _, item := range user.Gitlabs {
		text := item.Domain
		if item.Domain == form.domain {
			text += s.whiteCheck
		} else {
			buttons.AddButton(text, item.Domain)
		}
	}

	return tgmodel.EditMessageReply(buttons, messageId)
}

func (s *SubscribeProcessor) suggestRepositories(form *subscribeForm, chatId int64) (out tg.TgMessage) {
	repos, err := s.service.Subscription().GetRepositories(form.gitlab)
	if err != nil {
		logrus.Errorln(err)
	}
	switch len(repos) {
	case 0:
		s.Dump(chatId)
		return &tgmodel.MessageOut{
			Text: fmt.Sprintf("Для домена %s отсутствуют репозитории", form.gitlab.Domain),
		}
	case 1:
		form.currentStep++
		form.repositoryId = repos[0].Id

		out = &tg.MultipleMessage{
			chatId: []tg.TgMessage{
				&tgmodel.MessageOut{Text: fmt.Sprintf("Найден единственный репозиторий - %s.Этап пропускается", repos[0].Name)},
				s.suggestHookType(),
			},
		}
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

func (s *SubscribeProcessor) updateRepositoryMessage(form *subscribeForm, messageId int) (out tg.TgMessage) {
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
	return tgmodel.EditMessageReply(buttons, messageId)
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

func (s *SubscribeProcessor) updateHookTypeMessage(form *subscribeForm, messageId int) *tgmodel.MessageEdit {
	btns := &tgmodel.InlineKeyboard{Columns: 2}
	for i, item := range eventsNames {
		if selected, _ := form.hookTypes[HookType(i)]; selected {
			btns.AddButton(item+s.whiteCheck, strconv.Itoa(i))
		} else {
			btns.AddButton(item, strconv.Itoa(i))
		}
	}

	return tgmodel.EditMessageReply(btns, messageId)
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
