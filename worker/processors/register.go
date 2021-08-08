package processors

import (
	"context"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

const (
	StepToken interfaces.StepName = iota
	StepDomain
	StepEnd
)

type Register struct {
	services      interfaces.ServiceStorage
	dialogContext map[int64]*registrationProcess // [tgUserId]->regForm
}

type registrationProcess struct {
	GitlabName  string
	GitlabToken string
	Domain      string
	CurrentStep interfaces.StepName
}

func (r *registrationProcess) ToDto() model.GitUser {
	return model.GitUser{
		Username: r.GitlabName,
		Token:    r.GitlabToken,
		Domain:   r.Domain,
	}
}

var _ interfaces.Interceptor = (*Register)(nil)

func NewRegisterProcessor(services interfaces.ServiceStorage) *Register {
	return &Register{
		services:      services,
		dialogContext: map[int64]*registrationProcess{},
	}
}

func (r *Register) Process(ctx context.Context, message *tgbotapi.Message, bot *tgbotapi.BotAPI) (isEnd bool) {
	registration, ok := r.dialogContext[message.From.ID]
	if !ok {
		r.dialogContext[message.From.ID] = &registrationProcess{
			CurrentStep: StepToken,
		}
		_, _ = bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Введите токен для  Gitlab (необходимы права на использование API)"))
		return false
	}
	messageText := message.Text
	switch registration.CurrentStep {
	case StepToken:
		registration.GitlabToken = messageText
		_, _ = bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Введите домен Gitlab (стандартный gitlab.ru)"))
	case StepDomain:
		registration.Domain = messageText
	}
	logrus.Info("Для шага ", registration.CurrentStep, ". Используется значение:", messageText)
	registration.CurrentStep++

	if registration.CurrentStep >= StepEnd {
		err := r.services.User().AddGitAccount(message.From.ID, registration.ToDto())
		if err != nil {
			// TODO обработать
		}
		_, _ = bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Успешная регистрация!"))
		delete(r.dialogContext, message.From.ID)
		return true
	}
	return false
}

func (r *Register) IsInterceptor() bool {
	return true
}

func (r *Register) DumpChatSession(userId int64) {
	delete(r.dialogContext, userId)
}
