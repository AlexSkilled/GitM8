package processors

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-tg-bot/internal/interfaces"
)

const (
	StepUsername interfaces.StepName = iota
	StepToken
	StepEnd
)

type Register struct {
	services      interfaces.ServiceStorage
	dialogContext map[int64]*registrationProcess // [tgUserId]->[fieldName]->value
}

type registrationProcess struct {
	GitlabName  string
	GitlabToken string
	CurrentStep interfaces.StepName
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
	case StepUsername:
		registration.GitlabName = messageText
	case StepToken:
		registration.GitlabToken = messageText
	}
	registration.CurrentStep++

	if registration.CurrentStep == StepEnd {
		return true
	}
	return false
}

func (r *Register) IsInterceptor() bool {
	return true
}
