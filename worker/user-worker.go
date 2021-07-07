package worker

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-tg-bot/internal"
	"gitlab-tg-bot/internal/interfaces"
)

type Register struct {
	services      interfaces.ServiceStorage
	dialogContext map[int32]map[string]interface{} // [tgUserId]->[fieldName]->value
}

var registrationOrder = []string{
	"gitlabUsername",
	"gitlabToken",
}

type registrationForm struct {
	GitlabName  string
	GitlabToken string
}

func NewRegisterProcessor(services interfaces.ServiceStorage) internal.TgProcessor {
	return &Register{
		services:      services,
		dialogContext: map[int32]map[string]interface{}{},
	}
}

func (u *Register) Process(update tgbotapi.Update) (isEnd bool) {
	//dialogContext := u.dialogContext[int32(update.Message.From.ID)]
	return true
}

func (u *Register) IsInterceptor() bool {
	return true
}
