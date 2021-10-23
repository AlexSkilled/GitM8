package processors

import (
	"context"
	"strings"

	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"

	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
	"github.com/sirupsen/logrus"
)

const (
	StepToken interfaces.StepName = iota
	StepDomain
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

func NewRegisterProcessor(services interfaces.ServiceStorage) *Register {
	return &Register{
		services:      services,
		dialogContext: map[int64]*registrationProcess{},
	}
}

func (r *Register) Handle(_ context.Context, message *tgmodel.MessageIn) (out *tgmodel.MessageOut) {
	registration, ok := r.dialogContext[message.From.ID]
	if !ok {
		r.dialogContext[message.From.ID] = &registrationProcess{
			CurrentStep: StepToken,
		}
		return &tgmodel.MessageOut{
			Text: "Введите токен для  Gitlab (необходимы права на использование API)",
		}
	}

	switch registration.CurrentStep {
	case StepToken:
		registration.GitlabToken = message.Text
		registration.CurrentStep++
		return &tgmodel.MessageOut{
			Text: "Введите домен Gitlab (стандартный gitlab.com)",
		}
	case StepDomain:
		registration.Domain = message.Text
	}

	err := r.services.User().AddGitAccount(message.From.ID, registration.ToDto())
	response := "Успешная регистрация!"
	if err != nil {
		response = "Ошибка при регистрации!"
		if strings.Contains(err.Error(), "<401>") {
			response += "Авторизация не прошла. Скорее всего токен не годный"
		} else {
			response += "Неизвестная ошибка. Повторите попытку позже"
			logrus.Errorln(err)
		}
	}
	delete(r.dialogContext, message.From.ID)
	return &tgmodel.MessageOut{
		Text: response,
	}
}

func (r *Register) IsInterceptor() bool {
	return true
}

func (r *Register) DumpChatSession(userId int64) {
	delete(r.dialogContext, userId)
}
