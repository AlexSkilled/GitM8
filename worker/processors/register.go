package processors

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/model"
)

const (
	StepUsername interfaces.StepName = iota
	StepToken
	StepDomain
	StepEnd
)

type Register struct {
	services      interfaces.ServiceStorage
	dialogContext map[int64]*registrationProcess // [tgUserId]->[fieldName]->value
}

type registrationProcess struct {
	GitlabName  string
	GitlabToken string
	Domain      string
	CurrentStep interfaces.StepName
}

func (r *registrationProcess) ToDto() model.GitlabUser {
	return model.GitlabUser{
		Username: r.GitlabName,
		Token:    r.GitlabToken,
		Domain:   r.Domain,
	}
}

func NewRegisterProcessor(services interfaces.ServiceStorage) *Register {
	return &Register{
		services:      services,
		dialogContext: map[int64]*registrationProcess{},
	}
}

func (r *Register) Process(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI) (isEnd bool) {
	registration, ok := r.dialogContext[update.Message.From.ID]
	if !ok {
		r.dialogContext[update.Message.From.ID] = &registrationProcess{
			CurrentStep: StepUsername,
		}
		return false
	}
	messageText := update.Message.Text
	switch registration.CurrentStep {
	// TODO на каждом этапе нужно дать подсказки, что вводить дальше
	case StepUsername:
		registration.GitlabName = messageText
	case StepToken:
		// TODO Пока что нужен токен только со скопом на api
		registration.GitlabToken = messageText
	case StepDomain:
		registration.Domain = messageText
	}
	logrus.Info("Для шага ", registration.CurrentStep, ". Используется значение:", messageText)
	registration.CurrentStep++

	if registration.CurrentStep >= StepEnd {
		err := r.services.User().AddGitlabAccount(update.Message.From.ID, registration.ToDto())
		if err != nil {
			// TODO обработать
		}
		delete(r.dialogContext, update.Message.From.ID)
		return true
	}
	return false
}

func (r *Register) IsInterceptor() bool {
	return true
}
