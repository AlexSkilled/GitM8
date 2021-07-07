package processors

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-tg-bot/internal/interfaces"
)

const (
	StepUsername interfaces.StepName = iota
	StepToken
)

type Register struct {
	services      interfaces.ServiceStorage
	dialogContext map[int64]registrationProcess // [tgUserId]->[fieldName]->value
}

type registrationProcess struct {
	GitlabName  string
	GitlabToken string
	CurrentStep interfaces.StepName
}

func NewRegisterProcessor(services interfaces.ServiceStorage) *Register {
	return &Register{
		services:      services,
		dialogContext: map[int64]registrationProcess{},
	}
}

func (r *Register) Process(update tgbotapi.Update, bot *tgbotapi.BotAPI) (isEnd bool) {
	return true
}

func (u *Register) IsInterceptor() bool {
	return true
}
