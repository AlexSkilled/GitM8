package processors

import (
	"context"
	"fmt"
	"strings"

	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/message-handling/langs"
	"gitlab-tg-bot/internal/message-handling/register"
	"gitlab-tg-bot/service/model"

	tg "github.com/AlexSkilled/go_tg/pkg"
	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
	"github.com/sirupsen/logrus"
)

const (
	StepDomain interfaces.StepName = iota
	StepToken
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

func (r *Register) Handle(ctx context.Context, message *tgmodel.MessageIn) (out tg.TgMessage) {
	registration, ok := r.dialogContext[message.From.ID]
	if !ok {
		r.dialogContext[message.From.ID] = &registrationProcess{
			CurrentStep: StepDomain,
		}

		gits := &tgmodel.InlineKeyboard{Columns: 1}
		gits.AddButton(string(model.Gitlab), string(model.Gitlab)+".com")

		return &tgmodel.MessageOut{
			Text:          langs.Get(ctx, register.AskDomain),
			InlineButtons: gits,
		}
	}

	switch registration.CurrentStep {
	case StepDomain:
		registration.Domain = message.Text
		registration.CurrentStep++

		return &tgmodel.MessageOut{
			Text: langs.Get(ctx, register.AskToken),
		}
	case StepToken:
		registration.GitlabToken = message.Text
	}

	err := r.services.User().AddGitAccount(message.From.ID, registration.ToDto())
	var response string
	if err != nil {
		response = langs.Get(ctx, register.Error)
		if strings.Contains(err.Error(), "<401>") {
			response = fmt.Sprintf(response, langs.Get(ctx, register.ErrorInvalidToken))
		} else {
			response = fmt.Sprintf(response, fmt.Sprintf(langs.Get(ctx, register.ErrorUnknown), err.Error()))
			logrus.Errorln(err)
		}
	} else {
		response = langs.Get(ctx, register.Success)
	}

	delete(r.dialogContext, message.From.ID)
	return &tgmodel.MessageOut{
		Text: response,
	}
}

func (r *Register) Dump(userId int64) {
	delete(r.dialogContext, userId)
}
